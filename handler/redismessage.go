package handler

import (
	"cronservice/models"
)

const (
	REDIS_ADD = iota
	REDIS_DEL
	REDIS_MOTIFT
	REDIS_PAUSE
	REDIS_RESTORE
)

type RedisMsg struct {
	UUID         string     `json:"uuid"`
	Action 	     int        `json:"action"`
	*models.CronService
}
