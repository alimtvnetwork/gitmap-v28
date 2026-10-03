# Spec Authoring Rules

**Version:** 3.2.0
**Updated:** 2026-04-16
**Status:** Active
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## Overview

This specification sets out foundational, non-negotiable rules for authoring specifications, documentation, and module folders across the entire repository.

---

## 1. Absolute Ban on Overview and Index Files

> 🔴 **NON-NEGOTIABLE RULE: Strict Ban on `00-overview.md` and `01-index.md`**
>
> Files named `00-overview.md`, `01-index.md`, `overview.md`, `index.md`, or any generic overview/index variants are **STRICTLY BANNED** across all folders in the repository.
>
> Every directory, module, sub-module, or specification group summary **MUST** exclusively use `readme.md` in strict lowercase.

### Core Requirements:
1. **Canonical Entry Point:** Every module, subfolder, and documentation bundle must use `readme.md` as its single entry point.
2. **Strict Lowercase:** Filenames must strictly be `readme.md` (never `README.md`, `Readme.md`, or `00-readme.md`).
3. **No Prefix on README:** `readme.md` does not carry numeric prefixes like `00-` or `01-`. It is always unadorned `readme.md`.
4. **Immediate Remediation:** Any legacy references to `00-overview.md` or `01-index.md` discovered in tools, linters, or existing documentation must be eliminated and updated to point to `readme.md`.

---

## 2. File and Folder Naming Standards

1. **Strict Lowercase Kebab-Case:** All folders and files must adhere to lowercase kebab-case naming. No uppercase characters, spaces, or underscores are permitted.
2. **Two-Digit Numeric Prefixes:** All specification files (except `readme.md`) and specification subfolders must begin with a two-digit zero-padded number followed by a hyphen (e.g., `01-*.md`, `02-*.md`).
3. **Sequential Numbering:** Files must follow a logical numeric sequence without duplicate prefixes within the same folder.
4. **Mandatory `.md` Extension:** All specification and documentation files must use the `.md` extension.

---

## 3. Strict Relative Git Paths

1. **Relative Paths Only:** Markdown links, images, cross-references, and code pointers must use paths relative to the repository root or the current file.
2. **Total Ban on Absolute Paths:** Absolute filesystem paths (such as `C:\...`, `/Users/...`, `/home/...`) or `file:///` URIs are strictly forbidden.

---

## 4. Required Structural Files

Every specification module must provide:
- `readme.md` — The module index and entry point.
- `99-consistency-report.md` — Structural verification and link integrity audit.
- `97-acceptance-criteria.md` — Measurable, testable criteria for implementation-ready specifications.

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Spec Authoring Guide Overview | `./readme.md` |
| Folder Structure | `./02-folder-structure.md` |
| Naming Conventions | `./03-naming-conventions.md` |
| Required Files | `./04-required-files.md` |

### Rule 13: Strict Prohibition of Overview and Index Files (TOTAL BAN)

- **Strict Ban on `00-overview.md` and `01-index.md`:** Files named `00-overview.md`, `01-index.md`, or any generic overview/index filenames are strictly prohibited.
- **Mandatory `readme.md` Standard:** Every directory, package, prompt collection, and specification folder MUST use strictly `readme.md` as its primary entry point and table of contents.
- **Single Root Changelog & Single Consistency Report Mandate:** Individual folders MUST NOT maintain separate `98-changelog.md` or `99-consistency-report.md` files. All version changes live in the root `changelog.md`, and all cross-specification metrics live in `02-spec/99-consistency-report.md`.
- **Mandatory Acceptance Criteria:** Every technical specification file MUST conclude with an exhaustive `## Acceptance Criteria` section specifying verifiable binary test conditions.
