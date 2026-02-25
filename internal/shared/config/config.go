package config

import (
	"time"
)

type ApiConfig struct {
	JwtKey           string `mapstructure:"JWT_KEY"`
	SharedApiBankKey string `mapstructure:"SHARED_BANK_API"`
	DbPatch          string `mapstructure:"dsn"`
	Addr             string `mapstructure:"REDIS_ADDR"`
	RedisPassword    string `mapstructure:"REDIS_PASSWORD"`
	Port             string `mapstructure:"port"`
	Admins           Admin  `mapstructure:"admins"`
	TimeOut          map[string]time.Duration
}

type Admin []struct {
	NickName string `mapstructure:"nickName"`
	Name     string `mapstructure:"name"`
	Mail     string `mapstructure:"mail"`
	Password string `mapstructure:"password"`
}

func (c *ApiConfig) FillDefaults() {
	if c.TimeOut == nil {
		c.TimeOut = make(map[string]time.Duration)
	}
	c.TimeOut["POST /api/v1/user/order/{id}/payment"] = 10 * time.Second
}
