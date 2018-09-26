package app

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"cronservice/api"
	log "iceberg/frame/icelog"
)

var wg sync.WaitGroup

// 优雅关闭主程序
func Shutdown() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)

	wg.Add(1)
	go func() {
		for {
			signal := <- ch
			switch signal {
			case syscall.SIGTERM, syscall.SIGINT:
				log.Debugf("Receive signal：[%v], shutting down [cronservice] service", signal)
				close(ch)

				api.GlobalCronService.CronTask.Stop()
				wg.Done()
			case syscall.SIGHUP:
				return
			}

		}
	}()

	wg.Wait()
	os.Exit(0)
}