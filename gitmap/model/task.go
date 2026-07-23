// Package model — task.go defines the file-sync task record.
package model

// TaskEntry represents a named file-sync watch task.
type TaskEntry struct {
	Name     string `json:"name"`
	Source   string `json:"source"`
	Dest     string `json:"dest"`
	Interval int    `json:"interval,omitempty"`
}

// TaskFile represents the top-level tasks.json structure.
type TaskFile struct {
	Tasks []TaskEntry `json:"tasks"`
}
