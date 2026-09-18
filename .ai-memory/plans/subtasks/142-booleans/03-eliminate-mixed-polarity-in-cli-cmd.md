# Subtask 03: Eliminate Mixed Polarity in cli/cmd Subsystem

Parent Plan: [142-boolean-principles-negatives-and-complex-conditions.md](../../pending/142-boolean-principles-negatives-and-complex-conditions.md)

## Goals
1. Refactor `cli/cmd/backup.go`: Flatten `err == nil && !info.IsDir()` to guard clauses.
2. Refactor `cli/cmd/backup_cloud_ops.go`: Decompose `len(args) > 0 && !strings.HasPrefix(args[0], "-")`.
3. Refactor `cli/cmd/cdops.go`: Restructure `len(dflt) > 0 && !pick` into clean linear guard flow.
4. Refactor `cli/cmd/changelog.go`: Extract affirmative `isDefaultVersion` and compose with `!latest`.
5. Refactor `cli/cmd/changeloggen.go`: Normalize `err != nil && !os.IsNotExist(err)` to `if errors.Is(err, os.ErrNotExist) { err = nil }`.
6. Refactor `cli/cmd/clustersubcmd.go`: Extract `hasTokensMissing` and compose positive `isGitMissing` and `isProjMissing`; replace `!isStrippedEmpty` with `hasStrippedContent`.
7. Refactor `cli/cmd/commit_push.go`: Restructure `targetSha == "" && !strings.HasPrefix(arg, "-")` with guard clause `if strings.HasPrefix(...) continue`.
8. Refactor `cli/cmd/create_cmd.go`: Extract affirmative `isHeadlessError` composite.
9. Refactor `cli/cmd/dopendingretry.go`: Extract `isDirMissing := !pathExists(workDir)` and compose positive `isWorkDirNotFound`.
10. Refactor `cli/cmd/find_duplicates.go`: Restructure `len(args) > 0 && !strings.HasPrefix(args[0], "-")` into guard clauses.
11. Refactor `cli/cmd/find_files_match.go`: Extract positive `isExtDisallowed` and compose `isFilterRejected`.
12. Refactor `cli/cmd/fixauth.go`: Decompose `keyExists && !force` and `keyExists && !assumeYes && !confirmOverwrite(...)`.
13. Refactor `cli/cmd/historyrewrite_push.go`: Decompose `!opts.yes && !confirmHistoryPush(...)`.
14. Refactor `cli/cmd/ipchange_cmd.go`: Extract `isPingFailed` and compose positive `isRollbackNeeded`.

## Acceptance Criteria
- [x] All 14 mixed polarity conditions across target files eliminated.
- [x] Conditional nesting depth remains <= 1 across all modified blocks.
- [x] `go vet ./...` in `cli/` passes with 0 errors.
