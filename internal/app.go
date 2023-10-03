package internal

import (
	"deferredMessage/config"
	_ "deferredMessage/docs"
	"deferredMessage/internal/api/auth/user"
	"deferredMessage/internal/api/noauth"
	"deferredMessage/internal/db"

	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type DefferedMessageApp struct {
	Config config.Config
}

func NewApp(confPath string) (DefferedMessageApp, error) {
	conf, err := config.InitConfig(confPath)
	if err != nil {
		return DefferedMessageApp{
			Config: conf,
		}, err
	}
	return DefferedMessageApp{
		Config: conf,
	}, nil
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
func InitRouter(db db.DB) *gin.Engine {
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
	noauth.Init(db).Router(r.Group("/nauth"))
	user.Init(db).Router(r.Group("/auth/user/"))
	return router
}
func (app DefferedMessageApp) Run() error {
	db, err := db.ConnectDB(app.Config.DBHost, app.Config.DBName)
	if err != nil {
		return err
	}
	defer db.Disconnect()
	// tgBot, err := bot.InitBot(app.Config.TelegramBotToken, db)
	// if err != nil {
	// 	return err
	// }
	//defer tgBot.Stop()
	//tgBot.Mount()
	//go tgBot.Start()

	//	fmt.Printf("config: %#v\n", app.Config)
	router := InitRouter(db)
	err = router.Run(app.Config.HostPort)
	return err
}
