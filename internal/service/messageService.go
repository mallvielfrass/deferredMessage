package service

import (
	"deferredMessage/internal/models"
	"deferredMessage/internal/repository"
)

type messageService struct {
	repos *repository.Repository
}

func NewMessageService(repos *repository.Repository) *messageService {
	return &messageService{
		repos: repos,
	}
}
func (m messageService) GetListOfAllMessages(creatorId string, offset int, limit int) ([]models.Message, error) {
	return m.repos.Message.GetListOfAllMessages(creatorId, offset, limit)
}
