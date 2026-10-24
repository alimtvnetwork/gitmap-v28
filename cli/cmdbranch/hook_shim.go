package cmdbranch

import "github.com/alimtvnetwork/gitmap-v28/cli/gitutil"

// CheckHelpFn is wired by cmd/di_hooks.go to the canonical implementation.
var CheckHelpFn func(command string, args []string)

func checkHelp(command string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
	}
}

// PrintHelpAndExitFn is wired by cmd/di_hooks.go to the canonical implementation.
var PrintHelpAndExitFn func(command string, args []string)

func printHelpAndExit(command string, args []string) {
	if PrintHelpAndExitFn != nil {
		PrintHelpAndExitFn(command, args)
	}
}

// IsInsideWorkTreeFn allows overriding git work tree check in tests.
var IsInsideWorkTreeFn func() bool

func isInsideWorkTree() bool {
	if IsInsideWorkTreeFn != nil {
		return IsInsideWorkTreeFn()
	}

	return gitutil.IsInsideWorkTree()
}

// CheckoutBranchFn allows overriding git checkout execution in tests.
var CheckoutBranchFn func(repoPath, branch string) (string, error)

func checkoutBranch(repoPath, branch string) (string, error) {
	if CheckoutBranchFn != nil {
		return CheckoutBranchFn(repoPath, branch)
	}

	return gitutil.CheckoutBranch(repoPath, branch)
}

// ResolveDefaultBranchNameFn allows overriding default branch resolution in tests.
var ResolveDefaultBranchNameFn func(repoPath string) string

func resolveDefaultBranchName(repoPath string) string {
	if ResolveDefaultBranchNameFn != nil {
		return ResolveDefaultBranchNameFn(repoPath)
	}

	return gitutil.ResolveDefaultBranchName(repoPath)
}
