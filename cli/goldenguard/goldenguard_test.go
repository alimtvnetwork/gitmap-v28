package goldenguard

// Tests for the AllowUpdate dual gate. The function is tiny but its
// failure mode (silently allowing a CI fixture rewrite) is severe, so
// every branch is pinned: trigger-off, trigger-on+allow-on, and
// trigger-on+allow-bad cases that MUST fail loudly.
//
// Failure-path tests use a fake fatalReporter to capture the Fatalf
// invocation without aborting the outer test goroutine via
// runtime.Goexit (which would happen with a real *testing.T).

import (
	"fmt"
	"os"
	"testing"
)

// TestAllowUpdate_TriggerOff_IsFalse: when the per-test trigger is
// off the function must short-circuit to false WITHOUT consulting
// the env var. This is the hot path in CI — it must never touch
// os.Getenv-driven branches that could call t.Fatalf.
func TestAllowUpdate_TriggerOff_IsFalse(t *testing.T) {
	t.Setenv(AllowUpdateEnv, "1") // even with allow ON, trigger OFF wins
	if AllowUpdate(t, false) {
		t.Fatalf("AllowUpdate(false, allow=1) = true, want false "+
			"(trigger-off must short-circuit before reading %s)",
			AllowUpdateEnv)
	}
}

// TestAllowUpdate_BothOn_IsTrue: the only path that returns true.
// Documents the exact value pairing — trigger=true AND env=="1".
func TestAllowUpdate_BothOn_IsTrue(t *testing.T) {
	t.Setenv(AllowUpdateEnv, "1")
	if !AllowUpdate(t, true) {
		t.Fatalf("AllowUpdate(true, allow=1) = false, want true")
	}
}

// TestAllowUpdate_TriggerOnAllowMissing_Fails: the failure path that
// catches a stray -update flag or GITMAP_UPDATE_GOLDEN=1 in CI when
// the dedicated allow var was (correctly) NOT set.
func TestAllowUpdate_TriggerOnAllowMissing_Fails(t *testing.T) {
	_ = os.Unsetenv(AllowUpdateEnv)
	rec := &fakeFatalReporter{}
	_ = allowUpdate(rec, true)
	if !rec.fatalCalled {
		t.Fatalf("AllowUpdate(true, allow=<unset>) did NOT call Fatalf — " +
			"missing allow env var must abort the regenerate path")
	}
}

// TestAllowUpdate_TriggerOnAllowWrongValue_Fails: typo guard. The
// allow var is intentionally narrow (literal "1" only) so common
// misspellings ("true", "yes") fail closed instead of unlocking.
func TestAllowUpdate_TriggerOnAllowWrongValue_Fails(t *testing.T) {
	for _, bad := range []string{"true", "yes", "y", "TRUE", "0", " 1 "} {
		t.Run(bad, func(tt *testing.T) {
			tt.Setenv(AllowUpdateEnv, bad)
			rec := &fakeFatalReporter{}
			_ = allowUpdate(rec, true)
			if !rec.fatalCalled {
				tt.Fatalf("AllowUpdate accepted bogus allow=%q "+
					"(only literal \"1\" must unlock the gate)", bad)
			}
		})
	}
}

// fakeFatalReporter records Fatalf invocations without calling
// runtime.Goexit, so a single test goroutine can exercise the
// failure path many times without aborting itself.
type fakeFatalReporter struct {
	fatalCalled bool
	lastMessage string
}

func (f *fakeFatalReporter) Helper() {}

func (f *fakeFatalReporter) Fatalf(format string, args ...any) {
	f.fatalCalled = true
	f.lastMessage = fmt.Sprintf(format, args...)
}
