package fileutil

import (
	"path/filepath"
	"testing"

	"coding-guidelines/common/pkg/appfault"
)

type user struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestStreamJSON(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "massive.json")

	// Create a mock large JSON array
	users := []user{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 3, Name: "Charlie"},
	}
	ExportJSON(path, users, FilePermStandard)

	var count int
	res := StreamJSON(path, func(u user) *appfault.AppError {
		count++
		return nil
	})

	if res.IsFailure() {
		t.Fatalf("StreamJSON failed: %v", res.Fault().Error())
	}

	if count != 3 {
		t.Errorf("Expected to stream 3 elements, got %d", count)
	}
}
