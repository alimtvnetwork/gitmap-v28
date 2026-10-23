# Subtask 313-01: Core Architecture, Contracts & Concurrency Specifications

> **Subtask ID:** `313-01`  
> **Sequence ID:** `313`  
> **Status:** `ready`  
> **Master Plan:** [313-gitmap-execution-and-optimization-instructions.md](../../pending/313-gitmap-execution-and-optimization-instructions.md)  
> **Governing Spec:** [01-architecture-spec.md](../../../02-spec/21-app/313-gitmap-execution-and-optimization-instructions/01-architecture-spec.md)  
> **Assigned Role:** `Spec Author 01`  
> **Target Files (Owned):**  
> - `02-spec/21-app/313-gitmap-execution-and-optimization-instructions/01-architecture-spec.md`  
> - `.ai-memory/plans/subtasks/313-gitmap-execution-and-optimization-instructions/01-subtask-architecture-and-contracts.md`  

---

## 1. Objective & Strategic Scope

Formalize and document all structural contracts, Go interface signatures, SQLite Split-DB DDL, cryptographic data models, and terminal visualization protocols necessary to execute the 12 actionable items of Plan 313 without ambiguity or cross-agent collision.

This contract provides strict architectural invariants for downstream implementation workers (`Worker 01`, `Worker 02`, `Worker 03`):
- **Windows SSH Execution:** Bypassing `cmd.exe /c` quote-stripping and Win32 `ERROR_INVALID_NAME` (code 123) on remote node `w1` via Base64 UTF-16LE PowerShell encapsulation and `RunCommandSeparated`.
- **Dedicated Pull Split-DB (`gitmap-pull.db`):** Single-writer WAL mode schema for `PullTask`, `PullItemHistory`, and `PullEtaHistory` with PascalCase columns and positive boolean flags (`isActive`, `isCompleted`, `isFailed`, `hasActiveLock`, `isSuccess`, `hasChanges`, `isSkipped`).
- **Concurrency Locking & 10s Observer Loop:** Atomic machine lock state machine with 30s heartbeat expiry, stale lock stealing, and 10-second ticker tree-view progress rendering.
- **7-8 Run Rolling ETA Engine:** Real-time remaining duration estimation based on the last 7 to 8 runs.
- **Encrypted SSH Key Vault:** RSA-OAEP + AES-256-GCM envelope encryption persisting `EncryptedPrivateKey` in `SshKey` and `SSHConnection` tables with automatic re-login reconstitution upon Git authentication failure.
- **Heatmap Telemetry:** Atomic `--json '<payload>'` string ingestion with AI model (`Gemini`, `Antigravity IDE`, `Muse`) and tool metadata attribution.
- **Scan Disable/Enable:** Sequence-indexed CLI commands (`1-5, 8`), REST API routes, and React UI.

---

## 2. Concrete Structural Contracts & Go Type Signatures

### 2.1 Windows SSH Fleet Encapsulation & Stream Separation Contract (`cli/secrets/ssh_client.go`, `cli/cmdssh/`)

#### Function Signatures & Types
```go
package secrets

import (
	"golang.org/x/crypto/ssh"
)

// WrapPowerShellEncodedCommand transforms a PowerShell script into a Base64 UTF-16LE encoded command string.
// Eliminates quotes, semicolons, and metacharacters for safe execution under Windows OpenSSH cmd.exe /c.
func WrapPowerShellEncodedCommand(cmd string) string

// ResolveRemoteShell returns the canonical shell identifier ("ps", "bash", "sh") for a target OS string.
func ResolveRemoteShell(osType string) string

// RunCommandSeparated executes a command over SSH, routing stdout and stderr to discrete buffers.
// Prevents stderr diagnostic noise, PowerShell banners, or CLIXML objects from corrupting JSON payloads on stdout.
func RunCommandSeparated(client *ssh.Client, cmd, shellType string) (string, string, error)
```

#### Contract Specification
- **Input:** Any multi-line or chained PowerShell command string (e.g. `$env:GCM_CREDENTIAL_STORE='none'; gitmap pa --json`).
- **Processing:** Encodes command characters into Little-Endian UTF-16, then converts to standard Base64.
- **Output:** `powershell -NoProfile -ExecutionPolicy Bypass -EncodedCommand <Base64Payload>`.
- **Execution Invariant:** When `shellType == "ps"`, `secrets.RunCommand` and `secrets.RunCommandSeparated` invoke `WrapPowerShellEncodedCommand`. The resulting output is guaranteed to produce zero `The filename, directory name, or volume label syntax is incorrect.` errors under `cmd.exe /c`.
- **Stream Separation Invariant:** In `cli/cmdssh/ssh_pull_json.go`, `executeRemotePullWorker` invokes `RunCommandSeparated`. Stdout is used exclusively for `json.Unmarshal`. Stderr is captured and logged only if an error occurs.

