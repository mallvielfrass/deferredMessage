package middleware

import (
	"deferredMessage/internal/models"
	"deferredMessage/internal/service"
	sessionutils "deferredMessage/internal/utils/sessionUtils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	service *service.Service
}

func InitMiddleware(service *service.Service) Middleware {
	return Middleware{
		service: service,
	}
}
func (n Middleware) checkSession(c *gin.Context) (*gin.Context, models.SessionScheme, bool) {
	//get token from header
	token := c.GetHeader("Authorization")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "no token",
		})
		return c, models.SessionScheme{}, false
	}
	session, err := n.service.SessionService.CheckSession(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return c, models.SessionScheme{}, false
	}
	c.Set("session", session)
	return c, session, true
}
func (n Middleware) CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, _, ok := n.checkSession(c)
		c = ctx
		if !ok {
			if !c.IsAborted() {
				c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
					Error: "invalid session"})
			}
			return
		}
		c.Next()
	}
}
func (n Middleware) CheckAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := sessionutils.GetSession(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
				Error: err.Error()})
			return
		}
		userIsAdmin, err := n.service.UserService.UserIsAdmin(session.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
				Error: err.Error(),
			})
		}
		if !userIsAdmin {
			c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "you are not admin",
			})
			return
		}

		c.Next()
	}
}
