package main

import (
	"github.com/tucnak/telebot"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

// Timetable cahce
var timetable = [26][9]string{}

// The day ID for today
var today = 5

// The last time the timetable cache was updated
var lastupdate = timestamp()

// List of UIDs that are allowed to use the bot
var whitelist []int

// The markdown send options
var md *telebot.SendOptions

func main() {
	md = new(telebot.SendOptions)
	md.ParseMode = telebot.ModeMarkdown

	// Load the whitelist
	whitelist = loadWhitelist()

	// Connect to Telegram
	bot, err := telebot.NewBot("132300126:AAHps1NPAj9Y7qTBbDGlGsyuMGoMtk__Qa8")
	if err != nil {
		log.Printf("Error connecting to Telegram: %s", err)
		return
	}
	messages := make(chan telebot.Message)
	// Enable message listener
	bot.Listen(messages, 1*time.Second)
	// Print "connected" message
	log.Printf("Connected to Telegram!")

	// Update timetables
	updateTimes()

	// Listen to messages
	for message := range messages {
		handleCommand(bot, message)
	}
}

// Handle a command
func handleCommand(bot *telebot.Bot, message telebot.Message) {
	if !contains(whitelist, message.Sender.ID) {
		bot.SendMessage(message.Chat, "Et ole Päivölän Lukujärjestysbotin whitelistillä. "+
			"Voit tökkiä Tuliria päästäksesi whitelistille.\n"+
			"Telegram-käyttäjäsi ID on "+strconv.Itoa(message.Sender.ID), nil)
		bot.SendMessage(message.Chat, "", nil)
		return
	}
	log.Printf("%s (%d) sent command: %s", message.Sender.Username, message.Sender.ID, message.Text)
	if strings.HasPrefix(message.Text, "Mui.") {
		bot.SendMessage(message.Chat, "Mui. "+message.Sender.FirstName+".", nil)
	} else if strings.HasPrefix(message.Text, "/timetable") {
		if timestamp() > lastupdate+600 {
			bot.SendMessage(message.Chat, "Updating cached timetables...", md)
			updateTimes()
		}
		args := strings.Split(message.Text, " ")
		if len(args) > 1 {
			day := today
			if len(args) > 2 {
				shift, err := strconv.Atoi(args[2])
				if err != nil {
					bot.SendMessage(message.Chat, "I couldn't parse an integer from \"_"+args[2]+"_\"", md)
					return
				}
				day += shift
				if day < 0 || day >= len(timetable) {
					bot.SendMessage(message.Chat, "I'm limited to the data shown on http://ranssi.paivola.fi/lj.php, so I can't show the timetables that far away.", md)
					return
				}
			}
			if strings.EqualFold(args[1], "ventit") {
				bot.SendMessage(message.Chat,
					"Aamu: "+timetable[day][0]+
						"\nIP1: "+timetable[day][1]+
						"\nIP2: "+timetable[day][2]+
						"\n"+"Ilta: "+timetable[day][3],
					md)
				bot.SendMessage(message.Chat, "Muuta: "+timetable[day][4], md)
			} else if strings.EqualFold(args[1], "neliöt") {
				bot.SendMessage(message.Chat,
					"Aamu: "+timetable[day][5]+
						"\nIP1: "+timetable[day][6]+
						"\nIP2: "+timetable[day][7]+
						"\n"+"Ilta: "+timetable[day][8],
					md)
				bot.SendMessage(message.Chat, "Muuta: "+timetable[day][4], md)
			} else {
				bot.SendMessage(message.Chat, "*Usage:* /timetable <neliöt/ventit> [day offset]", md)
			}
		} else {
			bot.SendMessage(message.Chat, "*Usage:* /timetable <neliöt/ventit> [day offset]", md)
		}
	} else if message.Text == "/update" {
		updateTimes()
		bot.SendMessage(message.Chat, "Updated timetables successfully", nil)
	} else if strings.HasPrefix(message.Text, "/") {
		bot.SendMessage(message.Chat, "Komentoa ei tunnistettu.", nil)
	}
}

// Get the current UNIX timestamp
func timestamp() int64 {
	return int64(time.Now().Unix())
}

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
			// Get the next day node
			dayentry = dayentry.NextSibling.NextSibling
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

				// Check if the lesson contains anything
				if entry.FirstChild != nil {
					// Lesson data found; Try to parse it
					if entry.FirstChild.Type == html.TextNode {
						// Found lesson data directly under lesson node
						timetable[day][lesson] = entry.FirstChild.Data
					} else if entry.FirstChild.Type == html.ElementNode {
						// Didn't find data directly under lesson node
						// Check for a child element node.
						if entry.FirstChild.FirstChild != nil {
							// Child element node found. Check if the child of that child is text.
							if entry.FirstChild.FirstChild.Type == html.TextNode {
								// Child of child is text, use it as the data.
								timetable[day][lesson] = entry.FirstChild.FirstChild.Data
							}
						}
					} else {
						// Lesson data couldn't be parsed
						timetable[day][lesson] = "tyhjää"
					}
				} else {
					// Lesson is empty
					timetable[day][lesson] = "tyhjää"
				}
			}
		}
		lastupdate = timestamp()
	} else {
		// Node not found, print error
		log.Printf("Error updating timetables: Failed to find timetable table header node!")
		lastupdate = 0
	}
}

// Find a html element of the given type with the given key-value attribute from the given node
func findSpan(typ string, key string, val string, node *html.Node) *html.Node {
	if node.Type == html.ElementNode && node.Data == typ {
		for _, attr := range node.Attr {
			if attr.Key == key && attr.Val == val {
				return node
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		x := findSpan(typ, key, val, c)
		if x != nil {
			return x
		}
	}
	return nil
}

// Load the whitelist from file
func loadWhitelist() []int {
	// Read the file
	wldata, err := ioutil.ReadFile("whitelist.txt")
	// Check if there was an error
	if err != nil {
		// Error, print message and use hardcoded whitelist.
		log.Printf("Failed to load whitelist: %s; Using hardcoded version", err)
		return []int{
			84359547,  /* Tulir */
			67147746,  /* Ege   */
			128602828, /* Max   */
			124500539, /* Galax */
			54580303,  /* Antti */
			115187137, /* Å     */

		}
	}
	// No error, parse the data
	log.Printf("Loading whitelist...")
	// Split the file string to an array of lines
	wlraw := strings.Split(string(wldata), "\n")
	// Make the whitelist array
	whitelist := make([]int, len(wlraw), cap(wlraw))
	// Loop through the lines from the file
	for i := 0; i < len(wlraw); i++ {
		// Make sure the line is not empty
		if len(wlraw[i]) == 0 {
			continue
		}
		// Split the entry to UID and name
		entry := strings.Split(wlraw[i], "-")
		// Convert the UID string to an integer
		id, converr := strconv.Atoi(entry[0])
		// Make sure the conversion didn't fail
		if converr == nil {
			// No errors, add the UID to the whitelist
			whitelist[i] = id
			log.Printf("Added %s (ID %s) to the whitelist.", entry[1], entry[0])
		} else {
			// Error occured, print message
			log.Printf("Failed to parse %s: %s", wlraw[i], err)
		}
	}
	log.Printf("Finished loading whitelist")
	return whitelist
}
