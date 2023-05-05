package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

type TelegramBot struct {
	Bot *tgbotapi.BotAPI
}

type Message struct {
	Command string
	ChatId  int64
}

func NewTelegramBot(token string) *TelegramBot {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	return &TelegramBot{Bot: bot}
}

func (t *TelegramBot) SendMessage(chatId int64, message string) error {
	msg := tgbotapi.NewMessage(chatId, message)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := t.Bot.Send(msg)
	return err
}

func handleMessage(t *tgbotapi.Message) *Message {
	if !strings.HasPrefix(t.Text, "/") {
		return nil
	}

	return &Message{
		Command: t.Text,
		ChatId:  t.Chat.ID,
	}
}

func (t *TelegramBot) GetUpdates(ctx context.Context) chan Message {
	u := tgbotapi.NewUpdate(0)
	updates := t.Bot.GetUpdatesChan(u)
	messageChan := make(chan Message)
	u.Timeout = 60

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case update := <-updates:
				handleUpdate(update)
				messageChan <- *handleMessage(update.Message)
			}
		}
	}()

	return messageChan

}

func handleUpdate(update tgbotapi.Update) {
	handleMessage(update.Message)
}
