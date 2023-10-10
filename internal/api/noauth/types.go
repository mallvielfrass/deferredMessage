package noauth

type RegisterBody struct {
	// json tag to de-serialize json body
	Name     string `json:"name" binding:"required"`
	Mail     string `json:"mail" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type LoginBody struct {
	// json tag to de-serialize json body

	Mail     string `json:"mail" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type StatusResponse struct {
	Status string `json:"status"`
}
