// Package cmdchromeprofile provides Chrome profile management subcommands.
package cmdchromeprofile

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/helpdoc"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func checkHelp(command string, args []string) {
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			helpdoc.Print(command)
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

func isHelpFlag(arg string) bool {
	return slices.Contains([]string{"help", "--help", "-h"}, arg)
}

func resolveTempDir() string {
	db, err := store.OpenDefault()
	if err != nil {
		return filepath.Join(".", ".ai-memory", "temp")
	}

	defer db.Close()
	val := db.GetSetting("temp_dir")
	if len(val) > 0 {
		return val
	}

	return filepath.Join(".", ".ai-memory", "temp")
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
		return apperror.WrapSimple(err, "failed writing to "+targetPath)
	}

	fmt.Printf("  %s✓%s Output written to %s\n", constants.ColorGreen, constants.ColorReset, targetPath)

	return nil
}
