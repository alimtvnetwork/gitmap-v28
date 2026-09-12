// Package cmd — pull_table_style.go applies colors and icons for pull table.
package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func formatPullStatus(status string, isDirty bool) string {
	if isDirty {
		return constants.ColorYellow + "● dirty" + constants.ColorReset
	}

	if status == "active" || status == "ACTIVE" {
		return constants.ColorGreen + "✔ active" + constants.ColorReset
	}

	if status == "UP_TO_DATE" || status == "up-to-date" || status == "synced" {
		return constants.ColorGreen + "✔ active" + constants.ColorReset
	}

	return formatPullStatusOther(status)
}

func formatPullStatusOther(status string) string {
	switch status {
	case "SUCCESS", "ok", "updated":
		return constants.ColorGreen + "✔ updated" + constants.ColorReset
	case "FAILED", "fail", "error":
		return constants.ColorRed + "✖ failed" + constants.ColorReset
	case "clean":
		return constants.ColorGreen + "✔ clean" + constants.ColorReset
	default:
		return status
	}
}
