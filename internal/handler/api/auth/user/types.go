package user

type Bot struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
	BotLink    string `json:"botLink"`
	BotType    string `json:"botType"`
}
