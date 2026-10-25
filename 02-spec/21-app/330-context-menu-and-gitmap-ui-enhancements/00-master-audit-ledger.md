# 330 — Context Menu and GitMap UI Enhancements: Master Audit Ledger & Specification

> **Status:** Active Execution (Phase 1 Planning & Specification)  
> **Sequence ID:** 330  
> **Slug:** `context-menu-and-gitmap-ui-enhancements`  
> **Root Directory:** `02-spec/21-app/330-context-menu-and-gitmap-ui-enhancements/`  
> **Plan Ledger:** `.ai-memory/plans/pending/330-context-menu-and-gitmap-ui-enhancements.md`  
> **Rule Compliance:** Strict Leaf Error Sovereignty, Enums with `Type` suffix, Function limit <= 15 lines, Positive Booleans (`is*`/`has*`), Relative Git Paths only.

---

## 1. User Request (Verbatim)

```text
The context menu is created, that's fine. But there are several issues. First of all, the context menu does not appear to be in the install help, does not appear on the root level help where it should be. Actually, because context menu is very important. The context menu appears, that's all right, but it does not give the installation command when we right-click on a folder, and repos does not have a dropdown. Serious issue. Now, coming back to the menu. It should also have a menu where we could customize and fix other menus. For example, the PowerShell I have two times. I should be able to remove that using `gitmap` and also from `gitmap` UI, add that option. Now, if I click on a Git repository, right-click on that, where I could see the Tortoise Git diff, diff later, diff with previous version, show log, revision graph, check modified, stash. I see all these options. I don't see any option with `gitmap`. Why is that? Figure it out. And also in the `gitmap`, as I said, I should be able to see the branch comparison. I should see a log, which appears like Tortoise Git. None of these are done, and you are saying you have done it. How funny is that? It does not look very professional. Why you are lying to me, I don't understand. Whenever you do some work, make sure that it is consistent and it completes the work. Half-baked stuff, it does not look good. So, it's a trashy work, and I do not appreciate any trashy work. So make sure the previous requirement you understand, follow through. And make sure that each one of these steps, that works. Diff, branch diff, selecting by default the Beyond Compare if it is installed. And there is no install option in the `gitmap` menu. If we just right-click outside of the file, it should have install option for the repositories and things like that. So, no UI for this. You need to work on it as well. And also the SSH nodes, I want to have those commands as well so that this appears. And I could see that. And also the settings of the `gitmap` UI, these are also not there. So I think there's a lot of things that you need to work on, improve, and complete, and then you come back and say, "I have done it." It's not lying game, okay? No one is asking you to lie. And most of your buttons, if I click on the, let's say, list releases, it does nothing. You need to have visualization or probably keep open the terminal where I could see that for each one of the button. If I click on release latest branch, list versions, nothing appears. How you are doing this stuff? It's crazy. Yeah, also you can add make public, make private public option as it needs to be if it is a Git repo. So some of the options will be there if it is a Git repo. Some of the options will not be there if it is not a Git repo. Remember that. So please work on it. Try to understand the requirements as I have drafted in past. You did not work on it, but lied to me.
```

---

## 2. Visual Artifacts & Screenshot Ingestion

The user uploaded 4 critical visual screenshots preserved under `assets/screenshots/`:

| Artifact Name | Relative Path | Architectural Diagnosis & Finding |
| :--- | :--- | :--- |
| **Screenshot 01** | `![Folder Background Context Menu](assets/screenshots/context-menu-and-gitmap-ui-enhancements-01.png)` | Background context menu in `D:\work\cat-my` shows `gitmap ▸` with only `Scan`, `Clone`, `Release`, `Repos`. Items are truncated because Explorer limits multi-level cascades. Below it, duplicate `PowerShell 7` and `PowerShell` appear. |
| **Screenshot 02** | `![Repos Dropdown Failure](assets/screenshots/context-menu-and-gitmap-ui-enhancements-02.png)` | Hovering on `Repos ▸` fails to open a dropdown. Root cause: Explorer static cascades cannot nest 3 levels deep (`gitmap -> Repos -> Go repos`), and child actions had `Mode: Silent`. |
| **Screenshot 03** | `![Duplicate PowerShell Submenus](assets/screenshots/context-menu-and-gitmap-ui-enhancements-03.png)` | Explorer shows both `PowerShell 7` and legacy `PowerShell` submenus. Registry analysis shows duplicate keys `PowerShellMenu` and `PowerShell7x64` in `HKLM`. Needs cleaner tool in GitMap. |
| **Screenshot 04** | `![File Right-Click Missing GitMap](assets/screenshots/context-menu-and-gitmap-ui-enhancements-04.png)` | Right-clicking on `.prettierignore` displays TortoiseGit (`Diff`, `Diff later`, `Diff with previous version`, `Show log`, etc.) but **zero** GitMap entries! Root cause: `HKCU\Software\Classes\*\shell\gitmap` was omitted from registration. |

