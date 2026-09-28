// Package cmd — special_repos_ops.go provides filesystem sequencing, git auto-commit/push, and first-scan check operations for repo-secrets (rs) and repo-cache (rc).
package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// SpecialPutResult captures the outcome of storing a file, folder, or text note in rs/rc.
type SpecialPutResult struct {
	SpecialKey   string `json:"specialKey"`
	SpecialRoot  string `json:"specialRoot"`
	RepoName     string `json:"repoName"`
	RepoFolder   string `json:"repoFolder"`
	TargetPath   string `json:"targetPath"`
	RelativePath string `json:"relativePath"`
	IsCommitted  bool   `json:"isCommitted"`
	IsPushed     bool   `json:"isPushed"`
}

func resolveSpecialWorkBaseDir() string {
	if workDir, hasDefault := resolveDefaultWorkDirPath(); hasDefault {
		return workDir
	}
	if info, err := os.Stat(`D:\work`); err == nil && info.IsDir() {
		return `D:\work`
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return filepath.Dir(cwd)
}

func resolveSpecialRepoRoot(keyOrShort string) (string, *store.SpecialRepositoryRecord, error) {
	db, err := store.OpenSpecialReposSplitDB()
	if err != nil {
		return "", nil, err
	}
	defer db.Close()
	return resolveSpecialRepoRootWithDB(db, keyOrShort, resolveSpecialWorkBaseDir())
}

func resolveSpecialRepoRootWithDB(db *store.SpecialReposSplitDB, keyOrShort, workBaseDir string) (string, *store.SpecialRepositoryRecord, error) {
	rec, err := db.GetSpecialRepo(keyOrShort)
	if err != nil {
		return "", nil, err
	}
	targetPath := rec.LocalPath
	if len(strings.TrimSpace(targetPath)) == 0 {
		targetPath = filepath.Join(workBaseDir, rec.ConfiguredName)
	}
	if err := ensureSpecialGitRepo(targetPath); err != nil {
		return "", nil, err
	}
	return targetPath, rec, nil
}

func ensureSpecialGitRepo(repoDir string) error {
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		return apperror.WrapSimple(err, "mkdir special repo")
	}
	gitDir := filepath.Join(repoDir, ".git")
	if info, err := os.Stat(gitDir); err == nil && info.IsDir() {
		return nil
	}
	cmd := exec.Command("git", "init")
	cmd.Dir = repoDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return apperror.WrapSimple(err, "git init special repo: "+string(out))
	}
	return nil
}

func resolveCurrentRepoName(explicitRepo string) string {
	if len(strings.TrimSpace(explicitRepo)) > 0 {
		return stripNumericPrefix(strings.TrimSpace(explicitRepo))
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "gitmap"
	}
	return stripNumericPrefix(filepath.Base(cwd))
}

func stripNumericPrefix(name string) string {
	parts := strings.SplitN(name, "-", 2)
	if len(parts) < 2 || len(parts[0]) < 2 {
		return strings.ToLower(name)
	}
	if _, err := strconv.Atoi(parts[0]); err == nil {
		return strings.ToLower(parts[1])
	}
	return strings.ToLower(name)
}

func allocateSequencedItemPath(repoFolderDir, baseName string) string {
	if hasLeadingTwoDigitSeq(baseName) {
		return filepath.Join(repoFolderDir, baseName)
	}
	nextSeq := nextItemSeqInDir(repoFolderDir)
	seqName := fmt.Sprintf("%02d-%s", nextSeq, baseName)
	return filepath.Join(repoFolderDir, seqName)
}

func hasLeadingTwoDigitSeq(name string) bool {
	parts := strings.SplitN(name, "-", 2)
	if len(parts) < 2 || len(parts[0]) < 2 {
		return false
	}
	seq, err := strconv.Atoi(parts[0])
	return err == nil && seq > 0
}

func nextItemSeqInDir(dirPath string) int {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return 1
	}
	maxVal := 0
	for _, entry := range entries {
		seq := extractEntrySeq(entry.Name())
		if seq > maxVal {
			maxVal = seq
		}
	}
	return maxVal + 1
}

