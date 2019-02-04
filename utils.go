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

// FlagQueryResult - Describes result of a flag query
type FlagQueryResult struct {
	All         map[*Flag]bool
	User        map[*Flag]bool
	Channel     map[*Flag]bool
	UserChannel map[*Flag]bool
}
