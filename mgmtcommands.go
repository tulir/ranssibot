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
	sender := config.GetUserWithUID(message.Sender.ID)
	if !checkPerms(bot, sender, "server.stop") {
		return
	}
	Shutdown(fmt.Sprintf("%[1]s %[2]s (ID %[3]d)", message.Sender.FirstName, message.Sender.LastName, message.Sender.ID))
}

func handleConfig(bot *telebot.Bot, message telebot.Message, args []string) {
	sender := config.GetUserWithUID(message.Sender.ID)
	if len(args) > 0 {
		if util.CheckArgs(args[0], "save") {
			if !checkPerms(bot, sender, "config.save") {
				return
			}
			if !config.IndentConfig && len(args) > 1 && util.CheckArgs(args[0], "pretty", "indent", "readable", "human", "debug") {
				config.IndentConfig = true
				config.Save()
				config.IndentConfig = false
				bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.config.save.debug"), util.Markdown)
			} else {
				config.Save()
				bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.config.save"), util.Markdown)
			}
		} else if util.CheckArgs(args[0], "load") {
			if !checkPerms(bot, sender, "config.load") {
				return
			}
			config.Load()
			bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.config.load"), util.Markdown)
		} else {
			bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.config.usage"), util.Markdown)
		}
	} else {
		bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.config.usage"), util.Markdown)
	}
}

func handleWhitelist(bot *telebot.Bot, message telebot.Message, args []string) {
	sender := config.GetUserWithUID(message.Sender.ID)
	if len(args) > 0 {
		if util.CheckArgs(args[0], "add") {
			if !checkPerms(bot, sender, "whitelist.add") {
				return
			}
			if len(args) > 3 {
				uid, err := strconv.Atoi(args[1])
				if err != nil {
					bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.add.parseint", args[1]), util.Markdown)
					return
				}
				year, err := strconv.Atoi(args[3])
				if err != nil {
					bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.add.parseint", args[3]), util.Markdown)
					return
				}
				if config.AddUser(config.CreateUser(uid, args[2], year)) {
					bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.add", uid, args[2], year), util.Markdown)
				} else {
					bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.add.alreadyused", uid, args[2]), util.Markdown)
				}
			} else {
				bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.add.usage"), util.Markdown)
			}
		} else if util.CheckArgs(args[0], "remove", "delete", "rm", "del") {
			if !checkPerms(bot, sender, "whitelist.remove") {
				return
			}
			if len(args) > 1 {
				user := config.GetUser(args[1])
				if user.UID == config.NilUser.UID {
					bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.error.usernotfound", args[1]), util.Markdown)
				} else {
					config.RemoveUser(args[1])
					bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.remove", user.UID, user.Name), util.Markdown)
				}
			} else {
				bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.remove.usage"), util.Markdown)
			}
		} else if util.CheckArgs(args[0], "get") {
			if !checkPerms(bot, sender, "whitelist.get") {
				return
			}
			if len(args) > 1 {
				user := config.GetUser(args[1])
				if user.UID == config.NilUser.UID {
					bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.error.usernotfound", args[1]), util.Markdown)
				} else {
					bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.get", user.UID, user.Name, user.Year), util.Markdown)
				}
			} else {
				bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.get.usage"), util.Markdown)
			}
		} else if util.CheckArgs(args[0], "list") {
			if !checkPerms(bot, sender, "whitelist.list") {
				return
			}

			var buffer bytes.Buffer
			for _, user := range config.GetAllUsers() {
				buffer.WriteString(lang.Translatef(sender, "mgmt.whitelist.list.entry", user.Name, user.UID))
			}
			bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.list", buffer.String()), util.Markdown)
		} else if util.CheckArgs(args[0], "permissions", "perms") {
			if !handleWhitelistPerms(bot, message, args[1:]) {
				bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.permissions.usage"), util.Markdown)
			}
		} else if util.CheckArgs(args[0], "settings", "properties", "props", "preferences", "prefs") {
			if !handleWhitelistSettings(bot, message, args[1:]) {
				bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.settings.usage"), util.Markdown)
			}
		} else {
			bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.usage"), util.Markdown)
		}
	} else {
		bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.usage"), util.Markdown)
	}
}

