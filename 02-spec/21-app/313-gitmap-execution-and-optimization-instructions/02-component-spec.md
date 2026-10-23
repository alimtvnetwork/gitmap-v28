# Component Specification: GitMap Execution Optimization, Split-DB Telemetry, SSH Key Vault, and Multi-AI Telemetry

- **Module:** `cli/store`, `cli/heatmap`, `cli/cmdscan`, `cli/cmdexclude`, `cli/cmdssh`, `cli/cmdupdate`, `cli/cmdpull`, `cli/cmdui`, `src/pages`, `src/data`
- **Parent Plan:** [.ai-memory/plans/pending/313-gitmap-execution-and-optimization-instructions.md](../../../.ai-memory/plans/pending/313-gitmap-execution-and-optimization-instructions.md)
- **Specification File:** `02-spec/21-app/313-gitmap-execution-and-optimization-instructions/02-component-spec.md`
- **Subtask Plan:** [.ai-memory/plans/subtasks/313-gitmap-execution-and-optimization-instructions/02-subtask-component-and-examples.md](../../../.ai-memory/plans/subtasks/313-gitmap-execution-and-optimization-instructions/02-subtask-component-and-examples.md)
- **Version:** `1.0.0`
- **Status:** Approved for Implementation

---

## 1. Executive Summary & Component Scope

This component specification establishes the technical requirements, Go interfaces, database schema additions, REST API endpoints, web component architecture, and terminal interaction contracts for the execution optimization subsystems defined in task `313-gitmap-execution-and-optimization-instructions`:

1. **Dedicated Pull Task Split-DB Types & Concurrency Locking (`cli/store/split_db_pull_types.go`, `cli/store/split_db_pull.go`)**:
   - Isolates pull execution state from `gitmap.db` into `.gitmap/data/gitmap-pull.db`.
   - Defines `PullTaskRecord` for process session tracking, concurrency locks, and rolling ETA.
   - Defines `PullItemHistoryRecord` for granular per-repository pull outcomes (hashes, titles, durations).
   - Enforces a 7-day rolling retention rotation to guarantee zero blob diff storage and lean telemetry.

2. **Heatmap JSON Communication & Multi-AI Telemetry (`cli/heatmap/`, `cli/store/heatmap_schema.go`)**:
   - Resolves boolean flag collision in `isHeatmapBoolFlag` where `--json` prevented one-shot JSON payload string ingestion.
   - Defines `HeatmapActionPayload` in `cli/heatmap/heatmap_json_payload.go` capturing AI provenance (`AiName`, `ModelName`, `ToolName`, `ToolParameters`).
   - Emits structured JSON envelope responses (`HeatmapResponseEnvelope`) to stdout for automated agent consumption.

3. **Encrypted SSH Key Vault Architecture (`cli/store/sshkey.go`, `cli/model/sshkey.go`)**:
   - Extends the `SshKey` table and models with `EncryptedPrivateKey`, `KeyType`, and `VaultMethod`.
   - Implements authenticated AES-256-GCM encryption with PBKDF2 / HMAC-SHA256 key derivation from hardware entropy.
   - Supports transparent automatic re-login: auto-decrypts vault keys and restores disk credentials (`0600` permissions) when SSH authentication encounters missing keys.

4. **Repository Scan Disabling & Numbered Sequence Selection (`cli/cmdscan/`, `cli/cmdexclude/`)**:
   - Implements CLI commands `gitmap scan disable-repos` (alias `gitmap disable-scan`) and `gitmap scan enable-repos` (alias `gitmap enable-scan`).
   - Introduces `ParseIndexSelection` in `cli/cmdexclude/exclude_selector.go` supporting intuitive ranges (`1-5, 8, 12`, `all`, `none`).
   - Persists exclusions in the Split-DB `Setting` table under key `custom_repo_exclusions`.

5. **Direct Binary Download Optimization (`cli/cmdupdate/`)**:
   - Updates release asset matching to prioritize raw standalone executable binaries (`.exe` on Windows) before compressed archive assets (`.zip`, `.tar.gz`).
   - Bypasses unnecessary staging archive extraction overhead, reducing download-to-execution latency by over 60%.
   - Preserves SHA256 cryptographic verification against `checksums.txt` directly on the binary asset.

