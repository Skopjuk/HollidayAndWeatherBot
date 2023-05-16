package main

import (
	"context"
	"git.foxminded.ua/foxstudent104911/2.1about-me-bot/config"
	"git.foxminded.ua/foxstudent104911/2.1about-me-bot/holiday"
	"git.foxminded.ua/foxstudent104911/2.1about-me-bot/telegram"
	"github.com/enescakir/emoji"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

var (
	infoAboutMe             = "name: Kseniia\nage: 24\ngender: female"
	socNetworksLinks        = "Instagram: https://instagram.com/some_insta\nFacebook: https://facebook.com/any_facebook\nLinkedIn: https://linkedin.com/some_linkedin"
	help                    = "<b>List of comands:</b>\n/about -- Info about author\n/links -- links to social networks\n/start -- list of holidays by country"
	answerForUnknownCommand = "I have no clue what you are talking about"
	done                    = make(chan bool, 1)
	bot                     *telegram.TelegramBot
	holidayAPI              *holiday.HolidayAPI
	start                   = "Choose country \n"
	buttons                 = map[emoji.Emoji]string{
		emoji.FlagForUkraine:     "UA",
		emoji.FlagForAfghanistan: "AG",
		emoji.FlagForJapan:       "JP",
		emoji.FlagForMalaysia:    "ML",
	}
)

func main() {
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	config, err := config.NewConfig()
	if err != nil {
		logrus.Fatal("unable to load config")
	}

	logLevel, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		logrus.Error("Wrong log level in config")
		logLevel = logrus.InfoLevel
	}

	logrus.SetLevel(logLevel)

	bot, err = telegram.NewTelegramBot(config.Token, config.BotDebug)
	if err != nil {
		logrus.Fatal(err)
	}

	holidayAPI = holiday.NewHolidayAPI(config.HolidayAPI)

	ctx, cancel := context.WithCancel(context.Background())

	updates := bot.GetUpdates(ctx)

	go handleUpdates(updates)

	logrus.Info("Start listening for updates")

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
		logrus.Errorf("Got %s signal\n", s)
		done <- true
	}
}

func handleUpdates(update chan telegram.TelegramUpdate) {
	for {
		handleUpdate(<-update)
	}
}

func handleUpdate(update telegram.TelegramUpdate) {
	if update.Message != nil {
		err := handleMessage(*update.Message)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"chatId": update.Message.ChatId,
			}).Error("message unhandled")
		}
	} else if update.Callback != nil {
		err := handleCallback(*update.Callback)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"chatId": update.Callback.ChatId,
			}).Error("callback unhandled")
		}
	}
}

func handleCallback(callback telegram.Callback) error {
	var err error
	var holidayList string

	pressedButton := buttons[emoji.Emoji(callback.Button)]

	holidayList, err = holidayAPI.TransformListOfHolidaysToStr(pressedButton)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"chatId":         callback.ChatId,
			"button_pressed": callback.Button,
			"error":          err,
		}).Error(err)
	}

	bot.SendMessageWithCallback(callback.Message.Chat.ID, callback, holidayList)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"chatId":         callback.ChatId,
			"button_pressed": callback.Button,
			"error":          err,
		}).Error("Callback were not handled")
	} else {
		logrus.WithFields(logrus.Fields{
			"chatId": callback.ChatId,
			"button": callback.Button,
		}).Info("answer sent")
	}
	return err
}

func handleMessage(message telegram.Message) error {
	var err error

	switch message.Command {
	case "/help":
		err = sendHelp(message.ChatId)
	case "/links":
		err = sendLinks(message.ChatId)
	case "/about":
		err = sendInfo(message.ChatId)
	case "/start":
		err = sendStart(message.ChatId)
	default:
		err = handleUnknownMessage(message.ChatId)
	}

	if err != nil {
		logrus.Error("Message were not handled")
	} else {
		logrus.WithFields(logrus.Fields{
			"chatId":  message.ChatId,
			"message": message.Command,
		}).Info("answer sent")
	}
	return err
}

func sendInfo(chatId int64) error {
	err := bot.SendMessage(chatId, infoAboutMe)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"chatId": chatId,
		}).Error(err)
	}
	return err
}

func sendLinks(chatId int64) error {
	err := bot.SendMessage(chatId, socNetworksLinks)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"chatId": chatId,
		}).Error(err)
	}
	return err
}

func sendHelp(chatId int64) error {
	err := bot.SendMessage(chatId, help)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"chatId": chatId,
		}).Error(err)
	}
	return err
}

func sendStart(chatId int64) error {
	err := bot.SendMenu(chatId, start)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"chatId": chatId,
		}).Error(err)
	}
	return err
}

func handleUnknownMessage(chatId int64) error {
	err := bot.SendMessage(chatId, answerForUnknownCommand)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"chatId": chatId,
		}).Error(err)
	}
	return err
}
