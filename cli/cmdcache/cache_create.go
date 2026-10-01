package cmdcache

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

const maxCacheFileSize = 200 * 1024 // 200 KB

// CreateCache indexes the specified targets into the Split-DB cache.
func CreateCache(targets []string) error {
	repoRoot := findRepoRoot()
	rootDB, err := OpenRootCacheDB(repoRoot)
	if err != nil {
		return err
	}
	defer rootDB.Close()

	if len(targets) == 0 {
		targets = []string{"."}
	}
	totalIndexed := 0
	for _, t := range targets {
		indexed, _ := processSingleTarget(rootDB, t, repoRoot)
		totalIndexed += indexed
	}

	fmt.Printf("\n%s✓ Split-DB cache created: indexed %d file(s) (<= 200KB).%s\n",
		constants.ColorGreen, totalIndexed, constants.ColorReset)
	return nil
}

func processSingleTarget(rootDB *sql.DB, target, repoRoot string) (int, error) {
	absTarget, _ := filepath.Abs(target)
	info, err := os.Stat(absTarget)
	if err != nil {
		return 0, err
	}
	if info.IsDir() {
		return walkAndIndexDir(rootDB, absTarget, repoRoot)
	}
	return indexSingleFileIfEligible(rootDB, absTarget, repoRoot)
}

func walkAndIndexDir(rootDB *sql.DB, targetDir, repoRoot string) (int, error) {
	count := 0
	err := filepath.WalkDir(targetDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil || isSkippedPath(path, d) {
			return nil
		}
		if !d.IsDir() {
			n, _ := indexSingleFileIfEligible(rootDB, path, repoRoot)
			count += n
		}
		return nil
	})
	return count, err
}

func isSkippedPath(path string, d fs.DirEntry) bool {
	name := d.Name()
	return d.IsDir() && (name == ".git" || name == "node_modules" || name == ".gitmap")
}

func indexSingleFileIfEligible(rootDB *sql.DB, absPath, repoRoot string) (int, error) {
	info, err := os.Stat(absPath)
	if err != nil || info.IsDir() || info.Size() > maxCacheFileSize {
		return 0, nil
	}
	if isBinaryFile(absPath) {
		return 0, nil
	}

	relPath, _ := filepath.Rel(repoRoot, absPath)
	cleanRel := filepath.ToSlash(relPath)
	slug := resolveFileSlug(cleanRel)

	rec := CacheFileRecord{
		RelativePath: cleanRel,
		AbsolutePath: absPath,
		FileSize:     info.Size(),
		ModifiedTime: info.ModTime().Unix(),
		FolderSlug:   slug,
	}
	if insertErr := InsertCacheFile(rootDB, rec); insertErr != nil {
		return 0, insertErr
	}
	indexFileContent(slug, cleanRel, absPath, repoRoot)
	return 1, nil
}

func indexFileContent(slug, relPath, absPath, repoRoot string) {
	contentBytes, err := os.ReadFile(absPath)
	if err != nil {
		return
	}
	slugDB, openErr := OpenSlugCacheDB(slug, repoRoot)
	if openErr != nil {
		return
	}
	defer slugDB.Close()

	lines := strings.Split(string(contentBytes), "\n")
	_ = InsertCachedLines(slugDB, relPath, lines)
}

func isBinaryFile(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return true
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, _ := f.Read(buf)
	return bytes.Contains(buf[:n], []byte{0})
}

func resolveFileSlug(relPath string) string {
	parts := strings.Split(relPath, "/")
	if len(parts) > 1 {
		return store.SanitizeSlug(parts[0])
	}
	return "root"
}

func findRepoRoot() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return cwd
}
