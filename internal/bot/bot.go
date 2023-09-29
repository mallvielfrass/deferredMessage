package bot

import (
	"fmt"
	"time"

	tele "gopkg.in/telebot.v3"
)

type Bot struct {
	bot *tele.Bot
}

func InitBot(token string) (Bot, error) {

	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return Bot{}, err
	}
	return Bot{b}, nil
}

// mount
func (b *Bot) Mount() {
	b.bot.Handle("/hello", func(c tele.Context) error {
		return c.Send("Hello!")
	})
}

func (b *Bot) Stop() {
	b.bot.Stop()
}
func (b *Bot) Start() {
	fmt.Println("Bot started with name:", b.bot.Me.Username)
	b.bot.Start()

}
