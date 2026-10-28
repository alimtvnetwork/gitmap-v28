package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// SQLCreateRepoCacheSequenceIndex creates the sequence folder lookup cache in the root database.
const SQLCreateRepoCacheSequenceIndex = `
CREATE TABLE IF NOT EXISTS RepoCacheSequenceIndex (
    RepoSlug       TEXT PRIMARY KEY,
    SequenceFolder TEXT NOT NULL,
    CachedPath     TEXT NOT NULL,
    UpdatedAt      TEXT NOT NULL
);`

const sqlCreateIdxRepoCacheSeqFolder = `
CREATE INDEX IF NOT EXISTS IdxRepoCacheSequenceIndex_Folder ON RepoCacheSequenceIndex(SequenceFolder);`

func openRootDatabase() (*sql.DB, error) {
	rootDbPath := filepath.Join(BinaryDataDir(), "gitmap.db")
	wrap := OpenSQLiteDB(rootDbPath)
	if wrap.IsFailure() {
		return nil, wrap.AppError()
	}

	return initAndReturnDb(wrap.Value)
}

func initAndReturnDb(db *sql.DB) (*sql.DB, error) {
	initErr := initRepoCacheSequenceTable(db)
	if initErr != nil {
		_ = db.Close()

		return nil, initErr
	}

	return db, nil
}

func initRepoCacheSequenceTable(db *sql.DB) error {
	_, err := db.Exec(SQLCreateRepoCacheSequenceIndex)
	if err != nil {
		return apperror.WrapSimple(err, "init repo cache sequence table")
	}

	_, _ = db.Exec(sqlCreateIdxRepoCacheSeqFolder) // lint-allow: ignore-db-error

	return nil
}

func lookupCachedSequenceFolder(db *sql.DB, repoSlug string) (string, bool) {
	hasNilDb := db == nil
	if hasNilDb {
		return "", false
	}

	cachedPath := queryCachedPath(db, repoSlug)
	hasCached := len(cachedPath) > 0
	if !hasCached {
		return "", false
	}

	return verifyCachedPathExists(cachedPath)
}

func queryCachedPath(db *sql.DB, repoSlug string) string {
	var cachedPath string
	query := "SELECT CachedPath FROM RepoCacheSequenceIndex WHERE RepoSlug = ? LIMIT 1;"
	_ = db.QueryRow(query, repoSlug).Scan(&cachedPath) // lint-allow: ignore-db-error

	return cachedPath
}

func verifyCachedPathExists(cachedPath string) (string, bool) {
	info, statErr := os.Stat(cachedPath)
	isDir := statErr == nil && info.IsDir()
	if isDir {
		return filepath.ToSlash(cachedPath), true
	}

	return "", false
}

const sqlUpsertRepoCacheSequence = `
INSERT INTO RepoCacheSequenceIndex (RepoSlug, SequenceFolder, CachedPath, UpdatedAt)
VALUES (?, ?, ?, ?)
ON CONFLICT(RepoSlug) DO UPDATE SET
    SequenceFolder = excluded.SequenceFolder,
    CachedPath = excluded.CachedPath,
    UpdatedAt = excluded.UpdatedAt;`

func recordCachedSequenceFolder(db *sql.DB, repoSlug, seqFolder, cachedPath string) error {
	hasNilDb := db == nil
	if hasNilDb {
		return nil
	}

	nowStr := time.Now().UTC().Format(time.RFC3339)
	_, err := db.Exec(sqlUpsertRepoCacheSequence, repoSlug, seqFolder, cachedPath, nowStr)
	if err != nil {
		return apperror.WrapSimple(err, "record repo cache sequence")
	}

	return nil
}

func scanExistingSequenceFolder(baseDir, repoSlug string) (string, bool) {
	entries, err := os.ReadDir(baseDir)
	hasErr := err != nil
	if hasErr {
		return "", false
	}

	return findMatchingSequenceEntry(baseDir, entries, "-"+strings.ToLower(repoSlug))
}

func findMatchingSequenceEntry(baseDir string, entries []os.DirEntry, suffix string) (string, bool) {
	for _, entry := range entries {
		foundPath, hasFound := matchSequenceFolder(baseDir, entry, suffix)
		if hasFound {
			return foundPath, true
		}
	}

	return "", false
}

