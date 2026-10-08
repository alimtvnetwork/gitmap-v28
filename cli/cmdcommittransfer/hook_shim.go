package cmdcommittransfer

// CheckHelpFn is wired by cmd/di_hooks.go to the canonical implementation.
var CheckHelpFn func(command string, args []string)

func checkHelp(command string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
	}
}

// ResolveEndpointStringFn is wired by cmd/di_hooks.go to the canonical implementation.
var ResolveEndpointStringFn func(raw string) string

func resolveEndpointString(raw string) string {
	if ResolveEndpointStringFn != nil {
		return ResolveEndpointStringFn(raw)
	}
	return raw
}

// RunCommitPullFn is wired by cmd/di_hooks.go.
var RunCommitPullFn func(args []string) error

func runCommitPull(args []string) error {
	if RunCommitPullFn != nil {
		return RunCommitPullFn(args)
	}
	return nil
}

// RunMigrateWizardFn is wired by cmd/di_hooks.go.
var RunMigrateWizardFn func(args []string) error

func runMigrateWizard(args []string) error {
	if RunMigrateWizardFn != nil {
		return RunMigrateWizardFn(args)
	}
	return nil
}
