package main

import (
	"bytes"
	"github.com/tucnak/telebot"
	log "maunium.net/go/maulogger"
	"maunium.net/go/ranssibot/config"
	"maunium.net/go/ranssibot/food"
	"maunium.net/go/ranssibot/lang"
	"maunium.net/go/ranssibot/laundry"
	"maunium.net/go/ranssibot/posts"
	"maunium.net/go/ranssibot/timetables"
	"maunium.net/go/ranssibot/util"
	"strings"
)

// Handle a command
func handleCommand(bot *telebot.Bot, message telebot.Message) {
	sender := config.GetUserWithUID(message.Sender.ID)
	if sender.UID == 0 {
		bot.SendMessage(message.Chat, lang.GetLanguage("english").Translatef("whitelist.notwhitelisted", message.Sender.ID), util.Markdown)
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
		bot.SendMessage(message.Chat, "*Mui. "+message.Sender.FirstName+".*", util.Markdown)
	} else if util.CheckArgs(command, "/timetable", "/tt", "/timetables", "/tts") {
		timetables.HandleCommand(bot, message, args)
	} else if util.CheckArgs(command, "/posts", "/post", "/news") {
		posts.HandleCommand(bot, message, args)
	} else if util.CheckArgs(command, "/food", "/menu") {
		food.HandleCommand(bot, message, args)
	} else if util.CheckArgs(command, "/lang", "/language") {
		lang.HandleCommand(bot, message, args)
	} else if util.CheckArgs(command, "/laundry") {
		laundry.HandleCommand(bot, message, args)
	} else if util.CheckArgs(command, "/sauna") {
		bot.SendMessage(message.Chat, lang.Translatef(config.GetUserWithUID(message.Sender.ID), "saunatemp", strings.TrimSpace(util.HTTPGet("http://sauna.paivola.fi/saunatemp.cgi"))), util.Markdown)
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
		bot.SendMessage(message.Chat, lang.Translate(sender, "error.commandnotfound"), util.Markdown)
	}
}

func handleInstance(bot *telebot.Bot, message telebot.Message, args []string) {
	sender := config.GetUserWithUID(message.Sender.ID)
	// Create a bytebuffer to get a better performance with the instance info message.
	var buffer bytes.Buffer
	// Add the title.
	buffer.WriteString(lang.Translatef(sender, "instance.title", VersionLong))
	// Add debug mode status.
	if *debug {
		buffer.WriteString(lang.Translatef(sender, "instance.debug.active"))
	} else {
		buffer.WriteString(lang.Translatef(sender, "instance.debug.inactive"))
	}
	// Add amount of whitelisted users.
	buffer.WriteString(lang.Translatef(sender, "instance.users", len(config.GetAllUsers())))
	// Add the host machine's internal hostname.
	buffer.WriteString(lang.Translatef(sender, "instance.hostname", hostname))
	// Add the start timestamp.
	buffer.WriteString(lang.Translatef(sender, "instance.startedat", startedAt.Format("15:04:05 02.01.2006")))
	// Send the message to the user.
	bot.SendMessage(message.Chat, buffer.String(), util.Markdown)
}

func handleHelp(bot *telebot.Bot, message telebot.Message, args []string) {
	sender := config.GetUserWithUID(message.Sender.ID)
	if len(args) > 0 {
		if util.CheckArgs(args[0], "timetable", "timetables") {
			// Send help for /timetable
			bot.SendMessage(message.Chat, lang.Translate(sender, "help.timetable"), util.Markdown)
		} else if util.CheckArgs(args[0], "posts", "post", "news") {
			// Send help for /posts
			bot.SendMessage(message.Chat, lang.Translate(sender, "help.posts"), util.Markdown)
		} else if util.CheckArgs(args[0], "config", "configuration") {
			// Send help for /config
			bot.SendMessage(message.Chat, lang.Translate(sender, "help.config"), util.Markdown)
		} else if util.CheckArgs(args[0], "whitelist", "wl") {
			if len(args) > 1 {
				if util.CheckArgs(args[1], "permissions", "perms") {
					// Send help for /whitelist permissions
					bot.SendMessage(message.Chat, lang.Translate(sender, "help.whitelist.permissions"), util.Markdown)
				} else if util.CheckArgs(args[1], "settings", "preferences", "prefs", "properties", "props") {
					// Send help for /whitelist settings
					bot.SendMessage(message.Chat, lang.Translate(sender, "help.whitelist.settings"), util.Markdown)
				} else {
					// Send help for /whitelist
					bot.SendMessage(message.Chat, lang.Translate(sender, "help.whitelist"), util.Markdown)
				}
			} else {
				// Send help for /whitelist
				bot.SendMessage(message.Chat, lang.Translate(sender, "help.whitelist"), util.Markdown)
			}
		} else {
			// Unidentified help page, send standard help message.
			bot.SendMessage(message.Chat, lang.Translate(sender, "help"), util.Markdown)
		}
	} else {
		// No arguments were given, send the standard help message.
		bot.SendMessage(message.Chat, lang.Translate(sender, "help"), util.Markdown)
	}
}
