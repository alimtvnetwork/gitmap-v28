# Plan: Help Dashboard (HD) UI, Theme System, Terminal Runner & Help Documentation

**Task Slug:** `282-help-doc-ui-hd-theme-and-terminal-runner`  
**Parent Task ID:** `task-282`  
**Execution Budget:** N = 300, A = 2, H = 2  
**Spec Reference:** `02-spec/21-app/282-help-doc-ui-hd-theme-and-terminal-runner/spec.md`  
**Status:** `COMPLETED`

---

## User Request (Verbatim)

```text
# High Priority Instruction

A lot of new commands you have added. Have you added the help text into the UI and terminal confirm it? And also make sure that the help text UI using the HTML, CSS that has better themes now to display and find prompts. And also test out from the terminal section of the UI, that does not look very good. We cannot run commands and see the terminal view. That needs to be fixed into the HD section, help doc section from the `gitmap`. Can you please find the workflow and fix everything?

# Actionable Items Must Follow Non-Negotiable

0. Gitmap do `gitmap pull <repo path> or run into the repo path` Make sure we have the latest and conflict resolve before starting.
1. Write spec under 02-spec/21-app/<slug>/ and enqueue plan task in .ai-memory/plans/<slug>.md (subtasks in .ai-memory/plans/subtasks/<slug>/) first
2. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
3. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
4. Confirm that help text for new commands is added to the UI and terminal.
5. Ensure the help text UI uses HTML and CSS with improved themes for better display and prompt finding.
6. Test the terminal section of the UI to ensure commands can be run and the terminal view is functional.
7. Fix issues in the HD section, help doc section from the `gitmap`.
8. Find and fix the workflow issues related to the above tasks.

## Must follow and spawn agent using

@[.agents/skills/execute-parent-task-with-n-steps-v6]

## Additional Instructions

- /plan first before doing the work to reduce the credits.
- /learn from @[.agents/skills/gitmap] skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well.
```

---

## Subtask Decomposition (A = 2 Disjoint File Boxes)

### Wave 1: Parallel Core Implementation

- **Subtask 01: Frontend React UI, Commands Dataset & Live Terminal Runner**
  - **Owner:** Worker 01
  - **Plan File:** `.ai-memory/plans/subtasks/282-help-doc-ui-hd-theme-and-terminal-runner/01-ui-data-and-react-terminal.md`
  - **Owned Files:**
    - `src/data/commands.ts`
    - `src/data/troubleshootingIssues.ts`
    - `src/pages/HelpIndex.tsx`
    - `src/components/terminal/TerminalDrawer.tsx`
    - `src/pages/AGYPrompts.tsx`
  - **Deliverables:**
    1. Register `reclone` (`rec`) in `commands.ts` and `troubleshootingIssues.ts`.
    2. Add new command cards (`resolve`, `conflicts`, `split-large-files`, `lfs-convert`, `repo-flatten`, `reclone`) to `HelpIndex.tsx`.
    3. Upgrade `TerminalDrawer.tsx` from mock strings to live API execution (`POST /api/command/exec` with fallback to `/api/terminal/exec`), arrow history navigation, status badges, and theme support (`dark`, `ubuntu`, `win`).
    4. Implement prompt search and filtering on `AGYPrompts.tsx`.

- **Subtask 02: Backend Go UI Server, Terminal Top-Level Tab & HelpDoc Integration**
  - **Owner:** Worker 02
  - **Plan File:** `.ai-memory/plans/subtasks/282-help-doc-ui-hd-theme-and-terminal-runner/02-cmdui-terminal-and-helpdoc-server.md`
  - **Owned Files:**
    - `cli/cmdui/ui_assets_js_nodes.go`
    - `cli/cmdui/ui_assets_markup_ops.go`
    - `cli/cmdui/ui_assets_markup_misc.go`
    - `cli/cmdui/ui_assets_theme.go`
    - `cli/cmdui/ui_server.go`
    - `cli/cmd/hd_terminal.go`
  - **Deliverables:**
    1. Promote terminal in `cli/cmdui/` to top-level tab in nav and content area; fix direct routing for `/terminal`.
    2. Add prompt search and active instances list in `cli/cmdui/` `#tab-prompts`.
    3. Add dynamic embedded helpdoc search and topic viewer in `#tab-help`.
    4. Improve dark/light theme CSS tokens and terminal styling.
    5. Add dual API endpoint aliases (`/api/command/exec` $\leftrightarrow$ `/api/terminal/exec`) across both servers.

---

## Verification & Completion Gates

1. **FastGate Check:** `python 03-ai-scripts/50-fastgate.py` must pass with 0 errors.
2. **Terminal Help Text Verification:** Confirm `gitmap help --filter`, `gitmap help <cmd>`, `gitmap howto`, and `gitmap nc`.
3. **Web UI Verification:** Confirm command search, prompt finder, and live terminal execution.
4. **GitMap Atomic Commit:** Single commit via `gitmap cpf "hd-ui - improve help dashboard themes prompt finder and terminal view"`.
