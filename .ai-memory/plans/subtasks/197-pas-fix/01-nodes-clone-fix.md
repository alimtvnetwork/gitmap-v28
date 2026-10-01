# Subtask 01: Nodes Clone Table Alignment & W3 Resilience

**Objective**: Fix fleet nodes clone table alignment and W3 reachability diagnostics as specified in Spec 194 and Spec 196.

## Requirements
- Fix table column gutter padding in `cli/cmdnodes/nodes_clone_table.go`.
- Ensure W3 offline vs unreachable diagnostics are properly distinguished.
- Prevent ANSI escape code leakage into column measurement calculations.

## Constraints
- Bounding Box: `cli/cmdnodes/*.go`
- Coding Rules: Positive booleans only, AppError wrappers, functions <= 15 lines.
