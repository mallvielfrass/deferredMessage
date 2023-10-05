package user

import (
	"deferredMessage/internal/db"
	"deferredMessage/internal/db/mongo/session"
	"deferredMessage/internal/utils"
	"deferredMessage/internal/utils/dto"
	"fmt"
	"net/http"
	"os"
	"strconv"
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

func (n userApi) Router(router *gin.RouterGroup) *gin.RouterGroup {
	r := router.Group("/")
	r.Use(n.CheckAuth())
	r.POST("/admin", func(c *gin.Context) {
		userSession, ok := c.Get("session")
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no session"})
			return
		}
		session := userSession.(session.SessionScheme)
		fmt.Println(session)
		userID, err := primitive.ObjectIDFromHex(session.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
			return
		}
		_, exist, err := n.db.Collections.User.GetUserByID(userID)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !exist {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		}
		body, exist := dto.GetStruct[EncryptedData](c, EncryptedData{})
		if !exist {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no body"})
			return
		}
		if body.Token == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no token"})
			return
		}
		ADMIN_KEY, exist := os.LookupEnv("ADMIN_KEY")
		if !exist {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ADMIN_KEY not found"})
			return
		}

		if body.Token != ADMIN_KEY {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
			return
		}
		settedUser, isSetAdmin, err := n.db.Collections.User.SetUserAdmin(userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": gin.H{"name": settedUser.Name,
			"mail": settedUser.Mail, "admin": settedUser.Admin, "id": settedUser.ID}, "isAdmin": isSetAdmin})
	})
	r.GET("/ping", ping)

	r.GET("/networks", func(c *gin.Context) {
		networks, err := n.db.Collections.Network.GetAllNetworks()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var networksResponse []Network
		for _, network := range networks {
			networksResponse = append(networksResponse, Network{Name: network.Name, Identifier: network.Identifier})
		}

		c.JSON(http.StatusOK, gin.H{"networks": networksResponse})
	})
	r.POST("/networks", func(c *gin.Context) {
		body, exist := dto.GetStruct[Network](c, Network{})
		if !exist {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no body"})
			return
		}
		network, err := n.db.Collections.Network.CreateNetwork(body.Name, body.Identifier)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"network": gin.H{"name": network.Name, "identifier": network.Identifier}})
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
		userID, err := primitive.ObjectIDFromHex(session.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
			return
		}
		user, exist, err := n.db.Collections.User.GetUserByID(userID)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !exist {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user not found"})
			return
		}
		//get from url query params [count, offset]
		countString := c.DefaultQuery("count", "10")
		offsetString := c.DefaultQuery("offset", "0")
		count, err := strconv.Atoi(countString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		offset, err := strconv.Atoi(offsetString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		diff := len(user.Chats) - offset
		if count > diff {
			count = diff
		}
		if count < 0 {
			c.JSON(http.StatusOK, gin.H{"chats": []gin.H{}})
			return
		}
		userChatsId := user.Chats
		chatsId := userChatsId[offset : offset+count]
		if len(chatsId) == 0 {
			//empty array
			c.JSON(http.StatusOK, gin.H{"chats": []gin.H{}})
			return
		}
		chats, err := n.db.Collections.Chat.GetChatsByArrayID(chatsId)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var respArray []Chat
		for _, chat := range chats {
			respArray = append(respArray, Chat{
				ID:                chat.ID.Hex(),
				Name:              chat.Name,
				NetworkIdentifier: chat.NetworkIdentifier,
				NetworkID:         chat.NetworkID,
				LinkOrIdInNetwork: chat.LinkOrIdInNetwork,
				Verified:          chat.Verified,
			})
		}
		c.JSON(http.StatusOK, gin.H{"chats": respArray, "offset": offset, "count": count})
	})
	userChats.POST("/", func(c *gin.Context) {
		body, exist := dto.GetStruct[Chat](c, Chat{})
		if !exist {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no body"})
			return
		}
		network, isExist, err := n.db.Collections.Network.GetNetworkByIdentifier(body.NetworkIdentifier)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !isExist {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "network not found"})
			return
		}
		userSession, ok := c.Get("session")
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no session"})
			return
		}
		session := userSession.(session.SessionScheme)
		fmt.Println(session)
		userID, err := primitive.ObjectIDFromHex(session.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
			return
		}
		createdChat, err := n.db.Collections.Chat.CreateChat(body.Name, network.Identifier, network.ID, userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = n.db.Collections.User.AddChatToUser(createdChat.ID, userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		linker, err := utils.Encrypt(createdChat.ID.Hex())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"chat": gin.H{"linker": linker, "id": createdChat.ID, "name": createdChat.Name, "networkIdentifier": createdChat.NetworkIdentifier, "networkId": createdChat.NetworkID, "verified": createdChat.Verified}})
	})
	userChats.PUT("/:id", func(c *gin.Context) {
		var body map[string]interface{}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		}
		params := make(map[string]string)
		//check name in body
		if body["name"] != nil {
			switch body["name"].(type) {
			case string:
				params["name"] = body["name"].(string)
			}

		}
		//linkOrIdInNetwork
		if body["linkOrIdInNetwork"] != nil {
			switch body["linkOrIdInNetwork"].(type) {
			case string:
				params["linkOrIdInNetwork"] = body["linkOrIdInNetwork"].(string)
			}
		}
		fmt.Println(params)

		//user
		userSession, ok := c.Get("session")
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no session"})
			return
		}
		session := userSession.(session.SessionScheme)
		fmt.Println(session)
		userID, err := primitive.ObjectIDFromHex(session.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid token"})
			return
		}
		chatId := c.Param("id")
		fmt.Printf("chatId: %v\n", chatId)
		chatObjectID, err := primitive.ObjectIDFromHex(chatId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid chat id"})
			return
		}
		chat, exist, err := n.db.Collections.Chat.GetChatByID(chatObjectID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if !exist {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "chat not found"})
			return
		}
		if chat.Creator != userID {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "you are not creator"})
			return
		}
		err = n.db.Collections.Chat.UpdateChat(chatObjectID, params)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})

	})
	return r
}
