package whitelist

import (
	"io/ioutil"
	log "maunium.net/ranssibot/maulog"
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

const whitelistFile = "data/whitelist"

var whitelist []User

// IsWhitelisted checks if the user with the given UID has been whitelisted.
func IsWhitelisted(uid int) bool {
	for _, e := range whitelist {
		if e.UID == uid {
			return true
		}
	}
	return false
}

// GetYeargroupIndex gets the yeargroup index of the user using the given UID.
func GetYeargroupIndex(uid int) int {
	for _, e := range whitelist {
		if e.UID == uid {
			return e.Year
		}
	}
	return 0
}

// GetPermissionLevel gets the permission level of the user with the given UID.
func GetPermissionLevel(uid int) int {
	for _, e := range whitelist {
		if e.UID == uid {
			return e.PermissionLevel
		}
	}
	return 0
}

// GetUID gets the UID matching the given name.
func GetUID(name string) int {
	for _, e := range whitelist {
		if strings.EqualFold(name, e.Name) {
			return e.UID
		}
	}
	return 0
}

// GetRecipientByUID gets a Recipient-type interface of the given UID.
func GetRecipientByUID(uid int) SimpleUser {
	return SimpleUser{uid}
}

// GetRecipientByName gets a Recipient-type interface of the given name.
func GetRecipientByName(name string) SimpleUser {
	return SimpleUser{GetUID(name)}
}

// SimpleUser ...
type SimpleUser struct {
	uid int
}

// Destination ...
func (user SimpleUser) Destination() int {
	return user.uid
}

func init() {
	// Read the file
	wldata, err := ioutil.ReadFile(whitelistFile)
	// Check if there was an error
	if err != nil {
		// Error, print message and use hardcoded whitelist.
		log.Errorf("Failed to load whitelist: %[1]s; Using hardcoded version", err)
		whitelist = []User{
			User{84359547, "Tulir", 21, 9001},
		}
	}
	// No error, parse the data
	log.Infof("Reading whitelist data...")
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
		if converr1 == nil && converr2 == nil && converr3 == nil {
			// No errors, add the UID to the whitelist
			whitelist[i] = User{uid, entry[1], year, perms}
			log.Infof("Added %[1]s (ID %[2]d) to the whitelist.", whitelist[i].Name, whitelist[i].UID)
		} else {
			// Error occured, print message
			log.Errorf("Failed to parse %[1]s: %[2]s", wlraw[i], err)
		}
	}
	log.Infof("Successfully loaded whitelist from file!")
}
