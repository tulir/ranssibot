package main

import (
	"bytes"
	"fmt"
	"github.com/tucnak/telebot"
	"maunium.net/go/ranssibot/config"
	"maunium.net/go/ranssibot/lang"
	"maunium.net/go/ranssibot/util"
	"strconv"
)

func handleStop(bot *telebot.Bot, message telebot.Message, args []string) {
	if !checkPerms(bot, message.Sender.ID, "server.stop") {
		return
	}
	Shutdown(fmt.Sprintf("%[1]s %[2]s (ID %[3]d)", message.Sender.FirstName, message.Sender.LastName, message.Sender.ID))
}

func handleWhitelist(bot *telebot.Bot, message telebot.Message, args []string) {
	if len(args) > 0 {
		if util.CheckArgs(args[0], "add") {
			if !checkPerms(bot, message.Sender.ID, "whitelist.add") {
				return
			}
			if len(args) > 3 {
				uid, err := strconv.Atoi(args[1])
				if err != nil {
					bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.add.parseint", args[1]), util.Markdown)
					return
				}
				year, err := strconv.Atoi(args[3])
				if err != nil {
					bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.add.parseint", args[3]), util.Markdown)
					return
				}
				if config.AddUser(config.CreateUser(uid, args[2], year)) {
					bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.add", uid, args[2], year), util.Markdown)
				} else {
					bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.add.alreadyused", uid, args[2]), util.Markdown)
				}
			} else {
				bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.add.usage"), util.Markdown)
			}
		} else if util.CheckArgs(args[0], "remove", "delete", "rm", "del") {
			if !checkPerms(bot, message.Sender.ID, "whitelist.remove") {
				return
			}
			if len(args) > 1 {
				// TODO: Whitelist remove
			} else {
				bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.remove.usage"), util.Markdown)
			}
		} else if util.CheckArgs(args[0], "get") {
			if !checkPerms(bot, message.Sender.ID, "whitelist.get") {
				return
			}
			if len(args) > 1 {
				// TODO: Whitelist get
			} else {
				bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.get.usage"), util.Markdown)
			}
		} else if util.CheckArgs(args[0], "list") {
			if !checkPerms(bot, message.Sender.ID, "whitelist.list") {
				return
			}

			var buffer bytes.Buffer
			for _, user := range config.GetAllUsers() {
				buffer.WriteString(lang.Translatef("mgmt.whitelist.list.entry", user.Name, user.UID))
			}
			bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.list", buffer.String()), util.Markdown)
		} else if util.CheckArgs(args[0], "permissions", "perms") {
			if !handleWhitelistPerms(bot, message, args[1:]) {
				bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.permissions.usage"), util.Markdown)
			}
		} else if util.CheckArgs(args[0], "settings", "properties", "props", "preferences", "prefs") {
			if !handleWhitelistSettings(bot, message, args[1:]) {
				bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.settings.usage"), util.Markdown)
			}
		} else {
			bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.usage"), util.Markdown)
		}
	} else {
		bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.usage"), util.Markdown)
	}
}

