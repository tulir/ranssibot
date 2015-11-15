package main

import (
	"fmt"
	"github.com/tucnak/telebot"
	"log"
	"maunium.net/ranssibot/lang"
	"maunium.net/ranssibot/laundry"
	"maunium.net/ranssibot/timetables"
	"maunium.net/ranssibot/util"
	"maunium.net/ranssibot/whitelist"
	"strings"
	"time"
)

func main() {
	laundry.NotifierTick()
	if true {
		return
	}
	lang.Load()
	util.Init()
	whitelist.Load()

	// Connect to Telegram
	bot, err := telebot.NewBot("151651579:AAErjEHJw1bNs-iWlchFwHiroULpbha_Wz8")
	if err != nil {
		log.Printf(lang.Translate("telegram.connection.failed"), err)
		return
	}
	messages := make(chan telebot.Message)
	// Enable message listener
	bot.Listen(messages, 1*time.Second)
	// Print "connected" message
	log.Printf(lang.Translate("telegram.connection.success"))

	// Update timetables
	timetables.Update()

	bot.SendMessage(whitelist.GetRecipientByName("tulir"), "Ranssibot started up @ "+time.Now().Format("15:04:05 02.01.2006"), nil)

	// Listen to messages
	for message := range messages {
		handleCommand(bot, message)
	}
}

// Handle a command
func handleCommand(bot *telebot.Bot, message telebot.Message) {
	if !whitelist.IsWhitelisted(message.Sender.ID) {
		bot.SendMessage(message.Chat, fmt.Sprintf(lang.Translate("whitelist.notwhitelisted"), message.Sender.ID), nil)
		return
	}
	args := strings.Split(message.Text, " ")
	log.Printf(lang.Translate("telegram.commandreceived"), message.Sender.Username, message.Sender.ID, message.Text)
	if strings.HasPrefix(message.Text, "Mui.") || message.Text == "/start" {
		bot.SendMessage(message.Chat, "Mui. "+message.Sender.FirstName+".", nil)
	} else if strings.HasPrefix(message.Text, "/timetable") {
		timetables.HandleCommand(bot, message, args)
	} else if message.Text == "/spamme" {
		go spam(message.Sender.ID, bot)
	} else if strings.HasPrefix(message.Text, "/") {
		bot.SendMessage(message.Chat, lang.Translate("error.commandnotfound"), util.Markdown)
	}
}

func spam(uid int, bot *telebot.Bot) {
	bot.SendMessage(whitelist.GetRecipientByUID(uid), "OK! I'll spam you in 5 seconds.", nil)
	time.Sleep(5 * time.Second)
	bot.SendMessage(whitelist.GetRecipientByUID(uid), "Here's the spam you requested!", nil)
}
