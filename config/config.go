package config

import (
	"encoding/json"
	flag "github.com/ogier/pflag"
	"io/ioutil"
	log "maunium.net/go/maulogger"
)

// Configuration is the container for a config.
type Configuration struct {
	Whitelist    []User `json:"whitelist"`
	LastReadPost int    `json:"last-read-post"`
}

// IndentConfig determines if the config should be pretty-printed
var IndentConfig = false

// Autosave determines if the config should be automatically saved when something is changed.
var Autosave = flag.BoolP("config-autosave", "s", true, "Don't save config when something is changed.")

// ConfigFile is the configuration file to use.
const ConfigFile = "config.json"

var config = &Configuration{}

// Load loads the whitelist from file.
func Load() {
	config = &Configuration{}
	// Read the file
	data, err := ioutil.ReadFile(ConfigFile)
	// Check if there was an error
	if err != nil {
		loadFailed(err)
		return
	}
	// No error, parse the data
	log.Infof("[Config] Reading data...")
	err = json.Unmarshal(data, config)
	// Check if parsing failed
	if err != nil {
		loadFailed(err)
		return
	}
	log.Debugf("[Config] Successfully loaded data from disk")
}

// Save saves the whitelist data to file.
func Save() {
	log.Infof("[Config] Saving data to disk...")
	save()
}

func save() {
	var data []byte
	var err error
	if IndentConfig {
		data, err = json.MarshalIndent(config, "", "    ")
	} else {
		data, err = json.Marshal(config)
	}
	if err != nil {
		log.Errorf("[Config] Failed to save: %[1]s", err)
		return
	}
	err = ioutil.WriteFile(ConfigFile, data, 0600)
	if err != nil {
		log.Errorf("[Config] Failed to save: %[1]s", err)
		return
	}
	log.Debugf("[Config] Successfully saved to disk")
}

// ASave calls Save if Autosave is true
func ASave() {
	if *Autosave {
		log.Debugf("[Config] Autosaving...")
		save()
	}
}

// GetConfig gets the configuration
func GetConfig() *Configuration {
	return config
}

func loadFailed(err error) {
	log.Errorf("[Config] Failed to load: %[1]s; Using hardcoded version", err)
	*config = Configuration{
		Whitelist: []User{
			User{UID: 84359547, Name: "Tulir", Year: 1, Permissions: []string{"all"}},
		},
	}
}
