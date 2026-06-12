package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"taskmanager/internal/sse"
	"taskmanager/internal/utils"
)

// SSEHandler streams real-time task events to the authenticated user.
type SSEHandler struct {
	broker *sse.Broker
	jwt    *utils.JWTManager
}

// NewSSEHandler constructs an SSEHandler.
func NewSSEHandler(broker *sse.Broker, jwt *utils.JWTManager) *SSEHandler {
	return &SSEHandler{broker: broker, jwt: jwt}
}

// Stream handles GET /api/events?token=... . The browser EventSource API cannot
// set an Authorization header, so the JWT is supplied as a query parameter and
// validated here directly.
func (h *SSEHandler) Stream(c *gin.Context) {
	claims, err := h.jwt.Verify(c.Query("token"))
	if err != nil {
		utils.Fail(c, utils.ErrUnauthorized)
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no") // disable proxy buffering (nginx)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		utils.Fail(c, utils.ErrInternal)
		return
	}

	events, unsubscribe := h.broker.Subscribe(claims.UserID)
	defer unsubscribe()

	// Initial comment opens the stream and confirms connectivity to the client.
	_, _ = c.Writer.WriteString(": connected\n\n")
	flusher.Flush()

	ctx := c.Request.Context()
	for {
		select {
		case <-ctx.Done():
			return
		case ev, open := <-events:
			if !open {
				return
			}
			payload, err := json.Marshal(ev)
			if err != nil {
				continue
			}
			_, _ = c.Writer.WriteString("event: " + ev.Type + "\n")
			_, _ = c.Writer.WriteString("data: " + string(payload) + "\n\n")
			flusher.Flush()
		}
	}
}
