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
type CheckUserRequest struct {
	Mail string `json:"mail"`
}
type StatusResponse struct {
	Status string `json:"status"`
}
type User struct {
	Name string `json:"name"`
	Mail string `json:"mail"`
}
type Session struct {
	Id     string `json:"_id"`
	Expire int64  `json:"expire"`
}
type RegisterUserResponse struct {
	Status  string  `json:"status"`
	User    User    `json:"user"`
	Session Session `json:"session"`
}