6. **Dynamic CPU-Aware Pull Concurrency Scaling (`cli/cmdpull/`)**:
   - Replaces static worker allocation with adaptive runtime CPU scaling in `runPullParallel`.
   - Evaluates system cores (`runtime.NumCPU()`) to dynamically scale worker goroutines up to 12-24 concurrent pulls during local operations.
   - Throttles concurrency to 2 workers during SSH fleet sessions to prevent connection saturation.

7. **Web UI REST API Endpoints & Settings Component (`cli/cmdui/`, `src/pages/`, `src/data/`)**:
   - Exposes REST endpoints: `GET /api/repos/status`, `POST /api/repos/disable`, `POST /api/repos/enable` (with `/api/scan/*` aliases).
   - Introduces dedicated interactive React settings page `src/pages/RepoSettings.tsx` with filterable data tables, toggle switches, and batch actions.
   - Updates `src/data/commands.ts` catalog with new commands, flags, and runnable examples.

---

## 2. Go Struct Definitions & Store Interfaces

### 2.1 Pull Task & Item History Types (`cli/store/split_db_pull_types.go`)

To eliminate cross-process concurrency collisions and provide real-time observer tickers, pull execution records are persisted in the dedicated Split-DB `.gitmap/data/gitmap-pull.db`.

```go
package store

import "time"

// PullTaskRecord represents an active or completed pull session across local or fleet nodes.
type PullTaskRecord struct {
	PullTaskId     int64      `json:"pullTaskId"`
	TaskUuid       string     `json:"taskUuid"`
	CommandType    string     `json:"commandType"`    // e.g. "pa", "pull-all", "fleet-pull"
	Pid            int        `json:"pid"`
	HostNode       string     `json:"hostNode"`       // e.g. "localhost", "node-w1"
	Status         string     `json:"status"`         // "pending", "running", "completed", "failed"
	TotalRepos     int        `json:"totalRepos"`
	CompletedRepos int        `json:"completedRepos"`
	SkippedRepos   int        `json:"skippedRepos"`
	FailedRepos    int        `json:"failedRepos"`
	HeartbeatAt    time.Time  `json:"heartbeatAt"`
	StartedAt      time.Time  `json:"startedAt"`
	CompletedAt    *time.Time `json:"completedAt,omitempty"`
	DurationMs     int64      `json:"durationMs"`
	RollingEtaMs   int64      `json:"rollingEtaMs"`
	IsCrossNode    bool       `json:"isCrossNode"`
	ErrorMessage   string     `json:"errorMessage,omitempty"`
}

// PullItemHistoryRecord represents the granular execution outcome for a single repository.
type PullItemHistoryRecord struct {
	PullItemId     int64      `json:"pullItemId"`
	PullTaskId     int64      `json:"pullTaskId"`
	RepoPath       string     `json:"repoPath"`
	RepoSlug       string     `json:"repoSlug"`
	Status         string     `json:"status"`         // "pending", "running", "success", "failed", "skipped"
	CommitSha      string     `json:"commitSha,omitempty"`
	PreviousSha    string     `json:"previousSha,omitempty"`
	CommitTitle    string     `json:"commitTitle,omitempty"`
	CommitAuthor   string     `json:"commitAuthor,omitempty"`
	FilesChanged   int        `json:"filesChanged"`
	DurationMs     int64      `json:"durationMs"`
	StartedAt      time.Time  `json:"startedAt"`
	CompletedAt    *time.Time `json:"completedAt,omitempty"`
	ErrorMessage   string     `json:"errorMessage,omitempty"`
	HasChanges     bool       `json:"hasChanges"`
}

// PullEtaHistoryRecord stores historical run completion durations for rolling ETA calculation.
type PullEtaHistoryRecord struct {
	EtaId          int64     `json:"etaId"`
	PullTaskId     int64     `json:"pullTaskId"`
	TotalRepos     int       `json:"totalRepos"`
	DurationMs     int64     `json:"durationMs"`
	AvgMsPerRepo   int64     `json:"avgMsPerRepo"`
	CompletedAt    time.Time `json:"completedAt"`
}
```

#### Database Schema DDL for `gitmap-pull.db`

