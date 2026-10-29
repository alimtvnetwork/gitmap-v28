# Subtask 03: Machine Migration, Help Text Synchronization, and Release Ceremony

> **Sequence ID:** 347  
> **Parent Plan:** [.ai-memory/plans/pending/347-gitmap-user-change-command.md](../../pending/347-gitmap-user-change-command.md)  
> **Specification:** [02-spec/21-app/347-gitmap-user-change-command/03-migration-and-verification-spec.md](../../../../02-spec/21-app/347-gitmap-user-change-command/03-migration-and-verification-spec.md)  
> **Status:** Pending  
> **Category:** Feature, Migration & Release Management  

---

## 1. Scope & Objectives

This subtask governs the final phase of Plan 347 (`gitmap-user-change-command`), executing:
1. **Machine Git Identity Migration:** Switch the host machine's global Git user identity to `"AI automation bot"` with email `"local@local"` (or unset if empty email requested).
2. **Help Documentation Synchronization:** Integrate new CLI command aliases (`change-user`, `change-login`, `autobot`) into `cli/cmduser/user_help_menu.go` and the rich topic dispatcher in `cli/cmd/rich_help_dispatcher.go`.
3. **Automated Version Bump:** Advance repository version from `v6.627.0` to `v6.628.0` using `03-ai-scripts/37-bump-version.py`.
4. **Atomic Release Ceremony:** Finalize and seal the release using `gitmap cpr "v6.628.0 gitmap user change commands and autobot identity"`.

---

## 2. Step-by-Step Execution Plan

### Step 1: Rich Help Menu & Help Topic Dispatcher Integration
- Update `cli/cmduser/user_help_menu.go`:
  - Add usage lines for `gitmap change-user <name> [email]`, `gitmap change-login <name> [email]`, and `gitmap autobot [email]`.
  - Add entries in `buildGitUserSection()` detailing one-liner author changes and the autobot preset.
  - Enforce <= 15 lines per builder function.
- Update `cli/cmd/rich_help_dispatcher.go`:
  - Extend `tryRenderCoreRichTopic(topic string)` to handle `"change-user"`, `"change-login"`, and `"autobot"`.
  - Route directly to `cmduser.RenderUserHelp()`.
  - Ensure zero nested ifs and early returns.

### Step 2: Global Machine Git Identity Migration
- Execute global Git configuration commands on the local machine:
  ```bash
  git config --global user.name "AI automation bot"
  git config --global user.email "local@local"
  ```
- If an empty email is explicitly requested for unsetting:
  ```bash
  git config --global user.name "AI automation bot"
  git config --global --unset user.email
  ```
- Verify migration state:
  ```bash
  git config --global --get user.name
  git config --global --get user.email
  gitmap user info
  ```
- Confirm output reflects the active global identity card: `"AI automation bot <local@local>"`.

### Step 3: FastGate & Lint Verification
- Run FastGate validation script:
  ```bash
  gitmap py 03-ai-scripts/50-fastgate.py
  ```
- Audit codebase for strict adherence to coding guidelines:
  - Positive boolean naming (`is*`, `has*` only; zero negative booleans).
  - Leaf error sovereignty (return errors directly, no exit handlers in leaf functions).
  - Maximum 15 lines per function.
  - Relative Git paths only across all references.

### Step 4: Version Bump Automation
- Execute the SemVer minor bump script:
  ```bash
  gitmap py 03-ai-scripts/37-bump-version.py -t minor -s "gitmap user change commands and autobot identity"
  ```
- Confirm automated updates across all manifests:
  - `version.json` updated to `6.628.0`.
  - `package.json` updated to `6.628.0`.
  - `cli/version/constants.go` synchronized.
  - `latest.json` updated with tag `v6.628.0`.
  - `changelog.md` prepended with `## [6.628.0] - 2026-10-29`.
- Verify multi-manifest synchronization:
  ```bash
  gitmap py 03-ai-scripts/14-version-sync-checker.py
  ```

### Step 5: Atomic Release Ceremony
- Inspect working tree hygiene:
  ```bash
  git status
  ```
- Execute the atomic release commit via `gitmap cpr`:
  ```bash
  gitmap cpr "v6.628.0 gitmap user change commands and autobot identity"
  ```
- Inspect commit log to ensure compliance:
  ```bash
  git log -1 --oneline
  ```
  *(Format MUST match `"v6.628.0 <summary>"`; redundant prefixes like `VERSION -` or `RELEASE -` are banned).*
- Transition parent plan `.ai-memory/plans/pending/347-gitmap-user-change-command.md` status to `COMPLETED` and archive to `.ai-memory/plans/completed/`.

---

## 3. Acceptance Criteria Checklist

- [ ] `cli/cmduser/user_help_menu.go` contains entries for `change-user`, `change-login`, and `autobot`.
- [ ] `cli/cmd/rich_help_dispatcher.go` dispatches `gitmap help autobot` to `cmduser.RenderUserHelp()`.
- [ ] Global Git configuration on current machine is set to `user.name = "AI automation bot"` and `user.email = "local@local"`.
- [ ] `gitmap user info` displays the updated identity without errors.
- [ ] `37-bump-version.py` advances `version.json` to `6.628.0`.
- [ ] `14-version-sync-checker.py` reports green status across `version.json`, `package.json`, and `changelog.md`.
- [ ] Release commit is created strictly via `gitmap cpr "v6.628.0 ..."` with zero uncommitted changes remaining.
- [ ] All code adheres strictly to positive booleans, error sovereignty, <= 15-line functions, and relative paths.
