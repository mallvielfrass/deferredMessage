package message

import "deferredMessage/internal/models"

type MessageListResponse struct {
	Messages []models.Message `json:"messages"`
}
