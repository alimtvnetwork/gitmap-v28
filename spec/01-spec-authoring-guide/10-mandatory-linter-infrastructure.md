# Mandatory Linter Infrastructure

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**Type:** AI Instruction — Infrastructure Requirement  
**Priority:** Critical — Enforce Before Any Code Review or Validation

---

## Purpose

This file declares the **mandatory linter infrastructure** that must exist in every project repository that uses this spec system. If these folders or files are missing, the AI agent or contributor must **immediately flag the absence** and request that the files be uploaded or created before proceeding with code review, validation, or standards enforcement.

---

## AI Instruction

```
IMPORTANT — AI INSTRUCTION:

1. Before running any code review, validation, or standards check, verify that the
   mandatory linter infrastructure described below exists in the project repository.

2. If any required folder or file is MISSING:
   a. STOP — do not proceed with validation.
   b. Inform the user: "Required linter infrastructure is missing."
   c. List exactly which folders/files are absent.
   d. Ask the user to upload or create the missing files before continuing.

3. Do NOT attempt to generate linter scripts from scratch — they contain
   project-specific rules and must be provided by the project owner.

4. Do NOT skip linter checks because infrastructure is missing — the absence
   itself is the issue that must be resolved first.
```

---

## Required Folders

The following folders MUST exist at the **root of the project repository**:

| # | Folder | Purpose | Mandatory |
|---|--------|---------|-----------|
| 1 | `linter-scripts/` | Automated validation scripts — coding standards enforcement, link scanning, dashboard generation | ✅ Yes |

### Accepted Alternative Layout

Projects MAY use a unified `linters/` parent folder instead of placing `linter-scripts/` at the root:

```

# Option A — Root-level (default)

project-root/
├── linter-scripts/
│   ├── validate-guidelines.py
│   ├── validate-guidelines.go
│   ├── generate-dashboard-data.cjs
│   ├── check-axios-version.sh
│   ├── run.sh
│   └── run.ps1
└── ...

# Option B — Nested under linters/

project-root/
├── linters/
│   └── linter-scripts/
│       ├── validate-guidelines.py
│       ├── validate-guidelines.go
│       ├── generate-dashboard-data.cjs
│       ├── check-axios-version.sh
│       ├── run.sh
│       └── run.ps1
└── ...
```

Either layout is acceptable. The AI agent must check for **both** layouts before reporting the infrastructure as missing.

---

## Required Files Within `linter-scripts/`

The following files form the minimum linter infrastructure:

| # | File | Purpose | Mandatory |
|---|------|---------|-----------|
| 1 | `validate-guidelines.py` | Python-based coding standards validator — enforces CODE-RED and STYLE rules against `src/` | ✅ Yes |
| 2 | `validate-guidelines.go` | Go-based coding standards validator — alternative/companion to the Python version | ✅ Yes |
| 3 | `generate-dashboard-data.cjs` | Node.js script — validates cross-reference link integrity, generates system health dashboard | ✅ Yes |
| 4 | `check-axios-version.sh` | Shell script — verifies Axios dependency is pinned to a safe version | ✅ Yes |
| 5 | `run.sh` | Shell runner — executes the full validation suite (Linux/macOS) | ✅ Yes |
| 6 | `run.ps1` | PowerShell runner — executes the full validation suite (Windows) | ✅ Yes |
