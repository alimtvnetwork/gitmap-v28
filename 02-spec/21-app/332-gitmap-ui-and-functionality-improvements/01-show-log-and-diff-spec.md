# Specification: Show Log Modern UI, Diff Engine, and Exit Integration

> **Module:** `cli/cmdlogviewer/`, `cli/cmddiff/`, `cli/cmdgit/`  
> **Parent Program:** `332-gitmap-ui-and-functionality-improvements`  
> **File:** `02-spec/21-app/332-gitmap-ui-and-functionality-improvements/01-show-log-and-diff-spec.md`  
> **Status:** Draft / Active Specification  

---

## 1. Show Log Folder Inspection Fix
- **Path Normalization:** Ensure `ResolveTargetRepoAndPath()` cleans folder paths, converts Windows backslashes to forward slashes, and checks if target is a directory or file.
- **Git Log Arguments:** When inspecting a directory, append clean relative path `cli/cmdlogviewer` (without trailing slash) to `git log -- <path>`.

---

## 2. Modern, Responsive UI & Scrolling
- **Custom Scrollbars:** Add modern CSS scrollbars:
  ```css
  ::-webkit-scrollbar { width: 8px; height: 8px; }
  ::-webkit-scrollbar-track { background: var(--pane-bg); }
  ::-webkit-scrollbar-thumb { background: var(--border); border-radius: 4px; }
  ::-webkit-scrollbar-thumb:hover { background: var(--accent); }
  ```
- **Responsive Workspace Grid:**
  - Standard Desktop: 3-pane split (Commits 40%, Commit Detail 25%, Modified Files 35%).
  - Narrow Screens (`@media (max-width: 900px)`): Collapse panes gracefully with scrollable containers.
- **Visual Polish:**
  - Lipgloss Catppuccin-inspired modern color tokens.
  - Keyboard navigation (`ArrowUp` / `ArrowDown` to cycle commits).

---

## 3. Browser Exit & Local Server Shutdown
- Add an "Exit / Close" button in the top toolbar:
  ```html
  <button id="btn-exit" class="btn-danger" title="Stop local server and close window">✕ Exit</button>
  ```
- Backend Handler `/api/log/exit`:
  ```go
  func handleAPIExit(server *http.Server) http.HandlerFunc {
      return func(w http.ResponseWriter, r *http.Request) {
          writeJSON(w, http.StatusOK, map[string]string{"status": "shutting_down"})
          go func() {
              time.Sleep(100 * time.Millisecond)
              _ = server.Close()
              os.Exit(0)
          }()
      }
  }
  ```
- Frontend triggers `window.close()` upon successful shutdown call.

---

## 4. Web UI Right-Click Context Menu & Copy / Revert
- **Custom Context Menu:**
  - Right-clicking any commit row displays a context menu:
    - `Copy Commit Hash` (writes to clipboard)
    - `Copy Subject Line`
    - `Revert Commit` (calls `/api/log/revert` with confirm dialog)
    - `Compare with Previous`
  - Right-clicking any file row:
    - `Copy File Path`
    - `Diff with Previous (Beyond Compare / VS Code)`
    - `Revert File Changes`

---

## 5. Show Diff Single-Path & Working Tree Fallback with Stack Traces
- Update `cli/cmddiff/dispatchdiff.go` & `cli/cmddiff/diff.go`:
  - When invoked with 1 argument (`gitmap diff <path>` or `gitmap diff`):
    - If `path` is inside a Git repo, launch working tree comparison via `cmdgit.FindCompareTool()` (Beyond Compare / VS Code / Meld / WinMerge).
    - If compare tool fails or is not found, render colored terminal working tree diff with full AppError stack trace.
- When `/api/log/diff` encounters an error, include `stack_trace` and detailed error cause in the JSON envelope:
  ```json
  {
    "isSuccess": false,
    "error": "Failed to launch Beyond Compare",
    "stackTrace": "at cmdgit.FindCompareTool (compare_launcher.go:34)..."
  }
  ```
