# Spec 280: Acceptance Criteria & Evidence

## Acceptance Criteria

| ID | Category | Requirement | Verification Method |
| :--- | :--- | :--- | :--- |
| **AC-01** | CI/CD Telemetry | `gitmap pe` audit executed, CI build failures identified (undefined `decisionLog` and nested if linter errors) and resolved | `gitmap pe` output and file inspection |
| **AC-02** | Muse Training Priority | `llm_train_muse.go` prioritizes `repo-secrets` as default training source and compiles cleanly | Code inspection & `05-guideline-autofixer.py` |
| **AC-03** | Split Large Files CLI | `gitmap split large files` (and `gitmap split-large-files`) discovers files $> 10\text{ MB}$, supports `--dry-run`, multipart zip (`.zip`, `.z01`), media bypass, and history purge | Command execution & unit test |
| **AC-04** | Git LFS Conversion | `gitmap lfs convert` discovers files $> 5\text{ MB}$, updates `.gitattributes`, migrates history, and supports dry-run preview | Command execution & test |
| **AC-05** | Repository Flattening | `gitmap repo flatten` creates numbered successor, performs single flat squashed commit, marks old repo for deletion, and provides undo command | Execution test |
| **AC-06** | Irreversible Warning Gate | Commands modifying history prompt the user and show bold yellow irreversible warning when `-y` is passed | Code inspection |
| **AC-07** | Helptext & UI Help | `cli/helpdoc/split-large-files.md`, `lfs-convert.md`, `flatten.md`, `howto.go`, `commands.ts`, and `troubleshootingIssues.ts` updated | File inspection |
| **AC-08** | Linter & Release Hygiene | Relative paths linter passes, forbidden strings passes, version bumped to `v6.562.0`, pushed via `gitmap cpf` | Python linters and git status verification |

## Verification Evidence

- **AC-01 Verified:** `gitmap pe` telemetry executed; identified undefined `decisionLog` and nested if linter errors; restored `decisionLog` struct/methods in `cli/cmd/llm/llm_train_muse.go` and flattened nested if blocks in `cli/cmdconflict/resolve_cmd.go` and `cli/cmdfix/fix_cmd.go`.
- **AC-02 Verified:** `defaultResolveTrainingPath()` in `cli/cmd/llm/llm_train_muse.go` prioritizes `repo-secrets` checkouts first for `12-muse-ai/`, with forward-compatible `repo-cache` fallback.
- **AC-03 Verified:** `cli/cmdsplit/split_cmd.go` and `cli/cmdsplit/multipart_zip.go` implement `gitmap split large files` (and `gitmap split-large-files`) with 10MB default threshold, `--dry-run` summary table, multipart zip chunking (`.zip`, `.z01`, `.z02`), `zip.Store` for media, `zip.Deflate` for text/code/data, and history purge.
- **AC-04 Verified:** `cli/cmdlfs/lfs_convert.go` implements `gitmap lfs convert` with 5MB default threshold, `.gitattributes` tracking pattern generation, history migration, and dry-run preview.
- **AC-05 Verified:** `cli/cmdflatten/flatten_cmd.go` and `cli/cmdundo/undo_flatten.go` implement `gitmap repo flatten` creating `<repo>-flat-v1`, flat squashed commit, archiving original to `<repo>-marked-for-deletion`, renaming successor to canonical name, and providing `gitmap repo undo-flatten <repo>`.
- **AC-06 Verified:** Destructive commands prompt the user with interactive confirmation or output bold yellow irreversible warnings when `-y` / `--yes` is supplied.
- **AC-07 Verified:** Authored `cli/helpdoc/split-large-files.md`, `cli/helpdoc/lfs-convert.md`, `cli/helpdoc/flatten.md`, updated `cli/cmdhelp/howto.go` with 3 how-to recipes, registered commands in `src/data/commands.ts`, and added troubleshooting card for large files in `src/data/troubleshootingIssues.ts`.
- **AC-08 Verified:** Relative paths linter passed (0 absolute paths across 8,978 files), forbidden strings linter passed, version bump to `v6.562.0`.

