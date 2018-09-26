package api

import (
	"iceberg/frame"
	"cronservice/config"
	"cronservice/handler"
)

type CronService struct {
	CronTask *handler.CronTask
}

func NewCronService() *CronService {
	cronservice := new(CronService)
	frame.Prepare(func(c frame.Context) error {
		c.Response().SetHeader("Content-Type", "application/json;charset=utf-8")
		c.Response().SetHeader("Access-Control-Allow-Origin", "*")
		c.Response().SetHeader("Access-Control-Allow-Methods", "POST,OPTIONS")
		c.Response().SetHeader("Access-Control-Allow-Headers", "Accept, Content-Type, userId")
		return nil
	})

	return cronservice
}

func (cs *CronService) Init(cfg config.Config) {
	cs.CronTask = handler.NewCronTask(cfg)
	go cs.CronTask.Run()
	cs.CronTask.StartChan <- struct{}{}
}
