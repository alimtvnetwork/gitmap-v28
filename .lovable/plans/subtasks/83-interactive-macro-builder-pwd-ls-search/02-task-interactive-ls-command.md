# Subtask 02: In-Builder LS and DIR File Inspection

**Parent Plan:** [83-interactive-macro-builder-pwd-ls-search.md](../../completed/83-interactive-macro-builder-pwd-ls-search.md)  
**Status:** complete  
**Target:** `gitmap/cmd/macro_add_helpers.go`, `gitmap/cmd/macro_add_interactive.go`

---

## Objectives

1. Intercept `ls` or `dir` input inside `promptInteractiveMacroSteps`.
2. Inspect directory entries in current working directory:
   - Separate directories (blue/bold) and files (white/dim).
   - Display size (formatted in KB/MB) and modification date.
   - Display summary total files and directories count.
3. Allow user to either inspect only, or record `ls` as a step by answering prompt or typing `add ls`.
4. Ensure no exit, abort, or corruption of previously entered macro steps occurs.
