# Plan 21: Terminal Help Layout, LLM Specification URL & SSH Subcommands

## Overview

Comprehensive remediation of terminal help formatting (eliminating excessive column gaps and aligning descriptions close to commands), upgrading the LLM guideline command (`gitmap llm` and `gitmap llm --url`), and adding all missing SSH subcommands (`ssh join`, `ssh login`, `ssh alias`, `ssh exec`) to `gitmap ssh` and main help.

## Root Cause Analysis

1. **Excessive Column Gap & Mid-Screen Alignment:**
   - `maxHelpCmdLen` was computed globally across all 50+ command groups combined and capped at 38, forcing short command sections (like `History & Stats`, `Data, Profiles & Bookmarks`, `Author Amendment`, where commands are only 10–20 chars) to indent descriptions 30+ spaces away at column 42.
   - Certain constants had single-space separators or trailing spaces that prevented proper splitting.
   - **Solution:** Render help rows with compact, per-group dynamic alignment (or optimal 26-column base with max 30 cap), printing long commands on line 1 and indented description on line 2.

2. **LLM Guidelines & Public URL:**
   - `gitmap llm` provided minimal output and needed two clear modes: full rich markdown instructions and a clean public GitHub raw URL (`https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/llm.md`) for AI agents.
   - **Solution:** Enrich `gitmap/cmd/llm/llm.go` with full specification, command alternatives, and ensure `--url` flag returns the public raw link.

3. **Missing SSH Subcommands:**
   - `gitmap ssh` help output was missing `ssh join` (cluster SSH nodes), `ssh login` / `ssh login-install`, `ssh alias`, and `ssh exec`.
   - **Solution:** Update `constants.MsgSSHAvailableCommands`, `constants.HelpSSH`, and `gitmap/cmd/ssh.go` dispatcher.

## Subtasks Breakdown

- **Subtask 21.01:** Refactor `gitmap/cmd/rootusage.go` for compact per-group command width (26–30 max), eliminating mid-screen gaps and aligning bookmark/history/data sections.
- **Subtask 21.02:** Upgrade `gitmap/cmd/llm/llm.go` with rich AI instructions and public GitHub MD URL (`--url` flag).
- **Subtask 21.03:** Add missing SSH subcommands (`join`, `login`, `alias`, `exec`) to `constants_ssh.go`, `ssh.go`, and main help.
- **Subtask 21.04:** Verification, unit tests, and terminal output validation.

## Acceptance Criteria

- [ ] Terminal help column gap is compact (descriptions start at col 28–30 instead of col 42+).
- [ ] Bookmark, History, Data, and Version History rows are aligned with clean padding.
- [ ] `gitmap llm` outputs full specification + public URL; `gitmap llm --url` prints raw URL.
- [ ] `gitmap ssh` lists all subcommands including `create`, `list`, `status`, `copy`, `cat`, `delete`, `config`, `join`, `login`, `alias`, `exec`.
- [ ] `go test ./...` in `gitmap/` passes cleanly.
