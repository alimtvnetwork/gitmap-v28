# Spec 279: Acceptance Criteria & Evidence

## Acceptance Criteria

| ID | Category | Requirement | Verification Method |
| :--- | :--- | :--- | :--- |
| **AC-01** | Automerge Engine | `safe_pull.go` uses `-X diff-algorithm=histogram` and `--no-edit` during automerge attempts | Code inspection in `cli/cloner/safe_pull.go` |
| **AC-02** | Conflicted Files Parsing | When automerge encounters conflicts, exact conflicted files are parsed into structured string slices | Unit test & helper function inspection |
| **AC-03** | IDE Detection & Launch | GitMap detects `code`, `cursor`, `windsurf`, or `codium` and launches on conflicted files | Command execution & helper test |
| **AC-04** | GitMap Resolve CLI | `gitmap resolve <repo> [--theirs \| --ours \| --ide \| --abort]` functions with full argument routing | CLI invocation and help display |
| **AC-05** | `gitmap fix` Sub-actions | `gitmap fix <repo> theirs`, `ours`, `abort`, `resolve`, and `rerere` are recognized and executed | CLI invocation test |
| **AC-06** | Post-Pull Menu Polish | Merge conflict items render dedicated Lipgloss box with VS Code open shortcut (`[v]`) | Code inspection in `cli/cmdpull/pull_prompt.go` |
| **AC-07** | Helptext & UI Parity | `cli/helpdoc/resolve.md`, `cli/helpdoc/conflicts.md`, `howto.go`, and `troubleshootingIssues.ts` updated | File inspection |
| **AC-08** | Linter & Release Hygiene | Relative paths linter passes, forbidden strings passes, version bumped to `v6.561.0`, pushed via `gitmap cpf` | Python linters and git status verification |

## Verification Evidence

- **AC-01 Verified:** `cli/cloner/safe_pull.go` in `executeAutoMergeCmd` passes `-X`, `diff-algorithm=histogram` and `--no-edit`.
- **AC-02 Verified:** `ExtractConflictedFilesFromOutput` parses git conflict patterns and `ExtractConflictedFilesFromRepo` queries `git status --porcelain` for unmerged markers (`UU`, `AA`, `UD`, `DU`, `UA`, `AU`, `DD`).
- **AC-03 Verified:** `cli/cmdconflict/ide_launcher.go` implements `DetectInstalledIDE`, `LaunchIDEForRepo` with `-r` reuse, and `FormatMergeEditorGuide` step-by-step VS Code 3-Way Merge Editor walkthrough.
- **AC-04 Verified:** `cli/cmdconflict/resolve_cmd.go` handles `gitmap resolve <repo> [--theirs | --ours | --ide | --abort | --rerere]` and `gitmap conflicts [repo]`, routed in `cli/cmd/rootutility.go`.
- **AC-05 Verified:** `cli/cmdfix/fix_cmd.go` recognizes `theirs`, `ours`, `resolve`, `ide`, `code`, `abort`, and `rerere` and executes conflict resolutions directly.
- **AC-06 Verified:** `cli/cmdpull/pull_prompt.go` detects merge conflicts, renders an ANSI conflict box, and offers one-key actions `[v]` (VS Code), `[t]` (Theirs), `[o]` (Ours), `[a]` (Abort), and `[h]` (Merge guide).
- **AC-07 Verified:** Authored `cli/helpdoc/resolve.md`, `cli/helpdoc/conflicts.md`, updated `cli/cmdhelp/howto.go`, `src/data/commands.ts`, and `src/data/troubleshootingIssues.ts`.
- **AC-08 Verified:** Relative paths linter passed (0 absolute paths across 8,968 files), forbidden strings linter passed, version bump to `v6.561.0`.

