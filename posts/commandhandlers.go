package posts

import (
	"bytes"
	"github.com/tucnak/telebot"
	"golang.org/x/net/html"
	"maunium.net/go/ranssibot/config"
	"maunium.net/go/ranssibot/lang"
	"maunium.net/go/ranssibot/util"
	"strconv"
	"strings"
	"time"
)

func handleSubscribe(bot *telebot.Bot, message telebot.Message, args []string) {
	sender := config.GetUserWithUID(message.Sender.ID)
	if subscribe(sender) {
		bot.SendMessage(message.Chat, lang.Translate(sender, "posts.subscribed"), util.Markdown)
	} else {
		bot.SendMessage(message.Chat, lang.Translate(sender, "posts.alreadysubscribed"), util.Markdown)
	}
}
func handleUnsubscribe(bot *telebot.Bot, message telebot.Message, args []string) {
	sender := config.GetUserWithUID(message.Sender.ID)
	if unsubscribe(sender) {
		bot.SendMessage(message.Chat, lang.Translate(sender, "posts.unsubscribed"), util.Markdown)
	} else {
		bot.SendMessage(message.Chat, lang.Translate(sender, "posts.notsubscribed"), util.Markdown)
	}
}

func handleNews(bot *telebot.Bot, message telebot.Message, args []string) {
	sender := config.GetUserWithUID(message.Sender.ID)
	if util.Timestamp()-lastupdate > 1500 {
		updateNews()
	}
	var buffer bytes.Buffer
	for _, item := range news.Items {
		buffer.WriteString(lang.Translatef(sender, "posts.latest.entry", item.Title, item.Link, item.Summary, item.Date.Format("15:04:05 02.01.2006")))
	}
	bot.SendMessage(message.Chat, lang.Translatef(sender, "posts.latest", buffer.String()), util.Markdown)
}

func handleRead(bot *telebot.Bot, message telebot.Message, args []string) {
	sender := config.GetUserWithUID(message.Sender.ID)
	if len(args) < 1 {
		bot.SendMessage(message.Chat, lang.Translate(sender, "posts.read.usage"), util.Markdown)
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		bot.SendMessage(message.Chat, lang.Translatef(sender, "error.parse.integer", err), util.Markdown)
		return
	}
	post := getPost(id)
	if post == nil {
		bot.SendMessage(message.Chat, lang.Translatef(sender, "posts.notfound", id), util.Markdown)
		return
	}
	post = post.FirstChild

	title := strings.TrimSpace(post.FirstChild.Data)
	body := util.RenderText(post.NextSibling)
	time, _ := time.Parse("2006-01-02 15:04:05", strings.TrimSpace(post.NextSibling.NextSibling.FirstChild.NextSibling.Data))
	bot.SendMessage(message.Chat, lang.Translatef(sender, "posts.read", id, title, body, time.Format("15:04:05 02.01.2006")), util.Markdown)
}

func handleReadComments(bot *telebot.Bot, message telebot.Message, args []string) {
	bot.SendMessage(message.Chat, "Not yet implemented", nil)
}

func handleComment(bot *telebot.Bot, message telebot.Message, args []string) {
	sender := config.GetUserWithUID(message.Sender.ID)
	if len(args) < 2 {
		bot.SendMessage(message.Chat, lang.Translate(sender, "posts.spam.usage"), util.Markdown)
		return
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		bot.SendMessage(message.Chat, lang.Translatef(sender, "error.parse.integer", err), util.Markdown)
		return
	}

	data := util.HTTPGet("http://ranssi.paivola.fi/story.php?id=" + strconv.Itoa(id))
	if string(data) == "ID:tä ei ole olemassa." {
		bot.SendMessage(message.Chat, lang.Translatef(sender, "posts.notfound", id), util.Markdown)
		return
	}
	doc, _ := html.Parse(strings.NewReader(data))
	if util.FindSpan("div", "id", "comments", doc) == nil {
		bot.SendMessage(message.Chat, lang.Translatef(sender, "posts.spam.nospamlist", id), util.Markdown)
		return
	}
	var buffer bytes.Buffer
	for _, str := range args[1:] {
		buffer.WriteString(str)
		buffer.WriteString(" ")
	}
	msg := strings.TrimSpace(buffer.String())
	err = spam(id, msg)
	if err != nil {
		bot.SendMessage(message.Chat, err.Error(), nil)
	}
	bot.SendMessage(message.Chat, lang.Translatef(sender, "posts.spam.sent", id, msg), util.Markdown)
}
