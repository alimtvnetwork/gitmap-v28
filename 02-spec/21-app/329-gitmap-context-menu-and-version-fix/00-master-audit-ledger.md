# 329 — GitMap Context Menu and Version Fix: Master Audit Ledger & Specification

> **Status:** Active Execution (Phase 1 Planning & Specification)  
> **Sequence ID:** 329  
> **Slug:** `gitmap-context-menu-and-version-fix`  
> **Root Directory:** `02-spec/21-app/329-gitmap-context-menu-and-version-fix/`  
> **Plan Ledger:** `.ai-memory/plans/pending/329-gitmap-context-menu-and-version-fix.md`  
> **Rule Compliance:** Strict Leaf Error Sovereignty, Enums with `Type` suffix, Function limit <= 15 lines, Positive Booleans (`is*`/`has*`), Relative Git Paths only.

---

## 1. User Request (Verbatim)

```text
Hi there. In the `gitmap`, recently I have updated the code. It looks like there is a version problem. You need to fix that. I'm going to give you a screenshot from the update that looks weird. You need to fix that. Along with that, I want you to improve and add new context menu option, which needs to be installed by the user by installing `gitmap`, install context menu. And context menu will have a right mouse button click for Windows, Linux, Ubuntu, macOS, all of the OS. You need to write the spec first and then get into that. There I want to have a Git section. In the Git section, we will have some sort of menu like TortoiseGit. I'll give you some screenshot of the TortoiseGit. TortoiseGit can work in a folder. It would have this type of feature list for the TortoiseGit, which I'm giving you as a screenshot. I don't want to make every one of the features. Commit to main, yes, I want. It shouldn't be TortoiseGit, so it would be `gitmap`, and then inside this, it would have diff, diff later, diff later with previous, and show log. These are the three options I want. I also want a branch comparison option where I could compare with two or more branches. In the show log option I mostly use in TortoiseGit, and that actually gives a very nice UI view for that file or folder. Whatever changes happen for that folder or file shows it like this, very nice way. I think you can work on it. Also, we could do refresh button, we could have a view option, and also double-click will open the changes. Clicking on a section also shows which file got modified there, how it modified, and double-click will open the comparison in a compare tool. That is very important. Here you can see already I have selected. I can see which one has the modified. If I double-click it, it will open in editor tool. For my case, I have BeyondCompare. If the BeyondCompare is installed, it will automatically check and take it. You can also get my BeyondCompare settings, which I think we have given in the past. That needs to be there. Also, there it needs to have a install section in the sub-item. Sub-item will have Git and the install section. If you not select any file, then it would have install section, and install will open an install wizard, what we have installed, custom installer is there, which profile is installed, and there should be a uninstall button which can uninstall across all the OS.
Write the spec first and then execute through subagents.
```

---

## 2. Visual Artifacts & Screenshot Ingestion

The user uploaded 4 critical visual references preserved under `assets/screenshots/`:

| Artifact Name | Relative Path | Key Architectural Insight |
| :--- | :--- | :--- |
| **Screenshot 01** | `![Update Downgrade](assets/screenshots/gitmap-context-menu-and-version-fix-01.png)` | `gitmap update` running on `v6.613.0` resolved `v6.610.0` from `sibling-probe` and downgraded the binary. Strict SemVer check `targetVer > currentVer` was missing, and `sibling-probe` accepted an older repo release without validating monotonicity. |
| **Screenshot 02** | `![TortoiseGit Context Menu](assets/screenshots/gitmap-context-menu-and-version-fix-02.png)` | Explorer right-click context menu demonstrating hierarchical Git items: `Diff`, `Diff later`, `Diff with previous version`, `Show log`, `Branch compare`. |
| **Screenshot 03** | `![TortoiseGit Show Log](assets/screenshots/gitmap-context-menu-and-version-fix-03.png)` | TortoiseGit Log Messages dual-pane dialog: Upper pane lists commits (Graph, Revision, Author, Date, Message); Middle pane shows full commit message; Lower pane lists modified files (Action, Path). Double-click invokes external compare tool. |
| **Screenshot 04** | `![TortoiseGit Commit Diff](assets/screenshots/gitmap-context-menu-and-version-fix-04.png)` | Context actions on individual modified files within the log dialog, confirming immediate invocation of Beyond Compare for visual diffs. |

