package cmddb

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/helptext"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// DbFileInfo holds descriptive metadata for any discovered SQLite database.
type DbFileInfo struct {
	Name     string
	Path     string
	Size     int64
	Category string
	Purpose  string
	RepoID   int64
	RepoSlug string
}

// DBFileInfo is a backwards-compatible alias for DbFileInfo.
type DBFileInfo = DbFileInfo

// parseConfirmFlag parses the --confirm flag for a command, returning the boolean value directly without pointer dereference.
func parseConfirmFlag(cmdName string, args []string) bool {
	fs := flag.NewFlagSet(cmdName, flag.ExitOnError)
	var isConfirm bool
	fs.BoolVar(&isConfirm, constants.FlagConfirm, false, constants.FlagDescConfirm)
	_ = fs.Parse(args)

	return isConfirm || hasConfirmFlag(args)
}

func findSplitDbDirs() []string {
	binDataDir := store.BinaryDataDir()
	binRoot := filepath.Dir(binDataDir)
	raw := []string{
		filepath.Join(binDataDir, "repo_search"),
		filepath.Join(binRoot, constants.DefaultOutputDir, "repo_search"),
		filepath.Join(".", "data", "repo_search"),
		filepath.Join(".", constants.DefaultOutputDir, "repo_search"),
	}

	return dedupeDirs(raw)
}

// findSplitDBDirs is a backwards-compatible alias for findSplitDbDirs.
var findSplitDBDirs = findSplitDbDirs

func dedupeDirs(raw []string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, d := range raw {
		clean := filepath.Clean(d)
		if seen[clean] {
			continue
		}

		seen[clean] = true
		if isExistingDir(clean) {
			out = append(out, clean)
		}
	}

	return out
}

func isExistingDir(path string) bool {
	info, err := os.Stat(path)

	return err == nil && info.IsDir()
}

func formatBytes(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}

	if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024.0)
	}

	if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(bytes)/(1024.0*1024.0))
	}

	return fmt.Sprintf("%.2f GB", float64(bytes)/(1024.0*1024.0*1024.0))
}

func promptConfirm(msg string) (bool, error) {
	fmt.Print(msg)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	ans := strings.ToLower(strings.TrimSpace(line))
	hasConfirmed := ans == "y" || ans == "yes"

	return hasConfirmed, nil
}

func confirmOrSkip(msg string, args []string) bool {
	if hasConfirmFlag(args) {
		return true
	}

	if !isInteractiveStdin() {
		return false
	}

	hasConfirmed, err := promptConfirm(msg)

	return err == nil && hasConfirmed
}

func isInteractiveStdin() bool {
	if os.Getenv("CI") != "" || os.Getenv("GITMAP_NON_INTERACTIVE") == "1" {
		return false
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return (fi.Mode() & os.ModeCharDevice) != 0
}

func hasConfirmFlag(args []string) bool {
	for _, a := range args {
		trimmed := strings.TrimSpace(a)
		if trimmed == "-y" || trimmed == "--yes" || trimmed == "-f" || trimmed == "--force" || trimmed == "--confirm" {
			return true
		}
	}

	return false
}

func collectMainDbInfo() (DbFileInfo, bool) {
	mainPath := store.DefaultDBPath()
	info, err := os.Stat(mainPath)
	if err != nil {
		return DbFileInfo{
			Name:     filepath.Base(mainPath),
			Path:     mainPath,
			Category: "Primary Master DB",
			Purpose:  "Central SQLite database storing global tracked repositories, scan history, configurations, and profiles.",
		}, false
	}

	return DbFileInfo{
		Name:     filepath.Base(mainPath),
		Path:     mainPath,
		Size:     info.Size(),
		Category: "Primary Master DB",
		Purpose:  "Central SQLite database storing global tracked repositories, scan history, configurations, and profiles.",
	}, true
}

// collectMainDBInfo is a backwards-compatible alias for collectMainDbInfo.
var collectMainDBInfo = collectMainDbInfo

func hasArgFlag(args []string, flagName string) bool {
	for _, a := range args {
		if a == flagName || strings.HasPrefix(a, flagName+"=") {
			return true
		}
	}

	return false
}

func printJSON(v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}

func checkHelp(command string, args []string) {
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			helptext.Print(command)
			cliexit.Exit(0)
		}
	}
}
