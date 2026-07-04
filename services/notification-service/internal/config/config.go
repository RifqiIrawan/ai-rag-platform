package config

import "os"

type Config struct {
	Port      string
	RedisAddr string
}

func Load() *Config {
	return &Config{
		Port:      getEnv("PORT", "8086"),
		RedisAddr: getEnv("REDIS_ADDR", "redis:6379"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
