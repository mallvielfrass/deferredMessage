package sessionutils

import (
	"deferredMessage/internal/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetSession(c *gin.Context) (models.SessionScheme, error) {
	sessionFromCtx, ok := c.Get("session")
	if !ok {
		return models.SessionScheme{}, fmt.Errorf("no session")
	}
	switch sessionTyped := sessionFromCtx.(type) {
	case models.SessionScheme:
		return sessionTyped, nil
	default:
		return models.SessionScheme{}, fmt.Errorf("invalid session type")
	}

}
