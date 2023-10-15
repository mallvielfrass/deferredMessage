package noauth

import (
	"deferredMessage/internal/middleware"
	"deferredMessage/internal/models"
	"deferredMessage/internal/service"
	"deferredMessage/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type NoAuth struct {
	services   *service.Service
	middleware *middleware.Middleware
}

func Init(services *service.Service, middleware *middleware.Middleware) NoAuth {
	return NoAuth{
		services:   services,
		middleware: middleware,
	}
}

// HandleCheckUserExist checks if a user exists.
// @Summary Check if user exists
// @Description Checks if a user exists based on the provided registration information.
// @Tags NoAuth
// @Accept json
// @Produce json
// @Param body body CheckUserRequest true "Check user request body"
// @Success 200 {object} StatusResponse  "User does not exist"
// @Success 202 {object} models.ErrorResponse "User already exists"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Router /api/noauth/check [post]
func (n NoAuth) HandleCheckUserExist(c *gin.Context) {
	var body CheckUserRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error:  "no body",
			Reason: err.Error(),
		})
		return
	}

	isExist, err := n.services.UserService.CheckUserByMail(body.Mail)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error()})
		return
	}
	if isExist {
		c.JSON(http.StatusAccepted, models.ErrorResponse{
			Error: "user already exist"})
		return
	}
	c.JSON(http.StatusOK, StatusResponse{
		Status: "not exist",
	})

}

// HandleRegisterUser creates a new user.
// @Summary Create a new user
// @Description Creates a new user with the provided registration information.
// @Tags NoAuth
// @Accept json
// @Produce json
// @Param body body RegisterBody true "Registration information"
// @Success 200 {object} RegisterUserResponse "User created successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Router /api/noauth/register [post]
func (n NoAuth) HandleRegisterUser(c *gin.Context) {
	//fmt.Println("register")
	var body RegisterBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error:  "no body",
			Reason: err.Error(),
		})
		return
	}

	isExist, err := n.services.UserService.CheckUserByMail(body.Mail)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error()})
		return
	}
	if isExist {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "user already exist"})
		return
	}
	hash, err := utils.HashPassword(body.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error()})
		return
	}
	user, err := n.services.UserService.CreateUser(body.Name, body.Mail, hash)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error()})
		return
	}
	userIP := c.ClientIP()
	session, err := n.services.SessionService.CreateSession(user.ID, time.Now().Add(time.Hour*24*31).Unix(), userIP)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, RegisterUserResponse{
		Status: "success",
		User: User{
			Name: user.Name,
			Mail: user.Mail,
		},
		Session: Session{
			Id:     session.ID,
			Expire: session.Expire,
		},
	})
}

// HandleLoginUser logs in a user.
// @Summary Log in a user
// @Description Logs in a user with the provided login information.
// @Tags NoAuth
// @Accept json
// @Produce json
// @Param body body LoginBody true "Login information"
// @Success 200 {object} RegisterUserResponse "User logged in successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Router /api/noauth/login [post]
func (n NoAuth) HandleLoginUser(c *gin.Context) {
	var body LoginBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error:  "no body",
			Reason: err.Error(),
		})
		return
	}

	user, isExist, err := n.services.UserService.GetUserByMail(body.Mail)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "user or password incorrect"})

		return
	}
	if !isExist {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "user or password incorrect"})
		return
	}
	if !utils.CheckPasswordHash(body.Password, user.Hash) {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "user or password incorrect"})
		return
	}
	userIP := c.ClientIP()
	session, err := n.services.SessionService.CreateSession(user.ID, time.Now().Add(time.Hour*24*31).Unix(), userIP)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, RegisterUserResponse{
		Status: "success",
		User: User{
			Name: user.Name,
			Mail: user.Mail,
		},
		Session: Session{
			Id:     session.ID,
			Expire: session.Expire,
		},
	})
}
func (n NoAuth) Router(router *gin.RouterGroup) *gin.RouterGroup {
	r := router.Group("/")
	r.POST("/check", n.HandleCheckUserExist)
	r.POST("/register", n.HandleRegisterUser)
	r.POST("/login", n.HandleLoginUser)
	return r
}
