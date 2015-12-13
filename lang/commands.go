package lang

import (
	"github.com/tucnak/telebot"
	"maunium.net/go/ranssibot/config"
	"maunium.net/go/ranssibot/util"
)

const (
	langChangePerm      = "-lang.change"
	langChangeOtherPerm = "lang.change.other"
)

// HandleCommand handles a /language command
func HandleCommand(bot *telebot.Bot, message telebot.Message, args []string) {
	sender := config.GetUserWithUID(message.Sender.ID)
	if len(args) > 0 {
		lang := GetLanguage(args[0])
		if lang == nil {
			bot.SendMessage(message.Chat, UTranslatef(sender, "lang.notfound", lang.Name), util.Markdown)
			return
		}

		if len(args) > 1 {
			user := config.GetUser(args[1])
			if user.UID != config.NilUser.UID && !sender.HasPermission(langChangeOtherPerm) {
				bot.SendMessage(message.Chat, UTranslatef(sender, "error.noperms", langChangeOtherPerm), util.Markdown)
				return
			}

			user.SetSetting("language", lang.Name)
			bot.SendMessage(message.Chat, UTranslatef(sender, "lang.changed.other", lang.Name, user.Name), util.Markdown)
		} else {
			if sender.HasPermission(langChangePerm) {
				bot.SendMessage(message.Chat, UTranslatef(sender, "error.noperms", langChangePerm), util.Markdown)
				return
			}

			sender.SetSetting("language", lang.Name)
			bot.SendMessage(message.Chat, UTranslatef(sender, "lang.changed", lang.Name), util.Markdown)
		}
	} else {
		bot.SendMessage(message.Chat, UTranslatef(sender, "lang.usage"), util.Markdown)
	}
}
