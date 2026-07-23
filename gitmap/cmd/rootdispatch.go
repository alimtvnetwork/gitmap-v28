package cmd

import (
	"os"
)

// dispatchEntry maps one or more command aliases to a handler.
// All aliases share a single handler closure, which keeps each
// dispatcher a flat table instead of a deep if/else chain.
type dispatchEntry struct {
	names   []string
	handler func()
}

// runDispatchTable looks up command in entries and invokes the
// matching handler. Returns true when a handler was found.
//
// This helper exists so individual dispatchers (dispatchCore,
// dispatchData, dispatchTooling, ...) stay below gocyclo's 15
// complexity threshold by replacing N `if` branches with one loop.
func runDispatchTable(command string, entries []dispatchEntry) bool {
	for _, entry := range entries {
		if matchAny(command, entry.names) {
			entry.handler()

			return true
		}
	}

	return false
}

// matchAny reports whether command equals any name in names.
func matchAny(command string, names []string) bool {
	for _, name := range names {
		if command == name {
			return true
		}
	}

	return false
}

// argsTail returns the slice of CLI args after the subcommand,
// matching the previous `os.Args[2:]` pattern used by handlers.
func argsTail() []string {
	return os.Args[2:]
}
