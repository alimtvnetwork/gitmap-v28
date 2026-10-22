// Package llm — llm_train_muse.go implements the positional "muse" sub-target of
// `gitmap llm train`: it checks repo-cache first for the 12-muse-ai (or XX-muse-ai)
// directory with backward-compatible fallback to repo-secrets, diffs the last N days
// of the Muse memory digest, and prints a terminal summary plus a raw AI traversal
// checklist. The last-read watermark is tracked by gitmap itself (JSON under BinaryDataDir).
package llm

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// repoCacheRemote is the HTTPS clone URL for repo-cache.
const repoCacheRemote = "https://github.com/alimtvnetwork/repo-cache.git"

// repoSecretsRemote is the HTTPS clone URL for the memory-digest vault.
const repoSecretsRemote = "https://github.com/alimtvnetwork/repo-secrets.git"

// defaultDigestRootName is the default folder carrying the Muse memory digest.
const defaultDigestRootName = "12-muse-ai"

// repoSecretsPathResolver resolves the local training checkout and reports
// whether it still needs cloning. Tests override it with fixtures.
var repoSecretsPathResolver = defaultResolveTrainingPath

// decisionLog collects the why behind each step so an AI reading the
// terminal learns the rationales.
type decisionLog struct {
	entries []string
}

func (d *decisionLog) add(format string, args ...any) {
	d.entries = append(d.entries, fmt.Sprintf(format, args...))
}

func (d *decisionLog) print() {
	if len(d.entries) == 0 {
		return
	}
	fmt.Println("DECISION TRACE (why gitmap took each step):")
	for i, entry := range d.entries {
		fmt.Printf("  [%d] %s\n", i+1, entry)
	}
	fmt.Println()
}

// runTrainMuse syncs the training vault and renders the digest summary,
// logging every decision so an AI reading the terminal learns the why.
func runTrainMuse(opts TrainOptions) *apperror.AppError {
	log := &decisionLog{}
	path, needsClone, err := repoSecretsPathResolver()
	if err != nil {
		return err
	}
	if needsClone && isDirPath(path) {
		log.add("checkout: %s already exists -> DECISION reuse it, skip clone (why: directory is source of truth)", relForDisplay(path))
		needsClone = false
	}
	if cloneErr := ensureTrainingRepoCloned(log, path, needsClone); cloneErr != nil {
		return cloneErr
	}
	if pullErr := pullTrainingRepo(log, path); pullErr != nil {
		return pullErr
	}
	headSHA, headErr := gitHeadSHA(path)
	if headErr != nil {
		return headErr
	}
	prev := loadLastRead(path)
	since := diffSinceLastRead(log, path, prev, headSHA)

	digestDir, digestName, hasDigest := resolveMuseDigestDir(path)
	if !hasDigest {
		path, digestDir, digestName, hasDigest = fallbackToRepoSecretsDigest(log, path)
	}

	if !hasDigest {
		digestName = defaultDigestRootName
		log.add("digest: no %s/ (or XX-muse-ai/) dir at %s -> DECISION print guidance, exit 0 (why: nothing to train on yet; the digest is authored by the memory-IO program)", digestName, relForDisplay(path))
		saveLastRead(log, path, headSHA)
		log.print()
		printDigestMissingGuidance(path, digestName)
		return nil
	}

	files, commitCount, err := gitLogChangedFiles(path, opts.DaysN)
	if err != nil {
		return err
	}
	logTrainingDecisions(log, opts.DaysN)
	saveLastRead(log, path, headSHA)
	log.print()
	printDigestSummary(path, digestDir, digestName, opts.DaysN, files, commitCount, since)
	return nil
}

func ensureTrainingRepoCloned(log *decisionLog, path string, needsClone bool) *apperror.AppError {
	if !needsClone {
		log.add("checkout: training source present at %s -> DECISION reuse it, skip clone (why: never clone over existing local state)", relForDisplay(path))
		return nil
	}
	log.add("checkout: training source absent -> DECISION clone fresh into %s", relForDisplay(path))
	return cloneTrainingRepo(log, path)
}

func fallbackToRepoSecretsDigest(log *decisionLog, path string) (string, string, string, bool) {
	if strings.Contains(filepath.Base(path), "repo-secrets") {
		return path, "", "", false
	}
	for _, secPath := range probeRepoCandidates("repo-secrets", "rs") {
		dDir, dName, ok := resolveMuseDigestDir(secPath)
		if ok {
			log.add("digest: found in fallback repo-secrets at %s -> DECISION switch to repo-secrets", relForDisplay(secPath))
			return secPath, dDir, dName, true
		}
	}
	return path, "", "", false
}

