package telegram

import (
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Sender struct {
	bot    *tgbot.BotAPI
	chatID int64
}

func New(token string, chatID int64) (*Sender, error) {
	b, err := tgbot.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Sender{bot: b, chatID: chatID}, nil
}

func (s *Sender) Send(msg string) error {
	m := tgbot.NewMessage(s.chatID, msg)
	_, err := s.bot.Send(m)
	return err
}
