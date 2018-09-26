package base

import "fmt"

// 生成redis分布式锁的key
func RedisLockKey(tname string) string {
	return fmt.Sprintf( tname + ".lock")
}

// redis channel name
func RedisPubSubChannel() string {
	return "pubsub.crontask"
}
