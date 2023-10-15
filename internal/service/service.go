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
type SessionService interface {
	CheckSession(token string) (models.SessionScheme, error)
}
type Service struct {
	repository *repository.Repository
	SessionService
	AuthorizationService
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		repository:     repos,
		SessionService: NewSessionService(repos),
	}
}
