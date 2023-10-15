package service

import (
	"deferredMessage/internal/models"
	"deferredMessage/internal/repository"
)

type botService struct {
	repos *repository.Repository
}

func NewBotService(repos *repository.Repository) *botService {
	return &botService{
		repos: repos,
	}
}
func (b botService) GetBotByID(id string) (models.BotScheme, bool, error) {
	return b.repos.Bot.GetBotByID(id)
}
func (b botService) CreateBot(name string, botLink string, creator string, platform string, hashedToken string) (models.BotScheme, error) {
	return b.repos.Bot.CreateBot(name, botLink, creator, platform, hashedToken)
}
func (b botService) UpdateBot(botId string, data map[string]interface{}) (models.BotScheme, bool, error) {
	return b.repos.Bot.UpdateBot(botId, data)
}
func (b botService) GetAllBots() ([]models.BotScheme, error) {
	return b.repos.Bot.GetAllBots()
}
