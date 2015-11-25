package util

import (
	"bytes"
	"github.com/tdewolff/minify"
	mhtml "github.com/tdewolff/minify/html"
	"github.com/tucnak/telebot"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/http"
)

// Markdown contains the SendOptions for markdown messages
var Markdown *telebot.SendOptions

func init() {
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

// HTTPGetStream performs a HTTP GET request on the given URL and returns a io.Reader pointer
func HTTPGetStream(url string) (io.Reader, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

// HTTPGetAndParse performs a HTTP GET request on the given URL, parses the output and returns a html.Node pointer
func HTTPGetAndParse(url string) (*html.Node, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	node, err := html.Parse(response.Body)
	if err != nil {
		return nil, err
	}
	return node, nil
}

// HTTPGetMin performs a HTTP GET request on the given URL and minifies the output
func HTTPGetMin(url string) string {
	response, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer response.Body.Close()
	var b bytes.Buffer

	minifyh(response.Body, &b)
	return b.String()
}

// HTTPGetMinStream performs a HTTP GET request on the given URL, minifies the output and returns a io.Reader pointer
func HTTPGetMinStream(url string) (io.Reader, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var b bytes.Buffer
	minifyh(response.Body, &b)
	return &b, nil
}

// HTTPGetMinAndParse performs a HTTP GET request on the given URL, minifies and parses the output and returns a html.Node pointer
func HTTPGetMinAndParse(url string) (*html.Node, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var b bytes.Buffer
	minifyh(response.Body, &b)

	node, err := html.Parse(&b)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func minifyh(reader io.Reader, writer io.Writer) {
	mhtml.Minify(minify.New(), writer, reader, nil)
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
		if key == "" && val == "" {
			return node
		}
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
