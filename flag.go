package main

// Flag - Describes a flag in the database
type Flag struct {
	Group    string
	Name     string
	Flag     map[string]interface{}
	Priority int
}
