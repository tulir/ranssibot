package timetables

import (
	"errors"
	"fmt"
	"github.com/tucnak/telebot"
	"golang.org/x/net/html"
	log "maunium.net/maulogger"
	"maunium.net/ranssibot/lang"
	"maunium.net/ranssibot/util"
	"maunium.net/ranssibot/whitelist"
	"strconv"
	"strings"
)

// TimetableLesson is a struct that contains the data for a lesson in a timetable
type TimetableLesson struct {
	Subject  string
	TimeName string
	Date     util.Date
	Time     util.Time
}

var firstyear = [26][4]TimetableLesson{}
var secondyear = [26][4]TimetableLesson{}
var other = [26]TimetableLesson{}

// The day ID for today
var today = 5

// The last time the timetable cache was updated
var lastupdate = util.Timestamp()

// HandleCommand handles a /timetable command
func HandleCommand(bot *telebot.Bot, message telebot.Message, args []string) {
	if util.Timestamp() > lastupdate+600 {
		bot.SendMessage(message.Chat, "Updating cached timetables...", util.Markdown)
		Update()
	}

	// SetTime code
	/*if len(args) > 3 {
		lessonID, err := strconv.Atoi(args[2])
		if err != nil {
			bot.SendMessage(message.Chat, fmt.Sprintf(lang.Translate("error.parse.integer"), args[2]), util.Markdown)
		}
		dayShift, err := strconv.Atoi(args[3])
		if err != nil {
			bot.SendMessage(message.Chat, fmt.Sprintf(lang.Translate("error.parse.integer"), args[3]), util.Markdown)
		}
		time, err := timetables.StringToTime(args[4])
		if err != nil {
			bot.SendMessage(message.Chat, fmt.Sprintf(lang.Translate("error.parse.time"), args[4]), util.Markdown)
		}
		if args[1] == "ventit" {
			firstyear[today+dayShift][lessonID].Time = time
		} else if args[1] == "neliöt" {
			secondyear[today+dayShift][lessonID].Time = time
		} else if args[1] == "other" {
			other[today+dayShift].Time = time
		} else {
			return
		}
		bot.SendMessage(message.Chat, fmt.Sprintf(lang.Translate("settime.success"), args[1], lessonID, dayShift, TimeToString(time)), util.Markdown)
	} else {
		bot.SendMessage(message.Chat, lang.Translate("settime.usage"), util.Markdown)
	}*/

	day := today
	year := whitelist.GetYeargroupIndex(message.Sender.ID)
	if len(args) == 2 {
		if strings.EqualFold(args[1], lang.Translate("timetable.year.first")) {
			year = 1
		} else if strings.EqualFold(args[1], lang.Translate("timetable.year.second")) {
			year = 2
		} else if strings.EqualFold(args[1], "update") {
			Update()
			bot.SendMessage(message.Chat, lang.Translate("timetable.update.success"), util.Markdown)
		} else {
			dayNew, err := shift(day, args[1], 0, len(other), bot, message)
			if err != nil {
				return
			}
			day = dayNew
		}
	} else if len(args) == 3 {
		if strings.EqualFold(args[1], lang.Translate("timetable.year.first")) {
			year = 1
		} else if strings.EqualFold(args[1], lang.Translate("timetable.year.second")) {
			year = 2
		} else {
			bot.SendMessage(message.Chat, lang.Translate("timetable.usage"), util.Markdown)
		}
		dayNew, err := shift(day, args[2], 0, len(other), bot, message)
		if err != nil {
			return
		}
		day = dayNew
	}

	if day < 0 || day >= len(other) {
		bot.SendMessage(message.Chat, lang.Translate("timetable.nodata"), util.Markdown)
		return
	}

	if year == 1 {
		sendFirstYear(day, bot, message)
	} else if year == 2 {
		sendSecondYear(day, bot, message)
	} else {
		bot.SendMessage(message.Chat, lang.Translate("timetable.noyeargroup"), util.Markdown)
	}
}

func shift(toShift int, shiftBy string, min, max int, bot *telebot.Bot, message telebot.Message) (int, error) {
	shift, err := strconv.Atoi(shiftBy)
	if err == nil {
		toShift += shift
		if toShift < min || toShift > max {
			bot.SendMessage(message.Chat, lang.Translate("timetable.nodata"), util.Markdown)
			return -9999, errors.New("OOB")
		}
		return toShift, nil
	}
	bot.SendMessage(message.Chat, lang.Translate("timetable.usage"), util.Markdown)
	return -9999, errors.New("PARSEINT")
}

