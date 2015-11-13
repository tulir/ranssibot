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
	langdata, err := ioutil.ReadFile("ljbot.lang")
	// Check if there was an error
	if err != nil {
		// Error, print message and use hardcoded whitelist.
		log.Fatalf("Failed to load language: %s; Using hardcoded version", err)
	}
	// No error, parse the data
	log.Printf("Loading language...")
	// Split the file string to an array of lines
	langraw := strings.Split(string(langdata), "\n")
	for i := 0; i < len(langraw); i++ {
		// Make sure the line is not empty
		if len(langraw[i]) == 0 || strings.HasPrefix(langraw[i], "#") {
			continue
		}
		entry := strings.Split(langraw[i], "=")
		lang[entry[0]] = entry[1]
	}
}

// Translate translates the given key.
func Translate(key string) string {
	value, exists := lang[key]
	if exists {
		return strings.Replace(value, "<br>", "\n", -1)
	}
	return key
}