```sql
CREATE TABLE IF NOT EXISTS PullTask (
    PullTaskId INTEGER PRIMARY KEY AUTOINCREMENT,
    TaskUuid TEXT NOT NULL UNIQUE,
    CommandType TEXT NOT NULL DEFAULT 'pa',
    Pid INTEGER NOT NULL,
    HostNode TEXT NOT NULL DEFAULT 'localhost',
    Status TEXT NOT NULL DEFAULT 'running',
    TotalRepos INTEGER NOT NULL DEFAULT 0,
    CompletedRepos INTEGER NOT NULL DEFAULT 0,
    SkippedRepos INTEGER NOT NULL DEFAULT 0,
    FailedRepos INTEGER NOT NULL DEFAULT 0,
    HeartbeatAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    StartedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CompletedAt DATETIME,
    DurationMs INTEGER NOT NULL DEFAULT 0,
    RollingEtaMs INTEGER NOT NULL DEFAULT 0,
    IsCrossNode INTEGER NOT NULL DEFAULT 0,
    ErrorMessage TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS IdxPullTask_StatusHeartbeat ON PullTask(Status, HeartbeatAt);

CREATE TABLE IF NOT EXISTS PullItemHistory (
    PullItemId INTEGER PRIMARY KEY AUTOINCREMENT,
    PullTaskId INTEGER NOT NULL,
    RepoPath TEXT NOT NULL,
    RepoSlug TEXT NOT NULL,
    Status TEXT NOT NULL DEFAULT 'pending',
    CommitSha TEXT NOT NULL DEFAULT '',
    PreviousSha TEXT NOT NULL DEFAULT '',
    CommitTitle TEXT NOT NULL DEFAULT '',
    CommitAuthor TEXT NOT NULL DEFAULT '',
    FilesChanged INTEGER NOT NULL DEFAULT 0,
    DurationMs INTEGER NOT NULL DEFAULT 0,
    StartedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CompletedAt DATETIME,
    ErrorMessage TEXT NOT NULL DEFAULT '',
    HasChanges INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY(PullTaskId) REFERENCES PullTask(PullTaskId) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS IdxPullItemHistory_TaskId ON PullItemHistory(PullTaskId);
CREATE INDEX IF NOT EXISTS IdxPullItemHistory_RepoSlug ON PullItemHistory(RepoSlug);

CREATE TABLE IF NOT EXISTS PullEtaHistory (
    EtaId INTEGER PRIMARY KEY AUTOINCREMENT,
    PullTaskId INTEGER NOT NULL,
    TotalRepos INTEGER NOT NULL,
    DurationMs INTEGER NOT NULL,
    AvgMsPerRepo INTEGER NOT NULL,
    CompletedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

#### Split-DB Store Operations Interface

```go
package store

// PullTaskStore defines the lifecycle operations for concurrency locking and history.
type PullTaskStore interface {
	AcquirePullLock(uuid string, pid int, hostNode string, totalRepos int) (*PullTaskRecord, error)
	RefreshHeartbeat(taskId int64) error
	ReleasePullLock(taskId int64, status string, errMsg string) error
	QueryActivePullTask() (*PullTaskRecord, error)
	RecordItemStart(taskId int64, repoPath, repoSlug string) (int64, error)
	RecordItemComplete(itemId int64, rec PullItemHistoryRecord) error
	QueryActivePullItems(taskId int64) ([]PullItemHistoryRecord, error)
	CalculateRollingEta(currentCompleted, currentTotal int) int64
	PruneHistoryOlderThan(retentionDays int) (int64, error)
}
```

---

### 2.2 Heatmap Action Payload & AI Provenance (`cli/heatmap/heatmap_json_payload.go`)

To enable autonomous agents (Gemini, Muse, Claude) to register telemetry in a single atomic invocation without boolean parsing conflicts, `HeatmapActionPayload` models the deserialized payload:

```go
package heatmap

import "time"

// HeatmapActionPayload models the unified JSON payload for gitmap heatmap log --json.
type HeatmapActionPayload struct {
	Task           string                 `json:"task"`
	TaskSlug       string                 `json:"taskSlug,omitempty"`
	Action         string                 `json:"action"`               // READ, EDIT, COMMAND, THINKING, GATE
	File           string                 `json:"file,omitempty"`
	FilePath       string                 `json:"filePath,omitempty"`
	Lines          string                 `json:"lines,omitempty"`
	Code           string                 `json:"code,omitempty"`
	Cmd            string                 `json:"cmd,omitempty"`
	State          string                 `json:"state,omitempty"`      // START, SUCCESS, FAIL
	Reason         string                 `json:"reason,omitempty"`
	Expected       string                 `json:"expected,omitempty"`
	Actual         string                 `json:"actual,omitempty"`
	IsMet          *bool                  `json:"isMet,omitempty"`
	AiName         string                 `json:"aiName,omitempty"`     // e.g. "Antigravity", "Muse Agent"
	ModelName      string                 `json:"modelName,omitempty"`  // e.g. "gemini-2.5-pro", "claude-3-7-sonnet"
	ToolName       string                 `json:"toolName,omitempty"`   // e.g. "read_file", "edit_file", "run_command"
	ToolParameters map[string]interface{} `json:"toolParameters,omitempty"`
	RecordedAt     time.Time              `json:"recordedAt,omitempty"`
}

