package models

import "time"

type CronService struct {
	ID               int64           `gorm:"primary_key" json:"id"`
	TaskName         string          `gorm:"column:task_name; type:varchar(30); unique_index:idx_name_format_url; not null" json:"task_name"`
	TimeFormat       string          `gorm:"column:time_format; type:varchar(30); unique_index:idx_name_format_url; not null" json:"time_format,omitempty"`
	ServiceUrl       string          `gorm:"column:service_url; type:varchar(100); unique_index:idx_name_format_url; not null" json:"service_url,omitempty"`
	ServiceMethod    string          `gorm:"column:service_method; type:varchar(50); not null" json:"service_method,omitempty"`
	Method           string          `gorm:"column:method; type:varchar(10); not null" json:"method,omitempty"`
	ServiceHeader    string          `gorm:"column:service_header; type:varchar(255)" json:"service_header,omitempty"`
	Form             string          `gorm:"column:form; type:text" json:"form,omitempty"`
	ServiceBody      string          `gorm:"column:service_body; type:text" json:"service_body,omitempty"`
	Status           int             `gorm:"column:status; type:tinyint(1); not null" json:"status,omitempty"`
	LockKey          string          `gorm:"column:lock_key; type:varchar(40); unique; not null" json:"lock_key,omitempty"`
	TaskDesc         string          `grom:"column:task_desc; type:text" json:"task_desc,omitempty"`
	CreatedAt        time.Time       `gorm:"column:create_time; type:datetime; not null;" json:"create_time,omitempty"`
	UpdatedAt        time.Time       `gorm:"column:update_time; type:datetime" json:"update_time,omitempty"`
}

func (CronService) TableName() string {
	return "cron_service"
}