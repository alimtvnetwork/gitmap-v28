# Architecture Specification: GitMap Installation & Antigravity IDE Integration

- **Sequence**: 334
- **Component**: Installation Interface & Antigravity IDE Cross-Platform Integration
- **Status**: APPROVED

---

## 1. Overview
This specification details:
1. Enhancements to the context menu and CLI installation interface for `gitmap install context-menu` (`gitmap in ctx`).
2. Universal cross-platform integration for opening projects in Google Antigravity Desktop IDE via `gitmap open agy [path]` and top-level alias `gitmap open-agy [path]`.

---

## 2. CLI Command Contracts

### 2.1 Context Menu Installation Enhancements
- In `cli/cmdinstall/installctxentries.go`, under `installChildren()`:
  - Add `{KeyName: "05_install_ctx", MUIVerb: "Install Context Menu", Args: []string{"install", "context-menu"}, Mode: ctxmodetype.Terminal, Icon: constants.CtxIconGitmap}`
  - This ensures that users viewing the context menu or looking for installation commands see the explicit action to install/refresh the context menu on any directory.

### 2.2 Cross-Platform Open Antigravity Contract
- `gitmap open agy [path]` / `gitmap open antigravity [path]` / `gitmap open ag [path]`:
  - Recognizes `agy`, `antigravity`, and `ag` as first argument in `cli/cmdopen/open.go`.
  - Resolves target directory (defaults to cwd or Git top-level root).
  - Invokes `cmdagy.RunAgyOpen(targetPath)`.
- Top-level aliases:
  - `gitmap open-agy [path]`
  - `gitmap open-antigravity [path]`
  - Dispatches directly to `cmdagy.RunAgyOpen(path)`.
- Operating System Support:
  - **Windows**: Detects `Antigravity.exe` in `%LOCALAPPDATA%\Programs\Antigravity\`, `Program Files`, or PATH.
  - **macOS**: Detects `/Applications/Antigravity.app` or executes `open -a Antigravity <path>`.
  - **Linux**: Detects `/opt/antigravity/antigravity` or `antigravity` in PATH.
  - If not installed, prints clear installation suggestion: `gitmap install antigravity`.
