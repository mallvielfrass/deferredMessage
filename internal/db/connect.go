package db

import (
	"context"
	"deferredMessage/internal/db/mongo/chat"
	"deferredMessage/internal/db/mongo/network"
	"deferredMessage/internal/db/mongo/platform"
	"deferredMessage/internal/db/mongo/session"
	"deferredMessage/internal/db/mongo/user"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Network interface {
	GetNetworkByID(id primitive.ObjectID) (network.NetworkScheme, bool, error)
	GetNetworkByIdentifier(identifier string) (network.NetworkScheme, bool, error)
	CreateNetwork(name string, identifier string, botLink string, botType string, creator primitive.ObjectID) (network.NetworkScheme, error)
	UpdateNetwork(networkIdentifier string, data map[string]string) (network.NetworkScheme, bool, error)
	GetAllNetworks() ([]network.NetworkScheme, error)
}
type Chat interface {
	GetChatByID(id primitive.ObjectID) (chat.ChatScheme, bool, error)
	UpdateChat(chatId primitive.ObjectID, data map[string]interface{}) error
	CreateChat(name string, networkIdentifier string, networkID string, userID primitive.ObjectID) (chat.ChatScheme, error)
	GetChatsByArrayID(chats []primitive.ObjectID) ([]chat.ChatScheme, error)
}
type Platform interface {
	// GetPlatformByID(id primitive.ObjectID) (platform.PlatformScheme, bool, error)
	CreatePlatform(name string) (platform.PlatformScheme, error)
	// UpdatePlatform(platformIdentifier string, data map[string]string) (platform.PlatformScheme, bool, error)
	GetAllPlatforms() ([]platform.PlatformScheme, error)
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
type Collection struct {
	Chat     Chat
	User     User
	Session  Session
	Network  Network
	Platform Platform
}
type DB struct {
	driver      *mongo.Database
	Collections *Collection
	client      *mongo.Client
}

func ConnectDB(url, dbname string) (DB, error) {
	fmt.Printf("url: %s,dbname: %s\n", url, dbname)
	//var collection *mongo.Collection
	var ctx = context.TODO()
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return DB{}, err
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	databaseInstance := DB{
		driver: client.Database(dbname),
	}
	databaseInstance.mountSchemes()
	return databaseInstance, nil
	//client.Database("GoToster").Collection("tasks")
	//collection = client.Database("GoToster") //.Collection("tasks")
}
func (db *DB) mountSchemes() {
	db.Collections = &Collection{
		Chat:     chat.Init(db.driver),
		User:     user.Init(db.driver),
		Session:  session.Init(db.driver),
		Network:  network.Init(db.driver),
		Platform: platform.Init(db.driver),
	}
}

func (db *DB) Disconnect() {
	if err := db.client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}
