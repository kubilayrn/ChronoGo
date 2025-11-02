package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ListSentMessages godoc
// @Summary      Get list of sent messages
// @Description  Retrieve all messages that have been sent
// @Tags         messages
// @Accept       json
// @Produce      json
// @Success      200  {object}  ListSentMessagesResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /messages/sent [get]
func (h *Handler) ListSentMessages(c *gin.Context) {
	ctx := context.Background()

	messages, err := h.messageRepo.GetSentMessages(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to fetch sent messages",
		})
		return
	}

	messageResponses := make([]MessageResponse, len(messages))
	for i, msg := range messages {
		messageResponses[i] = MessageResponse{
			ID:        msg.ID,
			To:        msg.To,
			Content:   msg.Content,
			Status:    string(msg.Status),
			CreatedAt: msg.CreatedAt.Format(time.RFC3339),
			UpdatedAt: msg.UpdatedAt.Format(time.RFC3339),
		}
		if msg.SentAt != nil {
			messageResponses[i].SentAt = msg.SentAt.Format(time.RFC3339)
		}
		if msg.MessageID != nil {
			messageResponses[i].MessageID = msg.MessageID.String()
		}
	}

	c.JSON(http.StatusOK, ListSentMessagesResponse{
		Messages: messageResponses,
		Total:    len(messageResponses),
	})
}
