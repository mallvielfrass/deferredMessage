package platform

import (
	"deferredMessage/internal/db"
	"deferredMessage/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type platformApi struct {
	db db.DB
}

func Init(db db.DB) platformApi {
	return platformApi{
		db: db,
	}
}
func (n platformApi) Router(router *gin.RouterGroup) *gin.RouterGroup {
	r := router.Group("/")
	middlewares := middleware.InitMiddleware(n.db)
	//r.Use(middlewares.CheckAuth())
	r.GET("/list", func(c *gin.Context) {
		platformsBson, err := n.db.Collections.Platform.GetAllPlatforms()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		platforms := make([]PlatformResponse, 0)
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
	adminGroup.GET("/check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "admin pong",
		})
	})

	return r
}