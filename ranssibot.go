package main

import (
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
var disableSafeShutdown = flag.Bool("no-safe-shutdown", false, "Disable Interrupt/SIGTERM catching and handling.")

var bot *telebot.Bot

func init() {
	flag.Parse()

	log.PrintDebug = *debug
	log.Fileformat = func(now string, i int) string { return fmt.Sprintf("logs/%[1]s-%02[2]d.log", now, i) }
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
		startup = fmt.Sprintf("Ranssibot started up")
	}

	bot.SendMessage(config.GetUserWithName("tulir"), startup, nil)
	log.Infof(startup)

	// Listen to messages
	for message := range messages {
		handleCommand(bot, message)
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
