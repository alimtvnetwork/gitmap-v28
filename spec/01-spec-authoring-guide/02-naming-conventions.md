# Naming Conventions

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

All files and folders in the `spec/` and `.lovable/memories/` trees follow strict naming conventions. These rules are non-negotiable and enforced by the health dashboard scanner.

---

## Folder Naming

### Format

```
{NN}-{kebab-case-name}/
```

### Rules

1. **Two-digit numeric prefix** — Always zero-padded (e.g., `01`, `08`, `33`)
2. **Kebab-case** — All lowercase, words separated by hyphens
3. **No spaces, underscores, or camelCase**
4. **Descriptive slug** — The name should clearly identify the module's purpose

### Examples

✅ Correct:
```
09-gsearch-cli/
01-backend/
03-coding-guidelines/
99-archive/
```

❌ Incorrect:
```
gsearch-cli/          # Missing numeric prefix
08_gsearch_cli/       # Underscores instead of hyphens
08-GSearchCli/        # PascalCase
8-gsearch-cli/        # Single-digit prefix
```

---

## File Naming

### Format

```
{NN}-{kebab-case-name}.md
```

### Rules

1. **Two-digit numeric prefix** — Sequential within the folder
2. **Kebab-case** — All lowercase, words separated by hyphens
3. **`.md` extension** — All spec files are Markdown
4. **No gaps in sequence** — Content files (01–89) should be contiguous

### Examples

✅ Correct:
```
00-overview.md
01-architecture.md
07-error-codes.md
97-acceptance-criteria.md
99-consistency-report.md
```

❌ Incorrect:
```
overview.md               # Missing numeric prefix
01-Architecture.md        # Capital letter
01_architecture.md        # Underscore
architecture-01.md        # Number at end
1-architecture.md         # Single-digit prefix
```

---

## Reserved File Prefixes

These numeric prefixes have fixed meanings across the entire spec tree:

| Prefix | File | Purpose | Required? |
|--------|------|---------|-----------|
| `00` | `00-overview.md` | Module index — lists all files, provides metadata | ✅ Always required |
| `96` | `96-ai-context.md` | AI-specific context notes | Optional |
| `97` | `97-acceptance-criteria.md` | Testable acceptance criteria | ✅ Recommended for all modules |
| `98` | `98-changelog.md` | Chronological change log | Optional (recommended for active modules) |
| `99` | `99-consistency-report.md` | Structural health report | ✅ Always required at top-level |

---

## Non-Markdown Files

Some modules contain non-markdown files. These follow relaxed naming:

| File Type | Naming Rule | Example |
|-----------|-------------|---------|
| JSON data files | Kebab-case, no numeric prefix | `error-codes.json` |
| Configuration | Kebab-case | `config.json` |
| Diagrams | Kebab-case with prefix | `01-architecture-diagram.svg` |

---

## Metadata Header

Every `.md` file MUST begin with a standardized metadata header:

```markdown

# Title of the Document

**Version:** X.Y.Z  
**Updated:** YYYY-MM-DD

---
```

### Rules

- **H1 title** — First line, exactly one per file
- **Version** — Semantic versioning (Major.Minor.Patch)
- **Updated** — ISO date format (YYYY-MM-DD)
- **Horizontal rule** — Separates metadata from content

### Optional Metadata Fields

```markdown
**Status:** Draft | Active | Deprecated  
**Language:** Go | Rust | TypeScript | PHP  
**Priority:** Critical | High | Medium | Low
```

---

## Source Code Files (Non-Spec)

Source code files in `src/` follow **different** conventions:

| Language | Convention | Example |
|----------|-----------|---------|
| TypeScript/React | PascalCase for components | `SpecFileViewer.tsx` |
| TypeScript | camelCase for utilities | `utils.ts` |
| Entry points | Lowercase | `main.ts`, `index.tsx` |

**Important:** Spec files and source code files use DIFFERENT naming conventions. Do not mix them.
