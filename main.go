package main

import (
	"bufio"
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

var (
	infoAboutMe      = "name: Kseniia\nage: 24\ngender: female"
	socNetworksLinks = "Instagram: https://instagram.com/some_insta\nFacebook: https://facebook.com/any_facebook\nLinkedIn: https://linkedin.com/some_linkedin"
	start            = "<b>List of comands:</b>\n/about -- Info about author\n/links -- links to social networks"

	bot *tgbotapi.BotAPI
)

func main() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	bot, err = tgbotapi.NewBotAPI(os.Getenv("TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	updates := bot.GetUpdatesChan(u)

	go receiveUpdates(ctx, updates)

	log.Println("Start listening for updates. Press enter to stop")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	cancel()

}

func receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			handleUpdate(update)
		}
	}
}

func handleUpdate(update tgbotapi.Update) {
	handleMessage(update.Message)
}

func handleMessage(message *tgbotapi.Message) {
	user := message.From
	text := message.Text

	if user == nil {
		return
	}

	var err error
	if strings.HasPrefix(text, "/") {
		err = handleCommand(message.Chat.ID, text)
	}

	if err != nil {
		log.Printf("An error occured: %s", err.Error())
	}
}

func handleCommand(chatId int64, command string) error {
	var err error

	switch command {
	case "/start":
		err = sendStartAndHelp(chatId)
		break
	case "/help":
		err = sendStartAndHelp(chatId)
		break

	case "/links":
		err = sendLinks(chatId)
		break

	case "/about":
		err = sendInfo(chatId)
		break
	}

	return err
}

func sendInfo(chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, infoAboutMe)
	_, err := bot.Send(msg)
	return err
}

func sendLinks(chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, socNetworksLinks)
	_, err := bot.Send(msg)
	return err
}

func sendStartAndHelp(chatId int64) error {
	msg := tgbotapi.NewMessage(chatId, start)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := bot.Send(msg)
	return err
}
