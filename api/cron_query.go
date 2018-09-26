package api

import (
	"fmt"
	"iceberg/frame"
	"cronservice/pb"
	"cronservice/models"
)

func (cs *CronService) GetTask(c frame.Context) error {
	var req pb.GetTaskReq
	if err := c.Bind(&req); err != nil {
		return err
	}

	var cronservice []*models.CronService
	if req.GetTaskName() != "" {
		if cs.CronTask.DbRead.Where("task_name = ?", req.GetTaskName()).Find(&cronservice).RecordNotFound() {
			msg := fmt.Errorf("定时任务：%s 不存在！", req.GetTaskName())
			return c.JSON2(ERR_CODE_BAD_NOT_FOUND, msg.Error(), nil)
		}
	} else {
		if err := cs.CronTask.DbRead.Find(&cronservice).Error; err != nil {
			msg := fmt.Errorf("数据库查询失败！原因：%v", err.Error())
			return c.JSON2(ERR_CODE_INTERNAL_ERROR, msg.Error(), nil)
		}
	}

	return c.JSON2(Status_OK, "", &cronservice)
}
