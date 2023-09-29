package user

import (
	"deferredMessage/internal/db"
	"deferredMessage/internal/db/mongo/session"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userApi struct {
	db db.DB
}

func Init(db db.DB) userApi {
	return userApi{
		db: db,
	}
}

// check auth middleware

func (n userApi) CheckAuth() gin.HandlerFunc {
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
		fmt.Println(sessionID)

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

// @Summary Ping
// @Description Returns a "pong" message
// @Accept json
// @Tags user
// @Security Bearer
// @Produce json
// @Success 200 {object}  MessageResponse
// @Router /api/auth/user/ping [get]
func ping(c *gin.Context) {
	c.JSON(http.StatusOK, MessageResponse{
		Message: "pong",
	})
}

type Network struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
}

func (n userApi) Router(router *gin.RouterGroup) *gin.RouterGroup {
	r := router.Group("/")
	r.Use(n.CheckAuth())

	r.GET("/ping", ping)
	r.GET("/networks", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"networks": []Network{
			{Name: "Telegram official bot", Identifier: "telegram_official_bot"},
		}})
	})
	userChats := r.Group("/chats")
	userChats.GET("/", func(c *gin.Context) {
		userSession, ok := c.Get("session")
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no session"})
			return
		}
		session := userSession.(session.SessionScheme)
		fmt.Println(session)
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	return r
}
