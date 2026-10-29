# Context Menu Icons and Non-Git Error Handling Architecture Specification

## 1. Overview & Problem Statement
When invoking GitMap actions such as `gitmap show-log` or `gitmap diff` from the operating system right-click context menu (or terminal) within a non-git directory (such as drive roots `D:\` or non-tracked project directories):
1. The command failed with `✗ '.' is not a Git repository.` and tip text, and terminated without launching the Web UI, leaving users stranded at `PS D:\>`.
2. The context menu on directory background retained orphaned / stale git action entries (`01_show_log`, `02_diff`, `03_branch_compare`, `04_commit`, `05_pull`, `06_push`) because registry reinstallation did not wipe the old `shell` verb tree prior to installing.
3. Context menu items universally reused the default `{exe},0` GitMap icon or had no icon at all, producing a repetitive, low-contrast UI experience.
4. Users lacked a clean mechanism to toggle or hide git-specific actions in the context menu when working primarily outside of git repositories.

## 2. Architecture & Design Principles
- **Graceful UI Fallback**: When `show-log` or `diff` is triggered outside a git repository and browser launch is enabled (`!isNoBrowser`), GitMap deeply scans child repositories (depth 1 to 3) and launches the Web UI repository dashboard (`FallbackDashboardUIFn("search", 8080)`), eliminating dead-end terminal errors.
- **Surgical Stale Verb Purge**: Registry installation purges existing `shell` subtrees under `HKCU\Software\Classes\Directory\Background\shell\gitmap\shell`, `HKCU\Software\Classes\Directory\shell\gitmap\shell`, and `HKCU\Software\Classes\*\shell\gitmap\shell` before writing fresh entries, ensuring orphaned keys cannot linger.
- **Universal High-Fidelity Icons**: Every context menu entry and submenu category is assigned a distinct, standard Windows system icon index (`shell32.dll` and `imageres.dll`), ensuring crisp rendering across Windows 10, 11, and Server.
- **Git Action Visibility Toggle**: Introduce CLI flags (`--hide-git-actions`), CLI commands (`gitmap context-menu hide-git-actions`, `show-git-actions`, `toggle-git-actions`), and Settings Web UI integration (`/api/context-menu/toggle-git-actions`).
- **Coding Guideline Conformance**: Strict compliance with functions <= 15 lines, zero nested `if` statements, positive booleans (`is*`/`has*`), leaf error sovereignty, and strict relative file paths.
