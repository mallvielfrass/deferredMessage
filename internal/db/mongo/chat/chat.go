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
	Name              string `bson:"name"`
	ID                string `bson:"_id"`
	LinkOrIdInNetwork string `bson:"linkOrIdInNetwork"`
	Network           string `bson:"network"`
	Verified          bool   `bson:"verified"`
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
func (c Chat) UpdateChat(chat ChatScheme) error {
	_, err := c.ct.UpdateOne(context.TODO(), bson.M{"_id": chat.ID}, bson.M{"$set": chat})
	return err
}
