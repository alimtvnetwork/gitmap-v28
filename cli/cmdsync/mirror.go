package cmdsync

import (
	"bytes"
	"crypto/sha256"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

type MirrorStats struct {
	Added   int
	Updated int
	Removed int
}

// MirrorAssets synchronizes canonical assets from source to target repo adhering to boundaries.
func MirrorAssets(sourceRoot, targetRoot string) (MirrorStats, error) {
	var stats MirrorStats

	dirs := []struct {
		relSrc string
		relDst string
		addOnly bool
	}{
		{"01-prompts", "01-prompts", false},
		{filepath.Join(".agents", "skills"), filepath.Join(".agents", "skills"), false},
		{filepath.Join(".cursor", "skills"), filepath.Join(".cursor", "skills"), false},
		{"03-ai-scripts", "03-ai-scripts", true},
		{filepath.Join(".agents", "scripts"), filepath.Join(".agents", "scripts"), true},
	}

	specRoot := filepath.Join(sourceRoot, "02-spec")
	if info, err := os.Stat(specRoot); err == nil && info.IsDir() {
		entries, _ := os.ReadDir(specRoot)
		specNumRegex := regexp.MustCompile(`^(\d+)-`)
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			m := specNumRegex.FindStringSubmatch(e.Name())
			if len(m) > 1 {
				num, _ := strconv.Atoi(m[1])
				if num >= 1 && num <= 20 {
					rel := filepath.Join("02-spec", e.Name())
					dirs = append(dirs, struct {
						relSrc string
						relDst string
						addOnly bool
					}{rel, rel, false})
				}
			}
		}
	}

	for _, d := range dirs {
		srcDir := filepath.Join(sourceRoot, d.relSrc)
		dstDir := filepath.Join(targetRoot, d.relDst)

		if _, err := os.Stat(srcDir); err != nil {
			continue
		}

		s, err := syncDirectory(srcDir, dstDir, sourceRoot, targetRoot, d.addOnly)
		if err != nil {
			return stats, err
		}
		stats.Added += s.Added
		stats.Updated += s.Updated
		stats.Removed += s.Removed
	}

	return stats, nil
}

func syncDirectory(srcDir, dstDir, sourceRoot, targetRoot string, addOnly bool) (MirrorStats, error) {
	var stats MirrorStats

	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		relToSrc, _ := filepath.Rel(srcDir, path)
		if relToSrc == "." || relToSrc == "" {
			return nil
		}

		relFromRoot, _ := filepath.Rel(sourceRoot, path)
		normRel := filepath.ToSlash(relFromRoot)

		if IsSpec21OrHigher(normRel) || IsBumpScript(info.Name()) ||
			IsMemoryOrPlans(normRel) || IsArchive(normRel) || IsSecret(info.Name()) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		dstPath := filepath.Join(dstDir, relToSrc)

		if info.IsDir() {
			_ = os.MkdirAll(dstPath, 0755)
			return nil
		}

		// File handling
		dstInfo, dstErr := os.Stat(dstPath)
		if dstErr != nil {
			// Destination file missing -> Add
			if copyErr := copyFile(path, dstPath); copyErr == nil {
				stats.Added++
			}
			return nil
		}

		// Destination file exists
		if addOnly {
			// Additive-only mode: Never overwrite existing files in target
			return nil
		}

		// Compare content
		isEqual, cmpErr := filesEqual(path, dstPath)
		if cmpErr == nil && !isEqual {
			if copyErr := copyFile(path, dstPath); copyErr == nil {
				stats.Updated++
			}
		}
		_ = dstInfo

		return nil
	})

	return stats, err
}

func filesEqual(p1, p2 string) (bool, error) {
	h1, err1 := hashFile(p1)
	if err1 != nil {
		return false, err1
	}
	h2, err2 := hashFile(p2)
	if err2 != nil {
		return false, err2
	}
	return bytes.Equal(h1, h2), nil
}

func hashFile(p string) ([]byte, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func copyFile(src, dst string) error {
	_ = os.MkdirAll(filepath.Dir(dst), 0755)

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
