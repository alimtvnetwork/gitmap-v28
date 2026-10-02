# RCA-60: RootCore Commit Dispatch Shadowing and Stale Binary Macro Execution Failure

## 1. Symptom

1. **Macro Execution Failure (`rm test`):**
   When executing macro `alim1` (`gitmap alim1`), step 1 failed terminatingly on Windows PowerShell with `exit status 1` (`ItemNotFoundException` from `Remove-Item` cmdlet) because the target directory `C:\Users\Alim\test` did not exist.
   Even though `cli/macro/safe_rm.go` was previously designed, the installed runtime binary at `%LOCALAPPDATA%\gitmap-cli\gitmap.exe` was out-of-date and had not been recompiled/installed with the platform safe removal adapter.
2. **Commit Command Routing Collision:**
   When adding the new `gitmap commit` / `cm` command to `cli/cmd/roottooling.go`, the command failed to take effect because `cli/cmd/rootcore.go:113` contained an earlier entry `{[]string{"commit", "cm"}, func() error { return runCommitCLI(argsTail()) }}` that intercepted the command arguments and passed them to the legacy passthrough runner.

---

## 2. Root Cause

1. **Macro Execution:** Windows PowerShell's `rm` is an alias for `Remove-Item`, which throws a terminating `ItemNotFoundException` if the target is absent, and the locally active GitMap binary had not been updated to include the `AdaptCommandForPlatform` PowerShell loop wrapper.
2. **Routing Shadowing:** In GitMap's two-tier command dispatching table, `rootcore.go` is evaluated before `roottooling.go`. Having `"commit"` and `"cm"` registered in both tables resulted in silent shadowing of the new implementation.

---

## 3. Resolution

1. **Synchronize Binary Currency Across All User Profiles:** Rebuilt `gitmap.exe` from `cli/` and copied the newly compiled binary directly to all active target paths:
   - `./gitmap.exe` and `./bin/gitmap.exe`
   - `%LOCALAPPDATA%\gitmap-cli\gitmap.exe`
   - `%LOCALAPPDATA%\gitmap\gitmap.exe`
   This ensures that any active session executes the platform-adaptive `safe-rm` PowerShell loop without terminating `ItemNotFoundException`.
2. **Consolidate Commit Dispatching:** Updated `cli/cmd/rootcore.go` to use `constants.CmdCommit`, `CmdCommitAlias`, `CmdCommitAlias2`, and `CmdCommitAlias3` pointing to `runCommit(argsTail())`, removed the duplicate entry in `cli/cmd/roottooling.go`, and purged `cli/cmd/commit_cli.go`.
3. **Coding Guideline & Test Verification:** Decomposed `executeCommit` in `cli/cmd/commit_cmd.go` to keep all functions <= 15 lines, added package documentation, and implemented unit tests in `cli/cmd/commit_cmd_test.go` covering flag parsing and alias constants.

---

## 4. Prevention & Learnings

- **Single Dispatch Registry:** Never declare command aliases in multiple dispatch slice tables (`rootcore.go` vs `roottooling.go`); keep every command in a single authoritative table.
- **Runtime Binary Currency:** After introducing platform adapters or command improvements, always ensure the active runtime binaries (`gitmap.exe` and `%LOCALAPPDATA%\gitmap-cli\gitmap.exe`) are recompiled and verified.
- **Test-Backed Validation:** Always accompany new CLI commands with comprehensive flag and constant unit tests.
