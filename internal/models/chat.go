package models

type ChatScheme struct {
	Name          string `bson:"name"`
	ID            string `bson:"_id"`
	LinkOrIdInBot string `bson:"linkOrIdInBot"`
	BotIdentifier string `bson:"botIdentifier"`
	BotID         string `bson:"botId"`
	Verified      bool   `bson:"verified"`
	Creator       string `bson:"creator"`
	Hidden        bool   `bson:"hidden"`
}
