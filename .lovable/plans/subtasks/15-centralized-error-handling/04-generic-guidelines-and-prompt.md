# Subtask 04: Generic Guidelines and AI Prompt Creation

master-plan: 15-centralized-error-handling-and-exit-architecture
subtask: 04-generic-guidelines-and-prompt
status: pending

## Goal

Author two standalone, highly reusable documents:
1. Generic Error Handling Guidelines Document (shareable across teams/languages).
2. AI System Prompt File with Error Handling Checklist.

## Specifications

1. **File 1**: `.lovable/coding-guidelines/centralized-error-handling-architecture.md`
   - Detailed conceptual explanation of why `panic` and bare `os.Exit` are both anti-patterns.
   - Comprehensive pattern: Error Envelopes, Domain Errors, Handlers, Execution Strategies, Caller Attribution.
   - Cross-language illustrations (Go, TypeScript, Python).
   - "Never be silent" mandate.
2. **File 2**: `.lovable/prompts/05-coding-guidelines/02-centralized-error-handling-checklist.md`
   - A ready-to-use system prompt for AI coding assistants.
   - Rigid error management checklist.
   - Detailed failure mode breakdown based on visual diff patterns.
   - Actionable instructions preventing regressions.

## Verification

- Both files formatted cleanly, strictly lowercase filenames, zero placeholders.
