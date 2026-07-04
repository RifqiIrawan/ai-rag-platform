package router

import (
	"github.com/gin-gonic/gin"

	"github.com/RifqiIrawan/ai-rag-platform/services/rag-service/internal/clients"
	"github.com/RifqiIrawan/ai-rag-platform/services/rag-service/internal/config"
	"github.com/RifqiIrawan/ai-rag-platform/services/rag-service/internal/handlers"
)

func New(cfg *config.Config, qdrant *clients.QdrantClient, ollama *clients.OllamaClient) *gin.Engine {
	r := gin.Default()

	health := &handlers.HealthHandler{Qdrant: qdrant, Ollama: ollama}
	r.GET("/health", health.Health)
	r.GET("/health/ready", health.Ready)

	rag := &handlers.RagHandler{Qdrant: qdrant, Ollama: ollama, Cfg: cfg}
	r.POST("/api/v1/rag/query", rag.Query)

	return r
}