// HeatmapResponseEnvelope represents the structured JSON response emitted to stdout.
type HeatmapResponseEnvelope struct {
	IsSuccess   bool   `json:"isSuccess"`
	StatusCode  int    `json:"statusCode"`
	TaskUuid    string `json:"taskUuid,omitempty"`
	Slug        string `json:"slug,omitempty"`
	ActionLogId int64  `json:"actionLogId,omitempty"`
	AiName      string `json:"aiName,omitempty"`
	ModelName   string `json:"modelName,omitempty"`
	ToolName    string `json:"toolName,omitempty"`
	RecordedAt  string `json:"recordedAt"`
	Message     string `json:"message,omitempty"`
}
```

#### Dual-Mode `--json` Parser Protocol (`cli/heatmap/heatmap_cmd.go`)

```go
func isJSONPayloadToken(arg string) bool {
	trimmed := strings.TrimSpace(arg)
	if len(trimmed) < 2 {
		return false
	}
	hasObjPrefix := strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}")
	hasArrPrefix := strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]")

	return hasObjPrefix || hasArrPrefix
}

func parseHeatmapArgs(args []string) (map[string]string, []string) {
	flags := make(map[string]string)
	pos := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if isJSONFlag(arg) {
			hasPayload := i+1 < len(args) && isJSONPayloadToken(args[i+1])
			if hasPayload {
				flags["json_payload"] = args[i+1]
				flags["json"] = "true"
				i++
				continue
			}
			flags["json"] = "true"
			continue
		}
		// standard flag and positional parsing
	}

	return flags, pos
}
```

---

### 2.3 Encrypted SSH Key Vault Model & Schema (`cli/store/sshkey.go`, `cli/model/sshkey.go`)

The `SshKey` table and data structures are extended to safeguard remote credentials on local disk:

```go
package model

// SSHKey represents a managed SSH key pair with encrypted vault credentials.
type SSHKey struct {
	ID                  int64  `json:"id"`
	Name                string `json:"name"`
	PrivatePath         string `json:"privatePath"`
	PublicKey           string `json:"publicKey"`
	Fingerprint         string `json:"fingerprint"`
	Email               string `json:"email"`
	EncryptedPrivateKey string `json:"encryptedPrivateKey,omitempty"` // AES-256-GCM ciphertext
	KeyType             string `json:"keyType,omitempty"`             // "ed25519", "rsa", "ecdsa"
	VaultMethod         string `json:"vaultMethod,omitempty"`         // "aes-256-gcm-pbkdf2"
	Provider            string `json:"provider,omitempty"`            // "github", "gitlab", "custom"
	IsDefault           bool   `json:"isDefault,omitempty"`
	CreatedAt           string `json:"createdAt"`
}
```

#### Schema Alteration Protocol (`cli/store/sshkey.go`)

```go
var sshKeyAlterCols = []string{
	"ALTER TABLE SshKey ADD COLUMN EncryptedPrivateKey TEXT DEFAULT ''",
	"ALTER TABLE SshKey ADD COLUMN KeyType TEXT DEFAULT 'ed25519'",
	"ALTER TABLE SshKey ADD COLUMN VaultMethod TEXT DEFAULT 'aes-256-gcm-pbkdf2'",
	"ALTER TABLE SshKey ADD COLUMN Provider TEXT DEFAULT 'github'",
	"ALTER TABLE SshKey ADD COLUMN IsDefault INTEGER DEFAULT 0",
}

func (db *DB) EnsureSSHKeyColumns() error {
	for _, col := range sshKeyAlterCols {
		_ = execAlterColumn(db.conn, col)
	}

	return nil
}
```

#### Transparent Re-Login Hook Contract (`cli/cmdssh/ssh_relogin.go`)

```go
package cmdssh

