package templates

import (
	"bufio"
	"bytes"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

// TestEmbeddedCorpusHeaders enforces the audit-trail header contract on
// every embedded template file: source, kind, lang, version must be set.
func TestEmbeddedCorpusHeaders(t *testing.T) {
	walkErr := fs.WalkDir(FS, embedAssetsRoot, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if filepath.Base(p) == "README.md" {
			return nil
		}
		data, rErr := FS.ReadFile(p)
		if rErr != nil {
			t.Fatalf("read %s: %v", p, rErr)
		}
		assertHeader(t, p, data)

		return nil
	})
	if walkErr != nil {
		t.Fatalf("walk: %v", walkErr)
	}
}

func assertHeader(t *testing.T, path string, data []byte) {
	t.Helper()
	required := []string{
		templateHeaderSource,
		templateHeaderKind,
		templateHeaderLang,
		templateHeaderVersion,
	}
	seen := make(map[string]bool, len(required))

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		for _, prefix := range required {
			if strings.HasPrefix(line, prefix) {
				seen[prefix] = true
			}
		}
	}
	for _, prefix := range required {
		if !seen[prefix] {
			t.Errorf("%s: missing required header %q", path, prefix)
		}
	}
}

// TestSVGNotInLFSCommon locks in the deliberate exclusion: SVG is text and
// must NOT be marked as LFS in the curated lfs/common template.
func TestSVGNotInLFSCommon(t *testing.T) {
	r, err := Resolve(kindLFS, "common")
	if err != nil {
		t.Fatalf("Resolve lfs/common: %v", err)
	}
	if bytes.Contains(r.Content, []byte("*.svg")) &&
		bytes.Contains(r.Content, []byte("filter=lfs")) {
		t.Fatal("lfs/common.gitattributes must not LFS-track *.svg (SVG is text)")
	}
}
