package chat

import (
	"deferredMessage/internal/middleware"
	"deferredMessage/internal/service"
	"deferredMessage/internal/utils"
	"deferredMessage/internal/utils/dto"
	sessionutils "deferredMessage/internal/utils/sessionUtils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type chatApi struct {
	services   *service.Service
	middleware *middleware.Middleware
}

func Init(services *service.Service, middleware *middleware.Middleware) chatApi {
	return chatApi{
		services:   services,
		middleware: middleware,
	}
}

// HandleCreateChat
func (n chatApi) HandleCreateChat(c *gin.Context) {
	body, exist := dto.GetStruct[Chat](c, Chat{})
	if !exist {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no body"})
		return
	}
	bot, isExist, err := n.services.BotService.GetBotByID(body.BotID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !isExist {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bot not found"})
		return
	}
	session, err := sessionutils.GetSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdChat, err := n.services.ChatService.CreateChat(body.Name, bot.ID, bot.ID, session.UserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = n.services.UserService.AddChatToUser(createdChat.ID, session.UserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	linker, err := utils.Encrypt(createdChat.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"chat": gin.H{"linker": linker, "id": createdChat.ID, "name": createdChat.Name, "botIdentifier": createdChat.BotIdentifier, "botId": createdChat.BotID, "verified": createdChat.Verified}})

}

// HandleChatsList
func (n chatApi) HandleGetChatsList(c *gin.Context) {
	session, err := sessionutils.GetSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, exist, err := n.services.UserService.GetUserByID(session.UserID)

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
	chats, err := n.services.ChatService.GetChatsByArrayID(chatsId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var respArray []Chat
	for _, chat := range chats {
		respArray = append(respArray, Chat{
			ID:            chat.ID,
			Name:          chat.Name,
			BotIdentifier: chat.BotIdentifier,
			BotID:         chat.BotID,
			LinkOrIdInBot: chat.LinkOrIdInBot,
			Verified:      chat.Verified,
		})
	}
	c.JSON(http.StatusOK, gin.H{"chats": respArray, "offset": offset, "count": count})
}

// HandleUpdateChatSettings by /:id
func (n chatApi) HandleUpdateChatSettings(c *gin.Context) {
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

	session, err := sessionutils.GetSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chatId := c.Param("id")
	//	fmt.Printf("chatId: %v\n", chatId)

	chat, exist, err := n.services.ChatService.GetChatByID(chatId)
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
	err = n.services.ChatService.UpdateChat(chatId, params)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (n chatApi) Router(router *gin.RouterGroup) *gin.RouterGroup {
	r := router.Group("/")
	r.Use(n.middleware.CheckAuth())
	r.GET("/", n.HandleGetChatsList)
	r.POST("/", n.HandleCreateChat)
	r.PUT("/:id", n.HandleUpdateChatSettings)
	return r
}
