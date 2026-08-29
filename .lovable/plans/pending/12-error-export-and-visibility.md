# 12 Error Export and Visibility Settings

## Parent Task Goal

Give users control over error visibility (full traces vs simple messages) via a settings configuration, and provide a dedicated `gitmap error export <file>` command to extract the last failure (so users can easily copy/paste errors into AI tools or issues without manual terminal scraping).

## Architectural Plan

1. **Settings / Config:**
   - Add `ErrorDisplay string json:"errorDisplay"` to `model.Config`.
   - Update `model.DefaultConfig()` to set it to `"full"`.
   - Update `config/validate.go` and `config/validate_shape.go` to validate `"errorDisplay"` as either `"full"` or `"simple"`.

2. **Error Visibility (The "simple" mode):**
   - In `gitmap/cmd/root.go`'s `runDispatch` panic/error recovery, we currently call `cliexit.Reportf` which prints the entire error.
   - We need to load the config to check `ErrorDisplay`.
   - If `ErrorDisplay == "simple"`, `apperror` (or `cliexit`) should format the error as a single-line simple message without the stack/tree.
   - Wait, `cliexit.Reportf` delegates to `errreport.Format` or something? No, it just prints it. We will modify how the error is printed based on the config.

3. **The `error` Command (`error export`):**
   - Create `gitmap/cmd/error_cmd.go` handling `gitmap error export [file]`.
   - The CLI already logs command audits (`gitmap/cmd/root.go`). Does the DB store the actual string error of the last failed command?
   - If not, we will intercept errors at the top level (`root.go`) and write them to a temporary file: `.gitmap/last_error.log`.
   - `gitmap error export <file>` simply copies `.gitmap/last_error.log` to the specified path, or prints it to standard output if no path is given.

## Subtasks

- **Subtask 1:** Add `ErrorDisplay` configuration logic.
- **Subtask 2:** Modify the global error handler to respect `ErrorDisplay`, and write the last encountered error to a local file.
- **Subtask 3:** Implement the `gitmap error` command and its `export` sub-action.
- **Subtask 4:** Update the documentation and user help text.

## Code Review Guide

- Do not swallow errors; use `apperror`.
- Boolean variables must be prefixed with `is`, `has`, `should`, or `can`.
- No generic variable names (`temp`, `data`, `res`).
