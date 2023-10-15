package user

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

type UserScheme struct {
	Name  string             `bson:"name"`
	Mail  string             `bson:"mail"`
	Hash  string             `bson:"hash"`
	ID    primitive.ObjectID `bson:"_id"`
	Admin bool               `bson:"admin"`
	//Chats array with id of chats refferenced to ChatScheme
	Chats []primitive.ObjectID `bson:"chats"`
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
func (user User) CheckUserByMail(mail string) (bool, error) {
	var findedUser UserScheme
	err := user.ct.FindOne(context.TODO(), bson.M{"mail": mail}).Decode(&findedUser)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// find by ID
func (user User) GetUserByID(id string) (models.UserScheme, bool, error) {
	var findedUser UserScheme
	err := user.ct.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&findedUser)
	fmt.Printf("findedUser: %#v\n", findedUser)
	fmt.Printf("err: %#v\n", err)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.UserScheme{}, false, nil
		}
		return models.UserScheme{}, false, err
	}
	var chatsStr []string
	for _, v := range findedUser.Chats {
		chatsStr = append(chatsStr, v.Hex())
	}
	return models.UserScheme{
		Name:  findedUser.Name,
		Mail:  findedUser.Mail,
		Hash:  findedUser.Hash,
		ID:    findedUser.ID.Hex(),
		Admin: findedUser.Admin,
		Chats: chatsStr,
	}, true, nil
}
func (user User) getUserById(id primitive.ObjectID) (UserScheme, bool, error) {
	var findedUser UserScheme
	err := user.ct.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&findedUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return UserScheme{}, false, nil
		}
		return UserScheme{}, false, err
	}
	return findedUser, true, nil
}

// find by ID
func (user User) GetUserByMail(mail string) (models.UserScheme, bool, error) {
	var findedUser UserScheme
	err := user.ct.FindOne(context.TODO(), bson.M{"mail": mail}).Decode(&findedUser)
	//fmt.Printf("findedUser: %#v\n", findedUser)
	fmt.Printf("err: %#v\n", err)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.UserScheme{}, false, nil
		}
		return models.UserScheme{}, false, err
	}
	return models.UserScheme{
		Name:  findedUser.Name,
		Mail:  findedUser.Mail,
		Hash:  findedUser.Hash,
		ID:    findedUser.ID.Hex(),
		Admin: findedUser.Admin,
	}, true, nil
}

// create User (name, hash)
func (user User) CreateUser(name, mail, hash string) (models.UserScheme, error) {
	fmt.Printf("name: %s, mail: %s, hash: %s\n", name, mail, hash)
	res, err := user.ct.InsertOne(context.TODO(), bson.M{"name": name, "mail": mail, "hash": hash})
	if err != nil {
		return models.UserScheme{}, err
	}

	u, isExist, err := user.getUserById(res.InsertedID.(primitive.ObjectID))
	if err != nil {
		return models.UserScheme{}, err
	}
	if !isExist {
		return models.UserScheme{}, fmt.Errorf("User not created")
	}

	return models.UserScheme{
		Name:  u.Name,
		Mail:  u.Mail,
		Hash:  u.Hash,
		ID:    u.ID.Hex(),
		Admin: u.Admin,
	}, nil
}

// add chat to user
func (user User) AddChatToUser(chatID string, userID string) error {
	chatObjectID, err := primitive.ObjectIDFromHex(chatID)
	if err != nil {
		return err
	}
	useObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	res, err := user.ct.UpdateOne(context.TODO(), bson.M{"_id": useObjectID}, bson.M{"$push": bson.M{"chats": chatObjectID}})
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return fmt.Errorf("chat not added")
	}
	return nil
}

// SetUserAdmin(id primitive.ObjectID) (user.UserScheme, bool, error)
func (user User) SetUserAdmin(userID string) (models.UserScheme, bool, error) {
	useObjectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return models.UserScheme{}, false, err
	}
	_, err = user.ct.UpdateOne(context.TODO(), bson.M{"_id": useObjectID}, bson.M{"$set": bson.M{"admin": true}})
	if err != nil {
		return models.UserScheme{}, false, err
	}

	u, isExist, err := user.getUserById(useObjectID)
	if err != nil {
		return models.UserScheme{}, false, err
	}
	if !isExist {
		return models.UserScheme{}, false, fmt.Errorf("User not found")
	}

	return models.UserScheme{
		Name:  u.Name,
		Mail:  u.Mail,
		Hash:  u.Hash,
		ID:    u.ID.Hex(),
		Admin: u.Admin,
	}, true, nil
}
