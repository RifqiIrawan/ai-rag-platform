package config

import "os"

type Config struct {
	Port          string
	DatabaseURL   string
	MinioEndpoint string
	MinioRootUser string
	MinioRootPass string
	MinioBucket   string
	MinioUseSSL   bool
}

func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8082"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://raguser:ragpassword@postgres:5432/ai_rag_platform?sslmode=disable"),
		MinioEndpoint: getEnv("MINIO_ENDPOINT", "minio:9000"),
		MinioRootUser: getEnv("MINIO_ROOT_USER", "ragadmin"),
		MinioRootPass: getEnv("MINIO_ROOT_PASSWORD", "ragadminpassword"),
		MinioBucket:   getEnv("MINIO_BUCKET", "documents"),
		MinioUseSSL:   false,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
