// Package cmd — pull_table_summary.go renders final pull batch statistics.
package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
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
	fmt.Println("  --------------------------------------------------------------------------------")
}
