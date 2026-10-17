// Package llm — llm_train_muse.go implements the positional "muse" sub-target of
// `gitmap llm train`: it checks repo-cache first for the 12-muse-ai (or XX-muse-ai)
// directory with fallback to repo-secrets, diffs the last N days of the Muse memory
// digest, and prints a terminal summary plus a raw AI traversal checklist. The
// last-read watermark is tracked by gitmap itself (JSON under BinaryDataDir) — never
// in an MD file. It never imports cli/cmd (circular); git work is done with os/exec
// following the runGitInSpecialRepo pattern.
package llm

import (
	"bytes"
	"encoding/json"
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

// repoCacheRemote is the HTTPS clone URL for repo-cache.
const repoCacheRemote = "https://github.com/alimtvnetwork/repo-cache.git"

// repoSecretsRemote is the HTTPS clone URL for the fallback memory-digest vault.
const repoSecretsRemote = "https://github.com/alimtvnetwork/repo-secrets.git"

// defaultDigestRootName is the default folder carrying the Muse memory digest.
const defaultDigestRootName = "12-muse-ai"

// digestRootName matches legacy reference.
const digestRootName = defaultDigestRootName

// repoSecretsPathResolver resolves the local training checkout and reports
// whether it still needs cloning. Tests override it with fixtures.
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

// runTrainMuse syncs the training vault (repo-cache, fallback repo-secrets) and renders
// the digest summary, logging every decision so an AI reading the terminal learns the why.
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
	if needsClone {
		log.add("checkout: training source absent -> DECISION clone fresh into %s", relForDisplay(path))
		if cloneErr := cloneTrainingRepo(log, path); cloneErr != nil {
			return cloneErr
		}
	} else {
		log.add("checkout: training source present at %s -> DECISION reuse it, skip clone (why: never clone over existing local state)", relForDisplay(path))
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
		// If path was repo-cache without digest, check fallback repo-secrets
		if !strings.Contains(filepath.Base(path), "repo-secrets") {
			for _, secPath := range probeRepoCandidates("repo-secrets", "rs") {
				if dDir, dName, ok := resolveMuseDigestDir(secPath); ok {
					path = secPath
					digestDir = dDir
					digestName = dName
					hasDigest = true
					log.add("digest: found in fallback repo-secrets at %s -> DECISION switch to repo-secrets", relForDisplay(path))
					break
				}
			}
		}
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
	sinceDate := time.Now().AddDate(0, 0, -opts.DaysN).Format("2006-01-02")
	log.add("window: --days %d -> since %s -> DECISION git log --since + rev-list --count (why: bound the diff to the digest window; the commit count frames the summary)", opts.DaysN, sinceDate)
	log.add("order: DECISION README.md first, then sequenced files in lexical order, chat-seeds on demand (why: contract before detail; a deterministic order any LLM can replay)")
	log.add("search: DECISION scope gitmap aum search to the digest tree with --ext .md (why: keep the model inside the digest, no repo-wide noise)")
	saveLastRead(log, path, headSHA)
	log.print()
	printDigestSummary(path, digestDir, digestName, opts.DaysN, files, commitCount, since)
	return nil
}

// resolveMuseDigestDir inspects a candidate repo directory for 12-muse-ai or XX-muse-ai.
// Returns (digestDir, digestName, true) if found, or ("", "", false).
func resolveMuseDigestDir(repoRoot string) (string, string, bool) {
	if !isDirPath(repoRoot) {
		return "", "", false
	}

	// 1. Check exact 12-muse-ai first
	exact := filepath.Join(repoRoot, defaultDigestRootName)
	if isDirPath(exact) {
		return exact, defaultDigestRootName, true
	}

	// 2. Scan for XX-muse-ai or *-muse-ai
	entries, err := os.ReadDir(repoRoot)
	if err == nil {
		for _, e := range entries {
			if e.IsDir() {
				name := e.Name()
				if strings.HasSuffix(name, "-muse-ai") || name == "muse-ai" || name == "muse" {
					return filepath.Join(repoRoot, name), name, true
				}
			}
		}
	}

	return "", "", false
}

// probeRepoCandidates returns all existing local directories for a special repo name
// ("repo-cache" or "repo-secrets") by querying SQLite split-db and local relative paths.
func probeRepoCandidates(key, shortKey string) []string {
	var candidates []string
	seen := make(map[string]bool)
	addCandidate := func(p string) {
		if p == "" {
			return
		}
		clean := filepath.Clean(p)
		if !seen[clean] {
			seen[clean] = true
			candidates = append(candidates, clean)
		}
	}

	if db, err := store.OpenSpecialReposSplitDB(); err == nil {
		defer db.Close()
		if rec, recErr := db.GetSpecialRepo(shortKey); recErr == nil && rec != nil && rec.LocalPath != "" {
			addCandidate(rec.LocalPath)
		}
		if rec, recErr := db.GetSpecialRepo(key); recErr == nil && rec != nil && rec.LocalPath != "" {
			addCandidate(rec.LocalPath)
		}
	}

	cwd, _ := os.Getwd()
	if cwd != "" {
		addCandidate(filepath.Join(cwd, key))
		addCandidate(filepath.Join(filepath.Dir(cwd), key))
		addCandidate(filepath.Join(filepath.Dir(filepath.Dir(cwd)), key))
	}
	addCandidate(key)

	var valid []string
	for _, c := range candidates {
		if isDirPath(c) {
			valid = append(valid, c)
		}
	}
	return valid
}

// defaultResolveTrainingPath resolves the training source directory.
// It checks repo-secrets first for 12-muse-ai (or XX-muse-ai) for now,
// with repo-cache as forward-compatible source.
func defaultResolveTrainingPath() (string, bool, *apperror.AppError) {
	// 1. Check existing repo-secrets checkouts for 12-muse-ai or XX-muse-ai (primary for now)
	secretCandidates := probeRepoCandidates("repo-secrets", "rs")
	for _, path := range secretCandidates {
		if _, _, hasDigest := resolveMuseDigestDir(path); hasDigest {
			return path, false, nil
		}
	}

	// 2. Check existing repo-cache checkouts for 12-muse-ai or XX-muse-ai (future transition)
	cacheCandidates := probeRepoCandidates("repo-cache", "rc")
	for _, path := range cacheCandidates {
		if _, _, hasDigest := resolveMuseDigestDir(path); hasDigest {
			return path, false, nil
		}
	}

	// 3. Fallback to existing repo-secrets dir if present
	if len(secretCandidates) > 0 {
		return secretCandidates[0], false, nil
	}

	// 4. Fallback to existing repo-cache dir if present
	if len(cacheCandidates) > 0 {
		return cacheCandidates[0], false, nil
	}

	// 5. Default clone target: repo-secrets at work-folder root
	return defaultCloneTargetFor("repo-secrets", "rs"), true, nil
}

func defaultResolveRepoSecretsPath() (string, bool, *apperror.AppError) {
	return defaultResolveTrainingPath()
}

func defaultCloneTargetFor(key, shortKey string) string {
	name := key
	if db, err := store.OpenSpecialReposSplitDB(); err == nil {
		defer db.Close()
		if rec, recErr := db.GetSpecialRepo(shortKey); recErr == nil && rec != nil && rec.ConfiguredName != "" {
			name = rec.ConfiguredName
		}
	}
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	parent := filepath.Dir(cwd)
	if filepath.Base(cwd) == "gitmap" || !isDirPath(filepath.Join(parent, ".git")) {
		if isDirPath(parent) {
			return filepath.Join(parent, name)
		}
	}
	return filepath.Join(cwd, name)
}

// defaultCloneTarget builds the clone destination: <cwd>/<configured name>.
func defaultCloneTarget(rec *store.SpecialRepoRecord) string {
	name := "repo-cache"
	if rec != nil && rec.ConfiguredName != "" {
		name = rec.ConfiguredName
	}
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	parent := filepath.Dir(cwd)
	if isDirPath(parent) {
		return filepath.Join(parent, name)
	}
	return filepath.Join(cwd, name)
}

// isDirPath reports whether path exists and is a directory.
func isDirPath(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// cloneTrainingRepo clones the training repository (repo-cache or repo-secrets);
// fails fast with guidance when gh is missing or unauthenticated.
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
		return apperror.WrapSimple(fmt.Errorf("%w: %s", err, strings.TrimSpace(combined.String())), "clone "+repoName)
	}
	return nil
}

