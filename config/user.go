package config

import (
	"strconv"
	"strings"
)

// User struct that is used for whitelisted user entries
type User struct {
	UID         int               `json:"uid"`
	Name        string            `json:"name"`
	Year        int               `json:"year"`
	Permissions []string          `json:"permissions"`
	Settings    map[string]string `json:"settings"`
}

// NilUser is an empty user type.
var NilUser = User{}

// CreateUser creates an user.
func CreateUser(uid int, name string, year int) User {
	return User{UID: uid, Name: name, Year: year, Permissions: make([]string, 0), Settings: make(map[string]string, 0)}
}

// AddUser adds the given user to the whitelist.
func AddUser(user User) bool {
	if GetUserWithUID(user.UID).UID != NilUser.UID || GetUserWithName(user.Name).UID != NilUser.UID {
		return false
	}
	config.Whitelist = append(config.Whitelist, user)
	return true
}

// RemoveUser removes the given user from the whitelist.
func RemoveUser(identifier string) bool {
	for index, user := range config.Whitelist {
		if strings.EqualFold(user.Name, identifier) || strings.EqualFold(strconv.Itoa(user.UID), identifier) {
			config.Whitelist = append(config.Whitelist[:index], config.Whitelist[index+1:]...)
			return true
		}
	}
	return false
}

// GetAllUsers returns all the whitelisted users.
func GetAllUsers() []User {
	return config.Whitelist
}

// GetUserWithUID gets the User struct that has the given UID.
func GetUserWithUID(uid int) User {
	for _, user := range config.Whitelist {
		if user.UID == uid {
			return user
		}
	}
	return NilUser
}

// GetUserWithName gets the User struct that has the given name.
func GetUserWithName(name string) User {
	for _, user := range config.Whitelist {
		if strings.EqualFold(user.Name, name) {
			return user
		}
	}
	return NilUser
}

// GetUser gets the User struct that has the given value as name or UID.
func GetUser(identifier string) User {
	for _, user := range config.Whitelist {
		if strings.EqualFold(user.Name, identifier) || strings.EqualFold(strconv.Itoa(user.UID), identifier) {
			return user
		}
	}
	return NilUser
}

// GetUsersWithSetting get all the users that have the given setting.
func GetUsersWithSetting(setting string, value string) []User {
	setting = strings.ToLower(setting)
	var users []User
	for _, user := range config.Whitelist {
		val, ok := user.GetSetting(setting)
		if ok && (len(value) == 0 || strings.EqualFold(value, val)) {
			users = append(users, user)
		}
	}
	return users
}

// GetSetting gets the given setting from the user.
func (u User) GetSetting(key string) (string, bool) {
	key = strings.ToLower(key)
	val, ok := u.Settings[key]
	return val, ok
}

// HasSetting checks if the user has the given setting.
func (u User) HasSetting(key string) bool {
	key = strings.ToLower(key)
	_, ok := u.GetSetting(key)
	return ok
}

// SetSetting sets a setting
func (u User) SetSetting(key string, value string) {
	key = strings.ToLower(key)
	u.Settings[key] = value
}

// RemoveSetting removes a setting
func (u User) RemoveSetting(key string) {
	key = strings.ToLower(key)
	delete(u.Settings, key)
}

// GetLanguage gets the user display language
func (u User) GetLanguage() string {
	lang, ok := u.GetSetting("language")
	if !ok {
		return "en_US"
	}
	return lang
}

// Destination returns the UID for Telebot.
func (u User) Destination() int {
	return u.UID
}

// HasPermission checks if the user has the given permission.
func (u User) HasPermission(permission string) bool {
	permission = strings.ToLower(permission)
	for _, perm := range u.Permissions {
		if strings.EqualFold(perm, permission) || strings.EqualFold(perm, "all") {
			return true
		}
	}
	return false
}

// AddPermission adds the given permission to the user.
func (u User) AddPermission(permission string) bool {
	permission = strings.ToLower(permission)
	if !u.HasPermission(permission) {
		u.Permissions = append(u.Permissions, permission)
		return true
	}
	return false
}

// RemovePermission removes the given permission from the user.
func (u User) RemovePermission(permission string) bool {
	permission = strings.ToLower(permission)
	for i, perm := range u.Permissions {
		if strings.EqualFold(perm, permission) {
			u.Permissions = append(u.Permissions[:i], u.Permissions[i+1:]...)
			return true
		}
	}
	return false
}
