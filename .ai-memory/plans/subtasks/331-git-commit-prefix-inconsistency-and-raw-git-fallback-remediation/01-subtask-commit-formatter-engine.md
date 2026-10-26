# Subtask 331.1 — Standardize Commit Formatter Engine for Lowercase Modules & Conventional Syntax

> **Parent Plan:** `.ai-memory/plans/pending/331-git-commit-prefix-inconsistency-and-raw-git-fallback-remediation.md`  
> **Target Package:** `cli/cmdcommit/`  
> **Status:** Completed  

---

## 1. Objectives & Resolution
1. Modified `cli/cmdcommit/commit_formatter.go`:
   - Replaced uppercase module return in `resolveHyphenModuleAndItem` with `strings.ToLower(candidateMod)`.
   - All modules (`COMMIT`, `DOCS`, `BROWSE`, `UI`, `INSTALL`, `CONTEXTMENU`, `AUTH`, `SQLITE`) are strictly normalized to lowercase (`commit`, `docs`, `browse`, `ui`, `install`, `contextmenu`, `auth`, `sqlite`).
   - Added `VERSION` to `isCommitTypeToken` to strip redundant `VERSION - ` tokens.
   - Replaced `resolveUppercaseType` with `resolveConventionalType` returning standard lowercase conventional types (`feat`, `fix`, `chore`, `refactor`, `release`, `test`, `docs`, `perf`).
   - Conventional commit formatter strictly formats with lowercase type and lowercase module (`feat(module): summary`).
2. Updated all unit tests in `cli/cmdcommit/commit_formatter_test.go` to assert lowercase modules and lowercase conventional commits.
3. Verified 100% pass with `go test -C cli ./cmdcommit` (3.291s).
4. Recompiled and deployed `gitmap` binary to local application directory.
