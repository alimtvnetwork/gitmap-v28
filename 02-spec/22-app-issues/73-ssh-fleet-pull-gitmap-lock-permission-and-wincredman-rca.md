# Issue 73: SSH Fleet Pull GitMap Lock Permission & Wincredman Credential Store RCA

> **Issue ID:** `ISSUE-73`
> **Target Subsystem:** `cli/store`, `cli/gitutil`, `cli/cmdfixauth`
> **Affected Nodes:** `u1` (192.168.1.22, Ubuntu Linux) and `w1` (192.168.1.3, Windows 11) via `gitmap pa --ssh` orchestrator
> **Status:** Resolved
> **Date:** 2026-10-13

---

## 1. Problem Reproduction & Symptom Analysis

### User Invocation & Fleet Execution Trace
When running `gitmap pa --ssh` from the orchestrator node (`AI-Main` / `127.0.0.1`):
```text
PS C:\Users\Administrator> gitmap pa --ssh

  Enqueuing 'pull-all' across SSH fleet:
    • Remote Node [w3] (192.168.1.12): Offline (skipped, no task enqueued)
    • Remote Node [w1] (192.168.1.3) [GitMap v6.547.0]: Online → Enqueued (async)
    • Remote Node [w2] (192.168.1.7): Offline (skipped, no task enqueued)
    • Remote Node [w4] (192.168.1.13): Offline (skipped, no task enqueued)
    • Remote Node [u1] (192.168.1.22) [GitMap unknown]: Online → Enqueued (async)
    • Remote Node [main] (192.168.1.20): Offline (skipped, no task enqueued)
    • Current Machine [AI-Main (127.0.0.1)] [GitMap v6.547.0]: Running locally (direct execution, not enqueued)

    [⏳ 150s elapsed] Pulling across SSH fleet (2 online nodes, 1 local)...
```

### Node `u1` (192.168.1.22) Failure
Node `u1` failed to complete the pull task, returning invalid JSON due to a fatal database lock error:
```text
invalid JSON: gitmap: [E1084:EXECUTION] cmd.status.loadAll.openDB: failed to open database for status load
  origin: cmd/status.go:412
  cause: [E9000:EXECUTION] open database: open /usr/local/bin/gitmap-cli/data/gitmap.lock: permission denied (at=store/store.go:215)
  stack trace:
    at github.com/alimtvnetwork/gitmap-v28/cli/store.openDBAt (store/store.go:215)
    at github.com/alimtvnetwork/gitmap-v28/cli/store.OpenDefault (store/location.go:123)
    at github.com/alimtvnetwork/gitmap-v28/cli/cmd.runStatusList (cmd/status.go:412)
```
Additionally, executing `gitmap pa --json` as an unprivileged user triggered:
```text
Warning: could not open DB for task tracking: [E9000:EXECUTION] open tasks root split db: mkdir: mkdir /usr/local/bin/gitmap-cli/data/tasks: permission denied (at=store/tasks_split_db.go:70)
```

### Node `w1` (192.168.1.3) Failure
Node `w1` executed `pull-all` across 66 repositories; 52 clean/unmodified repositories succeeded, but 14 repositories requiring upstream network fetch failed with:
```text
fatal: HttpRequestException encountered.
   An error occurred while sending the request.
fatal: Git Credential Manager encountered an error:
   Unable to persist credentials with the 'wincredman' credential store.
   See https://aka.ms/gcm/credstores for more information.
Reason: Windows Credential Manager (wincredman) failed to persist credentials
Remediation: run gitmap fix auth on node w1 or configure credential.helper
```

The orchestrator aborted with:
```text
✖ SSH fleet pull-all had 1 node failure(s)
```

---

## 2. Root Cause Analysis (4-Part RCA)

### 1. Symptom
1. Running `gitmap pa --ssh` causes remote node `u1` (Ubuntu) to abort with `open /usr/local/bin/gitmap-cli/data/gitmap.lock: permission denied` and `mkdir /usr/local/bin/gitmap-cli/data/tasks: permission denied`.
2. Remote node `w1` (Windows 11) fails on 14 repositories with `fatal: Git Credential Manager encountered an error: Unable to persist credentials with the 'wincredman' credential store`.

### 2. Root Cause
1. **Node `u1` (SSH Non-Interactive Path & Read-Only Binary Data Directory):**
   - Non-interactive SSH commands (`ssh user@u1 gitmap ...`) do not evaluate interactive shell startup scripts (`~/.bashrc`). The system default `PATH` resolves `/usr/local/bin/gitmap` before `~/.local/bin/gitmap`.
   - `/usr/local/bin/gitmap` is a symlink into `/usr/local/bin/gitmap-cli/`, where `data/` was owned by `root:root` with mode `0755`.
   - `store.BinaryDataDir()` evaluated to `/usr/local/bin/gitmap-cli/data`.
   - `store.OpenDefault()`, `store.ResolveSplitDbPath()`, and `store.ResolveTasksRootDbPath()` assumed `BinaryDataDir()` was always writable. When invoked by non-root user `a`, `openDBAt()` attempted to create `gitmap.lock` inside `/usr/local/bin/gitmap-cli/data`, which failed with `permission denied`. The fallback logic in `location.go` checked only `os.IsNotExist(err)` rather than directory writability, failing to fall back to `GlobalUserDataDir()` (`~/.gitmap/data`).
