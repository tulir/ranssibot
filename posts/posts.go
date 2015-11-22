package posts

import (
	"errors"
	"fmt"
	"github.com/tucnak/telebot"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"maunium.net/ranssibot/lang"
	"maunium.net/ranssibot/util"
	"maunium.net/ranssibot/whitelist"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const lastreadpost = "data/lastreadpost"
const postsubs = "data/postsubs"

var subs = []int{}

// Loop is an infinite loop that checks for new Ranssi posts
func Loop(bot *telebot.Bot) {
	readSubs()
	for {
		lrData, _ := ioutil.ReadFile(lastreadpost)
		lastRead, _ := strconv.Atoi(strings.Split(string(lrData), "\n")[0])
		if lastRead == 0 {
			log.Println("Failed to find index of last read Ranssi post.")
			return
		}
		lastRead++

		node := getPost(lastRead)
		if node != nil {
			topic := strings.TrimSpace(node.FirstChild.NextSibling.FirstChild.Data)

			for _, uid := range subs {
				bot.SendMessage(whitelist.GetRecipientByUID(uid), fmt.Sprintf(lang.Translate("posts.new"), topic, lastRead), util.Markdown)
			}

			ioutil.WriteFile(lastreadpost, []byte(strconv.Itoa(lastRead)), 0700)
			time.Sleep(10 * time.Second)
		} else {
			time.Sleep(1 * time.Minute)
		}
	}
}

// Subscribe the given UID to the notification list.
func subscribe(uid int) {
	subs = append(subs, uid)
	writeSubs()
}

// Unsubscribe the given UID from the notification list.
func unsubscribe(uid int) {
	for i, subuid := range subs {
		if subuid == uid {
			subs[i] = subs[len(subs)-1]
			subs = subs[:len(subs)-1]
		}
	}
	writeSubs()
}

// Read the UIDs that are subscribed to the notification list.
func readSubs() {
	subsData, _ := ioutil.ReadFile(postsubs)
	subsRaw := strings.Split(string(subsData), "\n")
	for _, str := range subsRaw {
		if len(str) != 0 && !strings.HasPrefix(str, "#") {
			uid, err := strconv.Atoi(str)
			if err == nil {
				subs = append(subs, uid)
			} else {
				log.Println("Failed to parse subscription entry for " + str)
			}
		}
	}
}

// Write the UIDs that are subscribed to the notification list.
func writeSubs() {
	s := ""
	for _, uid := range subs {
		s += strconv.Itoa(uid) + "\n"
	}
	ioutil.WriteFile(postsubs, []byte(s), 0700)
}

func spam(id int, message string) error {
	resp, err := http.PostForm("http://ranssi.paivola.fi/story.php?id="+strconv.Itoa(id), url.Values{"comment": {message}})
	if err != nil {
		log.Println("Error posting message \"" + message + "\" to the Ranssi post with ID " + strconv.Itoa(id) + ":\n" + err.Error())
		return errors.New("Failed to post message")
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read response: " + err.Error())
		return errors.New("Failed to read response")
	}
	return nil
}

// Get the content of the Ranssi post with the given ID.
func getPost(id int) *html.Node {
	data := util.HTTPGet("http://ranssi.paivola.fi/story.php?id=" + strconv.Itoa(id))
	if string(data) != "ID:tä ei ole olemassa." {
		doc, _ := html.Parse(strings.NewReader(data))
		return util.FindSpan("div", "id", "story", doc)
	}
	return nil
}

// HandleCommand handles Ranssi post commands
func HandleCommand(bot *telebot.Bot, message telebot.Message, args []string) {
	if len(args) < 2 {
		bot.SendMessage(message.Chat, lang.Translate("posts.usage"), util.Markdown)
		return
	}

	if strings.EqualFold(args[1], "subscribe") || strings.EqualFold(args[1], "sub") {
		subscribe(message.Chat.ID)
		bot.SendMessage(message.Chat, lang.Translate("posts.subscribed"), util.Markdown)
	} else if strings.EqualFold(args[1], "unsubscribe") || strings.EqualFold(args[1], "unsub") {
		unsubscribe(message.Chat.ID)
		bot.SendMessage(message.Chat, lang.Translate("posts.unsubscribed"), util.Markdown)
	} else if strings.EqualFold(args[1], "get") || strings.EqualFold(args[1], "read") {
		bot.SendMessage(message.Chat, "*Error:* Reading Ranssi posts has not yet been implemented", util.Markdown)
	} else if strings.EqualFold(args[1], "comment") || strings.EqualFold(args[1], "message") || strings.EqualFold(args[1], "spam") {
		if len(args) < 4 {
			bot.SendMessage(message.Chat, lang.Translate("posts.spam.usage"), util.Markdown)
			return
		}

		id, err := strconv.Atoi(args[2])
		if err != nil {
			bot.SendMessage(message.Chat, fmt.Sprintf(lang.Translate("error.parse.integer"), err), util.Markdown)
			return
		}

		data := util.HTTPGet("http://ranssi.paivola.fi/story.php?id=" + strconv.Itoa(id))
		if string(data) == "ID:tä ei ole olemassa." {
			bot.SendMessage(message.Chat, fmt.Sprintf(lang.Translate("posts.spam.notfound"), id), util.Markdown)
			return
		}
		doc, _ := html.Parse(strings.NewReader(data))
		if util.FindSpan("div", "id", "comments", doc) == nil {
			bot.SendMessage(message.Chat, fmt.Sprintf(lang.Translate("posts.spam.nospamlist"), id), util.Markdown)
			return
		}

		msg := ""
		for _, str := range args[3:] {
			msg += str + " "
		}
		msg = strings.TrimSpace(msg)
		err = spam(id, msg)
		if err != nil {
			bot.SendMessage(message.Chat, err.Error(), nil)
		}
		bot.SendMessage(message.Chat, fmt.Sprintf(lang.Translate("posts.spam.sent"), id, msg), util.Markdown)
	} else {
		bot.SendMessage(message.Chat, lang.Translate("posts.usage"), util.Markdown)
	}
}
