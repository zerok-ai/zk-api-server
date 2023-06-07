package model

import (
	httpConfig "github.com/zerok-ai/zk-utils-go/http/config"
	zkLoggerConfig "github.com/zerok-ai/zk-utils-go/logs/config"
	zkPostgres "github.com/zerok-ai/zk-utils-go/storage/sqlDB/postgres/config"
)

type ServerConfig struct {
	Host string `yaml:"host" env:"SRV_HOST,HOST" env-description:"Server host" env-default:"localhost"`
	Port string `yaml:"port" env:"SRV_PORT,PORT" env-description:"Server port" env-default:"8080"`
}

// ZkApiServerConfig https://github.com/ilyakaznacheev/cleanenv/blob/master/example/simple_config/example.go
// Config is an application configuration structure
type ZkApiServerConfig struct {
	Server     ServerConfig              `yaml:"server"`
	LogsConfig zkLoggerConfig.LogsConfig `yaml:"logs"`
	Http       httpConfig.HttpConfig     `yaml:"http"`
	Postgres   zkPostgres.PostgresConfig `yaml:"postgres"`
}
