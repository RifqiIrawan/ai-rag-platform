package main

import (
	"log"

	"github.com/RifqiIrawan/ai-rag-platform/services/api-gateway/internal/config"
	"github.com/RifqiIrawan/ai-rag-platform/services/api-gateway/internal/router"
)

func main() {
	cfg := config.Load()
	r := router.New(cfg)

	log.Printf("api-gateway listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
