// Package llm — llm_train_muse.go implements the positional "muse" sub-target of
// `gitmap llm train`: it clone-or-pulls the repo-secrets vault, diffs the last
// N days of the Muse memory digest, and prints a terminal summary plus a raw
// AI traversal checklist. It never imports cli/cmd (circular); git work is
// done with os/exec following the runGitInSpecialRepo pattern.
package llm

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// repoSecretsRemote is the HTTPS clone URL for the memory-digest vault.
const repoSecretsRemote = "https://github.com/alimtvnetwork/repo-secrets.git"

// digestRootName is the repo-secrets folder carrying the Muse memory digest.
const digestRootName = "Muse"

// repoSecretsPathResolver resolves the local repo-secrets checkout and reports
// whether it still needs cloning. Tests override it with fixtures.
var repoSecretsPathResolver = defaultResolveRepoSecretsPath

// runTrainMuse syncs the repo-secrets vault and renders the digest summary,
// logging every decision so an AI reading the terminal learns the why.
func runTrainMuse(opts TrainOptions) *apperror.AppError {
	log := &decisionLog{}
	path, needsClone, err := repoSecretsPathResolver()
	if err != nil {
		return err
	}
	if needsClone && isDirPath(path) {
		log.add("checkout: no 'rs' record but %s already exists at the work-folder root -> DECISION reuse it, skip clone (why: the directory is the source of truth; the record is only a pointer)", relForDisplay(path))
		needsClone = false
	}
	if needsClone {
		log.add("checkout: no usable 'rs' special-repo record and no dir -> DECISION clone fresh into %s (why: the vault must sit at the work-folder root, beside the repos it describes)", relForDisplay(path))
		if cloneErr := cloneRepoSecrets(log, path); cloneErr != nil {
			return cloneErr
		}
	} else {
		log.add("checkout: vault present at %s -> DECISION reuse it, skip clone (why: never clone over existing local state)", relForDisplay(path))
	}
	if pullErr := pullRepoSecrets(log, path); pullErr != nil {
		return pullErr
	}
	return renderMuseDigest(log, opts, path)
}

// decisionLog collects the why behind each step so an AI reading the
// terminal learns the reasoning, not just the commands.
type decisionLog struct {
	entries []string
}

func (l *decisionLog) add(format string, args ...interface{}) {
	l.entries = append(l.entries, fmt.Sprintf(format, args...))
}

func (l *decisionLog) print() {
	if len(l.entries) == 0 {
		return
	}
	fmt.Println("----------------------------------------------------------------------")
	fmt.Println("DECISION TRACE - why each step ran (AI: learn the reasoning)")
	fmt.Println("----------------------------------------------------------------------")
	for i, entry := range l.entries {
		fmt.Printf("  %d. %s\n", i+1, entry)
	}
	fmt.Println()
}

// defaultResolveRepoSecretsPath resolves the checkout through the store
// special-repo entry (short key "rs"); absent entries resolve to a clone
// target in the current (work-folder) directory.
func defaultResolveRepoSecretsPath() (string, bool, *apperror.AppError) {
	db, err := store.OpenSpecialReposSplitDB()
	if err != nil {
		return "", false, apperror.WrapSimple(err, "open special-repos db")
	}
	defer db.Close()
	rec, recErr := db.GetSpecialRepo("rs")
	if recErr == nil && rec != nil && isDirPath(rec.LocalPath) {
		return rec.LocalPath, false, nil
	}
	return defaultCloneTarget(rec), true, nil
}

// defaultCloneTarget builds the clone destination: <cwd>/<configured name>.
func defaultCloneTarget(rec *store.SpecialRepoRecord) string {
	name := "repo-secrets"
	if rec != nil && rec.ConfiguredName != "" {
		name = rec.ConfiguredName
	}
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	return filepath.Join(cwd, name)
}

// isDirPath reports whether path exists and is a directory.
func isDirPath(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// cloneRepoSecrets clones the vault; fails fast with guidance when gh is
// missing or unauthenticated — stdin is non-interactive so gh never prompts.
func cloneRepoSecrets(log *decisionLog, dest string) *apperror.AppError {
	if err := checkGhReady(); err != nil {
		return err
	}
	log.add("clone: gh present and authenticated -> ran: gh repo clone %s %s (why: stdin is closed so gh can never prompt; fail fast with guidance instead)", repoSecretsRemote, relForDisplay(dest))
	fmt.Printf("[LLM-TRAIN-MUSE] Cloning repo-secrets into %s ...\n", relForDisplay(dest))
	cmd := exec.Command("gh", "repo", "clone", repoSecretsRemote, dest)
	cmd.Stdin = bytes.NewReader(nil)
	var combined bytes.Buffer
	cmd.Stdout = &combined
	cmd.Stderr = &combined
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(fmt.Errorf("%w: %s", err, strings.TrimSpace(combined.String())), "clone repo-secrets")
	}
	return nil
}

