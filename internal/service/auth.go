package service

import (
	"deferredMessage/internal/models"
	"deferredMessage/internal/repository"
)

type authService struct {
	repository *repository.Repository
}

func (a authService) Register(user models.RegisterUser) (models.Session, models.UserIdentify, error) {
	return models.Session{}, models.UserIdentify{}, nil
}
func (a authService) Login(user models.AuthUser) (models.Session, error) {
	return models.Session{}, nil
}
func (a authService) CheckUserExist(userMail string) (bool, error) {
	return false, nil
}
func NewAuthService(repos *repository.Repository) *authService {
	return &authService{
		repository: repos,
	}
}
