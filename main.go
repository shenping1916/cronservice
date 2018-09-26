package main

import (
	"os"
	"flag"
	"path/filepath"
	"cronservice/utility/app"
	log "iceberg/frame/icelog"
)

var (
	AppName = "cronservice"
)

var (
	cfgFile  = flag.String("config-path", "cronservice.json", "config file")
	logLevel = flag.String("loglevel", "debug", "log level")
	logPath  = flag.String("log-path", "", "log path")
)

func main() {
	// 设置进程的当前目录为程序所在的路径
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	os.Chdir(dir)

	// 解析命令行参数
	flag.Parse()
	if *logLevel != "" {
		log.SetLevel(*logLevel)
	}

	// 判断配置文件是否存在
	_, err := os.Stat(*cfgFile)
	if err != nil && os.IsNotExist(err) {
		os.Exit(-1)
	}

	log.Debugf("logpath=%s,loglevel=%s", *logPath, *logLevel)

	// 启动cronservic
	app.Start(*cfgFile)

	// 关闭cronservic
	app.Shutdown()
}