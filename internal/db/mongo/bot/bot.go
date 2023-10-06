package bot

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

type BotScheme struct {
	ID         string             `bson:"_id"`
	Name       string             `bson:"name"`
	Identifier string             `bson:"identifier"`
	BotLink    string             `bson:"botLink"`
	BotType    string             `bson:"botType"`
	Creator    primitive.ObjectID `bson:"creator"`
}

type Bot struct {
	ct *mongo.Collection
}

func Init(driver *mongo.Database) Bot {
	collectionName := "bot"
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

	return Bot{
		ct: driver.Collection(collectionName),
	}
}
func (c Bot) GetBotByID(id primitive.ObjectID) (BotScheme, bool, error) {
	var findedBot BotScheme
	err := c.ct.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&findedBot)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return BotScheme{}, false, nil
		}
		return BotScheme{}, false, err
	}
	return findedBot, true, nil
}

//	func (c Bot) UpdateBot(Bot BotScheme) error {
//		_, err := c.ct.UpdateOne(context.TODO(), bson.M{"_id": Bot.ID}, bson.M{"$set": Bot})
//		return err
//	}
//
// GetBotByIdentifier
func (c Bot) GetBotByIdentifier(identifier string) (BotScheme, bool, error) {
	var findedBot BotScheme
	err := c.ct.FindOne(context.TODO(), bson.M{"identifier": identifier}).Decode(&findedBot)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return BotScheme{}, false, nil
		}
		return BotScheme{}, false, err
	}
	return findedBot, true, nil
}
func (c Bot) CreateBot(name string, identifier string, botLink string, botType string, creator primitive.ObjectID) (BotScheme, error) {
	id := primitive.NewObjectID()
	Bot := BotScheme{
		ID:         id.Hex(),
		Name:       name,
		Identifier: identifier,
		BotLink:    botLink,
		BotType:    botType,
		Creator:    creator,
	}
	_, err := c.ct.InsertOne(context.TODO(), Bot)
	if err != nil {
		return BotScheme{}, err
	}
	return Bot, nil
}
func (c Bot) UpdateBot(botIdentifier string, data map[string]string) (BotScheme, bool, error) {
	_, err := c.ct.UpdateOne(context.TODO(), bson.M{"identifier": botIdentifier}, bson.M{"$set": data})
	if err != nil {
		return BotScheme{}, false, err
	}
	bot, isExist, err := c.GetBotByIdentifier(botIdentifier)
	return bot, isExist, err
}

// GetAllBots
func (c Bot) GetAllBots() ([]BotScheme, error) {
	var findedBots []BotScheme
	cur, err := c.ct.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var bot BotScheme
		err := cur.Decode(&bot)
		if err != nil {
			return nil, err
		}
		findedBots = append(findedBots, bot)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return findedBots, nil
}

// func (c Bot) GetBotsByArrayID(Bots []primitive.ObjectID) ([]BotScheme, error) {
// 	var findedBots []BotScheme
// 	cur, err := c.ct.Find(context.TODO(), bson.M{"_id": bson.M{"$in": Bots}})
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer cur.Close(context.TODO())
// 	for cur.Next(context.TODO()) {
// 		var Bot BotScheme
// 		err := cur.Decode(&Bot)
// 		if err != nil {
// 			return nil, err
// 		}
// 		findedBots = append(findedBots, Bot)
// 	}
// 	if err := cur.Err(); err != nil {
// 		return nil, err
// 	}
// 	return findedBots, nil
// }
