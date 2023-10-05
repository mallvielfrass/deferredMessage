package user

type Chat struct {
	ID                string `json:"_id"`
	Name              string `json:"name"`
	NetworkIdentifier string `json:"networkIdentifier"`
	NetworkID         string `json:"networkId"`
	LinkOrIdInNetwork string `json:"linkOrIdInNetwork"`
	Verified          bool   `json:"verified"`
}
type EncryptedData struct {
	Token string `json:"token"`
}
type Network struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
}
