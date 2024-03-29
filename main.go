package main

import (
	"context"
	"fmt"
	"git.foxminded.ua/foxstudent104911/2.1about-me-bot/config"
	"git.foxminded.ua/foxstudent104911/2.1about-me-bot/holiday"
	"git.foxminded.ua/foxstudent104911/2.1about-me-bot/telegram"
	"git.foxminded.ua/foxstudent104911/2.1about-me-bot/weather"
	"git.foxminded.ua/foxstudent104911/2.1about-me-bot/weather_subscription"
	"github.com/enescakir/emoji"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"
)

var (
	infoAboutMe                            = "name: Kseniia\nage: 24\ngender: female"
	socNetworksLinks                       = "Instagram: https://instagram.com/some_insta\nFacebook: https://facebook.com/any_facebook\nLinkedIn: https://linkedin.com/some_linkedin"
	weatherSubscriptionSetupInstruction    = "For setting up weather subscription folow these steps:\n1. Choose hour when you want to get update by pressing on /setup_weather_subscription_time\n2. Press /setup_weather_subscription_location \n3. Send location for which you want to receive updates"
	help                                   = "<b>List of comands:</b>\n/about -- Info about author\n/links -- links to social networks\n/start -- list of holidays by country"
	answerForUnknownCommand                = "I have no clue what are you talking about"
	answerForLocationUpdateInWeatherSubscr = "Please, send geolocation for which you want to receive updates"
	done                                   = make(chan bool, 1)
	bot                                    *telegram.TelegramBot
	holidayAPI                             *holiday.HolidayAPI
	weatherApi                             *weather.WeatherApi
	start                                  = "Choose country \n"
	database                               *mongo.Database
	subscriptionDB                         *weather_subscription.MongoSubscriptionConnection
	tickerMinutes                          int
)

func main() {
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	config, err := config.NewConfig()
	if err != nil {
		logrus.Fatal("unable to load config")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.MongoUrl))
	if err != nil {
		logrus.Fatal(err)
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		logrus.Fatal(err)
	}

	database = client.Database("subscriptions")

	logLevel, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		logrus.Error("Wrong log level in config")
		logLevel = logrus.InfoLevel
	}

	tickerMinutes = config.TickerMinutes

	logrus.SetLevel(logLevel)

	bot, err = telegram.NewTelegramBot(config.Token, config.BotDebug)
	if err != nil {
		logrus.Fatal(err)
	}

	holidayAPI = holiday.NewHolidayAPI(config.HolidayApiToken, config.HolidayApiUrlAddress)

	weatherApi = weather.NewWeatherApi(config.WeatherApiToken, config.WeatherApiUrlAddress)

	subscriptionDB = weather_subscription.NewMongoSubscriptionConnection(*database.Collection("subscriptions"))

	ctx, cancel := context.WithCancel(context.Background())
	logrus.Info("Start listening for updates")
	fmt.Println("listening")
	updates := bot.GetUpdates(ctx)

	go sendSubscriptionsByTicker()
	go handleUpdates(updates)

	go handleSignals()
	<-done

	cancel()

}

func sendSubscriptionsByTicker() {
	now := time.Now()
	ticker := time.NewTicker(time.Duration(tickerMinutes) * time.Second)

	for {
		<-ticker.C

		listOfSubscriptions, _ := subscriptionDB.GetSubscriptionDataFromMongoBySetTime(strconv.Itoa(now.Hour()))
		if listOfSubscriptions != nil {
			logrus.Info("no subscriptions for this hour")
		}

		for _, n := range listOfSubscriptions {
			err := handleMessageWithGeoToWeatherApi(n.ChatId, n.Location.Longitude, n.Location.Latitude)
			if err != nil {
				logrus.Error(err)
			}
		}
	}
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

		logrus.WithFields(logrus.Fields{
			"message": update.Message,
		}).Info("message handled")
	} else if !chooseClickedButton(update) {
		err := handleFlagButtonCallback(*update.Callback)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"chatId": update.Callback.CallbackId,
			}).Error("callback unhandled")
		}
	} else if chooseClickedButton(update) {
		err := handleHoursButtonCallback(*update.Callback)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"chatId": update.Callback.CallbackId,
			}).Error("callback unhandled")
		}
	}
}

func chooseClickedButton(update telegram.TelegramUpdate) bool {
	match, err := regexp.Match(`\d{1,2}`, []byte(update.Callback.Button))

	if err != nil {
		logrus.Error(err)
	}

	if match {
		return true
	} else {
		return false
	}
}

func handleFlagButtonCallback(callback telegram.Callback) error {
	var err error
	var holidayList string
	excuse := "Sorry, in reason of internal problem we can't show you list of holidays right now. Please, try to repeat in couple of minutes."

	pressedButton := telegram.FlagButtons[emoji.Emoji(callback.Button)]

	holidayList, err = holidayAPI.TransformListOfHolidaysToStr(pressedButton)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"chatId":         callback.CallbackId,
			"button_pressed": callback.Button,
			"error":          err,
		}).Error(err)
		bot.SendMessageWithCallback(callback.Message.Chat.ID, callback, excuse)
		return err
	}

	bot.SendMessageWithCallback(callback.Message.Chat.ID, callback, holidayList)

	logrus.WithFields(logrus.Fields{
		"callbackId": callback.CallbackId,
		"button":     callback.Button,
	}).Info("answer sent")

	return err
}

