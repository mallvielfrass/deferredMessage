package internal

import (
	"deferredMessage/config"
	"deferredMessage/internal/api/noauth"
	"deferredMessage/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
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
func (app DefferedMessageApp) Run() error {
	db, err := db.ConnectDB(app.Config.DBHost, app.Config.DBName)
	if err != nil {
		return err
	}
	defer db.Disconnect()
	router := gin.Default()

	r := router.Group("/api")
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	noAuthR := r.Group("/nauth")
	noAuthRouter := noauth.Init(db)
	noAuthRouter.Router(noAuthR)
	//	fmt.Printf("config: %#v\n", app.Config)
	err = router.Run(app.Config.HostPort)
	return err
}
