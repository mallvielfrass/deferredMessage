package bot

import (
	"deferredMessage/internal/db"
	"deferredMessage/internal/db/mongo/session"
	"deferredMessage/internal/middleware"
	"deferredMessage/internal/utils"
	"deferredMessage/internal/utils/dto"
	reqvalidator "deferredMessage/internal/utils/reqValidator"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type botApi struct {
	db db.DB
}

func Init(db db.DB) botApi {
	return botApi{
		db: db,
	}
}
func ping(c *gin.Context) {
	c.JSON(http.StatusOK, MessageResponse{
		Message: "pong",
	})
}

// HandleGetBotsList retrieves a list of all bots.
//
// @Summary Get a list of all bots
// @Description Retrieves a list of all bots for user
// @Tags Bot
// @Produce json
// @Security Bearer
// @Success 200 {object} BotStructArrayResponse "List of bots retrieved successfully"
// @Failure 400 {object} ErrorResponse "Error getting bots"
// @Router /bot [get]
func (n botApi) HandleGetBotsList(c *gin.Context) {
	bots, err := n.db.Collections.Bot.GetAllBots()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			Error:  "Error getting bots",
			Reason: err.Error()})
		return
	}
	botsResponse := []BotResponse{}
	for _, bot := range bots {
		botsResponse = append(botsResponse, BotResponse{
			Name:    bot.Name,
			Id:      bot.ID.Hex(),
			BotLink: bot.BotLink,

			Creator:  bot.Creator.Hex(),
			Platform: bot.Platform,
		})
	}

	c.JSON(http.StatusOK, BotStructArrayResponse{Bots: botsResponse})
}

// HandleCreateBot creates a new bot.

// @Summary Create a new bot
// @Description Create a new bot with the provided information
// @Tags Bot
// @Accept json
// @Produce json
// @Param botRequest body BotRequest true "Bot request body"
// @Security Bearer
// @Success 200 {object} BotStructResponse "Bot created successfully"
// @Failure 404 {object} ErrorResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/bot/ [post]
func (n botApi) HandleCreateBot(c *gin.Context) {
	body, exist := dto.GetStruct[BotRequest](c, BotRequest{})
	if !exist {
		c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{
			Error: "required params not found"})
		return
	}
	userSession, ok := c.Get("session")
	if !ok {

		c.AbortWithStatusJSON(http.StatusNotFound,
			ErrorResponse{
				Error: "session not found",
			})
		return
	}
	session := userSession.(session.SessionScheme)
	_, isExit, err := n.db.Collections.Platform.GetPlatformByName(body.Platform)
	if err != nil {
		fmt.Printf("HandleCreateBot> err: %v\n", err)
		c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{
			Error: "platform not found"})
		return
	}
	if !isExit {
		c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{
			Error: "platform not found"})
		return
	}
	hashToken, err := utils.Encrypt(body.Token)
	if err != nil {
		fmt.Printf("HandleCreateBot> err: %v\n", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			Error:  "token encryption failed",
			Reason: err.Error(),
		})
		return
	}
	bot, err := n.db.Collections.Bot.CreateBot(body.Name, body.BotLink, session.UserID, body.Platform, hashToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			Error:  "bot creation failed",
			Reason: err.Error()})
		return
	}

	obfuscatedToken := utils.ObfuscateToken(body.Token)

	resp := BotResponse{
		Name:     bot.Name,
		Id:       bot.ID.Hex(),
		BotLink:  bot.BotLink,
		Platform: bot.Platform,
		Creator:  bot.Creator.Hex(),
		Token:    obfuscatedToken,
	}

	c.JSON(http.StatusOK, BotStructResponse{Bot: resp})
}

// UpdateBot updates an existing bot.
//
// @Summary Update an existing bot
// @Description Update an existing bot with the provided information
// @Tags Bot
// @Accept json
// @Produce json
// @Param id path string true "Bot ID" Format(uuid)
// @Param botRequest body BotUpdateRequest true "Bot request body"
// @Security Bearer
// @Success 200 {object} BotStructResponse "Bot updated successfully"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Not found"
// @Router /api/bot/{id} [put]
func (n botApi) UpdateBot(c *gin.Context) {
	botId := c.Param("id")
	if botId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid id"})
		return
	}
	userSession, ok := c.Get("session")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			Error: "no session"})
		return
	}
	botObjectId, err := primitive.ObjectIDFromHex(botId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			Error: "invalid id"})
		return
	}
	session := userSession.(session.SessionScheme)
	bot, isExist, err := n.db.Collections.Bot.GetBotByID(botObjectId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{
			Error: err.Error()})
		return
	}
	if !isExist {
		c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{
			Error: "bot not found"})
		return
	}
	if bot.Creator != session.UserID {
		c.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{
			Error: "you are not creator"})
		return
	}

	var reqBody BotUpdateRequest
	body, err := reqvalidator.ValidateFlatMap(c, &reqBody, reqvalidator.GetTagsFromFlatStruct(BotUpdateRequest{}))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:  "validation error",
			Reason: err.Error()})
		return
	}
	_, isExit, err := n.db.Collections.Platform.GetPlatformByName(reqBody.Platform)
	if err != nil {
		fmt.Printf("get platform err: %+v\n", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			Error: "platform not found"})
		return
	}
	if !isExit {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			Error: "platform not found"})
		return
	}
	utils.HashTokenFromMap(&body)
	//fmt.Printf("body: %+v\n", body)
	botUpdated, _, err := n.db.Collections.Bot.UpdateBot(botObjectId, body)
	if err != nil {
		fmt.Printf("update err: %+v\n", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			Error: "error updating bot"})
		return
	}
	decryptToken, err := utils.Decrypt(botUpdated.HashedToken)
	if err != nil {
		fmt.Printf("decrypt err: %+v\n", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
			Error: "error decrypting token"})
		return
	}
	obfuscatedToken := utils.ObfuscateToken(decryptToken)
	resp := BotResponse{
		Name:     botUpdated.Name,
		Id:       botUpdated.ID.Hex(),
		BotLink:  botUpdated.BotLink,
		Platform: botUpdated.Platform,
		Creator:  botUpdated.Creator.Hex(),
		Token:    obfuscatedToken,
	}

	c.JSON(http.StatusOK, BotStructResponse{Bot: resp})
}
func (n botApi) Router(router *gin.RouterGroup) *gin.RouterGroup {
	r := router.Group("/")
	sessionMiddleware := middleware.InitMiddleware(n.db)
	r.Use(sessionMiddleware.CheckAuth())
	r.GET("/ping", ping)

	r.GET("/", n.HandleGetBotsList)
	r.POST("/", n.HandleCreateBot)
	r.PUT("/:id", n.UpdateBot)
	return r
}
