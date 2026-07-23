// Package model — bookmark.go defines the bookmark record for saved commands.
package model

// BookmarkRecord represents a saved command+flags combination.
type BookmarkRecord struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Command   string `json:"command"`
	Args      string `json:"args,omitempty"`
	Flags     string `json:"flags,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
}
