package message

import (
	"context"
	"deferredMessage/internal/models"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/mallvielfrass/fmc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageScheme struct {
	Message     string             `bson:"message"`
	ID          primitive.ObjectID `bson:"_id"`
	ChatId      primitive.ObjectID `bson:"chat"`
	CreatorId   primitive.ObjectID `bson:"creator"`
	Time        time.Time          `bson:"time"`
	IsProcessed bool               `bson:"isprocessed"`
	IsSended    bool               `bson:"issended"`
	Error       string             `bson:"error"`
}
type Message struct {
	ct *mongo.Collection
}

func Init(driver *mongo.Database) Message {
	collectionName := "message"
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

	return Message{
		ct: driver.Collection(collectionName),
	}
}
func (c Message) GetMessagesList(tm time.Time) ([]models.Message, error) {
	var msgList []models.Message
	//find messages where time < tm and isprocessed = false

	cur, err := c.ct.Find(context.TODO(), bson.M{"time": bson.M{"$lt": tm}, "isprocessed": false})
	if err != nil {
		return []models.Message{}, err
	}
	if err != nil {
		return []models.Message{}, err
	}

	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var msg MessageScheme
		err := cur.Decode(&msg)
		if err != nil {
			return []models.Message{}, err
		}
		msgList = append(msgList, models.Message{
			Message:     msg.Message,
			Time:        msg.Time,
			Id:          msg.ID.Hex(),
			ChatId:      msg.ChatId.Hex(),
			CreatorId:   msg.CreatorId.Hex(),
			IsProcessed: msg.IsProcessed,
			IsSended:    msg.IsSended,
			Error:       msg.Error,
		})

	}
	if err := cur.Err(); err != nil {
		return []models.Message{}, err
	}
	return msgList, nil
}
