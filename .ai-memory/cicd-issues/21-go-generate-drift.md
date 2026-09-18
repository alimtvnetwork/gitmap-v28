# 21-go-generate-drift

## Error Summary

The CI/CD pipeline failed with the error:
`Error: Generated files are out of sync with constants.`
`Error: Run 'cd gitmap && go generate ./...' locally and commit the result.`

## Root Cause Analysis

In a previous task, new CLI commands (`lowercase` and `fix-seq-files`) were added to the CLI constants (`gitmap/constants/constants_cli.go`). While the AST tests and help files were strictly updated, the `go generate ./...` command was not run. As a result, the auto-generated file `gitmap/completion/allcommands_generated.go` drifted out of sync with the true source code, which the CI anti-drift guard correctly blocked.

## Solution Applied

Ran `go generate ./...` inside the `gitmap/` directory, which successfully regenerated `gitmap/completion/allcommands_generated.go` to include the new commands.

## What NOT to Repeat

- NEVER add or modify CLI constants, shell commands, or help text without subsequently running `go generate ./...` to prevent auto-generated file drift.
