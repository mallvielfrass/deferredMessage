package network

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

type NetworkScheme struct {
	ID         string             `bson:"_id"`
	Name       string             `bson:"name"`
	Identifier string             `bson:"identifier"`
	BotLink    string             `bson:"botLink"`
	BotType    string             `bson:"botType"`
	Creator    primitive.ObjectID `bson:"creator"`
}

type Network struct {
	ct *mongo.Collection
}

func Init(driver *mongo.Database) Network {
	collectionName := "network"
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

	return Network{
		ct: driver.Collection(collectionName),
	}
}
func (c Network) GetNetworkByID(id primitive.ObjectID) (NetworkScheme, bool, error) {
	var findedNetwork NetworkScheme
	err := c.ct.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&findedNetwork)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return NetworkScheme{}, false, nil
		}
		return NetworkScheme{}, false, err
	}
	return findedNetwork, true, nil
}

//	func (c Network) UpdateNetwork(Network NetworkScheme) error {
//		_, err := c.ct.UpdateOne(context.TODO(), bson.M{"_id": Network.ID}, bson.M{"$set": Network})
//		return err
//	}
//
// GetNetworkByIdentifier
func (c Network) GetNetworkByIdentifier(identifier string) (NetworkScheme, bool, error) {
	var findedNetwork NetworkScheme
	err := c.ct.FindOne(context.TODO(), bson.M{"identifier": identifier}).Decode(&findedNetwork)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return NetworkScheme{}, false, nil
		}
		return NetworkScheme{}, false, err
	}
	return findedNetwork, true, nil
}
func (c Network) CreateNetwork(name string, identifier string, botLink string, botType string, creator primitive.ObjectID) (NetworkScheme, error) {
	id := primitive.NewObjectID()
	Network := NetworkScheme{
		ID:         id.Hex(),
		Name:       name,
		Identifier: identifier,
		BotLink:    botLink,
		BotType:    botType,
		Creator:    creator,
	}
	_, err := c.ct.InsertOne(context.TODO(), Network)
	if err != nil {
		return NetworkScheme{}, err
	}
	return Network, nil
}
func (c Network) UpdateNetwork(networkIdentifier string, data map[string]string) (NetworkScheme, bool, error) {
	_, err := c.ct.UpdateOne(context.TODO(), bson.M{"identifier": networkIdentifier}, bson.M{"$set": data})
	if err != nil {
		return NetworkScheme{}, false, err
	}
	network, isExist, err := c.GetNetworkByIdentifier(networkIdentifier)
	return network, isExist, err
}

// GetAllNetworks
func (c Network) GetAllNetworks() ([]NetworkScheme, error) {
	var findedNetworks []NetworkScheme
	cur, err := c.ct.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var network NetworkScheme
		err := cur.Decode(&network)
		if err != nil {
			return nil, err
		}
		findedNetworks = append(findedNetworks, network)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return findedNetworks, nil
}

// func (c Network) GetNetworksByArrayID(Networks []primitive.ObjectID) ([]NetworkScheme, error) {
// 	var findedNetworks []NetworkScheme
// 	cur, err := c.ct.Find(context.TODO(), bson.M{"_id": bson.M{"$in": Networks}})
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer cur.Close(context.TODO())
// 	for cur.Next(context.TODO()) {
// 		var Network NetworkScheme
// 		err := cur.Decode(&Network)
// 		if err != nil {
// 			return nil, err
// 		}
// 		findedNetworks = append(findedNetworks, Network)
// 	}
// 	if err := cur.Err(); err != nil {
// 		return nil, err
// 	}
// 	return findedNetworks, nil
// }
