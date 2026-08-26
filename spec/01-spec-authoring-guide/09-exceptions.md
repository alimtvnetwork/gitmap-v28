# Exceptions & Special Cases

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

While the spec tree enforces strict conventions, there are documented exceptions. Every exception has a historical reason and is explicitly tracked here. New exceptions should be rare and must be documented in this file.

---

## Exception 1: Non-Contiguous Module Numbers

**Rule:** Module numbers should be sequential.  
**Exception:** Several number gaps exist in the spec tree.

| Gap | Reason |
|-----|--------|
| 22 (skipped) | Reserved; previously used, now skipped for historical compatibility |
| 24–27 (range) | Modules 23–26 were consolidated into `03-coding-guidelines/` subfolders. Numbers retired. |

**Example:**
```
spec/
├── 21-wp-seo-publish-cli/       # Module 21 exists
├── 23-how-app-issues-track/     # Module 23 exists (22 is skipped)
├── ...
├── 28-wp-plugin-development/    # Module 28 exists (24-27 retired)
```

**Policy:** Do NOT fill gaps retroactively. Use the next number after the highest existing module.

---

## Exception 2: readme.md Files

**Rule:** All files must use `{NN}-{name}.md` format with numeric prefixes.  
**Exception:** `readme.md` files are allowed WITHOUT numeric prefixes.

**Where:**
- `.lovable/memories/readme.md` — Project memory overview
- Any project root `readme.md`

**Reason:** `readme.md` is a universal convention recognized by Git hosting platforms (GitHub, GitLab) for auto-rendering.

---

## Exception 3: Non-Markdown Files

**Rule:** All spec files use `.md` extension.  
**Exception:** Data files may use other extensions.

| File | Location | Reason |
|------|----------|--------|
| `error-codes.json` | Various CLI modules | Machine-readable error code registry |
| `config.json` | Some modules | Configuration data |

**Naming rule for exceptions:** Still use kebab-case, but numeric prefix is optional.

```
✅ error-codes.json
✅ 01-config.json
❌ ErrorCodes.json       # No PascalCase
❌ error_codes.json      # No underscores
```

---

## Exception 4: Legacy Suggestion File Naming

**Rule:** Files in `.lovable/memories/` use kebab-case.  
**Exception:** Completed suggestion files in `.lovable/memories/suggestions/completed/` use legacy `C-XXX` prefixes.

**Example:**
```
.lovable/memories/suggestions/completed/
├── C-001-suggestion-title.md
├── C-002-suggestion-title.md
├── ...
├── C-080-suggestion-title.md
```

**Reason:** These files were created before the naming convention was standardized. Renaming 80 files would break historical references with no practical benefit.

---

## Exception 5: Dual-Purpose Prefix 02

**Rule:** Each numeric prefix maps to exactly one module.  
**Exception:** Prefix `02` is used by both `02-spec-management-software/` (the folder) and the root file `02-prefix-disambiguation.md`.

**Resolution:** The file `spec/02-prefix-disambiguation.md` documents this overlap. Both are retained for backward compatibility.

---

## Exception 6: Memory Folders Without Numeric Prefixes

**Rule:** Spec folders require numeric prefixes (`{NN}-{name}/`).  
**Exception:** Memory folders (`/.lovable/memories/`) use plain kebab-case without numeric prefixes.

**Example:**
```
.lovable/memories/
├── architecture/       # No numeric prefix
├── workflow/           # No numeric prefix
├── constraints/        # No numeric prefix
```

**Reason:** Memory folders serve a different purpose (institutional knowledge vs. formal specifications). The simpler naming reflects their more flexible, less formal nature.

---

## Exception 7: CLI Module Without Frontend Folder

**Rule:** CLI modules MUST have `01-backend/`, `02-frontend/`, `03-deploy/`.  
**Exception:** Headless CLIs may omit `02-frontend/` when the UI is a separate module.

**Example:**
```
34-time-log-cli/                  # Headless CLI
├── 01-backend/                  # ✅ Present
├── 03-deploy/                   # ✅ Present (note: still uses 03, NOT 02)
├── 00-overview.md
└── 99-consistency-report.md
                                 # 02-frontend/ intentionally omitted

35-time-log-ui/                   # Separate UI module
├── 01-frontend/                 # Frontend specs live here instead
└── ...
```

