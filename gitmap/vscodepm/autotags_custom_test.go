package vscodepm

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// makeRepoFixture writes a few marker files so DetectTagsCustom has
// something concrete to inspect. Returns the temp dir path.
func makeRepoFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, name := range []string{".git", "go.mod", "Gemfile"} {
		if err := os.WriteFile(filepath.Join(root, name), []byte("x"), 0o644); err != nil {
			t.Fatalf("seed %q: %v", name, err)
		}
	}

	return root
}

// resetTagEnv unsets every tag env var so cases run in isolation.
func resetTagEnv(t *testing.T) {
	t.Helper()
	for _, k := range []string{
		constants.EnvVSCodeTagAdd,
		constants.EnvVSCodeTagSkip,
		constants.EnvVSCodeTagMarker,
	} {
		os.Unsetenv(k)
	}
}

// TestDetectTagsCustomNoEnvMatchesBuiltin asserts that with no env
// overrides, DetectTagsCustom returns DetectTags PLUS the canonical
// "gitmap" brand tag (always emitted so projects.json entries are
// self-identifying — see prependGitmapBrand in autotags_custom.go).
func TestDetectTagsCustomNoEnvMatchesBuiltin(t *testing.T) {
	resetTagEnv(t)
	root := makeRepoFixture(t)

	got := DetectTagsCustom(root)
	want := append([]string{constants.AutoTagGitmap}, DetectTags(root)...)
	if !reflect.DeepEqual(sortCopy(got), sortCopy(want)) {
		t.Errorf("custom = %v, want %v (= gitmap brand + builtin)", got, want)
	}
	if len(got) == 0 || got[0] != constants.AutoTagGitmap {
		t.Errorf("gitmap brand tag must lead the list, got %v", got)
	}
}

// TestDetectTagsCustomGitmapBrandAlwaysPresent confirms the brand tag
// is emitted even for paths with zero detected markers (an empty dir,
// a missing path, or an empty rootPath string). This guarantees
// every projects.json entry written by gitmap is greppable by tag.
func TestDetectTagsCustomGitmapBrandAlwaysPresent(t *testing.T) {
	resetTagEnv(t)
	cases := map[string]string{
		"empty-dir":   t.TempDir(),
		"missing":     filepath.Join(t.TempDir(), "nope"),
		"empty-input": "",
	}
	for name, root := range cases {
		got := DetectTagsCustom(root)
		if !containsString(got, constants.AutoTagGitmap) {
			t.Errorf("%s: gitmap brand missing from %v", name, got)
		}
	}
}

// TestDetectTagsCustomGitmapSkippable confirms users can opt out of
// the brand tag via --vscode-tag-skip gitmap (env: GITMAP_VSCODE_TAG_SKIP).
func TestDetectTagsCustomGitmapSkippable(t *testing.T) {
	resetTagEnv(t)
	defer resetTagEnv(t)
	os.Setenv(constants.EnvVSCodeTagSkip, constants.AutoTagGitmap)

	got := DetectTagsCustom(makeRepoFixture(t))
	if containsString(got, constants.AutoTagGitmap) {
		t.Errorf("expected gitmap brand to be skipped, got %v", got)
	}
}

// TestDetectTagsCustomGitmapNotDuplicated confirms that if the brand
// tag is already requested via --vscode-tag gitmap, it is emitted
// exactly once (no duplicates in the final tag array).
func TestDetectTagsCustomGitmapNotDuplicated(t *testing.T) {
	resetTagEnv(t)
	defer resetTagEnv(t)
	os.Setenv(constants.EnvVSCodeTagAdd, constants.AutoTagGitmap)

	got := DetectTagsCustom(makeRepoFixture(t))
	count := 0
	for _, tag := range got {
		if tag == constants.AutoTagGitmap {
			count++
		}
	}
	if count != 1 {
		t.Errorf("gitmap appeared %d times in %v, want 1", count, got)
	}
}

// TestDetectTagsCustomSkipsAndAdds covers the skip + add overlay,
// including the documented case where the same tag appears in both
// (always-add wins because it runs after skip).
func TestDetectTagsCustomSkipsAndAdds(t *testing.T) {
	resetTagEnv(t)
	defer resetTagEnv(t)
	root := makeRepoFixture(t)

	os.Setenv(constants.EnvVSCodeTagSkip, "git")
	os.Setenv(constants.EnvVSCodeTagAdd, "work"+constants.EnvVSCodeTagSeparator+"git")

	got := DetectTagsCustom(root)
	// Detected portion = everything except the trailing always-add slice.
	addCount := 2 // "work" and "git"
	cutoff := len(got) - addCount
	if cutoff < 0 {
		cutoff = 0
	}
	if containsString(got[:cutoff], "git") {
		t.Errorf("git should be skipped from detected portion, got %v", got)
	}
	if !containsString(got, "work") || !containsString(got, "git") {
		t.Errorf("always-add tags missing, got %v", got)
	}
}

// TestDetectTagsCustomCustomMarker registers a Gemfile→ruby rule and
// verifies it shows up after the canonical built-in tags.
func TestDetectTagsCustomCustomMarker(t *testing.T) {
	resetTagEnv(t)
	defer resetTagEnv(t)
	root := makeRepoFixture(t)

	os.Setenv(constants.EnvVSCodeTagMarker, "Gemfile=ruby")

	got := DetectTagsCustom(root)
	if !containsString(got, "ruby") {
		t.Errorf("custom marker tag missing, got %v", got)
	}
}

func containsString(s []string, v string) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}

	return false
}

func sortCopy(s []string) []string {
	c := append([]string{}, s...)
	sort.Strings(c)

	return c
}
