// Package cmdide — ide_scan_sync.go: post-scan IDE registration sync.
//
// Moved from cli/cmd/clihelpers.go during the spec-243 cmd split (Wave D):
// syncing scan records across IDEs is cmdide's concern.
package cmdide

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// SyncScanRecordsToIDEs registers scanned repos across the targeted IDEs
// (VS Code, Cursor, Antigravity, Desktop), honoring --skip-sync and the
// --sync-ide / --exclude-sync target filters.
func SyncScanRecordsToIDEs(records []model.ScanRecord, args []string, isSkipVSCodeSync, isQuiet bool) {
	opts := IDEOptions{
		IsVSCodeTargeted:      true,
		IsCursorTargeted:      true,
		IsAntigravityTargeted: true,
		IsDesktopTargeted:     true,
		IsQuiet:               isQuiet,
	}

	for i, arg := range args {
		if arg == "--skip-sync" || arg == "--no-ide-sync" || arg == "--no-sync" {
			return
		}
		if (arg == "--sync-ide" || arg == "--ide") && i+1 < len(args) {
			applyIDEScanTargetFilter(&opts, args[i+1])
		}
		if (arg == "--exclude-sync" || arg == "--skip-ide") && i+1 < len(args) {
			applyIDEScanExcludeFilter(&opts, args[i+1])
		}
	}

	if isSkipVSCodeSync {
		opts.IsVSCodeTargeted = false
	}

	repoPaths := make([]string, 0, len(records))
	for _, rec := range records {
		if rec.AbsolutePath != "" {
			repoPaths = append(repoPaths, rec.AbsolutePath)
		}
	}

	summary := SyncReposAcrossIDEsDirect(repoPaths, opts)
	if !isQuiet {
		fmt.Printf("  • IDE Registrations: %d repos checked across IDEs (VS Code: +%d, Cursor: +%d, Antigravity: +%d, Desktop: +%d)\n",
			summary.TotalRepos, summary.VSCodeAdded, summary.CursorAdded, summary.AntigravityAdded, summary.DesktopAdded)
	}
}

func applyIDEScanTargetFilter(opts *IDEOptions, target string) {
	low := strings.ToLower(target)
	if low != "all" {
		opts.IsVSCodeTargeted = strings.Contains(low, "vscode") || strings.Contains(low, "code")
		opts.IsCursorTargeted = strings.Contains(low, "cursor")
		opts.IsAntigravityTargeted = strings.Contains(low, "antigravity") || strings.Contains(low, "agy")
		opts.IsDesktopTargeted = strings.Contains(low, "desktop")
	}
}

func applyIDEScanExcludeFilter(opts *IDEOptions, excluded string) {
	low := strings.ToLower(excluded)
	opts.IsVSCodeTargeted = opts.IsVSCodeTargeted && !strings.Contains(low, "vscode") && !strings.Contains(low, "code")
	opts.IsCursorTargeted = opts.IsCursorTargeted && !strings.Contains(low, "cursor")
	opts.IsAntigravityTargeted = opts.IsAntigravityTargeted && !strings.Contains(low, "antigravity") && !strings.Contains(low, "agy")
	opts.IsDesktopTargeted = opts.IsDesktopTargeted && !strings.Contains(low, "desktop")
}
