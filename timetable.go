package main

import (
	"fmt"
	"github.com/tucnak/telebot"
	"golang.org/x/net/html"
	"log"
	"strconv"
	"strings"
)

// Date ...
type Date struct {
	Year  int
	Month int
	Day   int
}

// DateToString ...
func DateToString(date Date) string {
	dayStr := strconv.Itoa(date.Day)
	if len(dayStr) == 1 {
		dayStr = "0" + dayStr
	}
	monthStr := strconv.Itoa(date.Month)
	if len(monthStr) == 1 {
		monthStr = "0" + monthStr
	}
	yearStr := strconv.Itoa(date.Year)
	if len(yearStr) == 1 {
		yearStr = "0" + yearStr
	}

	return dayStr + "." + monthStr + "." + yearStr
}

// Time ...
type Time struct {
	Hours   int
	Minutes int
}

// TimeToString ...
func TimeToString(time Time) string {
	hourStr := strconv.Itoa(time.Hours)
	if len(hourStr) == 1 {
		hourStr = "0" + hourStr
	}
	minStr := strconv.Itoa(time.Minutes)
	if len(minStr) == 1 {
		minStr = "0" + minStr
	}
	return hourStr + ":" + minStr
}

// TimetableLesson is a struct that contains the data for a lesson in a timetable
type TimetableLesson struct {
	Subject  string
	TimeName string
	Date     Date
	Time     Time
}

var firstyear = [26][4]TimetableLesson{}
var secondyear = [26][4]TimetableLesson{}
var other = [26]TimetableLesson{}

// The last time the timetable cache was updated
var lastupdate = timestamp()

// Update the timetables from http://ranssi.paivola.fi/lj.php
func updateTimes() {
	// Get the timetable page and convert the string to a reader
	reader := strings.NewReader(httpGet("http://ranssi.paivola.fi/lj.php"))
	// Parse the HTML from the reader
	doc, err := html.Parse(reader)
	// Check if there was an error
	if err != nil {
		// Print the error
		log.Printf("%s", err)
		// Return
		return
	}

	// Find the timetable table header node
	ttnode := findSpan("tr", "class", "header", doc)
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

			// Get the date of this day
			dateraw := strings.Split(dayentry.FirstChild.NextSibling.LastChild.Data, ".")
			dateraw[0] = strings.Split(dateraw[0], "\n")[1]

			// Parse the day from the date
			dateday, err1 := strconv.Atoi(dateraw[0])
			// Parse the month from the date
			datemonth, err2 := strconv.Atoi(dateraw[1])
			// Parse the year from the date
			dateyear, err3 := strconv.Atoi(dateraw[2])
			var date Date
			// If no errors came in parsing, create a Date struct from the parsed data
			// If there were errors, set the date to 1.1.1970
			if err1 == nil && err2 == nil && err3 == nil {
				date = Date{dateyear, datemonth, dateday}
			} else {
				date = Date{1970, 1, 1}
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

				// Save the parsed data to the correct location.
				switch lesson {
				case 0:
					firstyear[day][0] = TimetableLesson{data, "Aamu", date, Time{9, 0}}
				case 1:
					firstyear[day][1] = TimetableLesson{data, "IP1", date, Time{12, 15}}
				case 2:
					firstyear[day][2] = TimetableLesson{data, "IP2", date, Time{15, 0}}
				case 3:
					firstyear[day][3] = TimetableLesson{data, "Ilta", date, Time{19, 0}}
				case 4:
					other[day] = TimetableLesson{data, "Muuta", date, Time{0, 0}}
				case 5:
					secondyear[day][0] = TimetableLesson{data, "Aamu", date, Time{9, 0}}
				case 6:
					secondyear[day][1] = TimetableLesson{data, "IP1", date, Time{12, 15}}
				case 7:
					secondyear[day][2] = TimetableLesson{data, "IP2", date, Time{15, 0}}
				case 8:
					secondyear[day][3] = TimetableLesson{data, "Ilta", date, Time{19, 0}}
				}
			}
		}
		lastupdate = timestamp()
	} else {
		// Node not found, print error
		log.Printf(translate("timetable.update.failed"))
		lastupdate = 0
	}
}

func sendFirstYearStandalone(day int, bot *telebot.Bot, message telebot.Message) {
	bot.SendMessage(message.Chat,
		fmt.Sprintf(translate("timetable.generic"),
			firstyear[day][0].Subject, firstyear[day][1].Subject, firstyear[day][2].Subject, firstyear[day][3].Subject,
			DateToString(firstyear[day][0].Date))+"\n"+fmt.Sprintf(translate("timetable.other"), other[day].Subject),
		md)
}

func sendSecondYearStandalone(day int, bot *telebot.Bot, message telebot.Message) {
	bot.SendMessage(message.Chat,
		fmt.Sprintf(translate("timetable.generic"),
			secondyear[day][0].Subject, secondyear[day][1].Subject, secondyear[day][2].Subject, secondyear[day][3].Subject,
			DateToString(secondyear[day][0].Date))+"\n"+fmt.Sprintf(translate("timetable.other"), other[day].Subject),
		md)
}

func sendBothYears(day int, bot *telebot.Bot, message telebot.Message) {
	bot.SendMessage(message.Chat, translate("timetable.firstyear")+fmt.Sprintf(translate("timetable.generic"),
		firstyear[day][0].Subject, firstyear[day][1].Subject, firstyear[day][2].Subject, firstyear[day][3].Subject), md)
	bot.SendMessage(message.Chat, translate("timetable.secondyear")+fmt.Sprintf(translate("timetable.generic"),
		secondyear[day][0].Subject, secondyear[day][1].Subject, secondyear[day][2].Subject, secondyear[day][3].Subject), md)
	bot.SendMessage(message.Chat, fmt.Sprintf(translate("timetable.other"), other[day].Subject), md)
}
