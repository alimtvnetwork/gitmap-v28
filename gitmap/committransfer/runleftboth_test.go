package committransfer

import (
	"strings"
	"testing"
)

// TestWithDirectionLabel locks in the prefix-suffix join so commit-both
// users can visually attribute each line to the right pass.
func TestWithDirectionLabel(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		prefix string
		suffix string
		want   string
	}{
		{"left to right", "[commit-both]", "(left→right)", "[commit-both] (left→right)"},
		{"right to left", "[commit-both]", "(right→left)", "[commit-both] (right→left)"},
		{"empty suffix", "[commit-right]", "", "[commit-right] "},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := withDirectionLabel(Options{LogPrefix: tc.prefix}, tc.suffix)
			if got.LogPrefix != tc.want {
				t.Fatalf("withDirectionLabel(%q, %q).LogPrefix = %q, want %q",
					tc.prefix, tc.suffix, got.LogPrefix, tc.want)
			}
		})
	}
}

// TestRunBothImmutableOptions guarantees that mutating LogPrefix in one
// pass does not leak into the caller's struct or the second pass. The
// regression we guard against: a previous draft mutated opts directly,
// so the second pass would have printed `[commit-both] (left→right) (right→left)`.
func TestRunBothImmutableOptions(t *testing.T) {
	t.Parallel()

	original := Options{LogPrefix: "[commit-both]"}
	first := withDirectionLabel(original, "(left→right)")
	second := withDirectionLabel(original, "(right→left)")

	if original.LogPrefix != "[commit-both]" {
		t.Fatalf("original mutated: got %q", original.LogPrefix)
	}
	if first.LogPrefix == second.LogPrefix {
		t.Fatalf("two passes share LogPrefix: %q", first.LogPrefix)
	}
	if !strings.Contains(first.LogPrefix, "left→right") {
		t.Fatalf("first pass missing direction tag: %q", first.LogPrefix)
	}
	if !strings.Contains(second.LogPrefix, "right→left") {
		t.Fatalf("second pass missing direction tag: %q", second.LogPrefix)
	}
}
