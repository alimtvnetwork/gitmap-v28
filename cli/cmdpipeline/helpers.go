package cmdpipeline

import (
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddb"
	"github.com/alimtvnetwork/gitmap-v28/cli/helptext"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/atotto/clipboard"
)

type clipboardWriterFunc func(text string) error

var globalAIFlag bool

// SetGlobalAIFlag explicitly marks AI session tracking mode.
func SetGlobalAIFlag(hasAI bool) {
	globalAIFlag = hasAI
}

// IsAIFlagActive checks whether AI mode or clipboard suppression is requested.
func IsAIFlagActive() bool {
	if globalAIFlag {
		return true
	}
	if os.Getenv("GITMAP_AI_TRACKING") == "1" || os.Getenv("GITMAP_NO_CLIPBOARD") == "1" {
		return true
	}

	return false
}

// IsClipboardWriteAllowed determines whether clipboard writes are allowed.
func IsClipboardWriteAllowed(hasAI bool) bool {
	if hasAI || IsAIFlagActive() {
		return false
	}

	return true
}

func isClipboardWriteAllowed(hasAI bool) bool {
	return IsClipboardWriteAllowed(hasAI)
}

var writeClipboard clipboardWriterFunc = func(text string) error {
	if !isClipboardWriteAllowed(false) {
		return nil
	}

	return clipboard.WriteAll(text)
}

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
