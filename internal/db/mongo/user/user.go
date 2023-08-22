package user

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

type UserScheme struct {
	Name string `bson:"name"`
	Mail string `bson:"mail"`
	Hash string `bson:"hash"`
	ID   string `bson:"_id"`
}
type User struct {
	ct *mongo.Collection
}

func Init(driver *mongo.Database) User {
	collectionName := "user"
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

	return User{
		ct: driver.Collection(collectionName),
	}
}
func (user User) CheckUser(mail string) (bool, error) {
	var findedUser UserScheme
	err := user.ct.FindOne(context.TODO(), bson.M{"mail": mail}).Decode(&findedUser)
	fmt.Printf("findedUser: %#v\n", findedUser)
	fmt.Printf("err: %#v\n", err)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// find by ID
func (user User) GetUser(id primitive.ObjectID) (UserScheme, error) {
	var findedUser UserScheme
	err := user.ct.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&findedUser)
	fmt.Printf("findedUser: %#v\n", findedUser)
	fmt.Printf("err: %#v\n", err)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return UserScheme{}, nil
		}
		return UserScheme{}, err
	}
	return findedUser, nil
}

// create User (name, hash)
func (user User) CreateUser(name, mail, hash string) (UserScheme, error) {
	fmt.Printf("name: %s, mail: %s, hash: %s\n", name, mail, hash)
	res, err := user.ct.InsertOne(context.TODO(), bson.M{"name": name, "mail": mail, "hash": hash})
	if err != nil {
		return UserScheme{}, err
	}
	u, err := user.GetUser(res.InsertedID.(primitive.ObjectID))
	if err != nil {
		return UserScheme{}, err
	}
	fmt.Println(res)
	fmt.Println(u)
	return u, nil
}
