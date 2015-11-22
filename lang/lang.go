package lang

import (
	"io/ioutil"
	"log"
	"strings"
)

var lang = make(map[string]string)

// Load loads the language from file.
func Load() {
	// Read the file
	langdata, err := ioutil.ReadFile("data/lang/en_US.lang")
	// Check if there was an error
	if err != nil {
		// Error, print message.
		log.Fatalf("Failed to load language: %s", err)
	}
	// No error, parse the data
	log.Printf("Loading language...")
	// Split the file string to an array of lines
	langraw := strings.Split(string(langdata), "\n")
	var appendTo string

	for i := 0; i < len(langraw); i++ {
		// Make sure the line is not empty
		if len(langraw[i]) == 0 || strings.HasPrefix(langraw[i], "#") {
			continue
		}
		if len(appendTo) != 0 {
			entry := langraw[i]
			entry = strings.TrimSpace(entry)
			appendToCache := appendTo
			if strings.HasSuffix(entry, "\\") {
				entry = trimSuffix(entry, "\\")
			} else {
				appendTo = ""
			}
			if len(lang[appendToCache]) == 0 {
				lang[appendToCache] = entry
			} else {
				lang[appendToCache] += "\n" + entry
			}
		} else {
			entry := strings.Split(langraw[i], "=")
			entry[1] = strings.TrimSpace(entry[1])
			if strings.HasSuffix(entry[1], "\\") {
				entry[1] = trimSuffix(entry[1], "\\")
				appendTo = entry[0]
			}
			lang[entry[0]] = entry[1]
		}
	}
}

func trimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

// Translate translates the given key.
func Translate(key string) string {
	value, exists := lang[key]
	if exists {
		return strings.Replace(value, "<br>", "\n", -1)
	}
	return key
}
