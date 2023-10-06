package middleware

import (
	"deferredMessage/internal/db"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SessionMiddleware struct {
	db db.DB
}

func InitSessionMiddleware(db db.DB) SessionMiddleware {
	return SessionMiddleware{
		db: db,
	}
}
func (n SessionMiddleware) CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//get token from header
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no token"})
			return
		}
		//validate and convert token to primitive.ObjectID
		sessionID, err := primitive.ObjectIDFromHex(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
			return
		}

		session, exist, err := n.db.Collections.Session.GetSessionByID(sessionID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !exist {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
			return
		}

		c.Set("session", session)
		//	fmt.Printf("session: %v\n", session)
		if session.Expire < time.Now().Unix() {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "session expired"})
			return
		}
		if !session.Valid {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
			return
		}
		c.Next()
	}
}
