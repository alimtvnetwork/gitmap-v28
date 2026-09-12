// Package cmdprompt — prompt_status_style.go applies colors to status strings.
package cmdprompt

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func formatPromptStatusText(meta model.PromptArchitectMetadata) string {
	if IsPromptArchitectInstalled(meta) {
		return constants.ColorGreen + "🟢 Active" + constants.ColorReset
	}

	return constants.ColorDim + "⚪ Not Installed" + constants.ColorReset
}
