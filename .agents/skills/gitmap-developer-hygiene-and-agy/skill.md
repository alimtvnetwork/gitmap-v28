---
name: gitmap-developer-hygiene-and-agy
description: Autonomously maintain and extend Google Antigravity (AGY) tools, developer cache cleaners, process managers, and Chrome profile synchronization.
---

# GitMap Developer Hygiene & Antigravity (AGY) Skill (`gitmap-developer-hygiene-and-agy`)

## Mission & Purpose
This skill provides authoritative architectural guidance, code navigation, and execution rules for managing Google Antigravity (AGY) integrations, developer cache purges, operating system process hygiene, and Chrome browser profile synchronizations.

---

## 1. Key Architectural Components & Code Map

| Component | Primary Location | Key Responsibilities |
|---|---|---|
| **Antigravity CLI (AGY)** | `cli/cmdagy/` | AGY prompt injection, conversation picker, empty conversation pruner (`prune --except`), help parity, and IPC focus. |
| **Developer Cache Cleaner** | `cli/osclean/`, `cli/cmdos/` | 10-category cache cleaner: Go build cache, npm/bun, pip, Vite, Antigravity logs, temp directories, and system junk. |
| **Process Manager** | `cli/osutil/`, `cli/osuser/` | Safe cross-platform process discovery and termination (`tasklist` on Windows, `ps` on Linux/macOS). |
| **Chrome Profile Picker** | `cli/cmdchromeprofile/` | Chromium `Local State` 13-attribute schema synchronization, OAuth refresh token vault, and orphan profile reconciliation. |
| **OS Action Decoupling** | `cli/cmdschedule/schedule_os.go` | Injectable OS executors (`DefaultOSActionExecutor`, `defaultFileRemover`) to prevent destructive actions during testing. |

---

## 2. Essential Commands

```bash
# Clean developer tools caches across all 10 categories
gitmap clean-dev
gitmap os dev-clean

# Clean Antigravity specific cache, transcripts, and temporary logs
gitmap agy clean-cache

# Prune empty Antigravity conversations while preserving active ones
gitmap agy prune --except <conversation-id>

# Inject prompt directly into active Antigravity window
gitmap agy prompt "<instruction>"

# Inspect and reconcile Google Chrome profiles
gitmap chrome profiles --reconcile
gitmap chrome picker
```

---

## 3. Core Invariants & Engineering Guardrails

1. **Hermetic Test Isolation (Coding Guideline 24):** Unit tests must NEVER trigger real operating system modifications, power state alterations (shutdown/restart), process termination, package manager upgrades, or actual file deletions. Always verify test coverage using mock injectors (`DefaultOSActionExecutor`, `defaultFileRemover`).
2. **Zero Data Loss in Chrome Reconciliation:** Never delete Chrome user data directories or cookie vaults during profile synchronization. Only reconcile orphan entries in `Local State` with existing filesystem directories.
3. **Safe Process Termination:** Always confirm process name, PID, and executable path before attempting termination to prevent accidental termination of critical system processes or parent IDE hosts.
4. **Preserve Transcript Integrity:** When cleaning AGY logs or temp files, preserve `.system_generated/logs/transcript.jsonl` files of active conversations.
5. **No Uppercase ID Acronyms:** All identifiers must follow `ProfileId`, `ConversationId`, `processId`. Uppercase `ID` is prohibited.

---

## 4. Chrome Profile Import/Export Verification Protocol

When conducting end-to-end testing of Chrome profile export and smart-import workflows:
1. **Secret & Log Isolation**: Isolate all test artifacts, emails, and execution traces under `repo-secrets/xx-<repo-name>/<sequence>-<task-name>/` (e.g. `repo-secrets/01-gitmap/01-chrome-import-export-test/`).
2. **Export Phase**: Execute `gitmap cpe "<Profile>" "<Path>.json"` to produce JSON and CSV companion snapshots.
3. **Removal Phase**: Completely delete the target profile from Chrome User Data on disk (`%LOCALAPPDATA%\Google\Chrome\User Data\<Profile>`) and verify removal (`Test-Path` returns `$false`).
4. **Smart-Import Phase**: Execute `gitmap cpi "<Path>.json" "<Profile>"`, verifying the 5-step pipeline (inspection, target resolution, content restoration, extension hints, and Local State registration).
5. **Post-Import Verification**: Verify on-disk recreation of `Bookmarks` and `Preferences` (with patched `profile.name`), split-DB tracking, and discovery in `gitmap cpl`.

