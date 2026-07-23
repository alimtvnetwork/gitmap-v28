package templates

import (
	"sort"
	"testing"
)

// requiredLangsBothKinds is the canonical list of languages that MUST
// ship with BOTH ignore/<lang>.gitignore AND attributes/<lang>.gitattributes
// in the embedded corpus. New entries here are a deliberate commitment:
// the corpus and the docs will start advertising the lang to users, so
// dropping one half of the pair would silently degrade `templates init`
// (the attributes step would soft-skip with a dim notice instead of
// scaffolding the file).
//
// This list is intentionally a single source of truth — the "what
// languages does gitmap support out of the box?" question gets one
// authoritative answer, and the tests below enforce it.
var requiredLangsBothKinds = []string{
	"common",
	"csharp",
	"go",
	"java",
	"kotlin",
	"node",
	"php",
	"python",
	"ruby",
	"rust",
	"swift",
}

// TestCorpusParityRequiredLangsHaveBothKinds asserts that every entry in
// requiredLangsBothKinds resolves cleanly for BOTH `ignore` and
// `attributes`. Regression guard for two failure modes:
//
//  1. Someone deletes one half of a pair (e.g. removes
//     attributes/java.gitattributes) and the corresponding `templates
//     init java` invocation silently degrades to a soft-skip.
//  2. Someone bumps the required list without actually adding the files.
func TestCorpusParityRequiredLangsHaveBothKinds(t *testing.T) {
	for _, lang := range requiredLangsBothKinds {
		for _, kind := range []string{kindIgnore, kindAttributes} {
			r, err := Resolve(kind, lang)
			if err != nil {
				t.Errorf("required %s/%s missing from embedded corpus: %v", kind, lang, err)

				continue
			}
			if r.Source != SourceEmbed {
				t.Errorf("required %s/%s should resolve to SourceEmbed, got %v (overlay leak?)", kind, lang, r.Source)
			}
			if len(r.Content) == 0 {
				t.Errorf("required %s/%s resolved to empty content", kind, lang)
			}
		}
	}
}

// TestListAdvertisesEveryRequiredLang locks in that the user-visible
// `templates list` output (which `gitmap templates list` prints
// verbatim) advertises every required language for both kinds. If a
// lang is in the corpus on disk but List() does not surface it, the
// `gitmap templates init <lang>` workflow is effectively undiscoverable.
func TestListAdvertisesEveryRequiredLang(t *testing.T) {
	withTempHome(t)
	entries, err := List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}

	have := map[string]bool{}
	for _, e := range entries {
		have[e.Kind+"/"+e.Lang] = true
	}

	var missing []string
	for _, lang := range requiredLangsBothKinds {
		for _, kind := range []string{kindIgnore, kindAttributes} {
			key := kind + "/" + lang
			if !have[key] {
				missing = append(missing, key)
			}
		}
	}
	sort.Strings(missing)
	if len(missing) > 0 {
		t.Fatalf("templates list output is missing required entries:\n  - %v", missing)
	}
}

// TestNewLanguagesPresent is a focused regression guard for the
// java/ruby/php/swift/kotlin extension batch. Catches a subtle revert
// (e.g. someone restoring an older corpus snapshot) faster than the
// broader parity test above, with a tighter error message.
func TestNewLanguagesPresent(t *testing.T) {
	for _, lang := range []string{"java", "ruby", "php", "swift", "kotlin"} {
		for _, kind := range []string{kindIgnore, kindAttributes} {
			if _, err := Resolve(kind, lang); err != nil {
				t.Errorf("expected %s/%s in corpus (added in v3.26.0): %v", kind, lang, err)
			}
		}
	}
}
