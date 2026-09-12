// Package cmdchromeprofile provides Chrome profile management subcommands.
package cmdchromeprofile

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/helptext"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func checkHelp(command string, args []string) {
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			helptext.Print(command)
			cliexit.Exit(0)
		}
	}
}

func hasJSONFlag(args []string) bool {
	for _, a := range args {
		if a == "--json" {
			return true
		}
	}

	return false
}

func isDirectoryPath(path string) bool {
	info, err := os.Stat(path)

	return err == nil && info.IsDir()
}

func isHelpFlag(arg string) bool {
	return arg == "help" || arg == "--help" || arg == "-h"
}

func resolveTempDir() string {
	db, err := store.OpenDefault()
	if err != nil {
		return filepath.Join(".", ".lovable", "temp")
	}

	defer db.Close()
	val := db.GetSetting("temp_dir")
	if len(val) > 0 {
		return val
	}

	return filepath.Join(".", ".lovable", "temp")
}

func truncateStr(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}

	return s
}

func writeContentToFile(targetPath, content string) error {
	dir := filepath.Dir(targetPath)
	if len(dir) > 0 {
		_ = os.MkdirAll(dir, 0755)
	}

	err := os.WriteFile(targetPath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed writing to %s: %w", targetPath, err)
	}

	fmt.Printf("  %s✓%s Output written to %s\n", constants.ColorGreen, constants.ColorReset, targetPath)

	return nil
}
