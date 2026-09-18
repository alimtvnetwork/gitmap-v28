# Subtask 04: AGY Commands & Detailed Terminal Help Parity Verification

## Scope
- Inspect `cli/cmdagy/` and `cli/cmdprompttemplate/` command suites.
- Verify subcommands:
  - `gitmap agy fix-pipeline` / `gitmap aef`: Direct injection to `agy.exe`, full error log embedding, absolute paths, multi-project batching with `--all`, `--projects`, `--limit`, `--reset-batch`.
  - `gitmap agy rerun`: Rerun by task ID, prefix templates, remote delegation across SSH/cluster.
  - `gitmap agy list-prompts`: Listing recent prompts, `--open-commit-multi-vscode`, inspect commit flags.
  - `gitmap agy scan` / `gitmap prompt-templates`: Finding, CRUD, import/export.
- Verify comprehensive terminal help text (`--help`) and markdown help docs (`cli/helptext/agy.md`, `cli/helptext/pipeline.md`, `cli/helptext/prompts-template.md`, `cli/helptext/ssh.md`).

## Acceptance Criteria
- [x] All AGY subcommands expose rich descriptive `--help` flags.
- [x] `cli/helptext/agy.md` contains dedicated sections for rerun, list-prompts, scan, templates, triad delegation, and fix-pipeline.
- [x] Synonyms (`aef`, `fp`, `pipeline errors agy fix`) route directly to the fix-pipeline engine.
- [x] Multi-project batching behavior is documented with clear usage examples.

## Completed Changes
- Confirmed AGY command suite registration in `cli/cmdagy/agy_cmd.go` and `cli/cmd/root.go`.
- Confirmed full terminal help rendering via `termhelp.RenderMenu` in `cli/cmdagy/agy_help.go` and `agy_help_automation.go`.
- Confirmed full parity with `cli/helptext/agy.md`, `cli/helptext/pipeline.md`, and `src/data/commands.ts`.
