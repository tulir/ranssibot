package posts

import (
	"errors"
	"github.com/tucnak/telebot"
	"golang.org/x/net/html"
	"io/ioutil"
	log "maunium.net/maulogger"
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
		lastRead, err := strconv.Atoi(strings.Split(string(lrData), "\n")[0])
		if lastRead == 0 || err != nil {
			log.Fatalf("Failed to find index of last read Ranssi post.")
			panic(err)
		}
		lastRead++

		node := getPost(lastRead)
		if node != nil {
			topic := strings.TrimSpace(node.FirstChild.FirstChild.Data)

			for _, uid := range subs {
				bot.SendMessage(whitelist.GetRecipientByUID(uid), lang.Translatef("posts.new", topic, lastRead), util.Markdown)
			}

			ioutil.WriteFile(lastreadpost, []byte(strconv.Itoa(lastRead)), 0700)
			time.Sleep(10 * time.Second)
		} else {
			time.Sleep(1 * time.Minute)
		}
	}
}

// Subscribe the given UID to the notification list.
func subscribe(uid int) bool {
	if isSubscribed(uid) {
		return false
	}
	subs = append(subs, uid)
	writeSubs()
	return true
}

// Unsubscribe the given UID from the notification list.
func unsubscribe(uid int) bool {
	if !isSubscribed(uid) {
		return false
	}
	for i, subuid := range subs {
		if subuid == uid {
			subs[i] = subs[len(subs)-1]
			subs = subs[:len(subs)-1]
		}
	}
	writeSubs()
	return true
}

func isSubscribed(uid int) bool {
	for _, id := range subs {
		if id == uid {
			return true
		}
	}
	return false
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
				log.Warnf("Failed to parse subscription entry %[1]s", str)
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
		log.Errorf("Error posting message \"%[1]s\" to the Ranssi post with ID %[2]d:\n%[3]s", message, id, err)
		return errors.New("Failed to post message")
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Failed to read response: %[1]s", err)
		return errors.New("Failed to read response")
	}
	return nil
}

// Get the content of the Ranssi post with the given ID.
func getPost(id int) *html.Node {
	data := util.HTTPGetMin("http://ranssi.paivola.fi/story.php?id=" + strconv.Itoa(id))
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
		if subscribe(message.Chat.ID) {
			bot.SendMessage(message.Chat, lang.Translate("posts.subscribed"), util.Markdown)
		} else {
			bot.SendMessage(message.Chat, lang.Translate("posts.alreadysubscribed"), util.Markdown)
		}
	} else if strings.EqualFold(args[1], "unsubscribe") || strings.EqualFold(args[1], "unsub") {
		if unsubscribe(message.Chat.ID) {
			bot.SendMessage(message.Chat, lang.Translate("posts.unsubscribed"), util.Markdown)
		} else {
			bot.SendMessage(message.Chat, lang.Translate("posts.notsubscribed"), util.Markdown)
		}
	} else if strings.EqualFold(args[1], "get") || strings.EqualFold(args[1], "read") {
		if len(args) < 3 {
			bot.SendMessage(message.Chat, lang.Translate("posts.read.usage"), util.Markdown)
			return
		}

		id, err := strconv.Atoi(args[2])
		if err != nil {
			bot.SendMessage(message.Chat, lang.Translatef("error.parse.integer", err), util.Markdown)
			return
		}
		post := getPost(id)
		if post == nil {
			bot.SendMessage(message.Chat, lang.Translatef("posts.notfound", id), util.Markdown)
			return
		}
		post = post.FirstChild

		title := strings.TrimSpace(post.FirstChild.Data)
		body := ""
		bodyNode := post.NextSibling.FirstChild
		prevbr := false
		for {
			if bodyNode == nil {
				break
			}
			if bodyNode.Data == "br" {
				if !prevbr {
					body += "\n"
					prevbr = true
					continue
				}
			} else {
				body += bodyNode.Data
			}
			prevbr = false
			bodyNode = bodyNode.NextSibling
			body = strings.TrimSpace(body)
		}

		bot.SendMessage(message.Chat, lang.Translatef("posts.read", id, title, body), util.Markdown)
	} else if strings.EqualFold(args[1], "comment") || strings.EqualFold(args[1], "message") || strings.EqualFold(args[1], "spam") {
		if len(args) < 4 {
			bot.SendMessage(message.Chat, lang.Translate("posts.spam.usage"), util.Markdown)
			return
		}

		id, err := strconv.Atoi(args[2])
		if err != nil {
			bot.SendMessage(message.Chat, lang.Translatef("error.parse.integer", err), util.Markdown)
			return
		}

		data := util.HTTPGet("http://ranssi.paivola.fi/story.php?id=" + strconv.Itoa(id))
		if string(data) == "ID:tä ei ole olemassa." {
			bot.SendMessage(message.Chat, lang.Translatef("posts.notfound", id), util.Markdown)
			return
		}
		doc, _ := html.Parse(strings.NewReader(data))
		if util.FindSpan("div", "id", "comments", doc) == nil {
			bot.SendMessage(message.Chat, lang.Translatef("posts.spam.nospamlist", id), util.Markdown)
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
		bot.SendMessage(message.Chat, lang.Translatef("posts.spam.sent", id, msg), util.Markdown)
	} else {
		bot.SendMessage(message.Chat, lang.Translate("posts.usage"), util.Markdown)
	}
}
