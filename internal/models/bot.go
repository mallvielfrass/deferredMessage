package models

type BotScheme struct {
	ID          string `bson:"_id"`
	Name        string `bson:"name"`
	BotLink     string `bson:"botLink"`
	BotType     string `bson:"botType"`
	Creator     string `bson:"creator"`
	Platform    string `bson:"platform"`
	HashedToken string `bson:"hashedToken"`
}
