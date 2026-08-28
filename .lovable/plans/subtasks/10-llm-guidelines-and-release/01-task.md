# Subtask 1: Documentation & Guidelines Update

1. Read `llm.md` and append a strict section defining that AI agents MUST NEVER use raw `git commit` or `git push`. They must use `gitmap feature <name>`, `gitmap bar`, `gitmap commit-in`, and `gitmap release`.
2. Open `.lovable/coding-guidelines/coding-guidelines.md` and `spec/02-coding-guidelines/00-overview.md`.
3. Add a strict "Error Architecture" section specifying that all core functions and command handlers MUST return `*apperror.AppError` and never call `os.Exit(1)` directly, to ensure `finishCommandAudit` logs them to the DB.
4. Mark task as `STATUS: DONE` in `.lovable/temp-agents/task-10-01.md`.
