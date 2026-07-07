package main

import (
	"log"

	"github.com/RifqiIrawan/ai-rag-platform/services/rag-service/internal/clients"
	"github.com/RifqiIrawan/ai-rag-platform/services/rag-service/internal/config"
	"github.com/RifqiIrawan/ai-rag-platform/services/rag-service/internal/router"
)

func main() {
	cfg := config.Load()

	qdrant := clients.NewQdrantClient(cfg.QdrantURL)
	ollama := clients.NewOllamaClient(cfg.OllamaURL)
	embedding := clients.NewEmbeddingClient(cfg.EmbeddingServiceURL)

	r := router.New(cfg, qdrant, ollama, embedding)

	log.Printf("rag-service listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
