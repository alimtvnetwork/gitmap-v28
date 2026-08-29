package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TestDecideDriftAction is the drift-guard integration test seam
// (step 34). It locks in the three-way policy that reverseOneRepo
// relies on so future refactors can't silently regress the guard.
func TestDecideDriftAction(t *testing.T) {
	pub := constants.VisibilityPublic
	pri := constants.VisibilityPrivate
	cases := []struct {
		name              string
		current, expected string
		force             bool
		want              driftAction
	}{
		{"no drift, no force → proceed", pub, pub, false, driftActionProceed},
		{"drift, no force → skip", pri, pub, false, driftActionSkip},
		{"no drift, force → force", pub, pub, true, driftActionForce},
		{"drift, force → force (override wins)", pri, pub, true, driftActionForce},
		{"empty current, no force → skip", "", pub, false, driftActionSkip},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := decideDriftAction(tc.current, tc.expected, tc.force)
			if got != tc.want {
				t.Fatalf("decideDriftAction(%q,%q,%v) = %q, want %q",
					tc.current, tc.expected, tc.force, got, tc.want)
			}
		})
	}
}
