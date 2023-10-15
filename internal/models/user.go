package models

type AuthUser struct {
	Mail     string
	Password string
	IP       string
}
type Session struct {
	ID     string
	Expire int64
}

type RegisterUser struct {
	Name     string
	Mail     string
	Password string
	IP       string
}
type UserIdentify struct {
	ID    string
	Name  string
	Mail  string
	Admin bool
}
type UserScheme struct {
	Name  string `bson:"name"`
	Mail  string `bson:"mail"`
	Hash  string `bson:"hash"`
	ID    string `bson:"_id"`
	Admin bool   `bson:"admin"`
	//Chats array with id of chats refferenced to ChatScheme
	Chats []string `bson:"chats"`
}
