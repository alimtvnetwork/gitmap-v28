# Specification 164: Pull-All Fast Mode, State Templates DB Engine, Pre-Compiled Variables & Declarative Commit-In Config

- **Status:** `active`
- **Date:** 2026-09-26
- **Category:** Application Architecture / Pull Engine / Templates State DB / Commit-In Replay

---

## 1. Verbatim User Request & Visual Evidence

### Verbatim User Request
> improve gitmap pa/pull all # to not to display the status anymore to do things faster
> the status will be display if
> gitmap pa --status or gitmap pat/gitmap pull all table # this will behave as current implementation, can you please do it
> also
> gitmap pa --json should give json output of the summary , clear???
> run gitmap pa first to proceed with implementation
>
> https://prnt.sc/n2CR1qfil2m-
> https://prnt.sc/NepOvEamHW5z
> https://prnt.sc/th5Wp3Bg38XS
>
> cli\helptext\seo-templates.md
>
> Here are a couple of SEO rules and a couple of things that we need to change during our commits. For example, if a title has changes, we will change that to-- So we will have some replacement techniques, where each one is going to go. When we have the prefix or postfix template, let's make sure the questions are in `#` and below is the answer, and that starts directly with the reasoning (`Because ...`). Give a detailed answer and combine the 20 things that we have, mix and match, and make it 40. Have good logical references and good numbers that validate those criteria.
> Also, do not feed the sponsor templates by default. Remove the SEO files (`seo_templates.go`, `seo-templates.md`) from the Git repo as a whole and keep the test JSON and test PowerShell scripts in `.ai-memory/temp/` so they don't get committed.
> We will feed the templates in 3 ways: CLI (`add`, `remove`, `edit`, `import`, `export`, `ui`), programmatic config JSON (`--config`), and Web UI (`gitmap templates ui` opening in browser). Store templates in the templates state DB (`templates.db`), not in the current project repo. Support variables (`variables`), pre-compile templates with variables before loops, and deduplicate imports via export hash ID (`exportId`), matching by `id` then `slug` to update.
> In the commit-in config JSON, support auto-importing template JSON files, line skipping (`starts_with`, `ends_with`, `contains`, `regex`) to strip lines like `Co-authored-by:`, blank line gap before templates, and replacing generic `"Changes"` commit titles using `$files.2.names` (first two filenames comma-separated, or single filename if 1 file). Also auto-create the target repo and support `--cd`.

### Visual Screenshots Reference
- [Screenshot 1: Commit Body Missing Blank Gap & Inline Question](../../assets/screenshots/templates-seo-01.png)
- [Screenshot 2: Generic "Changes" Titles & Unstripped Co-authored-by Line](../../assets/screenshots/templates-seo-02.png)
- [Screenshot 3: GitHub Desktop History Showing Repeated "Changes" Titles](../../assets/screenshots/templates-seo-03.png)

---

## 2. Pull-All Fast Mode (`pa`), Table Mode (`pat` / `--status`), and `--json`

1. **Default Fast Mode (`gitmap pa` / `gitmap pull all` / `gitmap pull-all`):**
   - Executes parallel git pull across all tracked repositories using the progress bar, then prints a concise active summary list (same fast rendering style as `pae`, without running slow post-pull `git branch -r`, PR status, and tag inspection per repo required by the full status table).
2. **Table Mode (`gitmap pa --status`, `gitmap pat`, `gitmap pull all table`, `gitmap pull-all-table`):**
   - Preserves the full post-pull status table (`RenderPullBatchTable`) inspecting latest remote branch, PR status, release tag, and SHA.
3. **JSON Mode (`gitmap pa --json` / `gitmap pull all --json`):**
   - Suppresses interactive banners/progress bars and emits a structured JSON summary (`Total`, `PulledCount`, `SuccessCount`, `FailedCount`, `States`, `DurationMs`) to `os.Stdout`.

---

## 3. Hardcoded SEO & Committed Test File Removal

- Remove `cli/cmd/commitin/seo_templates.go` from the repository.
- Remove `cli/helptext/seo-templates.md` from the repository.
- Remove `scripts/run-e2e-commit-pull.ps1` and `cli/cmd/commitin/e2e/commit_pull_tempe2e_test.go` from the repository.
- Do not inject hardcoded RiseUpAsia templates by default in `commitin.Parse`.
- Store the 40 mix-and-match SEO templates JSON in `.ai-memory/temp/seo-templates.json` and the 2-line migration runner PowerShell script in `.ai-memory/temp/run-migration-test.ps1` (never committed).

