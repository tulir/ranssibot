package main

import (
	"bufio"
	"github.com/tucnak/telebot"
	log "maunium.net/maulogger"
	"maunium.net/ranssibot/util"
	"maunium.net/ranssibot/whitelist"
	"os"
	"strconv"
	"strings"
)

func listen(bot *telebot.Bot) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-1]

		args := strings.Split(text, " ")

		onCommand(bot, strings.ToLower(args[0]), args[1:])
	}
}

func onCommand(bot *telebot.Bot, command string, args []string) {
	if command == "msg" {
		user := whitelist.GetUserWithName(args[0])
		if user.Destination() == 0 {
			i, err := strconv.Atoi(args[0])
			if err != nil {
				log.Errorf("Couldn't find an integer or a whitelisted user from %s", args[0])
				return
			}
			user = whitelist.GetUserWithUID(i)
		}

		msg := connect(args[1:])
		bot.SendMessage(user, "*[Sysadmin]* "+msg, util.Markdown)
		log.Infof("Sent message %[1]s to %[2]s", msg, args[0])
	} else if command == "stop" {
		Shutdown()
	}
}

func connect(array []string) string {
	var str string
	for _, val := range array {
		str = str + " " + val
	}
	return str[1:]
}
