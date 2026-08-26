// Package cmd — prompt_version_reader.go reads Prompt Architect metadata from version.json.
package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// ReadPromptArchitectMetadata extracts Prompt Architect metadata from a repo's version.json.
func ReadPromptArchitectMetadata(repoPath string) (model.PromptArchitectMetadata, error) {
	vPath := filepath.Join(repoPath, "version.json")
	data, errRead := os.ReadFile(vPath)
	if errRead != nil {
		return model.PromptArchitectMetadata{Status: "not_installed"}, errRead
	}

	var root map[string]any
	if errJSON := json.Unmarshal(data, &root); errJSON != nil {
		return model.PromptArchitectMetadata{Status: "corrupt"}, errJSON
	}

	var section any
	if val, ok := root[constants.PromptMetadataKey]; ok {
		section = val
	} else if valAlt, okAlt := root[constants.PromptMetadataAlt]; okAlt {
		section = valAlt
	} else {
		return model.PromptArchitectMetadata{Status: "not_installed"}, nil
	}

	secBytes, errMarshal := json.Marshal(section)
	if errMarshal != nil {
		return model.PromptArchitectMetadata{Status: "unknown"}, errMarshal
	}

	var meta model.PromptArchitectMetadata
	_ = json.Unmarshal(secBytes, &meta)
	if meta.Status == "" {
		meta.Status = "active"
	}
	return meta, nil
}
