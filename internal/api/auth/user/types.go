package user

type Chat struct {
	Name             string `json:"name"`
	NetworkIdentifer string `json:"networkIdentifer"`
}
type EncryptedData struct {
	Token string `json:"token"`
}
type Network struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
}