---

### 2.2 Dedicated Pull Task Split-DB & Concurrency Lock Contract (`cli/store/`, `cli/cmdpull/`)

#### SQLite Single-Writer WAL Configuration
Managed in `cli/store/split_db_pull.go`:
```sql
PRAGMA journal_mode = WAL;
PRAGMA busy_timeout = 5000;
PRAGMA synchronous = NORMAL;
PRAGMA foreign_keys = ON;
```

#### Go Data Structures (PascalCase & Positive Booleans)
```go
package store

// PullTaskRecord represents an execution session row in PullTask.
type PullTaskRecord struct {
	PullTaskId    int64  `json:"pullTaskId"`
	TaskUuid      string `json:"taskUuid"`
	CommandType   string `json:"commandType"`
	TargetNode    string `json:"targetNode"`
	MachineId     string `json:"machineId"`
	Status        string `json:"status"` // 'pending', 'running', 'completed', 'failed', 'cancelled'
	Pid           int    `json:"pid"`
	TotalRepos    int    `json:"totalRepos"`
	PulledRepos   int    `json:"pulledRepos"`
	EtaSeconds    int    `json:"etaSeconds"`
	IsActive      bool   `json:"isActive"`      // Positive boolean
	IsCompleted   bool   `json:"isCompleted"`   // Positive boolean
	IsFailed      bool   `json:"isFailed"`      // Positive boolean
	HasActiveLock bool   `json:"hasActiveLock"` // Positive boolean
	StartedAt     string `json:"startedAt"`
	CompletedAt   string `json:"completedAt"`
	ErrorMessage  string `json:"errorMessage"`
	HeartbeatAt   string `json:"heartbeatAt"`
	CreatedAt     string `json:"createdAt"`
}

// PullItemHistoryRecord represents a repository-level progress row in PullItemHistory.
type PullItemHistoryRecord struct {
	PullItemHistoryId int64  `json:"pullItemHistoryId"`
	PullTaskId        int64  `json:"pullTaskId"`
	RepoSlug          string `json:"repoSlug"`
	RepoPath          string `json:"repoPath"`
	PullStatus        string `json:"pullStatus"` // 'pending', 'pulling', 'success', 'skipped', 'failed'
	CommitSha         string `json:"commitSha"`
	CommitTitle       string `json:"commitTitle"`
	FilesChanged      int    `json:"filesChanged"`
	DurationMs        int64  `json:"durationMs"`
	IsSuccess         bool   `json:"isSuccess"`  // Positive boolean
	HasChanges        bool   `json:"hasChanges"` // Positive boolean
	IsSkipped         bool   `json:"isSkipped"`  // Positive boolean
	ErrorMessage      string `json:"errorMessage"`
	PulledAt          string `json:"pulledAt"`
	CreatedAt         string `json:"createdAt"`
}

// PullEtaHistoryRecord models a historical duration entry in PullEtaHistory.
type PullEtaHistoryRecord struct {
	PullEtaHistoryId int64   `json:"pullEtaHistoryId"`
	CommandType      string  `json:"commandType"`
	TargetNode       string  `json:"targetNode"`
	TotalRepos       int     `json:"totalRepos"`
	DurationSeconds  float64 `json:"durationSeconds"`
	CreatedAt        string  `json:"createdAt"`
}
```

#### Split-DB Interface Methods on `*PullSplitDB`
```go
package store

// TryAcquirePullTaskLock attempts to acquire the execution lock for a node.
// Returns existing active task if lock is held, or newly acquired task record.
func (s *PullSplitDB) TryAcquirePullTaskLock(cmdType, targetNode string, totalRepos int) (*PullTaskRecord, bool, error)

// UpdateTaskHeartbeat refreshes the HeartbeatAt timestamp to prevent stale lock recovery.
func (s *PullSplitDB) UpdateTaskHeartbeat(taskID int64) error

// RecordRepoPullResult inserts progress for a single repository under an active task.
func (s *PullSplitDB) RecordRepoPullResult(rec PullItemHistoryRecord) error

// CompletePullTask marks the task as completed, logs history to PullEtaHistory, and releases the lock.
func (s *PullSplitDB) CompletePullTask(taskID int64, durationMs int64) error

// FailPullTask marks the task as failed with an explanatory error message and releases the lock.
func (s *PullSplitDB) FailPullTask(taskID int64, errMsg string) error

// CalculateRollingEta projects remaining duration in seconds based on the last 7-8 runs.
func (s *PullSplitDB) CalculateRollingEta(cmdType, targetNode string, remainingRepos int) (int, error)

// PruneExpiredPullHistory removes PullTask and PullItemHistory records older than retentionDays (default: 7).
func (s *PullSplitDB) PruneExpiredPullHistory(retentionDays int) (int64, error)
```

