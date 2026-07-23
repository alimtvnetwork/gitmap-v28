// Package model — profile.go defines the profile configuration structure.
package model

// ProfileConfig holds the list of profiles and the active one.
type ProfileConfig struct {
	Active   string   `json:"active"`
	Profiles []string `json:"profiles"`
}
