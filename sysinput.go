package main

import (
	"bufio"
	"bytes"
	"github.com/tucnak/telebot"
	log "maunium.net/go/maulogger"
	"maunium.net/go/ranssibot/config"
	"maunium.net/go/ranssibot/util"
	"os"
	"strconv"
	"strings"
)

func listen(bot *telebot.Bot) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-1]

		log.Debugf("Sysinput: %s", text)

		args := strings.Split(text, " ")

		onCommand(bot, strings.ToLower(args[0]), args[1:])
	}
}

func onCommand(bot *telebot.Bot, command string, args []string) {
	if command == "msg" && len(args) > 1 {
		user := config.GetUserWithName(args[0])
		if user.Destination() == 0 {
			i, err := strconv.Atoi(args[0])
			if err != nil {
				log.Errorf("Couldn't find an integer or a whitelisted user from %s", args[0])
				return
			}
			user = config.GetUserWithUID(i)
		}

		msg := connect(args[1:])
		bot.SendMessage(user, "*[Sysadmin]* "+msg, util.Markdown)
		log.Infof("Sent message %[1]s to %[2]s", msg, args[0])
	} else if command == "broadcast" && len(args) > 0 {
		msg := connect(args)
		for _, user := range config.GetAllUsers() {
			bot.SendMessage(user, "*[Sysadmin Broadcast]* "+msg, util.Markdown)
		}
		log.Infof("Broadcasted message %[1]s", msg)
	} else if command == "config" && len(args) > 0 {
		if strings.EqualFold(args[0], "save") {
			config.Save()
		} else if strings.EqualFold(args[0], "load") {
			config.Load()
		}
	} else if command == "stop" {
		Shutdown("Sysadmin")
	}
}

func connect(array []string) string {
	var buffer bytes.Buffer
	for _, val := range array {
		buffer.WriteString(" ")
		buffer.WriteString(val)
	}
	return buffer.String()[1:]
}
