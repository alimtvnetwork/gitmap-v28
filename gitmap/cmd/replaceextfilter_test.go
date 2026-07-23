package cmd

import (
	"path/filepath"
	"sort"
	"testing"
)

// TestNormalizeExtListInsensitive pins the case-insensitive normalizer
// behavior: trim, lowercase, dot-prepend, dedup, drop empties.
func TestNormalizeExtListInsensitive(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []string
	}{
		{"empty", "", nil},
		{"single dotted", ".go", []string{".go"}},
		{"single bare", "go", []string{".go"}},
		{"mixed case folds", ".Go,MD", []string{".go", ".md"}},
		{"with spaces", "  .go , md ", []string{".go", ".md"}},
		{"dedup case-folded", ".go,go,.GO", []string{".go"}},
		{"drops empties and lone dot", ".,,.md, ", []string{".md"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := normalizeExtList(tc.in, true)
			if !equalStringSlice(got, tc.want) {
				t.Fatalf("got %v want %v", got, tc.want)
			}
		})
	}
}

// TestNormalizeExtListSensitive proves the case-sensitive variant
// preserves the user's original casing and treats `.GO` and `.go` as
// distinct entries (no dedup across cases).
func TestNormalizeExtListSensitive(t *testing.T) {
	got := normalizeExtList(".Go,GO,.go", false)
	want := []string{".Go", ".GO", ".go"}
	if !equalStringSlice(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

// TestMatchesExtFilterCaseModes locks both modes simultaneously so a
// future refactor cannot collapse them.
func TestMatchesExtFilterCaseModes(t *testing.T) {
	if !matchesExtFilter("/x/foo.go", nil, true) {
		t.Error("nil filter must match every file (insensitive)")
	}
	if !matchesExtFilter("/x/foo.go", nil, false) {
		t.Error("nil filter must match every file (sensitive)")
	}
	// insensitive: GO matches .go list
	if !matchesExtFilter("/x/foo.GO", []string{".go"}, true) {
		t.Error("insensitive: .GO should match .go list")
	}
	// sensitive: GO does NOT match .go list
	if matchesExtFilter("/x/foo.GO", []string{".go"}, false) {
		t.Error("sensitive: .GO must not match .go list")
	}
	// sensitive: GO matches .GO list
	if !matchesExtFilter("/x/foo.GO", []string{".GO"}, false) {
		t.Error("sensitive: .GO must match .GO list")
	}
	// no-extension file fails non-empty filter in both modes
	if matchesExtFilter("/x/Makefile", []string{".go"}, true) {
		t.Error("Makefile must not match a non-empty filter")
	}
	if matchesExtFilter("/x/Makefile", []string{".go"}, false) {
		t.Error("Makefile must not match a non-empty filter (sensitive)")
	}
}

// TestResolveExtCase pins the contract for --ext-case parsing: empty
// and "insensitive" both default to true; "sensitive" is the only
// false-producer; whitespace and casing are tolerated.
func TestResolveExtCase(t *testing.T) {
	cases := map[string]bool{
		"":              true,
		"insensitive":   true,
		"INSENSITIVE":   true,
		"  sensitive  ": false,
		"Sensitive":     false,
	}
	for in, want := range cases {
		if got := resolveExtCase(in); got != want {
			t.Errorf("resolveExtCase(%q) = %v, want %v", in, got, want)
		}
	}
}

// TestParseReplaceFlagsExtAndCase walks the parser end-to-end with the
// --ext / --ext-case combos that production code paths exercise.
func TestParseReplaceFlagsExtAndCase(t *testing.T) {
	cases := []struct {
		name       string
		args       []string
		wantExts   []string
		wantInsens bool
	}{
		{"default insensitive", []string{"--ext", ".Go,.MD", "old", "new"},
			[]string{".go", ".md"}, true},
		{"explicit insensitive", []string{"--ext", ".Go,.MD", "--ext-case", "insensitive", "old", "new"},
			[]string{".go", ".md"}, true},
		{"sensitive preserves case", []string{"--ext", ".Go,.MD", "--ext-case", "sensitive", "old", "new"},
			[]string{".Go", ".MD"}, false},
		{"equals form", []string{"--ext=.go", "--ext-case=sensitive", "old", "new"},
			[]string{".go"}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			opts, pos, err := parseReplaceFlags(tc.args)
			if err != nil {
				t.Fatalf("parseReplaceFlags: %v", err)
			}
			if !equalStringSlice(opts.exts, tc.wantExts) {
				t.Errorf("opts.exts = %v, want %v", opts.exts, tc.wantExts)
			}
			if opts.extCaseIns != tc.wantInsens {
				t.Errorf("opts.extCaseIns = %v, want %v", opts.extCaseIns, tc.wantInsens)
			}
			if len(pos) != 2 || pos[0] != "old" || pos[1] != "new" {
				t.Errorf("positional = %v, want [old new]", pos)
			}
		})
	}
}

// TestWalkRepoFilesExtCaseSensitive seeds files with mixed-case
// extensions and proves --ext-case=sensitive only picks up byte-exact
// matches (vs the insensitive case which also picks up CHANGELOG.MD).
func TestWalkRepoFilesExtCaseSensitive(t *testing.T) {
	root := t.TempDir()

	mustWriteFile(t, filepath.Join(root, "README.md"), []byte("doc\n"))
	mustWriteFile(t, filepath.Join(root, "CHANGELOG.MD"), []byte("upper-case ext\n"))
	mustWriteFile(t, filepath.Join(root, "src", "app.go"), []byte("package app\n"))

	got, err := walkRepoFiles(root, []string{".md", ".go"}, false)
	if err != nil {
		t.Fatalf("walkRepoFiles: %v", err)
	}
	rels := relativizeAll(t, root, got)
	sort.Strings(rels)

	want := []string{"README.md", "src/app.go"}
	if !equalStringSlice(rels, want) {
		t.Fatalf("sensitive walk returned %v, want %v", rels, want)
	}
}
