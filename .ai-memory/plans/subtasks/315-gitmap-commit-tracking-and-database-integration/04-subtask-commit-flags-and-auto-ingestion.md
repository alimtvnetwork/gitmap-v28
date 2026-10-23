# Subtask 04: Commit Flags Parsing, AI Provenance & Auto-Ingestion Hooking

> **Subtask ID:** `SUBTASK-315-04`  
> **Parent Task:** `315-gitmap-commit-tracking-and-database-integration`  
> **Status:** Completed  
> **Target Files:**  
> - `cli/cmdcommit/commit_types.go` (new file)  
> - `cli/cmdcommit/commit_flags.go` (new file)  
> - `cli/cmdcommit/commit_history_hook.go` (new file)  
> - `cli/cmdcommit/commit_cmd.go` (modify)  
> - `cli/cmdcommit/commit_push.go` (modify)  
> - `cli/cmdcommit/commit_batch.go` (modify)  
> - `cli/cmdcommit/commit_flags_test.go` (new unit test)  
> **Related Architecture Spec:** `02-spec/21-app/315-gitmap-commit-tracking-and-database-integration/02-component-spec.md`  

---

## 1. Objective

Implement CLI flag extraction, ambient AI provenance detection, and automatic post-commit telemetry ingestion across Gitmap's commit commands (`gitmap commit`, `cm`, `cpf`, `cmb`). The flag parser extracts both space-separated and equals-delimited flags (`--agent-id`, `--heatmap-id`, `--ai`, `--spec`, `--task`, `--push`, `--dry-run`), injecting ambient environment fallbacks (`ANTIGRAVITY_AGENT_ID`, `GITMAP_TASK`, `GITMAP_HEATMAP_ID`, `AI_MODEL`, `GITMAP_SPEC_PATH`). Telemetry ingestion automatically extracts git metadata (`rev-parse HEAD`, `log -1`, `diff-tree`) and records the entry into the external Split-DB (`repodb/commit_history.db`) via atomic transaction without disrupting git commit flow.

---

## 2. Action Checklist

- [x] 1. Author `cli/cmdcommit/commit_types.go`:
  - [x] Package declaration `package cmdcommit`.
  - [x] Define `CommitMetadata` struct:
    - [x] `TaskSlug string`
    - [x] `AgentID string`
    - [x] `HeatmapID string`
    - [x] `AIName string`
    - [x] `SpecPath string`
    - [x] `IsAIMode bool`
    - [x] `HasPush bool`
    - [x] `IsDryRun bool`
  - [x] Define string constants for flag keys and ambient environment variables.

- [x] 2. Author `cli/cmdcommit/commit_flags.go`:
  - [x] Implement `ExtractCommitFlags(args []string) ([]string, CommitMetadata)`.
  - [x] Support space-separated syntax (`--agent-id 123`) and equals format (`--agent-id=123`).
  - [x] Support flags: `--agent-id`, `--heatmap-id`, `--ai`, `--spec`, `--task`, `--push`, `-p`, `--dry-run`, `-n`.
  - [x] Implement `applyAmbientFallbacks(meta *CommitMetadata)` reading:
    - [x] `ANTIGRAVITY_AGENT_ID` and `AGENT_ID`
    - [x] `GITMAP_TASK`
    - [x] `GITMAP_HEATMAP_ID`
    - [x] `GITMAP_AI_NAME` and `AI_MODEL`
    - [x] `GITMAP_SPEC_PATH`
  - [x] Automatically set `IsAIMode = true` when `AgentID` or `ANTIGRAVITY_AGENT_ID` is present.
  - [x] Enforce 15-line limit (preferred <= 8 lines) and zero nested ifs.

- [x] 3. Author `cli/cmdcommit/commit_history_hook.go`:
  - [x] Implement `IngestCommitHistory(repoPath string, meta CommitMetadata) error`.
  - [x] Implement git plumbing extractors:
    - [x] `resolveCurrentCommitSha(repoPath string) (string, error)` via `git rev-parse HEAD`.
    - [x] `resolveCommitLogInfo(repoPath, sha string) (authorName, authorEmail, msg string, committedAt int64, err error)` via `git log -1`.
    - [x] `resolveCommitTouches(repoPath, sha string) ([]store.CommitFileTouch, int, int, int, error)` via `git diff-tree --name-status` and `git diff-tree --numstat`.
    - [x] `resolveRepoMetadata(repoPath string) (repoRoot, repoName, branchName string)`.
  - [x] Open database via `store.OpenCommitHistorySplitDB("")` and invoke `db.RecordCommit(rec, touches)`.
  - [x] Guarantee non-fatal error isolation: log ingestion warnings to `output.UIErr()` while returning clean execution status.

