package lang

import (
	"github.com/tucnak/telebot"
	"maunium.net/go/ranssibot/config"
	"maunium.net/go/ranssibot/util"
)

const (
	langChangePerm      = "lang.change"
	langChangeOtherPerm = "lang.change.other"
)

// HandleCommand handles a /language command
func HandleCommand(bot *telebot.Bot, message telebot.Message, args []string) {
	if len(args) > 0 {
		lang := GetLanguage(args[0])
		if lang == nil {
			bot.SendMessage(message.Chat, Translatef("lang.notfound", lang.Name), util.Markdown)
			return
		}

		if len(args) > 1 {
			user := config.GetUser(args[1])
			if user.UID != config.NilUser.UID && !config.GetUserWithUID(message.Sender.ID).HasPermission(langChangeOtherPerm) {
				bot.SendMessage(message.Chat, Translatef("error.noperms", langChangeOtherPerm), util.Markdown)
				return
			}

			user.SetSetting("language", lang.Name)
			bot.SendMessage(message.Chat, Translatef("lang.changed.other", lang.Name, user.Name), util.Markdown)
		} else {
			user := config.GetUserWithUID(message.Sender.ID)
			if !user.HasPermission(langChangePerm) {
				bot.SendMessage(message.Chat, Translatef("error.noperms", langChangePerm), util.Markdown)
				return
			}

			user.SetSetting("language", lang.Name)
			bot.SendMessage(message.Chat, Translatef("lang.changed", lang.Name), util.Markdown)
		}
	} else {
		bot.SendMessage(message.Chat, Translatef("lang.usage"), util.Markdown)
	}
}
