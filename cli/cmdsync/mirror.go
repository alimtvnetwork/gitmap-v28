package cmdsync

import (
	"bytes"
	"crypto/sha256"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsync/constants"
)

var specDirNumRe = regexp.MustCompile(`^(\d+)-`)

// MirrorStats records added, updated, and removed file counts during asset mirroring.
type MirrorStats struct {
	Added   int
	Updated int
	Removed int
}

type syncDirMapping struct {
	relSrc  string
	relDst  string
	addOnly bool
}

// MirrorAssets synchronizes canonical assets from source to target repo adhering to boundaries.
func MirrorAssets(sourceRoot, targetRoot string) (MirrorStats, error) {
	var stats MirrorStats
	dirs := buildSyncDirMappings(sourceRoot)

	for _, d := range dirs {
		if err := mirrorSingleDirectoryMapping(d, sourceRoot, targetRoot, &stats); err != nil {
			return stats, err
		}
	}

	return stats, nil
}

func buildSyncDirMappings(sourceRoot string) []syncDirMapping {
	dirs := []syncDirMapping{
		{constants.DirPrompts, constants.DirPrompts, false},
		{filepath.FromSlash(constants.DirAgentsSkills), filepath.FromSlash(constants.DirAgentsSkills), false},
		{filepath.FromSlash(constants.DirCursorSkills), filepath.FromSlash(constants.DirCursorSkills), false},
		{constants.DirAiScripts, constants.DirAiScripts, true},
		{filepath.FromSlash(constants.DirAgentsScripts), filepath.FromSlash(constants.DirAgentsScripts), true},
	}
	specRoot := filepath.Join(sourceRoot, constants.DirSpec)

	return appendSpecDirs(specRoot, dirs)
}

func mirrorSingleDirectoryMapping(d syncDirMapping, sourceRoot, targetRoot string, stats *MirrorStats) error {
	srcDir := filepath.Join(sourceRoot, d.relSrc)
	dstDir := filepath.Join(targetRoot, d.relDst)

	if _, err := os.Stat(srcDir); err != nil {
		return nil
	}

	s, err := syncDirectory(srcDir, dstDir, sourceRoot, targetRoot, d.addOnly)
	if err != nil {
		return err
	}

	stats.Added += s.Added
	stats.Updated += s.Updated
	stats.Removed += s.Removed

	return nil
}

func syncDirectory(srcDir, dstDir, sourceRoot, targetRoot string, addOnly bool) (MirrorStats, error) {
	var stats MirrorStats
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return nil
		}

		return processWalkItem(path, info, srcDir, dstDir, sourceRoot, addOnly, &stats)
	})

	return stats, err
}

func processWalkItem(path string, info os.FileInfo, srcDir, dstDir, sourceRoot string, addOnly bool, stats *MirrorStats) error {
	relToSrc, _ := filepath.Rel(srcDir, path)
	if relToSrc == "." || relToSrc == "" {
		return nil
	}

	relFromRoot, _ := filepath.Rel(sourceRoot, path)
	normRel := filepath.ToSlash(relFromRoot)

	if isExcludedSyncItem(normRel, info.Name()) {
		return handleExcludedItem(info)
	}

	dstPath := filepath.Join(dstDir, relToSrc)
	if info.IsDir() {
		_ = os.MkdirAll(dstPath, 0755)

		return nil
	}

	return syncSingleFile(path, dstPath, addOnly, stats)
}

func isExcludedSyncItem(normRel, name string) bool {
	return IsSpec21OrHigher(normRel) || IsBumpScript(name) ||
		IsMemoryOrPlans(normRel) || IsArchive(normRel) || IsSecret(name)
}

func handleExcludedItem(info os.FileInfo) error {
	if info.IsDir() {
		return filepath.SkipDir
	}

	return nil
}

func syncSingleFile(srcPath, dstPath string, addOnly bool, stats *MirrorStats) error {
	if _, dstErr := os.Stat(dstPath); dstErr != nil {
		recordFileCopy(srcPath, dstPath, &stats.Added)

		return nil
	}
	if addOnly {
		return nil
	}

	isEqual, cmpErr := filesEqual(srcPath, dstPath)
	if cmpErr == nil && !isEqual {
		recordFileCopy(srcPath, dstPath, &stats.Updated)
	}

	return nil
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

func recordFileCopy(src, dst string, count *int) {
	if copyErr := copyFile(src, dst); copyErr == nil {
		*count++
	}
}

func appendSpecDirs(specRoot string, dirs []syncDirMapping) []syncDirMapping {
	info, err := os.Stat(specRoot)
	if err != nil || !info.IsDir() {
		return dirs
	}

	entries, _ := os.ReadDir(specRoot)
	for _, e := range entries {
		dirs = appendMatchingSpecDir(e, dirs)
	}

	return dirs
}

func appendMatchingSpecDir(e os.DirEntry, dirs []syncDirMapping) []syncDirMapping {
	if !e.IsDir() {
		return dirs
	}

	m := specDirNumRe.FindStringSubmatch(e.Name())
	if len(m) <= 1 {
		return dirs
	}

	num, _ := strconv.Atoi(m[1])
	if num < 1 || num > 20 {
		return dirs
	}

	rel := filepath.Join(constants.DirSpec, e.Name())

	return append(dirs, syncDirMapping{relSrc: rel, relDst: rel, addOnly: false})
}
