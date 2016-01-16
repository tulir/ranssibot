package laundry

import (
	"github.com/tucnak/telebot"
	log "maunium.net/go/maulogger"
	"maunium.net/go/ranssibot/config"
	"maunium.net/go/ranssibot/lang"
	"maunium.net/go/ranssibot/util"
	"strings"
	"time"
)

var laundry = make(map[int][]string)

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
		minsNow := minutesInDay(now)
		if minsNow >= notifyMinutes[notified]+notifMaxDiff {
			// If the notified status is incorrect (as in the notify point was already over notifMaxDiff minutes ago)
			// increment the notified status and continue loop.
			notified++
		} else if minsNow >= notifyMinutes[notified]-notifMinDiff {
			notified++
			Notify(notified)
			time.Sleep(1 * time.Minute)
		}
	}
}

func minutesInDay(time time.Time) int {
	return time.Hour()*60 + time.Minute()
}

// Notify TODO: make comment
func Notify(time int) {
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
	laundrynode := util.FindSpan("tr", "class", "today", util.FindSpan("table", "class", "pyykit", doc)).FirstChild

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
		if attr.Key == "class" && attr.Val != "busy" {
			println(attr.Val)
			return
		}
	}
	println(util.Render(laundrynode.LastChild))
}

// HandleCommand handles laundry commands
func HandleCommand(bot *telebot.Bot, message telebot.Message, args []string) {
	sender := config.GetUserWithUID(message.Sender.ID)
	for i := 0; i < len(args); i++ {
		args[i] = strings.ToLower(args[i])
	}
	if len(args) > 2 {
		if util.CheckArgs(args[0], "listen", "watch", "sub", "subscribe") {
			laundry[sender.UID] = args[1:]
		}
	} else {
		bot.SendMessage(message.Chat, lang.Translatef(sender, "laundry.usage"), util.Markdown)
	}
}