#### 10-Second Observer Ticker Loop Contract (`cli/cmdpull/pull_watcher.go`)
```go
package cmdpull

import (
	"context"
	"time"
)

// WatchPullTaskProgress blocks the secondary terminal, polling PullTask and PullItemHistory every 10 seconds.
// Renders live tree-view progress until task reaches 'completed' or 'failed' status.
func WatchPullTaskProgress(ctx context.Context, taskID int64) error
```

---

### 2.3 Encrypted SSH Key Vault Contract (RSA-OAEP + AES-256-GCM Envelope Encryption)

#### Go Data Structures
```go
package store

// EncryptedKeyEnvelope encapsulates an encrypted private key using envelope encryption.
type EncryptedKeyEnvelope struct {
	KeyNonce       string `json:"keyNonce"`       // Base64 12-byte random nonce
	Ciphertext     string `json:"ciphertext"`     // Base64 AES-256-GCM ciphertext + 16-byte auth tag
	WrappedDataKey string `json:"wrappedDataKey"` // Base64 RSA-OAEP wrapped 256-bit DEK
	KeyVersion     int    `json:"keyVersion"`     // Key schema version (1)
	MasterKeySalt  string `json:"masterKeySalt"`  // Machine derivation salt
}
```

#### Vault Methods on `*DB` / `cmdssh`
```go
package store

// EncryptPrivateKeyEnvelope secures plaintext private key bytes using RSA-OAEP + AES-256-GCM.
func EncryptPrivateKeyEnvelope(plainKey []byte) (*EncryptedKeyEnvelope, error)

// DecryptPrivateKeyEnvelope decrypts an EncryptedKeyEnvelope using the machine-bound master key.
func DecryptPrivateKeyEnvelope(env *EncryptedKeyEnvelope) ([]byte, error)

// EnsureSSHKeyOnDisk checks if the key exists on disk; if missing, restores from encrypted vault.
func EnsureSSHKeyOnDisk(keyName string) (string, error)

// ExecuteTransparentRelogin handles SSH / Git authentication failures by restoring keys and re-authenticating.
func ExecuteTransparentRelogin(keyName string) error
```

#### Cryptographic Parameters
- **Symmetric Cipher:** AES-256-GCM with 96-bit (12-byte) random nonce from `crypto/rand`.
- **Asymmetric Wrapper:** RSA-OAEP with SHA-256.
- **Master Key Derivation:** Machine UUID + User SID + GitMap Salt.
- **Disk File Mode:** POSIX `0600` (`-rw-------`), Windows Restricted User DACL.
- **Tables Updated:** `SshKey` and `SSHConnection`.

---

### 2.4 Heatmap Single-Shot `--json '<payload>'` & AI Provenance Contract (`cli/heatmap/`, `cli/types/`)

#### Request Ingestion & Response Envelope
```go
package heatmap

// HeatmapLogJSONInput models single-shot JSON string input for gitmap heatmap log and agent heatmap.
type HeatmapLogJSONInput struct {
	Slug        string `json:"slug"`
	TaskUuid    string `json:"taskUuid"`
	SubtaskName string `json:"subtaskName"`
	Status      string `json:"status"`
	FilesRead   string `json:"filesRead"`
	FilesEdited string `json:"filesEdited"`
	ModelName   string `json:"modelName"` // e.g. "Gemini", "Antigravity IDE", "Muse", "gemini-1.5-pro"
	ToolName    string `json:"toolName"`  // e.g. "antigravity-ide", "cursor", "cli-runner"
	Roadmap     string `json:"roadmap"`
}

// HeatmapResponseEnvelope represents standard JSON response emitted to stdout.
type HeatmapResponseEnvelope struct {
	IsSuccess  bool   `json:"isSuccess"`
	HasChanges bool   `json:"hasChanges"`
	StatusCode int    `json:"statusCode"`
	TaskUuid   string `json:"taskUuid"`
	Slug       string `json:"slug"`
	ModelName  string `json:"modelName"`
	ToolName   string `json:"toolName"`
	RecordedAt string `json:"recordedAt"`
	Message    string `json:"message,omitempty"`
}
```

