// Package model — projecttype.go defines the ProjectType reference struct.
package model

// ProjectType represents a supported project type in the reference table.
type ProjectType struct {
	ID          int64  `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