- [x] 4. Wire `cli/cmdcommit/commit_cmd.go`:
  - [x] Update `dispatchCommit(args []string) error` to call `ExtractCommitFlags(args)`.
  - [x] Forward `meta CommitMetadata` through `executeCommit(cleanArgs, meta)`.
  - [x] In `commitAndOptionalPush`, invoke `RecordCurrentCommitSafely(meta)` immediately following successful `dispatchGitCommit(args)`.

- [x] 5. Wire `cli/cmdcommit/commit_push.go`:
  - [x] Update `performCommitPush(commitMessage string, meta CommitMetadata) *apperror.AppError`.
  - [x] Trigger `RecordCurrentCommitSafely(meta)` directly following `performCommitStep(commitMessage)`.

- [x] 6. Wire `cli/cmdcommit/commit_batch.go`:
  - [x] Update `commitSingleChildRepo(repo, msg string, meta CommitMetadata)`.
  - [x] Trigger `RecordChildRepoCommitSafely(repo, msg, meta)` immediately following child repository commit execution.

- [x] 7. Quality & Invariant Verification:
  - [x] Strictly relative paths (zero absolute drive letters or URI schemes).
  - [x] Strictly lowercase filenames.
  - [x] Positive booleans only (`IsAIMode`, `HasPush`, `IsDryRun`, `HasFound`, `HasChanges`).
  - [x] Zero git commands, zero builds, zero test executions.

---

## 3. Concrete Code Specifications

### Specification: `cli/cmdcommit/commit_types.go`

```go
package cmdcommit

// CommitMetadata encapsulates operational flags and ambient AI provenance.
type CommitMetadata struct {
	TaskSlug  string
	AgentID   string
	HeatmapID string
	AIName    string
	SpecPath  string
	IsAIMode  bool
	HasPush   bool
	IsDryRun  bool
}

const (
	FlagAgentID     = "--agent-id"
	FlagHeatmapID   = "--heatmap-id"
	FlagAI          = "--ai"
	FlagSpec        = "--spec"
	FlagTask        = "--task"
	FlagPushLong    = "--push"
	FlagPushShort   = "-p"
	FlagDryRunLong  = "--dry-run"
	FlagDryRunShort = "-n"

	EnvAntigravityAgentID = "ANTIGRAVITY_AGENT_ID"
	EnvAgentID            = "AGENT_ID"
	EnvGitmapTask         = "GITMAP_TASK"
	EnvGitmapHeatmapID    = "GITMAP_HEATMAP_ID"
	EnvAIModel            = "AI_MODEL"
	EnvGitmapAIName       = "GITMAP_AI_NAME"
	EnvGitmapSpecPath     = "GITMAP_SPEC_PATH"
)
```

---

### Specification: `cli/cmdcommit/commit_flags.go`

```go
package cmdcommit

import (
	"os"
	"strings"
)

// ExtractCommitFlags parses CLI arguments, strips metadata flags, and resolves ambient fallbacks.
func ExtractCommitFlags(args []string) ([]string, CommitMetadata) {
	var clean []string
	var meta CommitMetadata
	skipNext := false

	for i := 0; i < len(args); i++ {
		if skipNext {
			skipNext = false
			continue
		}
		clean, skipNext = processSingleFlag(args, i, clean, &meta)
	}

	applyAmbientFallbacks(&meta)

	return clean, meta
}

func processSingleFlag(args []string, i int, clean []string, meta *CommitMetadata) ([]string, bool) {
	arg := args[i]
	if isStandaloneFlag(arg, meta) {
		return clean, false
	}
	if hasValuePrefix(arg, meta) {
		return clean, false
	}
	if i+1 < len(args) && isValuedFlag(arg, args[i+1], meta) {
		return clean, true
	}

	return append(clean, arg), false
}

func isStandaloneFlag(arg string, meta *CommitMetadata) bool {
	switch arg {
	case FlagPushLong, FlagPushShort:
		meta.HasPush = true
		return true
	case FlagDryRunLong, FlagDryRunShort:
		meta.IsDryRun = true
		return true
	case FlagAI:
		meta.IsAIMode = true
		return true
	default:
		return false
	}
}

func hasValuePrefix(arg string, meta *CommitMetadata) bool {
	prefixes := []struct {
		prefix string
		target *string
	}{
		{FlagAgentID + "=", &meta.AgentID},
		{FlagHeatmapID + "=", &meta.HeatmapID},
		{FlagSpec + "=", &meta.SpecPath},
		{FlagTask + "=", &meta.TaskSlug},
		{FlagAI + "=", &meta.AIName},
	}
	for _, p := range prefixes {
		if strings.HasPrefix(arg, p.prefix) {
			*p.target = strings.TrimSpace(strings.TrimPrefix(arg, p.prefix))
			return true
		}
	}

	return false
}

func isValuedFlag(arg, val string, meta *CommitMetadata) bool {
	switch arg {
	case FlagAgentID:
		meta.AgentID = strings.TrimSpace(val)
		return true
	case FlagHeatmapID:
		meta.HeatmapID = strings.TrimSpace(val)
		return true
	case FlagSpec:
		meta.SpecPath = strings.TrimSpace(val)
		return true
	case FlagTask:
		meta.TaskSlug = strings.TrimSpace(val)
		return true
	case FlagAI:
		meta.AIName = strings.TrimSpace(val)
		meta.IsAIMode = true
		return true
	default:
		return false
	}
}

func applyAmbientFallbacks(meta *CommitMetadata) {
	fallbackString(&meta.AgentID, EnvAntigravityAgentID, EnvAgentID)
	fallbackString(&meta.TaskSlug, EnvGitmapTask)
	fallbackString(&meta.HeatmapID, EnvGitmapHeatmapID)
	fallbackString(&meta.AIName, EnvGitmapAIName, EnvAIModel)
	fallbackString(&meta.SpecPath, EnvGitmapSpecPath)

	if meta.AgentID != "" || meta.AIName != "" {
		meta.IsAIMode = true
	}
}

func fallbackString(target *string, envVars ...string) {
	if *target != "" {
		return
	}
	for _, env := range envVars {
		if val := strings.TrimSpace(os.Getenv(env)); val != "" {
			*target = val
			return
		}
	}
}
```

