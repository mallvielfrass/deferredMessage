package bot

type MessageResponse struct {
	Message string `json:"message" example:"pong"`
}

//	type Bot struct {
//		Name     string             `json:"name" binding:"required"`
//		Id       string             `json:"_id"`
//		BotLink  string             `json:"botLink"`
//		Creator  primitive.ObjectID `json:"creator"`
//		Platform string             `json:"platform"`
//	}
type BotResponse struct {
	Name     string `json:"name" binding:"required"`
	Id       string `json:"_id"`
	BotLink  string `json:"botLink"`
	Creator  string `json:"creator"`
	Platform string `json:"platform"  binding:"required"`
	Token    string `json:"token"`
}
type BotRequest struct {
	Name string `json:"name" binding:"required"`

	BotLink string `json:"botLink"`

	Platform string `json:"platform"  binding:"required"`
	Token    string `json:"token"`
}
type BotStructResponse struct {
	Bot BotResponse `json:"bot"`
}
type BotStructArrayResponse struct {
	Bots []BotResponse `json:"bots"`
}

type BotUpdateRequest struct {
	Name     string `json:"name"  `
	Platform string `json:"platform"  `
	Token    string `json:"token"  `
}
