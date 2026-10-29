# Specification: Machine Migration, Help Menu Integration, and Release Ceremony

> **Spec ID:** `SPEC-21-347-MIG`  
> **Subsystem:** `cli/cmduser`, `cli/cmd`, `cli/usercontext`, `03-ai-scripts`  
> **Status:** Active / Specification  
> **Date:** 2026-10-29  
> **Related Specs:**
> - [readme.md](readme.md)
> - [01-architecture-spec.md](01-architecture-spec.md)
> - [02-component-spec.md](02-component-spec.md)
> **Related Plans:**
> - [.ai-memory/plans/pending/347-gitmap-user-change-command.md](../../../.ai-memory/plans/pending/347-gitmap-user-change-command.md)
> - [.ai-memory/plans/subtasks/347-gitmap-user-change-command/03-machine-migration-and-release.md](../../../.ai-memory/plans/subtasks/347-gitmap-user-change-command/03-machine-migration-and-release.md)

---

## 1. Executive Summary & Migration Objective

This specification details the machine-level identity migration protocol, CLI rich help menu integration, and automated version bump release ceremony for Plan 347 (`gitmap-user-change-command`).

### Migration Objectives
1. **Machine Identity Migration:** Migrate the host developer machine's global Git identity from personal author credentials to the dedicated automation identity:
   - `user.name`: `"AI automation bot"`
   - `user.email`: `"local@local"` (or empty unset if configured without email)
2. **Rich Help Menu Registration:** Update `cli/cmduser/user_help_menu.go` and `cli/cmd/rich_help_dispatcher.go` so that user switching subcommands (`change-user`, `change-login`, `autobot`) render formatted two-column terminal menus and dispatch seamlessly from root help topics.
3. **Automated Minor Version Bump:** Advance GitMap from `v6.627.0` to `v6.628.0` via `03-ai-scripts/37-bump-version.py`.
4. **Semantic Release Ceremony:** Execute the atomic release commit via `gitmap cpr` strictly adhering to the formatting pattern `v<version> <summary>`.

---

## 2. Machine Git Identity Migration Blueprint

### 2.1 State Transition Matrix

The migration operates on the global Git configuration level (`~/.gitconfig` or OS-equivalent `%USERPROFILE%\.gitconfig`).

```
+---------------------------------------------------------------------------------------------------+
| 1. CURRENT GLOBAL GIT STATE                                                                       |
| user.name  = "MD ALIM UL KARIM"                                                                  |
| user.email = "devorg.bd@gmail.com"                                                                |
+---------------------------------------------------------------------------------------------------+
                                                  |
                                                  v
+---------------------------------------------------------------------------------------------------+
| 2. MIGRATION INVOCATION                                                                           |
| gitmap autobot "local@local" --global                                                             |
|   OR: git config --global user.name "AI automation bot"                                           |
|       git config --global user.email "local@local"                                                |
|   (If empty unsetting requested: git config --global --unset user.email)                          |
+---------------------------------------------------------------------------------------------------+
                                                  |
                                                  v
+---------------------------------------------------------------------------------------------------+
| 3. TARGET GLOBAL GIT STATE                                                                        |
| user.name  = "AI automation bot"                                                                  |
| user.email = "local@local" (or unset / empty string)                                              |
+---------------------------------------------------------------------------------------------------+
                                                  |
                                                  v
+---------------------------------------------------------------------------------------------------+
| 4. VERIFICATION                                                                                   |
| git config --global --get user.name  == "AI automation bot"                                       |
| git config --global --get user.email == "local@local"                                             |
| gitmap user info -> reflects active identity as "AI automation bot <local@local>"                 |
+---------------------------------------------------------------------------------------------------+
```

### 2.2 Unset vs Explicit Local Email Handling

The implementation in `cli/usercontext/git_identity.go` distinguishes between setting an explicit email address and unsetting/leaving email empty.

