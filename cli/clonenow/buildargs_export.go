package clonenow

// buildargs_export.go — thin exported wrapper around buildGitArgs so
// the cmd package's --verify-cmd-faithful checker can ask the
// executor "what argv would you actually pass to git for this row?"
// without duplicating the assembly logic. Single source of truth: if
// buildGitArgs ever changes (e.g. adds --filter), the verifier sees
// the change automatically and the printed cmd: stays in sync.
//
// Kept in its own file (vs. simply renaming buildGitArgs) to avoid
// touching the test file's reference and to keep the export surface
// trivially greppable: `git grep BuildGitArgs` returns this single
// shim across the repo.

// BuildGitArgs returns the argv (excluding the "git" binary) that
// Execute would pass to exec.Command for the given row. Mirrors the
// real call site in executeRow exactly — no recomputation.
func BuildGitArgs(r Row, url, dest string) []string {
	return buildGitArgs(r, url, dest)
}
