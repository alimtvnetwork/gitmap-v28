package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// chdirT is a t.Helper that cds into dir for the duration of the test
// and restores the prior CWD on cleanup. Required because runTemplatesInit
// resolves targets relative to os.Getwd(), and tests must not pollute the
// repo's working tree with stray .gitignore / .gitattributes files.
func chdirT(t *testing.T, dir string) {
	t.Helper()
	// t.Chdir registers cleanup AFTER any earlier t.TempDir cleanups, so we
	// never end up restoring CWD to a deleted dir (avoids os.Getwd → "/"
	// fallback on macOS that breaks schema-path resolution in sibling tests).
	t.Chdir(dir)
}

// TestParseTemplatesInitFlagsAcceptsMixedOrder locks in the
// reorderFlagsBeforeArgs contract: flags before, between, or after
// positionals all parse identically. Regression guard for when someone
// "simplifies" parseTemplatesInitFlags and breaks `init --lfs go` or
// `init go --lfs python`.
func TestParseTemplatesInitFlagsAcceptsMixedOrder(t *testing.T) {
	cases := [][]string{
		{"go", "--lfs", "--dry-run"},
		{"--lfs", "go", "--dry-run"},
		{"--dry-run", "--lfs", "go"},
		{"go", "node", "--force"},
	}
	for _, args := range cases {
		flags, err := parseTemplatesInitFlags(args)
		if err != nil {
			t.Fatalf("parse %v: %v", args, err)
		}
		if len(flags.langs) == 0 {
			t.Errorf("%v: expected at least one lang, got none", args)
		}
	}
}

// TestParseTemplatesInitFlagsLFSAndDryRunStored verifies the bool
// flags actually flow into the struct (not silently dropped).
func TestParseTemplatesInitFlagsLFSAndDryRunStored(t *testing.T) {
	flags, err := parseTemplatesInitFlags([]string{"go", "--lfs", "--dry-run", "--force"})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if !flags.lfs || !flags.dryRun || !flags.force {
		t.Fatalf("expected all flags true, got %+v", flags)
	}
	if len(flags.langs) != 1 || flags.langs[0] != "go" {
		t.Fatalf("expected langs=[go], got %v", flags.langs)
	}
}

// TestExecuteTemplatesInitWritesIgnoreAndAttributes is the happy-path
// integration test: scaffold a fresh dir for go and confirm both files
// land on disk with the marker block. Runs end-to-end through
// templates.Resolve + templates.Merge so the test catches breakage in
// either layer.
func TestExecuteTemplatesInitWritesIgnoreAndAttributes(t *testing.T) {
	dir := t.TempDir()
	chdirT(t, dir)

	results := executeTemplatesInit(dir, templatesInitFlags{langs: []string{"go"}})

	// Expect 2 results: ignore/go (required) + attributes/go (optional, present in embed).
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d: %+v", len(results), results)
	}
	for _, r := range results {
		if r.skipped {
			t.Errorf("step %s/%s should not be skipped (embed has it): %s", r.step.kind, r.step.lang, r.skipReason)
		}
	}

	for _, name := range []string{".gitignore", ".gitattributes"} {
		path := filepath.Join(dir, name)
		body, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		if !strings.Contains(string(body), "# >>> gitmap:") {
			t.Errorf("%s missing gitmap marker block:\n%s", name, body)
		}
	}
}

