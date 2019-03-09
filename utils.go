package main

import (
	"regexp"
	"strconv"
	"strings"
)

// Flag - Describes a flag in the database
type Flag struct {
	ID    int                    `json:"id"`
	Group string                 `json:"group"`
	Name  string                 `json:"name"`
	Flag  map[string]interface{} `json:"flag"`
}

// Group - describes a group in the database
type Group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GroupHolder - Describes a group holder in the database
type GroupHolder struct {
	ID      int    `json:"id"`
	Group   string `json:"group"`
	User    string `json:"user"`
	Channel string `json:"channel"`
}

// GroupHolderFQResult - Describes a group holder received from FlagQueryResult
type GroupHolderFQResult struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	ChannelID int    `json:"channel_id"`
	GroupID   int    `json:"group_id"`
	GroupName string `json:"group_name"`
}

// FlagQueryResult - Describes result of a flag query
type FlagQueryResult struct {
	GroupHolders []*GroupHolderFQResult
	FlagID       int
	FlagName     string
	Flag         map[string]interface{}
}

// ChatMessage - a message received from message history
type ChatMessage struct {
	Message   string   `json:"content"`
	Pings     []string `json:"pings"`
	ID        int      `json:"id"`
	Timestamp string   `json:"timestamp"`
	Author    string   `json:"author"`
	AuthorID  int      `json:"author_id"`
}

// User - describes a user in the database
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Owner    bool   `json:"owner"`
}

// UserStatus - describes the status of a user
type UserStatus struct {
	User   *User `json:"user"`
	Online bool  `json:"online"`
}

// Channel - describes a channel in the database
type Channel struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Topic string `json:"topic"`
}

// getPings - gets all mentioned/pinged user IDs from a message string.
func scanForPings(message string) []int {
	re := regexp.MustCompile(`\<([^\<\>]*)\>`)
	ids := []int{}
	submatchall := re.FindAllString(message, -1)
	for _, element := range submatchall {
		element = strings.Trim(element, "<")
		element = strings.Trim(element, ">")
		element = strings.Trim(element, "@")
		i, err := strconv.ParseInt(element, 10, 32)
		if err == nil {
			ids = append(ids, int(i))
		}
	}
	return ids
}

// contains - check if a slice contains a value
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
