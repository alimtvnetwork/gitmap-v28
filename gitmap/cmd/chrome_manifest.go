// Package cmd — chrome_manifest.go: SHA256 manifest written alongside
// every chrome backup tarball (`<tarball>.sha256.txt`) and verified
// automatically before `chrome restore` extracts anything.
package cmd

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// chromeManifestSuffix is appended to the tarball path to form the
// manifest filename. Kept centralized so backup + restore agree.
const chromeManifestSuffix = ".sha256.txt"

// ChromeManifestEntry is one row of the manifest: tar member name + sha.
type ChromeManifestEntry struct {
	Name string
	SHA  string
}

// buildChromeManifest streams a tar.gz and returns sorted sha256 entries
// for every regular file inside it. Used by both write (post-backup) and
// verify (pre-restore) so the hashes are guaranteed comparable.
func buildChromeManifest(tarballPath string) ([]ChromeManifestEntry, error) {
	f, err := os.Open(tarballPath) //nolint:gosec
	if err != nil {
		return nil, err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	var out []ChromeManifestEntry
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if hdr.FileInfo().IsDir() {
			continue
		}
		h := sha256.New()
		if _, err := io.Copy(h, tr); err != nil { //nolint:gosec
			return nil, err
		}
		out = append(out, ChromeManifestEntry{Name: hdr.Name, SHA: hex.EncodeToString(h.Sum(nil))})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

// encodeChromeManifest renders entries as `<sha>  <name>` lines (sha256sum format).
func encodeChromeManifest(entries []ChromeManifestEntry) string {
	var b strings.Builder
	for _, e := range entries {
		fmt.Fprintf(&b, "%s  %s\n", e.SHA, e.Name)
	}
	return b.String()
}

// decodeChromeManifest parses `<sha>  <name>` lines back into entries.
// Lines starting with `#` are treated as metadata headers (e.g. `# source: <path>`)
// and ignored here; see readChromeManifestSource for header lookup.
func decodeChromeManifest(raw string) []ChromeManifestEntry {
	var out []ChromeManifestEntry
	for _, ln := range strings.Split(raw, "\n") {
		ln = strings.TrimSpace(ln)
		if ln == "" || strings.HasPrefix(ln, "#") {
			continue
		}
		parts := strings.SplitN(ln, "  ", 2)
		if len(parts) != 2 {
			continue
		}
		out = append(out, ChromeManifestEntry{SHA: parts[0], Name: parts[1]})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// writeChromeManifest computes + writes the sidecar manifest next to tarballPath.
// sourcePath is the on-disk profile root that was captured; it is stamped as a
// `# source: <path>` header so restores can target the matching profile path in
// multi-profile layouts. Pass "" to omit the header.
func writeChromeManifest(tarballPath string) (string, error) {
	return writeChromeManifestWithSource(tarballPath, "")
}

func writeChromeManifestWithSource(tarballPath, sourcePath string) (string, error) {
	entries, err := buildChromeManifest(tarballPath)
	if err != nil {
		return "", err
	}
	manifestPath := tarballPath + chromeManifestSuffix
	var b strings.Builder
	if sourcePath != "" {
		fmt.Fprintf(&b, "# source: %s\n", filepath.ToSlash(sourcePath))
	}
	b.WriteString(encodeChromeManifest(entries))
	if err := os.WriteFile(manifestPath, []byte(b.String()), 0o644); err != nil {
		return "", err
	}
	return manifestPath, nil
}

// readChromeManifestSource returns the `# source: <path>` header from the
// sidecar manifest if present, or "" when missing/unreadable.
func readChromeManifestSource(tarballPath string) string {
	raw, err := os.ReadFile(tarballPath + chromeManifestSuffix) //nolint:gosec
	if err != nil {
		return ""
	}
	for _, ln := range strings.Split(string(raw), "\n") {
		ln = strings.TrimSpace(ln)
		if strings.HasPrefix(ln, "# source:") {
			return strings.TrimSpace(strings.TrimPrefix(ln, "# source:"))
		}
	}
	return ""
}

// verifyChromeManifest re-hashes the tarball and diffs it against the
// sidecar manifest. Returns (matched, mismatched-names, error). When the
// sidecar is missing it returns matched=false with a descriptive error so
// callers can warn instead of hard-failing on older backups.
func verifyChromeManifest(tarballPath string) (bool, []string, error) {
	manifestPath := tarballPath + chromeManifestSuffix
	raw, err := os.ReadFile(manifestPath) //nolint:gosec
	if err != nil {
		return false, nil, fmt.Errorf("manifest missing (%s): %w", manifestPath, err)
	}
	want := decodeChromeManifest(string(raw))
	got, err := buildChromeManifest(tarballPath)
	if err != nil {
		return false, nil, err
	}
	wantMap := map[string]string{}
	for _, e := range want {
		wantMap[e.Name] = e.SHA
	}
	var mismatches []string
	for _, e := range got {
		if wantMap[e.Name] != e.SHA {
			mismatches = append(mismatches, e.Name)
		}
		delete(wantMap, e.Name)
	}
	for name := range wantMap {
		mismatches = append(mismatches, name+" (missing from tarball)")
	}
	sort.Strings(mismatches)
	return len(mismatches) == 0, mismatches, nil
}