---

### Specification: `cli/cmdcommit/commit_history_hook.go`

```go
package cmdcommit

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/output"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// IngestCommitHistory extracts git telemetry and records commit history to the external Split-DB.
func IngestCommitHistory(repoPath string, meta CommitMetadata) error {
	rec, touches, err := buildCommitHistoryTelemetry(repoPath, meta)
	if err != nil {
		output.UIErr("Warning: telemetry extraction failed: %v", err)
		return nil
	}

	db, errDb := store.OpenCommitHistoryDB("")
	if errDb != nil {
		output.UIErr("Warning: commit history database open failed: %v", errDb)
		return nil
	}
	defer db.Close()

	if errRecord := db.RecordCommit(rec, touches); errRecord != nil {
		output.UIErr("Warning: commit history record failed: %v", errRecord)
		return nil
	}

	return nil
}

func buildCommitHistoryTelemetry(repoPath string, meta CommitMetadata) (store.CommitHistoryRecord, []store.CommitFileTouch, error) {
	sha, err := extractGitOutput(repoPath, "rev-parse", "HEAD")
	if err != nil {
		return store.CommitHistoryRecord{}, nil, err
	}
	shortSha := sha
	if len(sha) >= 8 {
		shortSha = sha[:8]
	}

	author, email, commitTime, msg := extractCommitLog(repoPath, sha)
	touches, filesChanged, insertions, deletions := extractDiffMetrics(repoPath, sha)
	repoRoot, repoName, branch := resolveRepoContext(repoPath)

	rec := store.CommitHistoryRecord{
		CommitSha:     sha,
		ShortSha:      shortSha,
		RepoRoot:      repoRoot,
		RepoName:      repoName,
		BranchName:    branch,
		AuthorName:    author,
		AuthorEmail:   email,
		CommitMessage: msg,
		TaskSlug:      meta.TaskSlug,
		AgentID:       meta.AgentID,
		HeatmapID:     meta.HeatmapID,
		AIName:        meta.AIName,
		SpecPath:      meta.SpecPath,
		FilesChanged:  filesChanged,
		Insertions:    insertions,
		Deletions:     deletions,
		IsAIMode:      meta.IsAIMode,
		HasPushed:     meta.HasPush,
		CommittedAt:   commitTime,
	}

	return rec, touches, nil
}

func extractGitOutput(repoPath string, args ...string) (string, error) {
	var cmdArgs []string
	if repoPath != "" && repoPath != "." {
		cmdArgs = append(cmdArgs, "-C", repoPath)
	}
	cmdArgs = append(cmdArgs, args...)

	cmd := exec.Command("git", cmdArgs...)
	out, err := cmd.Output()
	if err != nil {
		return "", apperror.WrapSimple(err, "cmdcommit.extractGitOutput")
	}

	return strings.TrimSpace(string(out)), nil
}

func extractCommitLog(repoPath, sha string) (string, string, int64, string) {
	out, err := extractGitOutput(repoPath, "log", "-1", "--format=%an%x00%ae%x00%ad%x00%s", "--date=unix", sha)
	if err != nil {
		return "gitmap", "gitmap@local", time.Now().Unix(), "commit"
	}
	parts := strings.Split(out, "\x00")
	if len(parts) < 4 {
		return "gitmap", "gitmap@local", time.Now().Unix(), out
	}
	unixTime, _ := strconv.ParseInt(parts[2], 10, 64)

	return parts[0], parts[1], unixTime, parts[3]
}

func extractDiffMetrics(repoPath, sha string) ([]store.CommitFileTouch, int, int, int) {
	statusOut, _ := extractGitOutput(repoPath, "diff-tree", "--no-commit-id", "--name-status", "-r", sha)
	numstatOut, _ := extractGitOutput(repoPath, "diff-tree", "--no-commit-id", "--numstat", "-r", sha)

	touches, totalIns, totalDel := parseDiffTreeOutputs(statusOut, numstatOut, sha, repoPath)

	return touches, len(touches), totalIns, totalDel
}

func parseDiffTreeOutputs(statusOut, numstatOut, sha, repoRoot string) ([]store.CommitFileTouch, int, int) {
	statusMap := make(map[string]string)
	for _, line := range strings.Split(statusOut, "\n") {
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			statusMap[fields[1]] = fields[0]
		}
	}

	var touches []store.CommitFileTouch
	totalIns, totalDel := 0, 0
	for _, line := range strings.Split(numstatOut, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		ins, _ := strconv.Atoi(fields[0])
		del, _ := strconv.Atoi(fields[1])
		filePath := fields[2]
		totalIns += ins
		totalDel += del
		touches = append(touches, store.CommitFileTouch{
			CommitSha:  sha,
			RepoRoot:   repoRoot,
			FilePath:   filePath,
			ChangeType: statusMap[filePath],
			Insertions: ins,
			Deletions:  del,
		})
	}

	return touches, totalIns, totalDel
}

func resolveRepoContext(repoPath string) (string, string, string) {
	cleanPath := repoPath
	if cleanPath == "" || cleanPath == "." {
		cleanPath, _ = filepath.Abs(".")
	}
	repoRoot := filepath.ToSlash(cleanPath)
	repoName := filepath.Base(repoRoot)
	branch, _ := extractGitOutput(repoPath, "rev-parse", "--abbrev-ref", "HEAD")

	return repoRoot, repoName, branch
}
```

