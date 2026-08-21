# Subtask: Inverted Booleans (TS)
Status: ✅ Done

## Steps
1. Edit src/components/docs/DocsTooltip.tsx: line 65, extract !isValidElement(child) to isInvalid := !isValidElement(child)
2. Edit src/components/docs/TerminalDemo.tsx: line 39, extract !isPlaying to isPaused := !isPlaying
3. Audit and fix all remaining inverted booleans across TypeScript codebase.
