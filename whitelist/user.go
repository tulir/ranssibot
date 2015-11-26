package whitelist

import (
	"strings"
)

// Whitelist is the container for a whitelisted group of users.
type Whitelist struct {
	Users []User `json:"users"`
}

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

var whitelist = &Whitelist{}

// GetUserWithUID gets the User struct that has the given UID.
func GetUserWithUID(uid int) User {
	for _, user := range whitelist.Users {
		if user.UID == uid {
			return user
		}
	}
	return NilUser
}

// GetUserWithName gets the User struct that has the given name.
func GetUserWithName(name string) User {
	for _, user := range whitelist.Users {
		if user.Name == name {
			return user
		}
	}
	return NilUser
}

// GetSetting gets the given setting from the user.
func (u User) GetSetting(key string) string {
	return u.Settings[key]
}

// SetSetting sets a setting
func (u User) SetSetting(key string, value string) {
	u.Settings[key] = value
}

// Destination returns the UID for Telebot.
func (u User) Destination() int {
	return u.UID
}

// HasPermission checks if the user has the given permission.
func (u User) HasPermission(permission string) bool {
	for _, perm := range u.Permissions {
		if strings.EqualFold(perm, permission) || strings.EqualFold(perm, "all") {
			return true
		}
	}
	return false
}

// AddPermission adds the given permission to the user.
func (u User) AddPermission(permission string) bool {
	if !u.HasPermission(permission) {
		u.Permissions = append(u.Permissions, permission)
		return true
	}
	return false
}

// RemovePermission removes the given permission from the user.
func (u User) RemovePermission(permission string) bool {
	for i, perm := range u.Permissions {
		if strings.EqualFold(perm, permission) {
			u.Permissions = append(u.Permissions[:i], u.Permissions[i+1:]...)
			return true
		}
	}
	return false
}
