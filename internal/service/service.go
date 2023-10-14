package service

import (
	"deferredMessage/internal/models"
	"deferredMessage/internal/repository"
)

type AuthorizationService interface {
	Register(user models.RegisterUser) (models.Session, models.UserIdentify, error)
	Login(user models.AuthUser) (models.Session, error)
	CheckUserExist(userMail string) (bool, error)
}
type Service struct {
	Repository *repository.Repository
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Repository: repos,
	}
}