---

## 3. The 5 Mandatory Pillars

### Pillar 1: Strategic Intent & Systems Reasoning
- **Update Engine Monotonicity:** Protect user binaries against regressions. Updating a tool must never silently install an older version. Both `sibling-probe` repository discovery and GitHub release resolution must enforce strict SemVer monotonicity (`targetVersion > currentVersion`). Downgrades require explicit user intent via `--allow-downgrade`.
- **Operating System Context Menu Integration:** Seamlessly bridge the file explorer with developer workflows on Windows (Registry cascades under `HKCU`), Linux/Ubuntu (Nautilus Python extensions & `.desktop` actions), and macOS (Automator Quick Actions).
- **Interactive Dual-Pane Log Viewer & Visual Diff Tool:** Provide instant developer ergonomics equivalent to TortoiseGit directly within GitMap, auto-detecting Beyond Compare (`BComp.exe`, `BCompare.exe`, `/usr/bin/bcompare`) with zero-configuration fallback to VS Code or terminal diffs.
- **Unified Install Wizard:** Enable users to audit installed packages, active profiles, and perform cross-platform clean uninstallation directly from the file explorer background menu.

### Pillar 2: Explicit Scope Boundaries
- **In-Scope:**
  1. `cli/cmdupdate/`: Update version target validation (`update_version_target.go`), sibling probing (`updateprobe.go`), release candidate filtering (`update_fallback_loop.go`), `--allow-downgrade` flag handling.
  2. `cli/cmdinstall/`: Context menu aliases (`context-menu`, `contextmenu`, `ctx`), Git submenu refinement (`installctxentries.go`), Install Wizard launch action.
  3. `cli/cmdcontextmenu/`: Dedicated context-menu command dispatcher (`install`, `uninstall`, `status`, `exec`), cross-platform OS handlers for Windows, Linux, and macOS.
  4. `cli/cmdlogviewer/`: Embedded dual-pane log viewer HTTP service & web UI (`gitmap log ui`, `gitmap show-log`), commit list API, commit files API, refresh action, Beyond Compare launcher integration (`cli/cmdgit/compare_launcher.go`).
  5. `cli/cmdinstall/`: Install wizard UI endpoint and CLI launcher (`gitmap install --wizard` or `gitmap wizard`).
- **Excluded / Deferred:**
  - Full reimplementation of all 40+ TortoiseGit features (only the requested core set: Commit to main, Diff, Diff later, Diff with previous, Show log, Branch comparison).
  - Modifying global system registry keys outside `HKCU` (all Windows registrations remain strictly scoped to current user `HKCU\Software\Classes`).

### Pillar 3: Non-Negotiable Constraints & Negative Guardrails
- **Search Ban:** Strictly no `grep`, `git grep`, `rg`, `ripgrep`, or `Select-String`. Use GitMap native tools (`gitmap find`, `gitmap aum search`, `gitmap cat`).
- **Path Hygiene:** Strictly relative Git paths only (`02-spec/...`, `.ai-memory/...`). No absolute paths, drive letters, or `file:///` URLs.
- **Leaf Error Sovereignty:** All leaf functions return `error` or `*apperror.AppError`. No `cliexit.HandleError` or `os.Exit` in leaves.
- **Boolean Naming:** Prefix all boolean variables and struct fields with `is` or `has` exclusively. Positive booleans only.
- **Function Size:** Functions <= 15 lines max (<= 8 lines preferred). Flatten all nested conditionals with early returns.
- **Rule R1:** No heavy compilation or global test runs during authoring. Run fast static linters (`fastgate.py`, `check-relative-paths.py`) and targeted tests only.

### Pillar 4: Blast Radius & Dependency Topology
- **Update Subsystem (`cli/cmdupdate/`):** Impacts `gitmap update`. Changes strictly guarded by unit tests in `cli/cmdupdate/`.
- **Install Subsystem (`cli/cmdinstall/`):** Impacts `gitmap install`, `gitmap install context-menu`. Windows registry changes only affect `HKCU\Software\Classes`.
- **Log Viewer & Diff Subsystem (`cli/cmdlogviewer/`, `cli/cmdui/`, `cli/cmdgit/`):** Standalone web UI and compare tool invocation. Reuses existing `compare_launcher.go` and `ui_server.go` without disrupting other commands.

