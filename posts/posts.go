package posts

import (
	"github.com/tucnak/telebot"
	"maunium.net/ranssibot/lang"
	"maunium.net/ranssibot/util"
	"strings"
)

// HandleCommand handles Ranssi post commands
func HandleCommand(bot *telebot.Bot, message telebot.Message, args []string) {
	if len(args) == 0 {
		handleNews(bot, message, args[1:])
	} else if strings.EqualFold(args[0], "subscribe") || strings.EqualFold(args[0], "sub") {
		handleSubscribe(bot, message, args[1:])
	} else if strings.EqualFold(args[0], "unsubscribe") || strings.EqualFold(args[0], "unsub") {
		handleUnsubscribe(bot, message, args[1:])
	} else if strings.EqualFold(args[0], "get") || strings.EqualFold(args[0], "read") {
		handleRead(bot, message, args[1:])
	} else if strings.EqualFold(args[0], "latest") || strings.EqualFold(args[0], "news") {
		handleNews(bot, message, args[1:])
	} else if strings.EqualFold(args[0], "comment") || strings.EqualFold(args[0], "message") || strings.EqualFold(args[0], "spam") {
		handleComment(bot, message, args[1:])
	} else {
		bot.SendMessage(message.Chat, lang.Translate("posts.usage"), util.Markdown)
	}
}
