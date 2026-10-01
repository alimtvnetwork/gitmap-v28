# RCA-26: RootCore Commit Dispatch Shadowing and Stale Binary Macro Execution Failure

Spec Reference: [02-spec/22-app-issues/60-rootcore-commit-dispatch-shadowing-and-stale-binary-rca.md](../../02-spec/22-app-issues/60-rootcore-commit-dispatch-shadowing-and-stale-binary-rca.md)

## 1. Symptom

1. **Macro Execution Failure (`rm test`):**
   Executing macro `alim1` (`gitmap alim1`) on Windows PowerShell failed with `exit status 1` (`ItemNotFoundException` from `Remove-Item`) on missing target directory `C:\Users\Alim\test` because the active runtime binary at `%LOCALAPPDATA%\gitmap-cli\gitmap.exe` was out of date.
2. **Commit Command Dispatch Collision:**
   Adding `gitmap commit` / `cm` in `cli/cmd/roottooling.go` failed to dispatch to `runCommit` because an earlier registration in `cli/cmd/rootcore.go` shadowed it with the legacy `runCommitCLI` passthrough.

---

## 2. Root Cause

Windows PowerShell's `Remove-Item` fails terminatingly on missing targets, the installed GitMap binary was built before safe removal shims were present, and `rootcore.go` evaluated before `roottooling.go`, shadowing the commit command entry point.

---

## 3. Resolution

1. Recompiled and installed the fresh GitMap binary to both `./gitmap.exe` and `%LOCALAPPDATA%\gitmap-cli\gitmap.exe`.
2. Consolidated `CmdCommit`, `CmdCommitAlias`, `CmdCommitAlias2`, `CmdCommitAlias3` in `cli/cmd/rootcore.go` routing to `runCommit(argsTail())`.
3. Deleted `cli/cmd/commit_cli.go` and removed the duplicate routing entry from `cli/cmd/roottooling.go`.
4. Decomposed `commit_cmd.go` to adhere to <= 15 line functions and added unit tests in `cli/cmd/commit_cmd_test.go`.

---

## 4. Prevention & Learnings

- Centralize dispatch table entries to prevent shadow dispatching.
- Recompile and verify runtime binaries whenever platform shims are introduced.
- Require dedicated unit tests for all new CLI command flags.
