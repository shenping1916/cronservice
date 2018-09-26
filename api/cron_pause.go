package api

import (
	"fmt"
	"encoding/json"
	"iceberg/frame"
	"cronservice/pb"
	"cronservice/models"
	"cronservice/handler"
	"cronservice/utility/base"
)

func (cs *CronService) PauseTask(c frame.Context) error {
	var req pb.PauseTaskReq
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
		if cronservice.Status == 0 {
			msg := fmt.Errorf("定时任务：%s 已暂停，无需再次执行！", req.GetTaskName())
			return c.JSON2(ERR_CODE_FORBIDDEN  , msg.Error(), nil)
		}
	}

	if err := cs.CronTask.DbWrite.Model(&cronservice).Where("task_name = ?", req.GetTaskName()).Update("status", 0).Error; err != nil {
		msg := fmt.Errorf("数据库更新失败！原因：%v", err.Error())
		return c.JSON2(ERR_CODE_INTERNAL_ERROR, msg.Error(), nil)
	}

	// 发布暂停定时任务消息到redis，通知其它节点
	redis_msg := &handler.RedisMsg{
		cs.CronTask.UUID,
		handler.REDIS_PAUSE,
		cronservice,
	}

	b, err := json.Marshal(redis_msg)
	if err != nil {
		return c.JSON2(ERR_CODE_INTERNAL_ERROR, err.Error(), nil)
	}

	redis_client := GlobalCronService.CronTask.RedisClient
	if err := redis_client.Publish(base.RedisPubSubChannel(), base.BytesToString(b)).Err(); err != nil {
		msg := fmt.Errorf("消息发布到redis失败！原因：%v", err.Error())
		return c.JSON2(ERR_CODE_INTERNAL_ERROR, msg.Error(), nil)
	}

	// 暂停原定时任务
	handler.StopSpeciTask(req.GetTaskName())
	msg := fmt.Sprintf("定时任务：%s 暂停成功", req.GetTaskName())
	return c.JSON2(Status_OK, msg, nil)
}