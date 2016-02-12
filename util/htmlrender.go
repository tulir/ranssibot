package util

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
)

type writer interface {
	io.Writer
	io.ByteWriter
	WriteString(string) (int, error)
}

func renderHTML(w writer, n *html.Node) {
	switch n.Type {
	case html.TextNode:
		w.WriteString(n.Data)
	case html.ElementNode:
		switch n.Data {
		case "br":
			w.WriteString("\n")
		case "b":
			w.WriteString("*")
			renderChildren(w, n)
			w.WriteString("*")
		case "tt":
			w.WriteString("`")
			renderChildren(w, n)
			w.WriteString("`")
		case "a":
			w.WriteString("[")
			renderChildren(w, n)
			w.WriteString("]")
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					w.WriteString(fmt.Sprintf("(%s)", attr.Val))
					break
				}
			}
		default:
			renderChildren(w, n)
		}
	case html.DocumentNode:
		renderChildren(w, n)
	}
}

func renderChildren(w writer, n *html.Node) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		renderHTML(w, c)
	}
}
