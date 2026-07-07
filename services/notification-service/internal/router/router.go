package router

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/RifqiIrawan/ai-rag-platform/services/notification-service/internal/config"
	"github.com/RifqiIrawan/ai-rag-platform/services/notification-service/internal/handlers"
)

func New(cfg *config.Config, rdb *redis.Client) *gin.Engine {
	r := gin.Default()

	health := &handlers.HealthHandler{Redis: rdb}
	r.GET("/health", health.Health)
	r.GET("/health/ready", health.Ready)

	notif := &handlers.NotificationHandler{Redis: rdb}
	r.POST("/api/v1/notifications", notif.Publish)

	stream := &handlers.StreamHandler{Redis: rdb}
	r.GET("/api/v1/notifications/ws", stream.Stream)

	return r
}
