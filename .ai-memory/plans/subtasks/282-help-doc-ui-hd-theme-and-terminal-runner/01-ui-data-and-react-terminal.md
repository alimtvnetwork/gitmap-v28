# Subtask 01: Frontend React UI, Commands Dataset & Live Terminal Runner

**Task Slug:** `282-help-doc-ui-hd-theme-and-terminal-runner`  
**Subtask ID:** `subtask-01`  
**Owner:** Worker 01  
**Status:** `COMPLETED`

---

## Scope of Work

1. **`src/data/commands.ts`**:
   - Add full command definition for `reclone` (alias `rec`, `cnow`).
   - Include description, category (`"repo"`), flags (`--manifest`, `--scan-root`, `--mode`, `--execute`, `-y`, `--no-summary`, `--quiet`), and single-repo re-cloning example.

2. **`src/data/troubleshootingIssues.ts`**:
   - Add troubleshooting issue card for `git-reclone-repair`:
     - Title: "Repository corruption, failed clone, or stale branch tracking"
     - Resolution steps using `gitmap reclone <url-or-path> -y`.

3. **`src/pages/HelpIndex.tsx`**:
   - Surface the new command cards:
     - `gitmap resolve` and `gitmap conflicts` in Merge & Conflict Resolution section.
     - `gitmap split large files`, `gitmap lfs convert`, `gitmap repo flatten`, and `gitmap reclone` in Large Files & Storage section.

4. **`src/components/terminal/TerminalDrawer.tsx`**:
   - Wire live execution calling `POST /api/command/exec` with fallback to `POST /api/terminal/exec`.
   - Implement command history navigation with ArrowUp and ArrowDown keys.
   - Respect user's active theme (`isDark` from `useTheme()`) and terminal color theme (`dark`, `ubuntu`, `win`).
   - Add clear exit code badges, duration timer, and output display.

5. **`src/pages/AGYPrompts.tsx`**:
   - Add interactive prompt search input and category filter chips.
   - Add one-click copy and terminal execution links.
