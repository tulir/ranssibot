package util

import (
	"bytes"
	"github.com/tucnak/telebot"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
)

// Markdown contains the SendOptions for markdown messages
var Markdown *telebot.SendOptions

// Init initializes things
func Init() {
	Markdown = new(telebot.SendOptions)
	Markdown.ParseMode = telebot.ModeMarkdown
}

// HTTPGet performs a HTTP GET request on the given URL
func HTTPGet(url string) string {
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

// Render renders the given HTML node to a string
func Render(node *html.Node) string {
	buf := new(bytes.Buffer)
	html.Render(buf, node)
	return buf.String()
}

// FindSpan finds a html element of the given type with the given key-value attribute from the given node
func FindSpan(typ string, key string, val string, node *html.Node) *html.Node {
	if node.Type == html.ElementNode && node.Data == typ {
		for _, attr := range node.Attr {
			if attr.Key == key && attr.Val == val {
				return node
			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		x := FindSpan(typ, key, val, c)
		if x != nil {
			return x
		}
	}
	return nil
}
