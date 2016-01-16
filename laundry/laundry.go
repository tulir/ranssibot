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

var laundry = make(map[string]int)

var notified = 0
var day = 0

var notifyMinutes = []int{
	8*60 + 00,
	11*60 + 15,
	14*60 + 30,
	17*60 + 00,
	18*60 + 45,
}

var notifMinDiff = 5
var notifMaxDiff = 5

// Loop TODO: make comment
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
		if minsNow >= notifyMinutes[notified]+notifMaxDiff {
			// If the notified status is incorrect (as in the notify point was already over notifMaxDiff minutes ago)
			// increment the notified status and continue loop.
			notified++
			log.Debugf("Skipping notifications for turn #%d", notified)
			continue
		} else if minsNow >= notifyMinutes[notified]-notifMinDiff {
			notified++
			log.Debugf("Sending notifications for turn #%d", notified)
			Notify(bot, notified)
		}
		time.Sleep(1 * time.Minute)
	}
}

func minutesInDay(time time.Time) int {
	return time.Hour()*60 + time.Minute()
}

// Notify TODO: make comment
func Notify(bot *telebot.Bot, time int) {
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

	// TODO uncomment to enable attribute checking.
	/*for _, attr := range laundrynode.Attr {
		// If the current node is not marked as busy, don't send any notifications.
		if attr.Key == "class" && attr.Val != "busy" {
			println(attr.Val)
			return
		}
	}*/
	curName := laundrynode.LastChild.Data
	for _, user := range config.GetAllUsers() {
		str, ok := user.GetSetting("laundry")
		if ok {
			names := strings.Split(str, ", ")
			for _, name := range names {
				if curName == name {
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
			var buf bytes.Buffer
			for n, arg := range args[1:] {
				buf.Write([]byte(arg))
				if n+2 < len(args) {
					buf.Write([]byte(", "))
				}
			}
			sender.SetSetting("laundry", buf.String())
			bot.SendMessage(message.Chat, lang.Translatef(sender, "laundry.subscribed", args[1]), util.Markdown)
		}
	} else {
		bot.SendMessage(message.Chat, lang.Translatef(sender, "laundry.usage"), util.Markdown)
	}
}
