package service

import (
	"deferredMessage/internal/models"
	"deferredMessage/internal/repository"
	"fmt"
	"time"
)

type poolService struct {
	repos *repository.Repository
}

func NewPoolService(repos *repository.Repository) *poolService {
	return &poolService{
		repos: repos,
	}
}
func (p poolService) GetMsgList(period time.Duration) []models.Message {
	msgs, err := p.repos.Message.GetMessagesList(time.Now().Add(period))
	if err != nil {
		fmt.Println("Error: GetMsgList:", err)
		return []models.Message{}
	}
	return msgs
}