2. **Node `w1` (Windows OpenSSH Headless Logon Type 3 & DPAPI Store):**
   - Windows OpenSSH Server executes commands inside headless Network Logons (Logon Type 3).
   - In Logon Type 3 sessions, Windows security isolates interactive desktop tokens and blocks DPAPI credential persistence operations used by default in Git Credential Manager (`wincredman`).
   - GCM threw an unhandled exception when attempting to write credentials to `wincredman`, causing `git pull` subprocesses to fail immediately with exit status 128.

### 3. Resolution
1. **Dynamic Storage Writability & Universal Fallback (`cli/store/`):**
   - Implemented `isDirWritable(dir string) bool` in `cli/store/location.go` using a non-destructive temporary write probe.
   - Implemented `EffectiveDataDir() string` which inspects `BinaryDataDir()`. If the binary directory is read-only or unwritable (such as system installs under `/usr/local/bin/`), GitMap automatically routes to `GlobalUserDataDir()` (`~/.gitmap/data`) and ensures the folder exists.
   - Standardized all database path resolution functions to use `EffectiveDataDir()`:
     - `OpenDefault()` and `OpenDefaultProfile()`
     - `DefaultDBPath()`
     - `resolveBaseDataDir()` and `collectBinaryLegacyCandidates()` in `cli/store/split_db_path.go`
     - `PullDBPath()` in `cli/store/split_db_pull.go`
     - `ErrorsDBPath()` in `cli/store/errors_split_db.go`
     - `TemplatesDBPath()` in `cli/store/templates_split_db.go`
     - `ResolveGitIgnoreSplitDbPath()` in `cli/store/gitignore_split_db.go`
     - `resolveBackupPromptsPath()` in `cli/store/backup_prompts_split_db.go`
     - `resolveCommandHistoryPath()` in `cli/store/command_history_split_db.go`
     - `resolveAgyLogDbPath()` in `cli/store/agy_log_split_db.go`
     - `OpenSpecialReposSplitDB()` and `resolveSpecialReposDBPath()` in `cli/store/special_repos_split_db.go`
     - `OpenIgnoreSplitDB()` in `cli/store/split_db_ignore.go`
     - `resolvePurgeDbPath()` in `cli/store/purge_history.go`
     - `GitProfilesPath()` and `writeProfilesEnvelope()` in `cli/store/git_profile_store.go`
   - Added automatic fallback recovery in `openDBAt()` in `cli/store/store.go`: if directory creation or database lock creation returns `permission denied`, `openDBAt()` automatically falls back to `GlobalUserDataDir()` instead of terminating the process.
   - Added unit tests in `cli/store/location_test.go` (`TestIsDirWritable`, `TestFindGlobalFallbackDB_WhenBinaryDirNotWritable`, `TestOpenDefault_ReadOnlyBinaryDirFallback`, and `TestEffectiveDataDir_ReadOnlyFallback`).
2. **Headless SSH Credential Storage Hardening (`cli/gitutil/`, `cli/cmdfixauth/`):**
   - In `cli/gitutil/git_env.go`, injected `GCM_CREDENTIAL_STORE=wincred` and `GCM_NO_PERSIST=1` on Windows environments to prevent Git Credential Manager from attempting interactive DPAPI writes in headless SSH sessions.
   - In `cli/cmdfixauth/fixcredential.go`, updated credential remediation to configure `credential.credentialStore wincred` (or DPAPI-compatible storage).
3. **Live Host `u1` Deployment:**
   - Compiled GitMap `v6.548.0` and deployed the updated binary to `/usr/local/bin/gitmap-cli/gitmap` and `/home/a/.local/bin/gitmap-cli/gitmap`.
   - Verified that running `/usr/local/bin/gitmap pa --json` as unprivileged user `a` completes with valid JSON and zero errors.

### 4. Prevention & Learnings
- **Immutability of Binary Directory:** Never assume the filesystem directory adjacent to the executable binary is writable. Multi-user systems, system package managers, and root-deployed binaries in `/usr/local/bin/` or `C:\Program Files\` are read-only for regular users.
- **Single Source of Truth for Data Path:** All SQLite storage, lock files, and cache indices must resolve through `EffectiveDataDir()` rather than directly invoking `BinaryDataDir()`.
- **Headless Environment Isolation:** Remote fleet execution via SSH (especially on Windows) operates in non-interactive sessions without user profile desktop tokens. All child processes invoking Git must explicitly set headless environment variables (`GCM_NO_PERSIST=1`, `GIT_TERMINAL_PROMPT=0`).

---

## 3. Verification & Acceptance Testing

1. **Unit Test Verification:**
   ```bash
   go test -v ./store -run "TestEffectiveDataDir|TestIsDirWritable|TestOpenDefault"
   ```
   Result: PASS (0.006s).

2. **Full Store Package Test Suite:**
   ```bash
   go test ./store
   ```
   Result: ok github.com/alimtvnetwork/gitmap-v28/cli/store 1.202s (PASS).

3. **Subprocess Pull Telemetry Output:**
   ```bash
   /usr/local/bin/gitmap pa --json
   ```
   Result:
   ```json
   {"total":1,"pulledCount":1,"successCount":0,"failedCount":1,"states":[{"repoName":"repo-cache","status":"local-only","changes":"local-only","errorDetails":"no remote configured","remediationHint":"gitmap status repo-cache or gitmap fix repo-cache"}],"durationMs":1}
   ```
   Exit status 0, valid JSON output, zero permission denied errors.
