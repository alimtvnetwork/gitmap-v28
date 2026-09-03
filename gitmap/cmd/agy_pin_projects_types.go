// Package cmd — agy_pin_projects_types.go defines types and storage for pinned projects.
package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

type PinnedProject struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	Branch   string `json:"branch,omitempty"`
	PinnedAt string `json:"pinnedAt"`
}

type PinnedProjectsStore struct {
	Version   string          `json:"version"`
	UpdatedAt string          `json:"updatedAt"`
	Projects  []PinnedProject `json:"projects"`
}

func getPinnedProjectsStorePath() (string, error) {
	homeDir, homeErr := os.UserHomeDir()

	if homeErr != nil {
		return "", homeErr
	}

	return filepath.Join(homeDir, ".gemini", "config", "pinned_projects.json"), nil
}

func loadPinnedProjectsStore() (*PinnedProjectsStore, *apperror.AppError) {
	path, pathErr := getPinnedProjectsStorePath()

	if pathErr != nil {
		return nil, apperror.WrapSimple(pathErr, "loadPinnedProjectsStore.path")
	}

	data, readErr := os.ReadFile(path)

	if readErr != nil {
		return &PinnedProjectsStore{
			Version:  "1.0.0",
			Projects: make([]PinnedProject, 0),
		}, nil
	}

	var store PinnedProjectsStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, apperror.WrapSimple(err, "loadPinnedProjectsStore.unmarshal")
	}

	return &store, nil
}

func savePinnedProjectsStore(store *PinnedProjectsStore) *apperror.AppError {
	path, pathErr := getPinnedProjectsStorePath()

	if pathErr != nil {
		return apperror.WrapSimple(pathErr, "savePinnedProjectsStore.path")
	}

	_ = os.MkdirAll(filepath.Dir(path), 0755)
	data, marshalErr := json.MarshalIndent(store, "", "  ")

	if marshalErr != nil {
		return apperror.WrapSimple(marshalErr, "savePinnedProjectsStore.marshal")
	}

	if writeErr := os.WriteFile(path, data, 0644); writeErr != nil {
		return apperror.WrapSimple(writeErr, "savePinnedProjectsStore.write")
	}

	return nil
}
