# Subtask 01: PWD Header Display & Toggle in Interactive Macro Builder

**Parent Plan:** [83-interactive-macro-builder-pwd-ls-search.md](../../completed/83-interactive-macro-builder-pwd-ls-search.md)  
**Status:** complete  
**Target:** `gitmap/cmd/macro_add_interactive.go`, `gitmap/uipref/`

---

## Objectives

1. In `promptInteractiveMacroSteps`, before printing `Step %d> `, retrieve and format current working directory via `os.Getwd()`.
2. Format as a clean dim/cyan header:
   ```text
     [PWD: /path/to/cwd]
     Step 1> 
   ```
3. Support dynamic toggling via `pwd on` and `pwd off` (or `:pwd on` / `:pwd off`) in the loop.
4. Support persistent storage of preference in `gitmap/uipref/` with getter/setter `GetMacroShowPwd()` / `SetMacroShowPwd(bool)`.
5. Support command line flag `--pwd` / `--no-pwd` on `gitmap macro add`.
