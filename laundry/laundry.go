package laundry

import (
	"bytes"
	"github.com/tucnak/telebot"
	log "maunium.net/go/maulogger"
	"maunium.net/go/ranssibot/config"
	"maunium.net/go/ranssibot/lang"
	"maunium.net/go/ranssibot/util"
	"strings"
	"time"
)

var notified = 0
var day = 0

var notifyMinutes = []int{
	8*60 + 00,
	11*60 + 15,
	14*60 + 30,
	17*60 + 00,
	18*60 + 45,
}

// Loop is an infinite loop that checks if the time is right to spam about laundry turns.
func Loop(bot *telebot.Bot) {
	for {
		now := time.Now()
		if day != now.Day() {
			notified = 0
			day = now.Day()
		}
		if notified == 5 {
			// All of the days notifications have been sent, so sleep times can be longer.
			time.Sleep(20 * time.Minute)
			continue
		}
		minsNow := minutesInDay(now)
		if minsNow >= notifyMinutes[notified]+5 {
			// If the notified status is incorrect (as in the notify point was already over notifMaxDiff minutes ago)
			// increment the notified status and continue loop.
			notified++
			log.Debugf("Skipping notifications for laundry turn #%d", notified)
			continue
		} else if minsNow >= notifyMinutes[notified]-5 {
			notified++
			notify(bot, notified)
		}
		time.Sleep(2 * time.Minute)
	}
}

func minutesInDay(time time.Time) int {
	return time.Hour()*60 + time.Minute()
}

func notify(bot *telebot.Bot, time int) {
	// Get the timetable page
	doc, err := util.HTTPGetMinAndParse("http://ranssi.paivola.fi/pyykit.php")
	// Check if there was an error
	if err != nil {
		// Print the error
		log.Errorf("Failed to read laundry lists: %s", err)
		// Return
		return
	}
	// Find the timetable table header node
	laundrynode := util.FindSpan("tr", "class", "today", doc).FirstChild

	switch time {
	case 5:
		laundrynode = laundrynode.NextSibling
		fallthrough
	case 4:
		laundrynode = laundrynode.NextSibling
		fallthrough
	case 3:
		laundrynode = laundrynode.NextSibling
		fallthrough
	case 2:
		laundrynode = laundrynode.NextSibling
		fallthrough
	case 1:
		laundrynode = laundrynode.NextSibling
	}

	for _, attr := range laundrynode.Attr {
		// If the current node is not marked as busy, don't send any notifications.
		if attr.Key == "class" && attr.Val == "free" {
			log.Debugf("Laundry turn #%d is marked as free, skipping.", notified)
			return
		}
	}

	curName := laundrynode.LastChild.Data
	curNameLower := strings.ToLower(curName)
	if len(curName) == 0 {
		log.Debugf("Laundry turn #%d is empty, skipping.", notified)
		return
	}
	log.Debugf("Sending notifications for laundry turn #%d", notified)
	for _, user := range config.GetAllUsers() {
		str, ok := user.GetSetting("laundry")
		if ok {
			names := strings.Split(str, ", ")
			for _, name := range names {
				if curNameLower == name {
					bot.SendMessage(user, lang.Translatef(user, "laundry.soon", curName), util.Markdown)
				}
			}
		}
	}
}

// HandleCommand handles laundry commands
func HandleCommand(bot *telebot.Bot, message telebot.Message, args []string) {
	sender := config.GetUserWithUID(message.Sender.ID)
	for i := 0; i < len(args); i++ {
		args[i] = strings.ToLower(args[i])
	}
	if len(args) > 1 {
		if util.CheckArgs(args[0], "listen", "watch", "sub", "subscribe", "spam") {
			if util.CheckArgs(args[1], "clear", "empty") {
				bot.SendMessage(message.Chat, lang.Translatef(sender, "laundry.subscriptionscleared"), util.Markdown)
				return
			}
			var buf bytes.Buffer
			for n, arg := range args[1:] {
				buf.Write([]byte(arg))
				bot.SendMessage(message.Chat, lang.Translatef(sender, "laundry.subscribed", arg), util.Markdown)
				if n+2 < len(args) {
					buf.Write([]byte(", "))
				}
			}
			listenFor := strings.ToLower(buf.String())
			sender.SetSetting("laundry", listenFor)
		}
	} else {
		bot.SendMessage(message.Chat, lang.Translatef(sender, "laundry.usage"), util.Markdown)
	}
}