---

## 4. State Templates & Variables Database Engine (`gitmap-templates.db`)

1. **Split State Database (`gitmap-templates.db` in `store.BinaryDataDir()`):**
   - `TemplateCategory`: `CategoryId`, `Slug`, `Name`, `ParentSlug`, `Description`, `IsDefault`, `CreatedAt`, `UpdatedAt`. Default categories: `seo`, `prompts`, `ui-ux`, `prefix`, `pr-descriptions`.
   - `TemplateItem`: `ItemId` (string UUID/hash or user ID), `CategorySlug`, `SubCategorySlug`, `Slug`, `Title`, `Text`, `AdditionalJson`, `CreatedAt`, `UpdatedAt`.
   - `TemplateVariable`: `VarKey`, `Scope`, `VarValue`, `Description`, `UpdatedAt`.
   - `TemplateImportHistory`: `ExportHashId` (PRIMARY KEY), `SourcePath`, `ItemCount`, `VarCount`, `ImportedAt`.
2. **Import / Export Hash Deduplication & Upsert:**
   - Exporting generates a deterministic SHA-256 `exportId` computed from the canonicalized items + referenced variables.
   - Exporting a template or category automatically includes all referenced `$VAR` / `${VAR}` definitions in the `variables` map of the exported JSON.
   - Importing checks `TemplateImportHistory` for `exportId`: if already present (and `--force` is not set), it skips immediately as a fast no-op.
   - When importing items, it matches by `id` (`ItemId`) first; if not found, it matches by `slug` (`Slug`) and updates the existing record in place. Variables in the payload are upserted into both `TemplateVariable` and `config.SetVariable`.
3. **Pre-Compilation Before Loops:**
   - `PrecompileCategoryTemplates(categoryOrSlug string, extraVars map[string]string) ([]CompiledTemplate, error)` loads all matching templates and variables once into memory and expands all static `$VAR` / `${VAR}` placeholders before loop execution (`commit-in` / `commit-pull`).
4. **CLI & Browser UI:**
   - `gitmap templates ls [--category <cat>] [--json]`
   - `gitmap templates add --category <cat> --slug <slug> --title <title> --text <text> [--id <id>] [--additional <json>]`
   - `gitmap templates edit <id|slug> [--title <t>] [--text <text>] [--category <cat>] [--additional <json>]`
   - `gitmap templates remove <id|slug>` (aliases: `rm`, `delete`)
   - `gitmap templates import <file.json> [--force]`
   - `gitmap templates export <dest.json> [--category <cat>] [--id <id|slug>]`
   - `gitmap templates var set <key> <value>`, `gitmap templates var ls`
   - `gitmap templates ui` (starts a local HTTP server with an interactive Templates & Variables CRUD + Import/Export Web UI and opens the default browser).

---

## 5. Declarative Commit-In / Commit-Pull Config JSON (`--config`)

`gitmap commit-in` and `gitmap commit-pull` accept `--config <file.json>` (or positional/flag config) supporting:
- `target`: Target repository path (automatically created locally + remotely via `create-repo` if missing).
- `inputs`: Array of source repositories or range patterns (e.g. `https://github.com/alimtvnetwork/gitmap-v{2..28}`).
- `cd`: Boolean (`true` writes shell handoff to `cd` into target repo).
- `tree`: Boolean (`true` prints the preflight PR/merge tree).
- `prMode`: e.g. `"merges"`.
- `finalSync`: Boolean.
- `imports`: Array of template JSON file paths to auto-import (hash-deduplicated) before running.
- `lineSkippers`: Array of rules `{ "mode": "starts_with" | "ends_with" | "contains" | "regex", "pattern": "...", "ignoreCase": true }` to strip lines like `Co-authored-by:` or `X-Lovable-Edit-ID:`.
- `newlineGap`: Boolean (default `true`) ensuring an empty line (`\n\n`) separates the commit body from any injected prefix/suffix template.
- `titleReplacements`: Array of rules `{ "matchMode": "equals" | "contains" | "regex", "match": "Changes", "replacement": "$files.2.names: $seo.title" }` where:
  - `$files.1.name` resolves to the base filename of the first changed file.
  - `$files.2.names` resolves to the base filename if 1 file changed, or `file1.ext, file2.ext` if 2 or more files changed.
  - `$seo.title` / `$template.title` resolves to the selected template's title/question (without leading `# `).
- `prefixTemplates` / `suffixTemplates`: Category slugs or template slugs to pre-compile from `gitmap-templates.db` and attach to commit bodies with a clean blank line gap.
