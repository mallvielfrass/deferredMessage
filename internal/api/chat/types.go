package chat

type Chat struct {
	ID            string `json:"_id"`
	Name          string `json:"name"`
	BotIdentifier string `json:"botIdentifier"`
	BotID         string `json:"botId"`
	LinkOrIdInBot string `json:"linkOrIdInBot"`
	Verified      bool   `json:"verified"`
}
