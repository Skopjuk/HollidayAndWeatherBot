package telegram

import (
	"context"
	"github.com/enescakir/emoji"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
)

var (
	afghanistanFlagButton = emoji.FlagForAfghanistan
	malaysiaFlagButton    = emoji.FlagForMalaysia
	japanFlagButton       = emoji.FlagForJapan
	ukrainianFlagButton   = emoji.FlagForUkraine
	startMenuMarkup       = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(string(afghanistanFlagButton), string(afghanistanFlagButton))),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(string(malaysiaFlagButton), string(malaysiaFlagButton))),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(string(japanFlagButton), string(japanFlagButton))),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(string(ukrainianFlagButton), string(ukrainianFlagButton))),
	)
)

type TelegramBot struct {
	bot *tgbotapi.BotAPI
}

type Message struct {
	Command string
	ChatId  int64
}

type Callback struct {
	ChatId  string
	Button  string
	Message *tgbotapi.Message
}

type TelegramUpdate struct {
	Message  *Message
	Callback *Callback
}

func NewTelegramBot(token string, botDebugStatus bool) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = botDebugStatus

	return &TelegramBot{bot: bot}, nil
}

func (t *TelegramBot) SendMessage(chatId int64, message string) error {
	msg := tgbotapi.NewMessage(chatId, message)
	msg.ParseMode = tgbotapi.ModeHTML
	_, err := t.bot.Send(msg)
	return err
}

func (t *TelegramBot) SendMenu(chatId int64, message string) error {
	msg := tgbotapi.NewMessage(chatId, message)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = startMenuMarkup
	_, err := t.bot.Send(msg)
	return err
}

func (t *TelegramBot) SendHoliday(queryId int64, callback Callback, holidayList string) {
	var err error

	callbackCfg := tgbotapi.NewCallback(callback.ChatId, "")
	t.bot.Send(callbackCfg)

	if err != nil {
		errorMsg := tgbotapi.NewMessage(queryId, "Please try again in few seconds")
		_, err = t.bot.Send(errorMsg)
		if err != nil {
			return
		}
		return
	}
	msg := tgbotapi.NewMessage(queryId, holidayList)
	_, err = t.bot.Send(msg)
	if err != nil {
		return
	}
}

func newUpdate(t *tgbotapi.Update) TelegramUpdate {
	update := TelegramUpdate{}
	var user *tgbotapi.User
	var text string

	if t.Message != nil {
		update.Message = &Message{
			Command: t.Message.Text,
			ChatId:  t.Message.Chat.ID,
		}
		log.WithFields(log.Fields{
			"chat_id": t.Message.Chat.ID,
			"user":    user,
			"text":    text,
		}).Info("message received")
	} else if t.CallbackQuery != nil {
		update.Callback = &Callback{
			ChatId:  t.CallbackQuery.ID,
			Button:  t.CallbackQuery.Data,
			Message: t.CallbackQuery.Message,
		}
		log.WithFields(log.Fields{
			"chat_id": t.CallbackQuery.ID,
			"user":    user,
		}).Info("message received")
	}

	return update
}

func (t *TelegramBot) GetUpdates(ctx context.Context) chan TelegramUpdate {
	u := tgbotapi.NewUpdate(0)
	updates := t.bot.GetUpdatesChan(u)
	updateChan := make(chan TelegramUpdate)
	u.Timeout = 60

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case update := <-updates:
				updateChan <- newUpdate(&update)
			}
		}
	}()

	return updateChan
}

func (t *TelegramBot) TransformListOfHolidaysToStr(holidayArray *[]string) (string, error) {
	var holidayListInString string
	if holidayArray != nil && len(*holidayArray) > 0 {
		for i := 0; i < len(*holidayArray); i++ {
			holidayListInString += (*holidayArray)[0]
		}
		return holidayListInString, nil
	} else {
		holidayListInString = ""
	}

	return "today is no holidays in this country", nil
}
