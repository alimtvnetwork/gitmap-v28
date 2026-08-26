// Package cmd — prompt_status_row.go formats individual prompt status rows.
package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func (l *PromptStatusTableLayout) PrintRow(repoPath string, meta model.PromptArchitectMetadata) {
	name := filepath.Base(repoPath)
	statusText := formatPromptStatusText(meta)
	ver := meta.Version
	if ver == "" {
		ver = "-"
	}
	date := meta.InstalledAt
	if date == "" {
		date = "-"
	}

	fmt.Printf("  %-*s   %-*s   %-*s   %s\n",
		l.MaxRepo, name,
		l.MaxStatus, statusText,
		l.MaxVersion, ver,
		date,
	)
}
