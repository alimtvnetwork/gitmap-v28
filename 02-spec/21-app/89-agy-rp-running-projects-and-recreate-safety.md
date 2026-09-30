# SPEC-APP-89: AGY RP Running-Projects Remap & Recreate-Project Safety Guard

## User Request (Verbatim)

```text
PS C:\Users\Administrator> gitmap agy rp

  ┌── Antigravity Recreate Project ──────────────────────────────────────┐
  │ Project:    Administrator                                             │
  │ Path:       C:\Users\Administrator                                    │
  │ Action:     Purge cache/convs → Re-add to AGY → Read Memory conv     │
  └───────────────────────────────────────────────────────────────────
```

User Expectation & Constraint:
- Running `gitmap agy rp` was expected to list running projects (`running-projects ls`), but it unexpectedly triggered `recreate-project` (`rp` alias) on `C:\Users\Administrator`, creating an ad-hoc project, purging conversations, and injecting prompts into active session.
- `gitmap agy rp` must strictly map to `running-projects` (equivalent to `gitmap agy running-projects ls`).
- Remove `rp` alias from `recreate-project`. Only allow explicit aliases like `recreate`, `rcp`, `rec`.
- Strict safety guard on `recreate-project`: it must NEVER run automatically on root drives (`C:\`, `D:\`, `/`), user home directories (`C:\Users\<user>`, `~`, `/home/<user>`), or system directories without explicit confirmation. Non-project folders must not be converted into ad-hoc projects.

---

## 1. Root Cause Analysis

1. **Subcommand Normalization Hijack:**
   - In `cli/cmdagy/agy_cmd.go` line 163 (`normalizeProjectSubcommands`), `rp` was listed as an alias for `recreate-project`:
     `if low == "recreate-project" || low == "recreate" || low == "rp" || low == "rec" { return "recreate-project" }`
   - Even though `cli/cmdagy/agy_running_projects.go` line 42 defined `Aliases: []string{"runningprojects", "rp"}`, `normalizeProjectSubcommands` runs before Cobra subcommand dispatch, rewriting `args[0]` from `"rp"` to `"recreate-project"`.
   - `normalizeWorkflowSubcommands` only handled `"running-projects"` and `"runningprojects"`, omitting `"rp"`.

2. **Top-Level Root Dispatch Shadowing:**
   - In `cli/cmd/root.go` line 712 (`dispatchAgySubsystem`), `"rp"` was included in the case list for `cmdagy.DispatchAgy`.
   - However, in `cli/constants/constants_cli.go` and `cli/cmd/rootrelease.go`, `"rp"` is canonical alias for `release-pending` (`constants.CmdReleasePendingAlias = "rp"`). Top-level `gitmap rp` was intercepted by `dispatchAgySubsystem` instead of executing `release-pending`.

3. **Unprotected Ad-Hoc Recreate Fallback:**
   - In `cli/cmdagy/agy_recreate_resolve.go` line 23 (`resolveCurrentDirTarget`), if no arguments are passed, it checks `os.Getwd()`. When run from `C:\Users\Administrator`, it did not match any registered Antigravity projects, and instead of failing safely, line 34 invoked `buildAdHocProject(targetDir)`.
   - This treated the entire user profile directory as an ad-hoc project, wiped prompt queues, and triggered a conversation creation on `C:\Users\Administrator`.

---

## 2. Technical Requirements & Architectural Design

### 2.1 Subcommand Routing & Normalization
1. **`cli/cmdagy/agy_cmd.go`:**
   - In `normalizeProjectSubcommands(low string)`:
     - Replace `"rp"` with `"rcp"`:
       `if low == "recreate-project" || low == "recreate" || low == "rcp" || low == "rec"` -> `return "recreate-project"`
   - In `normalizeWorkflowSubcommands(low string)`:
     - Add `"rp"` to `running-projects`:
       `if low == "running-projects" || low == "runningprojects" || low == "rp"` -> `return "running-projects"`

2. **`cli/cmdagy/agy_recreate_cmd.go`:**
   - Update `agyRecreateProjectCmd.Aliases`:
     `Aliases: []string{"recreate", "rcp", "rec"}` (remove `"rp"`).

3. **`cli/cmd/root.go`:**
   - In `dispatchAgySubsystem` (line 712):
     - Replace `"rp"` with `"rcp"`. Top-level `gitmap rp` strictly routes to `release-pending` via `releaseDispatchEntries()`.
     - Subcommand `gitmap agy rp` routes to `cmdagy.DispatchAgy` -> `running-projects`.

### 2.2 Recreate Project Safety Guards (`cli/cmdagy/agy_recreate_resolve.go` & `agy_recreate_ops.go`)
1. **`isRestrictedSystemOrHomeDir(dirPath string) bool`:**
   - Resolves absolute, cleaned path.
   - Detects Volume / Drive Roots:
     - Windows drive roots: `^[a-zA-Z]:[\\/]?$` or `filepath.Dir(clean) == clean`.
     - Unix root: `/`.
   - Detects User Home & User Parent:
     - `os.UserHomeDir()` (e.g. `C:\Users\Administrator`, `/home/user`, `/root`).
     - `filepath.Dir(homeDir)` (e.g. `C:\Users`, `/home`, `/Users`).
   - Detects System Directories:
     - Windows: `%WINDIR%`, `%SystemRoot%`, `%ProgramFiles%`, `%ProgramFiles(x86)%`, `%ProgramData%`.
     - Unix: `/etc`, `/usr`, `/bin`, `/sbin`, `/var`, `/System`, `/Library`.
   - Returns `true` if `dirPath` matches any restricted location.

2. **`isGitRepo(dirPath string) bool`:**
   - Checks if `.git` exists in `dirPath` or if `gitutil.RepoRoot(dirPath)` returns without error.

3. **Enforcement in Resolution & Execution:**
   - In `resolveCurrentDirTarget`:
     - If `isRestrictedSystemOrHomeDir(cwd)`: returns error `E9101: restricted system or home directory cannot be recreated as an Antigravity project: <cwd>`.
     - If `!isGitRepo(targetDir)` and not an existing AGY project: returns error `E9102: directory %q is neither an existing Antigravity project nor a valid Git repository; refusing to create ad-hoc project`.
   - In `resolveDirectoryOrFolderTarget` and `resolveLikelyPathToken`:
     - If path is restricted or non-git: reject with `E9101`/`E9102`.
   - In `processSingleRecreate`:
     - Guard: `if isRestrictedSystemOrHomeDir(p.GetPath())`: return `E9101`.

---

## 3. Verification & Validation Criteria

1. **`gitmap agy rp` Execution:**
   - Running `gitmap agy rp` from any directory (including `C:\Users\Administrator`) must execute `running-projects` (listing active or enqueued prompts across projects), NEVER `recreate-project`.
2. **`gitmap agy recreate` Safety Guard:**
   - Running `gitmap agy recreate` from `C:\Users\Administrator` must immediately abort with `E9101` error, without modifying, purging, or registering any project.
   - Running `gitmap agy recreate C:\` or `gitmap agy recreate C:\Users` must abort with `E9101`.
   - Running `gitmap agy recreate <non-git-dir>` must abort with `E9102`.
3. **Top-Level `gitmap rp` Verification:**
   - `gitmap rp` must execute `release-pending`, not `recreate-project`.
4. **Unit Tests:**
   - `cli/cmdagy/agy_recreate_test.go` must pass with tests for restricted directory detection and `rcp` alias resolution.
5. **Coding Guidelines:**
   - All Go code must adhere to coding guidelines: positive booleans, structured error handling (`*apperror.AppError`), and function size bounds.