func matchSequenceFolder(baseDir string, entry os.DirEntry, suffix string) (string, bool) {
	isDir := entry.IsDir()
	if !isDir {
		return "", false
	}

	nameLower := strings.ToLower(entry.Name())
	hasSuffix := strings.HasSuffix(nameLower, suffix)
	if hasSuffix {
		return filepath.ToSlash(filepath.Join(baseDir, entry.Name())), true
	}

	return "", false
}

func determineNextSequenceNumber(baseDir string) int {
	entries, err := os.ReadDir(baseDir)
	hasErr := err != nil
	if hasErr {
		return 1
	}

	return computeMaxSequence(entries) + 1
}

func computeMaxSequence(entries []os.DirEntry) int {
	maxSeq := 0
	for _, entry := range entries {
		seq := extractPrefixSequenceNumber(entry)
		hasHigher := seq > maxSeq
		if hasHigher {
			maxSeq = seq
		}
	}

	return maxSeq
}

func extractPrefixSequenceNumber(entry os.DirEntry) int {
	isDir := entry.IsDir()
	if !isDir {
		return 0
	}

	return parseDirPrefixDigits(entry.Name())
}

func parseDirPrefixDigits(name string) int {
	parts := strings.SplitN(name, "-", 2)
	hasPrefix := len(parts) == 2
	if !hasPrefix {
		return 0
	}

	num, numErr := strconv.Atoi(parts[0])
	isValid := numErr == nil && num > 0
	if isValid {
		return num
	}

	return 0
}

func createSequenceFolder(baseDir, repoSlug string, seqNum int) (string, string, error) {
	folderName := fmt.Sprintf("%02d-%s", seqNum, repoSlug)
	targetPath := filepath.ToSlash(filepath.Join(baseDir, folderName))
	mkdirErr := os.MkdirAll(targetPath, 0755)
	if mkdirErr != nil {
		return "", "", apperror.WrapSimple(mkdirErr, "create sequence folder")
	}

	return folderName, targetPath, nil
}

func cleanRepoSlug(repoSlug string) string {
	clean := strings.TrimSpace(repoSlug)
	clean = strings.TrimSuffix(clean, ".git")
	clean = strings.ToLower(clean)

	return filepath.Base(clean)
}

func resolveSequenceBaseDir(baseDir string) string {
	hasEmpty := len(strings.TrimSpace(baseDir)) == 0
	if hasEmpty {
		return ResolveRepoCacheRoot()
	}

	return baseDir
}

// ResolveRepoCacheSequenceFolder locates or creates the sequence directory for a repo in repo-cache.
func ResolveRepoCacheSequenceFolder(repoSlug string, repoCacheBaseDir string) (string, error) {
	slug := cleanRepoSlug(repoSlug)
	hasEmptySlug := len(slug) == 0
	if hasEmptySlug {
		return "", apperror.NewValidationError("repo slug is empty")
	}

	baseDir := resolveSequenceBaseDir(repoCacheBaseDir)

	return resolveAndCacheSequenceFolder(baseDir, slug)
}

func resolveAndCacheSequenceFolder(baseDir, slug string) (string, error) {
	db, _ := openRootDatabase()
	if db != nil {
		defer db.Close()
	}

	return resolveSequenceWithFallback(db, baseDir, slug)
}

func resolveSequenceWithFallback(db *sql.DB, baseDir, slug string) (string, error) {
	cachedPath, hasCache := lookupCachedSequenceFolder(db, slug)
	if hasCache {
		return cachedPath, nil
	}

	diskPath, hasDisk := scanExistingSequenceFolder(baseDir, slug)
	if hasDisk {
		_ = recordCachedSequenceFolder(db, slug, filepath.Base(diskPath), diskPath)

		return diskPath, nil
	}

	return provisionNewSequenceFolder(db, baseDir, slug)
}

func provisionNewSequenceFolder(db *sql.DB, baseDir, slug string) (string, error) {
	nextSeq := determineNextSequenceNumber(baseDir)
	folderName, targetPath, err := createSequenceFolder(baseDir, slug, nextSeq)
	if err != nil {
		return "", err
	}

	_ = recordCachedSequenceFolder(db, slug, folderName, targetPath)

	return targetPath, nil
}

// ResolveRepoConfusionDbPath returns the canonical path to ai-confusion.db inside the repo sequence folder.
func ResolveRepoConfusionDbPath(repoSlug string) (string, error) {
	seqFolder, err := ResolveRepoCacheSequenceFolder(repoSlug, "")
	if err != nil {
		return "", err
	}

	return filepath.ToSlash(filepath.Join(seqFolder, "ai-confusion.db")), nil
}