// cloneRepoSecrets retains legacy name as an alias to cloneTrainingRepo.
func cloneRepoSecrets(log *decisionLog, dest string) *apperror.AppError {
	return cloneTrainingRepo(log, dest)
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

// pullTrainingRepo fast-forward pulls; skips with a note when no origin remote.
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

// pullRepoSecrets retains legacy name as an alias to pullTrainingRepo.
func pullRepoSecrets(log *decisionLog, path string) *apperror.AppError {
	return pullTrainingRepo(log, path)
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

// lastReadCheckpoint is the watermark of the last successful train run,
// tracked by gitmap (JSON under BinaryDataDir) — never in an MD file.
type lastReadCheckpoint struct {
	VaultPath string `json:"vault_path"`
	HeadSHA   string `json:"head_sha"`
	ReadAt    string `json:"read_at"`
}

// lastReadFile is the checkpoint path, beside gitmap's other JSON state.
func lastReadFile() string {
	return filepath.Join(store.BinaryDataDir(), "llm-train-muse-last-read.json")
}

// loadLastRead returns the checkpoint for vaultPath, or nil when absent,
// unreadable, or recorded for a different vault path.
func loadLastRead(vaultPath string) *lastReadCheckpoint {
	raw, err := os.ReadFile(lastReadFile())
	if err != nil {
		return nil
	}
	var cp lastReadCheckpoint
	if json.Unmarshal(raw, &cp) != nil {
		return nil
	}
	if cp.VaultPath != vaultPath || cp.HeadSHA == "" {
		return nil
	}
	return &cp
}

// saveLastRead records vaultPath@sha as the new watermark. A bookkeeping
// failure never fails the command — it only warns.
func saveLastRead(log *decisionLog, vaultPath, sha string) {
	cp := lastReadCheckpoint{VaultPath: vaultPath, HeadSHA: sha, ReadAt: time.Now().UTC().Format(time.RFC3339)}
	raw, err := json.MarshalIndent(cp, "", "  ")
	if err != nil {
		return
	}
	if err := os.MkdirAll(filepath.Dir(lastReadFile()), 0o755); err != nil {
		return
	}
	if err := os.WriteFile(lastReadFile(), raw, 0o600); err != nil {
		log.add("last-read: checkpoint save failed (%v) -> DECISION warn and continue (why: never fail a read over bookkeeping)", err)
		fmt.Printf("[LLM-TRAIN-MUSE] warning: could not save last-read checkpoint: %v\n", err)
		return
	}
	log.add("last-read: watermark saved -> %s @ %s (why: the next run diffs against this SHA; no MD file involved)", shortSHA(sha), cp.ReadAt)
}

// gitHeadSHA returns the vault's current HEAD SHA.
func gitHeadSHA(path string) (string, *apperror.AppError) {
	raw, err := runGitCmd(path, "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(raw)), nil
}

// shortSHA renders the 12-char display prefix of a SHA.
func shortSHA(sha string) string {
	if len(sha) > 12 {
		return sha[:12]
	}
	return sha
}

// sinceRead describes what changed between the checkpoint and now.
type sinceRead struct {
	state  string // "baseline" | "unchanged" | "fallback" | "diff"
	readAt string
	count  int
	files  []string
}

// diffSinceLastRead logs the watermark decision and diffs checkpoint..HEAD.
func diffSinceLastRead(log *decisionLog, path string, prev *lastReadCheckpoint, headSHA string) sinceRead {
	if prev == nil {
		log.add("last-read: no checkpoint yet -> DECISION record this run as the baseline watermark (why: the first sync establishes the diff base)")
		return sinceRead{state: "baseline"}
	}
	if prev.HeadSHA == headSHA {
		log.add("last-read: vault HEAD %s matches checkpoint from %s -> DECISION report 'nothing new since last read' (why: never re-train on unchanged state)", shortSHA(headSHA), prev.ReadAt)
		return sinceRead{state: "unchanged", readAt: prev.ReadAt}
	}
	raw, err := runGitCmd(path, "log", prev.HeadSHA+"..HEAD", "--name-only", "--pretty=format:")
	if err != nil {
		log.add("last-read: checkpoint %s not in history -> DECISION fall back to the --days window (why: rewritten history invalidates the watermark)", shortSHA(prev.HeadSHA))
		return sinceRead{state: "fallback", readAt: prev.ReadAt}
	}
	countRaw, countErr := runGitCmd(path, "rev-list", "--count", prev.HeadSHA+"..HEAD")
	if countErr != nil {
		log.add("last-read: checkpoint %s not in history -> DECISION fall back to the --days window (why: rewritten history invalidates the watermark)", shortSHA(prev.HeadSHA))
		return sinceRead{state: "fallback", readAt: prev.ReadAt}
	}
	log.add("last-read: checkpoint %s (%s) -> now %s -> DECISION diff the range for the since-last-read list (why: the model re-reads only what changed)", shortSHA(prev.HeadSHA), prev.ReadAt, shortSHA(headSHA))
	return sinceRead{state: "diff", readAt: prev.ReadAt, count: parseCount(countRaw), files: uniqueNonEmptyLines(string(raw))}
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
func printDigestSummary(path, digestDir, digestName string, days int, files []string, commitCount int, since sinceRead) {
	sequenced := listSequencedDigestFiles(digestDir, digestName)
	fmt.Println("======================================================================")
	fmt.Printf("GITMAP LLM TRAIN — MUSE MEMORY SYNC (last %d days)\n", days)
	fmt.Println("======================================================================")
	fmt.Printf("Digest root: %s\n", relForDisplay(digestDir))
	fmt.Printf("Commits: %d | Files changed: %d | Sequenced digests: %d\n", commitCount, len(files), len(sequenced))
	printSinceLastRead(since)
	fmt.Println()
	printDigestFileList("SEQUENCED DIGEST FILES", sequenced)
	printDigestFileList("CHANGED FILES", files)
	printTraversalChecklist(path, digestDir, sequenced)
}

// printSinceLastRead prints the watermark diff section — the checkpoint is
// tracked by gitmap, never in an MD file.
func printSinceLastRead(since sinceRead) {
	fmt.Println("SINCE LAST READ (tracked by gitmap — no MD file):")
	switch since.state {
	case "baseline":
		fmt.Println("  First run: this sync is the baseline watermark.")
	case "unchanged":
		fmt.Printf("  Vault unchanged since last read at %s — 0 new commits.\n", since.readAt)
	case "fallback":
		fmt.Println("  Checkpoint SHA not in history — showing the --days window instead.")
	case "diff":
		fmt.Printf("  %d new commit(s) since last read at %s:\n", since.count, since.readAt)
		if len(since.files) == 0 {
			fmt.Println("    (file list empty)")
		}
		for _, f := range since.files {
			fmt.Printf("    %s\n", f)
		}
	}
}

// listSequencedDigestFiles returns the digest-root .md files, digest-relative.
func listSequencedDigestFiles(digestDir, digestName string) []string {
	entries, err := os.ReadDir(digestDir)
	if err != nil {
		return nil
	}
	var files []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		files = append(files, digestName+"/"+entry.Name())
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
func printDigestMissingGuidance(path, digestName string) {
	fmt.Println("======================================================================")
	fmt.Println("GITMAP LLM TRAIN — MUSE MEMORY SYNC")
	fmt.Println("======================================================================")
	fmt.Printf("[LLM-TRAIN-MUSE] No %s/ digest directory at %s (No Muse/ digest directory).\n", digestName, relForDisplay(path))
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
