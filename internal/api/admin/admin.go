package admin

import (
	db "deferredMessage/internal/repository"

	"deferredMessage/internal/models"
	"deferredMessage/internal/utils/dto"
	sessionutils "deferredMessage/internal/utils/sessionUtils"
	"deferredMessage/internal/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Admin struct {
	db db.DB
}

func Init(db db.DB) Admin {
	return Admin{
		db: db,
	}
}

// HandleSetAdmin
func (a Admin) HandleSetAdmin(c *gin.Context) {
	session, err := sessionutils.GetSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error()})
		return
	}

	body, exist := dto.GetStruct[EncryptedData](c, EncryptedData{})
	if !exist {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{Error: "no body"})
		return
	}
	err = token.ParseAndCheckToken(body.Token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	settedUser, isSetAdmin, err := a.db.Collections.User.SetUserAdmin(session.UserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, AdminResponse{
		User: UserResponse{
			Name:  settedUser.Name,
			Mail:  settedUser.Mail,
			Admin: settedUser.Admin,
			ID:    settedUser.ID.Hex(),
		},
		IsAdmin: isSetAdmin,
	})
}

func (a Admin) Router(router *gin.RouterGroup) *gin.RouterGroup {
	r := router.Group("/")

	r.POST("/", a.HandleSetAdmin)
	return r

}
