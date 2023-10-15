package platform

import (
	"context"
	"deferredMessage/internal/models"
	"fmt"
	"os"
	"slices"

	"github.com/mallvielfrass/fmc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PlatformScheme struct {
	Name string             `bson:"name"`
	ID   primitive.ObjectID `bson:"_id"`
}
type Platform struct {
	ct *mongo.Collection
}

func Init(driver *mongo.Database) Platform {
	collectionName := "platform"
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

	return Platform{
		ct: driver.Collection(collectionName),
	}
}

// GetAllPlatforms() ([]platform.PlatformScheme, error)
func (c Platform) GetAllPlatforms() ([]models.PlatformScheme, error) {
	var platforms []PlatformScheme
	cursor, err := c.ct.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.TODO(), &platforms)
	if err != nil {
		return nil, err
	}
	var foundPlatforms []models.PlatformScheme
	for _, platform := range platforms {
		foundPlatforms = append(foundPlatforms, models.PlatformScheme{
			Name: platform.Name,
			ID:   platform.ID.Hex(),
		})
	}
	return foundPlatforms, nil
}

// CreatePlatform(name string) (platform.PlatformScheme, error)
func (c Platform) CreatePlatform(name string) (models.PlatformScheme, error) {
	res, err := c.ct.InsertOne(context.TODO(), bson.M{"name": name})
	if err != nil {
		return models.PlatformScheme{}, err
	}

	return models.PlatformScheme{
		Name: name,
		ID:   res.InsertedID.(primitive.ObjectID).Hex(),
	}, nil
}

// GetPlatformByName(name string) (platform.PlatformScheme, bool, error)
func (c Platform) GetPlatformByName(name string) (models.PlatformScheme, bool, error) {
	var findedPlatform PlatformScheme
	err := c.ct.FindOne(context.TODO(), bson.M{"name": name}).Decode(&findedPlatform)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.PlatformScheme{}, false, nil
		}
		return models.PlatformScheme{}, false, err
	}
	return models.PlatformScheme{
		Name: name, ID: findedPlatform.ID.Hex(),
	}, true, nil
}
