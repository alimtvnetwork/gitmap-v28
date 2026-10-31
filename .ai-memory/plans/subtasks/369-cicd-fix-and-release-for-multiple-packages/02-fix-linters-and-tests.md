# Subtask 02 - Fix linters and tests

Date: 31-Oct-26
Owner: solo worker (parent task forbids subagents).

- Fix only what CI fails on, smallest diff, with WHY comments.
- Files outside .github/workflows/ are pushed normally.
- Any .github/workflows/ change is committed locally only and
  reported BLOCKED with its SHA (VM gh token lacks workflow scope).
- Pre-existing findings CI does not fail on are reported honestly,
  not chased.
