# Master Audit Ledger: GitMap Installation & UI Enhancement Instructions

- **Sequence ID**: 334
- **Task**: `GitMap Installation and UI Enhancement Instructions`
- **Slug**: `gitmap-installation-and-ui-enhancement-instructions`
- **Timestamp**: 2026-10-27T00:15:00+08:00
- **Status**: IN_PROGRESS

---

## 1. Problem Statement & User Requirements Forensic Decomposition

### Requirement 1: Installation Command Accessible from Right-Hand / Context Interface
- **Observation**:
  - The Explorer right-click context menu previously added an `Install ▸` submenu with sub-actions (wizard, dev profile, tools list, antigravity, clean-powershell).
  - However, when right-clicking on empty background or a directory, the user requested an explicit direct right-hand/context command to install or manage context menu and installation options directly (`gitmap install context-menu` or an `Install Context Menu` entry).
  - Furthermore, `gitmap install` right-hand column formatting in terminal help and CLI commands must clearly showcase the installation of the context menu (`gitmap install context-menu` / `gitmap in ctx`).
  - In `cli/cmdinstall/installctxentries.go`, ensure right-hand context entries include an explicit action `Install / Reinstall Context Menu` under `60_install/05_install_ctx` so users can re-register or update the context menu directly from the context menu itself.

### Requirement 2: Open Project in Antigravity IDE (Cross-Platform)
- **Observation**:
  - `gitmap open` currently only recognizes `code` and `vscode` as editor targets in `cli/cmdopen/open_target.go`.
  - When users run `gitmap open agy`, `gitmap open antigravity`, or `gitmap open ag`, it fails to recognize Antigravity as an editor and treats it as a non-existent directory.
  - While `gitmap agy` has `RunAgyOpen(path)`, `gitmap open [agy|antigravity|ag] [path]` should cleanly launch the Google Antigravity Desktop IDE across Windows, Linux, and macOS.
  - Additionally, add a top-level shorthand command `gitmap open-agy [path]` / `gitmap open-antigravity [path]` (and `gitmap agy open [path]`), and make sure `ResolveAntigravityIDE()` locates `Antigravity.exe` (Windows), `Antigravity.app` (macOS), or Linux desktop binary (`/opt/antigravity/antigravity` or `antigravity` in PATH).

### Requirement 3: App-Like Browser Window with Exit / Close Window Endpoint
- **Observation**:
  - When opening front-end UIs in the browser (`gitmap ui`, `gitmap show-log`, `gitmap ssh ui`), opening plain default browser tabs without window sizing feels like a website rather than a native standalone developer application.
  - Google Chrome and Edge support `--app=http://...` mode, which launches the web UI in an isolated frameless/standalone app window without browser URL bars, tabs, or bookmarks.
  - The browser opener (`openBrowserURL` in `cli/cmdui/ui_server.go`, `openBrowser` in `cli/cmdlogviewer/server.go`, and `openBrowserFallbackURL` in `cli/cmdssh/ssh_ui.go`) should prioritize launching in app-window mode (`--app=http://...`) if Chrome or Edge is detected, with graceful fallback to standard system browser start.
  - Every front-end UI (`cmdui` Fleet UI, `cmdlogviewer` Show Log UI) must have:
    - An explicit **Exit / Close Window** button in the header toolbar (`✕ Exit` / `Close Window`).
    - An active server shutdown endpoint (`/api/ui/exit` in `cmdui`, `/api/log/exit` in `cmdlogviewer`) that closes the HTTP server and exits the background Go process.
    - Client-side JavaScript that calls the exit endpoint, notifies the user with a graceful closing toast/overlay, and invokes `window.close()`.

### Requirement 4: Full-Screen Capability
- **Observation**:
  - Modern desktop web apps support instant full-screen toggling.
  - Add a dedicated Full-Screen toggle button (`⛶ Fullscreen`) in the top navigation bar / header of both the Fleet UI (`cmdui`) and Show Log UI (`cmdlogviewer`).
  - Wire JavaScript using standard HTML5 Fullscreen API (`document.documentElement.requestFullscreen()` and `document.exitFullscreen()`), with keyboard shortcut `F11` / `Alt+Enter` handling and visual icon state change.

