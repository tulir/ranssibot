package main

import (
	"fmt"
	"github.com/tucnak/telebot"
	"maunium.net/ranssibot/lang"
	log "maunium.net/ranssibot/maulog"
	//"maunium.net/ranssibot/laundry"
	flag "github.com/ogier/pflag"
	"maunium.net/ranssibot/posts"
	"maunium.net/ranssibot/timetables"
	"maunium.net/ranssibot/util"
	"maunium.net/ranssibot/whitelist"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func init() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Infof("Ranssibot interrupted. Cleaning up and exiting...")
		log.Shutdown()
		os.Exit(1)
	}()
	flag.Parse()
}

func main() {
	// Connect to Telegram
	bot, err := telebot.NewBot("151651579:AAErjEHJw1bNs-iWlchFwHiroULpbha_Wz8")
	if err != nil {
		log.Fatalf("Error connecting to Telegram: %[1]s", err)
		return
	}
	messages := make(chan telebot.Message)
	// Enable message listener
	bot.Listen(messages, 1*time.Second)
	// Print "connected" message
	log.Infof("Successfully connected to Telegram!")

	// Update timetables
	timetables.Update()

	//go laundry.NotifierTick()
	go posts.Loop(bot)

	startup := fmt.Sprintf("Ranssibot started up @ %s", time.Now().Format("15:04:05 02.01.2006"))
	bot.SendMessage(whitelist.GetRecipientByName("tulir"), startup, nil)
	log.Infof(startup)

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
	if message.Chat.IsGroupChat() {
		log.Infof("%[1]s (%[2]d) @ %[3]s (%[4]d) sent command: %[3]s", message.Sender.Username, message.Sender.ID, message.Chat.Title, message.Chat.ID, message.Text)
	} else {
		log.Infof("%[1]s (%[2]d) sent command: %[3]s", message.Sender.Username, message.Sender.ID, message.Text)
	}
	if strings.HasPrefix(message.Text, "Mui.") || message.Text == "/start" {
		bot.SendMessage(message.Chat, "Mui. "+message.Sender.FirstName+".", nil)
	} else if strings.HasPrefix(message.Text, "/timetable") {
		timetables.HandleCommand(bot, message, args)
	} else if strings.HasPrefix(message.Text, "/posts") {
		posts.HandleCommand(bot, message, args)
	} else if strings.HasPrefix(message.Text, "/") {
		bot.SendMessage(message.Chat, lang.Translate("error.commandnotfound"), util.Markdown)
	}
}
