# CI/CD Issue 46: Pipeline Smoke, Linter Baseline Diff, Nested If, Test Truncate, and Generate Drift

- JOB: Installer Smoke (release), Lint Baseline Diff, Go Vet & Test, Generate Drift
- TYPE: FAIL
- DETECTED: 2026-09-16
- STATUS: resolved
- PIPELINE RUN: #35072529679 / #35072522847

## 1. Symptoms
1. *_installer Smoke Failure (`task --help`)** :
   `open tasks.md: file does not exist` - apperror E1071 in helptext.PrintWithMode.
2. **Linter Baseline Diff Failures**:
   - macro/recurse.go:15:6: `contextKey` is unused.
   - cmd/desktopsync_ops.go:15:6: `syncResult` and consts unused.
   - cmd/replace_classify.go:10:6: `replaceMode` and consts unused.
   - cmd/task.go:12:6: func `runTask` is unused.
3. **Nested If Linter Violation**:
   `cli/pipelinedb/pipeline_split_db.go:90: Nested 'if' detected (depth 2)` in `migrateOrFallbackPipelineDb`.
4. **TestTruncateMiddle Failure**:
   `agy_ls_table_test.go:25: truncateMiddle("ai-empathy-prompt-tuner", 15) = "ai-emp...-tuner"; want "ai-emp...t-tuner"`.
5. **TestFormatRelativeDbPath Windows Backslash Failure**:
   `pipeline_history_test.go:80: expected forward slashes in relative db path, got \...`.
6. **Generate Drift Check Failure**:
   `Generated files are out of sync with constants. Run 'cd gitmap && go generate ./...' locally && commit the result.`

## 2. Root Cause Analysis
1. *Help Topic Lookup & Alias Gap*:
   `runTasks` invoked `checkHelp("tasks", args)`, but only `task.md` was embedded in helptext/ and `helpAliases` lacked `tasks -to- task` alias.
2. *Unused Types & Constants*:
   Refactorings had migrated type usages to PascalCase forms while leaving dead unexported aliases and constants.
3. **Nested Control Flow & File Size**:
   `migrateOrFallbackPipelineDb` nested `os.Rename` inside `if isFileExisting`. The file was also 125 lines, exceeding the 100-line guideline.
4. **String Truncation Math in Test**:
   `truncateMiddle("ai-empathy-prompt-tuner", 15)` computes 6 prefix + 3 dots + 6 suffix = 15 characters (`ai-emp...-tuner`). The test erroneously asserted 16 characters (`ai-emp...t-tuner`) which violated `got <= 15`.
5. **Platform Path Separator Normalization**:
   `FormatRelativeDbPath` returned `filepath.Clean(fullPath)` which on Windows produced backslashes.
6. **Out-of-Sync Generated Completion Catalog**:
   `CmdTasks` was added without running `go generate ./...` to sync `cli/completion/allcommands_generated.go`.

## 3. Resolution & Fixes Applied
1. Added `tasks` and `tk` to `helpAliases` in `cli/helptext/print.go`, registered `task` in `catalog.go`, created `cli/helptext/tasks.md`, updated `runTasks` to ask for `task`, and removed unused `runTask` from `cli/cmd/task.go`.
2. Removed unused `contextKey`, syncResult/syncAdded, and replaceMode/replaceModeUnknown declarations.
3. Extracted fs functions and flattened `pipeline_split_db_fs.go` (43 lines), reducing `pipeline_split_db.go` to 93 lines (guideline compliant).
4. Corrected want assertion in `agy_ls_table_test.go` to `ai-emp...-tuner`.
5. Normalized return path in `cli/cmdpipeline/pipeline_sync_cache.go` using `filepath.ToSlash(filepath.Clean(fullPath))`.
6. Re-generated `allcommands_generated.go` via `go generate ./...`.

## 4. Prevention & Learnings
- Run `go generate ./...` after adding or updating command constants.
- Pair singular and plural commands in `catalog.go` and `helpAliases`.
- Always normalize filesystem paths for terminal display with `filepath.ToSlash`.
- Keep all new or refactored Files strictly <= 100 lines, with zero nested ifs.
