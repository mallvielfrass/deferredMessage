package repository

import (
	"deferredMessage/internal/models"
	"deferredMessage/internal/repository/mongo/bot"
	"deferredMessage/internal/repository/mongo/chat"
	"deferredMessage/internal/repository/mongo/message"
	"deferredMessage/internal/repository/mongo/platform"
	"deferredMessage/internal/repository/mongo/session"
	"deferredMessage/internal/repository/mongo/user"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Bot interface {
	GetBotByID(id string) (models.BotScheme, bool, error)

	CreateBot(name string, botLink string, creator string, platform string, hashedToken string) (models.BotScheme, error)
	UpdateBot(botId string, data map[string]interface{}) (models.BotScheme, bool, error)
	GetAllBots() ([]models.BotScheme, error)
}
type Chat interface {
	GetChatByID(id string) (models.ChatScheme, bool, error)
	UpdateChat(chatId string, data map[string]interface{}) error
	CreateChat(name string, botIdentifier string, botID string, userID string) (models.ChatScheme, error)
	GetChatsByArrayID(chats []string) ([]models.ChatScheme, error)
	GetChatsListByCreatorWithLimits(userId string, count int, offset int) ([]models.ChatScheme, error)
}
type Platform interface {
	// GetPlatformByID(id string) (models.PlatformScheme, bool, error)
	CreatePlatform(name string) (models.PlatformScheme, error)
	// UpdatePlatform(platformIdentifier string, data map[string]string) (models.PlatformScheme, bool, error)
	GetAllPlatforms() ([]models.PlatformScheme, error)
	GetPlatformByName(name string) (models.PlatformScheme, bool, error)
}
type User interface {
	CheckUserByMail(mailOrUsername string) (bool, error)
	CreateUser(name, mail, hash string) (models.UserScheme, error)
	GetUserByID(id string) (models.UserScheme, bool, error)
	SetUserAdmin(id string) (models.UserScheme, bool, error)
	GetUserByMail(mail string) (models.UserScheme, bool, error)
	AddChatToUser(chatID string, userID string) error
}
type Session interface {
	GetSessionByID(id string) (models.SessionScheme, bool, error)
	CreateSession(UserID string, expire int64, ip string) (models.SessionScheme, error)
}
type Message interface {
	GetMessagesList(tm time.Time) ([]models.Message, error)
	GetMessageByID(id string) (models.Message, bool, error)
	SetMessageIsProcessed(id string) error
	SetMessageError(id string, errMsg string) error
	SetIsSended(id string) error
}
type Repository struct {
	Chat     Chat
	Platform Platform
	Bot      Bot
	User     User
	Session  Session
	Message  Message
}

func NewRepository(driver *mongo.Database) *Repository {
	return &Repository{
		Chat:     chat.Init(driver),
		Platform: platform.Init(driver),
		Bot:      bot.Init(driver),
		User:     user.Init(driver),
		Session:  session.Init(driver),
		Message:  message.Init(driver),
	}
}
