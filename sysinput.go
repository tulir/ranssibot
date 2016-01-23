package main

import (
	"bufio"
	"bytes"
	"github.com/tucnak/telebot"
	log "maunium.net/go/maulogger"
	"maunium.net/go/ranssibot/config"
	"maunium.net/go/ranssibot/util"
	"os"
	"strings"
)

func listen(bot *telebot.Bot) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-1]

		log.Debugf("[Syscmds] Received %s", text)

		args := strings.Split(text, " ")

		onCommand(bot, strings.ToLower(args[0]), args[1:])
	}
}

func onCommand(bot *telebot.Bot, command string, args []string) {
	if command == "msg" && len(args) > 1 {
		user := config.GetUser(args[0])
		if user.UID == config.NilUser.UID {
			log.Errorf("[Syscmds] Couldn't get an user with the name or UID %s", args[0])
		}

		msg := connect(args[1:])
		bot.SendMessage(user, "*[Sysadmin]* "+msg, util.Markdown)
		log.Infof("[Syscmds] Sent message %[1]s to %[2]s", msg, user.Name)
	} else if command == "broadcast" && len(args) > 0 {
		msg := connect(args)
		for _, user := range config.GetAllUsers() {
			bot.SendMessage(user, "*[Sysadmin Broadcast]* "+msg, util.Markdown)
		}
		log.Infof("[Syscmds] Broadcasted message %[1]s", msg)
	} else if command == "config" && len(args) > 0 {
		if strings.EqualFold(args[0], "save") {
			if !config.IndentConfig && len(args) > 1 && strings.EqualFold(args[0], "pretty") {
				config.IndentConfig = true
				config.Save()
				config.IndentConfig = false
			} else {
				config.Save()
			}
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