### Pillar 5: Multi-Agent Concurrency Topology
- **Capacity:** $A = 3$ concurrent subagents, $H = 2$ hands per agent (Capacity = 6).
- **Subagent Disjoint Domains:**
  - **Worker 01:** Version Update Downgrade Guard & Sibling Probe Resolution (`cli/cmdupdate/`).
  - **Worker 02:** Cross-Platform Context Menu Engine & Install Wizard (`cli/cmdcontextmenu/`, `cli/cmdinstall/`).
  - **Worker 03:** Log Viewer UI & Beyond Compare Diff Integration (`cli/cmdlogviewer/`, `cli/cmdgit/`).
- **Dynamic Join:** As soon as an agent finishes its primary subtask, it immediately joins companion linting, cross-platform validation, and documentation updates.

---

## 4. Discrete Task Extraction & 6-Part Verification Schema

### Task-01: Update Engine Monotonicity Guard & Sibling Probe Fix
- **Deliverable:** Prevent version downgrades in `gitmap update`. Add SemVer comparison (`targetVersion > currentVersion`) to `IsAlreadyUpdated`/`isUpdateSkipped`, add `--allow-downgrade` flag, and update `sibling-probe` to ignore sibling repos that do not contain releases strictly newer than the current version.
- **Understood / AI Comprehension:** [YES] — The update engine must enforce strict monotonic version ordering. If a discovered release is older than or equal to the running version, it must be treated as already up-to-date and skipped, unless the user explicitly passes `--allow-downgrade`. Sibling probe must never switch to a sibling repository with older releases.
- **Azure AI / Model Runtime Context:** N/A — Local CLI/GitMap Go runtime.

### Task-02: Cross-Platform Context Menu Engine (`gitmap install context-menu`)
- **Deliverable:** Implement `gitmap install context-menu` and `gitmap uninstall context-menu` (with aliases `contextmenu`, `ctx`). Wire Windows HKCU registry keys, Linux Nautilus scripts, and macOS Automator workflows.
- **Understood / AI Comprehension:** [YES] — Context menu setup must be easily callable via standard CLI verbs across all supported operating systems. It registers Explorer/Nautilus/Finder actions cleanly scoped to user space and supports clean uninstallation.
- **Azure AI / Model Runtime Context:** N/A — Local CLI/GitMap Go runtime.

### Task-03: Git Context Menu Actions & Install Wizard Sub-Item
- **Deliverable:** Update context menu structure with the TortoiseGit-inspired items: `Commit to main`, `Diff`, `Diff later`, `Diff with previous version`, `Show log`, `Branch comparison`. Include the Install Wizard sub-item on background/folder right-click showing installed components, profile, and uninstall button.
- **Understood / AI Comprehension:** [YES] — File and directory context menus must expose standard TortoiseGit-like operations while non-file/background clicks expose the GitMap Install Wizard for auditing installed components, active profiles, and uninstallation.
- **Azure AI / Model Runtime Context:** N/A — Local CLI/GitMap Go runtime.

### Task-04: Dual-Pane Log Viewer UI & External Diff Tool Integration
- **Deliverable:** Create interactive dual-pane Show Log UI (`gitmap show-log [path]` or `gitmap log ui [path]`). Upper pane lists commits with author/date/message; middle pane shows commit details; lower pane displays modified files. Double-clicking any file opens Beyond Compare (or configured compare tool) comparing the commit state against its parent or working tree.
- **Understood / AI Comprehension:** [YES] — The UI must replicate the ergonomics of TortoiseGit's Log Messages window with refresh, commit inspection, and double-click diff launching through discovered Beyond Compare binaries or fallback editors.
- **Azure AI / Model Runtime Context:** N/A — Local CLI/GitMap Go runtime.

### Task-05: Quality Verification, Documentation & Release Ceremony
- **Deliverable:** FastGate static verification (`50-fastgate.py`), relative path audit (`check-relative-paths.py`), targeted unit tests, version bump ceremony, and prompt/skill synchronization.
- **Understood / AI Comprehension:** [YES] — Code must adhere 100% to GitMap coding guidelines (no raw exits in leaves, positive booleans, <= 15 lines per function) and verify cleanly through FastGate before the release ceremony.
- **Azure AI / Model Runtime Context:** N/A — Local CLI/GitMap Go runtime.
