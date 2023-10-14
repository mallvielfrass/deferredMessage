package admin

type EncryptedData struct {
	Token string `json:"token"`
}
type UserResponse struct {
	Name  string `json:"name"`
	Mail  string `json:"mail"`
	Admin bool   `json:"admin"`
	ID    string `json:"id"`
}
type AdminResponse struct {
	User    UserResponse `json:"user"`
	IsAdmin bool         `json:"isAdmin"`
}
