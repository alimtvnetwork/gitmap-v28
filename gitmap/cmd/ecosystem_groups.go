// Package cmd — ecosystem_groups.go manages group persistence for AGY, VSCode, Chrome, and GitHub Desktop.
package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// EcosystemGroup represents a named group of targets (projects, workspaces, profiles).
type EcosystemGroup struct {
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Targets     []string `json:"targets"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
}

// EcosystemGroupStore holds groups indexed by name for a specific ecosystem.
type EcosystemGroupStore struct {
	Ecosystem string                    `json:"ecosystem"`
	Groups    map[string]EcosystemGroup `json:"groups"`
}

func getEcosystemGroupStorePath(ecosystem string) (string, *apperror.AppError) {
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return "", apperror.WrapSimple(homeErr, "user home dir")
	}

	dirPath := filepath.Join(homeDir, ".gemini", "config", "groups")
	if mkErr := os.MkdirAll(dirPath, constants.DirPermission); mkErr != nil {
		return "", apperror.WrapSimple(mkErr, "mkdir groups")
	}

	return filepath.Join(dirPath, ecosystem+"_groups.json"), nil
}

func loadEcosystemGroupStore(ecosystem string) (*EcosystemGroupStore, *apperror.AppError) {
	filePath, pathErr := getEcosystemGroupStorePath(ecosystem)
	if pathErr != nil {
		return nil, pathErr
	}

	store := &EcosystemGroupStore{
		Ecosystem: ecosystem,
		Groups:    make(map[string]EcosystemGroup),
	}

	data, readErr := os.ReadFile(filePath)
	if readErr != nil {
		return store, nil
	}

	if unmarshalErr := json.Unmarshal(data, store); unmarshalErr != nil {
		return store, nil
	}

	if store.Groups == nil {
		store.Groups = make(map[string]EcosystemGroup)
	}

	return store, nil
}

func saveEcosystemGroupStore(store *EcosystemGroupStore) *apperror.AppError {
	filePath, pathErr := getEcosystemGroupStorePath(store.Ecosystem)
	if pathErr != nil {
		return pathErr
	}

	data, marshalErr := json.MarshalIndent(store, "", "  ")
	if marshalErr != nil {
		return apperror.WrapSimple(marshalErr, "marshal group store")
	}

	if writeErr := os.WriteFile(filePath, data, constants.FilePermission); writeErr != nil {
		return apperror.WrapSimple(writeErr, "write group store")
	}

	return nil
}

func newEcosystemGroup(name string, desc string, targets []string) EcosystemGroup {
	now := time.Now().UTC().Format(time.RFC3339)

	return EcosystemGroup{
		Name:        name,
		Description: desc,
		Targets:     targets,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func sliceContains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}

	return false
}