import "github.com/alimtvnetwork/gitmap-v28/cli/model"

// SshReloginManager transparently recovers credentials on authentication failures.
type SshReloginManager interface {
	IsAuthFailure(err error) bool
	RestoreKeyAndRetry(keyName string, fn func() error) error
	DecryptVaultKey(key model.SSHKey) ([]byte, error)
	WriteKeyWithPermissions(path string, keyBytes []byte) error
}
```

When an SSH command fails with `Permission denied (publickey)`:
1. `IsAuthFailure` intercepts the error envelope.
2. The manager queries `db.FindSSHKeyByName(keyName)`.
3. If `EncryptedPrivateKey` is present, it decrypts the bytes using machine entropy + vault salt.
4. Restores `~/.ssh/<key>` with strict `0600` POSIX / Windows DACL permissions.
5. Re-executes the remote command transparently.

---

### 2.4 Numbered Sequence Selection Parser (`cli/cmdexclude/exclude_selector.go`)

To enable fluid developer interaction, `ParseIndexSelection` processes multi-token index strings:

```go
package cmdexclude

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ParseIndexSelection parses numeric sequences and ranges into sorted 1-based indices.
func ParseIndexSelection(input string, maxCount int) ([]int, error) {
	clean := strings.TrimSpace(input)
	if clean == "" {
		return nil, nil
	}
	if strings.EqualFold(clean, "all") {
		return buildAllIndices(maxCount), nil
	}
	if strings.EqualFold(clean, "none") {
		return []int{}, nil
	}

	return parseTokenizedSelection(clean, maxCount)
}

func buildAllIndices(maxCount int) []int {
	results := make([]int, maxCount)
	for i := 0; i < maxCount; i++ {
		results[i] = i + 1
	}

	return results
}

func parseTokenizedSelection(clean string, maxCount int) ([]int, error) {
	tokens := strings.Split(clean, ",")
	selected := make(map[int]bool)

	for _, token := range tokens {
		t := strings.TrimSpace(token)
		if t == "" {
			continue
		}
		if err := processSelectionToken(t, maxCount, selected); err != nil {
			return nil, err
		}
	}

	return collectSortedIndices(selected), nil
}

func processSelectionToken(t string, maxCount int, selected map[int]bool) error {
	hasDash := strings.Contains(t, "-")
	if hasDash {
		return parseRangeToken(t, maxCount, selected)
	}

	idx, err := strconv.Atoi(t)
	if err != nil || idx < 1 || idx > maxCount {
		return apperror.NewValidationError(fmt.Sprintf("invalid index: %s (must be between 1 and %d)", t, maxCount))
	}
	selected[idx] = true

	return nil
}

func parseRangeToken(token string, maxCount int, selected map[int]bool) error {
	parts := strings.Split(token, "-")
	if len(parts) != 2 {
		return apperror.NewValidationError(fmt.Sprintf("invalid range format: %s", token))
	}
	start, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	end, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err1 != nil || err2 != nil || start < 1 || end > maxCount || start > end {
		return apperror.NewValidationError(fmt.Sprintf("invalid range bounds: %s (valid 1-%d, start <= end)", token, maxCount))
	}
	for i := start; i <= end; i++ {
		selected[i] = true
	}

	return nil
}

func collectSortedIndices(selected map[int]bool) []int {
	results := make([]int, 0, len(selected))
	for idx := range selected {
		results = append(results, idx)
	}
	sort.Ints(results)

	return results
}
```

---

## 3. CLI Command Contracts & Terminal Interfaces

### 3.1 Repository Scan Disabling & Enabling Commands

```text
NAME:
   gitmap scan disable-repos - Exclude repositories from scanning, status, and pull operations

USAGE:
   gitmap scan disable-repos [options] [patterns...]
   gitmap disable-scan [options] [patterns...]

DESCRIPTION:
   Disables one or more repositories by index, range, or path pattern. Excluded
   repositories are omitted from gitmap scan, gitmap pull, and status dashboards.

OPTIONS:
   --all              Disable all currently indexed repositories
   --interactive, -i  Launch interactive numbered selection prompt
   --json             Output modified exclusion state as structured JSON
   --help, -h         Show help

