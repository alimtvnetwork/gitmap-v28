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
| **Antigravity CLI (AGY)** | `cli/cmdagy/` | AGY prompt injection, conversation management, cache cleaner (`cache-clear`, `clean-cache`, `ccko`, `cckf`), empty conversation pruner (`remove-empty-convs`), and IPC focus. |
| **Developer Cache Cleaner** | `cli/osclean/`, `cli/cmdos/` | 10-category cache cleaner: Go build cache, npm/bun, pip, Vite, Antigravity logs, temp directories, and system junk. |
| **Process Manager** | `cli/osutil/`, `cli/osuser/` | Safe cross-platform process discovery and termination (`tasklist` on Windows, `ps` on Linux/macOS). |
| **Chrome Profile Manager** | `cli/cmdchromeprofile/` | Chromium `Local State` 13-attribute schema synchronization, OAuth refresh token vault, and orphan profile reconciliation (`reconcile`, `ls`, `observe`, `open`, `reset`). |
| **OS Action Decoupling** | `cli/cmdschedule/schedule_os.go` | Injectable OS executors (`DefaultOSActionExecutor`, `defaultFileRemover`) to prevent destructive actions during testing. |

---

## 2. Essential Commands

```bash
# Clean developer tools caches across all 10 categories
gitmap clean-dev
gitmap os dev-clean

# Clean Antigravity specific cache, transcripts, and temporary logs
gitmap agy cache-clear
gitmap agy clean-cache

# Clean Antigravity cache keeping only top 1 (ccko) or top 5 (cckf) active conversations
gitmap agy ccko
gitmap agy cckf
gitmap agy cache-clear --keep 10

# Purge empty Antigravity conversations across all projects
gitmap agy remove-empty-convs

# List running Antigravity projects and last active project IDs
gitmap agy running-projects ls
gitmap agy last-active-projects 24

# Inject prompt directly into active Antigravity window
gitmap agy prompt "<instruction>"

# Inspect, list, and observe Google Chrome profiles and open tabs
gitmap chrome ls
gitmap chrome observe
gitmap chrome open "Profile 1" "https://github.com"

# Reconcile on-disk directories with Chromium Local State (repair orphan profiles)
gitmap chrome reconcile

# Reset Chrome profile caches safely without profile loss
gitmap chrome reset "Profile 1" --cache-only

# Compare configuration between two Chrome profiles
gitmap chrome diff "Profile 1" "Profile 2"
```

---

## 3. Core Invariants & Engineering Guardrails

1. **Hermetic Test Isolation (Coding Guideline 24):** Unit tests must NEVER trigger real operating system modifications, power state alterations (shutdown/restart), process termination, package manager upgrades, or actual file deletions. Always verify test coverage using mock injectors (`DefaultOSActionExecutor`, `defaultFileRemover`).
2. **Zero Data Loss in Chrome Reconciliation:** Never delete Chrome user data directories or cookie vaults during profile synchronization. Only reconcile orphan entries in `Local State` with existing filesystem directories.
3. **Safe Process Termination:** Always confirm process name, PID, and executable path before attempting termination to prevent accidental termination of critical system processes or parent IDE hosts.
4. **Preserve Transcript Integrity:** When cleaning AGY logs or temp files, preserve `.system_generated/logs/transcript.jsonl` files of active conversations.
5. **No Uppercase ID Acronyms:** All identifiers must follow `ProfileId`, `ConversationId`, `processId`. Uppercase `ID` is prohibited.
