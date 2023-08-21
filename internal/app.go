package internal

import (
	"deferredMessage/config"
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

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	//	fmt.Printf("config: %#v\n", app.Config)
	err := r.Run(app.Config.HostPort)
	return err
}