---

## 3. The 5 Mandatory Pillars

### Pillar 1: Strategic Intent & Systems Reasoning
- **File & Folder Context Menu Parity:** Developers interact with Git through the filesystem. When right-clicking on a file, GitMap MUST appear with file-scoped actions (`Diff`, `Diff with previous version`, `Show log`, `Diff later`, `Blame`). When right-clicking on a directory or background, GitMap MUST display repository actions directly without deep nesting failures.
- **PowerShell Duplicate Remediation:** Windows installations often leave redundant `PowerShellMenu` and `PowerShell7x64` registry entries. GitMap will provide a native cleaner via `gitmap context-menu clean-powershell` and the Web UI to disable duplicate entries via user-space `HKCU` overrides without requiring elevation.
- **Terminal & Visual Feedback Invariant:** Silent actions (`Mode: Silent`) that piped output to `msg.exe` fail silently on modern Windows. All informative commands (`list-releases`, `list-versions`, `logs`) must run in an interactive terminal (`pwsh -NoExit`) or open the Web UI dashboard so the user sees results immediately.
- **Web UI Completeness:** The Web UI (`cmdui`) must expose the SSH nodes commands, Settings configuration (Beyond Compare path, Context Menu cleaner), and direct links to the Show Log and Branch Compare tools.

### Pillar 2: Explicit Scope Boundaries
- **In-Scope:**
  1. `cli/cmdinstall/installctx.go`: Add `HKCU\Software\Classes\*\shell\gitmap` for file context menus. Flatten directory/background menus to eliminate nested cascade rendering bugs.
  2. `cli/cmdinstall/installctxentries.go`: Define discrete file entries and folder entries. Change `Silent` commands to `Terminal` (`pwsh -NoExit`) so output is always visible.
  3. `cli/cmdcontextmenu/clean_powershell.go`: Implement duplicate PowerShell detection and suppression in `HKCU\Software\Classes\Directory\Background\shell` and `Directory\shell`.
  4. `cli/cmd/roothelp_clusters.go` & `cli/cmdinstall/install_help_sections.go`: Add `context-menu (ctx)` to root help and install help.
  5. `cli/cmdui/`: Add Context Menu & Shell settings card, duplicate PowerShell cleaner toggle, and direct navigation to Log Viewer.
- **Excluded / Deferred:**
  - Machine-level registry deletions requiring elevation (`HKLM`): All overrides remain safely scoped to `HKCU`.

### Pillar 3: Non-Negotiable Constraints & Negative Guardrails
- **Search Ban:** Strictly no `grep`, `git grep`, `rg`, `ripgrep`, or `Select-String`. Use GitMap native tools (`gitmap find`, `gitmap aum search`, `gitmap cat`).
- **Path Hygiene:** Strictly relative Git paths only (`02-spec/...`, `.ai-memory/...`).
- **Leaf Error Sovereignty:** Leaf functions return `error` directly. No `cliexit.HandleError` or `os.Exit` in leaves.
- **Function Limit:** Functions <= 15 lines max. Zero nested ifs. Positive booleans only (`is*`/`has*`).
- **Rule R1:** No heavy compilation during authoring. Run fast static linters (`fastgate.py`, `check-relative-paths.py`) and targeted tests only.

