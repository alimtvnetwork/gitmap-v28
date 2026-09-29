package cmdfoldertree

import (
	"fmt"
	"strings"
)

func resolveTargetPath(remaining []string, targetDir string) string {
	if len(remaining) > 0 && remaining[0] != "" {
		return remaining[0]
	}
	if targetDir != "" {
		return targetDir
	}
	return "."
}

func printImportSummary(s *ImportSummary, dryRun bool) {
	prefix := "✔ Created"
	if dryRun {
		prefix = "✔ [DRY RUN] Would create"
	}
	fmt.Printf("%s %d directories and %d files in %s\n",
		prefix, s.DirsCreated, s.FilesCreated, s.TargetDir)
}

func isFolderArg(s string) bool {
	low := strings.ToLower(s)
	return low == "folder" || low == "folders" || low == "dir" || low == "dirs"
}
