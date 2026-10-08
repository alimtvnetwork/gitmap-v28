// Package cmd — special_repos_ops.go provides filesystem sequencing, git auto-commit/push, and first-scan check operations for repo-secrets (rs) and repo-cache (rc).
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddb"
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
	return fallbackWorkBaseDir()
}

func fallbackWorkBaseDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return filepath.Dir(cwd)
}

func resolveSpecialRepoCDPath(target string) (string, error) {
	db, err := store.OpenSpecialReposSplitDB()
	if err != nil {
		return "", err
	}
	defer db.Close()
	return resolveSpecialRepoCDPathWithDB(db, target, resolveSpecialWorkBaseDir())
}

func resolveSpecialRepoCDPathWithDB(db *store.SpecialReposSplitDB, target, workBaseDir string) (string, error) {
	rec, err := db.GetSpecialRepo(target)
	if err != nil {
		return "", err
	}
	if len(rec.LocalPath) > 0 && isDirValid(rec.LocalPath) {
		return rec.LocalPath, nil
	}
	return filepath.Join(workBaseDir, rec.ConfiguredName), nil
}

func isDirValid(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
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
	targetPath := resolveValidRepoPath(rec.LocalPath, workBaseDir, rec.ConfiguredName)
	if err := ensureSpecialGitRepo(targetPath); err != nil {
		return "", nil, err
	}
	return targetPath, rec, nil
}

func resolveValidRepoPath(localPath, workBaseDir, configuredName string) string {
	if len(localPath) > 0 && isDirValid(localPath) {
		return localPath
	}
	return filepath.Join(workBaseDir, configuredName)
}

func ensureSpecialGitRepo(repoDir string) error {
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		return apperror.WrapSimple(err, "mkdir special repo")
	}
	if hasGitSubdir(repoDir) {
		return nil
	}
	return runGitInitInDir(repoDir)
}

func hasGitSubdir(repoDir string) bool {
	info, err := os.Stat(filepath.Join(repoDir, ".git"))
	return err == nil && info.IsDir()
}

func runGitInitInDir(repoDir string) error {
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
	topDir := findEnclosingGitTopLevel()
	if len(topDir) > 0 {
		return stripNumericPrefix(filepath.Base(topDir))
	}
	return fallbackCurrentRepoName()
}

func fallbackCurrentRepoName() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "gitmap"
	}
	return stripNumericPrefix(filepath.Base(cwd))
}

func findEnclosingGitTopLevel() string {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return ""
	}
	trimmed := strings.TrimSpace(string(out))
	if len(trimmed) == 0 {
		return ""
	}
	return filepath.Clean(trimmed)
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
	nextSeq, _ := store.GetNextFileSeq(repoFolderDir, "")
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
		appendSanitizedSlugRune(&b, r)
	}
	return strings.Trim(collapseHyphens(b.String()), "-")
}

