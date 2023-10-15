package service

import (
	"deferredMessage/internal/models"
	"deferredMessage/internal/repository"
)

type chatService struct {
	repos *repository.Repository
}

func NewChatService(repos *repository.Repository) *chatService {
	return &chatService{
		repos: repos,
	}
}
func (c chatService) CreateChat(name string, botIdentifier string, botID string, userID string) (models.ChatScheme, error) {
	return c.repos.Chat.CreateChat(name, botIdentifier, botID, userID)
}

// GetChatByID(id string) (models.ChatScheme, bool, error)
func (c chatService) GetChatsByArrayID(chats []string) ([]models.ChatScheme, error) {
	return c.repos.Chat.GetChatsByArrayID(chats)
}

// UpdateChat(chatId string, data map[string]interface{}) error
func (c chatService) UpdateChat(chatId string, data map[string]interface{}) error {
	return c.repos.Chat.UpdateChat(chatId, data)
}

// GetChatsByArrayID(chats []string) ([]models.ChatScheme, error)
func (c chatService) GetChatByID(id string) (models.ChatScheme, bool, error) {
	return c.repos.Chat.GetChatByID(id)
}
