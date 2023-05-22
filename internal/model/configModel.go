package model

import (
	httpConfig "main/app/utils/http/config"
	logConfig "main/app/utils/logs/config"
	postgresConfig "main/app/utils/postgres/config"
)

type ServerConfig struct {
	Host string `yaml:"host" env:"SRV_HOST,HOST" env-description:"Server host" env-default:"localhost"`
	Port string `yaml:"port" env:"SRV_PORT,PORT" env-description:"Server port" env-default:"8080"`
}

// ZkApiServerConfig https://github.com/ilyakaznacheev/cleanenv/blob/master/example/simple_config/example.go
// Config is an application configuration structure
type ZkApiServerConfig struct {
	Server     ServerConfig                  `yaml:"server"`
	LogsConfig logConfig.LogsConfig          `yaml:"logs"`
	Http       httpConfig.HttpConfig         `yaml:"http"`
	Postgres   postgresConfig.PostgresConfig `yaml:"postgres"`
}
