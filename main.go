package main

import (
	"context"
	"fmt"
	"git.foxminded.ua/foxstudent104911/2.1about-me-bot/config"
	"git.foxminded.ua/foxstudent104911/2.1about-me-bot/telegram"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	infoAboutMe             = "name: Kseniia\nage: 24\ngender: female"
	socNetworksLinks        = "Instagram: https://instagram.com/some_insta\nFacebook: https://facebook.com/any_facebook\nLinkedIn: https://linkedin.com/some_linkedin"
	start                   = "<b>List of comands:</b>\n/about -- Info about author\n/links -- links to social networks"
	answerForUnknownCommand = "I have no clue what you are talking about"
	done                    = make(chan bool, 1)
	bot                     *telegram.TelegramBot
	err                     error
)

func main() {
	config := config.NewConfig()
	bot, err = telegram.NewTelegramBot(config.Token)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	messages := bot.GetUpdates(ctx)

	go handleMessages(messages)

	log.Println("Start listening for updates. Press enter to stop")

	go handleSignals()
	<-done

	cancel()

}

func handleSignals() {
	sign := make(chan os.Signal, 1)
	signal.Notify(sign,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	for {
		s := <-sign
		fmt.Printf("Got %s signal\n", s)
		done <- true
	}
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
	default:
		err = handleUnknownMessage(message.ChatId)
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

func handleUnknownMessage(chatId int64) error {
	err := bot.SendMessage(chatId, answerForUnknownCommand)
	if err != nil {
		return err
	}
	return nil
}
