package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const broadcastChannel = "notifications:broadcast"

func userChannel(userID string) string {
	return "notifications:user:" + userID
}

type NotificationHandler struct {
	Redis *redis.Client
}

type publishRequest struct {
	Message      string `json:"message" binding:"required"`
	TargetUserID string `json:"target_user_id"`
}

// Publish pushes a message onto Redis pub/sub: the shared broadcast channel,
// or a specific user's channel if target_user_id is set. Any client connected
// to StreamHandler's /ws endpoint for that channel receives it in real time.
func (h *NotificationHandler) Publish(c *gin.Context) {
	var req publishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	channel := broadcastChannel
	if req.TargetUserID != "" {
		channel = userChannel(req.TargetUserID)
	}

	if err := h.Redis.Publish(c.Request.Context(), channel, req.Message).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "published", "channel": channel})
}