func handleWhitelistSettings(bot *telebot.Bot, message telebot.Message, args []string) bool {
	sender := config.GetUserWithUID(message.Sender.ID)
	if len(args) < 2 {
		return false
	}
	user := config.GetUser(args[1])
	if user.UID == config.NilUser.UID {
		bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.error.usernotfound", args[1]), util.Markdown)
		return true
	}

	if util.CheckArgs(args[0], "view", "see", "list") {
		if !checkPerms(bot, sender, "whitelist.settings.view") {
			return true
		}
		if len(user.Settings) == 0 {
			bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.settings.view.empty", user.Name), util.Markdown)
		} else {
			var buffer bytes.Buffer
			for key, val := range user.Settings {
				buffer.WriteString(lang.Translatef(sender, "mgmt.whitelist.settings.view.entry", key, val))
			}
			bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.settings.view", user.Name, buffer.String()), util.Markdown)
		}
	} else if util.CheckArgs(args[0], "get") {
		if !checkPerms(bot, sender, "whitelist.settings.get") {
			return true
		}
		if len(args) > 2 {
			val, ok := user.GetSetting(args[2])
			if !ok {
				bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.settings.get.notfound", user.Name, args[2]), util.Markdown)
			} else {
				bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.settings.get", user.Name, args[2], val), util.Markdown)
			}
		} else {
			bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.settings.get.usage"), util.Markdown)
		}
	} else if util.CheckArgs(args[0], "add", "set") {
		if !checkPerms(bot, sender, "whitelist.settings.add") {
			return true
		}
		if len(args) > 3 {
			user.SetSetting(args[2], args[3])
			bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.settings.set", user.Name, args[2], args[3]), util.Markdown)
		} else {
			bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.settings.set.usage"), util.Markdown)
		}
	} else if util.CheckArgs(args[0], "remove", "rm", "delete", "del") {
		if !checkPerms(bot, sender, "whitelist.settings.remove") {
			return true
		}
		if len(args) > 2 {
			if !user.HasSetting(args[2]) {
				bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.settings.remove.fail", user.Name, args[2]), util.Markdown)
			} else {
				user.RemoveSetting(args[2])
				bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.settings.remove", user.Name, args[2]), util.Markdown)
			}
		} else {
			bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.settings.remove.usage"), util.Markdown)
		}
	} else {
		return false
	}
	return true
}

func handleWhitelistPerms(bot *telebot.Bot, message telebot.Message, args []string) bool {
	sender := config.GetUserWithUID(message.Sender.ID)
	if len(args) < 2 {
		return false
	}
	user := config.GetUser(args[1])
	if user.UID == config.NilUser.UID {
		bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.error.usernotfound", args[1]), util.Markdown)
		return true
	}

	if util.CheckArgs(args[0], "view", "see", "list") {
		if !checkPerms(bot, sender, "whitelist.permissions.view") {
			return true
		}
		if len(user.Permissions) == 0 {
			bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.permissions.view.empty", user.Name), util.Markdown)
		} else {
			var buffer bytes.Buffer
			for _, perm := range user.Permissions {
				buffer.WriteString(lang.Translatef(sender, "mgmt.whitelist.permissions.view.entry", perm))
			}
			bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.permissions.view", user.Name, buffer.String()), util.Markdown)
		}
	} else if util.CheckArgs(args[0], "add", "set") {
		if !checkPerms(bot, sender, "whitelist.permissions.add") {
			return true
		}
		if len(args) > 2 {
			if user.HasPermission(args[2]) {
				bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.permissions.add.fail", user.Name, args[2]), util.Markdown)
			} else {
				user.AddPermission(args[2])
				bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.permissions.add", user.Name, args[2]), util.Markdown)
			}
		} else {
			bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.permissions.add.usage"), util.Markdown)
		}
	} else if util.CheckArgs(args[0], "remove", "rm", "delete", "del") {
		if !checkPerms(bot, sender, "whitelist.permissions.remove") {
			return true
		}
		if len(args) > 2 {
			if !user.HasPermission(args[2]) {
				bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.permissions.remove.fail", user.Name, args[2]), util.Markdown)
			} else {
				user.RemovePermission(args[2])
				bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.permissions.remove", user.Name, args[2]), util.Markdown)
			}
		} else {
			bot.SendMessage(message.Chat, lang.Translatef(sender, "mgmt.whitelist.permissions.remove.usage"), util.Markdown)
		}
	} else {
		return false
	}
	return true
}

func checkPerms(bot *telebot.Bot, user config.User, permission string) bool {
	if user.HasPermission(permission) {
		return true
	}
	bot.SendMessage(user, lang.Translatef(user, "error.noperms", permission), util.Markdown)
	return false
}
