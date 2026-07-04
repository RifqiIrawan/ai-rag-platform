package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/RifqiIrawan/ai-rag-platform/services/rag-service/internal/clients"
	"github.com/RifqiIrawan/ai-rag-platform/services/rag-service/internal/config"
)

type RagHandler struct {
	Qdrant *clients.QdrantClient
	Ollama *clients.OllamaClient
	Cfg    *config.Config
}

type queryRequest struct {
	Query string `json:"query" binding:"required"`
}

// Query is a placeholder endpoint: the full retrieval (embed query via
// embedding-service, search Qdrant, assemble context) + generation
// (prompt Ollama) pipeline is out of scope for Sprint 1, which covers
// service scaffolding and infra wiring only.
func (h *RagHandler) Query(c *gin.Context) {
	var req queryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "rag pipeline not yet implemented",
		"query":   req.Query,
	})
}
