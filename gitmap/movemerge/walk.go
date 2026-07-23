package movemerge

import (
	"crypto/sha256"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FileMeta is a single file's identity used by the diff stage.
type FileMeta struct {
	RelPath string
	Info    os.FileInfo
	SHA     string
}

// IndexTree walks root and returns rel-path -> FileMeta for every
// non-ignored regular file. Symlinks are recorded but not followed.
func IndexTree(root string, opts Options) (map[string]FileMeta, error) {
	out := make(map[string]FileMeta)
	walkErr := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil || rel == "." {
			return relErr
		}
		if shouldSkipWalk(rel, opts) {
			if info.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}
		if info.IsDir() {
			return nil
		}
		out[filepath.ToSlash(rel)] = FileMeta{RelPath: filepath.ToSlash(rel), Info: info}

		return nil
	})

	return out, walkErr
}

// shouldSkipWalk applies the default ignore list (.git/, node_modules/,
// .gitmap/release-assets/) honoring the include-* opt-ins.
func shouldSkipWalk(rel string, opts Options) bool {
	base := filepath.Base(rel)
	if !opts.IncludeVCS && base == ".git" {
		return true
	}
	if !opts.IncludeNodeMods && base == "node_modules" {
		return true
	}

	return strings.HasPrefix(filepath.ToSlash(rel), ".gitmap/release-assets/")
}

// SortedKeys returns the union of two map keysets, sorted ascending.
func SortedKeys(a, b map[string]FileMeta) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sort.Strings(out)

	return out
}

// HashFile streams the file at path through SHA-256.
func HashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return "", err
	}

	return string(h.Sum(nil)), nil
}
