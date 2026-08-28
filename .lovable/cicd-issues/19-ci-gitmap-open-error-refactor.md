# 19-ci-gitmap-open-error-refactor.md

## RCA: Python `UnicodeDecodeError` in `bump_versions.py`
The build failed when reading `.lovable/cicd-issues/index.md` or `readme.md` during version bumping because Python on Windows defaulted to `cp1252` instead of UTF-8 when encountering emojis or unicode characters in markdown files.
**Resolution**: Added `encoding="utf-8"` to all `open()` calls in `bump_versions.py`.

## Massive Error Refactoring
The codebase heavily relied on `os.Exit(1)` inside command handlers (`run[A-Z]`), making error tracking, testing, and unified database logging impossible. 
**Resolution**: 
- A custom Go AST rewriting script was created and run across the entire `gitmap/cmd` package (80+ commands).
- Refactored `dispatchEntry` to use `handler func() error`.
- Refactored all `dispatchXxx` mappings to return `(bool, error)`.
- Replaced empty returns in `runXxx` functions with `return nil`.
- Bubbled errors all the way to `cmd.Run()` which now calls `cliexit.Reportf` to format and print errors uniformly.

## `gitmap open`
Implemented cross-platform file, folder, and URL opening via the host OS using `rundll32`, `open`, and `xdg-open`.

## Documentation
Added cross-platform search pipeline equivalents (like `Get-ChildItem -Path ... | Select-String`) to `gitmap search` and `gitmap repo-search` helptext files to ensure AI LLMs have training references.
