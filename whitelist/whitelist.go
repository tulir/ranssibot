package whitelist

import (
	"encoding/json"
	"io/ioutil"
	log "maunium.net/maulogger"
)

const whitelistFile = "data/whitelist.json"

// Load loads the whitelist from file.
func Load() {
	// Read the file
	data, err := ioutil.ReadFile(whitelistFile)
	// Check if there was an error
	if err != nil {
		loadFailed(err)
		return
	}
	// No error, parse the data
	log.Infof("Reading whitelist data...")
	err = json.Unmarshal(data, whitelist)
	// Check if parsing failed
	if err != nil {
		loadFailed(err)
		return
	}
	log.Debugf("Successfully loaded whitelist from file!")
}

// Save saves the whitelist data to file.
func Save(indent bool) {
	var data []byte
	var err error
	if indent {
		data, err = json.MarshalIndent(whitelist, "", "    ")
	} else {
		data, err = json.Marshal(whitelist)
	}
	if err != nil {
		log.Errorf("Failed to save whitelist: %[1]s", err)
		return
	}
	err = ioutil.WriteFile(whitelistFile, data, 0700)
	if err != nil {
		log.Errorf("Failed to save whitelist: %[1]s", err)
		return
	}
}

func loadFailed(err error) {
	log.Errorf("Failed to load whitelist: %[1]s; Using hardcoded version", err)
	*whitelist = Whitelist{
		Users: []User{
			User{UID: 84359547, Name: "Tulir", Year: 21, Permissions: []string{"all"}},
		},
	}
}
