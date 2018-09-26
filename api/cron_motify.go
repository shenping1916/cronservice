package api

import (
	"fmt"
	"time"
	"encoding/json"
	"iceberg/frame"
	"cronservice/pb"
	"cronservice/models"
	"cronservice/handler"
	"cronservice/utility/base"
)

func (cs *CronService) MotifyTask(c frame.Context) error {
	var req pb.MotifyTaskReq
	if err := c.Bind(&req); err != nil {
		return err
	}

	if req.GetTaskName() == "" {
		return c.JSON2(ERR_CODE_BAD_REQUEST, "缺少定时任务名：task_name字段！", nil)
	}

	cronservice := new(models.CronService)
	if cs.CronTask.DbRead.Where("task_name = ?", req.GetTaskName()).Find(&cronservice).RecordNotFound() {
		msg := fmt.Errorf("定时任务：%s 不存在！", req.GetTaskName())
		return c.JSON2(ERR_CODE_BAD_NOT_FOUND, msg.Error(), nil)
	}

	var m = make(map[string]interface{})
	if req.GetTimeFormat() != "" {
		m["time_format"] = req.GetTimeFormat()
	}
	if req.GetServiceUrl() != "" {
		m["service_url"] = req.GetServiceUrl()
	}
	if req.GetServiceMethod() != "" {
		m["service_method"] = req.GetServiceMethod()
	}
	if req.GetMethod() != "" {
		m["method"] = req.GetMethod()
	}
	if req.GetServiceHeader() != "" {
		m["service_header"] = req.GetServiceHeader()
	}
	if req.GetForm() != "" {
		m["form"] = req.GetForm()
	}
	if req.GetServiceBody() != "" {
		m["service_body"] = req.GetServiceBody()
	}
	if req.GetTaskDesc() != "" {
		m["task_desc"] = req.GetTaskDesc()
	}

	if err := cs.CronTask.DbWrite.Model(&cronservice).Updates(m).Error; err != nil {
		msg := fmt.Errorf("数据库更新失败！原因：%v", err.Error())
		return c.JSON2(ERR_CODE_INTERNAL_ERROR, msg.Error(), nil)
	}

	var time_format string
	if value, ok := m["time_format"]; ok {
		time_format = value.(string)
	} else {
		time_format = cronservice.TimeFormat
	}

	var service_url string
	if value, ok := m["service_url"]; ok {
		service_url = value.(string)
	} else {
		service_url = cronservice.ServiceUrl
	}

	var service_method string
	if value, ok := m["service_method"]; ok {
		service_method = value.(string)
	} else {
		service_method = cronservice.ServiceMethod
	}

	var method string
	if value, ok := m["method"]; ok {
		method = value.(string)
	} else {
		method = cronservice.Method
	}

	var header map[string]string
	if _, ok := m["service_header"]; !ok {
		if cronservice.ServiceHeader == "" {
			header = make(map[string]string)
		} else {
			if err := json.Unmarshal([]byte(cronservice.ServiceHeader), &header); err != nil {
				return c.JSON2(ERR_CODE_INTERNAL_ERROR, err.Error(), nil)
			}
		}
	} else {
		if err := json.Unmarshal([]byte(req.GetServiceHeader()), &header); err != nil {
			return c.JSON2(ERR_CODE_INTERNAL_ERROR, err.Error(), nil)
		}
	}

	var form map[string]string
	if _, ok := m["form"]; !ok {
		if cronservice.Form == "" {
			form = make(map[string]string)
		} else {
			if err := json.Unmarshal([]byte(cronservice.Form), &form); err != nil {
				return c.JSON2(ERR_CODE_INTERNAL_ERROR, err.Error(), nil)
			}
		}
	} else {
		if err := json.Unmarshal([]byte(req.GetForm()), &form); err != nil {
			return c.JSON2(ERR_CODE_INTERNAL_ERROR, err.Error(), nil)
		}
	}

	var body []byte
	if _, ok := m["service_body"]; !ok {
		if cronservice.ServiceBody == "" {
			body = make([]byte, 0)
		} else {
			body = base.StringToBytes(cronservice.ServiceBody)
		}
	} else {
		body = base.StringToBytes(req.GetServiceBody())
	}

	// 加redis锁
	redis_client := GlobalCronService.CronTask.RedisClient
	mm := make(map[string]interface{})
	mm["lock_key"] = base.RedisLockKey(req.GetTaskName())
	mm["redis_client"] = redis_client
	f := handler.Convert(req.GetTaskName(), service_url, service_method, method, header, form, body, mm)

	// 发布修改定时任务消息到redis，通知其它节点
	redis_msg := &handler.RedisMsg{
		cs.CronTask.UUID,
		handler.REDIS_MOTIFT,
		cronservice,
	}

	b, err := json.Marshal(redis_msg)
	if err != nil {
		return c.JSON2(ERR_CODE_INTERNAL_ERROR, err.Error(), nil)
	}
	if err := redis_client.Publish(base.RedisPubSubChannel(), base.BytesToString(b)).Err(); err != nil {
		msg := fmt.Errorf("消息发布到redis失败！原因：%v", err.Error())
		return c.JSON2(ERR_CODE_INTERNAL_ERROR, msg.Error(), nil)
	}

	// 删除原定时任务
	handler.StopSpeciTask(req.GetTaskName())

	// sleep 1s后将修改后的定时任务加入
	time.Sleep(1 * time.Second)

	// 加入新的定时任务
	handler.AddTask(req.GetTaskName(), time_format, f)

	msg := fmt.Sprintf("定时任务：%s 修改成功", req.GetTaskName())
	return c.JSON2(Status_OK, msg, nil)
}


