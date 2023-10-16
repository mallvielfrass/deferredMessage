package chat

import (
	"context"
	"deferredMessage/internal/models"
	"fmt"
	"os"

	"github.com/mallvielfrass/fmc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
func (c Chat) GetChatByID(id string) (models.ChatScheme, bool, error) {
	var findedChat ChatScheme
	err := c.ct.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&findedChat)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.ChatScheme{}, false, nil
		}
		return models.ChatScheme{}, false, err
	}
	return models.ChatScheme{
		Name:          findedChat.Name,
		ID:            findedChat.ID.Hex(),
		LinkOrIdInBot: findedChat.LinkOrIdInBot,
		BotIdentifier: findedChat.BotIdentifier,
		BotID:         findedChat.BotID,
		Verified:      findedChat.Verified,
		Creator:       findedChat.Creator.Hex(),
		Hidden:        findedChat.Hidden,
	}, true, nil
}
func (c Chat) UpdateChat(chatId string, data map[string]interface{}) error {
	chatObjectID, err := primitive.ObjectIDFromHex(chatId)
	if err != nil {
		return err
	}
	_, err = c.ct.UpdateOne(context.TODO(), bson.M{"_id": chatObjectID}, bson.M{"$set": data})
	return err
}

func (c Chat) CreateChat(name string, botIdentifier string, botID string, userID string) (models.ChatScheme, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return models.ChatScheme{}, err
	}
	chat := ChatScheme{

		Name:          name,
		BotIdentifier: botIdentifier,
		BotID:         botID,
		LinkOrIdInBot: "",
		Verified:      false,
		Creator:       userObjectID,
		Hidden:        false,
	}
	res, err := c.ct.InsertOne(context.TODO(), bson.M{"name": chat.Name, "botIdentifier": chat.BotIdentifier, "botID": chat.BotID,
		"linkOrIdInBot": chat.LinkOrIdInBot, "verified": chat.Verified, "creator": chat.Creator, "hidden": chat.Hidden})
	if err != nil {
		return models.ChatScheme{}, err
	}
	chat.ID = res.InsertedID.(primitive.ObjectID)

	return models.ChatScheme{
		Name:          chat.Name,
		ID:            chat.ID.Hex(),
		LinkOrIdInBot: chat.LinkOrIdInBot,
		BotIdentifier: chat.BotIdentifier,
		BotID:         chat.BotID,
		Verified:      chat.Verified,
		Creator:       chat.Creator.Hex(),
		Hidden:        chat.Hidden,
	}, nil
}
func (c Chat) GetChatsByArrayID(chats []string) ([]models.ChatScheme, error) {
	var findedChats []models.ChatScheme
	var chatsObjectID []primitive.ObjectID
	for _, chat := range chats {
		oID, err := primitive.ObjectIDFromHex(chat)
		if err != nil {
			return nil, err
		}
		chatsObjectID = append(chatsObjectID, oID)
	}
	cur, err := c.ct.Find(context.TODO(), bson.M{"_id": bson.M{"$in": chatsObjectID}})
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
		findedChats = append(findedChats, models.ChatScheme{
			Name:          chat.Name,
			ID:            chat.ID.Hex(),
			LinkOrIdInBot: chat.LinkOrIdInBot,
			BotIdentifier: chat.BotIdentifier,
			BotID:         chat.BotID,
			Verified:      chat.Verified,
			Creator:       chat.Creator.Hex(),
			Hidden:        chat.Hidden,
		})
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return findedChats, nil
}

// GetChatsListByCreatorWithLimits
func (c Chat) GetChatsListByCreatorWithLimits(userId string, count int, offset int) ([]models.ChatScheme, error) {
	userObjectID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	opts := options.Find().SetSkip(int64(offset)).SetLimit(int64(count))
	var findedChats []models.ChatScheme
	cur, err := c.ct.Find(context.TODO(), bson.M{"creator": userObjectID}, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var chat ChatScheme
		err := cur.Decode(&chat)
		if err != nil {
			return nil, err
		}
		findedChats = append(findedChats, models.ChatScheme{
			Name:          chat.Name,
			ID:            chat.ID.Hex(),
			LinkOrIdInBot: chat.LinkOrIdInBot,
			BotIdentifier: chat.BotIdentifier,
			BotID:         chat.BotID,
			Verified:      chat.Verified,
			Creator:       chat.Creator.Hex(),
			Hidden:        chat.Hidden,
		})
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return findedChats, nil
}