func appendSanitizedSlugRune(b *strings.Builder, r rune) {
	if isAlphaNumRune(r) {
		b.WriteRune(r)
	} else if r == '-' || r == '_' || r == ' ' || r == '.' {
		b.WriteRune('-')
	}
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

func autoCommitAndPushSpecialRepo(specialRepoRoot, relPath, shortKey, repoName, customMsg string, shouldPush bool) (bool, bool) {
	_ = ensureSpecialGitRepo(specialRepoRoot)
	_ = runGitInSpecialRepo(specialRepoRoot, "add", "--", relPath)
	commitMsg := resolveSpecialCommitMsg(shortKey, relPath, repoName, customMsg)
	isCommitted := runGitInSpecialRepo(specialRepoRoot, "commit", "-m", commitMsg) == nil
	if !shouldPush || !hasGitRemoteOrigin(specialRepoRoot) {
		return isCommitted, false
	}
	isPushed := tryPushSpecialRepo(specialRepoRoot)
	return isCommitted, isPushed
}

func resolveSpecialCommitMsg(shortKey, relPath, repoName, customMsg string) string {
	if len(strings.TrimSpace(customMsg)) > 0 {
		return strings.TrimSpace(customMsg)
	}
	return fmt.Sprintf("chore(%s): store %s for %s", shortKey, filepath.ToSlash(relPath), repoName)
}

func tryPushSpecialRepo(repoDir string) bool {
	if runGitInSpecialRepo(repoDir, "push") == nil {
		return true
	}
	return runGitInSpecialRepo(repoDir, "push", "-u", "origin", "HEAD") == nil
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
	fmt.Printf("  ✓ Detected special %s repository (%s) at %s\n", rec.Category, rec.ShortKey, targetDir)
}

func handleMissingSpecialRepoOnScan(db *store.SpecialReposSplitDB, rec store.SpecialRepositoryRecord, targetDir string, isAutoAccept bool) error {
	remoteURL := probeRemoteSpecialRepoURL(rec.ConfiguredName)
	if !isAutoAccept {
		printOneTimeSpecialRepoBanner(rec, targetDir, remoteURL)
	}
	if isAutoAccept || !cmddb.IsInteractiveStdin() {
		return applyAutoSpecialRepoResolution(db, rec, targetDir, remoteURL)
	}
	return promptInteractiveSpecialRepoResolution(db, rec, targetDir, remoteURL)
}

func probeRemoteSpecialRepoURL(repoName string) string {
	cmd := exec.Command("gh", "repo", "view", repoName, "--json", "url", "-q", ".url")
	out, err := cmd.Output()
	if err == nil && len(strings.TrimSpace(string(out))) > 0 {
		return strings.TrimSpace(string(out))
	}
	return ""
}

func promptInteractiveSpecialRepoResolution(db *store.SpecialReposSplitDB, rec store.SpecialRepositoryRecord, targetDir, remoteURL string) error {
	reader := bufio.NewReader(os.Stdin)
	promptMsg := formatSpecialRepoPrompt(rec.ConfiguredName, targetDir, remoteURL)
	fmt.Print(promptMsg)
	ans, _ := reader.ReadString('\n')
	cleanAns := strings.TrimSpace(ans)
	if strings.EqualFold(cleanAns, "n") || strings.EqualFold(cleanAns, "no") {
		fmt.Printf("  [info] Skipped %s setup. (Recorded in settings: gitmap settings)\n\n", rec.ConfiguredName)
		return db.MarkPromptAnswered(rec.ShortKey, "declined", "", remoteURL)
	}
	return executeSpecialRepoAccept(db, rec, targetDir, remoteURL)
}

func formatSpecialRepoPrompt(repoName, targetDir, remoteURL string) string {
	if len(remoteURL) > 0 {
		return fmt.Sprintf("Do you like to clone the %s repository from %s into %s? [Y/n]: ", repoName, remoteURL, targetDir)
	}
	return fmt.Sprintf("Do you like to create the %s repository at %s? [Y/n]: ", repoName, targetDir)
}

func executeSpecialRepoAccept(db *store.SpecialReposSplitDB, rec store.SpecialRepositoryRecord, targetDir, remoteURL string) error {
	if len(remoteURL) > 0 {
		return cloneAndRecordSpecialRepo(db, rec, targetDir, remoteURL)
	}
	return createAndRecordSpecialRepo(db, rec, targetDir, remoteURL)
}

func cloneAndRecordSpecialRepo(db *store.SpecialReposSplitDB, rec store.SpecialRepositoryRecord, targetDir, remoteURL string) error {
	cmd := exec.Command("git", "clone", remoteURL, targetDir)
	if out, err := cmd.CombinedOutput(); err != nil {
		return apperror.WrapSimple(err, "clone special repo: "+string(out))
	}
	fmt.Printf("  ✔ Cloned %s into %s\n\n", rec.ConfiguredName, targetDir)
	return db.MarkPromptAnswered(rec.ShortKey, "cloned", targetDir, remoteURL)
}

func createAndRecordSpecialRepo(db *store.SpecialReposSplitDB, rec store.SpecialRepositoryRecord, targetDir, remoteURL string) error {
	if err := ensureSpecialGitRepo(targetDir); err != nil {
		return err
	}
	fmt.Printf("  ✔ Created %s repository at %s\n\n", rec.ConfiguredName, targetDir)
	return db.MarkPromptAnswered(rec.ShortKey, "created", targetDir, remoteURL)
}

func applyAutoSpecialRepoResolution(db *store.SpecialReposSplitDB, rec store.SpecialRepositoryRecord, targetDir, remoteURL string) error {
	if len(remoteURL) > 0 {
		_ = exec.Command("git", "clone", remoteURL, targetDir).Run()
		return db.MarkPromptAnswered(rec.ShortKey, "auto-cloned", targetDir, remoteURL)
	}
	if err := ensureSpecialGitRepo(targetDir); err != nil {
		return err
	}
	return db.MarkPromptAnswered(rec.ShortKey, "auto-initialized", targetDir, remoteURL)
}

func printOneTimeSpecialRepoBanner(rec store.SpecialRepositoryRecord, targetDir, remoteURL string) {
	fmt.Printf("╭── Special Repository Discovery: %s (%s) ─────────────────────────╮\n", rec.ConfiguredName, rec.ShortKey)
	printSpecialRepoBannerPurpose(rec.ShortKey)
	fmt.Printf("│ Default  : %s is default (change anytime via 'gitmap settings')\n", rec.DefaultName)
	if len(remoteURL) > 0 {
		fmt.Printf("│ GitHub   : Found existing remote in your account: %s\n", remoteURL)
	}
	fmt.Printf("│ Location : %s\n", targetDir)
	fmt.Printf("╰────────────────────────────────────────────────────────────────────────────╯\n")
}

func printSpecialRepoBannerPurpose(shortKey string) {
	if shortKey == "rs" {
		fmt.Printf("│ Purpose  : Keep secret files (.env, passwords, tokens) out of public repos.│\n")
		fmt.Printf("│ Settings : gitmap settings set special_repos.secrets_name <custom-name>    │\n")
		fmt.Printf("│ Usage    : gitmap rs file|folder|text | gitmap cd rs                       │\n")
		return
	}
	fmt.Printf("│ Purpose  : Store reusable temporary scripts (.ps1, test harnesses, scratch)│\n")
	fmt.Printf("│ Settings : gitmap settings set special_repos.cache_name <custom-name>      │\n")
	fmt.Printf("│ Usage    : gitmap rc file|folder|text | gitmap cd rc                       │\n")
}

// RunSpecialRepoProbeForPull probes remote GitHub account and local workspace for companion repositories on pull.
func RunSpecialRepoProbeForPull(workBaseDir string, isAutoAccept bool) error {
	baseDir := workBaseDir
	if len(baseDir) == 0 {
		baseDir = resolveSpecialWorkBaseDir()
	}
	db, err := store.OpenSpecialReposSplitDB()
	if err != nil {
		return err
	}
	defer db.Close()
	return probeSpecialReposOnPullWithDB(db, baseDir, isAutoAccept)
}

func probeSpecialReposOnPullWithDB(db *store.SpecialReposSplitDB, workBaseDir string, isAutoAccept bool) error {
	repos, err := db.ListSpecialRepos()
	if err != nil {
		return err
	}
	for _, rec := range repos {
		if err := evaluateSingleSpecialRepoOnPull(db, rec, workBaseDir, isAutoAccept); err != nil {
			return err
		}
	}
	return nil
}

func evaluateSingleSpecialRepoOnPull(db *store.SpecialReposSplitDB, rec store.SpecialRepositoryRecord, workBaseDir string, isAutoAccept bool) error {
	targetDir := resolveValidRepoPath(rec.LocalPath, workBaseDir, rec.ConfiguredName)
	if hasExistingTargetDir(targetDir) {
		logDetectedSpecialRepo(rec, targetDir, false)
		return db.MarkPromptAnswered(rec.ShortKey, "detected", targetDir, rec.RemoteURL)
	}
	return handleMissingSpecialRepoOnScan(db, rec, targetDir, isAutoAccept)
}
