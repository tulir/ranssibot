package food

import (
	"encoding/json"
	"github.com/tucnak/telebot"
	log "maunium.net/go/maulogger"
	"maunium.net/go/ranssibot/config"
	"maunium.net/go/ranssibot/lang"
	"maunium.net/go/ranssibot/util"
)

const (
	subSetting = "food-subscription"
)

// Menu contains the menu entries for one day.
type Menu struct {
	Breakfast string `json:"aamu"`
	Lunch     string `json:"lounas"`
	Coffee    string `json:"kahvi"`
	Dinner    string `json:"paiva"`
}

var menu Menu
var lastUpdate int64

func init() {
	menustr := util.HTTPGetMin("https://ruoka.paivola.fi/api.php")
	json.Unmarshal([]byte(menustr), &menu)
	lastUpdate = util.Timestamp()
}

// HandleCommand handles Food/Menu commands
func HandleCommand(bot *telebot.Bot, message telebot.Message, args []string) {
	if len(args) == 0 {
		handleGetMenu(bot, message, args)
	} else if util.CheckArgs(args[0], "subscribe", "sub", "spam", "watch", "listen") {
		handleSubscribe(bot, message, args)
	} else if util.CheckArgs(args[0], "unsubscribe", "unsub", "nospam") {
		handleUnsubscribe(bot, message, args)
	}
}

func handleGetMenu(bot *telebot.Bot, message telebot.Message, args []string) {
	if lastUpdate+600 < util.Timestamp() {
		UpdateMenu()
	}
	sender := config.GetUserWithUID(message.Sender.ID)
	bot.SendMessage(message.Chat, lang.Translatef(sender, "food.menu", menu.Breakfast, menu.Lunch, menu.Coffee, menu.Dinner), util.Markdown)
}

func handleSubscribe(bot *telebot.Bot, message telebot.Message, args []string) {
	sender := config.GetUserWithUID(message.Sender.ID)
	if isSubscribed(sender) {
		bot.SendMessage(message.Chat, lang.Translate(sender, "food.alreadysubscribed"), util.Markdown)
		log.Debugf("%[1]s attempted to subscribe to the food notification list, but was already subscribed", sender.Name)
	} else {
		sender.SetSetting(subSetting, "true")
		config.ASave()
		bot.SendMessage(message.Chat, lang.Translate(sender, "food.subscribed"), util.Markdown)
		log.Debugf("%[1]s successfully subscribed to the food notifcation list", sender.Name)
	}
}

func handleUnsubscribe(bot *telebot.Bot, message telebot.Message, args []string) {
	sender := config.GetUserWithUID(message.Sender.ID)
	if !isSubscribed(sender) {
		bot.SendMessage(message.Chat, lang.Translate(sender, "food.notsubscribed"), util.Markdown)
		log.Debugf("%[1]s attempted to unsubscribe from the food notification list, but was not subscribed", sender.Name)
	} else {
		sender.RemoveSetting(subSetting)
		config.ASave()
		bot.SendMessage(message.Chat, lang.Translate(sender, "food.unsubscribed"), util.Markdown)
		log.Debugf("%[1]s successfully unsubscribed from the food notifcation list", sender.Name)
	}
}

// UpdateMenu updates the menu.
func UpdateMenu() {
	menustr := util.HTTPGetMin("https://ruoka.paivola.fi/api.php")
	json.Unmarshal([]byte(menustr), &menu)
	lastUpdate = util.Timestamp()
	log.Debugf("Successfully updated today's menu.")
}

func isSubscribed(u config.User) bool {
	return u.HasSetting(subSetting)
}
