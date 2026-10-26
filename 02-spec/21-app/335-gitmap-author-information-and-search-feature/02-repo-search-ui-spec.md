# Repository Search Web UI & AUM Integration Specification
**Sequence ID:** 335  
**Component:** `cli/cmdui`, `cli/cmdautomation`, `cli/cmdsearch`  

---

## 1. Architectural Changes
1. **API Endpoint `/api/search/repo`:**
   - Query parameters:
     - `q`: search query string.
     - `dir`: target directory path (defaults to current working directory / repo root).
     - `ext`: comma-separated file extensions filter (e.g., `.go,.ts,.json`).
     - `isRegex`: boolean flag (`"true"` / `"false"`).
     - `isCaseInsensitive`: boolean flag (`"true"` / `"false"`).
   - Core Engine Call:
     - Leverages `cmdautomation.RunSearch(opts SearchOptions) (SearchResult, *apperror.AppError)`.
   - Response Payload:
     - Returns JSON with `matches`: array of `{ path, lineNumber, lineContent }`, `totalFiles`, `totalHits`, and `duration`.

2. **Web UI Tab (`tab-search`):**
   - Added to navigation bar (`NAV_ITEMS`) in `cli/cmdui/ui_assets_js_settings.go`: `{ id: 'search', label: 'Repo Search', icon: 'search' }`.
   - Search input bar with instantaneous / enter-to-search triggers, regex toggle, case-insensitivity toggle, and extension filter chips.
   - Live result cards:
     - Displays matched file relative path, line number badges, highlighted match snippets.
     - Hover tooltips: previews 3 lines of contextual surrounding code on mouse hover.
     - Direct Actions:
       - "Open in Remote Editor" -> switches immediately to `tab-editor` with file loaded at that specific line.
       - "Open in IDE / Antigravity" -> calls `/api/project/open-antigravity?path=...`.

3. **CLI Launcher:**
   - Supports `gitmap search --ui` and `gitmap repo-search --ui` to open browser in full-screen/app mode pointing directly to the `#tab-search` view.
