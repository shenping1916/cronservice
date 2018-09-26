package app

import (
	"cronservice/pb"
	"cronservice/api"
	"cronservice/config"
	framecfg "iceberg/frame/config"
)

func Start(cfgfile string) {
	var cfg config.Config
	framecfg.Parseconfig(cfgfile, &cfg)

	// 调用init方法
	api.GlobalCronService.Init(cfg)

	// 注册grpc server
	pb.RegisterCronServiceServer(api.GlobalCronService, &cfg.Base)
}
