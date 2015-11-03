package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

// User struct that is used for whitelisted user entries
type User struct {
	UID       int
	Name      string
	Yeargroup int
}

var whitelist []User

func isWhitelisted(uid int) bool {
	for _, e := range whitelist {
		if e.UID == uid {
			return true
		}
	}
	return false
}

// Load the whitelist from file
func loadWhitelist() {
	// Read the file
	wldata, err := ioutil.ReadFile("whitelist.txt")
	// Check if there was an error
	if err != nil {
		// Error, print message and use hardcoded whitelist.
		log.Printf(translate("whitelist.load.failed"), err)
		whitelist = []User{
			User{84359547, "Tulir", 21},
			User{67147746, "Ege", 21},
			User{128602828, "Max", 21},
			User{124500539, "Galax", 21},
			User{54580303, "Antti", 21},
			User{115187137, "Å", 21},
		}
	}
	// No error, parse the data
	log.Printf(translate("whitelist.loading"))
	// Split the file string to an array of lines
	wlraw := strings.Split(string(wldata), "\n")
	// Make the whitelist array
	whitelist = make([]User, len(wlraw), cap(wlraw))
	// Loop through the lines from the file
	for i := 0; i < len(wlraw); i++ {
		// Make sure the line is not empty
		if len(wlraw[i]) == 0 || strings.HasPrefix(wlraw[i], "#") {
			continue
		}
		// Split the entry to UID and name
		entry := strings.Split(wlraw[i], "|")
		// Convert the UID string to an integer
		uid, converr1 := strconv.Atoi(entry[0])
		// Convert the yeargroup index string to an integer
		ygindex, converr2 := strconv.Atoi(entry[2])
		// Make sure the conversion didn't fail
		if converr1 == nil && converr2 == nil {
			// No errors, add the UID to the whitelist
			whitelist[i] = User{uid, entry[1], ygindex}
			log.Printf(translate("whitelist.add.success"), whitelist[i].Name, whitelist[i].Yeargroup)
		} else {
			// Error occured, print message
			log.Printf(translate("whitelist.add.failed"), wlraw[i], err)
		}
	}
	log.Printf(translate("whitelist.load.success"))
}