// checkGhReady verifies gh exists and is authenticated, printing actionable
// guidance otherwise.
func checkGhReady() *apperror.AppError {
	if _, lookErr := exec.LookPath("gh"); lookErr != nil {
		fmt.Println("[LLM-TRAIN-MUSE] 'gh' CLI not found on PATH. Install it from https://cli.github.com/ and re-run.")
		return apperror.NewSimple("llm.train.muse.clone", "E_GH_MISSING")
	}
	auth := exec.Command("gh", "auth", "status")
	auth.Stdin = bytes.NewReader(nil)
	if out, err := auth.CombinedOutput(); err != nil {
		fmt.Println("[LLM-TRAIN-MUSE] gh is not authenticated. Run 'gh auth login' (or set GH_TOKEN), then re-run.")
		fmt.Printf("[LLM-TRAIN-MUSE] gh auth status: %s\n", strings.TrimSpace(string(out)))
		return apperror.NewSimple("llm.train.muse.clone", "E_GH_AUTH_MISSING")
	}
	return nil
}

// pullRepoSecrets fast-forward pulls; skips with a note when no origin remote.
func pullRepoSecrets(log *decisionLog, path string) *apperror.AppError {
	if !hasGitRemote(path) {
		log.add("sync: no origin remote on %s -> DECISION skip pull (why: read-only local state, nothing to fetch)", relForDisplay(path))
		fmt.Println("[LLM-TRAIN-MUSE] No origin remote; skipping pull.")
		return nil
	}
	log.add("sync: origin remote present -> DECISION git pull --ff-only (why: take new digest commits fast-forward only; never merge, never rewrite history)")
	_, err := runGitCmd(path, "pull", "--ff-only")
	return err
}

// hasGitRemote reports whether dir has an origin remote configured.
func hasGitRemote(dir string) bool {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = dir
	cmd.Stdin = bytes.NewReader(nil)
	return cmd.Run() == nil
}

// runGitCmd runs git in dir and returns stdout; errors are wrapped with stderr.
func runGitCmd(dir string, args ...string) ([]byte, *apperror.AppError) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Stdin = bytes.NewReader(nil)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, apperror.WrapSimple(fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out))), "git "+strings.Join(args, " "))
	}
	return out, nil
}

// renderMuseDigest renders the summary + checklist, or missing-digest guidance.
func renderMuseDigest(log *decisionLog, opts TrainOptions, path string) *apperror.AppError {
	digestDir := filepath.Join(path, digestRootName)
	if !isDirPath(digestDir) {
		log.add("digest: no %s/ dir at %s -> DECISION print guidance, exit 0 (why: nothing to train on yet; the digest is authored by the memory-IO program)", digestRootName, relForDisplay(path))
		log.print()
		printDigestMissingGuidance(path)
		return nil
	}
	files, commitCount, err := gitLogChangedFiles(path, opts.DaysN)
	if err != nil {
		return err
	}
	since := time.Now().AddDate(0, 0, -opts.DaysN).Format("2006-01-02")
	log.add("window: --days %d -> since %s -> DECISION git log --since + rev-list --count (why: bound the diff to the digest window; the commit count frames the summary)", opts.DaysN, since)
	log.add("order: DECISION README.md first, then sequenced files in lexical order, chat-seeds on demand (why: contract before detail; a deterministic order any LLM can replay)")
	log.add("search: DECISION scope gitmap aum search to the digest tree with --ext .md (why: keep the model inside the digest, no repo-wide noise)")
	log.print()
	printDigestSummary(path, opts.DaysN, files, commitCount)
	return nil
}

// gitLogChangedFiles returns the deduplicated changed-file list and commit
// count for the last N days.
func gitLogChangedFiles(path string, days int) ([]string, int, *apperror.AppError) {
	since := fmt.Sprintf("%d days ago", days)
	raw, err := runGitCmd(path, "log", "--since="+since, "--name-only", "--pretty=format:")
	if err != nil {
		return nil, 0, err
	}
	countRaw, countErr := runGitCmd(path, "rev-list", "--count", "--since="+since, "HEAD")
	if countErr != nil {
		return nil, 0, countErr
	}
	return uniqueNonEmptyLines(string(raw)), parseCount(countRaw), nil
}

