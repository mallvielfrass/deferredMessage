package session

import (
	"context"
	"deferredMessage/internal/models"
	"fmt"
	"os"
	"time"

	"github.com/mallvielfrass/fmc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slices"
)

type SessionScheme struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Expire    int64              `bson:"expire"`
	IP        string             `bson:"ip"`
	Valid     bool               `bson:"valid"`
	AtCreated int64              `bson:"at_created"`
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
func (Session Session) GetSessionByID(id string) (models.SessionScheme, bool, error) {

	idObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.SessionScheme{}, false, err
	}
	var findedSession SessionScheme
	err = Session.ct.FindOne(context.TODO(), bson.M{"_id": idObjectID}).Decode(&findedSession)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.SessionScheme{}, false, nil
		}
		return models.SessionScheme{}, false, err
	}
	return models.SessionScheme{
		ID:        findedSession.ID.Hex(),
		UserID:    findedSession.UserID.Hex(),
		Expire:    findedSession.Expire,
		IP:        findedSession.IP,
		Valid:     findedSession.Valid,
		AtCreated: findedSession.AtCreated,
	}, true, nil
}

// create Session (name, hash)
func (Session Session) CreateSession(UserID string, expire int64, ip string) (models.SessionScheme, error) {
	userObjectID, err := primitive.ObjectIDFromHex(UserID)
	if err != nil {
		return models.SessionScheme{}, err
	}
	res, err := Session.ct.InsertOne(context.TODO(), bson.M{"user_id": userObjectID, "expire": expire, "ip": ip, "valid": true, "at_created": time.Now().Unix()})
	if err != nil {
		return models.SessionScheme{}, err
	}
	u, isExist, err := Session.GetSessionByID(res.InsertedID.(primitive.ObjectID).Hex())
	if err != nil {
		return models.SessionScheme{}, err
	}
	if !isExist {
		return models.SessionScheme{}, fmt.Errorf("Session not created")
	}
	// fmt.Println(res)
	// fmt.Println(u)
	return u, nil
}