---

### 2.5 Scan Disable/Enable CLI & REST API Contract (`cli/cmdscan/`, `cli/cmdui/`)

#### REST API Endpoints
- `POST /api/repos/scan/toggle`:
  ```json
  {
    "slug": "kubernetes-training-v1",
    "isScanEnabled": false
  }
  ```
- `GET /api/repos/scan/disabled`: Returns `["kubernetes-training-v1", ...]`
- `POST /api/repos/scan/batch-toggle`:
  ```json
  {
    "slugs": ["story-telling-v1", "kubernetes-training-v1"],
    "isScanEnabled": true
  }
  ```

#### CLI Syntax & Sequence Parser
- `gitmap scan disable-repos` / `gitmap disable-scan`: Displays numeric table (`[1] story-telling-v1`, `[2] kubernetes-training-v1`).
- `gitmap disable-scan 1-5, 8`: Parses comma-separated ranges and indices, resolving them to repository slugs.
- `gitmap enable-scan 1-5, 8`: Re-enables scanning for selected repositories.

---

## 3. Disjoint Implementation Breakdown for Worker Swarm

| Task Code | Subtask Title | Assigned Role | Target Files |
| :--- | :--- | :--- | :--- |
| `Task-02` | Windows SSH Fleet Quoting & Command Dispatch Fix | Worker 01 | `cli/cmdssh/ssh_pull_fleet_dispatch.go`, `cli/cmdssh/ssh_pull_json.go`, `cli/secrets/ssh_client.go` |
| `Task-03` | Pull Task Split-DB, Lock/Wait Engine & Dynamic ETA | Worker 02 | `cli/store/split_db_pull.go`, `cli/cmdpull/pull_watcher.go`, `cli/cmdpull/pull_eta.go`, `cli/constants/constants_pull.go` |
| `Task-04` | Encrypted SSH Key Vault & Auto-Relogin Engine | Worker 03 | `cli/store/sshkey.go`, `cli/store/migrations_ssh.go`, `cli/cmdssh/ssh_vault.go`, `cli/cmdssh/ssh_relogin.go` |
| `Task-05` | Heatmap `--json` Ingestion & AI Telemetry Attribution | Worker 01 | `cli/heatmap/heatmap_cmd.go`, `cli/heatmap/types.go`, `cli/heatmap/heatmap_render_json.go`, `cli/store/heatmap_schema.go` |
| `Task-06` | Heatmap Skill & Multi-AI Guidelines Documentation | Worker 02 | `.agents/skills/gitmap/skill.md`, `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`, `.agents/skills/muse-master-prompt/skill.md` |
| `Task-07` | CPU-Aware Adaptive Pull Concurrency Scaling | Worker 03 | `cli/cmdpull/pull_concurrency.go`, `cli/cmdpull/pull_worker_pool.go` |
| `Task-08` | Direct Binary Update Download Optimization | Worker 01 | `cli/cmdupdate/update_fallback_loop.go`, `cli/cmdupdate/update_binary_swap.go` |
| `Task-09` | Repository Disable/Enable Scan CLI & Web REST UI | Worker 02 | `cli/cmdexclude/exclude_toggle.go`, `cli/cmdscan/scan_disable.go`, `cli/cmdui/ui_server.go`, `src/pages/RepoSettings.tsx` |

---

## 4. Verification Evidence & Acceptance Contracts

1. **Relative Path Hygiene:**  
   Zero absolute Windows or POSIX paths and zero file URI schemes in any authored specification or plan file.
2. **Forbidden Strings Inspection:**  
   Verification passes with 0 occurrences of magic error numbers, negative booleans, or unhandled exits in leaf methods.
3. **Coding Guidelines Adherence:**  
   - Error Return Sovereignty: Leaf functions return `error` directly.
   - Boolean Naming: All booleans prefixed with `is` or `has` (`isActive`, `isCompleted`, `isFailed`, `hasActiveLock`, `isSuccess`, `hasChanges`, `isSkipped`, `isEncrypted`, `isScanEnabled`).
   - Function limit: 8-15 lines max with flat guard clauses.
   - Whitespace discipline: Blank line after `}` and before `return`.
