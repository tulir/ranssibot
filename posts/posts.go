package posts

import (
	"github.com/tucnak/telebot"
	"maunium.net/go/ranssibot/lang"
	"maunium.net/go/ranssibot/util"
)

// HandleCommand handles Ranssi post commands
func HandleCommand(bot *telebot.Bot, message telebot.Message, args []string) {
	if len(args) == 0 {
		handleNews(bot, message, args)
	} else if util.CheckArgs(args[0], "subscribe", "sub") {
		handleSubscribe(bot, message, args[1:])
	} else if util.CheckArgs(args[0], "unsubscribe", "unsub") {
		handleUnsubscribe(bot, message, args[1:])
	} else if util.CheckArgs(args[0], "read", "get", "view") {
		handleRead(bot, message, args[1:])
	} else if util.CheckArgs(args[0], "news", "latest") {
		handleNews(bot, message, args[1:])
	} else if util.CheckArgs(args[0], "comment", "message", "msg", "spam") {
		handleComment(bot, message, args[1:])
	} else if util.CheckArgs(args[0], "comments", "readcomments", "viewcomments", "getcomments") {
		handleReadComments(bot, message, args[1:])
	} else {
		bot.SendMessage(message.Chat, lang.Translate("posts.usage"), util.Markdown)
	}
}
