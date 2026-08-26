# CLI Module Template

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

Every CLI tool in the project follows a **standardized 3-folder structure**. This document provides a step-by-step guide to create a new CLI module from scratch, including all required files and their content templates.

---

## Standard CLI Folder Structure

```
{NN}-{cli-name}/
├── 00-overview.md                    # Module overview and file inventory
│
├── 01-backend/                       # Backend specifications
│   ├── 00-overview.md               # Backend overview and file listing
│   ├── 01-architecture.md           # Core system design, lifecycle, patterns
│   ├── 02-{component}.md           # Feature-specific spec (numbered sequentially)
│   ├── 03-{component}.md
│   ├── ...
│   ├── {NN}-error-codes.md          # Error code registry for this module
│   ├── 97-acceptance-criteria.md    # Backend acceptance criteria (optional if parent has them)
│   ├── 98-changelog.md             # Backend changelog (optional)
│   └── 99-consistency-report.md    # Backend consistency report
│
├── 02-frontend/                     # Frontend specifications (if applicable)
│   ├── 00-overview.md
│   ├── 01-architecture.md
│   ├── 02-components.md
│   ├── 03-pages.md
│   └── 99-consistency-report.md
│
├── 03-deploy/                       # Deployment & operations
│   ├── 00-overview.md
│   ├── 01-build-pipeline.md
│   ├── 02-{platform}.md           # Platform-specific deployment
│   └── 99-consistency-report.md
│
├── 97-acceptance-criteria.md        # Consolidated acceptance criteria
├── 98-changelog.md                  # Module changelog
└── 99-consistency-report.md         # Module consistency report
```

---

## Step-by-Step: Creating a New CLI Module

### Step 1: Choose the Module Number

Look at the current `spec/00-overview.md` master index. Use the **next available number** after the highest existing module. Do NOT reuse or insert numbers in the middle of the sequence.

### Step 2: Create the Folder

```bash
mkdir -p spec/{NN}-{cli-name}/01-backend
mkdir -p spec/{NN}-{cli-name}/02-frontend
mkdir -p spec/{NN}-{cli-name}/03-deploy
```

### Step 3: Write 00-overview.md (Root)

This is the **first file** to create. It must include:

````markdown

# {CLI Name}

**Version:** 3.2.0  
**Status:** Draft  
**Updated:** YYYY-MM-DD  
**Language:** {Go | Rust | TypeScript}

---

## Overview

[Description of what this CLI does, who uses it, and why it exists]

---

## Core Features

1. **Feature A** — Brief description
2. **Feature B** — Brief description
3. **Feature C** — Brief description

---

## Folder Structure

```
{NN}-{cli-name}/
├── 00-overview.md
├── 01-backend/
│   ├── 00-overview.md
│   ├── 01-architecture.md
│   └── ...
├── 02-frontend/
│   └── ...
├── 03-deploy/
│   └── ...
├── 97-acceptance-criteria.md
└── 99-consistency-report.md
```

---

## Coding Standards

All code must follow:

<!-- Example template links — replace {NN} and {language} with actual values -->

- [Cross-Language Guidelines](../02-coding-guidelines/01-cross-language/00-overview.md)
- `{Language} Guidelines` → `../02-coding-guidelines/{NN}-{language}/00-overview.md`
- PascalCase for all database columns, JSON fields, and API payloads

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Coding Guidelines | `../02-coding-guidelines/00-overview.md` |
| Error Code Registry | `../03-error-manage/03-error-code-registry/00-overview.md` |
````

### Step 4: Write Backend Specs

Start with `01-backend/00-overview.md`, then add numbered component specs:

- **01-architecture.md** — Always first: daemon lifecycle, patterns, event bus
- **02-{first-component}.md** — Next major feature
- Continue sequentially...
- **{NN}-error-codes.md** — Error codes are typically the LAST content file

### Step 5: Write Frontend and Deploy Specs

Follow the same pattern. If the CLI has no frontend (headless tool), you may omit `02-frontend/` but document this in the overview.

### Step 6: Write Acceptance Criteria

Consolidate all testable requirements from backend, frontend, and deploy specs:

```markdown
| # | Criterion | Source |
|---|-----------|--------|
| AC-001 | CLI starts daemon in background mode within 2 seconds | `01-backend/01-architecture.md` |
```

### Step 7: Create Consistency Reports

Create `99-consistency-report.md` at each level:
- Module root
- `01-backend/`
- `02-frontend/` (if present)
- `03-deploy/` (if present)

### Step 8: Update Master Index

Add the new module to `spec/00-overview.md` in the appropriate layer table.

---

## Real-World Example: Time Log CLI

```
34-time-log-cli/
├── 00-overview.md                    # Overview, features, folder structure
├── 01-backend/
│   ├── 00-overview.md               # Backend file listing
│   ├── 01-architecture.md           # Daemon lifecycle, collector pattern
│   ├── 02-os-integration.md         # Platform-specific hooks (Win/Linux/Mac)
│   ├── 03-browser-tracking.md       # Tab detection, URL capture
│   ├── 04-screenshot-capture.md     # Screenshot engine, storage
│   ├── 05-database-schema.md        # SQLite schema (7 tables)
│   ├── 06-api-interface.md          # CLI commands + HTTP API
│   ├── 07-error-codes.md            # Error codes 15000–15499
│   ├── 08-file-path-extraction.md   # Window title parsing
│   ├── 09-remote-sync.md            # Offline queue, API sync
│   ├── 10-remote-settings.md        # Admin settings fetch
│   ├── 11-time-slice-productivity.md # Productivity scoring
│   ├── 97-acceptance-criteria.md    # 64 backend criteria
│   ├── 98-changelog.md
│   └── 99-consistency-report.md
├── 03-deploy/
│   ├── 00-overview.md
│   ├── 01-build-pipeline.md
│   ├── 02-windows-installer.md
│   ├── 03-linux-packaging.md
│   ├── 04-macos-packaging.md
│   ├── 05-auto-update.md
│   └── 99-consistency-report.md
├── 97-acceptance-criteria.md         # 124 consolidated criteria
└── 99-consistency-report.md
```

**Key observations:**
- No `02-frontend/` — This CLI is headless; the UI is a separate module (`35-time-log-ui`)
- Backend has 11 content files (01–11) + meta files (97, 98, 99)
- Deploy covers 3 platforms + auto-update
- Acceptance criteria exist at both backend and root levels

---

## Checklist: New CLI Module

- [ ] Module number selected (next available after highest)
- [ ] Root `00-overview.md` created with full metadata
- [ ] `01-backend/` folder with `00-overview.md` and `01-architecture.md`
- [ ] `02-frontend/` folder (or documented why omitted)
- [ ] `03-deploy/` folder with `00-overview.md`
- [ ] `97-acceptance-criteria.md` with numbered, testable criteria
- [ ] `99-consistency-report.md` at root and each subfolder
- [ ] `spec/00-overview.md` master index updated
- [ ] Cross-references added to related modules
- [ ] Error codes registered in `spec/03-error-manage/03-error-code-registry/`
