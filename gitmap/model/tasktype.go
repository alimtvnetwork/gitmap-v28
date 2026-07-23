// Package model — tasktype.go defines the task type lookup record.
package model

// TaskTypeRecord represents a task category (Delete, Remove).
type TaskTypeRecord struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
