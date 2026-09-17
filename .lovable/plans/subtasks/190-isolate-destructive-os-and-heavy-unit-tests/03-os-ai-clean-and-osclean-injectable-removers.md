# Subtask 03: OS AI Clean & OSClean Injectable Removers

## Parent Plan
- Parent Plan: `.lovable/plans/pending/190-isolate-destructive-os-and-heavy-unit-tests.md`

## Objectives
1. Refactor `cli/cmdos/os_ai_clean.go`:
   - Introduce `type FileRemover func(filePath string) (int, int64)` with `var defaultFileRemover FileRemover = removeSingleFileSafely`.
   - Update `purgeCategoryFiles` to invoke `defaultFileRemover`.
2. Create `cli/cmdos/os_ai_clean_test.go`:
   - Test `RunOSAICleanCLI` with `--dry-run`, `--json`, and `--yes` using mock `defaultFileRemover` and `promptAICleanConfirmFn`.
   - Verify that NO real user files are deleted from `~/.gemini/antigravity/brain` or system directories.
3. Refactor `cli/osclean/clean.go`:
   - Introduce `type TempDirResolver func() []string` with `var defaultTempDirResolver TempDirResolver = resolveTempDirectories`.
   - Update `CleanTempDirectories` to use `defaultTempDirResolver`.
4. Create `cli/osclean/clean_test.go`:
   - Test `CleanTempDirectories` with a mock resolver returning `t.TempDir()`.
   - Verify file removal count, directory removal count, dry-run mode, and error handling safely.

## Coding Standards
- Functions <= 15 lines (target <= 8 lines).
- Affirmative booleans only (`is*`, `has*`).
- Strict Unix LF line endings.
- AppError error returns where applicable.
