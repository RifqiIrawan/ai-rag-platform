package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/RifqiIrawan/ai-rag-platform/services/rag-service/internal/clients"
)

type HealthHandler struct {
	Qdrant *clients.QdrantClient
	Ollama *clients.OllamaClient
}

func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"service":   "rag-service",
		"timestamp": time.Now().UTC(),
	})
}

func (h *HealthHandler) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	if err := h.Qdrant.Healthy(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable", "error": "qdrant: " + err.Error()})
		return
	}
	if err := h.Ollama.Healthy(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable", "error": "ollama: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}
