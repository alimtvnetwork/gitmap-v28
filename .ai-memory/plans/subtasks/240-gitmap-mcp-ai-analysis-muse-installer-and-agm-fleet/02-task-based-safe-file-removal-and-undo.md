# Subtask 02: Task-Based Safe File Removal & Undo

> **Subtask ID:** Subtask-02  
> **Parent Plan:** `.ai-memory/plans/240-gitmap-mcp-ai-analysis-muse-installer-and-agm-fleet.md`  
> **Target Subsystems:** `cli/cmdrm/`, `cli/cmd/rm.go`, `cli/cmd/roottooling.go`  
> **Owned Files:**  
> - `cli/cmdrm/rm_types.go`  
> - `cli/cmdrm/rm_manifest.go`  
> - `cli/cmdrm/rm_stage.go`  
> - `cli/cmdrm/rm_undo.go`  
> - `cli/cmdrm/rm_purge.go`  
> - `cli/cmdrm/rm_cmd.go`  
> - `cli/cmdrm/rm_test.go`  
> - `cli/cmd/rm.go`  

---

## 1. Concrete Objectives

1. **Eliminate Destructive File Deletion Anti-Patterns:**
   - Address the anti-pattern identified in `assets/screenshots/240-ai-rm-cleanup-patch.png` where AI agents execute raw shell commands (`Remove-Item patch_*.py -Force` or `rm -rf`).
   - Implement `gitmap rm <patterns...>` providing safe, recoverable file deletion for AI agents and human developers.

2. **Task Association & OS Temp Backup Staging:**
   - Every removal operation is scoped to a task ID.
   - If `--task <id>` is provided, use that ID; otherwise generate a deterministic identifier: `rm-task-<timestamp>-<uuid>`.
   - Before removing any file, copy it to an OS temporary backup staging area: `filepath.Join(os.TempDir(), "gitmap-rm", taskID, relativePath)`.
   - Preserve directory structure within the backup staging directory.

3. **Manifest Generation & Checksum Auditing (`manifest.json`):**
   - In `cli/cmdrm/rm_manifest.go`, author a structured manifest recording:
     - `taskId`: Unique removal session identifier.
     - `createdAt`: ISO-8601 UTC timestamp.
     - `repoRoot`: Normalized workspace root.
     - `fileCount`: Number of files staged.
     - `totalSizeBytes`: Aggregate size of removed files.
     - `files`: Array of file records containing `relativePath`, `stagedPath`, `sizeBytes`, `sha256`, and `fileMode`.
   - Persist manifest as `manifest.json` alongside staged backup files.

4. **Atomic Restoration Engine (`gitmap rm undo`):**
   - In `cli/cmdrm/rm_undo.go`, implement rollback:
     - If `task_id` is supplied, restore that specific task.
     - If omitted, identify and restore the most recent removal session in `$TEMP/gitmap-rm/`.
     - Re-create destination directories if deleted.
     - Copy files back to their exact original workspace locations, restoring original POSIX/Windows file modes.
     - Validate post-restore SHA-256 hashes against manifest entries.
     - Safely remove the staged temporary folder upon confirmed restore.
     - If the temporary directory was purged by OS cleanup, output a helpful error:
       `"Backup directory for task '<id>' not found in OS temporary storage. The files were permanently purged by OS temp cleaner and cannot be reverted."`

5. **Removal Inspection & Maintenance:**
   - Implement `gitmap rm list` displaying all active recoverable tasks in OS temp.
   - Implement `gitmap rm purge [--all | --older-than <duration>]` to permanently purge backup directories on demand.

6. **CLI Command Binding:**
   - In `cli/cmd/rm.go`, define top-level `RmCmd` (`gitmap rm`) and wire subcommands (`undo`, `list`, `purge`).
   - Bind flags `--task`, `-t`, `--undo`, `-u`, `--force`, `-f`, and `--dry-run`.
   - Register `RmCmd` into GitMap command dispatcher.

---

## 2. Core Domain Types & Structs (`cli/cmdrm/rm_types.go`)

```go
package cmdrm

import (
	"os"
	"time"
)

// RmFileRecord represents an individual file staged for removal.
type RmFileRecord struct {
	RelativePath string      `json:"relativePath"`
	StagedPath   string      `json:"stagedPath"`
	SizeBytes    int64       `json:"sizeBytes"`
	Sha256       string      `json:"sha256"`
	FileMode     os.FileMode `json:"fileMode"`
	DeletedAt    time.Time   `json:"deletedAt"`
}

// RmManifest represents the complete metadata for a task-based removal session.
type RmManifest struct {
	TaskId         string         `json:"taskId"`
	CreatedAt      time.Time      `json:"createdAt"`
	RepoRoot       string         `json:"repoRoot"`
	FileCount      int            `json:"fileCount"`
	TotalSizeBytes int64          `json:"totalSizeBytes"`
	Files          []RmFileRecord `json:"files"`
}

// RmOptions configures the safe removal execution.
type RmOptions struct {
	Patterns []string
	TaskId   string
	DryRun   bool
	Force    bool
	RepoRoot string
}

// RmUndoOptions configures restoration of a previous removal session.
type RmUndoOptions struct {
	TaskId   string
	RepoRoot string
}

// RmPurgeOptions configures manual purging of OS temporary backups.
type RmPurgeOptions struct {
	All       bool
	OlderThan time.Duration
	TaskId    string
}
```

