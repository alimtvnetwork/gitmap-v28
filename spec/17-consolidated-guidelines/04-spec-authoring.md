# Consolidated: Spec Authoring Guide — Complete Reference

**Version:** 3.2.0
**Updated:** 2026-04-16

---

## Purpose

This is the **standalone consolidated reference** for spec authoring conventions. An AI reading only this file must be able to create, organize, and validate specification files correctly — including folder structure, naming, required files, templates, cross-references, validation, exceptions, and the `.lovable/` context layer.

---

## Quick start (read this first)

1. **Copy the template:** `cp spec/_template.md spec/<NN-your-module>/01-your-file.md`
2. **Read the prohibitions:** [`03-strictly-avoid-quickref.md`](./03-strictly-avoid-quickref.md) — 30 seconds, prevents the most common CI failures.
3. **Before you commit:** run the four sync scripts in order — see [§X.2 Mandatory Execution Order](#x2-mandatory-execution-order). Skipping any step trips `Drift detected in version.json` in CI.
4. **Validate:** `python linter-scripts/validate-guidelines.py` — must report zero CODE-RED and zero STYLE violations.

---

## Folder Structure — Numbering Policy

- **`01–20`** — Core Fundamentals (principles, standards, integrations, research)
- **`21+`** — App-Specific content (`21-app`, `22-app-issues`, `23-app-db`, `24-app-ui-design-system`)

### Fixed Assignments

| # | Folder | Purpose |
|---|--------|---------|
| 01 | `01-spec-authoring-guide` | Authoring conventions |
| 02 | `02-coding-guidelines` | Cross-language coding standards |
| 03 | `03-error-manage` | Error management |
| 04 | `04-database-conventions` | Database naming and schema rules |
| 05 | `05-split-db-architecture` | Hierarchical SQLite pattern |
| 06 | `06-seedable-config-architecture` | Version-controlled config seeding |
| 07 | `07-design-system` | CSS variable-driven design system |
| 08 | `08-docs-viewer-ui` | Documentation viewer UI |
| 09 | `09-code-block-system` | Code block rendering |
| 10 | `11-powershell-integration` | PowerShell runner |
| 11 | `10-research` | Root-level research |
| 12 | `17-consolidated-guidelines` | AI-readable summaries (this folder) |
| 13 | `12-cicd-pipeline-workflows` | CI/CD pipeline specs |
| 14 | `14-update` | CLI self-update system |
| 15 | `18-wp-plugin-how-to` | WordPress plugin conventions |
| 21 | `21-app` | App-specific features |
| 22 | `22-app-issues` | App-specific issues |
| 23 | `23-app-db` | App-specific database schema |
| 24 | `24-app-ui-design-system` | App-specific UI and design system |

---

## Required Files

| File | Required When |
|------|---------------|
| `01-index.md` | **Always** — every module |
| `99-consistency-report.md` | Every module with 3+ files |
| `97-acceptance-criteria.md` | App/feature modules |
| `98-changelog.md` | Modules with version history |

---

## Naming Conventions

### Files

- Format: `NN-descriptive-name.md` (e.g., `02-architecture.md`)
- **Lowercase kebab-case** only
- No uppercase, no underscores, no spaces
- Numeric prefix is **mandatory** for sequencing
- Reserved prefixes: `00` (overview), `96` (AI context), `97` (acceptance criteria), `98` (changelog), `99` (consistency report)

### Folders

- Format: `NN-kebab-name/` (e.g., `05-split-db-architecture/`)
- Same lowercase kebab-case rule
- Two-digit zero-padded prefix (e.g., `01`, `08`, `33`)

### Non-Markdown Files

- Diagrams (`.mmd`), scripts, configs also use numeric-prefixed kebab-case
- Example: `01-self-update-workflow.mmd`, `02-update-cleanup-workflow.mmd`
- Data files (`.json`) use kebab-case, numeric prefix optional: `error-codes.json`

### Metadata Header

Every `.md` file **must** begin with:

```markdown

# Title of the Document

**Version:** X.Y.Z
**Updated:** YYYY-MM-DD

---
```

- Exactly one H1 title per file
- Semantic versioning (Major.Minor.Patch)
- ISO date format (YYYY-MM-DD)

---

## `01-index.md` — Required Content

Every `01-index.md` must include:

### Scoring Metrics

| Metric | Values |
|--------|--------|
| **AI Confidence** | `Production-Ready` / `High` / `Medium` / `Low` |
| **Ambiguity** | `None` / `Low` / `Medium` / `High` / `Critical` |
| **Health Score** | 0–100 from dashboard scanner |

### Other Required Sections

- **Version** and **Updated** date
- **Keywords** for searchability
- **Document Inventory** table listing all files with status
- **Cross-References** to related modules
- **Scoring table** with `01-index.md` present check

---

## Cross-References

- Use **file-relative paths only** — never root-relative `/spec/...`
- Always include `.md` extension
- Example: `[Split DB](../05-split-db-architecture/01-index.md)`
- Broken links are **blocking errors** in the dashboard scanner

---

## Module Templates

### CLI Module (3-Folder)

```
NN-module-name/
├── 01-index.md
├── 02-fundamentals.md
├── 02-features/
│   ├── 01-index.md
│   └── 01-feature-name.md
├── 03-issues/
│   └── 01-index.md
├── 97-acceptance-criteria.md
├── 98-changelog.md
└── 99-consistency-report.md
```

**Headless CLIs** may omit `02-frontend/` when UI is a separate module. Deploy folder **still uses prefix `03`** (not `02`).

### App / WordPress Project Template

```
{NN}-{app-name}/
├── 01-index.md
├── 02-fundamentals.md
├── 02-features/
│   ├── 01-index.md
│   └── 01-{feature-name}/
│       ├── 01-index.md
│       ├── 01-backend.md
│       ├── 02-frontend.md
│       ├── 03-wp-admin.md          # WordPress only
│       └── 99-consistency-report.md
├── 03-issues/
│   ├── 01-index.md
│   ├── 01-{issue-name}/            # Multi-file issues
│   │   ├── 01-index.md
│   │   ├── 01-investigation.md
│   │   └── 02-resolution.md
│   └── 02-{issue-name}.md          # Simple single-file issues
├── 97-acceptance-criteria.md
├── 98-changelog.md
└── 99-consistency-report.md
```

**Use when:** WordPress plugins/themes, web applications, platforms with user-facing features and tracked issues.

### Flat Module (Standards, Research)

```
NN-module-name/
├── 01-index.md
├── 01-topic-a.md
├── 02-topic-b.md
└── 99-consistency-report.md
```

### Template Comparison

| Aspect | CLI Template | App Template | Flat Module |
|--------|-------------|--------------|-------------|
| First content file | `01-backend/02-architecture.md` | `02-fundamentals.md` | `01-topic-a.md` |
| Feature organization | Flat files | `02-features/{NN}-{name}/` subfolders | N/A |
| Issue tracking | External (GitHub Issues) | `03-issues/` folder | N/A |
| Deploy specs | `03-deploy/` folder | Not included | N/A |

---

## `.lovable/` Folder — AI Context Layer

The `.lovable/` directory holds all AI-readable project knowledge.

### Canonical Structure

```
.lovable/
├── overview.md                     # AI onboarding — read FIRST
├── user-preferences                # Communication preferences
├── 29-plan.md                         # Current roadmap
├── suggestions.md                  # Pending improvement ideas
├── strictly-avoid.md               # ⛔ Quick-read hard prohibitions
├── prompt.md                       # AI prompt index
├── prompts/                        # Reusable AI prompts
├── memory/                         # Institutional knowledge
│   ├── index.md                    # Canonical index of all memory files
│   ├── architecture/               # System design decisions
│   ├── constraints/                # Hard constraints
│   ├── features/                   # Feature-specific knowledge
│   ├── issues/                     # Issue-specific knowledge
│   ├── processes/                  # Workflow processes
│   ├── project/                    # Project-level decisions
│   ├── standards/                  # Technical standards
│   └── style/                      # Code style rules
├── suggestions/                    # Suggestion details (one .md per suggestion)
│   └── completed/                  # Archived suggestions
├── pending-tasks/                  # Active work-in-progress
├── completed-tasks/                # Finished tasks
├── pending-issues/                 # Open issues
├── solved-issues/                  # Resolved issues
└── strictly-avoid/                 # ⛔ Detailed prohibition files
```

### Critical Rules

- **Memory folder is `.lovable/memory/`** — never `.lovable/memories/` (no trailing `s`)
- Memory subfolders use **kebab-case WITHOUT numeric prefixes**
- When adding a memory file, **always** update `memory/01-index.md`
- AI reading order: `overview.md` → `strictly-avoid.md` → `user-preferences` → `memory/01-index.md` → `29-plan.md`

### Task & Suggestion Lifecycle

- **Tasks:** `29-plan.md` → `pending-tasks/` → `completed-tasks/`
- **Suggestions:** `suggestions.md` → `suggestions/` → `suggestions/completed/`

---

## Authoring Workflow

### Creating a New Spec Module

1. **Select module number** — next available after highest existing (do NOT fill gaps)
2. **Create folder** — `{NN}-{module-name}/` using kebab-case
3. **Create `01-index.md`** — with all required sections (scoring, inventory, cross-refs)
4. **Write content files** — numbered sequentially (`01-`, `02-`, etc.)
5. **Create `99-consistency-report.md`** — structural health check
6. **Add `97-acceptance-criteria.md`** — if app/feature module
7. **Update synchronization targets** (see below)

### Updating an Existing Spec

1. **Read current file** — never modify without reading first
2. **Make changes** — preserve all unrelated content
3. **Bump version** — at minimum, bump the minor version
4. **Update `Updated` date** — ISO format (YYYY-MM-DD)
5. **Update `98-changelog.md`** — if the module has one
6. **Run validation** — zero violations before committing
7. **Update synchronization targets** — if structure changed

### Version Bump Rules

| Change Type | Version Bump | Example |
|-------------|-------------|---------|
| Fix typo, clarify wording | Patch (X.Y.+1) | 3.2.0 → 3.2.1 |
| Add new section, expand content | Minor (X.+1.0) | 3.2.0 → 3.3.0 |
| Restructure, rename, breaking change | Major (+1.0.0) | 3.2.0 → 4.0.0 |

---

## Validation & Linter Infrastructure

`linter-scripts/` is **mandatory**. AI must verify its presence before performing any validation.

### Required Files

| File | Purpose |
|------|---------|
| `validate-guidelines.py` | Main validator — zero CODE-RED or STYLE violations |
| `generate-dashboard-data.cjs` | Generates dashboard metrics |
| `run.sh` | Unix runner |
| `run.ps1` | Windows runner |

### Validation Command

```bash
python linter-scripts/validate-guidelines.py
```

### What the Validator Checks

| Check | Severity | Rule |
|-------|----------|------|
| `01-index.md` present in every module | CODE-RED | Blocking |
| Lowercase kebab-case naming | STYLE | Blocking |
| Unique numeric prefixes per folder | STYLE | Blocking |
| Metadata header (H1 + Version + Updated) | STYLE | Warning |
| Cross-reference link validity | CODE-RED | Blocking |
| Maximum folder depth (3 levels) | STYLE | Blocking |

### Validation Workflow

1. **Before committing:** Run `python linter-scripts/validate-guidelines.py`
2. **Zero violations required:** Both CODE-RED and STYLE violations must be zero
3. **Dashboard update:** Run `node linter-scripts/generate-dashboard-data.cjs` after structural changes
4. **CI enforcement:** The CI/CD pipeline runs validation automatically on every push

---

## Reproducibility Standard

All specifications must be **Blindly Reproducible** — sufficient for an AI to implement without human intervention. Target: **100% AI Reproducibility**.

### Reproducibility Checklist

- [ ] All technical terms defined or referenced
- [ ] All code examples are complete and runnable
- [ ] All schema definitions include types, constraints, and defaults
- [ ] All cross-references resolve to valid files
- [ ] No assumptions about reader's prior knowledge of the project

---

## Formatting Rules

- Triple-backticks in descriptive text are **prohibited** — use "triple-backtick blocks" as text
- 4-backtick fences required for nested code examples
- Abbreviations use full uppercase: `AI`, `DB`, `CI/CD`, `PHP`, `UI`
- Tables must have header row and alignment row
- Lists use `-` (not `*` or `+`)

---

## Documented Exceptions

These exceptions are permanent and tracked in `02-spec/01-spec-authoring-guide/11-exceptions.md`:

| # | Exception | Scope |
|---|-----------|-------|
| 1 | Non-contiguous module numbers (gaps allowed, never backfill) | `spec/` |
| 2 | `readme.md` without numeric prefix | Project-wide |
| 3 | Non-markdown data files (`.json`) | `spec/` modules |
| 4 | Legacy `C-XXX` suggestion file names | `suggestions/completed/` |
| 5 | Dual-purpose prefix `02` | `spec/` root |
| 6 | Memory folders without numeric prefixes | `.lovable/memory/` |
| 7 | CLI module without `02-frontend/` folder | Headless CLIs |
| 8 | Extra CLI subfolders beyond core 3 | Complex CLIs |
| 9 | Legacy `suggestions.md` without prefix | Memory root |
| 10 | 3-level depth in coding guidelines | `03-coding-guidelines/` |

**Policy:** New exceptions must be documented before creating non-conforming files.

---

## Synchronization

When specs change, update **all three** targets:

| # | Target | Purpose |
|---|--------|---------|
| 1 | `02-spec/spec-index.md` | Master index of all spec modules |
| 2 | `src/data/specTree.json` | UI tree (requires `content` field with full markdown) |
| 3 | `02-spec/17-consolidated-guidelines/` | Consolidated AI-readable summary |

### Synchronization Checklist

- [ ] New module added to `spec-index.md` with correct number and description
- [ ] `specTree.json` updated with new entry (including full `content` field)
- [ ] Consolidated guideline file created or updated in `17-consolidated-guidelines/`
- [ ] `17-consolidated-guidelines/01-index.md` file inventory updated
- [ ] `17-consolidated-guidelines/99-consistency-report.md` updated
- [ ] Gap analysis (`22-gap-analysis.md`) scores recalculated

---

---

## §X Version & Sync Workflow — The Mandatory Sequence

This section is the **single canonical reference** for how versions are bumped and how derived artifacts (`version.json`, `src/data/specTree.json`, `02-spec/dashboard-data.json`) stay synchronized. Skipping any step or running them out of order produces `Drift detected in version.json` and blocks CI.

### X.1 The Three Sync Scripts

| # | Script | Reads | Writes | When to Run |
|---|--------|-------|--------|-------------|
| 1 | `scripts/sync-version.mjs` | `package.json`, git HEAD | `version.json` (top-level + per-folder stats) | After **any** version bump or spec change |
| 2 | `scripts/sync-spec-tree.mjs` | `spec/**/*.md` | `src/data/specTree.json` | After **any** spec file add/rename/delete |
| 3 | `linter-scripts/generate-dashboard-data.cjs` | `version.json`, `spec/**/*.md` | `02-spec/dashboard-data.json` | After #1 and #2 (CI runs this automatically) |

### X.2 Mandatory Execution Order

```
1. Edit spec / source code
2. Bump package.json version (minor for code changes; patch for doc-only)
3. node scripts/sync-version.mjs
4. node scripts/sync-spec-tree.mjs
5. python3 linter-scripts/check-spec-cross-links.py --root spec --repo-root .
6. Commit all changed files together (one atomic commit)
```

**Order is enforced**: #2 must precede #3 (sync reads `package.json`); #3 must precede #4 (tree sync uses version metadata); #5 must run last (validates the synced state).

### X.3 What `sync-version.mjs` Computes

`version.json` fields are split into **computed** (do NOT hand-edit) vs **manual** (hand-edit allowed):

| Field | Source | Edit Mode |
|-------|--------|-----------|
| `version` | `package.json` | Auto |
| `name`, `description` | `package.json` | Manual (in `package.json`) |
| `git.commit`, `git.branch`, `git.updated` | `git rev-parse` | Auto |
| `stats.totalFiles`, `totalLines`, `totalFolders`, `totalBytes` | `spec/**` walk | Auto |
| `folders[].fileCount`, `lineCount`, `byteCount` | per-folder walk | Auto |
| `folders[].aiConfidence`, `ambiguity` | spec frontmatter | Auto (read from `## Scoring` table) |

**Hand-editing any auto field will be reverted on the next sync run** — and will trigger `Drift detected in version.json` in CI before that.

### X.4 Failure Modes

| Symptom | Cause | Fix |
|---------|-------|-----|
| `Drift detected in version.json` | Forgot to run sync after editing specs | Run #3 + #4, commit the diff |
| `git/updated fields ignored` (in CI diff) | Local sync used a different git commit than CI | Re-run sync on a clean working tree, then commit |
| `specTree.json` regenerated unexpectedly | Manual edit to `src/data/specTree.json` | Never hand-edit; always regenerate via #4 |
| Per-folder line count off by 1 | Trailing newline missing | Ensure every spec file ends with `\n` |

### X.5 Allowlist & Waiver Reference

When `check-spec-folder-refs.py` reports `Stale references found`, the fix is **not** always to fix the link. Three valid resolutions:

#### (a) Typo / rename — fix the path in markdown

Use this when the target folder was renamed or you mistyped the number.

#### (b) Sibling-repo reference — add to `[external]`

Use this when the path points to a folder in a different repo (e.g., `02-spec/15-domain-migration/` lives in a sibling repo).

```

# linter-scripts/spec-folder-refs.allowlist

[external]
02-spec/15-domain-migration/
02-spec/legacy-archive/
```

#### (c) Documentation-only — add to `[doc-only]`

Use this when the path is a redirect stub or a documentation pointer with no real content.

```
[doc-only]
02-spec/folder-structure-root.md
```

**Commit message convention**: When adding to either section, the commit message must include `allowlist: <reason>` so reviewers can audit.

### X.6 Cross-Link Allowlist — `spec-cross-links.allowlist`

Plain-text, one path per line, relative to repo root, no globs:

```
readme.md
LICENSE
.github/CODEOWNERS
changelog.md
```

Use this when a markdown file links to a non-spec file that exists outside the spec tree.

---

*Version & Sync Workflow appendix added — v3.3.0 — 2026-04-22*
