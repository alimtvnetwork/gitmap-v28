package cmdrm

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// GetBackupBaseDir returns the primary OS temporary backup staging directory.
func GetBackupBaseDir() string {
	return filepath.Join(os.TempDir(), "gitmap_rm_backups")
}

// GetLegacyBackupBaseDir returns the legacy backup directory for backward compatibility.
func GetLegacyBackupBaseDir() string {
	return filepath.Join(os.TempDir(), "gitmap-rm")
}

// GenerateDeterministicTaskID creates a timestamped task identifier if none provided.
func GenerateDeterministicTaskID() string {
	stamp := time.Now().UTC().Format("20060102-150405")
	nonce := time.Now().UnixNano() % 10000
	return fmt.Sprintf("rm-task-%s-%04d", stamp, nonce)
}

// ResolveTaskBackupDir computes the target backup directory for a given task.
func ResolveTaskBackupDir(taskID string) string {
	timestamp := time.Now().UTC().Format("20060102_150405")
	return filepath.Join(GetBackupBaseDir(), taskID, timestamp)
}

// IsProtectedPath checks if a path resides in critical repository management folders.
func IsProtectedPath(cleanRel string) bool {
	normalized := filepath.ToSlash(cleanRel)
	if normalized == ".git" || strings.HasPrefix(normalized, ".git/") {
		return true
	}
	if normalized == ".gitmap" || strings.HasPrefix(normalized, ".gitmap/") {
		return true
	}
	if normalized == "node_modules" || strings.HasPrefix(normalized, "node_modules/") {
		return true
	}
	return false
}

// CalculateSHA256 returns the hexadecimal SHA-256 checksum of a file.
func CalculateSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", apperror.WrapSimple(err, "rm: open file for hash")
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", apperror.WrapSimple(err, "rm: calculate hash")
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// CopyFilePreserve copies a file to destination while preserving mode permissions.
func CopyFilePreserve(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return apperror.WrapSimple(err, "rm: stat source file")
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return apperror.WrapSimple(err, "rm: create destination dir")
	}

	in, err := os.Open(src)
	if err != nil {
		return apperror.WrapSimple(err, "rm: open source for copy")
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return apperror.WrapSimple(err, "rm: open destination for copy")
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return apperror.WrapSimple(err, "rm: copy bytes")
	}

	return os.Chmod(dst, info.Mode())
}

// WriteManifest atomically serializes the removal manifest to disk.
func WriteManifest(dir string, manifest *RmManifest) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return apperror.WrapSimple(err, "rm: create manifest dir")
	}

	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "rm: marshal manifest")
	}

	tmpFile := filepath.Join(dir, "manifest.json.tmp")
	destFile := filepath.Join(dir, "manifest.json")

	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return apperror.WrapSimple(err, "rm: write manifest tmp")
	}

	if err := replaceFileWithFallback(tmpFile, destFile); err != nil {
		return apperror.WrapSimple(err, "rm: rename manifest")
	}

	return nil
}

func replaceFileWithFallback(tmpFile, destFile string) error {
	err := os.Rename(tmpFile, destFile)
	if err == nil {
		return nil
	}
	_ = os.Remove(destFile)
	return os.Rename(tmpFile, destFile)
}

// ReadManifest reads a removal manifest from a session directory or file path.
func ReadManifest(manifestPathOrDir string) (*RmManifest, error) {
	manifestPath := manifestPathOrDir
	info, err := os.Stat(manifestPathOrDir)
	if err != nil {
		return nil, apperror.WrapSimple(err, "rm: stat manifest path")
	}

	if info.IsDir() {
		manifestPath = filepath.Join(manifestPathOrDir, "manifest.json")
	}

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "rm: read manifest file")
	}

	var manifest RmManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, apperror.WrapSimple(err, "rm: unmarshal manifest")
	}

	return &manifest, nil
}

