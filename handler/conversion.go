package handler

import (
	"encoding/json"
	"iceberg/frame"
	"iceberg/frame/protocol"
	log "iceberg/frame/icelog"
	"cronservice/utility/base"
	"github.com/pkg/errors"
	"github.com/go-redis/redis"
	"github.com/bsm/redis-lock"
	"github.com/nobugtodebug/go-objectid"
)

//var option = &lock.Options{
//	//LockTimeout: time.Duration(10) * time.Minute,
//	LockTimeout: time.Duration(2) * time.Minute,
//	RetryCount: 5,
//	RetryDelay: time.Duration(200) * time.Millisecond,
//}

func Convert(tname string, svrURI,svrMethod,Method string, header,form map[string]string, body []byte, args ...map[string]interface{}) func() error {
	return func() error {
		bussiness := func() error {
			task := readyTask(header, form, body, svrURI, svrMethod, Method)
			resp, err:= frame.DeliverTo(task)
			if err != nil {
				log.Error(err.Error())
				return  err
			}

			data, err := json.Marshal(&resp)
			if err != nil {
				log.Error(err.Error())
				return err
			}

			log.Debugf("定时任务：%s 调度成功：%s", tname, base.BytesToString(data))
			return nil
		}

		if len(args) > 0 {
			arg := args[0]

			// redis锁实例
			lock_key := arg["lock_key"].(string)
			redis_client := arg["redis_client"].(*redis.Client)
			lock, err := lock.Obtain(redis_client, lock_key, nil)
			if err != nil {
				log.Errorf("定时任务：%s 无法获取Redis锁，原因：%v，放弃本次定时调度！", tname, err.Error())
				return err
			}
			log.Debugf("定时任务：%s 成功获取Redis锁，开始执行调度任务", tname)

			// 业务调度
			if err := bussiness(); err != nil {
				return err
			}

			// 释放redis锁
			if err := lock.Unlock(); err != nil {
				log.Errorf("Redis锁释放出错！原因：%v", err.Error())
				return err
			}

			// 检查redis锁是否已经释放
			if lock.IsLocked() {
				err := errors.New("Redis锁未被释放！！！")
				log.Error(err.Error())
				return err
			}

		} else {
			// 业务调度
			if err := bussiness(); err != nil {
				return err
			}
		}

		return nil
	}
}

func readyTask(header,form map[string]string, body []byte, svrURI,svrMethod,Method string) *protocol.Proto{
	var task protocol.Proto
	task.Bizid = objectid.New().String()
	task.Header = header
	task.Form = form
	task.RequestID =frame.GetInnerID()
	task.ServeURI = svrURI
	task.Format = protocol.RestfulFormat_JSON
	task.ServeMethod = svrMethod
	switch Method {
	case "POST":
		task.Method = protocol.RestfulMethod_POST
	case "PUT":
		task.Method = protocol.RestfulMethod_PUT
	case "GET":
		task.Method = protocol.RestfulMethod_GET
	case "DELETE":
		task.Method = protocol.RestfulMethod_DELETE
	}
	task.Body = body

	return &task
}