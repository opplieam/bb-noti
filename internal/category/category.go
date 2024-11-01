package category

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/opplieam/bb-noti-api/internal/state"
)

type Handler struct {
	ClientState *state.ClientState
}

func NewHandler(clientState *state.ClientState) *Handler {
	return &Handler{
		ClientState: clientState,
	}
}

func (h *Handler) SSE(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// TODO: Use UserID and Param to create Key
	// For now use a dummy value
	userKey := c.Query("os")
	clientCh := h.ClientState.AddClient(userKey)
	defer close(clientCh)
	defer h.ClientState.RemoveClient(userKey)
	clientClose := c.Writer.CloseNotify()
	// Handshake
	c.Writer.Flush()
loop:
	for {
		select {
		case v := <-clientCh:
			_, _ = fmt.Fprintf(c.Writer, "data: %s\n\n", v)
			c.Writer.Flush()
		case <-clientClose:
			//fmt.Println("Client disconnected")
			break loop
		}
	}
}
