# Issue 50 RCA: AGY RP Subcommand Alias Hijack, Unconfirmed Ad-Hoc Project Creation & CWD Pollution

## 1. Reproduction

1. **Triggering Subcommand from User Home Directory:**
   - Execute `gitmap agy rp` from the user profile root directory:
     ```powershell
     PS C:\Users\Administrator> gitmap agy rp
     ```
2. **Observed Erroneous Behavior:**
   - Instead of listing running projects (`running-projects ls`), the command executed `recreate-project` on `C:\Users\Administrator`.
   - The CLI rendered the `Antigravity Recreate Project` banner:
     - Name: `Administrator`
     - Path: `C:\Users\Administrator`
     - Action: Purged cache, registered `C:\Users\Administrator` as a new project in `~/.gemini/config/projects/b15c2df0...json`, and injected a prompt into the user's active session.
   - Subsequently, `gitmap agy running-projects` included `Administrator` with status `QUEUED` because candidate workspace discovery unconditionally scanned `cwd` (`os.Getwd()`).

## 2. Root Cause Analysis

1. **Subcommand Normalization Hijack:**
   - In `cli/cmdagy/agy_cmd.go` (`normalizeProjectSubcommands`), `rp` was listed as an alias for `recreate-project`:
     `if low == "recreate-project" || low == "recreate" || low == "rp" || low == "rec" { return "recreate-project" }`
   - Because `normalizeProjectSubcommands` executed before Cobra's subcommand lookup, any invocation of `gitmap agy rp` was rewritten to `recreate-project`.
   - Furthermore, `cli/cmd/root.go` (`dispatchAgySubsystem`) intercepted top-level `gitmap rp` instead of routing to `release-pending`.

2. **Unconditional Current Working Directory (`cwd`) Injection:**
   - In `cli/cmdagy/agy_queue_workspaces.go`, `collectCandidateWorkspaces()` unconditionally executed:
     ```go
     cwd, _ := os.Getwd()
     if cwd != "" {
         addCandidateWorkspace(wsMap, seen, cwd, filepath.Base(cwd))
     }
     ```
   - Running from `C:\Users\Administrator` injected the user home directory as a workspace candidate, causing `running-projects` to inspect and report `Administrator` as an active or queued workspace.

3. **Silent Ad-Hoc Project Registration Without User Confirmation:**
   - `resolveCurrentDirTarget` in `cli/cmdagy/agy_recreate_resolve.go` had a fallback: if `cwd` was not an existing Antigravity project, it invoked `buildAdHocProject(targetDir)`.
   - `workspacesync.SyncAntigravity` lacked restricted directory guards and wrote project registration JSONs into `~/.gemini/config/projects/` for any directory path, including user home and system paths.

## 3. Corrective Implementation

1. **Subcommand Routing & Alias Disambiguation:**
   - Removed `rp` from `normalizeProjectSubcommands` and added `rcp` (`recreate-project`, `recreate`, `rcp`, `rec`).
   - Mapped `rp` to `running-projects` in `normalizeWorkflowSubcommands`.
   - Removed `rp` from `agyRecreateProjectCmd.Aliases`.
   - Removed `rp` from `cli/cmd/root.go:dispatchAgySubsystem`, restoring top-level `gitmap rp` to `release-pending`.

2. **Elimination of CWD Injection in Running Projects Discovery:**
   - In `cli/cmdagy/agy_queue_workspaces.go`, removed `os.Getwd()` from `collectCandidateWorkspaces()`.
   - `collectCandidateWorkspaces()` strictly inspects only registered Antigravity projects that are valid Git repositories and not restricted paths.

3. **Mandatory Confirmation & Restricted Directory Guards:**
   - In `cli/cmdagy/agy_recreate_resolve.go`:
     - Added `IsRestrictedSystemOrHomeDir` to permanently block user home (`C:\Users\Administrator`), drive roots (`C:\`, `D:\`), and OS system directories (`E9101`).
     - Added `isGitRepo` validation to reject non-git folders (`E9102`).
     - Added `--confirm` (`-y`) requirement: if a folder is not an existing registered Antigravity project, it requires explicit `--confirm` to register and recreate (`E9103`).
   - In `cli/workspacesync/workspacesync.go`:
     - Added `isRestrictedPath` and `isGitRepoPath` guards in `SyncAntigravity` to refuse registering restricted paths or non-git folders into `~/.gemini/config/projects/`.

## 4. Prevention

1. **Unit Test Coverage:**
   - Added unit tests in `cli/cmdagy/agy_recreate_test.go` verifying restricted directory detection, non-git rejection, and `rcp` vs `rp` normalization.
2. **Deterministic Command Separation:**
   - Inspection commands (`running-projects` / `rp`) are strictly decoupled from lifecycle-modifying commands (`recreate-project` / `rcp`).
