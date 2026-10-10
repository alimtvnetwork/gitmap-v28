package cmdcommitin

// RunCommitPullBootstrapFn is wired by cmd/di_hooks.go.
var RunCommitPullBootstrapFn func(args []string) error

// PrintCommitPullTreeFn is wired by cmd/di_hooks.go.
var PrintCommitPullTreeFn func()
