// Package cmd — pull_table_style.go applies colors and icons for pull table.
package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

//nolint:unused
func formatPullStatus(status string, isDirty bool) string {
	if isDirty {
		return constants.ColorYellow + "● dirty" + constants.ColorReset
	}
	switch status {
	case "SUCCESS", "ok", "updated":
		return constants.ColorGreen + "✓ updated" + constants.ColorReset
	case "UP_TO_DATE", "up-to-date":
		return constants.ColorDim + "— up-to-date" + constants.ColorReset
	case "FAILED", "fail", "error":
		return constants.ColorRed + "✖ failed" + constants.ColorReset
	default:
		return status
	}
}
