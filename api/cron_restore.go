package api

import (
	"fmt"
	"iceberg/frame"
	"cronservice/pb"
	"cronservice/models"
	"cronservice/handler"
	"encoding/json"
	"cronservice/utility/base"
)

func (cs *CronService) RestoreTask(c frame.Context) error {
	var req pb.RestoreTaskReq
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

	if err := cs.CronTask.DbWrite.Model(&cronservice).Where("task_name = ?", req.GetTaskName()).Error; err == nil {
		if cronservice.Status == 1 {
			msg := fmt.Errorf("定时任务：%s 运行中，无需恢复！", req.GetTaskName())
			return c.JSON2(ERR_CODE_FORBIDDEN  , msg.Error(), nil)
		}
	}

	if err := cs.CronTask.DbWrite.Model(&cronservice).Where("task_name = ?", req.GetTaskName()).Update("status", 1).Error; err != nil {
		msg := fmt.Errorf("数据库更新失败！原因：%v", err.Error())
		return c.JSON2(ERR_CODE_INTERNAL_ERROR, msg.Error(), nil)
	}

	var header map[string]string
	if cronservice.ServiceHeader == "" {
		header = make(map[string]string)
	} else {
		if err := json.Unmarshal([]byte(cronservice.ServiceHeader), &header); err != nil {
			return c.JSON2(ERR_CODE_INTERNAL_ERROR, err.Error(), nil)
		}
	}

	var form map[string]string
	if cronservice.Form == "" {
		form = make(map[string]string)
	} else {
		if err := json.Unmarshal([]byte(cronservice.Form), &form); err != nil {
			return c.JSON2(ERR_CODE_INTERNAL_ERROR, err.Error(), nil)
		}
	}

	var body []byte
	if cronservice.ServiceBody == "" {
		body = make([]byte, 0)
	} else {
		body = base.StringToBytes(cronservice.ServiceBody)
	}

	// 加redis锁
	redis_client := GlobalCronService.CronTask.RedisClient
	m := make(map[string]interface{})
	m["lock_key"] = base.RedisLockKey(req.GetTaskName())
	m["redis_client"] = redis_client
	f := handler.Convert(req.GetTaskName(), cronservice.ServiceUrl, cronservice.ServiceMethod, cronservice.Method, header,
		form, body, m)

	// 发布暂停定时任务消息到redis，通知其它节点
	redis_msg := &handler.RedisMsg{
		cs.CronTask.UUID,
		handler.REDIS_RESTORE,
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

	time_format := cronservice.TimeFormat
	handler.AddTask(req.GetTaskName(), time_format, f)

	msg := fmt.Sprintf("定时任务：%s 恢复运行", req.GetTaskName())
	return c.JSON2(Status_OK, msg, nil)
}
