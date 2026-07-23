package movemerge

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// TestRunMerge_PreferNewer_BothSidesByteEqual creates two temp
// folders with one identical file, one only-left, one only-right,
// and two conflicting files (one newer on LEFT, one newer on RIGHT),
// runs RunMerge with DirBoth + PreferNewer, then asserts both sides
// converged byte-for-byte to the expected union tree.
//
// Spec: spec/01-app/97-move-and-merge.md (acceptance items
// "merge-both copies missing files both ways" + "--prefer-newer
// override the bypass default").
func TestRunMerge_PreferNewer_BothSidesByteEqual(t *testing.T) {
	left, right := t.TempDir(), t.TempDir()
	now := time.Now()
	older, newer := now.Add(-2*time.Hour), now

	// identical on both sides
	seed(t, left, "shared/identical.txt", "same-content", older)
	seed(t, right, "shared/identical.txt", "same-content", older)

	// only on one side -> copied to the other
	seed(t, left, "only-left.md", "L-content", older)
	seed(t, right, "only-right.md", "R-content", older)

	// conflict where LEFT is newer -> LEFT wins
	seed(t, left, "conflict-left-newer.txt", "L-WINS", newer)
	seed(t, right, "conflict-left-newer.txt", "R-loses", older)

	// conflict where RIGHT is newer -> RIGHT wins
	seed(t, left, "nested/conflict-right-newer.txt", "L-loses", older)
	seed(t, right, "nested/conflict-right-newer.txt", "R-WINS", newer)

	leftEP := Endpoint{DisplayName: left, WorkingDir: left, Kind: EndpointFolder, Existed: true}
	rightEP := Endpoint{DisplayName: right, WorkingDir: right, Kind: EndpointFolder, Existed: true}
	opts := Options{
		Yes: true, Prefer: PreferNewer, NoCommit: true, NoPush: true,
		CommandName: constants.CmdMergeBoth, LogPrefix: constants.LogPrefixMergeBoth,
	}
	if err := RunMerge(leftEP, rightEP, DirBoth, opts); err != nil {
		t.Fatalf("RunMerge: %v", err)
	}

	want := map[string]string{
		"shared/identical.txt":            "same-content",
		"only-left.md":                    "L-content",
		"only-right.md":                   "R-content",
		"conflict-left-newer.txt":         "L-WINS",
		"nested/conflict-right-newer.txt": "R-WINS",
	}
	assertTreeEquals(t, "LEFT", left, want)
	assertTreeEquals(t, "RIGHT", right, want)
}

// TestRunMerge_PreferNewer_LeftOnlyDoesNotTouchRight verifies that
// merge-left + PreferNewer never writes into RIGHT even when LEFT
// is the newer side.
func TestRunMerge_PreferNewer_LeftOnlyDoesNotTouchRight(t *testing.T) {
	left, right := t.TempDir(), t.TempDir()
	now := time.Now()
	seed(t, left, "x.txt", "L-newer", now)
	seed(t, right, "x.txt", "R-older", now.Add(-time.Hour))
	rightSnapshot := snapshot(t, right)

	leftEP := Endpoint{DisplayName: left, WorkingDir: left, Kind: EndpointFolder, Existed: true}
	rightEP := Endpoint{DisplayName: right, WorkingDir: right, Kind: EndpointFolder, Existed: true}
	opts := Options{
		Yes: true, Prefer: PreferNewer, NoCommit: true, NoPush: true,
		CommandName: constants.CmdMergeLeft, LogPrefix: constants.LogPrefixMergeLeft,
	}
	if err := RunMerge(leftEP, rightEP, DirLeftOnly, opts); err != nil {
		t.Fatalf("RunMerge: %v", err)
	}
	if got := snapshot(t, right); !mapsEqual(got, rightSnapshot) {
		t.Errorf("RIGHT was modified by merge-left:\nbefore=%v\nafter=%v", rightSnapshot, got)
	}
	// LEFT keeps its newer copy (L wins under PreferNewer).
	assertTreeEquals(t, "LEFT", left, map[string]string{"x.txt": "L-newer"})
}

// seed writes content at root/rel with the given mtime, creating dirs.
func seed(t *testing.T, root, rel, content string, mt time.Time) {
	t.Helper()
	full := filepath.Join(root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", full, err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", full, err)
	}
	if err := os.Chtimes(full, mt, mt); err != nil {
		t.Fatalf("chtimes %s: %v", full, err)
	}
}

// snapshot returns rel-path -> content for every file under root.
func snapshot(t *testing.T, root string) map[string]string {
	t.Helper()
	out := map[string]string{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		bytes, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		out[filepath.ToSlash(rel)] = string(bytes)

		return nil
	})
	if err != nil {
		t.Fatalf("walk %s: %v", root, err)
	}

	return out
}

// assertTreeEquals fails the test when root's tree != want.
func assertTreeEquals(t *testing.T, label, root string, want map[string]string) {
	t.Helper()
	got := snapshot(t, root)
	if !mapsEqual(got, want) {
		t.Errorf("%s tree mismatch\n  got:  %v\n  want: %v", label, got, want)
	}
}

// mapsEqual reports whether two string maps are byte-equal.
func mapsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}

	return true
}
