package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

// User struct that is used for whitelisted user entries
type User struct {
	UID             int
	Name            string
	Year            int
	PermissionLevel int
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

func getYeargroupIndex(uid int) int {
	for _, e := range whitelist {
		if e.UID == uid {
			return e.Year
		}
	}
	return 0
}

func getPermissionLevel(uid int) int {
	for _, e := range whitelist {
		if e.UID == uid {
			return e.PermissionLevel
		}
	}
	return 0
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
			User{84359547, "Tulir", 21, 9001},
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
		// Convert the year string to an integer
		year, converr2 := strconv.Atoi(entry[2])
		// Convert the permission level string to an integer
		perms, converr3 := strconv.Atoi(entry[3])
		// Make sure the conversion didn't fail
		if converr1 == nil && converr2 == nil && converr3 != nil {
			// No errors, add the UID to the whitelist
			whitelist[i] = User{uid, entry[1], year, perms}
			log.Printf(translate("whitelist.add.success"), whitelist[i].Name, whitelist[i].UID)
		} else {
			// Error occured, print message
			log.Printf(translate("whitelist.add.failed"), wlraw[i], err)
		}
	}
	log.Printf(translate("whitelist.load.success"))
}
