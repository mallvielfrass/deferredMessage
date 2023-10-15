package bot

import (
	"context"
	"deferredMessage/internal/models"
	"fmt"
	"os"

	"github.com/mallvielfrass/fmc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slices"
)

type BotScheme struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name"`
	BotLink     string             `bson:"botLink"`
	BotType     string             `bson:"botType"`
	Creator     primitive.ObjectID `bson:"creator"`
	Platform    string             `bson:"platform"`
	HashedToken string             `bson:"hashedToken"`
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
func (c Bot) GetBotByID(id string) (models.BotScheme, bool, error) {
	botObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.BotScheme{}, false, err
	}
	var findedBot BotScheme
	err = c.ct.FindOne(context.TODO(), bson.M{"_id": botObjectID}).Decode(&findedBot)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.BotScheme{}, false, nil
		}
		return models.BotScheme{}, false, err
	}
	return models.BotScheme{
		Name:     findedBot.Name,
		BotLink:  findedBot.BotLink,
		Creator:  findedBot.Creator.Hex(),
		Platform: findedBot.Platform,
		//HashedToken: bot.HashedToken,
		ID: findedBot.ID.Hex(),
	}, true, nil
}

//	func (c Bot) UpdateBot(Bot models.BotScheme) error {
//		_, err := c.ct.UpdateOne(context.TODO(), bson.M{"_id": Bot.ID}, bson.M{"$set": Bot})
//		return err
//	}
//

func (c Bot) CreateBot(name string, botLink string, creator string, platform string, hashedToken string) (models.BotScheme, error) {
	creatorObjectID, err := primitive.ObjectIDFromHex(creator)
	if err != nil {
		return models.BotScheme{}, fmt.Errorf("invalid creator id: %s", creator)
	}
	bot := BotScheme{
		Name:        name,
		BotLink:     botLink,
		Creator:     creatorObjectID,
		Platform:    platform,
		HashedToken: hashedToken,
	}
	res, err := c.ct.InsertOne(context.TODO(), bson.M{
		"name":        name,
		"botLink":     botLink,
		"creator":     creatorObjectID,
		"platform":    platform,
		"hashedToken": hashedToken,
	})
	if err != nil {
		return models.BotScheme{}, err
	}

	bot.ID = res.InsertedID.(primitive.ObjectID)
	//fmt.Printf("bot: %+v\n", bot)
	return models.BotScheme{
		Name:     name,
		BotLink:  botLink,
		Creator:  creator,
		Platform: platform,
		//HashedToken: hashedToken,
		ID: bot.ID.Hex(),
	}, nil
}
func (c Bot) UpdateBot(botId string, data map[string]interface{}) (models.BotScheme, bool, error) {
	botObjectID, err := primitive.ObjectIDFromHex(botId)
	if err != nil {
		return models.BotScheme{}, false, err
	}
	_, err = c.ct.UpdateOne(context.TODO(), bson.M{"_id": botObjectID}, bson.M{"$set": data})
	if err != nil {
		return models.BotScheme{}, false, err
	}
	bot, isExist, err := c.GetBotByID(botId)
	return bot, isExist, err
}

// GetAllBots
func (c Bot) GetAllBots() ([]models.BotScheme, error) {
	var findedBots []models.BotScheme
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
		findedBots = append(findedBots, models.BotScheme{
			Name:     bot.Name,
			BotLink:  bot.BotLink,
			Creator:  bot.Creator.Hex(),
			Platform: bot.Platform,
			//HashedToken: bot.HashedToken,
			ID: bot.ID.Hex(),
		})
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return findedBots, nil
}

// func (c Bot) GetBotsByArrayID(Bots []primitive.ObjectID) ([]models.BotScheme, error) {
// 	var findedBots []models.BotScheme
// 	cur, err := c.ct.Find(context.TODO(), bson.M{"_id": bson.M{"$in": Bots}})
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer cur.Close(context.TODO())
// 	for cur.Next(context.TODO()) {
// 		var Bot models.BotScheme
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
