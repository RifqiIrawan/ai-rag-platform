package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const notificationsChannel = "notifications"

type NotificationHandler struct {
	Redis *redis.Client
}

type publishRequest struct {
	Message string `json:"message" binding:"required"`
}

// Publish pushes a message onto the shared Redis pub/sub channel.
// Real notification delivery (per-user channels, WebSocket fan-out, email)
// is a follow-up sprint item.
func (h *NotificationHandler) Publish(c *gin.Context) {
	var req publishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Redis.Publish(c.Request.Context(), notificationsChannel, req.Message).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "published"})
}