### Requirement 5: Professional, Lucrative UI, Hover Effects & CSS3 Animations
- **Observation**:
  - The user noted: *"Currently, the animations or the UI you have created, these are very poor. Make sure that this has very nice and detailed UI systems."*
  - Refactor and elevate the styling across `cli/cmdui/ui_assets_theme.go` and `cli/cmdlogviewer/ui_assets_html.go`:
    - **Color Guiding & Theme Tokens**: Refine Catppuccin Mocha & Macchiato depth planes (`#1e1e2e`, `#181825`, `#11111b`, `#89b4fa`, `#a6e3a1`, `#fab387`, `#cba6f7`).
    - **Lucrative Micro-Interactions**: Smooth cubic-bezier transitions (`transition: all 0.22s cubic-bezier(0.16, 1, 0.3, 1)`), subtle glow hover effects (`box-shadow: 0 4px 20px -2px rgba(137, 180, 250, 0.25)`), dynamic border highlights, and scale hover feedback (`transform: translateY(-1px)`).
    - **CSS3 Animations**: Add silky keyframe animations:
      - `@keyframes appEntrance`: smooth opacity & slide-up entrance for cards and panes.
      - `@keyframes pulseGlow`: subtle pulsing glow on active status badges and connection dots.
      - `@keyframes shimmer`: elegant subtle loading shimmer on tables and data rows.
      - Smooth ripple/glow effect on button clicks.

---

## 2. Forensic Audit Matrix

| Component | Target File | Current State | Required State |
|---|---|---|---|
| Context Menu Installer | `cli/cmdinstall/installctxentries.go` | Submenu has wizard, dev, tools, agy, clean-pwsh | Add `05_install_ctx`: "Install / Reinstall Context Menu" (`gitmap install context-menu`) |
| Context Menu Argv Test | `cli/cmdinstall/installctxentries_argv_test.go` | Validates menu hierarchy | Add expected entry for `60_install/05_install_ctx` |
| Antigravity Open | `cli/cmdopen/open_target.go` | `isEditorTarget` only checks `code`/`vscode` | Add `isAntigravityTarget` (`agy`, `antigravity`, `ag`); invoke `cmdagy.RunAgyOpen(path)` |
| Open Antigravity CLI | `cli/cmdopen/open.go` | Only calls `launchEditor` | Route `isAntigravityTarget(args[0])` directly to `cmdagy.RunAgyOpen(path)` |
| Top-level CLI Route | `cli/cmd/roottooling.go` | No `open-agy` or `open-antigravity` | Add `open-agy`, `open-antigravity` routing to `cmdagy.RunAgyOpen` |
| Fleet UI Exit Endpoint | `cli/cmdui/ui_server.go` | No `/api/ui/exit` endpoint | Add `/api/ui/exit` handler that closes `http.Server` and exits process |
| Fleet UI Header Exit & Fullscreen | `cli/cmdui/ui_assets_markup_misc.go` | Header only has title & status indicator | Add Fullscreen (`⛶`) and Exit (`✕ Close App`) buttons in header |
| Fleet UI JS Actions | `cli/cmdui/ui_assets_js_settings.go` | No exit or fullscreen logic | Add `toggleAppFullscreen()` and `exitApp()` methods |
| Show Log Header Fullscreen | `cli/cmdlogviewer/ui_assets_html.go` | Has `#btn-exit` | Add `#btn-fullscreen` (`⛶ Fullscreen`) and enhanced app styling |
| Show Log JS Fullscreen | `cli/cmdlogviewer/ui_assets_script.go` | Exit button wired | Wire `#btn-fullscreen` to `document.fullscreenElement` API |
| App Mode Opener | `cli/cmdui/ui_server.go` & `cli/cmdlogviewer/server.go` | `openBrowserURL` launches default URL | Attempt `--app=<url>` via Chrome/Edge for true app window behavior |
| CSS3 Animation & Glow | `cli/cmdui/ui_assets_theme.go` & `cli/cmdlogviewer/ui_assets_html.go` | Basic static styling | Rich CSS3 keyframes, cubic-bezier transitions, glassmorphic cards, glow hovers |
