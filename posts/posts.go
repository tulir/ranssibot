package posts

import (
	"github.com/tucnak/telebot"
	"maunium.net/ranssibot/lang"
	"maunium.net/ranssibot/util"
	"strings"
)

// HandleCommand handles Ranssi post commands
func HandleCommand(bot *telebot.Bot, message telebot.Message, args []string) {
	if len(args) == 1 {
		handleNews(bot, message, args)
	} else if strings.EqualFold(args[1], "subscribe") || strings.EqualFold(args[1], "sub") {
		handleSubscribe(bot, message, args)
	} else if strings.EqualFold(args[1], "unsubscribe") || strings.EqualFold(args[1], "unsub") {
		handleUnsubscribe(bot, message, args)
	} else if strings.EqualFold(args[1], "get") || strings.EqualFold(args[1], "read") {
		handleRead(bot, message, args)
	} else if strings.EqualFold(args[1], "latest") || strings.EqualFold(args[1], "news") {
		handleNews(bot, message, args)
	} else if strings.EqualFold(args[1], "comment") || strings.EqualFold(args[1], "message") || strings.EqualFold(args[1], "spam") {
		handleComment(bot, message, args)
	} else {
		bot.SendMessage(message.Chat, lang.Translate("posts.usage"), util.Markdown)
	}
}
