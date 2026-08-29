# RCA: Top-Level Cmd Registry Missing Constants (TestTopLevelCmdRegistryMatchesAST)

## Error Summary

The `github.com/alimtvnetwork/gitmap-v28/gitmap/constants` test suite failed due to `TestTopLevelCmdRegistryMatchesAST`. The error reported:
`AST has 2 top-level Cmd constant(s) missing from topLevelCmds() registry: CmdCloneSync, CmdCloneSyncAlias`.

## Root Cause Analysis

1. We recently introduced the new `gitmap clone-sync` (`gitmap cs`) command.
2. During the implementation, we defined the constants `CmdCloneSync` and `CmdCloneSyncAlias` inside `gitmap/constants/constants_cli.go`.
3. The project contains a strict CLI parity test (`TestTopLevelCmdRegistryMatchesAST` in `cmd_constants_parity_test.go`) which uses the Go AST to scan `constants_cli.go` for all constants starting with `Cmd`.
4. It compares the discovered constants against the hardcoded `topLevelCmds()` map inside `gitmap/constants/cmd_constants_test.go` to ensure that every CLI command is explicitly tracked and accounted for.
5. Because we did not add `CmdCloneSync` and `CmdCloneSyncAlias` to the `topLevelCmds()` map, the AST scanner detected the discrepancy and failed the build. This caused the CI pipeline to fail.

## Solution Applied

Modified `gitmap/constants/cmd_constants_test.go` and added the missing constants to the `topLevelCmds()` map:
```go
"CmdCloneSync":             CmdCloneSync,
"CmdCloneSyncAlias":        CmdCloneSyncAlias,
```

## What NOT to Repeat

- Do not create new top-level `Cmd*` constants without also adding them to `topLevelCmds()` in `cmd_constants_test.go` (or marking them with `// gitmap:cmd skip`).
