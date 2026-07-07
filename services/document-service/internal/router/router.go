package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"

	"github.com/RifqiIrawan/ai-rag-platform/services/document-service/internal/clients"
	"github.com/RifqiIrawan/ai-rag-platform/services/document-service/internal/config"
	"github.com/RifqiIrawan/ai-rag-platform/services/document-service/internal/handlers"
)

func New(cfg *config.Config, db *pgxpool.Pool, minioClient *minio.Client) *gin.Engine {
	r := gin.Default()

	health := &handlers.HealthHandler{DB: db, Minio: minioClient}
	r.GET("/health", health.Health)
	r.GET("/health/ready", health.Ready)

	docs := &handlers.DocumentHandler{
		DB:        db,
		Minio:     minioClient,
		OCR:       clients.NewOCRClient(cfg.OCRServiceURL),
		Embedding: clients.NewEmbeddingClient(cfg.EmbeddingServiceURL),
		Cfg:       cfg,
	}
	v1 := r.Group("/api/v1/documents")
	{
		v1.POST("", docs.Upload)
		v1.GET("", docs.List)
	}

	return r
}
