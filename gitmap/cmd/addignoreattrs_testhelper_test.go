package cmd

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/templates"
)

// mustResolveForTest loads (kind, lang) pairs from the embedded
// templates and fails the test on any error. Used by add* tests so the
// production code path (templates.Resolve) is exercised, not a fixture.
func mustResolveForTest(t *testing.T, kind string, langs ...string) []templates.Resolved {
	t.Helper()
	out := make([]templates.Resolved, 0, len(langs))
	for _, lang := range langs {
		r, err := templates.Resolve(kind, lang)
		if err != nil {
			t.Fatalf("Resolve(%s,%s): %v", kind, lang, err)
		}
		out = append(out, r)
	}

	return out
}
