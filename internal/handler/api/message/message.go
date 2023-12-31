package message

import (
	"deferredMessage/internal/clocker"
	"deferredMessage/internal/middleware"
	"deferredMessage/internal/models"
	"deferredMessage/internal/service"
	"deferredMessage/internal/utils/dto"
	sessionutils "deferredMessage/internal/utils/sessionUtils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type messageApi struct {
	services   *service.Service
	middleware *middleware.Middleware
	clock      *clocker.Clocker
}

func Init(services *service.Service, middleware *middleware.Middleware, clock *clocker.Clocker) *messageApi {
	return &messageApi{
		services:   services,
		middleware: middleware,
		clock:      clock,
	}
}

func (h *messageApi) ping(c *gin.Context) {
	c.JSON(http.StatusOK, models.PingMessageResponse{
		Message: "pong",
	})
}

// HandleCreateNewMessage handles the API endpoint for creating a new message.
// @Summary Create new message
// @Description Creates a new message.
// @Accept json
// @Produce json
// @Param message body NewMessageRequest true "Message"
// @Success 200 {object} MessageResponse "Message"
// @Failure 400 {object} models.ErrorResponse "Error response"
// @Router /api/auth/messages [post]
func (h *messageApi) HandleCreateNewMessage(c *gin.Context) {
	session, err := sessionutils.GetSession(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	var msg NewMessageRequest
	err = c.BindJSON(&msg)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	//convert int64 to time

	tm := time.Unix(msg.Time, 0)
	message, err := h.services.MessageService.CreateNewMessage(session.UserID, models.Message{Message: msg.Message, ChatId: msg.ChatId, Time: tm})
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, MessageResponse{
		Message: message,
	})

}

// HandleListOfAllMessages handles the API endpoint for getting a list of all messages.
// @Summary Get list of all messages
// @Description Retrieves a list of all messages based on the provided limits.
// @Accept json
// @Produce json
// @Param limit query int false "Number of messages to retrieve (default is 10)"
// @Param offset query int false "Offset for pagination (default is 0)"
// @Security Bearer
// @Success 200 {object} MessageListResponse "List of all messages"
// @Failure 400 {object} models.ErrorResponse "Error response"
// @Router /api/auth/messages [get]
func (h *messageApi) HandleListOfAllMessages(c *gin.Context) {
	lims, err := dto.GetLimits(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	session, err := sessionutils.GetSession(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	msg, err := h.services.MessageService.GetListOfAllMessages(session.UserID, lims.Offset, lims.Count)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, MessageListResponse{
		Messages: msg,
	})
}
func (h *messageApi) Router(router *gin.RouterGroup) *gin.RouterGroup {
	r := router.Group("/")
	r.GET("/ping", h.ping)
	return r
}