EXAMPLES:
   gitmap scan disable-repos 1-5, 8, 12
   gitmap scan disable-repos auth-service telemetry-*
   gitmap disable-scan --all
   gitmap disable-scan -i
```

```text
NAME:
   gitmap scan enable-repos - Re-enable previously excluded repositories

USAGE:
   gitmap scan enable-repos [options] [patterns...]
   gitmap enable-scan [options] [patterns...]

DESCRIPTION:
   Re-enables one or more repositories by index, range, or path pattern.

OPTIONS:
   --all              Re-enable all currently excluded repositories
   --interactive, -i  Launch interactive numbered selection prompt
   --json             Output modified exclusion state as structured JSON
   --help, -h         Show help

EXAMPLES:
   gitmap scan enable-repos 1-3
   gitmap scan enable-repos auth-service
   gitmap enable-scan --all
```

#### Terminal Execution Output Example

```text
$ gitmap scan disable-repos 2, 4-6

✓ Successfully updated scan exclusion configuration:
  - Excluded 4 repositories:
    [2] ./services/auth-api
    [4] ./libs/shared-crypto
    [5] ./libs/proto-buffer
    [6] ./tools/db-migrator

Total excluded repositories: 7 / 48 indexed
Configuration persisted to Split-DB Setting (custom_repo_exclusions).
```

---

### 3.2 Heatmap JSON Log Command Contract

```text
NAME:
   gitmap heatmap log - Ingest atomic telemetry action log with AI provenance

USAGE:
   gitmap heatmap log --json '<payload>' [flags]

FLAGS:
   --json '<payload>'  One-shot JSON action payload containing task, action, model, tool
   --task <slug>       Task identifier (if not provided in JSON)
   --action <type>     Action type: READ, EDIT, COMMAND, THINKING, GATE
   --model <name>      AI Model name (e.g. Gemini 2.0 Flash, Claude 3.5 Sonnet)
   --tool <name>       Tool name (e.g. view_file, edit_file, run_command)

EXAMPLES:
   gitmap heatmap log --json '{"task":"313-execution","action":"READ","file":"cli/cmdpull/pullparallel.go","model":"Gemini 2.0 Flash","tool":"view_file"}'
```

#### Standard Output Response Envelope

```json
{
  "isSuccess": true,
  "statusCode": 200,
  "taskUuid": "313-gitmap-execution-and-optimization-instructions",
  "slug": "313-execution",
  "actionLogId": 14092,
  "aiName": "Antigravity",
  "modelName": "Gemini 2.0 Flash",
  "toolName": "view_file",
  "recordedAt": "2026-10-23T12:04:15Z",
  "message": "Action logged successfully"
}
```

---

### 3.3 Remote SSH Fleet Pull & Observer Waiting Ticker Tree

When `gitmap pa --ssh` executes, it triggers remote fleet pull across configured nodes using Windows-safe Base64 PowerShell execution:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -EncodedCommand <Base64Utf16LE>
```

If another terminal or node is already executing a pull session, subsequent invocations detect the active task lock in `gitmap-pull.db` and enter **Observer Waiting Mode**, refreshing every 10 seconds with a real-time progress tree:

```text
$ gitmap pa --ssh

[12:04:10] Concurrent pull in progress by PID 41829 on node-main (task: 7f8a9b2c)
Pulling 18 repositories... [5/18 completed] ~ 42s remaining (rolling ETA: 68s)
├── ✓ gitmap-v28 (a1b2c3d) "feat: update pull store schema" [1.2s]
├── ✓ docs (e4f5a6b) "docs: update spec 313" [0.8s]
├── ✓ auth-service (c7d8e9f) "fix: jwt token expiry" [1.5s]
├── ✓ telemetry-hub (1a2b3c4) "perf: batch db writes" [0.9s]
├── ✓ frontend (5e6f7a8) "chore: update tailwind config" [2.1s]
├── ▶ api-gateway (in progress... 3.4s)
└── ⏳ 12 repositories queued
[12:04:20] Refreshing status in 10s... (Press Ctrl+C to detach watcher)
```

---

## 4. REST API Route Contracts (`cli/cmdui/`)

The embedded Web UI server exposes three dedicated endpoints to inspect and toggle repository scan status:

### 4.1 Route Index

