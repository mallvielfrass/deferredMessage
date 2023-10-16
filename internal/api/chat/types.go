package chat

import "deferredMessage/internal/models"

type ChatRequest struct {
	Name          string `json:"name"`
	BotID         string `json:"botId"`
	LinkOrIdInBot string `json:"linkOrIdInBot"`
}

type createdChatType struct {
	Linker        string `json:"linker"`
	Name          string `json:"name"`
	BotIdentifier string `json:"botIdentifier"`
	BotID         string `json:"botId"`
	Verified      bool   `json:"verified"`
	Id            string `json:"_id"`
}
type CreateChatResponse struct {
	Chat createdChatType `json:"chat"`
}
type ChatsListResponse struct {
	Chats  []models.ChatScheme `json:"chats"`
	Offset int                 `json:"offset"`
	Count  int                 `json:"count"`
}
