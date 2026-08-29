# Subtask 01: VSCode Commands Error Handling Refactoring

> **Parent Plan:** `16-error-management-audit.md`  
> **Files:** `gitmap/cmd/vscodeworkspace.go`, `gitmap/cmd/vscode_cmd.go`

## Objective

Replace all `panic("error")` calls in `vscodeworkspace.go` and `vscode_cmd.go` with domain-specific `apperror.NewWithDetails` / `apperror.WrapWithDetails` and route to `cliexit.HandleError`.

## Action Steps

1. Inspect `vscodeworkspace.go` (lines 48, 162, 178) and `vscode_cmd.go` (lines 19, 39, 46, 58).
2. Construct structured errors with:
   - `Op`: e.g. `cmd.vscodeworkspace`, `cmd.vscode.add`, `cmd.vscode.rm`
   - `Code`: `E1001` - `E1005`
   - `Type`: `apperror.ErrorTypeValidation` / `apperror.ErrorTypeExecution`
   - `Severity`: `apperror.SeverityError`
   - `Creator`: `cmd.vscode`
   - `Ctx`: map with arguments and paths.
3. Call `cliexit.HandleError(appErr, 1)`.
