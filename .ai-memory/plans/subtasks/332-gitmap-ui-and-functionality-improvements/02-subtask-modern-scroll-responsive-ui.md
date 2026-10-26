# Subtask 332.2 — Modern Scrollable Responsive Show Log UI & Exit Integration

> **Parent Plan:** `.ai-memory/plans/pending/332-gitmap-ui-and-functionality-improvements.md`  
> **Target Packages:** `cli/cmdlogviewer/`  
> **Status:** Pending  

---

## 1. Objectives
1. In `cli/cmdlogviewer/ui_assets_html.go`:
   - Modernize visual design: modern scrollbar rules (`::-webkit-scrollbar`), sleek borders, and Catppuccin palette.
   - Make layout responsive: wrap flex containers with `@media (max-width: 900px)`.
   - Add "Exit" button in toolbar (`#btn-exit`).
   - Add HTML container for custom right-click context menu.
2. In `cli/cmdlogviewer/server.go`:
   - Add route `/api/log/exit` to gracefully close the HTTP listener and terminate the server process.
3. In `cli/cmdlogviewer/ui_assets_script.go`:
   - Wire exit button: call `/api/log/exit` and run `window.close()`.
   - Implement custom right-click context menu for commits and files:
     - Copy commit hash / subject.
     - Copy file path.
     - Revert commit / file changes.
   - Double-click on file row launches Beyond Compare or external diff tool with clear toast message and stack trace on error.
4. Add unit tests for exit endpoint and UI assets.