func extractEntrySeq(name string) int {
	parts := strings.SplitN(name, "-", 2)
	if len(parts) < 2 || len(parts[0]) < 2 {
		return 0
	}
	seq, err := strconv.Atoi(parts[0])
	if err != nil || seq < 0 {
		return 0
	}
	return seq
}

func deriveTextSlugAndFilename(text, customSlug, customExt string) string {
	slug := sanitizeKebabSlug(customSlug)
	if len(slug) == 0 {
		slug = deriveSlugFromWords(text)
	}
	ext := normalizeFileExt(customExt)
	if strings.HasSuffix(slug, ext) {
		return slug
	}
	return slug + ext
}

func deriveSlugFromWords(text string) string {
	words := strings.Fields(strings.ToLower(strings.TrimSpace(text)))
	if len(words) == 0 {
		return "note"
	}
	if len(words) > 5 {
		words = words[:5]
	}
	slug := sanitizeKebabSlug(strings.Join(words, "-"))
	if len(slug) == 0 {
		return "note"
	}
	return slug
}

func sanitizeKebabSlug(raw string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(raw)) {
		if isAlphaNumRune(r) {
			b.WriteRune(r)
		} else if r == '-' || r == '_' || r == ' ' || r == '.' {
			b.WriteRune('-')
		}
	}
	return strings.Trim(collapseHyphens(b.String()), "-")
}

func isAlphaNumRune(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
}

func normalizeFileExt(ext string) string {
	cleaned := strings.ToLower(strings.TrimSpace(ext))
	if len(cleaned) == 0 {
		return ".txt"
	}
	if strings.HasPrefix(cleaned, ".") {
		return cleaned
	}
	return "." + cleaned
}

func copySingleFile(srcPath, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return apperror.WrapSimple(err, "open source file")
	}
	defer srcFile.Close()
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return apperror.WrapSimple(err, "mkdir destination parent")
	}
	return writeStreamToDest(srcFile, destPath)
}

func writeStreamToDest(r io.Reader, destPath string) error {
	outFile, err := os.Create(destPath)
	if err != nil {
		return apperror.WrapSimple(err, "create destination file")
	}
	defer outFile.Close()
	if _, err := io.Copy(outFile, r); err != nil {
		return apperror.WrapSimple(err, "copy file contents")
	}
	return nil
}

func copyDirectoryRecursive(srcDir, destDir string) error {
	return filepath.Walk(srcDir, func(currentPath string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		return copyWalkEntry(srcDir, destDir, currentPath, info)
	})
}

func copyWalkEntry(srcDir, destDir, currentPath string, info os.FileInfo) error {
	rel, err := filepath.Rel(srcDir, currentPath)
	if err != nil {
		return apperror.WrapSimple(err, "compute relative path")
	}
	target := filepath.Join(destDir, rel)
	if info.IsDir() {
		return os.MkdirAll(target, 0755)
	}
	return copySingleFile(currentPath, target)
}

func autoCommitAndPushSpecialRepo(specialRepoRoot, relPath, shortKey, repoName string, shouldPush bool) (bool, bool) {
	_ = ensureSpecialGitRepo(specialRepoRoot)
	_ = runGitInSpecialRepo(specialRepoRoot, "add", "--", relPath)
	commitMsg := fmt.Sprintf("chore(%s): store %s for %s", shortKey, filepath.ToSlash(relPath), repoName)
	isCommitted := runGitInSpecialRepo(specialRepoRoot, "commit", "-m", commitMsg) == nil
	if !shouldPush || !hasGitRemoteOrigin(specialRepoRoot) {
		return isCommitted, false
	}
	isPushed := runGitInSpecialRepo(specialRepoRoot, "push") == nil
	return isCommitted, isPushed
}

func hasGitRemoteOrigin(repoDir string) bool {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = repoDir
	out, err := cmd.Output()
	return err == nil && len(strings.TrimSpace(string(out))) > 0
}

