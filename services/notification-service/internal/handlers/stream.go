package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var upgrader = websocket.Upgrader{
	// Origin checking is handled upstream by api-gateway; this service is
	// only reachable from the docker network plus the gateway proxy.
	CheckOrigin: func(r *http.Request) bool { return true },
}

type StreamHandler struct {
	Redis *redis.Client
}

// Stream upgrades the connection to a WebSocket and forwards any message
// published to the shared broadcast channel or to this user's own channel
// (see NotificationHandler.Publish) for as long as the client stays connected.
func (h *StreamHandler) Stream(c *gin.Context) {
	userID := c.GetHeader("X-User-Id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-User-Id header is required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("websocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sub := h.Redis.Subscribe(ctx, broadcastChannel, userChannel(userID))
	defer sub.Close()

	// The only way to detect the client closing the connection is to keep
	// reading from it; discard whatever it sends and cancel on any read error.
	go func() {
		for {
			if _, _, err := conn.NextReader(); err != nil {
				cancel()
				return
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-sub.Channel():
			if !ok {
				return
			}
			_ = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteJSON(gin.H{"channel": msg.Channel, "message": msg.Payload}); err != nil {
				return
			}
		}
	}
}
