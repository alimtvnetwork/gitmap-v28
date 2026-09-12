# 94-error-management-and-apperror-architecture

## Overview
Comprehensive audit and surgical refactoring of error management and CLI exit architecture across GitMap adhering to `spec/03-error-manage/`, `spec/02-coding-guidelines/`, and `.lovable/coding-guidelines.md`. This plan eliminated linter-evasion bypasses (`exitWith = os.Exit`), disguised bare exits (`cliexit.HandleError(nil, <magic-code>)`), naked dropped `apperror` expressions, swallowed database errors, and dual-handling in leaf functions. It implemented strongly-typed `ExitCodeType` enums, specialized exit helpers (`HandleValidationError`, `HandleUsageError`, `HandleGeneralError`, `HandleSuccess`), `*apperror.AppError` returns, and the Universal Response Envelope.

## Custom Task-Specific Rules Enforced
1. **Error Return Sovereignty (No Dual-Handling):** Leaf and service functions returning `error` or `*apperror.AppError` construct and return the error. No leaf functions call exit handlers or panics and return `nil`.
2. **No Magic Literal Exit Codes:** Raw integers passed to exit handlers replaced with `cliexit.ExitCodeType` or specialized helpers (`HandleValidationError`, `HandleUsageError`, `HandleNotFound`, `HandleSuccess`).
3. **No Disguised Exits or Linter Evasions:** `exitWith = os.Exit` permanently deleted across the entire repository. Passing `nil` to error handlers replaced with proper `*apperror.AppError` returns or `cliexit.Exit`.
4. **No Naked Error Statements:** Every call to `apperror.NewSimple` or `apperror.WrapSimple` is returned or handled; zero discarded error construction statements.
5. **No Swallowed DB Errors:** Every database `Scan`, `Exec`, `Query`, and row iteration inspects errors and wraps them in `*apperror.AppError`.

## Subtasks Executed and Verified
- `01-task-linter-and-cliexit-dispatcher-hardening.md`: Upgraded `linter-scripts/check-error-management.py` to AST-detect bare panics/exits, `exitWith`, naked `apperror` expressions, and unhandled DB errors; permanently deleted `exitWith` in `cmd/pendingtaskhelper.go`; hardened `cliexit/`.
- `02-task-cmd-naked-apperrors-and-reset-commands.md`: Fixed all 54 naked dropped `apperror` statements across `cmd/group*.go`, `cmd/env*.go`, `cmd/gomod*.go`, `cmd/historyreset.go`, `cmd/dbreset.go`, `cmd/dashboard.go`, `cmd/desktopsync.go`, `cmd/downloaderconfig.go`.
- `03-task-cmd-leaf-sovereignty-and-exit-migration.md`: Refactored `cmd/asops.go`, `cmd/as.go`, `cmd/cdops.go`, `cmd/amendexec.go`, `cmd/aliasops.go`, `cmd/bookmark*.go`, `cmd/status.go`, `cmd/clonenext.go`, `cmd/fixauth.go`, `cmd/whoami.go`, `cmd/sshbind.go`, `cmd/aliassuggest.go`, `cmd/installscripts.go`, `cmd/selfinstall.go` to return errors directly and use specialized exit helpers.
- `04-task-db-and-store-error-wrapping-and-swallows.md`: Converted `gitmap/db/` and `gitmap/searcher/` to return `*apperror.AppError`, eliminated swallowed errors in `finder.go`, `db_search.go`, `workdir_list.go`, `pipeline.go`, `scheduler.go`, `chromeprofile_export_all.go`, `chromeprofile_tokens.go`, `repo_db_ops.go`, `rm.go`, `schedule_export.go`.
- `05-task-frontend-api-envelopes-and-ci-gates.md`: Formatted `gitmap/cmd/hd_terminal.go` handlers with Universal Response Envelopes, aligned `src/types/result.ts` with `src/lib/queryWrapper.ts`, decomposed `Settings.tsx` into sub-80-line components with `queryWrapperSync`, and verified `check-error-management.py` with 100% green pass (0 violations across 2696 source files).
