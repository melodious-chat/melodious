package main

// Flag - Describes a flag in the database
type Flag struct {
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
