package service

import (
	"deferredMessage/internal/models"
	"deferredMessage/internal/repository"
	"fmt"
	"time"
)

type sessionService struct {
	repos *repository.Repository
}

func NewSessionService(repos *repository.Repository) *sessionService {
	return &sessionService{
		repos: repos,
	}
}
func (s *sessionService) CheckSession(token string) (models.SessionScheme, error) {
	session, exist, err := s.repos.Session.GetSessionByID(token)
	if err != nil {
		return models.SessionScheme{}, err
	}
	if !exist {
		return models.SessionScheme{}, fmt.Errorf("invalid token")
	}
	if session.Expire < time.Now().Unix() {
		return models.SessionScheme{}, fmt.Errorf("session expired")
	}
	if !session.Valid {
		return models.SessionScheme{}, fmt.Errorf("token revoked")
	}
	return session, nil
}

// CreateSession(UserID string, expire int64, ip string) (models.SessionScheme, error)
func (s *sessionService) CreateSession(UserID string, expire int64, ip string) (models.SessionScheme, error) {
	return s.repos.Session.CreateSession(UserID, expire, ip)
}

// GetUserByMail(mail string) (models.UserScheme, bool, error)
func (s *sessionService) GetUserByMail(mail string) (models.UserScheme, bool, error) {
	return s.repos.User.GetUserByMail(mail)
}
