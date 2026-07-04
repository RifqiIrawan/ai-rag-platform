package config

import "os"

type Config struct {
	Port                string
	QdrantURL           string
	OllamaURL           string
	OllamaModel         string
	EmbeddingServiceURL string
}

func Load() *Config {
	return &Config{
		Port:                getEnv("PORT", "8085"),
		QdrantURL:           getEnv("QDRANT_URL", "http://qdrant:6333"),
		OllamaURL:           getEnv("OLLAMA_URL", "http://ollama:11434"),
		OllamaModel:         getEnv("OLLAMA_MODEL", "llama3.2"),
		EmbeddingServiceURL: getEnv("EMBEDDING_SERVICE_URL", "http://embedding-service:8084"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
