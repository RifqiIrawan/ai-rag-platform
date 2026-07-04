package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
)

type HealthHandler struct {
	DB    *pgxpool.Pool
	Minio *minio.Client
}

func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"service":   "document-service",
		"timestamp": time.Now().UTC(),
	})
}

func (h *HealthHandler) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	if err := h.DB.Ping(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable", "error": "postgres: " + err.Error()})
		return
	}

	if _, err := h.Minio.ListBuckets(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable", "error": "minio: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}
