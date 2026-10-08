package cmdcommitpull
// RunCommitInFn is wired by cmd/di_hooks.go to cmdcommitin's runCommitIn.
var RunCommitInFn func(args []string) error

func runCommitIn(args []string) error {
	if RunCommitInFn != nil {
		return RunCommitInFn(args)
	}
	return nil
}
