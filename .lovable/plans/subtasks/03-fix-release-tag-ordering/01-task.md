# Subtask 1: Reorder Release Workflow Steps

## Instructions

1. Open `gitmap/release/workflow.go`.
2. Locate the `performRelease` function.
3. Move Step 3 (`localdirs.MigrateLegacyDirs`), Step 4 (`writeMetadataIfRequired`), and Step 5 (`AutoCommit`) **above** Step 1 (`executeSteps`).
4. Re-number the steps logically.
5. Update the error variable declarations:
   - `err := writeMetadataIfRequired(...)`
   - `err = executeSteps(...)`
6. Run `go build ./...` and `go test ./...` in the `gitmap` folder to verify compilation and tests pass.
