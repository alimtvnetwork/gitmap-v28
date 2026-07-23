package cmd

// Exit-code contract tests for `gitmap scan`.
//
// Asserts the documented codes from spec/04-generic-cli/07-error-handling.md:
//
//   0 -- success: scan an empty tempdir (zero repos is still success)
//   1 -- failure: scan a path that doesn't exist on disk
//
// Scan has no interactive confirmation prompt, so there is no
// "user-canceled" exit path to assert. The clone-now suite covers
// that scenario for the clone family.

import (
	"path/filepath"
	"testing"
)

// TestScanCLI_ExitCodes drives `gitmap scan` end-to-end against a
// real built binary and asserts the documented exit codes.
func TestScanCLI_ExitCodes(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		args     func(t *testing.T) []string
		wantCode int
		wantTag  string // substring expected in stderr/stdout for failure cases
		stream   string // "stderr" or "" (skip content check)
	}{
		{
			name: "success_empty_dir",
			args: func(t *testing.T) []string {
				return []string{"scan", "--quiet", scanEmptyDir(t)}
			},
			wantCode: 0,
		},
		{
			name: "failure_missing_dir",
			args: func(t *testing.T) []string {
				return []string{"scan", "--quiet", filepath.Join(scanEmptyDir(t), "does-not-exist")}
			},
			wantCode: 1,
			wantTag:  "scan",
			stream:   "stderr",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name == "failure_missing_dir" {
				skipOnWindowsSubprocess(t)
			}
			t.Parallel()
			code, stdout, stderr := runGitmap(t, tc.args(t), "")
			if code != tc.wantCode {
				t.Fatalf("exit=%d want=%d\nstdout=%s\nstderr=%s",
					code, tc.wantCode, stdout, stderr)
			}
			// Some Windows CI configurations split or buffer
			// stdout/stderr unpredictably for short-lived
			// subprocesses; assert against the combined output
			// so the label check stays meaningful either way.
			combined := stdout + "\n" + stderr
			if tc.wantTag != "" && tc.stream == "stderr" && !containsCI(combined, tc.wantTag) {
				t.Fatalf("output missing %q\nstdout=%s\nstderr=%s",
					tc.wantTag, stdout, stderr)
			}
		})
	}
}

// scanEmptyDir returns a fresh tempdir for a scan target. Isolated
// helper so cases stay declarative.
func scanEmptyDir(t *testing.T) string {
	t.Helper()

	return t.TempDir()
}
