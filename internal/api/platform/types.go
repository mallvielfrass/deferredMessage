package platform

type PlatformResponse struct {
	Name string `json:"name"`
	// ID   string `json:"_id"`
}
type CreatePlatformRequest struct {
	Name string `json:"name" binding:"required"`
}
