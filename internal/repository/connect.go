package repository

import (
	"context"
	"deferredMessage/internal/repository/mongo/bot"
	"deferredMessage/internal/repository/mongo/chat"
	"deferredMessage/internal/repository/mongo/platform"
	"deferredMessage/internal/repository/mongo/session"
	"deferredMessage/internal/repository/mongo/user"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection struct {
	Chat     Chat
	User     User
	Session  Session
	Bot      Bot
	Platform Platform
}
type DB struct {
	Driver      *mongo.Database
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
		Driver: client.Database(dbname),
	}
	//	databaseInstance.mountSchemes()
	return databaseInstance, nil
	//client.Database("GoToster").Collection("tasks")
	//collection = client.Database("GoToster") //.Collection("tasks")
}
func (db *DB) mountSchemes() {
	db.Collections = &Collection{
		Chat:     chat.Init(db.Driver),
		User:     user.Init(db.Driver),
		Session:  session.Init(db.Driver),
		Bot:      bot.Init(db.Driver),
		Platform: platform.Init(db.Driver),
	}
}

func (db *DB) Disconnect() {
	if err := db.client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}
