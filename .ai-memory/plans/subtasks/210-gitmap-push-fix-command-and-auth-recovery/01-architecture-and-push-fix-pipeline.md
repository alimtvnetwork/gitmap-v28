# Subtask 210.1: Architecture & Push-Fix Pipeline Core

- **Parent Plan:** [80-gitmap-push-fix-command-and-auth-recovery.md](../../80-gitmap-push-fix-command-and-auth-recovery.md)
- **Spec Reference:** [02-spec/21-app/210-gitmap-push-fix-command-and-auth-recovery/01-architecture-spec.md](../../../../02-spec/21-app/210-gitmap-push-fix-command-and-auth-recovery/01-architecture-spec.md)
- **Status:** In-Progress
- **Target Area:** `cli/cmd`, `cli/cmdpull`, `cli/constants`

---

## Objective
Establish the core command routing, CLI flag handling, anti-hang subprocess execution wrapper, and Phase 1 pre-flight Git inspection pipeline for `gitmap push-fix` (and its subcommands/aliases: `gitmap push fix`, `gitmap pushfix`, and `gitmap pf`).

---

## Detailed Implementation Breakdown

### 1. CLI Constants & Flag Registration
- Define command constants in `cli/constants/constants_cli.go`:
  ```go
  CmdPushFix      = "push-fix"
  CmdPushFixAlias = "pushfix"
  CmdPushFixShort = "pf"
  ```
- Define flag descriptions in `cli/constants/constants_cli.go`:
  - `--remote` / `-r`: Remote repository name (default: `origin`).
  - `--branch` / `-b`: Target branch name (default: current active branch).
  - `--dry-run` / `-n`: Execute diagnostics and auth checks without performing push.
  - `--no-rebase`: Skip auto-rebase on non-fast-forward push rejection.

### 2. Root & Subcommand Dispatch Integration
- In `cli/cmd/rootcore.go`:
  - Register `CmdPushFix`, `CmdPushFixAlias`, `CmdPushFixShort` to route directly to `runPushFix(argsTail())`.
- In `cli/cmdpull/push.go`:
  - In `runPush(args []string)`, detect if `len(args) > 0 && args[0] == "fix"`, delegating execution to `RunPushFix(args[1:])`.

### 3. Pre-flight & Git State Inspection (`cli/cmdpull/push_fix.go`)
- Create `cli/cmdpull/push_fix.go` orchestrating the 4-phase recovery pipeline.
- Implement pre-flight inspection functions:
  1. `checkGitRepoCWD() (bool, error)`: Verify current directory is inside a Git repository.
  2. `detectCurrentBranch(cwd string) (string, error)`: Query active branch name using `git symbolic-ref --short HEAD` (fallback: `git rev-parse --abbrev-ref HEAD`).
  3. `parseRemoteURL(cwd, remoteName string) (string, error)`: Retrieve remote URL using `git remote get-url <remote>`.
  4. `countUnpushedCommits(cwd, branch string) (int, bool, error)`:
     - Check if branch has an upstream tracking reference: `git rev-parse --abbrev-ref @{u}`.
     - If upstream exists: execute `git rev-list --count @{u}..HEAD`.
     - If upstream does not exist: execute `git rev-list --count HEAD --not --remotes`.
     - Return commit count and `hasUpstream` boolean.

### 4. Anti-Hang Subprocess Execution Environment
- Implement `runIsolatedGitCmd(cwd string, args ...string) (string, string, int, error)`:
  - Isolate environment variables to guarantee headless, non-interactive execution:
    ```go
    cmd.Env = append(os.Environ(),
        "GIT_TERMINAL_PROMPT=0",
        "GCM_INTERACTIVE=never",
        "GIT_SSH_COMMAND=ssh -o BatchMode=yes -o ConnectTimeout=5",
    )
    ```
  - Bind execution context with hard timeout (`30 * time.Second`).
  - Capture standard output, standard error, and process exit code separately.

### 5. Early Exit for Clean & Up-to-Date State
- If `unpushedCommits == 0` and working directory is clean (`git status --porcelain` is empty):
  - Emit clear informative notice: `✓ Everything up-to-date (0 unpushed commits)`.
  - Exit gracefully with status `0`.
