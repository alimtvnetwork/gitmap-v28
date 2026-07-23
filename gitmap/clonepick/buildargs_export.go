package clonepick

import "strconv"

// buildargs_export.go — exported wrapper that mirrors gitClonePartial's
// argv assembly so cmd/--verify-cmd-faithful can diff the rendered
// `cmd:` line against the real invocation.
//
// Why a shim instead of refactoring gitClonePartial to call us:
// gitClonePartial owns the runGit side-effect; the verifier only
// wants the pure argv. Splitting concerns keeps the executor's call
// site unchanged (low blast radius) while still giving the verifier
// a byte-identical builder. A unit test in cmd/ pins both code paths
// to the same expected argv so a future drift between this shim and
// gitClonePartial fails CI.

// BuildGitArgs returns the argv (excluding the "git" binary) that
// gitClonePartial would pass to runGit for the given plan/dest.
// Order matches gitClonePartial in sparse.go exactly:
// `clone --filter=blob:none --no-checkout [--branch B] [--depth N] URL DEST`.
func BuildGitArgs(plan Plan, dest string) []string {
	args := []string{"clone", "--filter=blob:none", "--no-checkout"}
	if len(plan.Branch) > 0 {
		args = append(args, "--branch", plan.Branch)
	}
	if plan.Depth > 0 {
		args = append(args, "--depth", strconv.Itoa(plan.Depth))
	}
	args = append(args, plan.RepoUrl, dest)

	return args
}
