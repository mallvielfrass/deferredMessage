package handler

import (
	"deferredMessage/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /tapi/example/helloworld [get]
func Helloworld(g *gin.Context) {
	g.JSON(http.StatusOK, "helloworld")
}
func (h *Handler) newRouter() *gin.Engine {
	router := gin.Default()
	r := router.Group("/api")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	//docs.SwaggerInfo.BasePath = "/api/"
	v1 := router.Group("/tapi")
	{
		eg := v1.Group("/example")
		{
			eg.GET("/helloworld", Helloworld)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	//html file serve
	router.GET("/front", func(c *gin.Context) {
		c.File("./internal/static/index.html")
	})
	///static files
	router.Static("/static", "./internal/static")

	// sessionMiddleware := middleware.InitMiddleware(db)
	// authUserApi := r.Group("/auth")
	// authUserApi.Use(sessionMiddleware.CheckAuth())

	// noauth.Init(db).Router(r.Group("/nauth"))
	// user.Init(db).Router(authUserApi.Group("/user"))
	// platform.Init(db).Router(authUserApi.Group("/platform"))
	// bot.Init(db).Router(authUserApi.Group("/bot"))
	// admin.Init(db).Router(authUserApi.Group("/admin"))
	return router
}
func (h *Handler) Run(port string) error {
	router := h.newRouter()
	return router.Run(port)
}
