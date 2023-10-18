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
func (c Message) GetMessageByID(id string) (models.Message, bool, error) {
	idObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Message{}, false, err
	}
	var findedMessage MessageScheme
	err = Message{}.ct.FindOne(context.TODO(), bson.M{"_id": idObjectID}).Decode(&findedMessage)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Message{}, false, nil
		}
		return models.Message{}, false, err
	}
	return models.Message{
		Message:     findedMessage.Message,
		Time:        findedMessage.Time,
		Id:          findedMessage.ID.Hex(),
		ChatId:      findedMessage.ChatId.Hex(),
		CreatorId:   findedMessage.CreatorId.Hex(),
		IsProcessed: findedMessage.IsProcessed,
		IsSended:    findedMessage.IsSended,
		Error:       findedMessage.Error,
	}, true, nil
}

func (c Message) SetMessageIsProcessed(id string) error {
	idObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = c.ct.UpdateOne(context.TODO(), bson.M{"_id": idObjectID}, bson.M{"$set": bson.M{"isprocessed": true}})
	if err != nil {
		return err
	}
	return nil
}

// SetMessageError
func (c Message) SetMessageError(id string, errMsg string) error {
	idObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = c.ct.UpdateOne(context.TODO(), bson.M{"_id": idObjectID}, bson.M{"$set": bson.M{"error": errMsg}})
	if err != nil {
		return err
	}
	return nil
}

// SetIsSended
func (c Message) SetIsSended(id string) error {
	idObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = c.ct.UpdateOne(context.TODO(), bson.M{"_id": idObjectID}, bson.M{"$set": bson.M{"issended": true}})
	if err != nil {
		return err
	}
	return nil
}
func (c Message) GetListOfAllMessages(creatorId string, offset int, limit int) ([]models.Message, error) {
	creatorObjectID, err := primitive.ObjectIDFromHex(creatorId)
	if err != nil {
		return []models.Message{}, err
	}
	var msgList []models.Message
	cur, err := c.ct.Find(context.TODO(), bson.M{"creator": creatorObjectID})
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
func (c Message) getMessageByObjectID(id primitive.ObjectID) (MessageScheme, error) {
	var msg MessageScheme
	err := c.ct.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&msg)
	if err != nil {
		return MessageScheme{}, err
	}
	return msg, nil
}
func (c Message) CreateNewMessage(creatorId string, msg models.Message) (models.Message, error) {
	creatorObjectID, err := primitive.ObjectIDFromHex(creatorId)
	if err != nil {
		return models.Message{}, err
	}
	chatIdObjectID, err := primitive.ObjectIDFromHex(msg.ChatId)
	if err != nil {
		return models.Message{}, err
	}
	res, err := c.ct.InsertOne(context.TODO(), bson.M{
		"message":     msg.Message,
		"time":        msg.Time,
		"creator":     creatorObjectID,
		"chat":        chatIdObjectID,
		"isprocessed": false,
		"issended":    false,
		"error":       "",
	})
	if err != nil {
		return models.Message{}, err
	}
	idObjectID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return models.Message{}, err
	}
	findedMsg, err := c.getMessageByObjectID(idObjectID)
	if err != nil {
		return models.Message{}, err
	}

	return models.Message{
		Message:     findedMsg.Message,
		Time:        findedMsg.Time,
		Id:          findedMsg.ID.Hex(),
		ChatId:      findedMsg.ChatId.Hex(),
		CreatorId:   findedMsg.CreatorId.Hex(),
		IsProcessed: findedMsg.IsProcessed,
		IsSended:    findedMsg.IsSended,
		Error:       findedMsg.Error,
	}, nil

}
