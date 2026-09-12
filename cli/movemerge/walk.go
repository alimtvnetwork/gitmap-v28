package movemerge

import (
	"crypto/sha256"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// IndexTree walks root and returns rel-path -> FileMeta for every
// non-ignored regular file. Symlinks are recorded but not followed.
func IndexTree(root string, opts Options) FileMetaMapResult {
	out := make(map[string]FileMeta)
	walkErr := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, relErr := filepath.Rel(root, path)
		if relErr != nil || rel == "." {
			return relErr
		}

		if IsSkipWalk(rel, opts) && info.IsDir() {
			return filepath.SkipDir
		}

		if IsSkipWalk(rel, opts) {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		out[filepath.ToSlash(rel)] = FileMeta{RelPath: filepath.ToSlash(rel), Info: info}

		return nil
	})
	if walkErr != nil {
		appErr := apperror.WrapSimple(walkErr, "walk root index tree")

		return result.FailMap[string, FileMeta](appErr)
	}

	return result.OkMap(out)
}

// IsSkipWalk applies the default ignore list (.git/, node_modules/,
// .gitmap/release-assets/) honoring the include-* opt-ins.
func IsSkipWalk(rel string, opts Options) bool {
	base := filepath.Base(rel)
	if base == ".git" {
		return !opts.IsIncludeVCS
	}

	if base == "node_modules" {
		return !opts.IsIncludeNodeMods
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
