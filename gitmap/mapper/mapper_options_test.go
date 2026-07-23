package mapper

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TestResolveDefaultBranch_EmptyFallsBackToConstant locks the contract
// that an empty BuildOptions.DefaultBranch resolves to the package
// default — preserving legacy behavior for callers (cmd/as.go,
// cmd/releaseautoregister.go) that don't set the field.
func TestResolveDefaultBranch_EmptyFallsBackToConstant(t *testing.T) {
	if got := resolveDefaultBranch(""); got != constants.DefaultBranch {
		t.Errorf("resolveDefaultBranch(\"\") = %q, want %q",
			got, constants.DefaultBranch)
	}
}

// TestResolveDefaultBranch_ExplicitOverrideHonored proves the CLI
// surface (gitmap scan --default-branch <name>) actually wins over
// the compiled-in default.
func TestResolveDefaultBranch_ExplicitOverrideHonored(t *testing.T) {
	cases := []string{"master", "trunk", "develop", "release/2026"}
	for _, name := range cases {
		if got := resolveDefaultBranch(name); got != name {
			t.Errorf("resolveDefaultBranch(%q) = %q, want %q", name, got, name)
		}
	}
}
