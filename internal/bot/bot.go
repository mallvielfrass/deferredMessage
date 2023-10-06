package bot

import (
	"deferredMessage/internal/db"
	"deferredMessage/internal/utils"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type Bot struct {
	bot *tele.Bot
	db  db.DB
}

func InitBot(token string, db db.DB) (Bot, error) {

	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return Bot{}, err
	}
	return Bot{b, db}, nil
}

// mount
func (b *Bot) Mount() {
	b.bot.Handle("/hello", func(c tele.Context) error {
		return c.Send("Hello!")
	})
	adminOnly := b.bot.Group()
	middleware.Logger()
	// adminOnly.Use(func(c tele.Context) error {
	// 	admins, err := b.bot.AdminsOf(c.Chat()).ID(c.Sender().ID)
	// 	fmt.Println(admins)
	// 	return err
	// })
	adminOnly.Use(func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if c.Sender().ID == c.Chat().ID {
				return fmt.Errorf("is private chat")
			}
			admins, err := b.bot.AdminsOf(c.Chat())
			if err != nil {
				return err
			}
			for _, admin := range admins {
				if admin.User.ID == c.Sender().ID {
					return next(c)
				}
			}
			return fmt.Errorf("not admin")
		}
	})
	adminOnly.Handle("/link", func(c tele.Context) error {
		//get params msg
		params := c.Message().Text
		fmt.Println(params)
		groupIdEncrypted := strings.TrimSpace(params)
		groupId, err := utils.Decrypt(groupIdEncrypted)
		if err != nil {
			fmt.Println("/link: ERR:", err)
			return c.Reply("invalid group id")
		}
		sessionObjectID, err := primitive.ObjectIDFromHex(groupId)
		if err != nil {
			return c.Reply("invalid group id")
		}
		fmt.Println(sessionObjectID)
		chat, exist, err := b.db.Collections.Chat.GetChatByID(sessionObjectID)
		if err != nil {
			fmt.Println("/link: ERR:", err)
			return c.Reply("invalid group id")
		}
		if !exist {
			return c.Reply("invalid group id")
		}
		if chat.Verified {
			return c.Reply("Chat already verified")
		}

		chat.Verified = true
		chat.LinkOrIdInNetwork = strconv.FormatInt(c.Chat().ID, 10)
		paramsMap := map[string]interface{}{
			"verified":          chat.Verified,
			"linkOrIdInNetwork": chat.LinkOrIdInNetwork,
		}
		err = b.db.Collections.Chat.UpdateChat(chat.ID, paramsMap)
		if err != nil {
			fmt.Println("/link: ERR:", err)
			return c.Reply("Server error")
		}
		return c.Send("Chat verified")
	})
}

func (b *Bot) Stop() {
	b.bot.Stop()
}
func (b *Bot) Start() {
	fmt.Println("Bot started with name:", b.bot.Me.Username)
	b.bot.Start()

}
