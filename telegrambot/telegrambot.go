package telegrambot

import (
	"strconv"

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
	//Telegbot.Debug = true
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
			// check if message is command
			if update.Message.Chat.ID == 185222660 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				if update.Message.IsCommand() {
					switch update.Message.Command() {
					case "register":
						msg.Text = "Thank you. I've been registered"
					case "unregister":
						msg.Text = "Well...Ok, you've been unregistered"
					case "status":
						msg.Text = "I'm ok."
					default:
						msg.Text = "Whazz up man! How are you?"
					}
				} else {
					// get Username who sent the message
					UserName := update.Message.From.UserName
					// get Chat ID
					ChatID := update.Message.Chat.ID
					// get text message
					Text := update.Message.Text
					loggers.Info.Printf("Telegram bot [%s] %d %s", UserName, ChatID, Text)
					// compose reply text
					msg.Text = Text + "   Reply for ChatID = " + strconv.FormatInt(ChatID, 10)
				}
				// send message back
				Telegbot.Send(msg)
			}
		}
	}
}
