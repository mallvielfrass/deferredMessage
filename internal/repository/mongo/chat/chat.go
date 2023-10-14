package chat

import (
	"context"
	"fmt"
	"os"

	"github.com/mallvielfrass/fmc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slices"
)

type ChatScheme struct {
	Name          string             `bson:"name"`
	ID            primitive.ObjectID `bson:"_id"`
	LinkOrIdInBot string             `bson:"linkOrIdInBot"`
	BotIdentifier string             `bson:"botIdentifier"`
	BotID         string             `bson:"botID"`
	Verified      bool               `bson:"verified"`
	Creator       primitive.ObjectID `bson:"creator"`
	Hidden        bool               `bson:"hidden"`
}

type Chat struct {
	ct *mongo.Collection
}

func Init(driver *mongo.Database) Chat {
	collectionName := "chat"
	names, err := driver.ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {

		fmt.Println("List:", err)
		os.Exit(1)
	}
	//	fmt.Printf("names: %+v\n", names)
	if !slices.Contains(names, collectionName) {
		command := bson.M{"create": collectionName}
		var result bson.M
		if err := driver.RunCommand(context.TODO(), command).Decode(&result); err != nil {
			fmt.Printf("%s: %+v\n", fmc.WhoCallerIs(), err)
			os.Exit(1)
		}
	}

	return Chat{
		ct: driver.Collection(collectionName),
	}
}
func (c Chat) GetChatByID(id primitive.ObjectID) (ChatScheme, bool, error) {
	var findedChat ChatScheme
	err := c.ct.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&findedChat)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ChatScheme{}, false, nil
		}
		return ChatScheme{}, false, err
	}
	return findedChat, true, nil
}
func (c Chat) UpdateChat(chatId primitive.ObjectID, data map[string]interface{}) error {

	_, err := c.ct.UpdateOne(context.TODO(), bson.M{"_id": chatId}, bson.M{"$set": data})
	return err
}

func (c Chat) CreateChat(name string, botIdentifier string, botID string, userID primitive.ObjectID) (ChatScheme, error) {

	chat := ChatScheme{

		Name:          name,
		BotIdentifier: botIdentifier,
		BotID:         botID,
		LinkOrIdInBot: "",
		Verified:      false,
		Creator:       userID,
		Hidden:        false,
	}
	res, err := c.ct.InsertOne(context.TODO(), bson.M{"name": chat.Name, "botIdentifier": chat.BotIdentifier, "botID": chat.BotID,
		"linkOrIdInBot": chat.LinkOrIdInBot, "verified": chat.Verified, "creator": chat.Creator, "hidden": chat.Hidden})
	if err != nil {
		return ChatScheme{}, err
	}
	chat.ID = res.InsertedID.(primitive.ObjectID)
	return chat, nil
}
func (c Chat) GetChatsByArrayID(chats []primitive.ObjectID) ([]ChatScheme, error) {
	var findedChats []ChatScheme
	cur, err := c.ct.Find(context.TODO(), bson.M{"_id": bson.M{"$in": chats}})
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.TODO())
	fmt.Printf("cur: %#v\n", cur)
	for cur.Next(context.TODO()) {
		var chat ChatScheme
		err := cur.Decode(&chat)
		fmt.Printf("chat: %#v\n", chat)
		if err != nil {
			return nil, err
		}
		findedChats = append(findedChats, chat)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return findedChats, nil
}
