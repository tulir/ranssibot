package main

import (
	"fmt"
	"github.com/tucnak/telebot"
	"golang.org/x/net/html"
	"log"
	"strconv"
	"strings"
	"time"
)

// Timetable cahce
var timetable = [26][9]string{}

// The day ID for today
var today = 5

// The last time the timetable cache was updated
var lastupdate = timestamp()

// The markdown send options
var md *telebot.SendOptions

func main() {
	loadLanguage()

	md = new(telebot.SendOptions)
	md.ParseMode = telebot.ModeMarkdown

	// Load the whitelist
	loadWhitelist()

	// Connect to Telegram
	bot, err := telebot.NewBot("132300126:AAHps1NPAj9Y7qTBbDGlGsyuMGoMtk__Qa8")
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
		bot.SendMessage(message.Chat, "", nil)
		return
	}
	log.Printf(translate("telegram.commandreceived"), message.Sender.Username, message.Sender.ID, message.Text)
	if strings.HasPrefix(message.Text, "Mui.") {
		bot.SendMessage(message.Chat, "Mui. "+message.Sender.FirstName+".", nil)
	} else if strings.HasPrefix(message.Text, "/timetable") {
		if timestamp() > lastupdate+600 {
			bot.SendMessage(message.Chat, "Updating cached timetables...", md)
			updateTimes()
		}
		args := strings.Split(message.Text, " ")
		if len(args) > 1 {
			day := today
			if len(args) > 2 {
				shift, err := strconv.Atoi(args[2])
				if err != nil {
					bot.SendMessage(message.Chat, fmt.Sprintf(translate("error.parse.integer"), args[2]), md)
					return
				}
				day += shift
				if day < 0 || day >= len(timetable) {
					bot.SendMessage(message.Chat, translate("timetable.nodata"), md)
					return
				}
			}
			if strings.EqualFold(args[1], "ventit") {
				bot.SendMessage(message.Chat,
					fmt.Sprintf(translate("timetable.generic"), timetable[day][0], timetable[day][1], timetable[day][2], timetable[day][3])+
						"\n"+fmt.Sprintf(translate("timetable.other"), timetable[day][4]),
					md)
			} else if strings.EqualFold(args[1], "neliöt") {
				bot.SendMessage(message.Chat,
					fmt.Sprintf(translate("timetable.generic"), timetable[day][5], timetable[day][6], timetable[day][7], timetable[day][8])+
						"\n"+fmt.Sprintf(translate("timetable.other"), timetable[day][4]),
					md)
			} else {
				bot.SendMessage(message.Chat, translate("timetable.usage"), md)
			}
		} else {
			bot.SendMessage(message.Chat, translate("timetable.firstyear")+fmt.Sprintf(translate("timetable.generic"),
				timetable[today][0], timetable[today][1], timetable[today][2], timetable[today][3]), md)
			bot.SendMessage(message.Chat, translate("timetable.secondyear")+fmt.Sprintf(translate("timetable.generic"),
				timetable[today][5], timetable[today][6], timetable[today][7], timetable[today][8]), md)
			bot.SendMessage(message.Chat, fmt.Sprintf(translate("timetable.other"), timetable[today][4]), md)
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

// Update the timetables from http://ranssi.paivola.fi/lj.php
func updateTimes() {
	// Get the timetable page and convert the string to a reader
	reader := strings.NewReader(httpGet("http://ranssi.paivola.fi/lj.php"))
	// Parse the HTML from the reader
	doc, err := html.Parse(reader)
	// Check if there was an error
	if err != nil {
		// Print the error
		log.Printf("%s", err)
		// Return
		return
	}

	// Find the timetable table header node
	ttnode := findSpan("tr", "class", "header", doc)
	// Check if the node was found
	if ttnode != nil {
		dayentry := ttnode
		// Loop through the days in the timetable
		for day := 0; ; day++ {
			// Make sure the next day exists
			if dayentry.NextSibling == nil ||
				dayentry.NextSibling.NextSibling == nil ||
				dayentry.NextSibling.NextSibling.FirstChild == nil {
				break
			}
			// Get the next day node
			dayentry = dayentry.NextSibling.NextSibling
			// Get the first lesson node in the day node
			entry := dayentry.FirstChild.NextSibling
			// Loop through the lessons on the day
			for lesson := 0; lesson < 9; lesson++ {
				// Make sure the next lesson exists
				if entry == nil ||
					entry.NextSibling == nil ||
					entry.NextSibling.NextSibling == nil {
					break
				}
				// Get the next lesson node
				entry = entry.NextSibling.NextSibling

				// Check if the lesson contains anything
				if entry.FirstChild != nil {
					// Lesson data found; Try to parse it
					if entry.FirstChild.Type == html.TextNode {
						// Found lesson data directly under lesson node
						timetable[day][lesson] = entry.FirstChild.Data
					} else if entry.FirstChild.Type == html.ElementNode {
						// Didn't find data directly under lesson node
						// Check for a child element node.
						if entry.FirstChild.FirstChild != nil {
							// Child element node found. Check if the child of that child is text.
							if entry.FirstChild.FirstChild.Type == html.TextNode {
								// Child of child is text, use it as the data.
								timetable[day][lesson] = entry.FirstChild.FirstChild.Data
							}
						}
					} else {
						// Lesson data couldn't be parsed
						timetable[day][lesson] = "tyhjää"
					}
				} else {
					// Lesson is empty
					timetable[day][lesson] = "tyhjää"
				}
			}
		}
		lastupdate = timestamp()
	} else {
		// Node not found, print error
		log.Printf(translate("timetable.update.failed"))
		lastupdate = 0
	}
}
