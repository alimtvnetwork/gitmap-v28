# Subtask 03: Automatic Clipboard Integration & Terminal Confirmation

## Objective
Automatically copy the pipeline error logs / report to the system clipboard upon execution and print a clear confirmation notice in the terminal.

## Implementation Steps
1. Capture or format the text report generated for the terminal or log file.
2. Use `github.com/atotto/clipboard.WriteAll(content)` to copy the text to the clipboard.
3. If clipboard write is successful, print:
   `\n  📋 Copied pipeline error logs to clipboard\n` (or `\n  📋 Copied pipeline status to clipboard\n` when clean).
4. Guard against headless/CI environments where clipboard may be unavailable, ensuring zero crashes.

## Constraints
- Max 15 lines per function.
- Blank line before every return.