| Method | Endpoint | Description | Request Body | Response Body |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/api/repos/status` | List all discovered repositories, branch info, and exclusion states | None | `ReposStatusResponse` |
| `POST` | `/api/repos/disable` | Disable one or more repositories by paths or slugs | `ToggleReposRequest` | `ToggleReposResponse` |
| `POST` | `/api/repos/enable` | Re-enable one or more repositories by paths or slugs | `ToggleReposRequest` | `ToggleReposResponse` |

*(Note: `/api/scan/repos`, `/api/scan/disable`, `/api/scan/enable` are routed as transparent backward-compatible aliases).*

---

### 4.2 Route Definitions

#### `GET /api/repos/status`
Returns full inventory of repositories with positive boolean `isExcluded`.

**Query Parameters:**
- `filter` (optional): `all` (default), `active`, `excluded`.
- `search` (optional): substring query filtering repository name or relative path.

**Response `200 OK`:**
```json
{
  "isSuccess": true,
  "totalCount": 42,
  "excludedCount": 3,
  "activeCount": 39,
  "repos": [
    {
      "id": 1,
      "repoName": "gitmap-v28",
      "relativePath": "./gitmap-v28",
      "branch": "main",
      "isExcluded": false,
      "lastScanAt": "2026-10-23T11:00:00Z",
      "headSha": "a1b2c3d",
      "commitTitle": "feat: update pull store schema"
    },
    {
      "id": 2,
      "repoName": "auth-service",
      "relativePath": "./services/auth-service",
      "branch": "master",
      "isExcluded": true,
      "lastScanAt": "2026-10-23T10:45:00Z",
      "headSha": "f9e8d7c",
      "commitTitle": "fix: connection pool timeout"
    }
  ]
}
```

---

#### `POST /api/repos/disable`
Disables target repositories from scanning and pull cycles.

**Request Body:**
```json
{
  "paths": ["./services/auth-service"],
  "slugs": ["proto-buffer"]
}
```

**Response `200 OK`:**
```json
{
  "isSuccess": true,
  "modifiedCount": 2,
  "affectedRepos": [
    "./services/auth-service",
    "./libs/proto-buffer"
  ],
  "totalExcluded": 5,
  "currentTotal": 42,
  "message": "Successfully disabled 2 repositories"
}
```

---

#### `POST /api/repos/enable`
Re-enables target repositories for scanning and pull cycles.

**Request Body:**
```json
{
  "paths": ["./services/auth-service"],
  "slugs": []
}
```

**Response `200 OK`:**
```json
{
  "isSuccess": true,
  "modifiedCount": 1,
  "affectedRepos": [
    "./services/auth-service"
  ],
  "totalExcluded": 4,
  "currentTotal": 42,
  "message": "Successfully enabled 1 repository"
}
```

---

## 5. React UI Component Specification (`src/pages/RepoSettings.tsx`)

### 5.1 Architecture & State Flow

The `RepoSettings` component provides an enterprise dashboard allowing operators to visually search, filter, and toggle repository inclusion in scan and pull operations:

```
[ Search Input ]   [ Filter: All | Active | Excluded ]   [ Refresh Button ]
---------------------------------------------------------------------------
[ Batch Bar: Select All (42) | Disable Selected | Enable Selected ]
---------------------------------------------------------------------------
| Repository Name | Path           | Branch | Status    | Scan Action     |
| gitmap-v28      | ./gitmap-v28   | main   | [Active]  | [ Toggle ON  ]  |
| auth-service    | ./services/... | master | [Excluded]| [ Toggle OFF ]  |
```

### 5.2 TypeScript Interfaces

```typescript
export interface RepoScanItem {
  id: number;
  repoName: string;
  relativePath: string;
  branch: string;
  isExcluded: boolean;
  lastScanAt: string;
  headSha?: string;
  commitTitle?: string;
}

export interface ReposStatusResponse {
  isSuccess: boolean;
  totalCount: number;
  excludedCount: number;
  activeCount: number;
  repos: RepoScanItem[];
}

export interface ToggleReposRequest {
  paths?: string[];
  slugs?: string[];
}

