// Package cmd — pull_table_summary.go renders final pull batch statistics.
package cmdpull

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func RenderPullBatchTable(rows []model.PullTableRow) {
	if len(rows) == 0 {
		return
	}

	layout := NewPullTableLayout(rows)
	fmt.Println()
	layout.PrintHeader()
	for _, r := range rows {
		layout.PrintRow(r)
	}

	fmt.Printf("  %s\n", strings.Repeat("-", layout.DividerLen))
}
