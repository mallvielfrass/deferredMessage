package user

type Chat struct {
	ID            string `json:"_id"`
	Name          string `json:"name"`
	BotIdentifier string `json:"botIdentifier"`
	BotID         string `json:"botId"`
	LinkOrIdInBot string `json:"linkOrIdInBot"`
	Verified      bool   `json:"verified"`
}

type Bot struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
	BotLink    string `json:"botLink"`
	BotType    string `json:"botType"`
}
