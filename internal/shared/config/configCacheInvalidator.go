package config

type ConfigCacheInvalidator struct {
	Addr          string   `mapstructure:"REDIS_ADDR"`
	RedisPassword string   `mapstructure:"REDIS_PASSWORD"`
	Broker        []string `mapstructure:"BROKER"`
}
