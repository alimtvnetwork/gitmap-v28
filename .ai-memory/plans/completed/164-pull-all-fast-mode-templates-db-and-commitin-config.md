# Plan 164: Pull-All Fast Mode, State Templates DB Engine, Pre-Compiled Variables & Declarative Commit-In Config

- **Status:** `completed`
- **Completed Date:** 2026-09-26
- **Canonical Spec:** [02-spec/21-app/164-pull-all-fast-mode-templates-db-and-commitin-config.md](../../../02-spec/21-app/164-pull-all-fast-mode-templates-db-and-commitin-config.md)

---

## Verbatim User Request
> improve gitmap pa/pull all # to not to display the status anymore to do things faster
> the status will be display if
> gitmap pa --status or gitmap pat/gitmap pull all table # this will behave as current implementation, can you please do it
> also
> gitmap pa --json should give json output of the summary , clear???
> run gitmap pa first to proceed with implementation
> https://prnt.sc/n2CR1qfil2m-
> https://prnt.sc/NepOvEamHW5z
> https://prnt.sc/th5Wp3Bg38XS
> cli\helptext\seo-templates.md
> When we have the prefix or postfix template, let's make sure the questions are in `#` and below is the answer, and that starts directly with the reasoning (`Because ...`). Give a detailed answer and combine the 20 things that we have, mix and match, and make it 40. Have good logical references and good numbers that validate those criteria.
> Also, do not feed the sponsor templates by default. Remove the SEO files (`seo_templates.go`, `seo-templates.md`) from the Git repo as a whole and keep the test JSON and test PowerShell scripts in `.ai-memory/temp/` so they don't get committed.
> Store templates in the templates state DB (`templates.db`), not in the current project repo. Support variables (`variables`), pre-compile templates with variables before loops, and deduplicate imports via export hash ID (`exportId`), matching by `id` then `slug` to update.
> In the commit-in config JSON, support auto-importing template JSON files, line skipping (`starts_with`, `ends_with`, `contains`, `regex`) to strip lines like `Co-authored-by:`, blank line gap before templates, and replacing generic `"Changes"` commit titles using `$files.2.names`. Also auto-create the target repo and support `--cd`.

---

## Consolidated Subtask Summary & Verified Outcomes

- [x] **Subtask 01 — Pull-All Fast Mode, Status Table (`--status` / `pat`), and `--json` Summary (`cli/cmdpull/pull.go`, `cli/cmdpull/pullall.go`, `cli/cmdpull/pull_batch_json.go`, `cli/cmd/rootcore.go`, `cli/cmd/rootgit.go`):**
  - Default `gitmap pa` / `gitmap pull all` now skips slow post-pull table inspection and prints a fast concise active summary.
  - `gitmap pa --status`, `gitmap pat`, and `gitmap pull all table` preserve full post-pull status table rendering.
  - `gitmap pa --json` emits a structured JSON summary of the pull batch to stdout.
- [x] **Subtask 02 — Hardcoded SEO & Committed Test File Removal:**
  - Removed `cli/cmd/commitin/seo_templates.go`, `cli/helptext/seo-templates.md`, `02-spec/21-app/148-riseup-asia-templates-and-prompt-enhancement.md`, `scripts/run-e2e-commit-pull.ps1`, `cli/cmd/commitin/e2e/commit_pull_tempe2e_test.go`, `migrate.ps1`, and `scripts/migrate-gitmap-v28.ps1` from Git.
- [x] **Subtask 03 — State Templates DB (`gitmap-templates.db`), Variables, Pre-Compilation & Web UI (`cli/store/templates_split_*.go`, `cli/cmd/templates_state_cli.go`, `cli/cmd/templates_ui_server.go`, `cli/cmd/templatescli.go`):**
  - Implemented `gitmap-templates.db` split state database with `TemplateCategory`, `TemplateItem`, `TemplateVariable`, and `TemplateImportHistory`.
  - Implemented SHA-256 `exportId` import deduplication, ID-first then Slug upsert matching, referenced `$VAR` export extraction, and `PrecompileTemplates` in-memory variable expansion prior to loop execution.
  - Added CLI subcommands (`ls`, `add`, `edit`, `remove`, `import`, `export`, `var`) and interactive browser studio (`gitmap templates ui`).
- [x] **Subtask 04 — Declarative Commit-In / Commit-Pull `--config`, Line Skippers, Blank Line Gap & `$files.2.names` (`cli/cmd/commitin/config_json.go`, `cli/cmd/commitin/message/strip.go`, `cli/cmd/commitin/message/affix.go`, `cli/cmd/commitin/message/title_replace.go`, `cli/cmd/commitin/message/pipeline.go`):**
  - Added `--config <file.json>` (and single-positional JSON file auto-detection) with automatic template import, variable pre-compilation, `starts_with`/`ends_with`/`contains`/`regex` line skippers, guaranteed `\n\n` blank line separation before suffix templates, and `$files.1.name` / `$files.2.names` + `$seo.title` dynamic title replacement for generic `"Changes"` commits.
- [x] **Subtask 05 — 40 Mix-and-Match SEO Templates JSON (`.ai-memory/temp/`) & Documentation:**
  - Generated uncommitted `.ai-memory/temp/seo-templates.json` (40 `# Why ...?\nBecause ...` templates with concrete metrics and variables), `.ai-memory/temp/commit-pull-config.json`, and 1-line `.ai-memory/temp/run-migration-test.ps1`.
  - Updated `cli/helptext/pull-all.md`, `cli/helptext/templates.md`, `cli/helptext/commit-pull.md`, and `src/pages/CommitInExamples.tsx`.
