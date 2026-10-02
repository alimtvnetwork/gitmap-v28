# Subtask 02: Zero-Allocation strutil.EqualFoldAny and CLI Refactoring

- **Subtask ID:** Task-02
- **Assigned Worker:** Worker 02
- **Owned Files:**
  - `cli/strutil/strutil.go`
  - `cli/cmdvscode/vscode_cmd.go`
  - `cli/cmdvmware/vmware.go`

## Instructions
1. Create `cli/strutil/strutil.go` with package `strutil`:
   - `EqualFoldAny(target string, candidates ...string) bool`: zero-allocation check matching `target` against any candidate using `strings.EqualFold`.
   - `EqualFoldAnyTrim(target string, candidates ...string) bool`: trims whitespace on `target` before evaluating `EqualFoldAny`.
   - `NormalizeLowerTrim(s string) string`: returns trimmed lowercased string.
2. In `cli/cmdvscode/vscode_cmd.go`, update `isVSCodeProjectSubcommand` and `routeVSCodeMaintenanceAction` to use `strutil.EqualFoldAny` instead of lowercasing arguments.
3. In `cli/cmdvmware/vmware.go`, update subcommand dispatch to use `strutil.EqualFoldAny`.
