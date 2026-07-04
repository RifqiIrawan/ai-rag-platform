package config

import "os"

type Config struct {
	Port                   string
	AuthServiceURL         string
	DocumentServiceURL     string
	RagServiceURL          string
	NotificationServiceURL string
}

func Load() *Config {
	return &Config{
		Port:                   getEnv("PORT", "8080"),
		AuthServiceURL:         getEnv("AUTH_SERVICE_URL", "http://auth-service:8081"),
		DocumentServiceURL:     getEnv("DOCUMENT_SERVICE_URL", "http://document-service:8082"),
		RagServiceURL:          getEnv("RAG_SERVICE_URL", "http://rag-service:8085"),
		NotificationServiceURL: getEnv("NOTIFICATION_SERVICE_URL", "http://notification-service:8086"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
