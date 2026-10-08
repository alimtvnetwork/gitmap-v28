# Subtask 04: Zero-Loss File Removal, Temp Staging & SHA-256 Revert Engine

- **Parent Task:** `240-mcp-server-ai-analysis-agm-nodes-and-git-cache`
- **Subtask ID:** `04-task-file-removal`
- **Spec Reference:** `02-spec/21-app/240-mcp-server-ai-analysis-agm-nodes-and-git-cache/01-architecture-spec.md` (Sections 2.2, 3.1, 3.2, and 5)
- **Status:** `PENDING`
- **Assigned Subagent:** Implementation Subagent (Phase 2)

---

## 1. Objective & Scope

Implement the **Zero-Loss File Removal & Cryptographic Revert Engine** for the GitMap AI Analysis subsystem.

### Problem Statement
When autonomous AI agents (such as Google Antigravity, Cursor, or Claude) perform mass refactoring, code deduplication, or architecture simplification, they frequently remove obsolete files. If an AI agent hallucinates or makes a mistaken deduction, unlinking files directly with `os.Remove()` causes catastrophic permanent data loss.

### The Solution: Zero-Loss Vault Architecture
1. **Pre-Unlink Staging:** Before any file is unlinked, GitMap stages the file into an isolated temporary directory vault:
   ```text
   <temp>/gitmap/removed/<taskId>/<relPath>
   ```
2. **Cryptographic Integrity:** A pre-removal SHA-256 checksum (`BeforeSha256`) is computed and stored alongside the file metadata in `AiTaskFile`.
3. **Database Ledger:** The file is cataloged in `ai-analysis.db` with `Action = "removed"`, `BackupPath = "<staged_path>"`, and affirmative flag `IsRemoved = 1`.
4. **Instant Bit-for-Bit Reversion:** Executing `gitmap ai-analysis revert <taskId> [relPath]` verifies that the staged file's hash matches `BeforeSha256`, restores the file to the repository working tree, creates missing parent directories, and updates `AiTaskFile` (`IsRemoved = 0`, `Action = "restored"`).

---

## 2. Concrete Files to Create / Modify

| File | Nature | Purpose |
| :--- | :--- | :--- |
| `cli/tempdir/ai_removed_vault.go` | Create | Vault manager implementing staging path resolution, copy-with-mode, SHA-256 hashing, and verification. |
| `cli/cmdai/ai_analysis_rm.go` | Create | CLI command `gitmap ai-analysis rm <taskId> <paths...>` staging files and unlinking safely. |
| `cli/cmdai/ai_analysis_revert.go` | Create | CLI command `gitmap ai-analysis revert <taskId> [relPath]` verifying checksums and restoring files to the tree. |
| `cli/cmdai/ai_analysis_rm_test.go` | Create | End-to-end integration tests verifying file deletion, staging checksum match, and complete revert restoration. |

---

## 3. Detailed Component Architecture

### 3.1 Vault Engine (`cli/tempdir/ai_removed_vault.go`)
- `ResolveAiRemovedVaultDir(taskId string) string`: Resolves `<temp>/gitmap/removed/<taskId>/`.
- `StageFileForRemoval(taskId, repoRoot, relPath string) (*StagedFileResult, error)`:
  - Validates `relPath` is within `repoRoot`.
  - Reads source file bytes and computes SHA-256.
  - Writes bytes to `<temp>/gitmap/removed/<taskId>/<relPath>` preserving file permissions.
  - Re-reads staged file to guarantee byte-for-byte fidelity.
- `RestoreStagedFile(taskId, repoRoot, relPath, expectedSha string) error`:
  - Locates staged backup file.
  - Verifies SHA-256 matches `expectedSha`.
  - Reconstructs parent directory hierarchy in `repoRoot`.
  - Atomically writes file to `relPath`.

### 3.2 CLI Commands
- `gitmap ai-analysis rm <taskId> <path...> [--reason <text>] [--json]`
  - Emits JSON envelope with status `staged_and_removed`, list of affected files, and SHA-256 hashes.
- `gitmap ai-analysis revert <taskId> [relPath] [--force] [--json]`
  - Emits JSON envelope with status `restored`, list of restored files, and verification outcomes.

---

## 4. Acceptance Criteria

- [ ] File removal command `gitmap ai-analysis rm` NEVER deletes a file without prior staging in `<temp>/gitmap/removed/<taskId>/`.
- [ ] Pre-deletion SHA-256 hash is accurately calculated and committed to `AiTaskFile.BeforeSha256`.
- [ ] Database record committed with `IsRemoved = 1` and `Action = "removed"`.
- [ ] Revert command `gitmap ai-analysis revert` restores all removed files (or a targeted relative path) bit-for-bit.
- [ ] Revert aborts with an informative error if the staged backup file does not match `BeforeSha256`.
- [ ] Full machine-readable `--json` envelope supported for all commands.
- [ ] Unit and integration tests in `cli/cmdai/ai_analysis_rm_test.go` pass 100%.

---

## 5. Verification Commands

```powershell
# 1. Run integration tests for file staging and reversion
go test -v ./cli/tempdir -run "TestAiRemovedVault"
go test -v ./cli/cmdai -run "TestAiAnalysisRmAndRevert"

# 2. Verify staging directory isolation manually
# Verify that <temp>/gitmap/removed/<taskId>/ is cleaned up or created properly
```
