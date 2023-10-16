package chat

import (
	"deferredMessage/internal/middleware"
	"deferredMessage/internal/models"
	"deferredMessage/internal/service"
	"deferredMessage/internal/utils"
	"deferredMessage/internal/utils/dto"
	sessionutils "deferredMessage/internal/utils/sessionUtils"
	"net/http"

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

// HandleCreateChat creates a new chat.
// @Summary Create a chat
// @Description Creates a new chat with the provided chat details.
// @Security Bearer
// @Tags Chat
// @Accept json
// @Produce json
// @Param body body ChatRequest true "Chat details"
// @Success 200 {object} CreateChatResponse "Chat created successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Router /api/chat [post]
func (n chatApi) HandleCreateChat(c *gin.Context) {
	body, exist := dto.GetStruct[ChatRequest](c, ChatRequest{})
	if !exist {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "no body"})
		return
	}
	bot, isExist, err := n.services.BotService.GetBotByID(body.BotID)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error()})
		return
	}
	if !isExist {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "bot not found"})
		return
	}
	session, err := sessionutils.GetSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error()})
		return
	}

	createdChat, err := n.services.ChatService.CreateChat(body.Name, bot.ID, bot.ID, session.UserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error()})
		return
	}

	err = n.services.UserService.AddChatToUser(createdChat.ID, session.UserID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error()})
		return
	}
	linker, err := utils.Encrypt(createdChat.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK,
		CreateChatResponse{
			Chat: createdChatType{
				Linker:        linker,
				Name:          createdChat.Name,
				BotIdentifier: createdChat.BotIdentifier,
				BotID:         createdChat.BotID,
				Verified:      createdChat.Verified,
				Id:            createdChat.ID,
			},
		},
	)
}

// HandleGetChatsList retrieves a list of chats created by the authenticated user.
// @Summary Get chats list
// @Description Retrieves a list of chats created by the authenticated user.
// @Security Bearer
// @Tags Chat
// @Accept json
// @Produce json
// @Success 200 {object} ChatsListResponse "List of chats"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Router /api/chat [get]
func (n chatApi) HandleGetChatsList(c *gin.Context) {
	session, err := sessionutils.GetSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error()})
		return
	}
	lims, err := dto.GetLimits(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
	}

	chats, err := n.services.ChatService.GetChatsListByCreatorWithLimits(session.UserID, lims.Offset, lims.Count)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, ChatsListResponse{
		Chats:  chats,
		Offset: lims.Offset,
		Count:  lims.Count,
	})
}

// HandleUpdateChatSettings updates the settings of a chat.
// @Summary Update chat settings
// @Description Updates the settings of a chat identified by the provided ID.
// @Router /api/chat/{id} [put]
// @Security Bearer
// @Tags Chat
// @Accept json
// @Param id path string true "Chat ID"
// @Param name body string false "Chat name"
// @Param linkOrIdInBot body string false "Link or ID in bot"
// @Success 200 {object} models.SuccessResponse "Settings updated successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
func (n chatApi) HandleUpdateChatSettings(c *gin.Context) {
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error()})
	}
	params := dto.TypedMap([]dto.ParamTyped{{Name: "name", Type: "string"}, {Name: "linkOrIdInBot", Type: "string"}}, body)
	session, err := sessionutils.GetSession(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error()})
		return
	}

	chatId := c.Param("id")
	chat, exist, err := n.services.ChatService.GetChatByID(chatId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error()})
		return
	}
	if !exist {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "chat not found"})
		return
	}
	if chat.Creator != session.UserID {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "you are not creator"})
		return
	}
	err = n.services.ChatService.UpdateChat(chatId, params)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse{
		Status: "ok",
	})
}

func (n chatApi) Router(router *gin.RouterGroup) *gin.RouterGroup {
	r := router.Group("/")
	r.Use(n.middleware.CheckAuth())
	r.GET("/", n.HandleGetChatsList)
	r.POST("/", n.HandleCreateChat)
	r.PUT("/:id", n.HandleUpdateChatSettings)
	return r
}