// ExpandFilePatterns expands globs, directories, and paths into clean relative files.
func ExpandFilePatterns(patterns []string, repoRoot string) ([]string, error) {
	root := repoRoot
	if root == "" {
		root = "."
	}

	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, apperror.WrapSimple(err, "rm: abs repo root")
	}

	fileMap := make(map[string]bool)
	for _, pattern := range patterns {
		cleanPat := strings.TrimSpace(pattern)
		if cleanPat == "" {
			continue
		}
		if err := collectPatternFiles(absRoot, cleanPat, fileMap); err != nil {
			return nil, err
		}
	}

	results := make([]string, 0, len(fileMap))
	for f := range fileMap {
		results = append(results, f)
	}
	sort.Strings(results)
	return results, nil
}

func collectPatternFiles(absRoot, pattern string, fileMap map[string]bool) error {
	hasWildcard := strings.ContainsAny(pattern, "*?[")
	if !hasWildcard {
		return collectExactPattern(absRoot, pattern, fileMap)
	}
	return collectWildcardPattern(absRoot, pattern, fileMap)
}

func collectExactPattern(absRoot, pattern string, fileMap map[string]bool) error {
	targetPath := pattern
	if !filepath.IsAbs(targetPath) {
		targetPath = filepath.Join(absRoot, pattern)
	}

	info, err := os.Stat(targetPath)
	if err != nil {
		return nil
	}

	if info.IsDir() {
		return walkDirectoryFiles(absRoot, targetPath, fileMap)
	}

	rel, err := filepath.Rel(absRoot, targetPath)
	if err != nil || IsProtectedPath(rel) {
		return nil
	}

	fileMap[filepath.ToSlash(rel)] = true
	return nil
}

func collectWildcardPattern(absRoot, pattern string, fileMap map[string]bool) error {
	return filepath.Walk(absRoot, func(current string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info == nil {
			return nil
		}

		rel, relErr := filepath.Rel(absRoot, current)
		if relErr != nil || rel == "." {
			return nil
		}

		if info.IsDir() && IsProtectedPath(rel) {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return nil
		}

		if isPatternMatch(pattern, rel, filepath.Base(current)) {
			fileMap[filepath.ToSlash(rel)] = true
		}
		return nil
	})
}

func walkDirectoryFiles(absRoot, dirPath string, fileMap map[string]bool) error {
	return filepath.Walk(dirPath, func(current string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info == nil {
			return nil
		}

		rel, relErr := filepath.Rel(absRoot, current)
		if relErr != nil || rel == "." {
			return nil
		}

		if info.IsDir() && IsProtectedPath(rel) {
			return filepath.SkipDir
		}

		if !info.IsDir() && !IsProtectedPath(rel) {
			fileMap[filepath.ToSlash(rel)] = true
		}
		return nil
	})
}

func isPatternMatch(pattern, relPath, baseName string) bool {
	normPat := filepath.ToSlash(pattern)
	normRel := filepath.ToSlash(relPath)

	if matched, _ := filepath.Match(normPat, normRel); matched {
		return true
	}
	if matched, _ := filepath.Match(normPat, baseName); matched {
		return true
	}
	return false
}

// StageAndRemoveFiles executes safe staging and removal of matched files.
func StageAndRemoveFiles(opts RmOptions) (*RmManifest, error) {
	root := opts.RepoRoot
	if root == "" {
		root = "."
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, apperror.WrapSimple(err, "rm: resolve abs root")
	}

	matchedFiles, err := ExpandFilePatterns(opts.Patterns, absRoot)
	if err != nil {
		return nil, err
	}

	if len(matchedFiles) == 0 {
		return nil, apperror.NewValidationError("no files matched specified removal pattern")
	}

	taskID := opts.TaskId
	if taskID == "" {
		taskID = GenerateDeterministicTaskID()
	}

	sessionDir := ResolveTaskBackupDir(taskID)
	manifest := &RmManifest{
		TaskId:    taskID,
		CreatedAt: time.Now().UTC(),
		RepoRoot:  filepath.ToSlash(absRoot),
		Reason:    opts.Reason,
		Files:     make([]RmFileRecord, 0, len(matchedFiles)),
	}

	if err := processStagingRecords(absRoot, sessionDir, matchedFiles, opts.DryRun, manifest); err != nil {
		return nil, err
	}

	manifest.FileCount = len(manifest.Files)
	if opts.DryRun {
		return manifest, nil
	}

	if err := WriteManifest(sessionDir, manifest); err != nil {
		return nil, err
	}

	return manifest, nil
}