// defaultResolveTrainingPath resolves the training source directory.
// It checks repo-cache first for 12-muse-ai (or XX-muse-ai); repo-secrets
// stays as a backward-compatible fallback (the digest moved out in 2026-10).
func defaultResolveTrainingPath() (string, bool, *apperror.AppError) {
	cacheCandidates := probeRepoCandidates("repo-cache", "rc")
	for _, path := range cacheCandidates {
		if _, _, hasDigest := resolveMuseDigestDir(path); hasDigest {
			return path, false, nil
		}
	}

	secretCandidates := probeRepoCandidates("repo-secrets", "rs")
	for _, path := range secretCandidates {
		if _, _, hasDigest := resolveMuseDigestDir(path); hasDigest {
			return path, false, nil
		}
	}

	if len(cacheCandidates) > 0 {
		return cacheCandidates[0], false, nil
	}

	if len(secretCandidates) > 0 {
		return secretCandidates[0], false, nil
	}

	return defaultCloneTargetFor("repo-cache", "rc"), true, nil
}

func defaultCloneTargetFor(key, shortKey string) string {
	name := key
	if configured := querySpecialRepoConfiguredName(shortKey); configured != "" {
		name = configured
	}
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	parent := filepath.Dir(cwd)
	isGitmapParent := filepath.Base(cwd) == "gitmap" || !isDirPath(filepath.Join(parent, ".git"))
	if isGitmapParent && isDirPath(parent) {
		return filepath.Join(parent, name)
	}
	return filepath.Join(cwd, name)
}

func querySpecialRepoConfiguredName(shortKey string) string {
	db, err := store.OpenSpecialReposSplitDB()
	if err != nil {
		return ""
	}
	defer db.Close()
	rec, recErr := db.GetSpecialRepo(shortKey)
	if recErr != nil || rec == nil {
		return ""
	}
	return rec.ConfiguredName
}

func isDirPath(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func cloneTrainingRepo(log *decisionLog, dest string) *apperror.AppError {
	if err := checkGhReady(); err != nil {
		return err
	}
	remote := repoCacheRemote
	repoName := filepath.Base(dest)
	if strings.Contains(repoName, "repo-secrets") {
		remote = repoSecretsRemote
	}
	log.add("clone: gh present and authenticated -> ran: gh repo clone %s %s (why: stdin is closed so gh can never prompt; fail fast with guidance instead)", remote, relForDisplay(dest))
	fmt.Printf("[LLM-TRAIN-MUSE] Cloning %s into %s ...\n", repoName, relForDisplay(dest))
	cmd := exec.Command("gh", "repo", "clone", remote, dest)
	cmd.Stdin = bytes.NewReader(nil)
	var combined bytes.Buffer
	cmd.Stdout = &combined
	cmd.Stderr = &combined
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "clone "+repoName).WithContext("output", strings.TrimSpace(combined.String()))
	}
	return nil
}

func checkGhReady() *apperror.AppError {
	if _, lookErr := exec.LookPath("gh"); lookErr != nil {
		fmt.Println("[LLM-TRAIN-MUSE] 'gh' CLI not found on PATH. Install it from https://cli.github.com/ and re-run.")
		return apperror.NewSimple("llm.train.muse.clone", "E_GH_MISSING")
	}
	auth := exec.Command("gh", "auth", "status")
	auth.Stdin = bytes.NewReader(nil)
	out, err := auth.CombinedOutput()
	if err != nil {
		fmt.Println("[LLM-TRAIN-MUSE] gh is not authenticated. Run 'gh auth login' (or set GH_TOKEN), then re-run.")
		fmt.Printf("[LLM-TRAIN-MUSE] gh auth status: %s\n", strings.TrimSpace(string(out)))
		return apperror.NewSimple("llm.train.muse.clone", "E_GH_AUTH_MISSING")
	}
	return nil
}

func pullTrainingRepo(log *decisionLog, path string) *apperror.AppError {
	if !hasGitRemote(path) {
		log.add("sync: no origin remote on %s -> DECISION skip pull (why: read-only local state, nothing to fetch)", relForDisplay(path))
		fmt.Println("[LLM-TRAIN-MUSE] No origin remote; skipping pull.")
		return nil
	}
	log.add("sync: origin remote present -> DECISION git pull --ff-only (why: take new digest commits fast-forward only; never merge, never rewrite history)")
	_, err := runGitCmd(path, "pull", "--ff-only")
	return err
}

func hasGitRemote(dir string) bool {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = dir
	cmd.Stdin = bytes.NewReader(nil)
	return cmd.Run() == nil
}

func runGitCmd(dir string, args ...string) ([]byte, *apperror.AppError) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Stdin = bytes.NewReader(nil)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, apperror.WrapSimple(err, "git "+strings.Join(args, " ")).WithContext("output", strings.TrimSpace(string(out)))
	}
	return out, nil
}
