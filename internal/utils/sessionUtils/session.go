package sessionutils

import (
	"deferredMessage/internal/db/mongo/session"
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetSession(c *gin.Context) (session.SessionScheme, error) {
	sessionFromCtx, ok := c.Get("session")
	if !ok {
		return session.SessionScheme{}, fmt.Errorf("no session")
	}
	switch sessionTyped := sessionFromCtx.(type) {
	case session.SessionScheme:
		return sessionTyped, nil
	default:
		return session.SessionScheme{}, fmt.Errorf("invalid session type")
	}

}
