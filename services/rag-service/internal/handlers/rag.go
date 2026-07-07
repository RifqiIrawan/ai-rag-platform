package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/RifqiIrawan/ai-rag-platform/services/rag-service/internal/clients"
	"github.com/RifqiIrawan/ai-rag-platform/services/rag-service/internal/config"
)

type RagHandler struct {
	Qdrant    *clients.QdrantClient
	Ollama    *clients.OllamaClient
	Embedding *clients.EmbeddingClient
	Cfg       *config.Config
}

type queryRequest struct {
	Query string `json:"query" binding:"required"`
}

type querySource struct {
	Score float64 `json:"score"`
	Text  string  `json:"text,omitempty"`
}

// Query embeds the question, retrieves the closest chunks from Qdrant (if
// any have been indexed), and asks Ollama to answer grounded in that context.
func (h *RagHandler) Query(c *gin.Context) {
	var req queryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()

	vector, err := h.Embedding.Embed(ctx, req.Query)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "embedding-service: " + err.Error()})
		return
	}

	exists, err := h.Qdrant.CollectionExists(ctx, h.Cfg.QdrantCollection)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "qdrant: " + err.Error()})
		return
	}

	var results []clients.SearchResult
	if exists {
		results, err = h.Qdrant.Search(ctx, h.Cfg.QdrantCollection, vector, h.Cfg.RetrievalLimit)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "qdrant: " + err.Error()})
			return
		}
	}

	sources := make([]querySource, 0, len(results))
	contextParts := make([]string, 0, len(results))
	for _, r := range results {
		var payload struct {
			Text string `json:"text"`
		}
		_ = json.Unmarshal(r.Payload, &payload)
		sources = append(sources, querySource{Score: r.Score, Text: payload.Text})
		if payload.Text != "" {
			contextParts = append(contextParts, payload.Text)
		}
	}

	prompt := req.Query
	if len(contextParts) > 0 {
		prompt = fmt.Sprintf(
			"Answer the question using only the context below. If the context doesn't contain the answer, say you don't know.\n\nContext:\n%s\n\nQuestion: %s",
			strings.Join(contextParts, "\n---\n"), req.Query,
		)
	}

	answer, err := h.Ollama.Generate(ctx, h.Cfg.OllamaModel, prompt)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "ollama: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"query":   req.Query,
		"answer":  answer,
		"sources": sources,
	})
}
