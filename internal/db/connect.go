package db

import (
	"context"
	"deferredMessage/internal/db/mongo/chat"
	"deferredMessage/internal/db/mongo/session"
	"deferredMessage/internal/db/mongo/user"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Chat interface {
}
type User interface {
	CheckUser(mail string) (bool, error)
	CreateUser(name, mail, hash string) (user.UserScheme, error)
	GetUserByID(id primitive.ObjectID) (user.UserScheme, bool, error)
	GetUserByMail(mail string) (user.UserScheme, bool, error)
}
type Session interface {
	GetSessionByID(id primitive.ObjectID) (session.SessionScheme, bool, error)
	CreateSession(UserID string, expire int64, ip string) (session.SessionScheme, error)
}
type Collection struct {
	Chat    Chat
	User    User
	Session Session
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
		Chat:    chat.Init(db.driver),
		User:    user.Init(db.driver),
		Session: session.Init(db.driver),
	}
}

func (db *DB) Disconnect() {
	if err := db.client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}