- **Explicit Email (`"local@local"`):**
  When an email argument is supplied, GitMap executes:
  ```bash
  git config --global user.name "AI automation bot"
  git config --global user.email "local@local"
  ```
- **Empty / Unset Email:**
  When email is explicitly passed as empty or omitted during autobot initialization, GitMap preserves existing unsetting semantics:
  ```bash
  git config --global user.name "AI automation bot"
  git config --global --unset user.email
  ```
- **Boolean Guard Conventions:**
  All conditions evaluating configuration states MUST use positive booleans (`is*`, `has*`):
  - `hasEmail`: Evaluates `len(strings.TrimSpace(email)) > 0`.
  - `isGlobal`: Distinguishes global from project/local repository configuration.
  - `isConfigured`: Verifies identity fields are populated.

---

## 3. Rich Help Menu & Help Dispatcher Architecture

### 3.1 Help Menu Updates in `cli/cmduser/user_help_menu.go`

`cli/cmduser/user_help_menu.go` constructs the styled two-column terminal help display via `termout.HelpMenu`. The menu must be updated to feature the new identity commands:

1. **Usage Lines:**
   ```go
   UsageLines: []string{
       "gitmap user <subcommand> [flags]",
       "gitmap change-user <name> [email] [--global|--project]",
       "gitmap change-login <name> [email] [--global|--project]",
       "gitmap autobot [email] [--global|--project]",
       "gitmap user info [--json]",
       "gitmap user switch <alias> [--global|--project]",
   }
   ```

2. **Command Entries in `buildGitUserSection()`:**
   - `change-user / change-login <name> [email]`: "Change Git author identity via fast one-liner"
   - `autobot [email]`: "Switch Git identity to 'AI automation bot' (optional email or local@local)"

3. **Footer Flags in `buildUserFooterFlags()`:**
   - `--global`: Apply identity or profile change globally across the machine.
   - `--project`: Apply identity or profile change to current repository only.

4. **15-Line Function Limit Conformance:**
   Each helper function must remain within the 15-line limit:
   ```go
   func buildGitUserSection() termout.HelpSection {
       return termout.HelpSection{
           Title: "Git Profile & Context Actions",
           Entries: []termout.CommandEntry{
               {Command: "change-user / change-login <name> [email]", Description: "Change Git author identity via fast one-liner"},
               {Command: "autobot [email]", Description: "Switch Git identity to AI automation bot"},
               {Command: "info / status", Description: "Display consolidated Git and GitHub CLI user identity card"},
               {Command: "list / ls", Description: "List configured Git profiles with active and bound indicators"},
               {Command: "switch / use <alias>", Description: "Switch active Git profile (--global or --project)"},
               {Command: "project [bind|unbind]", Description: "Manage per-repository profile bindings"},
               {Command: "config, set [scope]", Description: "Inspect or configure Git user.name and user.email"},
               {Command: "sync", Description: "Apply bound project user profile to current repository"},
           },
       }
   }
   ```

### 3.2 Help Topic Dispatching in `cli/cmd/rich_help_dispatcher.go`

Root command help calls (`gitmap help <topic>` or `gitmap <topic> --help`) query `tryRenderRichTopic(topic string)`.

`cli/cmd/rich_help_dispatcher.go` registers topics in `tryRenderCoreRichTopic(topic string)`:
```go
func tryRenderCoreRichTopic(topic string) bool {
    switch topic {
    // ... existing topic cases ...
    case "user", "os-user", "change-user", "change-login", "autobot":
        cmduser.RenderUserHelp()

        return true
    }

    return false
}
```

This guarantees that invoking `gitmap help autobot`, `gitmap help change-user`, or `gitmap help change-login` immediately surfaces the rich terminal user identity dashboard.

---

## 4. Version Bump Automation & Manifest Synchronization

### 4.1 Single Source of Truth: `version.json`

GitMap adheres to root `version.json` as the single canonical source of truth for version metadata.

