package telegrambot

import (
	"log"

	"github.com/duoflow/yc-snapshot/loggers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	// Telegbot - telegram bot api interface
	Telegbot *tgbotapi.BotAPI
	// ChatID - Chat ID for message sending
	ChatID int64
)

// Init - Initializing of bot
func Init(tgtoken string) {
	Telegbot, err := tgbotapi.NewBotAPI(tgtoken)
	if err != nil {
		loggers.Error.Printf("Error while Telegram bot init: %s", err.Error())
	}
	//
	Telegbot.Debug = true
	loggers.Info.Printf("Authorized on account %s", Telegbot.Self.UserName)
	//
	// channel initialization for updates from API
	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	uchannel, err := Telegbot.GetUpdatesChan(ucfg)
	// read update channel messages
	for {
		select {
		case update := <-uchannel:
			// get Username who sent the message
			UserName := update.Message.From.UserName
			// get Chat ID
			ChatID := update.Message.Chat.ID
			// get text message
			Text := update.Message.Text
			loggers.Info.Printf("Telegram bot [%s] %d %s", UserName, ChatID, Text)
			// compose reply text
			reply := Text
			// compose reply message (chat ID + text)
			msg := tgbotapi.NewMessage(ChatID, reply+"ChatID: %s"+string(ChatID))
			// send message back
			Telegbot.Send(msg)
		}
	}
}

// Initv2 - test fucntion
func Initv2(token string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
