package constants

import (
	"strings"
	"testing"
)

// TestHostedDocsFallbackContract pins the two constants that the runtime
// `gitmap hd` hosted-docs fallback depends on. If either drifts, `hd`
// silently regresses to the pre-fix "  ✗ Docs site directory not found"
// hard-exit (see spec/02-app-issues/34-hd-hosted-docs-fallback.md).
func TestHostedDocsFallbackContract(t *testing.T) {
	if DocsURL == "" {
		t.Fatal("DocsURL must be non-empty — hd fallback opens this URL")
	}

	if !strings.HasPrefix(DocsURL, "https://") {
		t.Fatalf("DocsURL must be an https:// URL, got %q", DocsURL)
	}

	if got := strings.Count(MsgHDHostedFallback, "%s"); got != 1 {
		t.Fatalf("MsgHDHostedFallback must contain exactly one %%s verb (for DocsURL), got %d in %q", got, MsgHDHostedFallback)
	}

	if strings.Contains(MsgHDHostedFallback, "%d") || strings.Contains(MsgHDHostedFallback, "%v") {
		t.Fatalf("MsgHDHostedFallback must not contain %%d or %%v verbs, got %q", MsgHDHostedFallback)
	}
}
