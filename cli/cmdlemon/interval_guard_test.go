package cmdlemon

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestIntervalGuard_Default(t *testing.T) {
	buf := &bytes.Buffer{}
	guard := NewIntervalGuard(buf)

	resolved := guard.ResolveInterval(0)
	if resolved != DefaultPollingInterval {
		t.Fatalf("expected %v, got %v", DefaultPollingInterval, resolved)
	}

	if buf.Len() > 0 {
		t.Fatalf("expected no output for default interval, got: %s", buf.String())
	}
}

func TestIntervalGuard_ClampingNotice(t *testing.T) {
	testCases := []struct {
		requested time.Duration
		expected  time.Duration
	}{
		{2 * time.Second, MinPollingInterval},
		{3 * time.Second, MinPollingInterval},
		{500 * time.Millisecond, MinPollingInterval},
		{4999 * time.Millisecond, MinPollingInterval},
	}

	for _, tc := range testCases {
		buf := &bytes.Buffer{}
		guard := NewIntervalGuard(buf)

		resolved := guard.ResolveInterval(tc.requested)
		if resolved != tc.expected {
			t.Errorf("for %v: expected %v, got %v", tc.requested, tc.expected, resolved)
		}

		if !strings.Contains(buf.String(), "Notice: interval clamped to minimum 5s safety threshold") {
			t.Errorf("expected clamp notice for %v, got: %q", tc.requested, buf.String())
		}
	}
}

func TestIntervalGuard_ValidIntervals(t *testing.T) {
	testCases := []time.Duration{
		5 * time.Second,
		8 * time.Second,
		15 * time.Second,
		1 * time.Minute,
	}

	for _, dur := range testCases {
		buf := &bytes.Buffer{}
		guard := NewIntervalGuard(buf)

		resolved := guard.ResolveInterval(dur)
		if resolved != dur {
			t.Errorf("expected %v, got %v", dur, resolved)
		}

		if buf.Len() > 0 {
			t.Errorf("expected no output for valid interval %v, got: %s", dur, buf.String())
		}
	}
}

func TestIntervalGuard_ParseAndResolve(t *testing.T) {
	buf := &bytes.Buffer{}
	guard := NewIntervalGuard(buf)

	dur, err := guard.ParseAndResolveInterval("2s")
	if err != nil || dur != MinPollingInterval {
		t.Fatalf("expected 5s without error, got %v (err: %v)", dur, err)
	}

	dur, err = guard.ParseAndResolveInterval("10s")
	if err != nil || dur != 10*time.Second {
		t.Fatalf("expected 10s without error, got %v (err: %v)", dur, err)
	}

	dur, err = guard.ParseAndResolveInterval("")
	if err != nil || dur != DefaultPollingInterval {
		t.Fatalf("expected default 8s, got %v (err: %v)", dur, err)
	}
}
