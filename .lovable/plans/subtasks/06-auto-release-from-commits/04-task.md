# Subtask 04: CG Update Verbosity

## Goal

Enhance `gitmap cg update` and `gitmap install` to print exactly which files were modified.

## Requirements

1. **File**: `gitmap/cmd/cg_worker.go` (and wherever install writes files).
   - Ensure that when files are written or overwritten, their relative paths are collected into a `[]string` array.
   - Render this array in the Lipgloss UI output as a bulleted list or a comma-separated list of explicitly modified files.
2. **Constraints**:
   - 15 line max per function.
