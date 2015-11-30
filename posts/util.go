package posts

import (
	"errors"
	"github.com/SlyMarbo/rss"
	"github.com/tucnak/telebot"
	"golang.org/x/net/html"
	"io/ioutil"
	log "maunium.net/go/maulogger"
	"maunium.net/go/ranssibot/config"
	"maunium.net/go/ranssibot/lang"
	"maunium.net/go/ranssibot/util"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	subSetting = "posts-subscription"
)

var lastupdate int64
var news *rss.Feed

func init() {
	var err error
	news, err = rss.Fetch("http://ranssi.paivola.fi/rss.php")
	if err != nil {
		panic(err)
	}
	lastupdate = util.Timestamp()
}

func updateNews() {
	err := news.Update()
	if err != nil {
		log.Errorf("Failed to update Ranssi News: %s", err)
	}
}

// Loop is an infinite loop that checks for new Ranssi posts
func Loop(bot *telebot.Bot, noNotifAtInit bool) {
	for {
		readNow := config.GetConfig().LastReadPost + 1

		node := getPost(readNow)
		if node != nil {
			topic := strings.TrimSpace(node.FirstChild.FirstChild.Data)

			log.Infof("New Ranssi post detected: %s (ID %d)", topic, readNow)

			if !noNotifAtInit {
				for _, user := range config.GetUsersWithSetting(subSetting, "true") {
					bot.SendMessage(user, lang.Translatef("posts.new", topic, readNow), util.Markdown)
				}
			}

			config.GetConfig().LastReadPost = readNow
			config.ASave()
			updateNews()
			time.Sleep(5 * time.Second)
			continue
		}
		noNotifAtInit = false
		time.Sleep(1 * time.Minute)
	}
}

// Subscribe the given UID to the notification list.
func subscribe(uid int) bool {
	if isSubscribed(uid) {
		log.Debugf("%[1]d attempted to subscribe to the notification list, but was already subscribed", uid)
		return false
	}
	log.Debugf("[Posts] %[1]d successfully subscribed to the notifcation list", uid)
	config.GetUserWithUID(uid).SetSetting(subSetting, "true")
	config.ASave()
	return true
}

// Unsubscribe the given UID from the notification list.
func unsubscribe(uid int) bool {
	if !isSubscribed(uid) {
		log.Debugf("%[1]d attempted to unsubscribe from the notification list, but was not subscribed", uid)
		return false
	}
	log.Debugf("%[1]d successfully unsubscribed from the notifcation list", uid)
	config.GetUserWithUID(uid).RemoveSetting(subSetting)
	config.ASave()
	return true
}

func isSubscribed(uid int) bool {
	return config.GetUserWithUID(uid).HasSetting(subSetting)
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
