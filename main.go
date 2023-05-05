package main

import (
	"bufio"
	"context"
	"git.foxminded.ua/foxstudent104911/2.1about-me-bot/config"
	"git.foxminded.ua/foxstudent104911/2.1about-me-bot/telegram"
	"log"
	"os"
)

var (
	infoAboutMe      = "name: Kseniia\nage: 24\ngender: female"
	socNetworksLinks = "Instagram: https://instagram.com/some_insta\nFacebook: https://facebook.com/any_facebook\nLinkedIn: https://linkedin.com/some_linkedin"
	start            = "<b>List of comands:</b>\n/about -- Info about author\n/links -- links to social networks"
	token            string
	bot              *telegram.TelegramBot
)

func main() {
	config := config.NewConfig()
	bot = telegram.NewTelegramBot(config.Token)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	messages := bot.GetUpdates(ctx)

	go handleMessages(messages)

	log.Println("Start listening for updates. Press enter to stop")

	bufio.NewReader(os.Stdin).ReadBytes('\n')

	cancel()

}

func handleMessages(messages chan telegram.Message) {
	for {
		err := handleCommand(<-messages)
		if err != nil {
			return
		}
	}
}

func handleCommand(message telegram.Message) error {
	var err error

	switch message.Command {
	case "/start":
		err = sendStartAndHelp(message.ChatId)
	case "/help":
		err = sendStartAndHelp(message.ChatId)
	case "/links":
		err = sendLinks(message.ChatId)
	case "/about":
		err = sendInfo(message.ChatId)
	}

	return err
}

func sendInfo(chatId int64) error {
	err := bot.SendMessage(chatId, infoAboutMe)
	if err != nil {
		return err
	}
	return nil
}

func sendLinks(chatId int64) error {
	err := bot.SendMessage(chatId, socNetworksLinks)
	if err != nil {
		return err
	}
	return nil
}

func sendStartAndHelp(chatId int64) error {
	err := bot.SendMessage(chatId, start)
	if err != nil {
		return err
	}
	return nil
}
