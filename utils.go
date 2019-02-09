package main

// Flag - Describes a flag in the database
type Flag struct {
	ID       int
	HasID    bool
	Group    string
	Name     string
	Flag     map[string]interface{}
	Priority int
}

// GroupHolder - Describes a group holder in the database
type GroupHolder struct {
	Group string

	User    string
	Channel string
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
	ID       int
	Username string
	Owner    bool
}

// UserStatus - describes the status of a user
type UserStatus struct {
	User   *User
	Online bool
}
