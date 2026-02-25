package app

import (
	"log"
	"os"
	"strings"

	"github.com/patato8984/Shop/internal/shared/config"
	"github.com/spf13/viper"
)

func LoadApiConfig() (*config.ApiConfig, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AddConfigPath("..")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("error read api.yml: %s", err)
	}
	var cfg config.ApiConfig
	v.SetConfigFile("admins.yml")
	v.SetConfigFile("/app/admins.yml")
	if err := v.MergeInConfig(); err != nil {
		log.Printf("admins config not found (skipping): %v", err)
	}
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("error unmarshal config.yml: %s", err)
	}
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	cfg.JwtKey = getSecret("/run/secrets/jwt_key")
	cfg.SharedApiBankKey = getSecret("/run/secrets/shared_bank_key")
	return &cfg, nil
}
func getSecret(secretPath string) string {
	data, err := os.ReadFile(secretPath)
	if err != nil {
		log.Fatalf("error save secret %s", err)
	}
	return strings.TrimSpace(string(data))
}
func LoadCDCConfig() *config.ChangeDataCaptureConfig {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AddConfigPath("..")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("error read config.yml: %s", err)
	}
	var cfg config.ChangeDataCaptureConfig
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("error unmarshal config.yml: %s", err)
	}
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	return &cfg
}

func LoadCacheConfig() *config.ConfigCacheInvalidator {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AddConfigPath("..")
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("error read config.yml: %s", err)
	}
	var cfg config.ConfigCacheInvalidator
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("error unmarshal config.yml: %s", err)
	}
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	return &cfg
}