func handleWhitelistSettings(bot *telebot.Bot, message telebot.Message, args []string) bool {
	if len(args) < 2 {
		return false
	}
	user := config.GetUser(args[1])
	if user.UID == config.NilUser.UID {
		bot.SendMessage(message.Chat, lang.Translatef("mgmt.error.usernotfound", args[1]), util.Markdown)
		return true
	}

	if util.CheckArgs(args[0], "view", "see", "list") {
		if !checkPerms(bot, message.Sender.ID, "whitelist.settings.view") {
			return true
		}
		if len(user.GetSettings()) == 0 {
			bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.settings.view.empty", user.Name), util.Markdown)
		} else {
			var buffer bytes.Buffer
			for key, val := range user.GetSettings() {
				buffer.WriteString(lang.Translatef("mgmt.whitelist.settings.view.entry", key, val))
			}
			bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.settings.view", user.Name, buffer.String()), util.Markdown)
		}
	} else if util.CheckArgs(args[0], "get") {
		if !checkPerms(bot, message.Sender.ID, "whitelist.settings.get") {
			return true
		}
		if len(args) > 2 {
			val, ok := user.GetSetting(args[2])
			if !ok {
				bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.settings.get.notfound", user.Name, args[2]), util.Markdown)
			} else {
				bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.settings.get", user.Name, args[2], val), util.Markdown)
			}
		} else {
			bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.settings.get.usage"), util.Markdown)
		}
	} else if util.CheckArgs(args[0], "add", "set") {
		if !checkPerms(bot, message.Sender.ID, "whitelist.settings.add") {
			return true
		}
		if len(args) > 3 {
			user.SetSetting(args[2], args[3])
			bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.settings.set", user.Name, args[2], args[3]), util.Markdown)
		} else {
			bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.settings.set.usage"), util.Markdown)
		}
	} else if util.CheckArgs(args[0], "remove", "rm", "delete", "del") {
		if !checkPerms(bot, message.Sender.ID, "whitelist.settings.remove") {
			return true
		}
		if len(args) > 2 {
			if !user.HasSetting(args[2]) {
				bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.settings.remove.fail", user.Name, args[2]), util.Markdown)
			} else {
				user.RemoveSetting(args[2])
				bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.settings.remove", user.Name, args[2]), util.Markdown)
			}
		} else {
			bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.settings.remove.usage"), util.Markdown)
		}
	} else {
		return false
	}
	return true
}

func handleWhitelistPerms(bot *telebot.Bot, message telebot.Message, args []string) bool {
	if len(args) < 2 {
		return false
	}
	user := config.GetUser(args[1])
	if user.UID == config.NilUser.UID {
		bot.SendMessage(message.Chat, lang.Translatef("mgmt.error.usernotfound", args[1]), util.Markdown)
		return true
	}

	if util.CheckArgs(args[0], "view", "see", "list") {
		if !checkPerms(bot, message.Sender.ID, "whitelist.permissions.view") {
			return true
		}
		if len(user.GetPermissions()) == 0 {
			bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.permissions.view.empty", user.Name), util.Markdown)
		} else {
			var buffer bytes.Buffer
			for _, perm := range user.GetPermissions() {
				buffer.WriteString(lang.Translatef("mgmt.whitelist.permissions.view.entry", perm))
			}
			bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.permissions.view", user.Name, buffer.String()), util.Markdown)
		}
	} else if util.CheckArgs(args[0], "add", "set") {
		if !checkPerms(bot, message.Sender.ID, "whitelist.permissions.add") {
			return true
		}
		if len(args) > 2 {
			if user.HasPermission(args[2]) {
				bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.permissions.add.fail", user.Name, args[2]), util.Markdown)
			} else {
				user.AddPermission(args[2])
				bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.permissions.add", user.Name, args[2]), util.Markdown)
			}
		} else {
			bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.permissions.add.usage"), util.Markdown)
		}
	} else if util.CheckArgs(args[0], "remove", "rm", "delete", "del") {
		if !checkPerms(bot, message.Sender.ID, "whitelist.permissions.remove") {
			return true
		}
		if len(args) > 2 {
			if !user.HasPermission(args[2]) {
				bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.permissions.remove.fail", user.Name, args[2]), util.Markdown)
			} else {
				user.RemovePermission(args[2])
				bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.permissions.remove", user.Name, args[2]), util.Markdown)
			}
		} else {
			bot.SendMessage(message.Chat, lang.Translatef("mgmt.whitelist.permissions.remove.usage"), util.Markdown)
		}
	} else {
		return false
	}
	return true
}

func checkPerms(bot *telebot.Bot, uid int, permission string) bool {
	user := config.GetUserWithUID(uid)
	if user.HasPermission(permission) {
		return true
	}
	bot.SendMessage(user, lang.Translatef("error.noperms", permission), util.Markdown)
	return false
}
