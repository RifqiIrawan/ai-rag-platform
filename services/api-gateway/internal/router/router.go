package router

import (
	"github.com/gin-gonic/gin"

	"github.com/RifqiIrawan/ai-rag-platform/services/api-gateway/internal/config"
	"github.com/RifqiIrawan/ai-rag-platform/services/api-gateway/internal/handlers"
	"github.com/RifqiIrawan/ai-rag-platform/services/api-gateway/internal/middleware"
)

// proxyGroup registers a downstream target at both the bare group path
// (e.g. "/documents") and every path below it ("/documents/foo"). gin's
// RedirectTrailingSlash means a lone "*proxyPath" wildcard never matches
// the bare path, so both registrations are required.
func proxyGroup(r gin.IRoutes, path, target string, extra ...gin.HandlerFunc) {
	handlers_ := append(append([]gin.HandlerFunc{}, extra...), handlers.ProxyTo(target))
	r.Any(path, handlers_...)
	r.Any(path+"/*proxyPath", handlers_...)
}

func New(cfg *config.Config) *gin.Engine {
	r := gin.Default()
	r.GET("/health", handlers.Health)

	api := r.Group("/api/v1")
	{
		proxyGroup(api, "/auth", cfg.AuthServiceURL)
		proxyGroup(api, "/documents", cfg.DocumentServiceURL, middleware.RequireAuth(cfg.JWTSecret))
		proxyGroup(api, "/rag", cfg.RagServiceURL)
		proxyGroup(api, "/notifications", cfg.NotificationServiceURL, middleware.RequireAuth(cfg.JWTSecret))
	}
	return r
}
