package internal

import (
	"deferredMessage/config"
	_ "deferredMessage/docs"
	"deferredMessage/internal/handler"
	"deferredMessage/internal/repository"
	db "deferredMessage/internal/repository"
	"deferredMessage/internal/service"
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
	// tgBot, err := bot.InitBot(app.Config.TelegramBotToken, db)
	// if err != nil {
	// 	return err
	// }
	//defer tgBot.Stop()
	//tgBot.Mount()
	//go tgBot.Start()

	//	fmt.Printf("config: %#v\n", app.Config)
	repos := repository.NewRepository(db.Driver)
	services := service.NewService(repos)
	handler := handler.NewHandler(services)
	err = handler.Run(app.Config.HostPort)

	return err
}