// uniqueNonEmptyLines dedupes and sorts non-empty lines.
func uniqueNonEmptyLines(raw string) []string {
	seen := map[string]bool{}
	var out []string
	for _, line := range strings.Split(raw, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || seen[trimmed] {
			continue
		}
		seen[trimmed] = true
		out = append(out, trimmed)
	}
	sort.Strings(out)
	return out
}

// parseCount parses a git rev-list count, defaulting to 0 on garbage.
func parseCount(raw []byte) int {
	count, err := strconv.Atoi(strings.TrimSpace(string(raw)))
	if err != nil || count < 0 {
		return 0
	}
	return count
}

// printDigestSummary prints the terminal summary and traversal checklist.
func printDigestSummary(path string, days int, files []string, commitCount int) {
	digestDir := filepath.Join(path, digestRootName)
	sequenced := listSequencedDigestFiles(digestDir)
	fmt.Println("======================================================================")
	fmt.Printf("GITMAP LLM TRAIN — MUSE MEMORY SYNC (last %d days)\n", days)
	fmt.Println("======================================================================")
	fmt.Printf("Digest root: %s\n", relForDisplay(digestDir))
	fmt.Printf("Commits: %d | Files changed: %d | Sequenced digests: %d\n\n", commitCount, len(files), len(sequenced))
	printDigestFileList("SEQUENCED DIGEST FILES", sequenced)
	printDigestFileList("CHANGED FILES", files)
	printTraversalChecklist(path, digestDir, sequenced)
}

// listSequencedDigestFiles returns the digest-root .md files, digest-relative.
func listSequencedDigestFiles(digestDir string) []string {
	entries, err := os.ReadDir(digestDir)
	if err != nil {
		return nil
	}
	var files []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		files = append(files, digestRootName+"/"+entry.Name())
	}
	sort.Strings(files)
	return files
}

// printDigestFileList prints a titled file list with an empty-state fallback.
func printDigestFileList(title string, files []string) {
	fmt.Println(title + ":")
	if len(files) == 0 {
		fmt.Println("  (none)")
		fmt.Println()
		return
	}
	for _, f := range files {
		fmt.Printf("  %s\n", f)
	}
	fmt.Println()
}

// printTraversalChecklist prints the raw AI traversal checklist with paths
// relative to the working directory.
func printTraversalChecklist(path, digestDir string, sequenced []string) {
	fmt.Println("======================================================================")
	fmt.Println("AI TRAVERSAL CHECKLIST (execute in order)")
	fmt.Println("======================================================================")
	fmt.Printf("1. Read the digest root contract first (why: it defines the say-READ protocol and the file order):\n   gitmap cat %s\n", relForDisplay(filepath.Join(digestDir, "README.md")))
	fmt.Println("2. Walk the sequenced digests in order (why: lexical order is the authored narrative):")
	for _, f := range sequenced {
		fmt.Printf("   gitmap cat %s\n", relForDisplay(filepath.Join(path, filepath.FromSlash(f))))
	}
	fmt.Printf("3. Scope searches to the digest tree (why: stay inside the digest, avoid repo-wide noise):\n   gitmap aum search \"<keyword>\" %s --ext .md\n", relForDisplay(digestDir))
	fmt.Printf("4. Inspect chat seeds when a conversation needs full context (why: a seed rebuilds one chat without re-reading everything):\n   gitmap lf %s\n", relForDisplay(filepath.Join(digestDir, "chat-seeds")))
	fmt.Println()
}

// printDigestMissingGuidance handles the absent-digest case: guidance, exit 0.
func printDigestMissingGuidance(path string) {
	fmt.Println("======================================================================")
	fmt.Println("GITMAP LLM TRAIN — MUSE MEMORY SYNC")
	fmt.Println("======================================================================")
	fmt.Printf("[LLM-TRAIN-MUSE] No %s/ digest directory at %s.\n", digestRootName, relForDisplay(path))
	fmt.Println("The 7-day Muse digest is authored by the Muse memory-IO program (another worker).")
	fmt.Println("Nothing to traverse yet — once the digest lands, re-run:")
	fmt.Println("  gitmap llm train muse --days 7")
	fmt.Println()
}

// relForDisplay renders abs relative to the working directory when possible.
func relForDisplay(abs string) string {
	cwd, err := os.Getwd()
	if err != nil {
		return filepath.ToSlash(abs)
	}
	rel, relErr := filepath.Rel(cwd, abs)
	if relErr != nil {
		return filepath.ToSlash(abs)
	}
	return filepath.ToSlash(rel)
}
