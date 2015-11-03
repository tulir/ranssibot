package main

import (
	"fmt"
	"github.com/tucnak/telebot"
	"log"
	"strconv"
	"strings"
	"time"
)

// The day ID for today
var today = 5

// The markdown send options
var md *telebot.SendOptions

func main() {
	loadLanguage()

	md = new(telebot.SendOptions)
	md.ParseMode = telebot.ModeMarkdown

	// Load the whitelist
	loadWhitelist()

	// Connect to Telegram
	bot, err := telebot.NewBot("170943030:AAE8O_pJ2nHeWCTDTHOBD3Wy-ryFmNkxOOQ")
	if err != nil {
		log.Printf(translate("telegram.connection.failed"), err)
		return
	}
	messages := make(chan telebot.Message)
	// Enable message listener
	bot.Listen(messages, 1*time.Second)
	// Print "connected" message
	log.Printf(translate("telegram.connection.success"))

	// Update timetables
	updateTimes()

	// Listen to messages
	for message := range messages {
		handleCommand(bot, message)
	}
}

// Handle a command
func handleCommand(bot *telebot.Bot, message telebot.Message) {
	if !isWhitelisted(message.Sender.ID) {
		bot.SendMessage(message.Chat, fmt.Sprintf(translate("whitelist.notwhitelisted"), message.Sender.ID), nil)
		return
	}
	args := strings.Split(message.Text, " ")
	log.Printf(translate("telegram.commandreceived"), message.Sender.Username, message.Sender.ID, message.Text)
	if strings.HasPrefix(message.Text, "Mui.") || message.Text == "/start" {
		bot.SendMessage(message.Chat, "Mui. "+message.Sender.FirstName+".", nil)
	} else if strings.HasPrefix(message.Text, "/timetable") {
		if timestamp() > lastupdate+600 {
			bot.SendMessage(message.Chat, "Updating cached timetables...", md)
			updateTimes()
		}
		if len(args) > 1 {
			day := today
			if len(args) > 2 {
				shift, err := strconv.Atoi(args[2])
				if err != nil {
					bot.SendMessage(message.Chat, fmt.Sprintf(translate("error.parse.integer"), args[2]), md)
					return
				}
				day += shift
				if day < 0 || day >= len(other) {
					bot.SendMessage(message.Chat, translate("timetable.nodata"), md)
					return
				}
			}
			if strings.EqualFold(args[1], "ventit") {
				sendFirstYear(day, bot, message)
			} else if strings.EqualFold(args[1], "neliöt") {
				sendSecondYear(day, bot, message)
			} else {
				bot.SendMessage(message.Chat, translate("timetable.usage"), md)
			}
		} else {
			year := getYeargroupIndex(message.Sender.ID)
			if year == 1 {
				sendFirstYear(today, bot, message)
			} else if year == 2 {
				sendSecondYear(today, bot, message)
			} else {
				bot.SendMessage(message.Chat, translate("timetable.noyeargroup"), md)
			}
		}
	} else if strings.HasPrefix(message.Text, "/settime") {
		if len(args) > 3 {
			lessonID, err := strconv.Atoi(args[2])
			if err != nil {
				bot.SendMessage(message.Chat, fmt.Sprintf(translate("error.parse.integer"), args[2]), md)
			}
			dayShift, err := strconv.Atoi(args[3])
			if err != nil {
				bot.SendMessage(message.Chat, fmt.Sprintf(translate("error.parse.integer"), args[3]), md)
			}
			time, err := StringToTime(args[4])
			if err != nil {
				bot.SendMessage(message.Chat, fmt.Sprintf(translate("error.parse.time"), args[4]), md)
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
			bot.SendMessage(message.Chat, fmt.Sprintf(translate("settime.success"), args[1], lessonID, dayShift, TimeToString(time)), md)
		} else {
			bot.SendMessage(message.Chat, translate("settime.usage"), md)
		}
	} else if message.Text == "/update" {
		updateTimes()
		bot.SendMessage(message.Chat, translate("timetable.update.success"), md)
	} else if strings.HasPrefix(message.Text, "/") {
		bot.SendMessage(message.Chat, translate("error.commandnotfound"), md)
	}
}

// Get the current UNIX timestamp
func timestamp() int64 {
	return int64(time.Now().Unix())
}
