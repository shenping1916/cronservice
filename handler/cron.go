package handler

import (
	"sync"
	"time"
	"reflect"
	"strings"
	"strconv"
	"context"
	"encoding/json"
	"cronservice/config"
	"cronservice/models"
	"cronservice/utility/base"
	log "iceberg/frame/icelog"
	"github.com/pkg/errors"
	"github.com/jinzhu/gorm"
	"github.com/go-redis/redis"
	"github.com/astaxie/beego/toolbox"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type CronTask struct {
	Mu                    *sync.RWMutex
	UUID                  string
	Ctx                   context.Context
	Cancel                context.CancelFunc
	DbRead                *gorm.DB
	DbWrite               *gorm.DB
	RedisClient           *redis.Client
	RedisPubSub           *redis.PubSub
	StartChan             chan struct{}
	ErrChan               chan error
}

var (
	Driver = "mysql"
	Tables = []string{"cron_service"}
)

func NewCronTask(cfg config.Config) *CronTask {
	crontask := &CronTask{}
	crontask.Mu = new(sync.RWMutex)
	crontask.UUID = base.NewV4().String()
	crontask.Ctx, crontask.Cancel = context.WithCancel(context.Background())
	crontask.StartChan = make(chan struct{}, 1)

	var err error
	port := strconv.Itoa(cfg.Mysql.Port)
	read := []string {
		cfg.Mysql.User,
		":",
		cfg.Mysql.Psw,
		"@tcp(",
		cfg.Mysql.Host.Read,
		":",
		port,
		")/",
		cfg.Mysql.DbName,
		"?charset=utf8mb4&parseTime=True&loc=Local",
	}
	read_dsn := base.StringSplice(read)
	crontask.DbRead, err = gorm.Open(Driver, read_dsn)
	if err != nil {
		panic(err.Error())
	}

	write := []string {
		cfg.Mysql.User,
		":",
		cfg.Mysql.Psw,
		"@tcp(",
		cfg.Mysql.Host.Write,
		":",
		port,
		")/",
		cfg.Mysql.DbName,
		"?charset=utf8mb4&parseTime=True&loc=Local",
	}
	write_dsn := base.StringSplice(write)
	crontask.DbWrite, err = gorm.Open(Driver, write_dsn)
	if err != nil {
		panic(err.Error())
	}

	if cfg.Env != "prod" {
		crontask.DbRead.LogMode(true)
		crontask.DbWrite.LogMode(true)
	} else {
		crontask.DbRead.LogMode(false)
		crontask.DbWrite.LogMode(false)
	}

	// 初始化redis客户端
	crontask.RedisClient = redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr: cfg.Redis.Addr,
		Password: cfg.Redis.Psw,
		DB: cfg.Redis.DBNo,
	})

	// 初始化redis pubsub
	crontask.RedisPubSub = crontask.RedisClient.Subscribe(base.RedisPubSubChannel())

	return crontask
}

