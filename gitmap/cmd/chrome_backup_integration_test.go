// Package cmd — integration tests covering the chrome backup → mutate →
// restore loop. They prove that tampering with a single tar member trips
// the SHA256 manifest check, and that --no-verify lets a power user
// bypass that check intentionally.
package cmd

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// seedChromeProfileTree lays down a tiny but realistic Chrome user-data
// layout (Local State + one profile with Preferences + Bookmarks) so the
// backup walker has something deterministic to capture.
func seedChromeProfileTree(t *testing.T, root string) {
	t.Helper()
	mustWrite := func(rel, body string) {
		full := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	mustWrite("Local State", `{"profile":{"info_cache":{}}}`)
	mustWrite(filepath.Join("Default", "Preferences"), `{"profile":{"name":"Default"}}`)
	mustWrite(filepath.Join("Default", "Bookmarks"), `{"roots":{"bookmark_bar":{"children":[]}}}`)
}

// mutateTarMember rewrites a single regular-file member inside a .tar.gz,
// preserving every other byte. Used to simulate on-disk corruption between
// backup and restore.
func mutateTarMember(t *testing.T, path, target, newBody string) {
	t.Helper()
	in, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	gzr, err := gzip.NewReader(in)
	if err != nil {
		t.Fatal(err)
	}
	tr := tar.NewReader(gzr)

	var buf bytes.Buffer
	gzw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gzw)

	found := false
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		body, err := io.ReadAll(tr)
		if err != nil {
			t.Fatal(err)
		}
		if hdr.Name == target {
			body = []byte(newBody)
			hdr.Size = int64(len(body))
			found = true
		}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write(body); err != nil {
			t.Fatal(err)
		}
	}
	_ = tw.Close()
	_ = gzw.Close()
	_ = gzr.Close()
	_ = in.Close()
	if !found {
		t.Fatalf("mutateTarMember: %q not found in %s", target, path)
	}
	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestChromeBackupRestoreDetectsTamperedMember(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "src")
	seedChromeProfileTree(t, src)

	tarball := filepath.Join(tmp, "snap.tar.gz")
	if _, err := writeChromeBackup(src, tarball); err != nil {
		t.Fatalf("writeChromeBackup: %v", err)
	}
	if _, err := writeChromeManifestWithSource(tarball, src); err != nil {
		t.Fatalf("writeChromeManifestWithSource: %v", err)
	}

	// Baseline: untouched tarball verifies cleanly.
	ok, miss, err := verifyChromeManifest(tarball)
	if err != nil || !ok || len(miss) != 0 {
		t.Fatalf("baseline verify failed: ok=%v miss=%v err=%v", ok, miss, err)
	}

	// Tamper with one member; verify must now report mismatch.
	mutateTarMember(t, tarball, "Default/Preferences", `{"profile":{"name":"tampered"}}`)
	ok, miss, err = verifyChromeManifest(tarball)
	if err != nil {
		t.Fatalf("verify after tamper errored: %v", err)
	}
	if ok {
		t.Fatalf("tampered tarball verified clean; expected mismatch")
	}
	foundPref := false
	for _, m := range miss {
		if m == "Default/Preferences" {
			foundPref = true
		}
	}
	if !foundPref {
		t.Fatalf("tampered member missing from mismatch list: %v", miss)
	}

	// --no-verify bypass: readChromeBackup still extracts the (corrupted)
	// payload without consulting the manifest, mirroring the CLI's
	// --no-verify branch which skips verifyChromeManifest entirely.
	dst := filepath.Join(tmp, "restored")
	if err := os.MkdirAll(dst, 0o755); err != nil {
		t.Fatal(err)
	}
	n, err := readChromeBackup(tarball, dst)
	if err != nil {
		t.Fatalf("readChromeBackup (no-verify path): %v", err)
	}
	if n == 0 {
		t.Fatalf("expected files to be restored under --no-verify, got 0")
	}
	got, err := os.ReadFile(filepath.Join(dst, "Default", "Preferences"))
	if err != nil {
		t.Fatalf("read restored Preferences: %v", err)
	}
	if string(got) != `{"profile":{"name":"tampered"}}` {
		t.Fatalf("restored payload mismatch: %s", string(got))
	}
}

func TestChromeBackupRecordsSourcePathHeader(t *testing.T) {
	tmp := t.TempDir()
	src := filepath.Join(tmp, "MultiProfileRoot")
	seedChromeProfileTree(t, src)

	tarball := filepath.Join(tmp, "snap.tar.gz")
	if _, err := writeChromeBackup(src, tarball); err != nil {
		t.Fatal(err)
	}
	if _, err := writeChromeManifestWithSource(tarball, src); err != nil {
		t.Fatal(err)
	}
	got := readChromeManifestSource(tarball)
	if got != filepath.ToSlash(src) {
		t.Fatalf("readChromeManifestSource = %q, want %q", got, filepath.ToSlash(src))
	}
}
