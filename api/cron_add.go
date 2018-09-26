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

func (cs *CronService) RegisterTask(c frame.Context) error {
	var req pb.RegisterReq
	if err := c.Bind(&req); err != nil {
		return err
	}

	if req.GetTaskName() == "" {
		return c.JSON2(ERR_CODE_BAD_REQUEST, "缺少定时任务名：task_name字段！", nil)
	}
	
	if req.GetTimeFormat() == "" {
		return c.JSON2(ERR_CODE_BAD_REQUEST, "缺少时间配置：time_format字段！", nil)
	}
	
	if req.GetServiceUrl() == "" {
		return c.JSON2(ERR_CODE_BAD_REQUEST, "缺少路由：service_url字段！", nil)
	}

	if req.GetServiceMethod() == "" {
		return c.JSON2(ERR_CODE_BAD_REQUEST, "缺少被调用方方法：service_method字段！", nil)
	}

	if req.GetMethod() == "" {
		return c.JSON2(ERR_CODE_BAD_REQUEST, "缺少http方法：method字段！", nil)
	}

	//f := handler.Convert(req.GetTaskName(), "/services/v1/gift", "giftsinfo", "POST", map[string]string{"Content-Type": "application/json; charset=UTF-8"}, make(map[string]string), []byte(`{"g_id": [50, 51, 52]}`))
	var header map[string]string
	if req.GetServiceHeader() == "" {
		header = make(map[string]string)
	} else {
		if err := json.Unmarshal([]byte(req.GetServiceHeader()), &header); err != nil {
			return c.JSON2(ERR_CODE_INTERNAL_ERROR, err.Error(), nil)
		}
	}

	var form map[string]string
	if req.GetForm() == "" {
		form = make(map[string]string)
	} else {
		if err := json.Unmarshal([]byte(req.GetForm()), &form); err != nil {
			return c.JSON2(ERR_CODE_INTERNAL_ERROR, err.Error(), nil)
		}
	}

	var body []byte
	if req.GetServiceBody() == "" {
		body = make([]byte, 0)
	} else {
		body = base.StringToBytes(req.GetServiceBody())
	}

	// 加redis锁
	redis_client := GlobalCronService.CronTask.RedisClient
	m := make(map[string]interface{})
	m["lock_key"] = base.RedisLockKey(req.GetTaskName())
	m["redis_client"] = redis_client
	f := handler.Convert(req.GetTaskName(), req.GetServiceUrl(), req.GetServiceMethod(), req.GetMethod(), header,
		form, body, m)
	cronservice := new(models.CronService)
	if cs.CronTask.DbRead.Where("task_name = ? AND time_format = ? AND service_url = ?",
		req.GetTaskName(), req.GetTimeFormat(), req.GetServiceUrl()).Find(&cronservice).RecordNotFound() {
		cronservice.TaskName = req.GetTaskName()
		cronservice.TimeFormat = req.GetTimeFormat()
		cronservice.ServiceUrl = req.GetServiceUrl()
		cronservice.ServiceMethod = req.GetServiceMethod()
		cronservice.ServiceHeader = req.GetServiceHeader()
		cronservice.Form = req.GetForm()
		cronservice.Method = req.GetMethod()
		cronservice.ServiceBody = req.GetServiceBody()
		cronservice.Status = 1
		cronservice.LockKey = base.RedisLockKey(req.GetTaskName())
		if req.GetTaskDesc() != "" {
			cronservice.TaskDesc = req.GetTaskDesc()
		}
		cronservice.CreatedAt = time.Now().Local()
		if err := cs.CronTask.DbWrite.Create(&cronservice).Error; err != nil {
			msg := fmt.Errorf("定时任务: %v 写入DB失败！", &cronservice)
			return c.JSON2(ERR_CODE_INTERNAL_ERROR, msg.Error(), nil)
		}

		// 发布注册定时任务消息到redis，通知其它节点
		redis_msg := &handler.RedisMsg{
			cs.CronTask.UUID,
			handler.REDIS_ADD,
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

		// 加入全局定时任务
		handler.AddTask(req.GetTaskName(), req.GetTimeFormat(), f)
		msg := fmt.Sprintf("定时任务：%s 注册成功", req.GetTaskName())
		return c.JSON2(Status_OK, msg, nil)
	} else {
		msg := fmt.Errorf("定时任务：%s 已存在，请勿重复注册！", req.GetTaskName())
		return c.JSON2(ERR_CODE_FORBIDDEN, msg.Error(), nil)
	}
}
