package cloneconcurrency

import (
	"runtime"
	"testing"
)

// Resolve contract — see resolve.go for the full description.
// These tests pin the three branches (negative / zero / positive)
// so a future "let's just clamp negatives" refactor cannot quietly
// degrade the CLI's hard-fail behavior.

func TestResolve_NegativeIsRejected(t *testing.T) {
	cases := []int{-1, -10, -1000}
	for _, n := range cases {
		got, ok := Resolve(n)
		if ok {
			t.Errorf("Resolve(%d) ok=true, want false", n)
		}
		if got != 0 {
			t.Errorf("Resolve(%d) workers=%d, want 0 on rejection", n, got)
		}
	}
}

func TestResolve_ZeroMeansAuto(t *testing.T) {
	got, ok := Resolve(0)
	if !ok {
		t.Fatalf("Resolve(0) ok=false, want true")
	}
	want := runtime.NumCPU()
	if want < 1 {
		want = 1
	}
	if got != want {
		t.Errorf("Resolve(0) = %d, want NumCPU=%d", got, want)
	}
}

func TestResolve_PositivePassesThrough(t *testing.T) {
	cases := []int{1, 2, 4, 32, 1024}
	for _, n := range cases {
		got, ok := Resolve(n)
		if !ok {
			t.Errorf("Resolve(%d) ok=false, want true", n)
		}
		if got != n {
			t.Errorf("Resolve(%d) = %d, want verbatim %d", n, got, n)
		}
	}
}
