package user

import (
	db "deferredMessage/internal/repository"
	"deferredMessage/internal/repository/mongo/session"
	"fmt"
	"net/http"
	"strconv"

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

	r.GET("/ping", ping)
	r.GET("/platforms/list", func(c *gin.Context) {
		botTypes, err := n.db.Collections.Platform.GetAllPlatforms()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"botTypes": botTypes})
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

		user, exist, err := n.db.Collections.User.GetUserByID(session.UserID)

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
				ID:            chat.ID.Hex(),
				Name:          chat.Name,
				BotIdentifier: chat.BotIdentifier,
				BotID:         chat.BotID,
				LinkOrIdInBot: chat.LinkOrIdInBot,
				Verified:      chat.Verified,
			})
		}
		c.JSON(http.StatusOK, gin.H{"chats": respArray, "offset": offset, "count": count})
	})
	// userChats.POST("/", func(c *gin.Context) {
	// 	body, exist := dto.GetStruct[Chat](c, Chat{})
	// 	if !exist {
	// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no body"})
	// 		return
	// 	}
	// 	bot, isExist, err := n.db.Collections.Bot.GetBotByIdentifier(body.BotIdentifier)

	// 	if err != nil {
	// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// 	if !isExist {
	// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bot not found"})
	// 		return
	// 	}
	// 	userSession, ok := c.Get("session")
	// 	if !ok {
	// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no session"})
	// 		return
	// 	}
	// 	session := userSession.(session.SessionScheme)
	// 	fmt.Println(session)

	// 	createdChat, err := n.db.Collections.Chat.CreateChat(body.Name, bot.Id , bot.ID, session.UserID)
	// 	if err != nil {
	// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	err = n.db.Collections.User.AddChatToUser(createdChat.ID, session.UserID)
	// 	if err != nil {
	// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// 	linker, err := utils.Encrypt(createdChat.ID.Hex())
	// 	if err != nil {
	// 		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	// 		return
	// 	}

	// 	c.JSON(http.StatusOK, gin.H{"chat": gin.H{"linker": linker, "id": createdChat.ID, "name": createdChat.Name, "botIdentifier": createdChat.BotIdentifier, "botId": createdChat.BotID, "verified": createdChat.Verified}})
	// })
	userChats.PUT("/:id", func(c *gin.Context) {
		var body map[string]interface{}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		}
		params := make(map[string]interface{})
		//check name in body
		if body["name"] != nil {
			switch body["name"].(type) {
			case string:
				params["name"] = body["name"].(string)
			}

		}
		//linkOrIdInBot
		if body["linkOrIdInBot"] != nil {
			switch body["linkOrIdInBot"].(type) {
			case string:
				params["linkOrIdInBot"] = body["linkOrIdInBot"].(string)
			}
		}
		//fmt.Println(params)

		//user
		userSession, ok := c.Get("session")
		if !ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no session"})
			return
		}
		session := userSession.(session.SessionScheme)
		fmt.Println(session)

		chatId := c.Param("id")
		//	fmt.Printf("chatId: %v\n", chatId)
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
		if chat.Creator != session.UserID {
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
