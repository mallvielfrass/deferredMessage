package bot

import (
	"deferredMessage/internal/db"
	"deferredMessage/internal/db/mongo/session"
	"deferredMessage/internal/middleware"
	"deferredMessage/internal/utils/dto"
	reqvalidator "deferredMessage/internal/utils/reqValidator"
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
func (n botApi) HandleGetBotsList(c *gin.Context) {
	bots, err := n.db.Collections.Bot.GetAllBots()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	c.JSON(http.StatusOK, gin.H{"bots": botsResponse})
}
func (n botApi) HandleCreateBot(c *gin.Context) {
	body, exist := dto.GetStruct[BotResponse](c, BotResponse{})
	if !exist {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "required params not found"})
		return
	}
	userSession, ok := c.Get("session")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no session"})
		return
	}
	session := userSession.(session.SessionScheme)

	bot, err := n.db.Collections.Bot.CreateBot(body.Name, body.Id, body.BotLink, session.UserID, body.Platform)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"bot": gin.H{"name": bot.Name,
		"_id": bot.ID.Hex(), "botLink": bot.BotLink,
		"creator": bot.Creator.Hex(), "platform": bot.Platform}})
}
func (n botApi) UpdateBot(c *gin.Context) {
	botId := c.Param("id")
	if botId == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	userSession, ok := c.Get("session")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no session"})
		return
	}
	botObjectId, err := primitive.ObjectIDFromHex(botId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	session := userSession.(session.SessionScheme)
	bot, isExist, err := n.db.Collections.Bot.GetBotByID(botObjectId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !isExist {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "bot not found"})
		return
	}
	if bot.Creator != session.UserID {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "you are not creator"})
		return
	}

	type BodyValidate struct {
		Name     string `json:"name" validate:"required"`
		Platform string `json:"platform" validate:"required"`
	}

	body, err := reqvalidator.ValidateFlatMap(c, &BodyValidate{}, reqvalidator.GetTagsFromFlatStruct(BodyValidate{}))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"bot": body})
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
