// Package cmd — prompt_status_style.go applies colors to status strings.
package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func formatPromptStatusText(meta model.PromptArchitectMetadata) string {
	if IsPromptArchitectInstalled(meta) {
		return constants.ColorGreen + "🟢 Active" + constants.ColorReset
	}
	return constants.ColorDim + "⚪ Not Installed" + constants.ColorReset
}
