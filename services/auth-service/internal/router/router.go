package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/RifqiIrawan/ai-rag-platform/services/auth-service/internal/config"
	"github.com/RifqiIrawan/ai-rag-platform/services/auth-service/internal/handlers"
)

func New(cfg *config.Config, db *pgxpool.Pool) *gin.Engine {
	r := gin.Default()

	health := &handlers.HealthHandler{DB: db}
	r.GET("/health", health.Health)
	r.GET("/health/ready", health.Ready)

	auth := &handlers.AuthHandler{DB: db, Cfg: cfg}
	v1 := r.Group("/api/v1/auth")
	{
		v1.POST("/register", auth.Register)
		v1.POST("/login", auth.Login)
	}

	return r
}
