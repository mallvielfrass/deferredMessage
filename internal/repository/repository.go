package repository

import (
	"deferredMessage/internal/repository/mongo/bot"
	"deferredMessage/internal/repository/mongo/chat"
	"deferredMessage/internal/repository/mongo/platform"
	"deferredMessage/internal/repository/mongo/session"
	"deferredMessage/internal/repository/mongo/user"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Bot interface {
	GetBotByID(id primitive.ObjectID) (bot.BotScheme, bool, error)

	CreateBot(name string, botLink string, creator primitive.ObjectID, platform string, hashedToken string) (bot.BotScheme, error)
	UpdateBot(botId primitive.ObjectID, data map[string]interface{}) (bot.BotScheme, bool, error)
	GetAllBots() ([]bot.BotScheme, error)
}
type Chat interface {
	GetChatByID(id primitive.ObjectID) (chat.ChatScheme, bool, error)
	UpdateChat(chatId primitive.ObjectID, data map[string]interface{}) error
	CreateChat(name string, botIdentifier string, botID string, userID primitive.ObjectID) (chat.ChatScheme, error)
	GetChatsByArrayID(chats []primitive.ObjectID) ([]chat.ChatScheme, error)
}
type Platform interface {
	// GetPlatformByID(id primitive.ObjectID) (platform.PlatformScheme, bool, error)
	CreatePlatform(name string) (platform.PlatformScheme, error)
	// UpdatePlatform(platformIdentifier string, data map[string]string) (platform.PlatformScheme, bool, error)
	GetAllPlatforms() ([]platform.PlatformScheme, error)
	GetPlatformByName(name string) (platform.PlatformScheme, bool, error)
}
type User interface {
	CheckUserByMail(mailOrUsername string) (bool, error)
	CreateUser(name, mail, hash string) (user.UserScheme, error)
	GetUserByID(id primitive.ObjectID) (user.UserScheme, bool, error)
	SetUserAdmin(id primitive.ObjectID) (user.UserScheme, bool, error)
	GetUserByMail(mail string) (user.UserScheme, bool, error)
	AddChatToUser(chatID primitive.ObjectID, userID primitive.ObjectID) error
}
type Session interface {
	GetSessionByID(id primitive.ObjectID) (session.SessionScheme, bool, error)
	CreateSession(UserID primitive.ObjectID, expire int64, ip string) (session.SessionScheme, error)
}
type Repository struct {
	Chat     Chat
	Platform Platform
	Bot      Bot
	User     User
	Session  Session
}

func NewRepository(driver *mongo.Database) *Repository {
	return &Repository{
		Chat:     chat.Init(driver),
		Platform: platform.Init(driver),
		Bot:      bot.Init(driver),
		User:     user.Init(driver),
		Session:  session.Init(driver),
	}
}