func (ct *CronTask) Run() {
	for {
		select {
		case <- ct.Ctx.Done():
			// 资源释放
			StopAll()
			ct.DbRead.Close()
		    ct.DbWrite.Close()
		    ct.RedisPubSub.Close()
		    ct.RedisClient.Close()
		    log.Close()

			return
		case msg := <- ct.RedisPubSub.Channel():
			log.Debugf("收到redis订阅消息：%s", msg.String())

			// 处理接收到的redis sub信息
			var m = make(map[string]interface{})
		    if err := json.Unmarshal(base.StringToBytes(msg.Payload), &m); err != nil {
		    	log.Error(err.Error())
				return
			}

			task_name := m["task_name"].(string)
			time_format := m["time_format"].(string)
			service_url := m["service_url"].(string)
			service_method := m["service_method"].(string)
			method := m["method"].(string)

			var header map[string]string
			if value, ok := m["service_header"]; !ok {
				header = make(map[string]string)
			} else {
				if err := json.Unmarshal([]byte(value.(string)), &header); err != nil {
					log.Error(err.Error())
					return
				}
			}

			var form map[string]string
			if value, ok := m["form"]; !ok {
				form = make(map[string]string)
			} else {
				if err := json.Unmarshal([]byte(value.(string)), &form); err != nil {
					log.Error(err.Error())
					return
				}
			}

			var body []byte
			if value, ok := m["service_body"]; !ok {
				body = make([]byte, 0)
			} else {
				body = base.StringToBytes(value.(string))
			}

			// 加redis锁
			mm := make(map[string]interface{})
			mm["lock_key"] = base.RedisLockKey(task_name)
			mm["redis_client"] = ct.RedisClient

			switch m["action"].(float64) {
			case REDIS_ADD:
				if !reflect.DeepEqual(m["uuid"].(string), ct.UUID) {
					// 添加定时任务
					f := Convert(task_name, service_url, service_method, method, header, form, body, mm)
					AddTask(task_name, time_format, f)
					log.Infof("定时任务：%s 注册成功", task_name)
				}
			case REDIS_DEL:
				if !reflect.DeepEqual(m["uuid"].(string), ct.UUID) {
					// 删除定时任务
					StopSpeciTask(task_name)
					log.Infof("删除定时任务：%s 成功", task_name)
				}
			case REDIS_MOTIFT:
				if !reflect.DeepEqual(m["uuid"].(string), ct.UUID) {
					// 先删除定时任务
					StopSpeciTask(task_name)

					// sleep 1s
					time.Sleep(1 * time.Second)

					// 重新添加新的定时任务
					f := Convert(task_name, service_url, service_method, method, header, form, body, mm)
					AddTask(task_name, time_format, f)
					log.Infof("定时任务：%s 修改成功", task_name)
				}
			case REDIS_PAUSE:
				if !reflect.DeepEqual(m["uuid"].(string), ct.UUID) {
					// 暂停原定时任务
					StopSpeciTask(task_name)
					log.Infof("定时任务：%s 暂停成功", task_name)
				}
			case REDIS_RESTORE:
				if !reflect.DeepEqual(m["uuid"].(string), ct.UUID) {
					// 恢复原定时任务
					f := Convert(task_name, service_url, service_method, method, header, form, body, mm)
					AddTask(task_name, time_format, f)
					log.Infof("定时任务：%s 恢复运行", task_name)
				}
			default:
				return
			}

		case <- ct.StartChan:
			// 初始化建表
			if err := ct.InitDB(); err != nil {
				panic(err.Error())
			}

			// 从数据库加载所有定时任务，并开始运行
			var cronservice []*models.CronService
			if err := ct.DbRead.Where("status = ?", 1).Find(&cronservice).Error; err != nil {
				log.Error(err.Error())
				return
			}

			for _, v := range cronservice {
				var header map[string]string
				if v.ServiceHeader == "" {
					header = make(map[string]string)
				} else {
					if err := json.Unmarshal([]byte(v.ServiceHeader), &header); err != nil {
						log.Error(err.Error())
						return
					}
				}

				var form map[string]string
				if v.Form == "" {
					form = make(map[string]string)
				} else {
					if err := json.Unmarshal([]byte(v.Form), &form); err != nil {
						log.Error(err.Error())
						return
					}
				}

				var body []byte
				if v.ServiceBody == "" {
					body = make([]byte, 0)
				} else {
					body = base.StringToBytes(v.ServiceBody)
				}

				// 生成调用函数，加入定时任务
				m := make(map[string]interface{})
				m["lock_key"] = v.LockKey
				m["redis_client"] = ct.RedisClient
				f := Convert(v.TaskName, v.ServiceUrl, v.ServiceMethod, v.Method, header, form, body, m)
				task_ := toolbox.NewTask(v.TaskName, v.TimeFormat, f)
				toolbox.AddTask(v.TaskName, task_)

				log.Infof("定时任务：%s 加载成功", v.TaskName)
			}

			RunTask()
		}
	}
}

func (ct *CronTask) Stop() {
	ct.Cancel()
}

func(ct *CronTask) InitDB() error {
	for i := 0; i < len(Tables); i ++ {
		table := Tables[i]
		ok := ct.DbWrite.HasTable(table)
		if !ok {
			switch tablename := strings.ToLower(table);tablename {
			case "cron_service":
				if err := ct.DbWrite.DB().Ping(); err != nil {
					return errors.Errorf("数据库状态异常!!! 原因: %v !", err)
				}

				cron_service := models.CronService{}
				if err := ct.DbWrite.CreateTable(cron_service).Error; err != nil {
					return errors.Errorf("创建数据库表：[%s] 失败！原因：%v !", tablename, err)
				}
				log.Debugf("创建数据库表：[%s] 成功", tablename)
			}
		}
	}

	return nil
}