func processStagingRecords(absRoot, sessionDir string, files []string, isDryRun bool, manifest *RmManifest) error {
	for _, relPath := range files {
		record, err := stageSingleFile(absRoot, sessionDir, relPath, isDryRun)
		if err != nil {
			return err
		}
		manifest.Files = append(manifest.Files, *record)
		manifest.TotalSizeBytes += record.SizeBytes
	}
	return nil
}

func stageSingleFile(absRoot, sessionDir, relPath string, isDryRun bool) (*RmFileRecord, error) {
	srcPath := filepath.Join(absRoot, filepath.FromSlash(relPath))
	info, err := os.Stat(srcPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "rm: stat file for staging")
	}

	sha, err := CalculateSHA256(srcPath)
	if err != nil {
		return nil, err
	}

	stagedPath := filepath.Join(sessionDir, "files", filepath.FromSlash(relPath))
	if err := executeStageCopyAndRemove(srcPath, stagedPath, isDryRun); err != nil {
		return nil, err
	}

	return &RmFileRecord{
		RelativePath: relPath,
		StagedPath:   filepath.ToSlash(stagedPath),
		SizeBytes:    info.Size(),
		Sha256:       sha,
		FileMode:     info.Mode(),
		DeletedAt:    time.Now().UTC(),
	}, nil
}

func executeStageCopyAndRemove(srcPath, stagedPath string, isDryRun bool) error {
	if isDryRun {
		return nil
	}
	if err := CopyFilePreserve(srcPath, stagedPath); err != nil {
		return err
	}
	if err := os.Remove(srcPath); err != nil {
		return apperror.WrapSimple(err, "rm: remove file from workspace")
	}
	return nil
}

// ListRmBackups scans all stored backup sessions across OS temp storage.
func ListRmBackups() ([]RmManifest, error) {
	var manifests []RmManifest
	baseDirs := []string{GetBackupBaseDir(), GetLegacyBackupBaseDir()}

	for _, base := range baseDirs {
		found := scanBaseDirManifests(base)
		manifests = append(manifests, found...)
	}

	sort.Slice(manifests, func(i, j int) bool {
		return manifests[i].CreatedAt.After(manifests[j].CreatedAt)
	})

	return manifests, nil
}

func scanBaseDirManifests(base string) []RmManifest {
	var manifests []RmManifest
	if info, err := os.Stat(base); err != nil || !info.IsDir() {
		return manifests
	}

	_ = filepath.Walk(base, func(current string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() {
			return nil
		}
		if info.Name() != "manifest.json" {
			return nil
		}
		m, readErr := ReadManifest(current)
		if readErr == nil {
			manifests = append(manifests, *m)
		}
		return nil
	})

	return manifests
}

// PurgeRmBackups deletes backup directories based on options and returns count of removed sessions.
func PurgeRmBackups(opts RmPurgeOptions) (int, error) {
	manifests, err := ListRmBackups()
	if err != nil {
		return 0, err
	}

	purgedCount := 0
	cutoff := time.Now().UTC().Add(-opts.OlderThan)

	for _, m := range manifests {
		shouldPurge := resolvePurgeTarget(m, opts, cutoff)
		if !shouldPurge {
			continue
		}
		if err := purgeSessionDirs(m.TaskId); err == nil {
			purgedCount++
		}
	}

	return purgedCount, nil
}

func resolvePurgeTarget(m RmManifest, opts RmPurgeOptions, cutoff time.Time) bool {
	if opts.All {
		return true
	}
	if opts.TaskId != "" && m.TaskId == opts.TaskId {
		return true
	}
	if opts.OlderThan > 0 && m.CreatedAt.Before(cutoff) {
		return true
	}
	return false
}

func purgeSessionDirs(taskID string) error {
	baseDirs := []string{GetBackupBaseDir(), GetLegacyBackupBaseDir()}
	for _, base := range baseDirs {
		target := filepath.Join(base, taskID)
		if info, err := os.Stat(target); err == nil && info.IsDir() {
			_ = os.RemoveAll(target)
		}
	}
	return nil
}
