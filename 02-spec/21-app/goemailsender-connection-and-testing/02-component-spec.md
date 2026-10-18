# Component Specification: GoEmailSender, Multi-Profile Fallback & Fleet Delegation

## 1. Binary Resolution Component (`cli/cmdspace/send_email.go`, `cli/cmdspace/send_email_windows.go`)

### 1.1 Root Cause & Architectural Fix for Windows Resolution
* **Defect**: On Windows, `os.Stat` does not populate POSIX execute bits (`info.Mode().Perm()&0111`). `isExecutableFile` in `cli/cmdspace/send_email.go` returned `false` unconditionally. Furthermore, Go installs compiled binaries on Windows with a `.exe` extension, whereas `resolveGoEmailReaderBin` only probed bare `"go-email-reader"`.
* **Component Design**:
  - Extract platform-specific executable validation into `isExecutableFile(path string) bool`.
  - On Windows: verify regular file and check that file extension is `.exe`, `.cmd`, `.bat`, or empty with an existing counterpart `.exe`.
  - Expand candidates: `[]string{"go-email-reader", "GoEmailSender", "goemailsender"}`.
  - Probe locations in order:
    1. `exec.LookPath(candidate)` and `exec.LookPath(candidate + ".exe")`.
    2. `filepath.Join(os.Getenv("GOPATH"), "bin", candidate)`.
    3. `filepath.Join(userHome, "go", "bin", candidate)`.
    4. GitMap configuration setting: `email.go_email_reader_bin`.

## 2. Multi-Login Profile Fallback Component (`cli/cmdclone/clone_auth_resolve.go`, `cli/cmdclone/gh_resolver.go`, `cli/cmdclone/clone.go`)

### 2.1 Fallback Pipeline & Data Flow
* **Flags**: Add `--profile` / `-p` to `CloneFlags` and `cloneFixRepoFlags`.
* **Step 1: Primary Profile Probe**:
  - If `--profile <alias>` is specified, bind and use that profile exclusively.
  - Otherwise, resolve default profile from `store.LoadGitProfiles()` (`cfg.Default` or `isDefault == true`).
  - Attempt clone with primary profile credentials.
* **Step 2: Concurrent Secondary Probing on 401/404**:
  - If primary attempt encounters repository missing or unauthorized (`Repository not found`, `401 Unauthorized`, `404 Not Found`, `could not read Username`, `Permission denied (publickey)`):
  - Enumerate secondary profiles from `cfg.Profiles` where profile name != default.
  - Launch parallel lightweight verification using `git ls-remote --heads <url>` with 3-second context timeout.
  - The first profile returning exit code 0 is selected via atomic cancellation.
* **Step 3: Execution & Binding**:
  - Clone proceeds using the winning secondary profile.
  - Automatically invoke `usercontext.BindProject(absPath, winningProfile.Name)`.
  - If all secondary profiles fail, return aggregated diagnostic error.

## 3. Terminal Email UI & Scheduling Component (`cli/cmdspace/send_email_ui.go`, `cli/cmdspace/schedule_send_email.go`)

### 3.1 Terminal Drafting & Account Viewer UI
* **Command**: `gitmap space email ui` (or `gitmap space send-email --ui`).
* **Capabilities**:
  - Renders account status from `go-email-reader list` or `repo-secrets/11-go-email-reader/01-email-setup.json`.
  - Displays registered notifier recipient from `02-default-notifier.json`.
  - Interactive prompts for recipient (`--to`), subject (`--subject`), and body (`--body`) with preview.
  - Interactive schedule selector: prompts for trigger task (`release`, `pull`, `clone`, `scan`, `push`, `backup`).

## 4. Database Relative Path & Repo-Secrets URL Component (`cli/heatmap/heatmap_view.go`, `cli/cmdagent/agent_task_sqlite.go`)

### 4.1 Relative Path Hygiene & URL Resolution
* **Relative Path Calculation**:
  - Internal DB: `filepath.Rel(repoRoot, dbPath)` -> `.ai-memory/heatmap.db` or `.ai-memory/temp-agents/<slug>/agent-task.db`.
  - External/RC DB: Traversal from repo root -> `../repo-cache/.heatmap/...` or `../repo-secrets/...` using clean forward slashes.
* **Repo-Secrets URL Extraction**:
  - Helper `ResolveRepoSecretsURL()` inspects special repos split DB or `repo-secrets` git configuration (`remote.origin.url`).
  - Included as `"repoSecretsUrl"` in all JSON envelopes for `gitmap heatmap --json` and `gitmap agent task status --json`.

## 5. Fleet Nodes SSH Remote Dispatch & Delegation Component (`cli/cmdnodes/nodes_remote.go`, `cli/cluster/exec_delegate.go`, `cli/cmdssh/ssh_exec_command.go`)

### 5.1 Command Recognition & Remote Execution
* **Fix in `cli/cmdssh/ssh_exec_command.go`**: Add `"space"`, `"heatmap"`, `"hm"`, `"agent"`, `"agents"`, `"ai-agents"` to `isGitmapCoreCommand()`.
* **Nodes CLI**: `gitmap nodes run [target] "<command>" [flags]` (aliases `nodes exec`, `nodes r`).
* **Peer Delegation**: `cli/cluster/exec_delegate.go` defines `PeerDelegationPayload` and `PeerDelegationResponse` allowing one GitMap node to dispatch and collect telemetry from peer instances.
