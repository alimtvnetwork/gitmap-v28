package cmd

// fixrepo_strict_packages_test.go — unit tests for the pure helper
// that maps modified .go files to `go test` package patterns. Pure
// (no I/O, no exec) so the table-driven cases are deterministic on
// every OS. Production behavior assertions:
//
//   1. Files at the repo root produce the literal "." pattern (NOT
//      "./") because `go test .` is the canonical form for the root
//      package.
//   2. Files in subdirectories produce "./<rel-dir>" with FORWARD
//      slashes regardless of host OS — Windows-derived paths must
//      normalize so log readers and golden files stay portable.
//   3. The output is sorted + deduplicated: feeding 5 files from the
//      same package produces exactly one entry, and the iteration
//      order is stable across calls (map-backed dedup is randomized,
//      sort.Strings pins it).
//   4. Files outside repoRoot (would only happen via a symlink the
//      scan layer should already reject) are silently skipped — the
//      function never returns absolute paths or "../"-prefixed
//      patterns, both of which would force `go test` into modes the
//      caller cannot predict.
//   5. Empty input returns nil (not an empty slice) so callers can
//      use `len(...) == 0` AND distinguish nil from explicit empty
//      via the standard Go idiom.

import (
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// TestDerivePackagesFromGoFiles_EmptyInput pins contract #5: nil
// input → nil output (not a zero-length non-nil slice).
func TestDerivePackagesFromGoFiles_EmptyInput(t *testing.T) {
	if got := derivePackagesFromGoFiles("/repo", nil); got != nil {
		t.Fatalf("nil input: want nil, got %#v", got)
	}
	if got := derivePackagesFromGoFiles("/repo", []string{}); got != nil {
		t.Fatalf("empty input: want nil, got %#v", got)
	}
}

// TestDerivePackagesFromGoFiles_RootAndSubdirs covers the canonical
// happy path: a file at the root, two files in the same subdir
// (dedup), one file in a deeper subdir, all from a Unix-style
// repoRoot. Asserts both the patterns AND the sort order.
func TestDerivePackagesFromGoFiles_RootAndSubdirs(t *testing.T) {
	repoRoot := filepath.Join(string(filepath.Separator), "tmp", "repo")
	files := []string{
		filepath.Join(repoRoot, "main.go"),
		filepath.Join(repoRoot, "gitmap", "cmd", "fixrepo.go"),
		filepath.Join(repoRoot, "gitmap", "cmd", "fixrepo_strict.go"),
		filepath.Join(repoRoot, "gitmap", "constants", "constants_fixrepo.go"),
	}
	got := derivePackagesFromGoFiles(repoRoot, files)
	want := []string{".", "./gitmap/cmd", "./gitmap/constants"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("packages mismatch\nwant: %#v\n got: %#v", want, got)
	}
}

// TestDerivePackagesFromGoFiles_SkipsOutsideRoot covers contract #4:
// a path that filepath.Rel resolves to "../something" must be
// silently dropped, never appear in the output.
func TestDerivePackagesFromGoFiles_SkipsOutsideRoot(t *testing.T) {
	repoRoot := filepath.Join(string(filepath.Separator), "tmp", "repo")
	outside := filepath.Join(string(filepath.Separator), "tmp", "elsewhere", "foo.go")
	inside := filepath.Join(repoRoot, "pkg", "x.go")
	got := derivePackagesFromGoFiles(repoRoot, []string{outside, inside})
	want := []string{"./pkg"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("outside-root not skipped\nwant: %#v\n got: %#v", want, got)
	}
}

// TestDerivePackagesFromGoFiles_WindowsBackslashes covers contract #2:
// on Windows, filepath.Join uses backslashes; the emitted pattern MUST
// still use forward slashes. Skipped on non-Windows runners because
// filepath.Join on Unix would produce forward slashes already and the
// test would tautologically pass without exercising the normalization.
func TestDerivePackagesFromGoFiles_WindowsBackslashes(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("normalization only observable on Windows-style separators")
	}
	repoRoot := `C:\repo`
	files := []string{
		`C:\repo\gitmap\cmd\fixrepo.go`,
		`C:\repo\gitmap\cmd\fixrepo_strict.go`,
	}
	got := derivePackagesFromGoFiles(repoRoot, files)
	want := []string{"./gitmap/cmd"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("forward-slash normalization failed\nwant: %#v\n got: %#v", want, got)
	}
}

// TestDerivePackagesFromGoFiles_DeterministicOrder covers contract #3:
// identical inputs must produce identical outputs across calls. The
// dedup map is randomized in Go, so without the trailing sort.Strings
// this assertion would flake intermittently — and a flaky package
// list would make CI logs and downstream goldens nondeterministic.
func TestDerivePackagesFromGoFiles_DeterministicOrder(t *testing.T) {
	repoRoot := filepath.Join(string(filepath.Separator), "tmp", "repo")
	files := []string{
		filepath.Join(repoRoot, "z", "z.go"),
		filepath.Join(repoRoot, "a", "a.go"),
		filepath.Join(repoRoot, "m", "m.go"),
	}
	first := derivePackagesFromGoFiles(repoRoot, files)
	for i := 0; i < 20; i++ {
		got := derivePackagesFromGoFiles(repoRoot, files)
		if !reflect.DeepEqual(got, first) {
			t.Fatalf("nondeterministic output on iteration %d:\n first: %#v\n got:   %#v",
				i, first, got)
		}
	}
	want := []string{"./a", "./m", "./z"}
	if !reflect.DeepEqual(first, want) {
		t.Fatalf("not sorted\nwant: %#v\n got: %#v", want, first)
	}
}
