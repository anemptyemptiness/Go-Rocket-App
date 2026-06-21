package config

import (
	"fmt"
	"time"
)

type redisConfig struct {
	Host              string        `yaml:"host"     env:"REDIS_HOST" env-default:"localhost"`
	Port              string        `yaml:"port"     env:"REDIS_PORT" env-default:"6379"`
	Password          string        `yaml:"password" env:"REDIS_PASSWORD"`
	Database          int32         `yaml:"database" env:"REDIS_DATABASE"`
	ConnectionTimeout time.Duration `yaml:"connection_timeout" env:"REDIS_CONNECTION_TIMEOUT" env-default:"5s"`
	MaxIdle           int           `yaml:"max_idle"           env:"REDIS_MAX_IDLE"           env-default:"10"`
	IdleTimeout       time.Duration `yaml:"idle_timeout"       env:"REDIS_IDLE_TIMEOUT"       env-default:"10s"`
}

func (c *redisConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