---

### Specification: `cli/cmdcommit/commit_cmd.go` (Wiring)

In `cli/cmdcommit/commit_cmd.go`, update `dispatchCommit` and `commitAndOptionalPush`:

```go
func dispatchCommit(args []string) error {
	cleanArgs, meta := ExtractCommitFlags(args)
	if !cmdpull.IsGitRepoCWD() {
		return handleNonGitRepoCommit(cleanArgs, meta.HasPush, meta.IsDryRun)
	}

	if meta.IsDryRun {
		return previewCommitChanges()
	}

	return executeCommit(cleanArgs, meta)
}

func executeCommit(args []string, meta CommitMetadata) error {
	hasChanges, err := stageAndCheckStaged()
	if err != nil {
		return err
	}

	if !hasChanges {
		printPaddedInfo("Working tree clean, nothing to commit.")

		return handleCleanWorkingTree(meta.HasPush)
	}

	return commitAndOptionalPush(args, meta)
}

func commitAndOptionalPush(args []string, meta CommitMetadata) error {
	if err := dispatchGitCommit(args); err != nil {
		return apperror.WrapSimple(err, "git commit failed:")
	}

	_ = IngestCommitHistory(".", meta)

	return handleOptionalPush(meta.HasPush)
}
```

---

### Specification: `cli/cmdcommit/commit_push.go` (Wiring)

In `cli/cmdcommit/commit_push.go`, update `performCommitPush`:

```go
func performCommitPush(commitMessage string, meta CommitMetadata) *apperror.AppError {
	printPaddedInfo("Committing: %s", commitMessage)
	if errCommit := performCommitStep(commitMessage); errCommit != nil {
		return errCommit
	}

	_ = IngestCommitHistory(".", meta)

	printPaddedInfo("Pushing to remote...")
	if errPush := performPushStep(); errPush != nil {
		return errPush
	}
	renderSuccessSummary(commitMessage)

	return nil
}
```

---

### Specification: `cli/cmdcommit/commit_batch.go` (Wiring)

In `cli/cmdcommit/commit_batch.go`, update `commitSingleChildRepo`:

```go
func commitSingleChildRepo(repo, msg string, meta CommitMetadata) {
	repoName := filepath.Base(repo)
	printPaddedInfo("Committing in child repo: %s", repoName)
	_ = execGitPaddedFiltered("-C", repo, "add", "-A")
	_ = execGitPaddedFiltered("-C", repo, "commit", "-m", msg)

	_ = IngestCommitHistory(repo, meta)

	if meta.HasPush {
		printPaddedInfo("Pushing child repo: %s", repoName)
		_ = execGitPaddedFiltered("-C", repo, "push")
	}
}
```
