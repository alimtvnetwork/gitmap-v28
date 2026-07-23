package cmd

// Exit-code contract tests for the clone command family:
//
//   gitmap clone        -- direct-URL clone wrapper
//   gitmap clone-from   -- manifest-driven re-clone (dry-run + execute)
//   gitmap clone-now    -- scan-output-driven re-clone (alias: reclone)
//   gitmap clone-next   -- next-batch clone with version flatten
//   gitmap clone-pick   -- sparse-checkout pick
//
// We assert the three documented scenarios per command:
//
//   success (0)        -- a `--help` invocation, which every command
//                          short-circuits to exit 0 via checkHelp().
//                          That gives us a stable success signal that
//                          doesn't need network or git access.
//   user-canceled (2)  -- only `clone-now` has a documented user-cancel
//                          exit code (CloneNowExitConfirmAborted=2)
//                          reachable without a TTY: feed it a manifest
//                          with an existing destination and no --yes,
//                          and the non-TTY branch aborts with code 2.
//   failure (1 or 2)   -- a missing positional / missing manifest, which
//                          every command rejects before doing real work.

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TestCloneCLI_HelpExitsZero asserts the success contract: every
// clone-family entrypoint prints help and exits 0 when --help is
// passed, no matter what other (potentially invalid) flags follow.
func TestCloneCLI_HelpExitsZero(t *testing.T) {
	t.Parallel()
	for _, cmd := range cloneFamilyCmds() {
		t.Run("help_"+cmd, func(t *testing.T) {
			t.Parallel()
			code, stdout, stderr := runGitmap(t, []string{cmd, "--help"}, "")
			if code != 0 {
				t.Fatalf("`%s --help` exit=%d want=0\nstdout=%s\nstderr=%s",
					cmd, code, stdout, stderr)
			}
		})
	}
}

// TestCloneCLI_FailureExitCodes asserts each command rejects the
// trivial bad-invocation with a non-zero code. The exact code is
// per-command (see comments in the source files); we assert the
// command-specific value rather than a generic >0 so a regression
// from "exit 2 (bad usage)" to "exit 1 (operational failure)" is
// surfaced as a contract violation, not silently absorbed.
func TestCloneCLI_FailureExitCodes(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name     string
		args     []string
		wantCode int
	}{
		{"clone_no_args", []string{"clone"}, 1},
		{"clonefrom_missing_file", []string{"clone-from", "/no/such/manifest.json"}, 1},
		{"clonenow_missing_file", []string{"clone-now", "/no/such/manifest.json"}, 1},
		{"clonepick_no_args", []string{"clone-pick"}, 2},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			code, stdout, stderr := runGitmap(t, tc.args, "")
			if code != tc.wantCode {
				t.Fatalf("%s: exit=%d want=%d\nstdout=%s\nstderr=%s",
					tc.name, code, tc.wantCode, stdout, stderr)
			}
		})
	}
}

// TestCloneNowCLI_UserCanceledNonTTY asserts the documented user-
// canceled exit (CloneNowExitConfirmAborted = 2) for the only
// command in the family with a non-TTY-safe cancel path. We seed a
// manifest whose row's RelativePath already exists on disk; with no
// --yes and no TTY (the subprocess inherits /dev/null for stdin via
// the test harness), the confirm gate must abort with exactly code
// 2. We also assert the gate's stderr message is present so a
// future regression that exits 2 from an unrelated path can't
// silently satisfy this test.
func TestCloneNowCLI_UserCanceledNonTTY(t *testing.T) {
	skipOnWindowsSubprocess(t)
	t.Parallel()
	root := t.TempDir()
	existing := filepath.Join(root, "alpha")
	if err := os.MkdirAll(existing, 0o755); err != nil {
		t.Fatalf("seed existing dest: %v", err)
	}
	manifest := writeCloneNowManifest(t, root)

	args := []string{
		"clone-now",
		"--execute",
		"--quiet",
		"--cwd", root,
		manifest,
	}
	code, stdout, stderr := runGitmap(t, args, "")
	if code != constants.CloneNowExitConfirmAborted {
		t.Fatalf("clone-now non-TTY confirm: exit=%d want %d (CloneNowExitConfirmAborted)\nstdout=%s\nstderr=%s",
			code, constants.CloneNowExitConfirmAborted, stdout, stderr)
	}
	// Combine streams so Windows CI configurations that buffer
	// or redirect stderr differently still satisfy the substring
	// contract — the message-presence check is what matters.
	combined := stdout + "\n" + stderr
	if !strings.Contains(combined, constants.MsgCloneNowConfirmNonTTY) {
		t.Fatalf("clone-now non-TTY confirm: output missing non-TTY gate message\nwant substring=%q\nstdout=%s\nstderr=%s",
			constants.MsgCloneNowConfirmNonTTY, stdout, stderr)
	}
}

// cloneFamilyCmds is the canonical list of command IDs whose help
// surface this suite verifies. Lifted into a helper so adding a new
// clone variant is one line.
func cloneFamilyCmds() []string {
	return []string{
		"clone",
		"clone-from",
		"clone-now",
		"clone-next",
		"clone-pick",
	}
}

// writeCloneNowManifest writes a minimal JSON manifest with one row
// whose RelativePath is "alpha", matching the existing dir seeded by
// TestCloneNowCLI_UserCanceledNonTTY.
func writeCloneNowManifest(t *testing.T, root string) string {
	t.Helper()
	path := filepath.Join(root, "manifest.json")
	body := strings.Join([]string{
		`[`,
		`  {`,
		`    "repoName": "alpha",`,
		`    "httpsUrl": "https://example.invalid/alpha.git",`,
		`    "sshUrl": "git@example.invalid:alpha.git",`,
		`    "relativePath": "alpha",`,
		`    "branch": "main"`,
		`  }`,
		`]`,
	}, "\n")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	return path
}

// containsCI is a case-insensitive substring helper used by the scan
// suite. Lives here (next to the other clone-suite helpers) so the
// scan file stays focused on its own cases.
func containsCI(haystack, needle string) bool {
	return strings.Contains(strings.ToLower(haystack), strings.ToLower(needle))
}
