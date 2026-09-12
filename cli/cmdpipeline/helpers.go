package cmdpipeline

import (
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddb"
	"github.com/alimtvnetwork/gitmap-v28/cli/helptext"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func checkHelp(command string, args []string) {
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			helptext.Print(command)
			cliexit.Exit(0)
		}
	}
}

func openDB() (*store.DB, error) {
	return store.OpenDefault()
}

func findRepoRoot(path string) string {
	current := path
	for {
		if _, err := os.Stat(filepath.Join(current, "version.json")); err == nil {
			return current
		}
		parent := filepath.Dir(current)
		if parent == current {
			return ""
		}
		current = parent
	}
}

var formatBytes = cmddb.FormatBytes
var confirmOrSkip = cmddb.ConfirmOrSkip
var runDB = cmddb.RunDB
