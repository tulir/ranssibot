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

var timetable = [26][9]string{}
var today = 5
var lastupdate = timestamp()

//var timetable = make([]string, 9)

func main() {
	md := new(telebot.SendOptions)
	md.ParseMode = telebot.ModeMarkdown

	// Load the whitelist
	whitelist := loadWhitelist()

	// Connect to Telegram
	bot, err := telebot.NewBot("132300126:AAHps1NPAj9Y7qTBbDGlGsyuMGoMtk__Qa8")
	if err != nil {
		log.Printf("Error connecting to Telegram!\n" + err.Error())
		return
	}

	messages := make(chan telebot.Message)
	bot.Listen(messages, 1*time.Second)

	log.Printf("Connected to Telegram!")

	updateTimes()

MainLoop:
	for message := range messages {
		if !contains(whitelist, message.Sender.ID) {
			bot.SendMessage(message.Chat, "Et ole Päivölän Lukujärjestysbotin whitelistillä. "+
				"Voit tökkiä Tuliria päästäksesi whitelistille.\n"+
				"Telegram-käyttäjäsi ID on "+strconv.Itoa(message.Sender.ID), nil)
			bot.SendMessage(message.Chat, "", nil)
			continue MainLoop
		}
		log.Printf(message.Sender.FirstName + " is spamming me with " + message.Text)
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
						continue MainLoop
					}
					day += shift
					if day < 0 || day >= len(timetable) {
						bot.SendMessage(message.Chat, "I'm limited to the data shown on http://ranssi.paivola.fi/lj.php, so I can't show the timetables that far away.", md)
						continue MainLoop
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
}

func render(node *html.Node) string {
	buf := new(bytes.Buffer)
	html.Render(buf, node)
	return buf.String()
}

func timestamp() int64 {
	return int64(time.Now().Unix())
}

func updateTimes() {
	reader := strings.NewReader(httpGet("http://ranssi.paivola.fi/lj.php"))
	doc, err := html.Parse(reader)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	ttnode := findSpan("tr", "class", "header", doc)
	if ttnode != nil {
		dayentry := ttnode
		for day := 0; ; day++ {
			if dayentry.NextSibling == nil || dayentry.NextSibling.NextSibling == nil || dayentry.NextSibling.NextSibling.FirstChild == nil {
				break
			}
			dayentry = dayentry.NextSibling.NextSibling
			//println(render(dayentry))
			entry := dayentry.FirstChild.NextSibling
			for lesson := 0; lesson < 9; lesson++ {
				entry = entry.NextSibling.NextSibling
				if entry == nil {
					break
				}

				if entry.FirstChild != nil {
					if entry.FirstChild.Type == html.TextNode {
						timetable[day][lesson] = entry.FirstChild.Data
					} else if entry.FirstChild.Type == html.ElementNode {
						if entry.FirstChild.FirstChild != nil {
							if entry.FirstChild.FirstChild.Type == html.TextNode {
								timetable[day][lesson] = entry.FirstChild.FirstChild.Data
							}
						}
					}
				} else {
					timetable[day][lesson] = "tyhjää"
				}
			}
			print("\n")
		}
	}
}

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
