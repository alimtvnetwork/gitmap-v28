package cmdinject

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/flagutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func openDB() (*store.DB, error) {
	return store.OpenDefault()
}

func reorderFlagsBeforeArgs(args []string) []string {
	return flagutil.ReorderFlagsBeforeArgs(args)
}

// WriteShellHandoffFn is wired by cmd/di_hooks.go.
var WriteShellHandoffFn func(targetPath string)

func WriteShellHandoff(targetPath string) {
	if WriteShellHandoffFn != nil {
		WriteShellHandoffFn(targetPath)
	}
}

// ExtractPositionalArgsFn is wired by cmd/di_hooks.go.
var ExtractPositionalArgsFn func(args []string) []string

func extractPositionalArgs(args []string) []string {
	if ExtractPositionalArgsFn != nil {
		return ExtractPositionalArgsFn(args)
	}
	return nil
}
