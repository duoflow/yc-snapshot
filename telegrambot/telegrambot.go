package telegrambot

import (
	"strconv"

	"github.com/duoflow/yc-snapshot/loggers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	// Tgbot - telegram bot
	Tgbot Telegbot
)

// Telegbot  - telegram bot interface
type Telegbot struct {
	// Telegbot - telegram bot api interface
	bot *tgbotapi.BotAPI
	// ChatID - Chat ID for message sending
	ChatID int64
	// InitFlag - if Init func was called
	InitFlag bool
}

// Init - Initializing of bot
func Init(tgtoken string) {
	Tgbot.ChatID = 185222660
	Tgbot.InitFlag = false
	tg, err := tgbotapi.NewBotAPI(tgtoken)
	if err != nil {
		loggers.Error.Printf("Telegram Init() Error while Telegram bot init: %s", err.Error())
		Tgbot.bot = nil
	} else {
		Tgbot.bot = tg
		Tgbot.InitFlag = true
		loggers.Info.Println("Telegram Init() Telegram API initialised")
	}
}

// Serve - function for message exchange
func (t Telegbot) Serve() {
	if t.InitFlag == false {
		return
	}
	// channel initialization for updates from API
	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	uchannel, err := t.bot.GetUpdatesChan(ucfg)
	if err != nil {
		loggers.Error.Printf("Telegram Serve() Error while Telegram bot init: %s", err.Error())
	}
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
				t.bot.Send(msg)
			}
		}
	}
	// end
}

// SendMessage - send message to admin
func (t Telegbot) SendMessage(m string) {
	if t.InitFlag == false {
		return
	}
	msg := tgbotapi.NewMessage(t.ChatID, m)
	t.bot.Send(msg)
}
