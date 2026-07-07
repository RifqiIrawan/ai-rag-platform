package main

import (
	"context"
	"log"

	"github.com/RifqiIrawan/ai-rag-platform/services/document-service/internal/config"
	"github.com/RifqiIrawan/ai-rag-platform/services/document-service/internal/db"
	"github.com/RifqiIrawan/ai-rag-platform/services/document-service/internal/router"
	"github.com/RifqiIrawan/ai-rag-platform/services/document-service/internal/storage"
)

func main() {
	cfg := config.Load()

	pool, err := db.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	minioClient, err := storage.NewMinioClient(cfg)
	if err != nil {
		log.Fatalf("failed to create minio client: %v", err)
	}

	if err := storage.EnsureBucket(context.Background(), minioClient, cfg.MinioBucket); err != nil {
		log.Fatalf("failed to ensure minio bucket %q: %v", cfg.MinioBucket, err)
	}

	r := router.New(cfg, pool, minioClient)

	log.Printf("document-service listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
