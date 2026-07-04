package router

import (
	"github.com/gin-gonic/gin"

	"github.com/RifqiIrawan/ai-rag-platform/services/api-gateway/internal/config"
	"github.com/RifqiIrawan/ai-rag-platform/services/api-gateway/internal/handlers"
)

func New(cfg *config.Config) *gin.Engine {
	r := gin.Default()
	r.GET("/health", handlers.Health)

	api := r.Group("/api/v1")
	{
		api.Any("/auth/*proxyPath", handlers.ProxyTo(cfg.AuthServiceURL))
		api.Any("/documents/*proxyPath", handlers.ProxyTo(cfg.DocumentServiceURL))
		api.Any("/rag/*proxyPath", handlers.ProxyTo(cfg.RagServiceURL))
		api.Any("/notifications/*proxyPath", handlers.ProxyTo(cfg.NotificationServiceURL))
	}
	return r
}
