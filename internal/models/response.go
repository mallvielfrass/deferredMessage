package models

type ErrorResponse struct {
	Error  string `json:"error"`
	Reason string `json:"reason"`
}
type SuccessResponse struct {
	Status string `json:"status"`
}
type PingMessageResponse struct {
	Message string `json:"message" example:"pong"`
}