func runGitInSpecialRepo(repoDir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = repoDir
	_, err := cmd.CombinedOutput()
	return err
}

// CheckSpecialReposOnScan performs the one-time first-scan verification and SQLite persistence for repo-secrets and repo-cache.
func CheckSpecialReposOnScan(workBaseDir string, isAutoAccept bool) ([]store.SpecialRepositoryRecord, error) {
	db, err := store.OpenSpecialReposSplitDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	return CheckSpecialReposOnScanWithDB(db, workBaseDir, isAutoAccept)
}

// CheckSpecialReposOnScanWithDB executes the one-time first-scan check using the given SpecialReposSplitDB instance.
func CheckSpecialReposOnScanWithDB(db *store.SpecialReposSplitDB, workBaseDir string, isAutoAccept bool) ([]store.SpecialRepositoryRecord, error) {
	repos, err := db.ListSpecialRepos()
	if err != nil {
		return nil, err
	}
	for _, rec := range repos {
		if err := evaluateSingleSpecialRepoOnScan(db, rec, workBaseDir, isAutoAccept); err != nil {
			return nil, err
		}
	}
	return db.ListSpecialRepos()
}

func evaluateSingleSpecialRepoOnScan(db *store.SpecialReposSplitDB, rec store.SpecialRepositoryRecord, workBaseDir string, isAutoAccept bool) error {
	if rec.HasAnsweredPrompt() {
		return nil
	}
	targetDir := filepath.Join(workBaseDir, rec.ConfiguredName)
	if !hasExistingTargetDir(targetDir) {
		return handleMissingSpecialRepoOnScan(db, rec, targetDir, isAutoAccept)
	}
	logDetectedSpecialRepo(rec, targetDir, isAutoAccept)
	return db.MarkPromptAnswered(rec.ShortKey, "detected", targetDir, rec.RemoteURL)
}

func hasExistingTargetDir(targetDir string) bool {
	info, err := os.Stat(targetDir)
	return err == nil && info.IsDir()
}

func logDetectedSpecialRepo(rec store.SpecialRepositoryRecord, targetDir string, isAutoAccept bool) {
	if isAutoAccept {
		return
	}
	fmt.Printf("  ✓ Detected special %s repository (%s) at %s\n", rec.Category, rec.ShortKey, targetDir)
}

func handleMissingSpecialRepoOnScan(db *store.SpecialReposSplitDB, rec store.SpecialRepositoryRecord, targetDir string, isAutoAccept bool) error {
	if !isAutoAccept {
		printOneTimeSpecialRepoBanner(rec, targetDir)
	}
	decision := "prompted-initialized"
	if isAutoAccept {
		decision = "auto-initialized"
	}
	if err := ensureSpecialGitRepo(targetDir); err != nil {
		return err
	}
	return db.MarkPromptAnswered(rec.ShortKey, decision, targetDir, rec.RemoteURL)
}

func printOneTimeSpecialRepoBanner(rec store.SpecialRepositoryRecord, targetDir string) {
	fmt.Printf("╭── Special Repository Discovery: %s (%s) ─────────────────────────╮\n", rec.ConfiguredName, rec.ShortKey)
	if rec.ShortKey == "rs" {
		fmt.Printf("│ Purpose  : Keep secret files (.env, passwords, tokens) out of public repos.│\n")
		fmt.Printf("│ Settings : gitmap settings set special_repos.secrets_name <custom-name>    │\n")
		fmt.Printf("│ Usage    : gitmap rs file|folder|text | gitmap cd rs                       │\n")
	} else {
		fmt.Printf("│ Purpose  : Store reusable temporary scripts (.ps1, test harnesses, scratch)│\n")
		fmt.Printf("│ Settings : gitmap settings set special_repos.cache_name <custom-name>      │\n")
		fmt.Printf("│ Usage    : gitmap rc file|folder|text | gitmap cd rc                       │\n")
	}
	fmt.Printf("│ Location : %s\n", targetDir)
	fmt.Printf("╰────────────────────────────────────────────────────────────────────────────╯\n")
}
