package user

import (
	"deferredMessage/internal/middleware"
	"deferredMessage/internal/models"
	"deferredMessage/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userApi struct {
	services   *service.Service
	middleware *middleware.Middleware
}

func Init(services *service.Service, middleware *middleware.Middleware) userApi {
	return userApi{
		services:   services,
		middleware: middleware,
	}
}

// check auth middleware

// @Summary Ping
// @Description Returns a "pong" message
// @Accept json
// @Tags user
// @Security Bearer
// @Produce json
// @Success 200 {object}  models.PingMessageResponse
// @Router /api/auth/user/ping [get]
func ping(c *gin.Context) {
	c.JSON(http.StatusOK, models.PingMessageResponse{
		Message: "pong",
	})
}

func (n userApi) Router(router *gin.RouterGroup) *gin.RouterGroup {
	r := router.Group("/")

	r.GET("/ping", ping)

	return r
}
