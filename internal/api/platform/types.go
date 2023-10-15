package platform

type PlatformResponse struct {
	Name string `json:"name"`
	// ID   string `json:"id"`
}
type CreatePlatformRequest struct {
	Name string `json:"name" binding:"required"`
}
