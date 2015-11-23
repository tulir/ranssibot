package laundry

import (
	"github.com/tucnak/telebot"
	"golang.org/x/net/html"
	"maunium.net/ranssibot/lang"
	"maunium.net/ranssibot/log"
	"maunium.net/ranssibot/util"
	"strings"
	"time"
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

	// Get the timetable page and convert the string to a reader
	reader := strings.NewReader(util.HTTPGet("http://ranssi.paivola.fi/pyykit.php"))
	// Parse the HTML from the reader
	doc, err := html.Parse(reader)
	// Check if there was an error
	if err != nil {
		// Print the error
		log.Errorf("Failed to read laundry lists: %s", err)
		// Return
		return
	}
	// Find the timetable table header node
	laundrynode := util.FindSpan("tr", "class", "today",
		util.FindSpan("table", "class", "pyykit", doc)).
		FirstChild.NextSibling.NextSibling.NextSibling

	minutes := minutesInDay()
	switch minutes {
	case 21 * 60:
		print("1")
		laundrynode = laundrynode.NextSibling.NextSibling
		fallthrough
	case 17*60 + 00:
		print("2")
		laundrynode = laundrynode.NextSibling.NextSibling
		fallthrough
	case 14*60 + 45:
		print("3")
		laundrynode = laundrynode.NextSibling.NextSibling
		fallthrough
	case 12*60 + 00:
		print("4")
		laundrynode = laundrynode.NextSibling.NextSibling
		fallthrough
	case 8*60 + 45:
		print("5")
		laundrynode = laundrynode.NextSibling.NextSibling
	}
	println(util.Render(laundrynode))
}

func minutesInDay() int {
	t := time.Now()
	return t.Hour()*60 + t.Minute()
}

// HandleCommand handles laundry commands
func HandleCommand(bot *telebot.Bot, message telebot.Message, args []string) {
	if len(args) > 1 {
		laundry[message.Chat.ID] = args[1:]
	} else {
		bot.SendMessage(message.Chat, lang.Translate("laundry.usage"), util.Markdown)
	}
}
