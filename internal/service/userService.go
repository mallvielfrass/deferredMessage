package service

import "deferredMessage/internal/repository"

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
