// Package cmd — pull_table_style.go applies colors and icons for pull table.
package cmdpull

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func formatPullStatus(status string, isDirty bool) string {
	if isDirty {
		return constants.ColorYellow + "● dirty" + constants.ColorReset
	}

	st := constants.PullChangesStatusType(strings.ToLower(status))
	if st.IsActive() {
		return constants.ColorGreen + "✔ active" + constants.ColorReset
	}

	return formatPullStatusOther(st, status)
}

func formatPullStatusOther(st constants.PullChangesStatusType, raw string) string {
	if st.IsUpdated() {
		return constants.ColorGreen + "✔ updated" + constants.ColorReset
	}

	if st.IsFailed() {
		return constants.ColorRed + "✖ failed" + constants.ColorReset
	}

	if st == constants.PullChangesStatusClean {
		return constants.ColorGreen + "✔ clean" + constants.ColorReset
	}

	return raw
}
