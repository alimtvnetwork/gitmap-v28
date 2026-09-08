# Subtask 04: Unit Testing, Linting & CI/CD Verification

**Parent Plan:** [83-interactive-macro-builder-pwd-ls-search.md](../../completed/83-interactive-macro-builder-pwd-ls-search.md)  
**Status:** complete  
**Target:** `gitmap/cmd/macro_add_helpers_test.go`, `gitmap/cmd/macro_add_test.go`

---

## Objectives

1. Write unit tests for:
   - PWD formatting and toggle state (`TestMacroPromptPwdToggle`).
   - Directory listing output formatting (`TestMacroInteractiveLs`).
   - File search, find, and replacement algorithms (`TestMacroInteractiveFindSearchReplace`).
2. Run quality linters:
   - `python linter-scripts/check-nested-ifs.py`
   - `python linter-scripts/check-boolean-guidelines.py`
   - `python linter-scripts/check-relative-paths.py`
   - `python linter-scripts/check-error-management.py`
3. Execute local CI/CD pipeline:
   - `python 03-ai-scripts/06-cicd-local-runner.py`
