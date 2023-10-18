package message

import "deferredMessage/internal/models"

type MessageListResponse struct {
	Messages []models.Message `json:"messages"`
}
type NewMessageRequest struct {
	Message string `json:"message" binding:"required"`
	ChatId  string `json:"chatId" binding:"required"`
	Time    int64  `json:"time" binding:"required"`
}
type MessageResponse struct {
	Message models.Message `json:"message"`
}
