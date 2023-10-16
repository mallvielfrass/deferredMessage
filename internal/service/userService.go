package service

import (
	"deferredMessage/internal/models"
	"deferredMessage/internal/repository"
	"deferredMessage/internal/utils"
	"fmt"
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
func (u *userService) CreateUser(name string, mail string, password string) (models.UserScheme, error) {
	isExist, err := u.CheckUserByMail(mail)
	if err != nil {
		return models.UserScheme{}, err
	}
	if isExist {
		return models.UserScheme{}, fmt.Errorf("user already exist")
	}
	hash, err := utils.HashPassword(password)
	if err != nil {
		return models.UserScheme{}, err
	}
	return u.repos.User.CreateUser(name, mail, hash)
}

// LoginUser
func (u *userService) LoginUser(mail string, password string) (models.UserScheme, error) {
	user, isExist, err := u.GetUserByMail(mail)
	if err != nil {
		return models.UserScheme{}, err
	}
	if !isExist {
		return models.UserScheme{}, fmt.Errorf("user not found")
	}
	if !utils.CheckPasswordHash(password, user.Hash) {
		return models.UserScheme{}, fmt.Errorf("user or password incorrect")
	}
	return user, nil
}

// GetUserByMail
func (u *userService) GetUserByMail(mail string) (models.UserScheme, bool, error) {
	return u.repos.User.GetUserByMail(mail)
}

// GetUserByID
func (u *userService) GetUserByID(id string) (models.UserScheme, bool, error) {
	return u.repos.User.GetUserByID(id)
}

// AddChatToUser
func (u *userService) AddChatToUser(chatID string, userID string) error {
	return u.repos.User.AddChatToUser(chatID, userID)
}
