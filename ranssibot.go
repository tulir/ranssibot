package main

import (
	"fmt"
	"github.com/tucnak/telebot"
	log "maunium.net/maulogger"
	"maunium.net/ranssibot/lang"
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

var token = flag.StringP("token", "t", "", "The Telegram bot token to use.")
var debug = flag.BoolP("debug", "d", false, "Enable debug mode")
var disableSafeShutdown = flag.Bool("no-safe-shutdown", false, "Disable Interrupt/SIGTERM catching and handling.")

func init() {
	flag.Parse()

	log.PrintDebug = *debug
	log.Fileformat = "logs/%[1]s-%[2]d.log"
	log.Init()
	lang.Init()
	whitelist.Load()

	if !*disableSafeShutdown {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			Shutdown()
		}()
	}
}

func main() {
	start := util.TimestampMS()
	// Connect to Telegram
	bot, err := telebot.NewBot(*token)
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
	go listen(bot)

	var startup string
	if *debug {
		startup = fmt.Sprintf("Ranssibot started up in %[1]dms @ %[2]s (Debug mode)", util.TimestampMS()-start, time.Now().Format("15:04:05 02.01.2006"))
	} else {
		startup = fmt.Sprintf("Ranssibot started up @ %[1]s", time.Now().Format("15:04:05 02.01.2006"))
	}

	bot.SendMessage(whitelist.GetUserWithName("tulir"), startup, nil)
	log.Infof(startup)

	// Listen to messages
	for message := range messages {
		handleCommand(bot, message)
	}
}

// Handle a command
func handleCommand(bot *telebot.Bot, message telebot.Message) {
	if whitelist.GetUserWithUID(message.Sender.ID).UID != 0 {
		bot.SendMessage(message.Chat, lang.Translatef("whitelist.notwhitelisted", message.Sender.ID), util.Markdown)
		return
	}
	args := strings.Split(message.Text, " ")[1:]
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
	} else if strings.HasPrefix(message.Text, "/help") {
		if len(args) == 0 {
			bot.SendMessage(message.Chat, lang.Translate("help"), util.Markdown)
		} else if len(args) > 0 {
			if strings.EqualFold(args[0], "timetable") {
				bot.SendMessage(message.Chat, lang.Translate("help.timetable"), util.Markdown)
			} else if strings.EqualFold(args[0], "posts") {
				bot.SendMessage(message.Chat, lang.Translate("help.posts"), util.Markdown)
			} else {
				bot.SendMessage(message.Chat, lang.Translate("help.usage"), util.Markdown)
			}
		} else {
			bot.SendMessage(message.Chat, lang.Translate("help.usage"), util.Markdown)
		}
	} else if strings.HasPrefix(message.Text, "/") {
		bot.SendMessage(message.Chat, lang.Translate("error.commandnotfound"), util.Markdown)
	}
}

// Shutdown shuts down the Ranssibot.
func Shutdown() {
	log.Infof("Ranssibot cleaning up and exiting...")
	whitelist.Save(*debug)
	log.Shutdown()
	os.Exit(0)
}
