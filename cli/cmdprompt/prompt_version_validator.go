// Package cmdprompt — prompt_version_validator.go validates metadata sanity.
package cmdprompt

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func IsPromptArchitectInstalled(meta model.PromptArchitectMetadata) bool {
	return meta.Status == "active" || meta.Status == "installed" || len(meta.Version) > 0
}
