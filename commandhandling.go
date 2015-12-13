package main

import (
	"bytes"
	"github.com/tucnak/telebot"
	log "maunium.net/go/maulogger"
	"maunium.net/go/ranssibot/config"
	"maunium.net/go/ranssibot/lang"
	"maunium.net/go/ranssibot/posts"
	"maunium.net/go/ranssibot/timetables"
	"maunium.net/go/ranssibot/util"
	"strings"
)

// Handle a command
func handleCommand(bot *telebot.Bot, message telebot.Message) {
	sender := config.GetUserWithUID(message.Sender.ID)
	if sender.UID == 0 {
		bot.SendMessage(message.Chat, lang.LTranslatef("english", "whitelist.notwhitelisted", message.Sender.ID), util.Markdown)
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
		handleHelp(bot, message, args)
	} else if util.CheckArgs(command, "/instanceinfo", "/instinfo", "/insinfo", "/instance", "/info") {
		handleInstance(bot, message, args)
	} else if strings.HasPrefix(message.Text, "/") {
		bot.SendMessage(message.Chat, lang.UTranslate(sender, "error.commandnotfound"), util.Markdown)
	}
}

func handleInstance(bot *telebot.Bot, message telebot.Message, args []string) {
	sender := config.GetUserWithUID(message.Sender.ID)
	var buffer bytes.Buffer
	buffer.WriteString(lang.UTranslatef(sender, "instance.title", VersionLong))
	if *debug {
		buffer.WriteString(lang.UTranslatef(sender, "instance.debug.active"))
	} else {
		buffer.WriteString(lang.UTranslatef(sender, "instance.debug.inactive"))
	}
	buffer.WriteString(lang.UTranslatef(sender, "instance.users", len(config.GetAllUsers())))
	buffer.WriteString(lang.UTranslatef(sender, "instance.hostname", hostname))
	buffer.WriteString(lang.UTranslatef(sender, "instance.startedat", startedAt.Format("15:04:05 02.01.2006")))
	bot.SendMessage(message.Chat, buffer.String(), util.Markdown)
}

func handleHelp(bot *telebot.Bot, message telebot.Message, args []string) {
	sender := config.GetUserWithUID(message.Sender.ID)
	if len(args) == 0 {
		bot.SendMessage(message.Chat, lang.UTranslate(sender, "help"), util.Markdown)
	} else if len(args) > 0 {
		if util.CheckArgs(args[0], "timetable", "timetables") {
			bot.SendMessage(message.Chat, lang.UTranslate(sender, "help.timetable"), util.Markdown)
		} else if util.CheckArgs(args[0], "posts", "post") {
			bot.SendMessage(message.Chat, lang.UTranslate(sender, "help.posts"), util.Markdown)
		} else if util.CheckArgs(args[0], "config", "configuration") {
			bot.SendMessage(message.Chat, lang.UTranslate(sender, "help.config"), util.Markdown)
		} else if util.CheckArgs(args[0], "whitelist", "wl") {
			if len(args) > 1 {
				if util.CheckArgs(args[1], "permissions", "perms") {
					bot.SendMessage(message.Chat, lang.UTranslate(sender, "help.whitelist.permissions"), util.Markdown)
				} else if util.CheckArgs(args[1], "settings", "preferences", "prefs", "properties", "props") {
					bot.SendMessage(message.Chat, lang.UTranslate(sender, "help.whitelist.settings"), util.Markdown)
				} else {
					bot.SendMessage(message.Chat, lang.UTranslate(sender, "help.whitelist"), util.Markdown)
				}
			} else {
				bot.SendMessage(message.Chat, lang.UTranslate(sender, "help.whitelist"), util.Markdown)
			}
		} else {
			bot.SendMessage(message.Chat, lang.UTranslate(sender, "help.usage"), util.Markdown)
		}
	} else {
		bot.SendMessage(message.Chat, lang.UTranslate(sender, "help.usage"), util.Markdown)
	}
}
