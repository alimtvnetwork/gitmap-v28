// Package goldenguard centralizes the policy that gates EVERY
// golden-fixture rewrite path in the gitmap test suite.
//
// Why this exists: each golden test has a "regenerate" mode triggered
// by either an env var (GITMAP_UPDATE_GOLDEN=1) or a -update test
// flag. Both are easy to flip on accidentally — by a CI matrix that
// inherits a developer's shell env, a stray flag in a Makefile, or a
// typo in a `go test` invocation. A silent regenerate is the worst
// kind of failure because it makes the very test that should catch
// drift WRITE the drifted bytes to disk and pass.
//
// The fix is a defense-in-depth gate: a regenerate happens ONLY when
// the per-test trigger AND a dedicated allow-env-var are BOTH set.
// CI must never set the allow-env-var; humans set it explicitly when
// they intend to rewrite fixtures. If the trigger is on but the allow
// var is off, AllowUpdate fails the test with a clear remediation
// message so the CI failure points straight at the misconfiguration.
//
// Usage from any golden test:
//
//	if goldenguard.AllowUpdate(t, *updateFlag) {
//	    os.WriteFile(path, got, 0o644)
//	    return
//	}
package goldenguard

import (
	"os"
	"testing"
)

// AllowUpdateEnv is the dedicated env var that must be set to "1"
// to permit ANY golden fixture rewrite, regardless of the per-test
// trigger (env var or -update flag). Centralized so CI configs and
// docs reference a single name. Do NOT add a second name — the
// whole point is that there's exactly one switch CI can pin OFF.
const AllowUpdateEnv = "GITMAP_ALLOW_GOLDEN_UPDATE"

// allowUpdateValue is the literal value AllowUpdateEnv must hold to
// unlock regenerate mode. "1" is intentionally narrow — "true",
// "yes", "y" etc. are all treated as OFF so a typo can't accidentally
// flip the gate on.
const allowUpdateValue = "1"

// fatalReporter is the minimal subset of *testing.T that AllowUpdate
// needs. Defined as an interface so internal tests can supply a
// recorder that captures the Fatalf invocation without aborting the
// outer test goroutine via runtime.Goexit. Production callers pass
// a real *testing.T, which satisfies this interface.
type fatalReporter interface {
	Helper()
	Fatalf(format string, args ...any)
}

// AllowUpdate reports whether the caller may rewrite golden fixtures.
// Returns true ONLY when the per-test trigger is on AND the dedicated
// allow-env-var is set to "1". When the trigger is on but the allow
// var is off, the test fails loudly via t.Fatalf — silent skipping
// would let CI green-light a misconfigured regenerate attempt.
//
// trigger: the per-test signal (e.g. *flag.Bool("update") for the cmd
// package, os.Getenv("GITMAP_UPDATE_GOLDEN")=="1" for the formatter
// package). Caller computes it.
func AllowUpdate(t *testing.T, trigger bool) bool {
	t.Helper()

	return allowUpdate(t, trigger)
}

// allowUpdate is the interface-typed core so internal tests can
// inject a fake reporter; production code keeps the *testing.T
// signature on AllowUpdate to avoid changing every caller.
func allowUpdate(t fatalReporter, trigger bool) bool {
	t.Helper()
	if !trigger {

		return false
	}
	allow := os.Getenv(AllowUpdateEnv)
	if allow == allowUpdateValue {

		return true
	}
	t.Fatalf("golden update requested but %s is not set to %q "+
		"(got %q). This double-gate prevents accidental fixture "+
		"rewrites in CI. To regenerate locally, run with BOTH the "+
		"per-test trigger AND %s=%s exported in your shell.",
		AllowUpdateEnv, allowUpdateValue, allow,
		AllowUpdateEnv, allowUpdateValue)

	return false
}
