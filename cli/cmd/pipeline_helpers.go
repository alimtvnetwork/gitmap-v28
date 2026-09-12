package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func resolveTempDir() string {
	db, err := openDB()
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

func hasArgFlag(args []string, flagName string) bool {
	for _, a := range args {
		if a == flagName || strings.HasPrefix(a, flagName+"=") {
			return true
		}
	}

	return false
}

func extractFlagVal(args []string, flagName string) string {
	for i, arg := range args {
		if arg == flagName && i+1 < len(args) {
			return args[i+1]
		}

		if strings.HasPrefix(arg, flagName+"=") {
			return strings.TrimPrefix(arg, flagName+"=")
		}
	}

	return ""
}

func printJSON(v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}
