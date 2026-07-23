package cmd

// runPullAll is the explicit batch-pull entry point. It forwards to
// runPull with the --all flag injected, so the heavy lifting
// (resolution, parallelism, stop-on-fail, pending-task accounting)
// stays in one place. Useful for users and the right-click context
// menu's power-user "pull-all" action; equivalent to
// `gitmap pull --all <forwarded flags>`.
func runPullAll(args []string) {
	runPull(prependAll(args))
}

// prependAll injects --all at the front of args unless the caller
// already passed it (long or short form). Keeping it idempotent lets
// `gitmap pull-all --all` behave identically to `gitmap pull-all`.
func prependAll(args []string) []string {
	for _, a := range args {
		if a == "--all" || a == "-all" || a == "-a" {
			return args
		}
	}

	return append([]string{"--all"}, args...)
}
