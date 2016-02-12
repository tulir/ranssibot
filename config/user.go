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
	// Make sure that neither the UID nor the name of the new user is used.
	if GetUserWithUID(user.UID).UID != NilUser.UID || GetUserWithName(user.Name).UID != NilUser.UID {
		return false
	}
	// Append the user to the whitelist.
	config.Whitelist = append(config.Whitelist, user)
	return true
}

// RemoveUser removes the given user from the whitelist.
func RemoveUser(identifier string) bool {
	// Loop through the whitelist
	for index, user := range config.Whitelist {
		// Check if the current user in the loop has the given identifier (name or UID).
		if strings.EqualFold(user.Name, identifier) || strings.EqualFold(strconv.Itoa(user.UID), identifier) {
			if len(config.Whitelist) == 1 {
				// If the user to be removed is the only user in the whitelist, simply delete the whole list.
				config.Whitelist = []User{}
			} else if len(config.Whitelist)-1 == index || index == 0 {
				// If the user to be removed is the first or last user in the whitelist, take a slice of the list
				// leaving out the first or last index correspondingly.
				config.Whitelist = config.Whitelist[:index]
			} else {
				// In all other cases (as in there are users on both sides of the user to be removed), take a slice
				// from the first user to the previous user, another slice from the next user to the last user and
				// append the two slices.
				config.Whitelist = append(config.Whitelist[:index], config.Whitelist[index+1:]...)
			}
			// Match found and user removed from whitelist.
			return true
		}
	}
	// No match found, return false.
	return false
}

// GetAllUsers returns all the whitelisted users.
func GetAllUsers() []User {
	return config.Whitelist
}

// GetUserWithUID gets the User struct that has the given UID.
func GetUserWithUID(uid int) User {
	// Loop through the whitelist
	for _, user := range config.Whitelist {
		// If the UID of the current user in the loop matches the given UID, return the user.
		if user.UID == uid {
			return user
		}
	}
	// No match found, return the empty user.
	return NilUser
}

// GetUserWithName gets the User struct that has the given name.
func GetUserWithName(name string) User {
	// Loop through the whitelist
	for _, user := range config.Whitelist {
		// If the name of the current user in the loop matches the given name, return the user.
		if strings.EqualFold(user.Name, name) {
			return user
		}
	}
	// No match found, return the empty user.
	return NilUser
}

// GetUser gets the User struct that has the given value as name or UID.
func GetUser(identifier string) User {
	// Loop through the whitelist
	for _, user := range config.Whitelist {
		// If the current user in the loop has the given identifier (name or UID), return the user.
		if strings.EqualFold(user.Name, identifier) || strings.EqualFold(strconv.Itoa(user.UID), identifier) {
			return user
		}
	}
	// No match found, return the empty user.
	return NilUser
}

// GetUsersWithSetting get all the users that have the given setting.
func GetUsersWithSetting(setting string, values ...string) []User {
	// Make the setting key lowercase
	setting = strings.ToLower(setting)
	var users []User
	// Loop through the whitelist
	for _, user := range config.Whitelist {
		// Check if the current user in the loop has the given setting.
		val, ok := user.GetSetting(setting)
		// If a required value was passed, make sure the value the user has matches.
		if ok {
			if len(values) > 0 {
				for _, valc := range values {
					if strings.EqualFold(valc, val) {
						users = append(users, user)
					}
				}
			} else {
				users = append(users, user)
			}
		}
	}
	// Return the list of accepted users.
	return users
}

// GetSetting gets the given setting from the user.
func (u User) GetSetting(key string) (string, bool) {
	// Make the setting key lowercase
	key = strings.ToLower(key)
	val, ok := u.Settings[key]
	return val, ok
}

// HasSetting checks if the user has the given setting.
func (u User) HasSetting(key string) bool {
	_, ok := u.GetSetting(key)
	return ok
}

// SetSetting sets a setting
func (u User) SetSetting(key string, value string) {
	// Make the setting key lowercase
	key = strings.ToLower(key)
	u.Settings[key] = value
}

// RemoveSetting removes a setting
func (u User) RemoveSetting(key string) {
	// Make the setting key lowercase
	key = strings.ToLower(key)
	delete(u.Settings, key)
}

// GetLanguage gets the user display language
func (u User) GetLanguage() string {
	// Get the language setting of the user.
	lng, ok := u.GetSetting("language")
	// If none set, return "english"
	if !ok {
		return "english"
	}
	// Return the language of the user.
	return lng
}

// Destination returns the UID for Telebot.
func (u User) Destination() string {
	return strconv.Itoa(u.UID)
}

// HasPermission checks if the user has the given permission.
func (u User) HasPermission(permission string) bool {
	permission = strings.ToLower(permission)
	minus := strings.HasPrefix(permission, "-")
	// Loop through the permissions of the user.
	for _, perm := range u.Permissions {
		if !minus && strings.EqualFold(perm, "all") {
			// If the requested permission is NOT a negative one and the user has the "all" permission, return true.
			return true
		} else if strings.EqualFold(perm, "-all") {
			if minus {
				// If the requested permission is a negative one and the user has the negative "all" permission, return true.
				return true
			}
			// If the requested permission is NOT a negative one and the user has the negative "all" permission, return false.
			return false
		} else if strings.EqualFold(perm, permission) {
			// If the user has a permission that equals the given permission, return true.
			return true
		}
	}
	// The user does not have the permission, return false.
	return false
}

// AddPermission adds the given permission to the user.
func (u User) AddPermission(permission string) bool {
	permission = strings.ToLower(permission)
	// Make sure that the user doesn't have the permission.
	if !u.HasPermission(permission) {
		// Append the permission to the users permissions.
		u.Permissions = append(u.Permissions, permission)
		return true
	}
	return false
}

// RemovePermission removes the given permission from the user.
func (u User) RemovePermission(permission string) bool {
	permission = strings.ToLower(permission)
	// Loop through the permissions of the user.
	for index, perm := range u.Permissions {
		if strings.EqualFold(perm, permission) {
			if len(u.Permissions) == 1 {
				// If the user has no other permissions, simply delete the whole list.
				u.Permissions = []string{}
			} else if len(u.Permissions)-1 == index || index == 0 {
				// If the permission to be removed is the first or last permission in the list, take a slice of the list
				// leaving out the first or last index correspondingly.
				u.Permissions = u.Permissions[:index]
			} else {
				// In all other cases (as in there are permissions on both sides of the permission to be removed), take
				// a slice from the first permission to the previous permission, another slice from the next permission
				// to the last permission and append the two slices.
				u.Permissions = append(u.Permissions[:index], u.Permissions[index+1:]...)
			}
			return true
		}
	}
	return false
}
