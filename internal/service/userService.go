package service

import (
	"deferredMessage/internal/models"
	"deferredMessage/internal/repository"
)

type userService struct {
	repos *repository.Repository
}

func NewUserService(repos *repository.Repository) *userService {
	return &userService{
		repos: repos,
	}
}
func (u *userService) UserIsAdmin(userID string) (bool, error) {
	user, isExist, err := u.repos.User.GetUserByID(userID)
	if err != nil {
		return false, err
	}
	if !isExist {
		return false, nil
	}
	return user.Admin, nil
}
func (u *userService) SetUserAdmin(userID string) (models.UserScheme, bool, error) {
	return u.repos.User.SetUserAdmin(userID)
}

// CheckUserByMail
func (u *userService) CheckUserByMail(mail string) (bool, error) {
	return u.repos.User.CheckUserByMail(mail)
}

// CreateUser
func (u *userService) CreateUser(name string, mail string, hash string) (models.UserScheme, error) {
	return u.repos.User.CreateUser(name, mail, hash)
}

// GetUserByMail
func (u *userService) GetUserByMail(mail string) (models.UserScheme, bool, error) {
	return u.repos.User.GetUserByMail(mail)
}

// GetUserByID
func (u *userService) GetUserByID(id string) (models.UserScheme, bool, error) {
	return u.repos.User.GetUserByID(id)
}
