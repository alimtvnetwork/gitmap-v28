# Plan 77: AGY Scripts-Fixer Parity Verification, SSH-First Access Probe & Interactive Terminal Auth Prompt

## Task Execution Header
- **Initial Trigger**: User request to verify AGY clear implementation parity against `D:\work\scripts-fixer`, verify AGY commands/prompts/end-to-end tests, and implement an SSH-first access check for parallel pulls/clones with terminal authentication fallback (PAT entry or browser login, repository attribution, and global token reuse).
- **Execution Lifecycle**: Completed in 1 continuous multi-phase loop (Phases 1A, 1B, 2, 3) across 10 modified/created files.
- **Verification Gates**: Targeted quality linting (`go vet ./...` exit 0), strict file size limits (<= 100 lines per file), function limits (<= 15 lines), zero nested ifs.

---

## User Request (Verbatim)
```text
https://prnt.sc/2NP6Jvwb5d-v

Can you please read the scripts fixture and check the AGY new implementation, and have we implemented this in our code base for AGY, clear commands from scripts fixture? Please confirm. And also, make sure that the AGY commands, prompting, and others are good with the end-to-end testing and working fine. Can you please confirm that? In some cases, this is another issue that we have. We see that the GitHub token that needs to be pushed outside, right? Because it's trying to pull parallely. In those cases, it does not work or did not work or require an access token. The first thing the Git map will try is to see that it has the access using SSH. If it has the access using SSH, it will forget about the access token. But in case it cannot access it using the SSH, then in the terminal, it will seek for the access token or the browser login. So it will have the two options for the user to fix it, and user will prompt and select which option they are doing. And also on top of it, it will mention which repository's token, because without this information, how user can provide the access token. And also there will be an option to reuse this access token everywhere. So these are the things I want. Do you understand the task first, confirm, and can you please complete this task with confidence? Okay. Do you think is there anything that is out of your hand that you cannot do, then let me know. If not, then complete it without a failure. And then at the end, minor bump and release it properly.
```

---

## Consolidated Subtasks & Delivered Features

### 1. Subtask 01: Audit AGY Clear Implementation & Verify E2E Parity (Task-01)
- Verified `D:\work\scripts-fixer\scripts\69-install-antigravity\helpers\agy_optimizer.py` and confirmed full native parity in GitMap's `cli/cmdagy/`:
  - `gitmap agy clean-cache` / `cache-clear` with configurable retention (`--keep 10` by default).
  - Direct convenience shortcuts: `gitmap ccko` (`cache-clear-keep-one`), `gitmap cckf` (`cache-clear-keep-five`).
  - Preflight simulation flags: `--pre`, `--precheck`, `--preflight`.
  - Transaction backup and rollback via `gitmap agy undo`.
  - Cross-platform cache discovery for Windows (`%APPDATA%\Antigravity`, `%LOCALAPPDATA%\Antigravity\Cache`, `antigravity-updater`), macOS, and Linux.
  - End-to-end unit tests in `cli/cmdagy/*_test.go` verified passing.

### 2. Subtask 02: SSH-First Pre-Flight Access Probe (Task-02)
- Implemented `ProbeSSHAccess(repoURL string)` in `cli/cmdclone/clone_auth_probe.go`.
- Converts HTTPS URLs to SSH shorthand (`git@<host>:<owner>/<repo>.git`).
- Executes fast non-interactive SSH probe via `git ls-remote --exit-code -h <sshURL> HEAD` with `BatchMode=yes` and 5s timeout.
- If SSH works, switches clone/pull remote to SSH URL automatically and completely bypasses access tokens.
- Results cached per SSH URL using `sync.Map` to eliminate redundant probes.

### 3. Subtask 03: Interactive Terminal Auth Prompt with Global Token Reuse (Task-03)
- Prevented background GUI modal popups (like Git Credential Manager) by injecting `GCM_INTERACTIVE=never` and `GIT_TERMINAL_PROMPT=0` into git commands (`cli/constants/constants_git.go`, `cli/cloner/cloner.go`).
- Implemented `PromptTerminalAuth(repoName, repoURL string)` in `cli/cmdclone/clone_auth_prompt.go` displaying the exact repository name and 2 resolution choices:
  - Option 1: Personal Access Token (PAT) input.
  - Option 2: Browser login via `gh auth login` / device flow.
- Implemented `askTokenReuseConfirmation()` and thread-safe session store in `cli/cmdclone/clone_auth_store.go`.
- If user confirms token reuse, the token is saved into `globalToken` and automatically injected into subsequent repositories (`https://<token>@github.com/...`) without re-prompting.

### 4. Subtask 04: Integration & Minor Release (Task-04)
- Wired `ResolveRepoAuth` into `runCloneExecution`, `executeDirectClone`, and `executeDirectCloneOne` in `cli/cmdclone/clone.go` and `cli/cmdclone/clonemulti.go`.
- Verified 0 linter violations via `go vet ./...`.
- Executed minor version release bump from `v6.305.0` to `v6.306.0`.
