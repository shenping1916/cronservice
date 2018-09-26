package config

import "iceberg/frame/config"

type Config struct {
	Env   string          `json:"env"`
	Base  config.BaseCfg  `json:"baseCfg"`
	Mysql config.MysqlCfg `json:"mysqlCfg"`
	Redis config.RedisCfg `json:"redisCfg"`
}