---

## 3. Step-by-Step Implementation Plan

### Step 1: Staging Directory & Path Resolution

- In `cli/cmdrm/rm_stage.go`, implement `ResolveTaskBackupDir(taskID string) string`:
  - Returns `filepath.Join(os.TempDir(), "gitmap-rm", taskID)`.
  - Ensures path normalization using `filepath.Clean`.

### Step 2: File Matcher & Pattern Expander

- In `cli/cmdrm/rm_stage.go`, implement `ExpandFilePatterns(patterns []string, repoRoot string) ([]string, error)`:
  - Supports file paths, relative glob patterns (`patch_*.py`, `tmp/**/*.log`), and directory paths.
  - Excludes `.git/`, `.gitmap/`, and node modules from accidental deletion.
  - Asserts matching files exist on disk before proceeding.

### Step 3: Staging Engine & SHA-256 Calculation

- In `cli/cmdrm/rm_stage.go`, implement `StageAndRemoveFiles(opts RmOptions) (*RmManifest, error)`:
  - Generates `taskID` if empty.
  - Creates backup staging directory.
  - For each file:
    - Reads file attributes and calculates SHA-256 hash.
    - Copies file to staged path preserving subdirectory hierarchy.
    - Deletes file from workspace filesystem.
    - Records entry in `RmManifest`.

### Step 4: Manifest Serialization

- In `cli/cmdrm/rm_manifest.go`, implement:
  - `WriteManifest(dir string, manifest *RmManifest) error`
  - `ReadManifest(dir string) (*RmManifest, error)`
  - Uses atomic file write (`manifest.json.tmp` -> `manifest.json`).

### Step 5: Undo & Restoration Engine

- In `cli/cmdrm/rm_undo.go`, implement `ExecuteUndo(opts RmUndoOptions) error`:
  - If `opts.TaskId` is blank, scans `$TEMP/gitmap-rm/` for the most recent directory containing `manifest.json`.
  - Reads manifest.
  - Verifies each staged file exists.
  - Restores each file to original `relativePath`, creating parent dirs if needed.
  - Validates restored file size and SHA-256 against manifest.
  - Deletes staged directory upon confirmed restoration.

### Step 6: Backup Listing & Purging

- In `cli/cmdrm/rm_purge.go`, implement:
  - `ListRmBackups() ([]RmManifest, error)`: Scans `$TEMP/gitmap-rm/` and returns active sessions.
  - `PurgeRmBackups(opts RmPurgeOptions) (int, error)`: Removes expired or all staging folders.

### Step 7: Cobra CLI Command Integration

- Create `cli/cmdrm/rm_cmd.go` with `RmCmd`:
  - Subcommands: `undo`, `list`, `purge`.
  - Wire flags `--task`, `--dry-run`, `--force`, `--undo`.
  - Bind to `cli/cmd/rm.go` and export to root CLI tooling.

### Step 8: Unit Testing Suite

- Create `cli/cmdrm/rm_test.go`:
  - Test pattern expansion with temporary test files.
  - Test safe staging and manifest generation.
  - Test atomic undo and file byte parity.
  - Test missing temp folder error handling.
  - Test purge with duration threshold.

---

## 4. Acceptance Criteria

- [x] `gitmap rm <pattern>` safely removes files from workspace without permanently destroying them.
- [x] Removed files are preserved in `$TEMP/gitmap-rm/<task_id>/` with accurate directory structure.
- [x] `manifest.json` correctly stores file sizes, modes, and SHA-256 hashes.
- [x] `gitmap rm undo <task_id>` restores files with original permissions and checksums.
- [x] Running bare `gitmap rm undo` automatically restores the most recently deleted session.
- [x] If temporary staging directory is missing, a clear, descriptive error is returned.
- [x] `gitmap rm list` renders active backup sessions with task ID, date, and file counts.
- [x] `gitmap rm purge --all` safely frees temporary disk space.
- [x] Zero absolute paths or `file:///` URIs exist in any owned file.
- [x] All unit tests in `cli/cmdrm/rm_test.go` pass with 100% success rate.

---

## 5. Verification Commands

```bash
  # 1. Run unit tests for safe removal and undo
  go test -v ./cli/cmdrm/...

  # 2. Create temporary test files
  echo "test content 1" > temp_test_file1.tmp
  echo "test content 2" > temp_test_file2.tmp

  # 3. Test safe removal with task association
  gitmap rm "temp_test_*.tmp" --task "task-240-verify"

  # 4. Verify files are no longer in workspace
  gitmap find "temp_test_*.tmp"

  # 5. List recoverable tasks in temp storage
  gitmap rm list

  # 6. Test undo restoration
  gitmap rm undo task-240-verify

  # 7. Verify restored files and clean up
  gitmap rm "temp_test_*.tmp" --task "task-240-cleanup"
  gitmap rm purge --task "task-240-cleanup"

  # 8. Run relative paths linter
  python linter-scripts/check-relative-paths.py
```
