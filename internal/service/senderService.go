package service

import (
	"deferredMessage/internal/models"
	"deferredMessage/internal/repository"
	"deferredMessage/internal/utils"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type senderService struct {
	repos *repository.Repository
}

func NewSenderService(repos *repository.Repository) *senderService {
	return &senderService{
		repos: repos,
	}
}
func (s *senderService) tgSend(msg models.Message, bot models.BotScheme, chat models.ChatScheme) error {
	tokenBot, err := utils.Decrypt(bot.HashedToken)
	if err != nil {
		return err
	}
	botInstance, err := tgbotapi.NewBotAPI(tokenBot)
	if err != nil {
		return err
	}
	int64ChatId, err := strconv.ParseInt(chat.ID, 10, 64)
	if err != nil {
		return err
	}
	_, err = botInstance.Send(tgbotapi.NewMessage(int64ChatId, msg.Message))
	if err != nil {
		return err
	}
	//	fmt.Println(retunMsg)
	return nil
}
func (s *senderService) Send(msg models.Message) error {
	message, isExist, err := s.repos.Message.GetMessageByID(msg.Id)
	if err != nil {
		return err
	}
	if !isExist {
		return fmt.Errorf("message not found")
	}
	err = s.repos.Message.SetMessageIsProcessed(msg.Id)
	if err != nil {
		return err
	}

	chat, isExist, err := s.repos.Chat.GetChatByID(message.ChatId)
	if err != nil {
		return err
	}
	if !isExist {
		return fmt.Errorf("chat not found")
	}
	bot, isExist, err := s.repos.Bot.GetBotByID(chat.BotID)
	if err != nil {
		return err
	}
	if !isExist {
		return fmt.Errorf("bot not found")
	}
	platform, isExist, err := s.repos.Platform.GetPlatformByName(bot.Platform)
	if err != nil {
		return err
	}
	if !isExist {
		return fmt.Errorf("platform not found")
	}
	if platform.Name == "telegram" {
		err := s.tgSend(msg, bot, chat)
		if err != nil {
			s.repos.Message.SetMessageError(msg.Id, err.Error())
			return err
		}
	} else {
		s.repos.Message.SetMessageError(msg.Id, fmt.Errorf("platform not found").Error())
		return fmt.Errorf("platform not found")
	}

	return s.repos.Message.SetIsSended(msg.Id)
}
