package posts

import (
	"errors"
	"github.com/SlyMarbo/rss"
	"github.com/tucnak/telebot"
	"golang.org/x/net/html"
	"io/ioutil"
	log "maunium.net/maulogger"
	"maunium.net/ranssibot/config"
	"maunium.net/ranssibot/lang"
	"maunium.net/ranssibot/util"
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
func Loop(bot *telebot.Bot) {
	for {
		lastRead := config.GetConfig().LastReadPost
		/*lrData, _ := ioutil.ReadFile(lastreadpost)
		lastRead, err := strconv.Atoi(strings.Split(string(lrData), "\n")[0])
		if lastRead == 0 || err != nil {
			log.Fatalf("Failed to find index of last read Ranssi post.")
			panic(err)
		}*/
		lastRead++

		node := getPost(lastRead)
		if node != nil {
			topic := strings.TrimSpace(node.FirstChild.FirstChild.Data)

			for _, user := range config.GetUsersWithSetting(subSetting, "true") {
				bot.SendMessage(user, lang.Translatef("posts.new", topic, lastRead), util.Markdown)
			}

			/*ioutil.WriteFile(lastreadpost, []byte(strconv.Itoa(lastRead)), 0700)*/
			config.GetConfig().LastReadPost = lastRead
			config.ASave()
			time.Sleep(10 * time.Second)
			updateNews()
		} else {
			time.Sleep(1 * time.Minute)
		}
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
