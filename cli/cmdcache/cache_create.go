package cmdcache

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

const maxCacheFileSize = 200 * 1024 // 200 KB

// CreateCache indexes the specified targets into the Split-DB cache.
func CreateCache(targets []string) *appfault.AppError {
	repoRoot := findRepoRoot()
	rootDB, err := store.OpenRootCacheDB(repoRoot)
	if err != nil {
		return err
	}
	defer rootDB.Close()

	if len(targets) == 0 {
		targets = []string{"."}
	}

	processTargets(rootDB, targets, repoRoot)
	return nil
}

func processTargets(rootDB *sql.DB, targets []string, repoRoot string) {
	totalIndexed := 0
	for _, t := range targets {
		indexed := processSingleTarget(rootDB, t, repoRoot)
		totalIndexed += indexed
	}
	fmt.Printf("\n%s✓ Split-DB cache created: indexed %d file(s) (<= 200KB).%s\n",
		constants.ColorGreen, totalIndexed, constants.ColorReset)
}

func processSingleTarget(rootDB *sql.DB, target, repoRoot string) int {
	absTarget, _ := filepath.Abs(target)
	info, err := os.Stat(absTarget)
	if err != nil {
		return 0
	}
	if info.IsDir() {
		return walkAndIndexDir(rootDB, absTarget, repoRoot)
	}
	return indexSingleFileIfEligible(rootDB, absTarget, repoRoot)
}

func walkAndIndexDir(rootDB *sql.DB, targetDir, repoRoot string) int {
	count := 0
	_ = filepath.WalkDir(targetDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil || isSkippedPath(path, d) {
			return nil
		}
		if !d.IsDir() {
			count += indexSingleFileIfEligible(rootDB, path, repoRoot)
		}
		return nil
	})
	return count
}

func isSkippedPath(path string, d fs.DirEntry) bool {
	name := d.Name()
	isSkip := name == ".git" || name == "node_modules" || name == ".gitmap"
	return d.IsDir() && isSkip
}

func isFileEligible(info fs.FileInfo, absPath string) bool {
	if info.IsDir() || info.Size() > maxCacheFileSize {
		return false
	}
	return !isBinaryFile(absPath)
}

func indexSingleFileIfEligible(rootDB *sql.DB, absPath, repoRoot string) int {
	info, err := os.Stat(absPath)
	if err != nil || !isFileEligible(info, absPath) {
		return 0
	}

	relPath, _ := filepath.Rel(repoRoot, absPath)
	cleanRel := filepath.ToSlash(relPath)
	slug := resolveFileSlug(cleanRel)

	rec := store.CacheFileRecord{
		RelativePath: cleanRel,
		AbsolutePath: absPath,
		FileSize:     info.Size(),
		ModifiedTime: info.ModTime().Unix(),
		FolderSlug:   slug,
	}

	if insertErr := store.InsertCacheFile(rootDB, rec); insertErr != nil {
		return 0
	}
	indexFileContent(slug, cleanRel, absPath, repoRoot)
	return 1
}

func indexFileContent(slug, relPath, absPath, repoRoot string) {
	contentBytes, err := os.ReadFile(absPath)
	if err != nil {
		return
	}
	slugDB, openErr := store.OpenSlugCacheDB(slug, repoRoot)
	if openErr != nil {
		return
	}
	defer slugDB.Close()

	lines := strings.Split(string(contentBytes), "\n")
	_ = store.InsertCachedLines(slugDB, relPath, lines)
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
