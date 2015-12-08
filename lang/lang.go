package lang

import (
	"fmt"
	"io/ioutil"
	log "maunium.net/go/maulogger"
	"strings"
)

// Language is a language.
type Language struct {
	Name string
	Data map[string]string
}

var languages []*Language

// Load loads the language data.
func Load() {
	files, _ := ioutil.ReadDir("./lang")
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".lang") {
			continue
		}
		filename := "lang/" + f.Name()
		langname := trimSuffix(f.Name(), ".lang")
		log.Infof("Loading language %s...", langname)
		// Read the file
		langdata, err := ioutil.ReadFile(filename)
		// Check if there was an error
		if err != nil {
			// Error, print message.
			log.Fatalf("Failed to load the language %[1]s: %[2]s", langname, err)
			panic(err)
		}
		// Split the file string to an array of lines
		langraw := strings.Split(string(langdata), "\n")
		// Parse the data and save it to the language map.
		languages = append(languages, &Language{Name: langname, Data: parseLangData(langraw)})
		log.Debugf("Successfully loaded the language %s", langname)
	}
}

func parseLangData(langraw []string) map[string]string {
	lang := make(map[string]string)
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
	return lang
}

func trimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

// Translatef translates the given key and then formats the translated text with the given arguments.
func Translatef(key string, args ...interface{}) string {
	return fmt.Sprintf(Translate(key), args...)
}

// Translate translates the given key.
func Translate(key string) string {
	return GetLanguage("english").Translate(key)
}

// GetLanguage returns a language by the given name.
func GetLanguage(lang string) *Language {
	lang = strings.ToLower(lang)
	for _, l := range languages {
		if lang == l.Name {
			return l
		}
	}
	return nil
}

// Translatef translates the given key and then formats the translated text with the given arguments.
func (lng Language) Translatef(key string, args ...interface{}) string {
	return fmt.Sprintf(lng.Translate(key), args...)
}

// Translate translates the given key.
func (lng Language) Translate(key string) string {
	value, exists := lng.Data[key]
	if exists {
		return strings.Replace(value, "<br>", "\n", -1)
	}
	return key
}
