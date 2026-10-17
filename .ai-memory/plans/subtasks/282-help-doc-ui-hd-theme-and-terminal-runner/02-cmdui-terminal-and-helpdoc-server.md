# Subtask 02: Backend Go UI Server, Terminal Top-Level Tab & HelpDoc Integration

**Task Slug:** `282-help-doc-ui-hd-theme-and-terminal-runner`  
**Subtask ID:** `subtask-02`  
**Owner:** Worker 02  
**Status:** `COMPLETED`

---

## Scope of Work

1. **`cli/cmdui/ui_assets_js_nodes.go`**:
   - Add `{ id: 'terminal', label: 'Terminal', icon: 'scroll-text' }` to `NAV_ITEMS`.
   - Ensure `showTab('terminal')` activates the dedicated terminal tab, fixing direct `/terminal` routing.

2. **`cli/cmdui/ui_assets_markup_ops.go`**:
   - Promote terminal from settings sub-pane into top-level `<div id="tab-terminal" class="content-area">`.
   - Add prompt search input and interactive results container in `<div id="tab-prompts">`.

3. **`cli/cmdui/ui_assets_markup_misc.go`**:
   - Add dynamic embedded helpdoc search and topic viewer in `<div id="tab-help">`.

4. **`cli/cmdui/ui_assets_theme.go`**:
   - Enhance contrast tokens, clean card borders, and theme definitions for `--terminal-bg` and `--terminal-text`.

5. **`cli/cmdui/ui_server.go`**:
   - Mount `/api/command/exec` alias pointing to `handleAPITerminalExec`.
   - Mount `/api/help/topics` and `/api/help/doc` calling `cli/helpdoc` engine.
   - Optimize local terminal execution with faster command shell (`cmd.exe /c` on Windows).

6. **`cli/cmd/hd_terminal.go`**:
   - Mount `/api/terminal/exec` alias pointing to `handleCommandExec`.
