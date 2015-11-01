package main

import (
	"bytes"
	"fmt"
	"github.com/tucnak/telebot"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var timetable = make([]string, 9)

func main() {
	md := new(telebot.SendOptions)
	md.ParseMode = telebot.ModeMarkdown

	// Load the whitelist
	whitelist := loadWhitelist()

	// Connect to Telegram
	bot, err := telebot.NewBot("132300126:AAHps1NPAj9Y7qTBbDGlGsyuMGoMtk__Qa8")
	if err != nil {
		log.Printf("Error connecting to Telegram!\n")
		return
	}

	messages := make(chan telebot.Message)
	bot.Listen(messages, 1*time.Second)

	for message := range messages {
		if !contains(whitelist, message.Sender.ID) {
			bot.SendMessage(message.Chat, "Et ole Päivölän Lukujärjestysbotin whitelistillä. "+
				"Voit tökkiä Tuliria päästäksesi whitelistille.\n"+
				"Telegram-käyttäjäsi ID on "+strconv.Itoa(message.Sender.ID), nil)
			bot.SendMessage(message.Chat, "", nil)
			continue
		}
		log.Printf(message.Sender.FirstName + " is spamming me with " + message.Text)
		if strings.HasPrefix(message.Text, "Mui.") {
			bot.SendMessage(message.Chat, "Mui. "+message.Sender.FirstName+".", nil)
		} else if strings.HasPrefix(message.Text, "/today") {
			args := strings.Split(message.Text, " ")
			if len(args) > 1 {
				if strings.EqualFold(args[1], "ventit") {
					bot.SendMessage(message.Chat,
						"Aamu: "+timetable[0]+
							"\nIP1: "+timetable[1]+
							"\nIP2: "+timetable[2]+
							"\n"+"Ilta: "+timetable[3],
						md)
					bot.SendMessage(message.Chat, "Muuta: "+timetable[4], md)
				} else if strings.EqualFold(args[1], "neliöt") {
					bot.SendMessage(message.Chat,
						"Aamu: "+timetable[5]+
							"\nIP1: "+timetable[6]+
							"\nIP2: "+timetable[7]+
							"\n"+"Ilta: "+timetable[8],
						md)
					bot.SendMessage(message.Chat, "Muuta: "+timetable[4], md)
				} else {
					bot.SendMessage(message.Chat, "*Usage:* /today <neliöt/ventit>", md)
				}
			} else {
				bot.SendMessage(message.Chat, "*Usage:* /today <neliöt/ventit>", md)
			}
		} else if message.Text == "/update" {
			updateTimes()
			bot.SendMessage(message.Chat, "Updated timetables successfully", nil)
		} else if strings.HasPrefix(message.Text, "/") {
			bot.SendMessage(message.Chat, "Komentoa ei tunnistettu.", nil)
		}
	}
}

func render(node *html.Node) string {
	buf := new(bytes.Buffer)
	html.Render(buf, node)
	return buf.String()
}

func updateTimes() {
	reader := strings.NewReader(httpGet("http://ranssi.paivola.fi/lj.php"))
	doc, err := html.Parse(reader)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	today := findTodaysTimetables(doc)
	if today != nil {
		entry := today.FirstChild.NextSibling
		for i := 0; i < 9; i++ {
			entry = entry.NextSibling
			if entry == nil {
				break
			}

			if entry.FirstChild != nil {
				if entry.FirstChild.Type == html.TextNode {
					timetable[i] = entry.FirstChild.Data
				} else if entry.FirstChild.Type == html.ElementNode {
					if entry.FirstChild.FirstChild != nil {
						if entry.FirstChild.FirstChild.Type == html.TextNode {
							timetable[i] = entry.FirstChild.FirstChild.Data
						}
					}
				}
			} else if entry.Type == html.ElementNode {
				timetable[i] = "tyhjää"
			} else {
				i--
			}
		}
	}
}

func findTodaysTimetables(node *html.Node) *html.Node {
	if node.Type == html.ElementNode && node.Data == "tr" {
		for _, attr := range node.Attr {
			if attr.Key == "class" && attr.Val == "today" {
				return node
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		x := findTodaysTimetables(c)
		if x != nil {
			return x
		}
	}
	return nil
}

func loadWhitelist() []int {
	wldata, err := ioutil.ReadFile("whitelist.txt")
	if err != nil {
		fmt.Printf("Failed to load whitelist: %s; Using hardcoded version", err)
		return []int{
			84359547,  /* Tulir */
			67147746,  /* Ege */
			128602828, /* Max */
			124500539, /* Galax */
		}
	}
	println("Loading whitelist...")
	wlraw := strings.Split(string(wldata), "\n")
	whitelist := make([]int, len(wlraw), cap(wlraw))
	for i := 0; i < len(wlraw); i++ {
		if len(wlraw[i]) == 0 {
			continue
		}
		entry := strings.Split(wlraw[i], "-")
		id, converr := strconv.Atoi(entry[0])
		if converr == nil {
			whitelist[i] = id
			println("Added " + entry[1] + " (ID " + entry[0] + ") to the whitelist.")
		} else {
			fmt.Printf("Failed to parse "+wlraw[i]+": %s", err)
		}
	}
	println("Finished loading whitelist")
	return whitelist
}

func contains(list []int, i int) bool {
	for _, ii := range list {
		if ii == i {
			return true
		}
	}
	return false
}

func httpGet(url string) string {
	response, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return ""
	}
	return string(contents)
}
