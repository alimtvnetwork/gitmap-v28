package cmdrelease

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func openDB() (*store.DB, error) {
	return store.OpenDefault()
}

// RequireOnlineFn is wired by cmd/di_hooks.go.
var RequireOnlineFn func()

func requireOnline() {
	if RequireOnlineFn != nil {
		RequireOnlineFn()
	}
}

// SplitSemverFn is wired by cmd/di_hooks.go.
var SplitSemverFn func(v string) [3]int

func splitSemver(v string) [3]int {
	if SplitSemverFn != nil {
		return SplitSemverFn(v)
	}
	return [3]int{}
}
