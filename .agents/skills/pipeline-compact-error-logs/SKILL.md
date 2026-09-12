---
name: pipeline-compact-error-logs
description: Autonomously filter ok lines from pipeline error logs by default and support detailed/verbose flags for full logs.
---

# Pipeline Compact Error Logs

## Overview
Filters passing `ok` lines, test status lines, and non-failure noise from `gitmap pipeline error-logs` (and its aliases `errorlogs`, `errors`, `err`) by default.
Supports `--detailed`, `--verbose`, `--v`, and `-v` flags to display full, uncompressed logs.

## Flag Matrix
- Default (compact): removes `ok lines` (`ok\t`, `ok `, `? `, `--- PASS`, `=== RUN`, `PASS`, `✔ ok`, `✔ Macro`).
- Non-compact (`--detailed`, `--verbose`, `--v`, `-v`): preserves all lines verbatim.
