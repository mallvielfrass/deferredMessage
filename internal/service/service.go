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
	CreateSession(UserID string, expire int64, ip string) (models.SessionScheme, error)
}
type UserService interface {
	UserIsAdmin(userID string) (bool, error)
	SetUserAdmin(userID string) (models.UserScheme, bool, error)
	CheckUserByMail(mail string) (bool, error)
	CreateUser(name string, mail string, hash string) (models.UserScheme, error)
	GetUserByMail(mail string) (models.UserScheme, bool, error)
	GetUserByID(userID string) (models.UserScheme, bool, error)
	AddChatToUser(chatID string, userID string) error
}
type PlatformService interface {
	CreatePlatform(name string) (models.PlatformScheme, error)
	GetAllPlatforms() ([]models.PlatformScheme, error)
	GetPlatformByName(name string) (models.PlatformScheme, bool, error)
}
type BotService interface {
	GetBotByID(id string) (models.BotScheme, bool, error)
	CreateBot(name string, botLink string, creator string, platform string, hashedToken string) (models.BotScheme, error)
	UpdateBot(botId string, data map[string]interface{}) (models.BotScheme, bool, error)
	GetAllBots() ([]models.BotScheme, error)
}
type ChatService interface {
	GetChatByID(id string) (models.ChatScheme, bool, error)
	UpdateChat(chatId string, data map[string]interface{}) error
	CreateChat(name string, botIdentifier string, botID string, userID string) (models.ChatScheme, error)
	GetChatsByArrayID(chats []string) ([]models.ChatScheme, error)
}
type Service struct {
	repository *repository.Repository
	SessionService
	AuthorizationService
	UserService
	PlatformService
	BotService
	ChatService
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		repository:           repos,
		SessionService:       NewSessionService(repos),
		AuthorizationService: NewAuthService(repos),
		UserService:          NewUserService(repos),
		PlatformService:      NewPlatformService(repos),
		BotService:           NewBotService(repos),
		ChatService:          NewChatService(repos),
	}
}
