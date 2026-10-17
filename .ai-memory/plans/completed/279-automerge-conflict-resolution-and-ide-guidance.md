# Plan: 279-automerge-conflict-resolution-and-ide-guidance

## Status: COMPLETED
- **Task ID:** 279-automerge-conflict-resolution-and-ide-guidance
- **Milestone:** Milestone 35 (Automerge Optimization, IDE Merge Guidance & Conflict CLI)
- **Specification:** `02-spec/21-app/279-automerge-conflict-resolution-and-ide-guidance/`

## User Request (Verbatim)
```text
auto merging failed, try to see how to we can improve our automerge to avoid issues in the future, if issues we should provide manual suggestion for the user to fix using IDE like vscode and provide gitmap based commands to run and resolve conflicts, can you please do that and make sure to have helptext and ui help
```

## Subtask Breakdown

- [x] **Subtask 01:** Automerge optimization in `cli/cloner/safe_pull.go` (histogram diff, rerere auto-activation helper, and conflicted file parser).
- [x] **Subtask 02:** IDE detection & launcher helper in `cli/cmdconflict/ide_launcher.go` (`code`, `cursor`, `windsurf`, `codium`).
- [x] **Subtask 03:** CLI commands `gitmap resolve <repo>` and `gitmap conflicts [repo]` in `cli/cmdconflict/resolve_cmd.go` and wiring in `cli/cmd/`.
- [x] **Subtask 04:** Extend `gitmap fix <repo>` actions (`resolve`, `theirs`, `ours`, `abort`, `rerere`) in `cli/cmdfix/fix_cmd.go`.
- [x] **Subtask 05:** Post-pull conflict prompt enhancement in `cli/cmdpull/pull_prompt.go` with Lipgloss card, conflicted file badges, and one-key `[v]` VS Code launcher.
- [x] **Subtask 06:** CLI help documentation (`cli/helpdoc/resolve.md`, `cli/helpdoc/conflicts.md`), `cli/cmdhelp/howto.go`, and UI data (`src/data/commands.ts`, `src/data/troubleshootingIssues.ts`).
- [x] **Subtask 07:** Quality verification (linters, Spec 279 update, version bump to `v6.561.0`, and atomic commit via `gitmap cpf`).
