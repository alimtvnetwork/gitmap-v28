# CI Issue 82: Gencommands Drift, Top-Level Command Parity, and Antigravity Projects Dir RCA

## 1. Reproduction & Symptoms
- **CI Run ID:** [36279515487](https://github.com/alimtvnetwork/gitmap-v28/actions/runs/36279515487)
- **Failing Workflows:**
  - `CI / Generate Drift Check (Step: Run go generate)`
  - `CI / Top-Level Command Parity (TestTopLevelCmdRegistryMatchesAST)`
  - `CI / AGY Unit Tests (TestExecuteAgyRecreate_RealLifecycle)`
- **Symptoms:**
  1. `Generate Drift Check` failed at `Run go generate`:
     ```text
     gencommands: found 1 Cmd* constant(s) in const block(s) lacking the `// gitmap:cmd top-level` marker:
       constants.go:164  CmdCommon, CmdCommonAlias
     ```
  2. `cli/constants: TestTopLevelCmdRegistryMatchesAST` failed with:
     ```text
     AST has 6 top-level Cmd constant(s) missing from topLevelCmds() registry:
       CmdMigrate
       CmdRecreateRepo
       CmdRecreateRepoAlias
       CmdRestEnable
       CmdVar
       CmdVarAlias
     ```
  3. `cli/cmdagy: TestExecuteAgyRecreate_RealLifecycle` failed on headless Linux CI runners:
     ```text
     real recreate failed: [E9000:EXECUTION] all project recreate operations failed: [E9003:EXECUTION] failed to register project with Antigravity
     ```

## 2. Root Cause Analysis
1. **Unmarked Const Block in `constants.go`**: `CmdCommon` and `CmdCommonAlias` were added to `cli/constants/constants.go` in an unmarked const block. The completion code generator `cli/completion/internal/gencommands/main.go` scans all `constants/*.go` files and strictly rejects any const block defining `Cmd*` string constants without an explicit `// gitmap:cmd top-level` doc comment.
2. **Missing Registry Entries in `cmd_constants_test.go`**: Six newly added CLI commands (`migrate`, `recreate-repo`, `recreate`, `rest-enable`, `var`, `variable`) were declared in `cli/constants/constants_cli.go` under an opted-in block, but were never registered in `topLevelCmds()` in `cmd_constants_test.go`. The AST parity test ensures every top-level command constant is accounted for in `topLevelCmds()`.
3. **Missing Antigravity Projects Dir on Fresh/Headless Runners**: `workspacesync.SyncAntigravity` returns `false` if `~/.gemini/config/projects` does not exist on disk. On headless CI environments where Antigravity is not installed, the directory is absent, causing `registerFreshProject` and `TestExecuteAgyRecreate_RealLifecycle` to fail.

## 3. Code Fix
1. **Added Generator Markers**: Added `// gitmap:cmd top-level` to the const block in `cli/constants/constants.go` and marked `CmdCommon` and `CmdCommonAlias` with `// gitmap:cmd skip` to prevent alias collision with `CmdCommonsAlias` ("co").
2. **Synchronized Command Registry & Generated Completions**:
   - Added `CmdMigrate`, `CmdRecreateRepo`, `CmdRecreateRepoAlias`, `CmdRestEnable`, `CmdVar`, and `CmdVarAlias` to `topLevelCmds()` in `cli/constants/cmd_constants_test.go`.
   - Executed `go generate ./...` in `cli`, updating `cli/completion/allcommands_generated.go` with the newly registered commands.
3. **Initialized Projects Directory for Recreate Operations**:
   - In `cli/cmdagy/agy_recreate_ops.go` (`registerFreshProject`), ensured `configDir` is created via `os.MkdirAll(configDir, 0755)` before delegating to `workspacesync.SyncAntigravity`.
   - In `cli/cmdagy/agy_recreate_test.go` (`TestExecuteAgyRecreate_RealLifecycle`), ensured `configDir` exists prior to executing the test lifecycle.

## 4. Prevention
- Always execute `go generate ./...` locally whenever top-level `Cmd*` constants are declared or modified.
- Include all top-level commands in `topLevelCmds()` in `cli/constants/cmd_constants_test.go`.
- Ensure directory-dependent helper methods gracefully create or tolerate missing parent directories in clean test environments.
