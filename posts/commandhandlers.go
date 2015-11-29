package posts

import (
	"github.com/tucnak/telebot"
	"golang.org/x/net/html"
	"maunium.net/ranssibot/lang"
	"maunium.net/ranssibot/util"
	"strconv"
	"strings"
)

func handleSubscribe(bot *telebot.Bot, message telebot.Message, args []string) {
	if subscribe(message.Chat.ID) {
		bot.SendMessage(message.Chat, lang.Translate("posts.subscribed"), util.Markdown)
	} else {
		bot.SendMessage(message.Chat, lang.Translate("posts.alreadysubscribed"), util.Markdown)
	}
}
func handleUnsubscribe(bot *telebot.Bot, message telebot.Message, args []string) {
	if unsubscribe(message.Chat.ID) {
		bot.SendMessage(message.Chat, lang.Translate("posts.unsubscribed"), util.Markdown)
	} else {
		bot.SendMessage(message.Chat, lang.Translate("posts.notsubscribed"), util.Markdown)
	}
}

func handleNews(bot *telebot.Bot, message telebot.Message, args []string) {
	if util.Timestamp()-lastupdate > 1500 {
		updateNews()
	}
	var entries string
	for _, item := range news.Items {
		entries += lang.Translatef("posts.latest.entry", item.Title, item.Link, item.Summary, item.Date.Format("15:04:05 02.01.2006"))
	}
	bot.SendMessage(message.Chat, lang.Translatef("posts.latest", entries), util.Markdown)
}

func handleRead(bot *telebot.Bot, message telebot.Message, args []string) {
	if len(args) < 1 {
		bot.SendMessage(message.Chat, lang.Translate("posts.read.usage"), util.Markdown)
		return
	}

	id, err := strconv.Atoi(args[0])
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
	body := util.RenderText(post.NextSibling)

	bot.SendMessage(message.Chat, lang.Translatef("posts.read", id, title, body), util.Markdown)
}

func handleReadComments(bot *telebot.Bot, message telebot.Message, args []string) {
	bot.SendMessage(message.Chat, "Not yet implemented", nil)
}

func handleComment(bot *telebot.Bot, message telebot.Message, args []string) {
	if len(args) < 2 {
		bot.SendMessage(message.Chat, lang.Translate("posts.spam.usage"), util.Markdown)
		return
	}

	id, err := strconv.Atoi(args[0])
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
	for _, str := range args[1:] {
		msg += str + " "
	}
	msg = strings.TrimSpace(msg)
	err = spam(id, msg)
	if err != nil {
		bot.SendMessage(message.Chat, err.Error(), nil)
	}
	bot.SendMessage(message.Chat, lang.Translatef("posts.spam.sent", id, msg), util.Markdown)
}
