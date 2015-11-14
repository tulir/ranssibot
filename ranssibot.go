package main

import (
	"fmt"
	"github.com/tucnak/telebot"
	"log"
	"maunium.net/ranssibot/lang"
	"maunium.net/ranssibot/timetables"
	"maunium.net/ranssibot/util"
	"maunium.net/ranssibot/whitelist"
	"strings"
	"time"
)

func main() {
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
	} else if strings.HasPrefix(message.Text, "/settime") {
		/*if len(args) > 3 {
			lessonID, err := strconv.Atoi(args[2])
			if err != nil {
				bot.SendMessage(message.Chat, fmt.Sprintf(lang.Translate("error.parse.integer"), args[2]), util.Markdown)
			}
			dayShift, err := strconv.Atoi(args[3])
			if err != nil {
				bot.SendMessage(message.Chat, fmt.Sprintf(lang.Translate("error.parse.integer"), args[3]), util.Markdown)
			}
			time, err := timetables.StringToTime(args[4])
			if err != nil {
				bot.SendMessage(message.Chat, fmt.Sprintf(lang.Translate("error.parse.time"), args[4]), util.Markdown)
			}
			if args[1] == "ventit" {
				firstyear[today+dayShift][lessonID].Time = time
			} else if args[1] == "neliöt" {
				secondyear[today+dayShift][lessonID].Time = time
			} else if args[1] == "other" {
				other[today+dayShift].Time = time
			} else {
				return
			}
			bot.SendMessage(message.Chat, fmt.Sprintf(lang.Translate("settime.success"), args[1], lessonID, dayShift, TimeToString(time)), util.Markdown)
		} else {
			bot.SendMessage(message.Chat, lang.Translate("settime.usage"), util.Markdown)
		}*/
	} else if message.Text == "/update" {
		timetables.Update()
		bot.SendMessage(message.Chat, lang.Translate("timetable.update.success"), util.Markdown)
	} else if strings.HasPrefix(message.Text, "/") {
		bot.SendMessage(message.Chat, lang.Translate("error.commandnotfound"), util.Markdown)
	}
}
