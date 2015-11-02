package main

import (
	"bytes"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
)

// Check if the given integer array contains the given integer.
func contains(list []int, i int) bool {
	for _, ii := range list {
		if ii == i {
			return true
		}
	}
	return false
}

// Perform a HTTP GET request on the given URL
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

// Render the given HTML node to a string
func render(node *html.Node) string {
	buf := new(bytes.Buffer)
	html.Render(buf, node)
	return buf.String()
}
