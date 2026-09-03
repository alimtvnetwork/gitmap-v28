// Package installer — ui_ls_test.go tests table formatting.
package installer

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func TestFormatTable(t *testing.T) {
	empty := FormatInstallerTable(nil)
	if empty != "No installer scripts found." {
		t.Errorf("unexpected empty table: %s", empty)
	}

	items := []model.InstallerScript{
		{Name: "Tool A", Slug: "tool-a", TargetOS: "all", Version: "v1.0.0", Description: "Desc A"},
	}

	table := FormatInstallerTable(items)
	if !strings.Contains(table, "Tool A") || !strings.Contains(table, "tool-a") {
		t.Errorf("table missing items: %s", table)
	}
}
