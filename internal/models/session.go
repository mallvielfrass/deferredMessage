package models

type SessionScheme struct {
	ID        string `bson:"_id"`
	UserID    string `bson:"user_id"`
	Expire    int64  `bson:"expire"`
	IP        string `bson:"ip"`
	Valid     bool   `bson:"valid"`
	AtCreated int64  `bson:"at_created"`
}
