package laundry

import (
	"github.com/tucnak/telebot"
	"maunium.net/ranssibot/lang"
	"maunium.net/ranssibot/util"
)

var laundry = make(map[int][]string)

// NotifierTick notifies the people who have a laundry turn coming up soon
func NotifierTick() {
	// TODO: Load next laundry turns, check if any of the reservation names
	// are found in the laundry name registry and if found, send a notification.
	//
	// ALSO: Don't forget to make something call this method in a separate
	// thread on a specific interval.
	//
	// ALSO²: If this method is directly put into a new thread,
	// remember to add a sleep call.
}

// HandleCommand handles laundry commands
func HandleCommand(bot *telebot.Bot, message telebot.Message, args []string) {
	if len(args) > 1 {
		laundry[message.Chat.ID] = args[1:]
	} else {
		bot.SendMessage(message.Chat, lang.Translate("laundry.usage"), util.Markdown)
	}
}
