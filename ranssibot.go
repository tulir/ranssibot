package main

import (
	"bytes"
	"fmt"
	flag "github.com/ogier/pflag"
	"github.com/tucnak/telebot"
	"io/ioutil"
	log "maunium.net/go/maulogger"
	"maunium.net/go/ranssibot/config"
	"maunium.net/go/ranssibot/lang"
	"maunium.net/go/ranssibot/posts"
	"maunium.net/go/ranssibot/timetables"
	"maunium.net/go/ranssibot/util"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// VersionLong is the human-readable form of the version.
const VersionLong = "0.1 Beta 1"

// Version is the computer-readable form of the version.
const Version = "0.1-B1"

var startedAt time.Time
var hostname string

var token = flag.StringP("token", "t", "", "The Telegram bot token to use.")
var debug = flag.BoolP("debug", "d", false, "Enable debug mode")
var bot *telebot.Bot
var disableSafeShutdown = flag.Bool("no-safe-shutdown", false, "Disable Interrupt/SIGTERM catching and handling.")

func init() {
	flag.Parse()

	log.PrintDebug = *debug
	log.Fileformat = "logs/%[1]s-%02[2]d.log"
	log.Init()
	lang.Load()
	config.IndentConfig = *debug
	config.Load()

	if !*disableSafeShutdown {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			Shutdown("Interrupt/SIGTERM")
		}()
	}

	data, err := ioutil.ReadFile("/etc/hostname")
	if err != nil {
		log.Fatalln("Failed to read hostname: %s", err)
		return
	}
	hostname = strings.TrimSpace(string(data))
}

func main() {
	start := util.TimestampMS()
	// Connect to Telegram
	var err error
	bot, err = telebot.NewBot(*token)
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

	go posts.Loop(bot, *debug)
	go listen(bot)

	startedAt = time.Now()

	var startup string
	if *debug {
		startup = fmt.Sprintf("Ranssibot started up in %[1]dms @ %[2]s (Debug mode)", util.TimestampMS()-start, startedAt.Format("15:04:05 02.01.2006"))
	} else {
		startup = fmt.Sprintf("Ranssibot started up @ %[1]s", startedAt.Format("15:04:05 02.01.2006"))
	}

	bot.SendMessage(config.GetUserWithName("tulir"), startup, nil)
	log.Infof(startup)

	// Listen to messages
	for message := range messages {
		handleCommand(bot, message)
	}
}

// Handle a command
func handleCommand(bot *telebot.Bot, message telebot.Message) {
	if config.GetUserWithUID(message.Sender.ID).UID == 0 {
		bot.SendMessage(message.Chat, lang.Translatef("whitelist.notwhitelisted", message.Sender.ID), util.Markdown)
		return
	}
	args := strings.Split(message.Text, " ")
	command := args[0]
	args = args[1:]
	if message.Chat.IsGroupChat() {
		log.Infof("%[1]s (%[2]d) @ %[3]s (%[4]d) sent command: %[3]s", message.Sender.Username, message.Sender.ID, message.Chat.Title, message.Chat.ID, message.Text)
	} else {
		log.Infof("%[1]s (%[2]d) sent command: %[3]s", message.Sender.Username, message.Sender.ID, message.Text)
	}
	if strings.HasPrefix(message.Text, "Mui.") || message.Text == "/start" {
		bot.SendMessage(message.Chat, "Mui. "+message.Sender.FirstName+".", nil)
	} else if util.CheckArgs(command, "/timetable", "/tt", "/timetables", "/tts") {
		timetables.HandleCommand(bot, message, args)
	} else if util.CheckArgs(command, "/posts", "/post") {
		posts.HandleCommand(bot, message, args)
	} else if util.CheckArgs(command, "/lang", "/language") {
		lang.HandleCommand(bot, message, args)
	} else if util.CheckArgs(command, "/config", "/configuration") {
		handleConfig(bot, message, args)
	} else if util.CheckArgs(command, "/stop", "/shutdown", "/poweroff") {
		handleStop(bot, message, args)
	} else if util.CheckArgs(command, "/whitelist", "/wl") {
		handleWhitelist(bot, message, args)
	} else if util.CheckArgs(command, "/help", "help", "?", "/?") {
		if len(args) == 0 {
			bot.SendMessage(message.Chat, lang.Translate("help"), util.Markdown)
		} else if len(args) > 0 {
			if util.CheckArgs(args[0], "timetable", "timetables") {
				bot.SendMessage(message.Chat, lang.Translate("help.timetable"), util.Markdown)
			} else if util.CheckArgs(args[0], "posts", "post") {
				bot.SendMessage(message.Chat, lang.Translate("help.posts"), util.Markdown)
			} else if util.CheckArgs(args[0], "config", "configuration") {
				bot.SendMessage(message.Chat, lang.Translate("help.config"), util.Markdown)
			} else if util.CheckArgs(args[0], "whitelist", "wl") {
				if len(args) > 1 {
					if util.CheckArgs(args[1], "permissions", "perms") {
						bot.SendMessage(message.Chat, lang.Translate("help.whitelist.permissions"), util.Markdown)
					} else if util.CheckArgs(args[1], "settings", "preferences", "prefs", "properties", "props") {
						bot.SendMessage(message.Chat, lang.Translate("help.whitelist.settings"), util.Markdown)
					} else {
						bot.SendMessage(message.Chat, lang.Translate("help.whitelist"), util.Markdown)
					}
				} else {
					bot.SendMessage(message.Chat, lang.Translate("help.whitelist"), util.Markdown)
				}
			} else {
				bot.SendMessage(message.Chat, lang.Translate("help.usage"), util.Markdown)
			}
		} else {
			bot.SendMessage(message.Chat, lang.Translate("help.usage"), util.Markdown)
		}
	} else if util.CheckArgs(command, "/instanceinfo", "/instinfo", "/insinfo", "/instance", "/info") {
		var buffer bytes.Buffer
		buffer.WriteString(lang.Translatef("instance.title", VersionLong))
		if *debug {
			buffer.WriteString(lang.Translatef("instance.debug.active"))
		} else {
			buffer.WriteString(lang.Translatef("instance.debug.inactive"))
		}
		buffer.WriteString(lang.Translatef("instance.users", len(config.GetAllUsers())))
		buffer.WriteString(lang.Translatef("instance.hostname", hostname))
		buffer.WriteString(lang.Translatef("instance.startedat", startedAt.Format("15:04:05 02.01.2006")))
		bot.SendMessage(message.Chat, buffer.String(), util.Markdown)
	} else if strings.HasPrefix(message.Text, "/") {
		bot.SendMessage(message.Chat, lang.Translate("error.commandnotfound"), util.Markdown)
	}
}

// Shutdown shuts down the Ranssibot.
func Shutdown(by string) {
	log.Infof("Ranssibot cleaning up and exiting...")
	config.Save()
	log.Shutdown()

	shutdown := fmt.Sprintf("Ranssibot shut down by %[2]s @ %[1]s", time.Now().Format("15:04:05 02.01.2006"), by)
	log.Debugf(shutdown)
	bot.SendMessage(config.GetUserWithName("tulir"), shutdown, nil)

	os.Exit(0)
}
