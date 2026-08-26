STATUS: DONE

# Task 04: Macro Steps Composition Tree Output Execution Record

## Overview
Implemented `gitmap/cmd/macro_tree.go` to provide UTF-8 box-drawing tree representation of macro steps and their composition hierarchy during execution and inspection.

## Created / Modified Files
- `gitmap/cmd/macro_tree.go`: Implemented `printMacroTree(macroName string)`, `printMacroStepsTree(targetMacro *macro.Macro)`, `printMacroTreeHeader`, `renderMacroStepList`, `selectTreeConnector`, `printStepBranch`, `renderMacroStepNode`, and `formatStepDescription`.
- `gitmap/cmd/macro_cmd.go`: Updated `runExecuteCmd` to call `printMacroStepsTree` before execution when steps > 0, refactored execution helpers into <=15 line functions.
- `gitmap/cmd/macro_tree_test.go`: Added comprehensive unit tests for macro tree rendering, step formatting, and branch connectors.

## Validation
- `go build ./...` passed cleanly with 0 errors.
- `go test -v ./cmd -run "TestPrintMacro|TestFormatStepDescription|TestSelectTreeConnector|TestPrintStepBranch"` passed cleanly.
- `go test ./cmd` passed cleanly.
- Functions strictly adhere to <= 15 lines max.
- Booleans follow positive `is`/`has` naming rules (`isLastStep`).
- Semantic naming followed without generic terms.
- UTF-8 box-drawing characters (`├──`, `└──`) and ANSI terminal colors (`constants.ColorCyan`, `constants.ColorWhite`, `constants.ColorDim`, `constants.ColorReset`) properly utilized.
