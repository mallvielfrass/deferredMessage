package service

import (
	"deferredMessage/internal/models"
	"deferredMessage/internal/repository"
)

type platformService struct {
	repos *repository.Repository
}

func NewPlatformService(repos *repository.Repository) *platformService {
	return &platformService{
		repos: repos,
	}
}
func (p platformService) CreatePlatform(name string) (models.PlatformScheme, error) {
	return p.repos.Platform.CreatePlatform(name)
}

// GetAllPlatforms() ([]models.PlatformScheme, error)
func (p platformService) GetAllPlatforms() ([]models.PlatformScheme, error) {
	return p.repos.Platform.GetAllPlatforms()
}

// GetPlatformByName(name string) (models.PlatformScheme, bool, error)
func (p platformService) GetPlatformByName(name string) (models.PlatformScheme, bool, error) {
	return p.repos.Platform.GetPlatformByName(name)
}
