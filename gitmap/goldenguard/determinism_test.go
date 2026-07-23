package goldenguard

// Tests for the determinism pre-check. Failures are captured via a
// fake fataler instead of t.Run sub-tests because t.Fatalf calls
// runtime.Goexit, which would skip post-call assignments and also
// propagate failure to the parent test. The fake records the call
// and returns control to the test, which is the cleanest way to
// assert "this helper would have fatalled".

import (
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
)

// fakeFataler implements the unexported fataler interface. It
// records whether Fatalf was called and the formatted message, so
// tests can assert both that a failure occurred and (when useful)
// that the message mentions the right cause.
type fakeFataler struct {
	fataled bool
	msg     string
}

func (f *fakeFataler) Helper() {}
func (f *fakeFataler) Fatalf(format string, args ...interface{}) {
	f.fataled = true
	f.msg = fmt.Sprintf(format, args...)
}

func TestAssertWriterDeterministic_StableWriter_Passes(t *testing.T) {
	fake := &fakeFataler{}
	writer := func() ([]byte, error) { return []byte("stable-bytes"), nil }
	assertWriterDeterministicOn(fake, "stable", writer)
	if fake.fataled {
		t.Fatalf("stable writer must NOT trigger a failure; got: %s", fake.msg)
	}
}

func TestAssertWriterDeterministic_DriftingWriter_Fails(t *testing.T) {
	fake := &fakeFataler{}
	var counter int32
	writer := func() ([]byte, error) {
		n := atomic.AddInt32(&counter, 1)
		return []byte{byte(n)}, nil
	}
	assertWriterDeterministicOn(fake, "drifting", writer)
	if !fake.fataled {
		t.Fatalf("drifting writer must trigger a determinism failure")
	}
}

func TestAssertWriterDeterministic_ErroringWriter_Fails(t *testing.T) {
	fake := &fakeFataler{}
	writer := func() ([]byte, error) { return nil, errors.New("boom") }
	assertWriterDeterministicOn(fake, "erroring", writer)
	if !fake.fataled {
		t.Fatalf("erroring writer must trigger a failure (cannot be deterministic)")
	}
}

func TestAssertWriterDeterministic_FailureMessageMentionsLabel(t *testing.T) {
	fake := &fakeFataler{}
	var counter int32
	writer := func() ([]byte, error) {
		return []byte{byte(atomic.AddInt32(&counter, 1))}, nil
	}
	assertWriterDeterministicOn(fake, "my-writer", writer)
	if !fake.fataled || !containsString(fake.msg, "my-writer") {
		t.Fatalf("failure message must include the writer label; got: %q", fake.msg)
	}
}

func TestAllowUpdateAfterDeterminism_TriggerOff_NeverInvokesWriter(t *testing.T) {
	var calls int32
	writer := func() ([]byte, error) {
		atomic.AddInt32(&calls, 1)
		return []byte("x"), nil
	}
	if AllowUpdateAfterDeterminism(t, false, "off", writer) {
		t.Fatalf("trigger=false must return false")
	}
	if got := atomic.LoadInt32(&calls); got != 0 {
		t.Fatalf("writer must not run when trigger=false; got %d calls", got)
	}
}

func TestAllowUpdateAfterDeterminism_StableWriterGateOn_ReturnsTrue(t *testing.T) {
	t.Setenv(AllowUpdateEnv, allowUpdateValue)
	writer := func() ([]byte, error) { return []byte("ok"), nil }
	if !AllowUpdateAfterDeterminism(t, true, "stable", writer) {
		t.Fatalf("stable writer + gate ON must return true")
	}
}

// containsString is a tiny readability shim — the test only needs
// substring matching, not the heft of the strings package import in
// every assertion.
func containsString(haystack, needle string) bool {
	return len(haystack) >= len(needle) && indexOf(haystack, needle) >= 0
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
