# Parent Task: UI Terminal & Antigravity/VS Code Management

## Overview
This specification details the implementation of a new web-based terminal UI within the GitMap help dashboard (`hd`), as well as comprehensive Antigravity and VS Code project management commands via the GitMap CLI.

## Architectural Goals & Coding Guidelines
- **Booleans**: Must use `is`, `has`, `can`, `should` (e.g., `isReady`, `hasTerminal`). No negatives (e.g., `isNotReady` is banned).
- **Naming**: Strict semantic naming. No `temp`, `data`, `obj`, `Input1`, etc.
- **Functions**: Max 15 lines. Arguments wrapped at 100 characters.
- **Error Handling**: Wrap all errors with `AppError` or equivalent domain-specific wrapper.
- **UI Terminal**: Must support multiple instances (tabs/splits), Ubuntu Mono font, Ubuntu aesthetic, and robust auto-completion.
- **SQLite Suggestion**: Partial matches (e.g., `GH`) during remove/move operations must suggest matching repositories from the SQLite database backup.

## Subtasks Execution Plan

### Task 1: UI Terminal Frontend (Browser)
- **File**: `gitmap/web/terminal.js` / `gitmap/web/terminal.html` (or appropriate web UI directories)
- **Features**:
  - Implement a web-based terminal UI (e.g. using `xterm.js`).
  - Aesthetics: Ubuntu terminal colors (dark eggplant/purple background, white text) and Ubuntu Mono font.
  - Multi-terminal support (tabs or split view).
  - WebSockets or SSE for backend communication.

### Task 2: UI Terminal Backend & Autocomplete
- **File**: `gitmap/cmd/hd_server.go` (or wherever `hd` is hosted)
- **Features**:
  - Expose API endpoints / WebSockets to proxy terminal commands securely.
  - Implement the auto-complete provider for terminal interactions.

### Task 3: SQLite Suggestion Engine
- **File**: `gitmap/store/suggestions.go`
- **Features**:
  - Implement `GetRepoSuggestions(partialSlug string) ([]string, error)` pulling from the SQLite database.
  - Integrate into the `rm` and `mv` commands so typing partial names (e.g., `GH`) offers interactive suggestions.

### Task 4: Antigravity CLI Management (`agy`)
- **Files**: `gitmap/cmd/agy_cmd.go`
- **Features**:
  - `gitmap agy add <repo>`: Add to `~/.gemini/config/projects/`.
  - `gitmap agy rm <repo>` / `del`: Remove project.
  - `gitmap agy ls`: List all projects.
  - `gitmap agy stats`: Show current user info, email, and limitations.
  - `gitmap agy update`: Automatically update the Antigravity instance.

### Task 5: VS Code CLI Management (`vscode`)
- **Files**: `gitmap/cmd/vscode_cmd.go`
- **Features**:
  - `gitmap vscode add <repo>`: Add to `projects.json`.
  - `gitmap vscode rm <repo>` / `del`: Remove project.
  - `gitmap vscode ls`: List VS Code projects.

### Task 6: Release & Architecture Map
- **Files**: `version.json`, `.lovable/memory/release-architecture-map.md`
- **Features**:
  - Document release architecture in `.lovable/memory/release-architecture-map.md`.
  - Bump MINOR version in `version.json`.
  - Perform release ceremony.

## Verification Checklist (Pre-Commit)
- [ ] Coding Guidelines & Master Consolidated File enforced.
- [ ] Boolean Examples & Fixations strictly followed.
- [ ] Anti-Garbage Naming enforced (no generic temp names).
- [ ] Semantic Tests.
- [ ] Function Size <= 15 lines.
- [ ] Error Handling uses wrappers.
- [ ] Code adheres to explicit booleans, Type-suffixed Enums.
- [ ] Formatting & Acronyms strictly PascalCase (e.g., `SwapIpWindows`).
- [ ] Temp-Scripts ignored.


### Task 7: Update Output Summary
- **Files**: `gitmap/cmd/update.go`, `gitmap/cmd/release.go`
- **Features**:
  - Summarize version change (`vOLD -> vNEW`).
  - Read and print the last 2 lines of the changelog.
  - Append exactly 2 empty lines at the end.
