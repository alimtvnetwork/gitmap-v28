// Package llm — llm_train_muse_digest.go implements digest resolution,
// watermark checkpoints, and diff calculations.
package llm

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func logTrainingDecisions(log *decisionLog, daysN int) {
	sinceDate := time.Now().AddDate(0, 0, -daysN).Format("2006-01-02")
	log.add("window: --days %d -> since %s -> DECISION git log --since + rev-list --count (why: bound the diff to the digest window; the commit count frames the summary)", daysN, sinceDate)
	log.add("order: DECISION README.md first, then sequenced files in lexical order, chat-seeds on demand (why: contract before detail; a deterministic order any LLM can replay)")
	log.add("search: DECISION scope gitmap aum search to the digest tree with --ext .md (why: keep the model inside the digest, no repo-wide noise)")
}

func resolveMuseDigestDir(repoRoot string) (string, string, bool) {
	if !isDirPath(repoRoot) {
		return "", "", false
	}
	exact := filepath.Join(repoRoot, defaultDigestRootName)
	if isDirPath(exact) {
		return exact, defaultDigestRootName, true
	}
	entries, err := os.ReadDir(repoRoot)
	if err != nil {
		return "", "", false
	}
	return findMuseDirInEntries(repoRoot, entries)
}

func findMuseDirInEntries(repoRoot string, entries []os.DirEntry) (string, string, bool) {
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, "-muse-ai") || name == "muse-ai" || name == "muse" {
			return filepath.Join(repoRoot, name), name, true
		}
	}
	return "", "", false
}

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

	collectSpecialRepoCandidates(shortKey, key, addCandidate)
	collectLocalSpecialRepoCandidates(key, addCandidate)
	return candidates
}

func collectSpecialRepoCandidates(shortKey, key string, addCandidate func(string)) {
	db, err := store.OpenSpecialReposSplitDB()
	if err != nil {
		return
	}
	defer db.Close()
	addCandidate(querySpecialRepoPath(db, shortKey))
	addCandidate(querySpecialRepoPath(db, key))
}

func querySpecialRepoPath(db *store.SpecialReposSplitDB, key string) string {
	rec, recErr := db.GetSpecialRepo(key)
	if recErr != nil || rec == nil {
		return ""
	}
	return rec.LocalPath
}

func collectLocalSpecialRepoCandidates(key string, addCandidate func(string)) {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}
	for _, cand := range []string{
		filepath.Join(cwd, key),
		filepath.Join(filepath.Dir(cwd), key),
		filepath.Join("D:\\work", key),
	} {
		if isDirPath(cand) {
			addCandidate(cand)
		}
	}
}

type lastReadCheckpoint struct {
	VaultPath string `json:"vault_path"`
	HeadSHA   string `json:"head_sha"`
	ReadAt    string `json:"read_at"`
}

func lastReadFile() string {
	return filepath.Join(store.BinaryDataDir(), "llm-train-muse-last-read.json")
}

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

func saveLastRead(log *decisionLog, vaultPath, sha string) {
	cp := lastReadCheckpoint{VaultPath: vaultPath, HeadSHA: sha, ReadAt: time.Now().UTC().Format(time.RFC3339)}
	raw, err := json.MarshalIndent(cp, "", "  ")
	if err != nil {
		return
	}
	if mkErr := os.MkdirAll(filepath.Dir(lastReadFile()), 0o755); mkErr != nil {
		return
	}
	if writeErr := os.WriteFile(lastReadFile(), raw, 0o600); writeErr != nil {
		log.add("last-read: checkpoint save failed (%v) -> DECISION warn and continue (why: never fail a read over bookkeeping)", writeErr)
		fmt.Printf("[LLM-TRAIN-MUSE] warning: could not save last-read checkpoint: %v\n", writeErr)
		return
	}
	log.add("last-read: watermark saved -> %s @ %s (why: the next run diffs against this SHA; no MD file involved)", shortSHA(sha), cp.ReadAt)
}

func gitHeadSHA(path string) (string, *apperror.AppError) {
	raw, err := runGitCmd(path, "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(raw)), nil
}

func shortSHA(sha string) string {
	if len(sha) > 12 {
		return sha[:12]
	}
	return sha
}

type sinceRead struct {
	state  string
	readAt string
	count  int
	files  []string
}

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

func parseCount(raw []byte) int {
	count, err := strconv.Atoi(strings.TrimSpace(string(raw)))
	if err != nil || count < 0 {
		return 0
	}
	return count
}

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
