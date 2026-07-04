package main

import (
	"context"
	"log"

	"github.com/RifqiIrawan/ai-rag-platform/services/auth-service/internal/config"
	"github.com/RifqiIrawan/ai-rag-platform/services/auth-service/internal/db"
	"github.com/RifqiIrawan/ai-rag-platform/services/auth-service/internal/router"
)

func main() {
	cfg := config.Load()

	pool, err := db.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	r := router.New(cfg, pool)

	log.Printf("auth-service listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
