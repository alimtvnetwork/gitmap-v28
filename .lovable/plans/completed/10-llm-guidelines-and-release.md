# 10 LLM Commit Guidelines and Final Release

## Parent Task Goal

Define explicit instructions in `llm.md` prohibiting AI agents from using raw `git commit` or `git push` commands. Instead, enforce the usage of `gitmap feature`, `gitmap bar`, `gitmap commit-in`, and `gitmap release` for grouped, orchestrated commits. Following this, persist the strict error handling rules from the recent `apperror` refactoring into the coding guidelines, and finalize the turn by orchestrating a grouped commit and release sequence using the newly mandated commands.

## Subtasks Execution Strategy

1. **Subtask 1: Documentation & Guidelines Update**
   - Update `llm.md` to forbid raw `git commit/push` and document Gitmap's orchestration commands.
   - Inject the new Error Architecture rules (the `apperror` returning rules) into `.lovable/coding-guidelines/coding-guidelines.md` and `spec/02-coding-guidelines/00-overview.md`.

2. **Subtask 2: Grouped Commit via Gitmap**
   - Execute `gitmap feature "error-architecture-and-llm-guidelines"`.
   - Add all modified files.
   - Execute `gitmap commit-in "fix: resolve 150+ apperror compile issues and update guidelines"`.
   - (Note: if `gitmap feature` is not available as a binary yet, simulate the grouping strategy using standard git but with explicit node markers).

3. **Subtask 3: Final Release**
   - Bump version using `.lovable/release/bump_versions.py --type minor` (or fallback script).
   - Generate release notes, pin to root `readme.md`.
   - Commit the release bump and push.
