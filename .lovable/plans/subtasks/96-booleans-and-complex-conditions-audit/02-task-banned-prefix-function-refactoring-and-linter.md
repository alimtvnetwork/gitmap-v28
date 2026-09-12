# Subtask 02: Banned Function Prefix Refactoring & Linter Hardening

**Parent Plan:** `.lovable/plans/pending/96-booleans-and-complex-conditions-audit.md`  
**Status:** In Progress  
**Target Files:**
- `gitmap/cliexit/handle.go`
- `gitmap/clonefrom/execute.go` & `execute_dest.go`
- `gitmap/cmd/` (`agy_scan.go`, `audit.go`, `clone.go`, `clonenext.go`, `clonenextbatchdispatch.go`, `inject.go`, `inject_idempotency.go`, `installer_ls.go`, `install_add.go`, `install_prompt.go`, `installscripts.go`, `prune.go`, `pull.go`, `push.go`, `release.go`, `releaseautobump.go`, `root.go`, `root_url_shortcut_test.go`, `scanresolve.go`, `schedule_export.go`, `schemaregistry_assert_test.go`, `schemaregistry_contract_test.go`, `selfuninstall.go`, `seowritegit.go`, `seowriteloop.go`, `updaterepo.go`)
- `gitmap/committransfer/pr_engine.go` & `replay.go`
- `gitmap/detector/detector.go` & `csharpparser.go`
- `gitmap/diff/tree.go`
- `gitmap/release/selfrelease_resolve.go`
- `gitmap/vscodepm/optimize.go`
- `linter-scripts/check-boolean-guidelines.py`
- `03-ai-scripts/01-index.md`

---

## 1. Objectives & Grounded Rules

1. **Rename Banned Prefix Functions (`can*` / `should*` returning `bool`) & Callers:**
   - `shouldPrintStackTrace` -> `isStackTraceEnabled`
   - `shouldSkip` -> `isSkippable`
   - `shouldSkipScanDir` -> `isSkippableScanDir`
   - `shouldAuditCommand` -> `isAuditableCommand`
   - `shouldUseMultiClone` -> `isMultiCloneEnabled`
   - `shouldRunBatch` -> `isBatchRunnable`
   - `shouldRunDesktop` -> `isDesktopRunnable`
   - `shouldRunVSCode` -> `isVSCodeRunnable`
   - `shouldSkipInstallerRow` -> `isSkippableInstallerRow`
   - `shouldPromptInteractively` -> `isInteractivePromptNeeded`
   - `shouldProceedInstall` -> `isInstallApproved`
   - `shouldProceedWithPrune` -> `isPruneConfirmed`
   - `shouldPullCWD` -> `isPullCWDEnabled`
   - `shouldPushCWD` -> `isPushCWDEnabled`
   - `shouldAutoBumpMinor` -> `isAutoBumpEligible`
   - `shouldRewriteToClone` -> `isCloneRewriteRequired`
   - `shouldAnnounceResolve` -> `isAnnounceResolveNeeded`
   - `shouldSkipSchedule` -> `isSkippableSchedule`
   - `shouldUpdateSchema` -> `isSchemaUpdateNeeded`
   - `shouldHandoffSelfUninstall` -> `isSelfUninstallHandoffNeeded`
   - `shouldStop` -> `isStopRequested`
   - `canPromptForRepoPath` -> `isRepoPathPromptAllowed`
   - `shouldCreatePR` -> `isPRCreationEnabled`
   - `shouldSkipPath` -> `isSkippablePath`
   - `shouldExcludeDir` -> `isExcludedDir`
   - `shouldIgnore` -> `isIgnoredPath`
   - `canPromptForPath` -> `isPathPromptAllowed`
   - `shouldWrite` -> `isWriteEnabled`

2. **Linter Hardening in `linter-scripts/check-boolean-guidelines.py`:**
   - Update `INVERTED_SUCCESS_REGEX` to check `!\s*(?:[a-zA-Z0-9_$.->]+\.)?[iI]sSuccess\b` (case-insensitive for Go).
   - Add check for banned function prefixes returning `bool` (`\bfunc\s+(?:(?:\([a-zA-Z0-9_*]+\)\s+)?)(can[A-Z]\w*|should[A-Z]\w*)\s*\([^)]*\)\s*bool\b`).
   - Ensure linter exits with code 0 only when all files are compliant.

3. **Coding Standards:**
   - Functions <= 8 lines preferred (hard cap 15 lines).
   - UNIX LF line endings and UTF-8 without BOM.
   - Blank lines before `return` and after `}` blocks.
   - Verify with `go test -v -short ./...` and `python linter-scripts/check-boolean-guidelines.py`.
