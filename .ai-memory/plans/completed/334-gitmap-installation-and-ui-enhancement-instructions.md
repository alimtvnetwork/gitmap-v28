# Master Plan: GitMap Installation and UI Enhancement Instructions

- **Sequence ID**: 334
- **Task**: `GitMap Installation and UI Enhancement Instructions`
- **Slug**: `gitmap-installation-and-ui-enhancement-instructions`
- **Status**: PENDING
- **Created**: 2026-10-27T00:20:00+08:00

---

## 1. Subtasks Overview

1. [`subtask-01-install-ctx-and-antigravity-open.md`](subtasks/334-gitmap-installation-and-ui-enhancement-instructions/01-install-ctx-and-antigravity-open.md)
   - Add `05_install_ctx` to context menu install submenu (`gitmap install context-menu`).
   - Update `installctxentries_argv_test.go`.
   - Add Antigravity target resolution to `cmdopen` (`gitmap open agy [path]`) and top-level alias `gitmap open-agy [path]`.
2. [`subtask-02-app-mode-browser-launcher.md`](subtasks/334-gitmap-installation-and-ui-enhancement-instructions/02-app-mode-browser-launcher.md)
   - Enhance browser opener (`cmdui/ui_server.go`, `cmdlogviewer/server.go`, `cmdssh/ssh_ui.go`) to launch in `--app=<url>` window mode via Chrome/Edge when available.
   - Fall back smoothly to default browser on other platforms or headless modes.
3. [`subtask-03-fullscreen-exit-and-css3-animations.md`](subtasks/334-gitmap-installation-and-ui-enhancement-instructions/03-fullscreen-exit-and-css3-animations.md)
   - Add `/api/ui/exit` endpoint in `cmdui/ui_server.go`.
   - Add Fullscreen (`⛶`) and Exit (`✕ Close`) header buttons in Fleet UI and Show Log UI.
   - Refactor CSS styling in `cmdui/ui_assets_theme.go` and `cmdlogviewer/ui_assets_html.go` with rich CSS3 animations, hover effects, and Catppuccin theme tokens.
4. [`subtask-04-verification-and-release.md`](subtasks/334-gitmap-installation-and-ui-enhancement-instructions/04-verification-and-release.md)
   - Run FastGate and unit tests.
   - Minor version release ceremony (`v6.617.0`).
   - Recompile local binary and reinstall context menu.
