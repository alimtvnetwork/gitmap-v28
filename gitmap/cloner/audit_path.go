package cloner

import "os"

// pathExists is a stat-only existence check used by the audit planner.
// Kept separate from IsGitRepo so the audit never has to reason about
// whether a stat error means "missing" vs "permission denied" — both are
// surfaced as "not present" to keep the report deterministic on
// permission-restricted shares.
func pathExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}
