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
func (n platformApi) Router(router *gin.RouterGroup) *gin.RouterGroup {
	r := router.Group("/")
	middlewares := middleware.InitMiddleware(n.services)
	//r.Use(middlewares.CheckAuth())
	r.GET("/", func(c *gin.Context) {
		platformsBson, err := n.services.PlatformService.GetAllPlatforms()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, models.ErrorResponse{
				Error: err.Error()})
		}
		var platforms []PlatformResponse
		for _, platform := range platformsBson {
			platforms = append(platforms, PlatformResponse{
				Name: platform.Name,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"platforms": platforms,
		})
	})
	adminGroup := r.Group("/admin")
	adminGroup.Use(middlewares.CheckAdmin())
	adminGroup.POST("/create", func(c *gin.Context) {
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
		c.JSON(http.StatusOK, gin.H{
			"platform": platform,
			"_id":      platform.ID,
		})
	})
	adminGroup.GET("/check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "admin pong",
		})
	})

	return r
}