// TestExecuteTemplatesInitIsIdempotent guarantees the second run is a
// byte-stable no-op (Changed=false on every step). This is the whole
// point of routing through templates.Merge — re-runs MUST NOT churn the
// file or cause a spurious git diff.
func TestExecuteTemplatesInitIsIdempotent(t *testing.T) {
	dir := t.TempDir()
	chdirT(t, dir)

	_ = executeTemplatesInit(dir, templatesInitFlags{langs: []string{"go"}})

	before, err := os.ReadFile(filepath.Join(dir, ".gitignore"))
	if err != nil {
		t.Fatalf("read .gitignore after first run: %v", err)
	}

	results := executeTemplatesInit(dir, templatesInitFlags{langs: []string{"go"}})
	for _, r := range results {
		if r.skipped {
			continue
		}
		if r.merge.Changed {
			t.Errorf("second run mutated %s (block=%s): expected Changed=false", r.merge.Path, r.merge.BlockTag)
		}
	}

	after, err := os.ReadFile(filepath.Join(dir, ".gitignore"))
	if err != nil {
		t.Fatalf("read .gitignore after second run: %v", err)
	}
	if !bytes.Equal(before, after) {
		t.Fatalf("idempotent run changed bytes:\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

// TestExecuteTemplatesInitDryRunWritesNothing verifies the --dry-run
// contract: zero filesystem mutations, but results still describe what
// WOULD happen so the summary printer can render a useful preview.
func TestExecuteTemplatesInitDryRunWritesNothing(t *testing.T) {
	dir := t.TempDir()
	chdirT(t, dir)

	results := executeTemplatesInit(dir, templatesInitFlags{
		langs:  []string{"go"},
		dryRun: true,
	})

	for _, name := range []string{".gitignore", ".gitattributes"} {
		if _, err := os.Stat(filepath.Join(dir, name)); !os.IsNotExist(err) {
			t.Fatalf("--dry-run wrote %s; expected it to remain absent", name)
		}
	}
	if len(results) == 0 {
		t.Fatal("expected dry-run to still return result rows for the printer")
	}
	for _, r := range results {
		if !r.skipped && !r.dryRun {
			t.Errorf("dry-run result for %s/%s missing dryRun flag", r.step.kind, r.step.lang)
		}
	}
}

// TestExecuteTemplatesInitForceReplacesExisting locks in --force semantics:
// pre-existing target files (with arbitrary user content) are wiped and
// replaced by a fresh gitmap-managed block. Without --force, that user
// content would survive (Merge would just append a marker block alongside).
func TestExecuteTemplatesInitForceReplacesExisting(t *testing.T) {
	dir := t.TempDir()
	chdirT(t, dir)

	stale := []byte("# my hand-written gitignore\nnode_modules/\n*.log\n")
	if err := os.WriteFile(filepath.Join(dir, ".gitignore"), stale, 0o644); err != nil {
		t.Fatalf("seed stale .gitignore: %v", err)
	}

	_ = executeTemplatesInit(dir, templatesInitFlags{
		langs: []string{"go"},
		force: true,
	})

	body, err := os.ReadFile(filepath.Join(dir, ".gitignore"))
	if err != nil {
		t.Fatalf("read .gitignore: %v", err)
	}
	if strings.Contains(string(body), "my hand-written gitignore") {
		t.Errorf("--force should have discarded stale content; got:\n%s", body)
	}
	if !strings.Contains(string(body), "# >>> gitmap:ignore/go >>>") {
		t.Errorf("--force result missing fresh gitmap block:\n%s", body)
	}
}

// TestExecuteTemplatesInitWithoutForcePreservesUserContent is the
// safety counterpart to the --force test: without --force, hand-written
// content OUTSIDE the gitmap marker block must survive untouched.
// Regression guard against accidentally making --force the default.
func TestExecuteTemplatesInitWithoutForcePreservesUserContent(t *testing.T) {
	dir := t.TempDir()
	chdirT(t, dir)

	hand := []byte("# my hand-written gitignore\nsecret.env\n")
	if err := os.WriteFile(filepath.Join(dir, ".gitignore"), hand, 0o644); err != nil {
		t.Fatalf("seed hand .gitignore: %v", err)
	}

	_ = executeTemplatesInit(dir, templatesInitFlags{langs: []string{"go"}})

	body, err := os.ReadFile(filepath.Join(dir, ".gitignore"))
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !strings.Contains(string(body), "secret.env") {
		t.Errorf("non-force run lost user content; got:\n%s", body)
	}
	if !strings.Contains(string(body), "# >>> gitmap:ignore/go >>>") {
		t.Errorf("non-force run missing appended gitmap block:\n%s", body)
	}
}

// TestExecuteTemplatesInitLFSMergesIntoAttributes confirms --lfs adds
// the lfs/common block to .gitattributes alongside the per-lang
// attributes block. Two distinct marker tags must coexist in the same
// file without colliding.
func TestExecuteTemplatesInitLFSMergesIntoAttributes(t *testing.T) {
	dir := t.TempDir()
	chdirT(t, dir)

	_ = executeTemplatesInit(dir, templatesInitFlags{
		langs: []string{"go"},
		lfs:   true,
	})

	body, err := os.ReadFile(filepath.Join(dir, ".gitattributes"))
	if err != nil {
		t.Fatalf("read .gitattributes: %v", err)
	}
	want := []string{
		"# >>> gitmap:attributes/go >>>",
		"# >>> gitmap:lfs/common >>>",
	}
	for _, w := range want {
		if !strings.Contains(string(body), w) {
			t.Errorf("missing marker %q in .gitattributes:\n%s", w, body)
		}
	}
}

// TestExecuteTemplatesInitSoftSkipsMissingAttributes locks in the
// "ignore is required, attributes is optional" contract by using a lang
// the embed corpus does NOT have an attributes file for. Today every
// embedded lang has both, so we use a synthetic check: temporarily this
// asserts the result struct's skip path works end-to-end via a known
// missing kind/lang pair.
//
// We exercise the skip path directly via runTemplatesInitStep with
// required=false, since faking a missing embed entry without modifying
// the binary is more brittle than testing the function's contract.
func TestRunTemplatesInitStepSoftSkipsMissingOptional(t *testing.T) {
	dir := t.TempDir()
	chdirT(t, dir)

	step := templatesInitStep{
		kind:   "attributes",
		lang:   "definitely-does-not-exist-zzz",
		tag:    "attributes/definitely-does-not-exist-zzz",
		target: filepath.Join(dir, ".gitattributes"),
	}

	r := runTemplatesInitStep(step, templatesInitFlags{}, false)
	if !r.skipped {
		t.Fatalf("expected skip for missing optional template, got merged result: %+v", r)
	}
	if r.skipReason == "" {
		t.Error("skip result should carry a non-empty skipReason for the summary printer")
	}
	if _, err := os.Stat(step.target); !os.IsNotExist(err) {
		t.Error("soft-skip path must not create the target file")
	}
}
