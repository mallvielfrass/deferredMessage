package platform

import (
	"deferredMessage/internal/middleware"
	"deferredMessage/internal/models"
	"deferredMessage/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type platformApi struct {
	services   *service.Service
	middleware *middleware.Middleware
}

func Init(services *service.Service, middleware *middleware.Middleware) platformApi {
	return platformApi{
		services:   services,
		middleware: middleware,
	}
}

// HandleGetPlatformsList retrieves a list of all platforms.
// @Summary Get platforms list
// @Description Retrieves a list of all platforms.
// @Security Bearer
// @Tags Platform
// @Accept json
// @Produce json
// @Success 200 {object} PlatformListResponse "List of platforms"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Router /api/platform [get]
func (n platformApi) HandleGetPlatformsList(c *gin.Context) {
	platforms, err := n.services.PlatformService.GetAllPlatforms()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error:  "Error getting platforms",
			Reason: err.Error()})
		return
	}
	platformsResponse := []PlatformResponse{}
	for _, platform := range platforms {
		platformsResponse = append(platformsResponse, PlatformResponse{
			Name: platform.Name,
		})
	}
	c.JSON(http.StatusOK, PlatformListResponse{
		Platforms: platformsResponse,
	})
}

// HandleCreatePlatform creates a new platform.
// @Security Bearer
// @Tags Platform
// @Accept json
// @Produce json
// @Param body body CreatePlatformRequest true "Create platform request"
// @Success 200 {object} models.PlatformScheme "Platform created successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Router /api/platform [post]
func (n platformApi) HandleCreatePlatform(c *gin.Context) {
	var request CreatePlatformRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error()})
		return
	}
	platform, err := n.services.PlatformService.CreatePlatform(request.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.PlatformScheme{
		Name: platform.Name,
		ID:   platform.ID,
	})
}

// HandleCheckAdminAccess
// @Security Bearer
// @Tags Platform
// @Produce json
// @Router /api/platform/check [get]
// @Success 200 {object} MessageResponse "pong"
func (n platformApi) HandleCheckAdminAccess(c *gin.Context) {
	c.JSON(http.StatusOK, MessageResponse{
		Message: "pong",
	})
}
func (n platformApi) Router(router *gin.RouterGroup) *gin.RouterGroup {
	r := router.Group("/")

	r.GET("/", n.HandleGetPlatformsList)
	adminGroup := r.Group("/")
	adminGroup.Use(n.middleware.CheckAdmin())
	adminGroup.POST("/create", n.HandleCreatePlatform)
	adminGroup.GET("/check", n.HandleCheckAdminAccess)

	return r
}
