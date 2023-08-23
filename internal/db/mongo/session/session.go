package session

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

type SessionScheme struct {
	ID     string `bson:"_id"`
	UserID string `bson:"user_id"`
	Expire int64  `bson:"expire"`
	IP     string `bson:"ip"`
}
type Session struct {
	ct *mongo.Collection
}

func Init(driver *mongo.Database) Session {
	collectionName := "session"
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

	return Session{
		ct: driver.Collection(collectionName),
	}
}

// find by ID
func (Session Session) GetSessionByID(id primitive.ObjectID) (SessionScheme, bool, error) {
	var findedSession SessionScheme
	err := Session.ct.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&findedSession)
	fmt.Printf("findedSession: %#v\n", findedSession)
	fmt.Printf("err: %#v\n", err)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return SessionScheme{}, false, nil
		}
		return SessionScheme{}, false, err
	}
	return findedSession, true, nil
}

// create Session (name, hash)
func (Session Session) CreateSession(UserID string, expire int64, ip string) (SessionScheme, error) {

	res, err := Session.ct.InsertOne(context.TODO(), bson.M{"user_id": UserID, "expire": expire, "ip": ip})
	if err != nil {
		return SessionScheme{}, err
	}
	u, isExist, err := Session.GetSessionByID(res.InsertedID.(primitive.ObjectID))
	if err != nil {
		return SessionScheme{}, err
	}
	if !isExist {
		return SessionScheme{}, fmt.Errorf("Session not created")
	}
	// fmt.Println(res)
	// fmt.Println(u)
	return u, nil
}
