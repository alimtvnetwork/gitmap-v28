// Package cmd — prompt_version_validator.go validates metadata sanity.
package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func IsPromptArchitectInstalled(meta model.PromptArchitectMetadata) bool {
	return meta.Status == "active" || meta.Status == "installed" || len(meta.Version) > 0
}