// Update the timetables from http://ranssi.paivola.fi/lj.php
func Update() {
	// Get the timetable page and convert the string to a reader
	reader := strings.NewReader(util.HTTPGet("http://ranssi.paivola.fi/lj.php"))
	// Parse the HTML from the reader
	doc, err := html.Parse(reader)
	// Check if there was an error
	if err != nil {
		// Print the error
		log.Errorf("Failed to update timetable cache: %s", err)
		// Return
		return
	}

	// Find the timetable table header node
	ttnode := util.FindSpan("tr", "class", "header", doc)
	// Check if the node was found
	if ttnode != nil {
		dayentry := ttnode
		// Loop through the days in the timetable
		for day := 0; ; day++ {
			// Make sure the next day exists
			if dayentry.NextSibling == nil ||
				dayentry.NextSibling.NextSibling == nil ||
				dayentry.NextSibling.NextSibling.FirstChild == nil {
				break
			}
			// Get the first day node
			dayentry = dayentry.NextSibling.NextSibling

			var date util.Date
			// Get the date of this day
			dateraw := strings.Split(dayentry.FirstChild.NextSibling.LastChild.Data, ".")
			dateraw[0] = strings.Split(dateraw[0], "\n")[1]

			// Parse the day from the date
			dateday, err1 := strconv.Atoi(dateraw[0])
			// Parse the month from the date
			datemonth, err2 := strconv.Atoi(dateraw[1])
			// Parse the year from the date
			dateyear, err3 := strconv.Atoi(dateraw[2])
			// If no errors came in parsing, create a Date struct from the parsed data
			// If there were errors, set the date to 1.1.1970
			if err1 == nil && err2 == nil && err3 == nil {
				date = util.Date{Year: dateyear, Month: datemonth, Day: dateday}
			} else {
				date = util.Date{Year: 1970, Month: 1, Day: 1}
			}
			// Get the first lesson node in the day node
			entry := dayentry.FirstChild.NextSibling
			// Loop through the lessons on the day
			for lesson := 0; lesson < 9; lesson++ {
				// Make sure the next lesson exists
				if entry == nil ||
					entry.NextSibling == nil ||
					entry.NextSibling.NextSibling == nil {
					break
				}
				// Get the next lesson node
				entry = entry.NextSibling.NextSibling
				data := "tyhjää"
				// Check if the lesson contains anything
				if entry.FirstChild != nil {
					// Lesson data found; Try to parse it
					if entry.FirstChild.Type == html.TextNode {
						// Found lesson data directly under lesson node
						data = entry.FirstChild.Data
					} else if entry.FirstChild.Type == html.ElementNode {
						// Didn't find data directly under letimetable[day][lesson]sson node
						// Check for a child element node.
						if entry.FirstChild.FirstChild != nil {
							// Child element node found. Check if the child of that child is text.
							if entry.FirstChild.FirstChild.Type == html.TextNode {
								// Child of child is text, use it as the data.
								data = entry.FirstChild.FirstChild.Data
							}
						}
					}
				}
				lsn := ParseLesson(data)
				if lsn != nil {
					data = lsn.ToString()
				}
				// Save the parsed data to the correct location.
				switch lesson {
				case 0:
					firstyear[day][0] = TimetableLesson{data, "Aamu", date, util.Time{Hours: 9, Minutes: 0}}
				case 1:
					firstyear[day][1] = TimetableLesson{data, "IP1", date, util.Time{Hours: 12, Minutes: 15}}
				case 2:
					firstyear[day][2] = TimetableLesson{data, "IP2", date, util.Time{Hours: 15, Minutes: 0}}
				case 3:
					firstyear[day][3] = TimetableLesson{data, "Ilta", date, util.Time{Hours: 19, Minutes: 0}}
				case 4:
					other[day] = TimetableLesson{data, "Muuta", date, util.Time{Hours: 0, Minutes: 0}}
				case 5:
					secondyear[day][0] = TimetableLesson{data, "Aamu", date, util.Time{Hours: 9, Minutes: 0}}
				case 6:
					secondyear[day][1] = TimetableLesson{data, "IP1", date, util.Time{Hours: 12, Minutes: 15}}
				case 7:
					secondyear[day][2] = TimetableLesson{data, "IP2", date, util.Time{Hours: 15, Minutes: 0}}
				case 8:
					secondyear[day][3] = TimetableLesson{data, "Ilta", date, util.Time{Hours: 19, Minutes: 0}}
				}
			}
		}
		lastupdate = util.Timestamp()
	} else {
		// Node not found, print error
		log.Errorf("Error updating timetables: Failed to find timetable table header node!")
		lastupdate = 0
	}
}

func sendFirstYear(day int, bot *telebot.Bot, message telebot.Message) {
	bot.SendMessage(message.Chat,
		fmt.Sprintf(lang.Translate("timetable.generic"),
			firstyear[day][0].Subject, firstyear[day][1].Subject, firstyear[day][2].Subject, firstyear[day][3].Subject,
			util.DateToString(firstyear[day][0].Date))+"\n"+fmt.Sprintf(lang.Translate("timetable.other"), other[day].Subject),
		util.Markdown)
}

func sendSecondYear(day int, bot *telebot.Bot, message telebot.Message) {
	bot.SendMessage(message.Chat,
		fmt.Sprintf(lang.Translate("timetable.generic"),
			secondyear[day][0].Subject, secondyear[day][1].Subject, secondyear[day][2].Subject, secondyear[day][3].Subject,
			util.DateToString(secondyear[day][0].Date))+"\n"+fmt.Sprintf(lang.Translate("timetable.other"), other[day].Subject),
		util.Markdown)
}