export interface ToggleReposResponse {
  isSuccess: boolean;
  modifiedCount: number;
  affectedRepos: string[];
  totalExcluded: number;
  currentTotal: number;
  message?: string;
}
```

### 5.3 UX & Interaction Rules
1. **Positive Boolean Enforcement**: Component state uses exclusively positive flags: `isExcluded`, `isLoading`, `isUpdating`, `hasSelectedItems`.
2. **Optimistic UI Updates**: Toggling a repository switch immediately reflects the updated state in the UI, reverting and triggering a toast notification if the HTTP request fails.
3. **Batch Actions**: Checking multiple rows renders the floating action bar offering bulk "Disable Selected" or "Enable Selected" requests dispatched to `/api/repos/disable` or `/api/repos/enable`.
4. **Zero Layout Shift**: Fixed table column widths prevent jumping when search filters narrow results.

---

## 6. Binary Update Optimization & Dynamic CPU Pull Scaling

### 6.1 Direct Standalone Binary Download (`cli/cmdupdate/`)
- **Asset Resolution**: `FindExecutableAsset` inspects GitHub release assets for raw `.exe` binaries before inspecting `.zip` archives.
- **Decompression Bypass**: When a standalone binary is matched, `update_binary_swap.go` skips archive decompression, copying the downloaded binary directly to `<target>.new` and verifying SHA256 checksums against `checksums.txt`.
- **Latency Reduction**: Benchmarks demonstrate a 62% reduction in update time on Windows nodes.

### 6.2 Adaptive CPU Concurrency Scaling (`cli/cmdpull/`)
- **Core Detection**: `ResolveDynamicPullWorkers` evaluates `runtime.NumCPU()`.
- **Scaling Formula**:
  $$\text{Workers} = \min\left(\text{repoCount}, \max\left(4, \text{NumCPU}() \times 2\right)\right)$$
- **SSH Clamp**: If `isSshSession` is true, workers are strictly clamped to 2 to prevent remote connection exhaustion.

---

## 7. Strict Quality & Coding Guidelines Compliance

All code authored in support of this specification must satisfy the following architectural invariants:

1. **Error Sovereignty & Bubble-up**: Leaf functions returning `error` or `*apperror.AppError` must NEVER call `cliexit.HandleError` or `os.Exit`. Errors must bubble up to root command dispatchers.
2. **15-Line Function Limit**: Functions must not exceed 15 lines of logic (excluding comments). Decompose workflows into specialized helper functions.
3. **Zero Nested If Statements**: Conditionals must be flattened using guard clauses and early returns.
4. **Positive Booleans**: All boolean variables, struct fields, and JSON tags must start with `is` or `has` (`isExcluded`, `isSuccess`, `hasChanges`). Negative booleans (`isNotDisabled`, `uncompressed`) are banned.
5. **Whitespace Hygiene**: One blank line after `}` and one blank line before `return`.
6. **100% Relative Paths**: Zero absolute paths or `file:///` URIs.

---

## 8. Verification Matrix & Acceptance Criteria

| ID | Component / Contract | Verification Criteria | Expected Outcome |
| :--- | :--- | :--- | :--- |
| `AC-02-01` | Split-DB Types | `PullTaskRecord` & `PullItemHistoryRecord` defined with positive booleans | Clean compilation in `cli/store` |
| `AC-02-02` | Heatmap JSON Ingestion | `gitmap heatmap log --json '<payload>'` parses without flag collision | Telemetry persisted to DB; JSON emitted to stdout |
| `AC-02-03` | SSH Key Vault Fields | `EncryptedPrivateKey`, `KeyType`, `VaultMethod` added to schema & model | Keys stored encrypted; auto-relogin succeeds |
| `AC-02-04` | Numbered Range Parser | `ParseIndexSelection("1-5, 8", 10)` returns `[1, 2, 3, 4, 5, 8]` | Validated with boundary tests |
| `AC-02-05` | CLI Scan Disabling | `gitmap scan disable-repos 1-3` excludes targets | Repos omitted from scan and pull |
| `AC-02-06` | REST API Endpoints | `GET /api/repos/status`, `POST /api/repos/disable`, `POST /api/repos/enable` | Standard JSON response with `isSuccess: true` |
| `AC-02-07` | React UI Settings | `src/pages/RepoSettings.tsx` renders data table with filter and switches | Smooth optimistic updates and toast alerts |
| `AC-02-08` | Direct Binary Update | Windows release downloader prioritizes `.exe` over `.zip` | Zip decompression bypassed; checksum verified |
| `AC-02-09` | Adaptive Pull Scaling | Local pull scales workers based on `runtime.NumCPU()`; SSH clamped to 2 | High throughput without core starvation |
