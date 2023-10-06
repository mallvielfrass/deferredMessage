package middleware

import (
	"deferredMessage/internal/db"
	"deferredMessage/internal/db/mongo/session"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Middleware struct {
	db db.DB
}

func InitMiddleware(db db.DB) Middleware {
	return Middleware{
		db: db,
	}
}
func (n Middleware) checkSession(c *gin.Context) (*gin.Context, interface{}, bool) {
	//get token from header
	token := c.GetHeader("Authorization")
	if token == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no token"})
		return c, nil, false
	}
	//validate and convert token to primitive.ObjectID
	sessionID, err := primitive.ObjectIDFromHex(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return c, nil, false
	}

	session, exist, err := n.db.Collections.Session.GetSessionByID(sessionID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return c, nil, false
	}
	if !exist {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return c, nil, false
	}

	c.Set("session", session)

	if session.Expire < time.Now().Unix() {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "session expired"})
		return c, nil, false
	}
	if !session.Valid {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
		return c, nil, false
	}
	return c, session, true
}
func (n Middleware) CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, _, ok := n.checkSession(c)
		c = ctx
		if !ok {
			if !c.IsAborted() {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid session"})
			}
			return
		}
		c.Next()
	}
}
func (n Middleware) CheckAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionFromCtx, isExist := c.Get("session")
		if !isExist {
			ctx, session, ok := n.checkSession(c)
			c = ctx
			sessionFromCtx = session
			if !ok {
				if !c.IsAborted() {
					c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid session"})
				}
				return
			}
		}

		session, ok := sessionFromCtx.(session.SessionScheme)

		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
			return
		}
		user, isExist, err := n.db.Collections.User.GetUserByID(session.UserID)
		fmt.Printf("user: %#v\n", user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !isExist {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
			return
		}

		if !user.Admin {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "you are not admin"})
			return
		}
		c.Next()
	}
}
