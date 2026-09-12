# 114-macro-live-execution-copy-explorer-browser.md: Macro Live Execution & Edit, Memory/Clipboard Copy-Paste, Explorer & Browser URL Openers

**Status: completed**

## 1. Executive Summary

Empowers GitMap macro creation and editing with real-time terminal feedback, cross-platform clipboard/memory persistence, OS file manager integration, and default/Chrome browser URL launching:
1. **Macro Live Execution by Default**: During interactive macro creation (`gitmap macro add/create/new`), every step command entered executes live with streamed terminal output (stdout/stderr) so users visually inspect what happens in real time.
2. **Interactive Macro Edit Mode**: Introduces `gitmap macro edit <name>` (with aliases `edit`, `modify`, `macro-edit`, `macro-modify`) allowing users to view existing steps, remove, replace, insert, or append steps with live command execution and terminal visibility.
3. **Memory & Clipboard Copy-Paste System**: Introduces `gitmap copy <item|text|--file <path>>` (aliases: `copy`, `cp-mem`, `copy-to-memory`) and `gitmap paste [--file <path>]` (aliases: `paste`, `paste-mem`). Persists copied items to both the OS clipboard and an internal repository memory buffer (`.gitmap/memory/clipboard.txt`) so items can be pasted into files, clipboard, or downstream macro steps.
4. **Cross-Platform Explorer Opener**: Introduces `gitmap explorer [path]` (aliases: `open-explorer`, `folder`, `open-folder`, `browse-folder`) that opens the native graphical file manager on Windows (`explorer.exe`), macOS (`open`), and Linux (`xdg-open`).
5. **Cross-Platform Browser & Chrome URL Opener**: Introduces `gitmap open-url <url> [--chrome]` (aliases: `browse`, `open-url`, `browse-url`, `open-browser`) allowing opening any URL in either the OS default browser or specifically in Google Chrome.

## 2. Invariants & Guidelines Verified

1. **Strict Relative Git Paths**: All paths in markdown files, artifacts, and documentation are strictly relative to repository root.
2. **Strict Function Sizing**: Every Go function strictly <= 15 lines, with mandatory blank line before every return statement.
3. **Affirmative Booleans**: All boolean variables and parameters use affirmative prefixes (`is*`, `has*`).
4. **Zero Swallowed Errors**: AppError wrapping with context for all failing operations.
5. **Cache Tracking**: Recorded 14 modified files in `.lovable/temp/recent-file-changes.json` under atomic lock via `python 03-ai-scripts/33-test-inventory-generator.py --record`.
6. **Constraint Enforcement**: No auto-releases, no unit test execution, no `06-cicd-local-runner.py`.

## 3. Subtasks Consolidated

- **Subtask 01**: Default live execution in interactive macro creation and interactive editor in `gitmap/cmd/macro_edit.go`.
- **Subtask 02**: `gitmap copy` and `gitmap paste` commands with file/text memory buffering and OS clipboard integration in `gitmap/cmd/copy_paste.go`.
- **Subtask 03**: Cross-platform desktop file manager opener in `gitmap/cmd/explorer.go`.
- **Subtask 04**: Cross-platform default browser and Google Chrome URL opener in `gitmap/cmd/browse.go`.
- **Subtask 05**: Command registration in `rootdata.go` and `rootutility.go`, help documentation, code formatting, and verification.