- **Current Version:** `6.627.0`
- **Target Version:** `6.628.0`
- **SemVer Tier:** `minor` (incrementing minor version for new command feature additions per Rule 0).

### 4.2 Automation Script Execution

The synchronization script `03-ai-scripts/37-bump-version.py` coordinates multi-manifest propagation:

```bash
gitmap py 03-ai-scripts/37-bump-version.py -t minor -s "gitmap user change commands and autobot identity"
```

### 4.3 Propagated Files
1. `version.json`: Updates `"Version": "6.628.0"`.
2. `package.json`: Updates `"version": "6.628.0"`.
3. `latest.json`: Updates version, tag (`v6.628.0`), and release branch (`release/v6.628.0`).
4. `changelog.md`: Prepends new section:
   ```markdown
   ## [6.628.0] - 2026-10-29
   ### Added
   - Added `gitmap change-user` and `gitmap change-login` commands for fast one-liner Git identity updates.
   - Added `gitmap autobot` command to switch identity to "AI automation bot" with optional email or local@local.
   - Migrated host machine Git identity globally to "AI automation bot".
   - Integrated rich help menus for user change commands in `cli/cmduser` and `cli/cmd`.
   ```
5. `cli/version/constants.go`: Regenerated to embed version constants.
6. `readme.md`: Synchronized version badges and docs.

### 4.4 Verification Check

Post-bump synchronization is validated using:
```bash
gitmap py 03-ai-scripts/14-version-sync-checker.py
```
This guarantees 100% parity across `version.json`, `package.json`, and `changelog.md`.

---

## 5. Release Ceremony & Atomic Commit Invariants

### 5.1 Dedicated Commit Verb: `gitmap cpr`

All release commits in GitMap MUST use the dedicated release verb:
```bash
gitmap cpr "v6.628.0 gitmap user change commands and autobot identity"
```

### 5.2 Strict Commit Discipline Invariants
- **TOTAL BAN on Generic Commit Verbs:** Never use `git commit`, `git add -A`, `gitmap cpf`, or `gitmap cpc` for a release.
- **Strict Format:** Must be `"v<version> <summary>"`. Redundant prefixes such as `VERSION - ` or `RELEASE - ` are strictly banned.
- **Relative Paths Only:** Commit summaries and release notes must reference only relative repository paths (e.g. `cli/cmduser/`, `02-spec/21-app/`). Absolute paths (e.g., drive letters `d:\...`) are prohibited.
- **Clean Working Tree:** The working tree must be fully clean before and after release execution.

---

## 6. Verification Matrix & Acceptance Criteria

| Check # | Requirement | Verification Command / Target | Expected Result |
| :--- | :--- | :--- | :--- |
| **VM-01** | Global Git Name Migration | `git config --global --get user.name` | Returns `"AI automation bot"` |
| **VM-02** | Global Git Email Migration | `git config --global --get user.email` | Returns `"local@local"` (or empty if unset) |
| **VM-03** | User Identity Card Output | `gitmap user info` | Displays "AI automation bot <local@local>" |
| **VM-04** | User Help Menu Command Registration | `gitmap user --help` | Includes `change-user`, `change-login`, `autobot` |
| **VM-05** | Rich Help Dispatcher Routing | `gitmap help autobot` | Renders styled user help menu |
| **VM-06** | FastGate Validation | `gitmap py 03-ai-scripts/50-fastgate.py` | All gates pass (0 errors) |
| **VM-07** | Version Sync Verification | `gitmap py 03-ai-scripts/14-version-sync-checker.py` | `✅ Version synchronization verified: v6.628.0` |
| **VM-08** | Release Commit Audit | `git log -1 --oneline` | Contains `v6.628.0 gitmap user change commands...` |
| **VM-09** | Function Line Length | Code audit across modified Go files | <= 15 lines per function |
| **VM-10** | Boolean Hygiene | Code audit across modified Go files | Positive booleans only (`is*`, `has*`) |
