package clonefrom

// Test-only helper for writing fixture files. Pulled into its own
// `_test.go` file so multiple test files in this package can share
// it without circular sharing through a non-test source file.

import "os"

// writeFile is a tiny os.WriteFile wrapper with a fixed permission.
// 0o644 matches what Git itself writes for tracked files — keeps
// the fixture's permission bits boring so platform-specific umask
// surprises can't perturb tests.
func writeFile(path, body string) error {
	return os.WriteFile(path, []byte(body), 0o644)
}
