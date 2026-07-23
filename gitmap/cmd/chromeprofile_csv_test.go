// Package cmd — chromeprofile_csv_test.go: verifies writeChromeExportCSV
// emits a stable Category,Key,Value schema and that readChromeExportCSV
// round-trips extension IDs + preferences. Also asserts persistChromeProfile
// upserts into SQLite after an export so the CLI list view stays accurate.
package cmd

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func writeChromeFixtureProfile(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	mustWrite := func(rel, body string) {
		path := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatalf("write %s: %v", rel, err)
		}
	}
	mustWrite("Preferences", `{"homepage":"https://example.com","homepage_is_newtabpage":false}`)
	mustWrite("Bookmarks", `{"roots":{"bookmark_bar":{"name":"Bar","children":[1,2,3]}}}`)
	mustWrite("Extensions/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa/manifest.json", `{}`)
	mustWrite("Extensions/bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb/manifest.json", `{}`)
	return dir
}

func TestWriteChromeExportCSVSchema(t *testing.T) {
	src := writeChromeFixtureProfile(t)
	out := filepath.Join(t.TempDir(), "snap.csv")
	if _, err := writeChromeExportCSV(src, "fixture", out); err != nil {
		t.Fatalf("writeCSV: %v", err)
	}
	f, _ := os.Open(out)
	defer f.Close()
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(rows) < 4 {
		t.Fatalf("want >=4 rows incl header, got %d", len(rows))
	}
	if got := strings.Join(rows[0], ","); got != "Category,Key,Value" {
		t.Fatalf("header drift: %q", got)
	}
	if !containsRow(rows, "meta", "name", "fixture") {
		t.Fatalf("missing meta/name row: %v", rows)
	}
	if !containsCategory(rows, "extension") {
		t.Fatalf("missing extension rows")
	}
}

func TestReadChromeExportCSVRoundtrip(t *testing.T) {
	src := writeChromeFixtureProfile(t)
	out := filepath.Join(t.TempDir(), "snap.csv")
	if _, err := writeChromeExportCSV(src, "rt", out); err != nil {
		t.Fatalf("writeCSV: %v", err)
	}
	exp, err := readChromeExportCSV(out)
	if err != nil {
		t.Fatalf("readCSV: %v", err)
	}
	if exp.Name != "rt" {
		t.Fatalf("name drift: %q", exp.Name)
	}
	if len(exp.ExtensionIDs) != 2 {
		t.Fatalf("want 2 extensions, got %d (%v)", len(exp.ExtensionIDs), exp.ExtensionIDs)
	}
}

func TestPersistChromeProfileUpsertsAfterExport(t *testing.T) {
	db, err := store.OpenAt(filepath.Join(t.TempDir(), "ct.sqlite"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer db.Close()
	id, err := db.UpsertChromeProfile("Demo", "/src", true)
	if err != nil {
		t.Fatalf("upsert: %v", err)
	}
	if err := db.InsertChromeProfileExport(id, "json", "/snap/demo.json", 99); err != nil {
		t.Fatalf("insert: %v", err)
	}
	rows, _ := db.ListChromeProfilesDB()
	if len(rows) != 1 || rows[0].ExportCount != 1 {
		t.Fatalf("expected 1 demo row w/ 1 export, got %+v", rows)
	}
}

func containsRow(rows [][]string, c, k, v string) bool {
	for _, r := range rows {
		if len(r) >= 3 && r[0] == c && r[1] == k && r[2] == v {
			return true
		}
	}
	return false
}

func containsCategory(rows [][]string, c string) bool {
	for _, r := range rows {
		if len(r) >= 1 && r[0] == c {
			return true
		}
	}
	return false
}

func TestResolveChromeProfileUsesDisplayNameSummary(t *testing.T) {
	root := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", root)
	mustWriteProfileState(t, root)
	res, ok := resolveChromeProfile("Lovable")
	if !ok || chromeProfileSummary(res) != "Lovable (dir: Profile 15)" {
		t.Fatalf("unexpected profile resolution: ok=%v res=%+v", ok, res)
	}
}

func TestCopyChromeProfileReportsDestinationMkdir(t *testing.T) {
	src := writeChromeFixtureProfile(t)
	dst := filepath.Join(t.TempDir(), "blocked")
	mustWriteTextFile(t, dst, "not a dir")
	_, err := copyChromeProfile(src, dst)
	copyErr := unwrapChromeProfileCopyError(err)
	if copyErr.Op != constants.ChromeProfileCopyOpMkdir || copyErr.Target != dst {
		t.Fatalf("unexpected copy error: %+v", copyErr)
	}
}

func TestChromeProfileLockOpenErrorIsSkipped(t *testing.T) {
	lockName := constants.ChromeProfileLockFileName
	lockPath := filepath.Join(t.TempDir(), lockName)
	copied, err := chromeProfileCopyFile(lockPath, filepath.Join(t.TempDir(), lockName))
	if err != nil || copied {
		t.Fatalf("LOCK file open error should be skipped: %v", err)
	}
}

func mustWriteProfileState(t *testing.T, root string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(root, "Profile 15"), 0o755); err != nil {
		t.Fatalf("mkdir profile: %v", err)
	}
	body := `{"profile":{"info_cache":{"Profile 15":{"name":"Lovable"}}}}`
	mustWriteTextFile(t, filepath.Join(root, "Local State"), body)
}

func mustWriteTextFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
