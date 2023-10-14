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
