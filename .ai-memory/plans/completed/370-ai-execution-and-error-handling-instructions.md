# Master Plan: 370-ai-execution-and-error-handling-instructions

- **Parent Task:** AI Execution and Error Handling Instructions
- **Parent Task Slug:** `370-ai-execution-and-error-handling-instructions`
- **Specification Directory:** `02-spec/21-app/370-ai-execution-and-error-handling-instructions/`
- **Subtasks Directory:** `.ai-memory/plans/subtasks/370-ai-execution-and-error-handling-instructions/`
- **Target Release Version:** `v6.645.0`
- **Task DB:** `.ai-memory/temp-agents/370-ai-execution-and-error-handling-instructions/agent-task.db`

---

## Executive Summary

When AI agents attempt to inspect directories or files, they occasionally confuse `gitmap ps` (PowerShell process execution) with directory listing (`ls`, `dir`) or task inspection commands, leading to:
```
gitmap ps .ai-memory/plans/subtasks/240-running-and-queued-prompt-detection-in-multi-instance
.ai-memory/plans/...: The term '.ai-memory/plans/...' is not recognized as a name of a cmdlet, function, script file, or executable program.
gitmap ps: execute failed: [E9000:EXECUTION] execute powershell: exit status 1
```
This task implements:
1. **Pre-flight argument inspection** in `gitmap ps`: detecting existing directories, non-script document files, and raw paths before invoking PowerShell.
2. **AI diagnostic inquiry**: asking the agent "Why did you run this? What was your expectation?" to trigger self-reflection.
3. **PowerShell execution explanation**: clarifying where and how `gitmap ps` runs (`pwsh` / `powershell.exe`) and why folders/documents cannot be executed directly.
4. **Actionable replacement suggestions**: pointing to `gitmap dir`, `gitmap ls`, `gitmap agent subtask ls`, `gitmap heatmap`, and properly quoted PowerShell commands (`gitmap ps "Get-ChildItem '<path>'"`).
5. **Two-to-three step learning solutions and examples**: giving concrete syntax examples so AI agents learn immediately without polluting logs with Go stack traces.
6. **Post-failure fallback interceptor**: catching cmdlet-not-found errors during execution and decorating with the same educational guidance.
7. **Skill & guideline synchronization**: updating `.agents/skills/gitmap/skill.md` and anti-pattern tables, then syncing across fleet repositories.

---

## Subtask Breakdown

| Code | Subtask ID | Title | Status | Owned Files |
|------|------------|-------|--------|-------------|
| `TASK-01` | `sub-121922-task-01` | Architecture & Diagnostic Specifications Authoring | DONE | `02-spec/21-app/370-ai-execution-and-error-handling-instructions/` |
| `TASK-02` | `sub-121928-task-02` | PowerShell Argument Pre-Flight Inspection & AI Reflection Engine | PENDING | `cli/cmd/cmd_powershell.go`, `cli/cmd/powershell_advice.go` |
| `TASK-03` | `sub-121928-task-03` | PowerShell Error Post-Processing & Fallback Interceptor | PENDING | `cli/cmd/powershell_interceptor.go`, `cli/cmd/cmd_powershell_test.go` |
| `TASK-04` | `sub-121928-task-04` | Update Skills, Cheatsheets & Anti-Pattern Guidelines | PENDING | `.agents/skills/gitmap/skill.md`, `02-spec/02-coding-guidelines/` |
| `TASK-05` | `sub-121928-task-05` | High-Speed Fleet Synchronization Across Repositories | PENDING | `cli/cmdsync/`, fleet repos |
| `TASK-06` | `sub-121929-task-06` | Compilation, Live Terminal Verification & Release Ceremony (v6.645.0) | PENDING | `cli/constants/constants.go`, `version.json`, `changelog.md` |

---

## Verification Criteria
1. `gitmap ps <existing-directory>` intercepts cleanly with Tokyo Night AI guidance card and exits code 2 without raw Go stack trace.
2. `gitmap ps <markdown-or-json-file>` intercepts cleanly with guidance suggesting `gitmap heatmap <file>` or `gitmap cat <file>`.
3. `gitmap ps "Get-ChildItem -Path '.'"` continues to execute valid PowerShell expressions normally.
4. AI inquiry ("Why did you run this? What was your expectation?") is surfaced prominently.
5. All 6 subtasks marked DONE in SQLite task DB.
6. Synchronized across fleet repos and released via `gitmap cpr "v6.645.0 ..."`.