### Pillar 4: Blast Radius & Dependency Topology
- **Context Menu Registrations (`cli/cmdinstall/`):** Scoped to `HKCU\Software\Classes\*\shell\gitmap`, `Directory\shell\gitmap`, `Directory\Background\shell\gitmap`.
- **CLI Help Output (`cli/cmd/`, `cli/cmdinstall/`):** Purely presentation updates.
- **PowerShell Cleaner (`cli/cmdcontextmenu/`):** Only sets `LegacyDisable` in user registry keys; completely reversible.

### Pillar 5: Multi-Agent Concurrency Topology ($A = 3$, $H = 2$, Capacity = 6)
- **Worker 01:** Context Menu Registry Architecture & File Right-Click Integration (`HKCU\Software\Classes\*\shell\gitmap`, directory flattening, terminal mode conversion).
- **Worker 02:** Duplicate PowerShell Cleaner & Context Menu Customization CLI (`cli/cmdcontextmenu/clean_powershell.go`, `gitmap context-menu clean-powershell`).
- **Worker 03:** CLI Help Visibility & Web UI Enhancements (Root/Install help menus, Settings card, SSH node commands, and Log Viewer linkage).

---

## 4. Discrete Task Extraction & 6-Part Verification Schema

### Task-01: File-Level Context Menu & Cascade Flattening
- **Deliverable:** Register `HKCU\Software\Classes\*\shell\gitmap` with file diff, diff with previous, show log, and blame. Flatten folder/background context menus to prevent cascade truncation in Windows Explorer. Switch all informational actions to `Terminal` mode (`pwsh -NoExit`).
- **Understood / AI Comprehension:** [YES] — When right-clicking on any file, GitMap must offer TortoiseGit-like file operations using Beyond Compare. In folders, menus must not exceed 2 cascade levels to prevent Windows Explorer from cutting off items or breaking dropdowns. All silent commands that produced invisible output must open an interactive terminal.
- **Azure AI / Model Runtime Context:** N/A — Local CLI/GitMap Go runtime.

### Task-02: Duplicate PowerShell Context Menu Cleaner
- **Deliverable:** Implement `gitmap context-menu clean-powershell` (and `gitmap context-menu customize`). Detect duplicate `PowerShellMenu` / `PowerShell7x64` keys and disable legacy duplicates via `LegacyDisable` in `HKCU\Software\Classes`.
- **Understood / AI Comprehension:** [YES] — Windows Explorer exhibits duplicate PowerShell entries when both PowerShell 7 and legacy PowerShell menus are registered. GitMap must provide an automated, reversible mechanism to suppress duplicate keys in user space.
- **Azure AI / Model Runtime Context:** N/A — Local CLI/GitMap Go runtime.

### Task-03: CLI Help Integration (Root & Install Help)
- **Deliverable:** Add `context-menu (ctx)` to root CLI help (`cli/cmd/roothelp_clusters.go`) under System/OS Tooling, and to installer help (`cli/cmdinstall/install_help_sections.go`, `install_help_menu.go`).
- **Understood / AI Comprehension:** [YES] — The context menu command is a primary developer entrypoint and must be prominently documented in both `gitmap --help` and `gitmap install --help`.
- **Azure AI / Model Runtime Context:** N/A — Local CLI/GitMap Go runtime.

### Task-04: GitMap Web UI Enhancements (Settings, SSH Commands & Log Viewer Link)
- **Deliverable:** Expand the Web UI with a Context Menu & Shell Settings section, SSH Fleet CLI command cards, and direct navigation to the Git Log & Diff viewer.
- **Understood / AI Comprehension:** [YES] — The Web UI must allow configuring context menu settings, cleaning duplicate PowerShell menus via a 1-click button, inspecting SSH nodes and fleet commands, and navigating between log/diff and repo management.
- **Azure AI / Model Runtime Context:** N/A — Local CLI/GitMap Go runtime.

### Task-05: Quality Verification, Documentation & Minor Release Ceremony
- **Deliverable:** FastGate verification, relative paths audit, targeted unit tests, version bump to `v6.615.0`, and release ceremony.
- **Understood / AI Comprehension:** [YES] — Code must pass FastGate 100%, adhere strictly to coding guidelines, and complete with a clean release ceremony.
- **Azure AI / Model Runtime Context:** N/A — Local CLI/GitMap Go runtime.
