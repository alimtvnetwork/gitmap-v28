# Subtask 23-02: Root Help Layout & Compaction Formatting

## Goal

Standardize the CLI root help table rendering in `gitmap/cmd/rootusage.go` to eliminate excessive column gaps:
- `maxCmdColumnWidth = 26`
- Descriptions aligned cleanly at column 30
- Long multiline commands indent descriptions by exactly 4 spaces on line 2 (`renderLongHelpRow`)

## Status: DONE

- Compacted terminal table width and verified help formatting tests.