func handleHoursButtonCallback(callback telegram.Callback) error {
	var err error
	excuse := "Sorry, in reason of internal problem we can't set up your subscription for weather updates. Please, try to repeat in couple of minutes."

	buttonInInt, err := strconv.Atoi(callback.Button)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"username": callback.User.UserName,
		}).Error(err)
	}

	pressedButton := telegram.HoursMap()[buttonInInt]
	username := callback.User.UserName

	err = subscriptionDB.UpdateSubscriptionWithTime(username, pressedButton)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"callback_id":    callback.CallbackId,
			"button_pressed": callback.Button,
			"error":          err,
		}).Error(err)

		bot.SendMessageWithCallback(callback.Message.Chat.ID, callback, excuse)
	}

	bot.SendMessageWithCallback(callback.Message.Chat.ID, callback, "Subscription time set up")

	logrus.WithFields(logrus.Fields{
		"callback_id": callback.CallbackId,
		"button":      callback.Button,
	}).Info("answer sent")

	return err
}

func handleMessageWithGeoForWeatherSubscription(message telegram.Message) error {
	var err error
	newMessage := weather_subscription.Message{
		ChatId:    message.ChatId,
		Longitude: message.Location.Longitude,
		Latitude:  message.Location.Latitude,
		Username:  message.Username.UserName,
	}

	err = subscriptionDB.UpdateSubscriptionWithLocation(newMessage)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"chat_id": message.ChatId,
			"user":    message.Username,
		}).Error(err)
	}

	err = bot.SendMessage(message.ChatId, "Location added")
	if err != nil {
		return err
	}

	return nil
}

func checkIfUserAlreadyExist(message telegram.Message) bool {
	usernameInBson := bson.D{{"username", message.Username.UserName}}

	cursor, err := database.Collection("subscriptions").Find(context.TODO(), usernameInBson)
	if err != nil {
		logrus.Error(err)
	}

	if cursor != nil {
		return true
	}

	return false
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
	case "/setup_weather_subscription":
		err = sendWeatherSubscriptionInstruction(message.ChatId)
	case "/setup_weather_subscription_time":
		err = sendHoursMenu(message.ChatId)
	case "/setup_weather_subscription_location":
		err = sendWeatherSubscriptionLocation(message.ChatId)
	default:
		if checkIfUserAlreadyExist(message) && message.Location != nil {
			err = handleMessageWithGeoForWeatherSubscription(message)
		} else if message.Location != nil {
			err = sendStart(message.ChatId)
		} else {
			err = handleUnknownMessage(message.ChatId)
		}
	}

	if err != nil {
		logrus.Error("Message were not handled")
	}
	logrus.WithFields(logrus.Fields{
		"chatId":  message.ChatId,
		"message": message.Command,
	}).Info("answer sent")

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

func sendHoursMenu(chatId int64) error {
	err := bot.SendHoursMenu(chatId, "Choose hour")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"chatId": chatId,
		}).Error(err)
	}
	return err
}

func sendWeatherSubscriptionLocation(chatId int64) error {
	err := bot.SendMessage(chatId, answerForLocationUpdateInWeatherSubscr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"chatId": chatId,
		}).Error(err)
	}
	return err
}

func sendStart(chatId int64) error {
	err := bot.SendFlagsMenu(chatId, start)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"chatId": chatId,
		}).Error(err)
	}
	return err
}

func sendWeatherSubscriptionInstruction(chatId int64) error {
	err := bot.SendMessage(chatId, weatherSubscriptionSetupInstruction)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"chatId": chatId,
		}).Error(err)
	}
	return err
}

func handleMessageWithGeoToWeatherApi(chatId int64, longitude float64, latitude float64) error {
	var newMessage string

	weatherMessage, err := weatherApi.MakeRequest(longitude, latitude)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"chatId": chatId,
		}).Error("Problems with making request")
		return err
	}

	newMessage = fmt.Sprintf(
		"<b>Real Temperature:</b> %.2f\n<b>Feels Like: </b>%.2f\n<b>Main: </b>%s\n"+
			"<b>Minimal Temperature: </b>%.2f\n<b>Maximum Temperature: </b>%.2f\n<b>Humidity: </b>%.2f",
		weatherMessage.MainWeather.Temp-272.15,
		weatherMessage.MainWeather.FeelLike-272.15,
		weatherMessage.Weather[0].Main,
		weatherMessage.MainWeather.TempMin-272.15,
		weatherMessage.MainWeather.TempMax-272.15,
		weatherMessage.MainWeather.Humidity,
	)

	err = bot.SendMessage(int64(chatId), newMessage)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"chatId": chatId,
		}).Error("Problems with sending message")
	}

	return nil
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
