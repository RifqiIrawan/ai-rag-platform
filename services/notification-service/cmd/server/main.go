package main

import (
	"log"

	"github.com/redis/go-redis/v9"

	"github.com/RifqiIrawan/ai-rag-platform/services/notification-service/internal/config"
	"github.com/RifqiIrawan/ai-rag-platform/services/notification-service/internal/router"
)

func main() {
	cfg := config.Load()

	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	defer rdb.Close()

	r := router.New(cfg, rdb)

	log.Printf("notification-service listening on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
