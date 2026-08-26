// Package cmd — prompt_version_writer.go writes Prompt Architect metadata to version.json.
package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// WritePromptArchitectMetadata injects or updates Prompt Architect metadata in version.json.
func WritePromptArchitectMetadata(repoPath string, meta model.PromptArchitectMetadata) error {
	vPath := filepath.Join(repoPath, "version.json")
	var root map[string]any

	data, errRead := os.ReadFile(vPath)
	if errRead == nil {
		_ = json.Unmarshal(data, &root)
	}
	if root == nil {
		root = make(map[string]any)
	}

	root[constants.PromptMetadataKey] = meta
	out, errMarshal := json.MarshalIndent(root, "", "  ")
	if errMarshal != nil {
		return errMarshal
	}

	return os.WriteFile(vPath, out, 0644)
}