**Important:** When `02-frontend/` is omitted, the deploy folder STILL uses prefix `03` (not `02`). This maintains consistency with the standard CLI pattern.

---

## Exception 8: Extensions Folder in CLI Modules

**Rule:** CLI modules have exactly 3 subfolders (backend, frontend, deploy).  
**Exception:** Some CLI modules have additional subfolders for extensions or configs.

**Example:**
```
09-gsearch-cli/
├── 01-backend/
├── 02-frontend/
├── 03-deploy/
├── 04-extensions/               # Additional: plugin/extension specs
├── 06-configs/                  # Additional: configuration specs
└── ...
```

**Policy:** Additional subfolders are permitted when a CLI has significant feature areas beyond the core 3. They MUST follow the same naming convention and contain `00-overview.md`.

---

## Exception 9: The `suggestions.md` Legacy Tracker

**Rule:** All files in `.lovable/memories/` use kebab-case with optional numeric prefix.  
**Exception:** `suggestions.md` at the memory root is a legacy file without a numeric prefix.

**Reason:** Created before the convention was established. Maintained for backward compatibility.

---

## Exception 10: Coding Guidelines Nested Sub-Module

**Rule:** Maximum folder depth is 3 levels.  
**Exception:** `03-coding-guidelines/03-golang/01-enum-specification/` reaches the maximum depth with its own internal files.

```
spec/
└── 03-coding-guidelines/            # Level 1
    └── 03-golang/                    # Level 2
        └── 01-enum-specification/    # Level 3 (maximum)
            ├── 00-overview.md
            ├── 01-{file}.md
            └── ...
```

**Policy:** Do NOT create folders deeper than 3 levels. If content requires more granularity, use additional files instead of deeper nesting.

---

## Summary Table

| # | Exception | Scope | Permanent? |
|---|-----------|-------|------------|
| 1 | Non-contiguous module numbers | spec/ | ✅ Yes |
| 2 | readme.md without prefix | Project-wide | ✅ Yes |
| 3 | Non-markdown data files | spec/ modules | ✅ Yes |
| 4 | Legacy C-XXX suggestion names | memories/suggestions/ | ✅ Yes (frozen) |
| 5 | Dual-purpose prefix 02 | spec/ root | ✅ Yes |
| 6 | Memory folders without prefixes | .lovable/memories/ | ✅ Yes |
| 7 | CLI without frontend folder | CLI modules | ✅ Yes (case-by-case) |
| 8 | Extra CLI subfolders | CLI modules | ✅ Yes (case-by-case) |
| 9 | Legacy suggestions.md | memories/ root | ✅ Yes (frozen) |
| 10 | 3-level depth in coding guidelines | 03-coding-guidelines/ | ✅ Yes |

---

## Adding New Exceptions

1. **Document here first** — Before creating a non-conforming file/folder
2. **Include reason** — Why the standard cannot be followed
3. **Mark as permanent or temporary** — Temporary exceptions should have a remediation plan
4. **Update the consistency report** — Note the exception so it doesn't flag as an error

---

## Strictly Prohibited

> 🔴 Hard prohibitions for spec authoring. Numbered for stable reference —
> never renumber, only append. Mirrors `.lovable/strictly-avoid.md` →
> "Strictly Prohibited" and `.lovable/memory/avoid/`.
>
> 🧠 **Memory directive:** Once an AI agent reads this section, it MUST
> retain every numbered prohibition in its strictly-avoid working memory
> for the rest of the session and across future sessions. Never
> re-propose, re-suggest, or "double-check" any prohibited behaviour.
> Silence is the required response.

### Sequence 01 — `readme.txt` time suggestions

- **Never** suggest, recommend, hint at, ask about, or auto-add anything
  related to timestamps inside `readme.txt`.
- **Never** mention "git update time", commit time, last-modified time,
  scheduled regeneration, cron jobs, or any time-flavoured workflow in
  or about `readme.txt`.
- **Never** propose helper scripts, hooks, or CI tasks whose purpose is
  to keep `readme.txt` time fresh.
- **Only** modify `readme.txt` when the user explicitly asks in that
  exact turn, using the format the user dictates.
- **No follow-up offers** about time after writing `readme.txt`.
- Cross-reference: `.lovable/memory/avoid/02-no-time-suggestions-in-readme-txt.md`
  and `.lovable/strictly-avoid.md` → "Strictly Prohibited" → Sequence 01.
