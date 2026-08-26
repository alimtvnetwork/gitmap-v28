// Package installer — ui_ls.go formats the installer scripts list table.
package installer

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// FormatInstallerTable produces a formatted ASCII table string of the installer scripts.
func FormatInstallerTable(list []model.InstallerScript) string {
	if len(list) == 0 {
		return "No installer scripts found."
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("%-20s %-15s %-10s %-10s %s\n", "NAME", "SLUG", "OS", "VERSION", "DESCRIPTION"))
	b.WriteString(strings.Repeat("-", 75) + "\n")

	for _, item := range list {
		b.WriteString(fmt.Sprintf("%-20s %-15s %-10s %-10s %s\n",
			item.Name, item.Slug, item.TargetOS, item.Version, item.Description))
	}

	return b.String()
}
