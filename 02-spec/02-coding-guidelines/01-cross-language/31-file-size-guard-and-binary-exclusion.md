# Rule R19: Repository File Size Guard, Binary Probing & Large JSON Exclusion Standards

**Status:** Active
**Scope:** Universal (All Polyglot Repositories, AI Scripts, Tools & Workflows)
**Strictness:** Mandatory Quality Gate

---

## 1. Overview & Core Motivation

Unrestricted file growth, unindexed binary files, and massive serialized data blobs degrade developer tooling, inflate Git repositories, and exhaust AI agent context windows:
1. **Search & Grep Degraded:** Standard text search tools (grep, ripgrep, AST parsers) slow down drastically when scanning multi-megabyte JSON trees or binary blobs.
2. **Context Window Exhaustion:** AI coding assistants (Antigravity, Claude Code, Cursor) that blindly ingest large JSON trees (like `specTree.json` or `test-inventory.json`) consume hundreds of thousands of tokens, triggering truncation and degraded reasoning.
3. **Repository Hygiene:** Uncommitted binaries, cache databases, and accidental media uploads bloat clone operations and repository storage.
4. **Documentation Inconsistencies:** Broken sequence numbering (e.g. `01`, `02`, `04` with missing `03`) and mismatched H1 titles confuse readers, automated indexers, and documentation site generators.

---

## 2. Core Architectural Standards

### 2.1 File Size Ceilings & De-bloat Rules

- **Standard Source Files:** Any source code (`.go`, `.ts`, `.py`, `.rs`, `.php`, `.cs`) or specification markdown (`.md`) MUST NOT exceed **500 KB** (warning at 250 KB).
- **Line Count Limits:** Files must observe modular decomposition limits (100–200 lines target, maximum 8–15 lines per function).
- **Allowed Waivers (`ALLOWED_LARGE_FILES`):** Explicitly pre-approved generated metadata trees (e.g. `src/data/specTree.json`, `.ai-memory/test-inventory.json`) are permitted only when cataloged in GitMap's built-in waiver list.

### 2.2 Large JSON Exclusion Mandate ("Out of the List")

- **Automatic Exclusion:** Any `.json` file exceeding **500 KB** is automatically excluded ("out of the list") from general repository search, file streaming, and AI coding prompts.
- **Explicit Targeted Access:** Large JSON files must be targeted individually via dedicated schema queries (`gitmap automation search -e .json --include-large-json`) rather than broad repository sweeps.

### 2.3 Memory-Safe Binary Probing & User Confirmation

- **Two-Pass Binary Detection:**
  1. **Extension Check:** Standard binary extensions (`.exe`, `.dll`, `.bin`, `.zip`, `.tar.gz`, `.png`, `.jpg`, `.ico`, `.wasm`, `.db`, `.sqlite`, `.pyc`, `.class`, etc.) are recognized instantly without disk reads.
  2. **8KB Chunk Null-Byte Probe:** Unrecognized extensions are inspected for null bytes (`0x00`) within the first 8,192 bytes, avoiding full file RAM allocations.
- **Interactive Exclusion Request:**
  When an un-excluded binary or oversized asset is discovered during an interactive session:
  ```
  [Binary File Detected] assets/icon.png (57 KB)
  Exclude this file from future searches? [y/N]:
  ```
  - Selecting `y` / `yes` records the exclusion into `.gitmap/data/<repo-slug>/automation/sql.db`.
  - Selecting `n` / `no` keeps the file permitted.
- **Non-Interactive / CI Execution:** In headless environments (CI pipelines, automation scripts), common binaries and large JSONs are excluded automatically by default without blocking standard input.

### 2.4 Numbered Documentation & Title Header Integrity

- **Contiguous Sequences:** Numbered markdown files (e.g. `01-intro.md`, `02-auth.md`, `03-api.md`) within any directory must begin at `00` or `01` and increment contiguously without numbering gaps or duplicate prefixes.
- **H1 Header Alignment:** The primary top-level H1 header MUST match the numeric prefix of the filename:
  - File: `04-newline-fixer.md` ➔ Title: `# 04 Newline Normalization`
  - File: `126-cargo-command.md` ➔ Title: `# 126 Cargo Command & Toolchain Runner`
- **Automated Fix Verification:** Run `gitmap automation sequence --fix` to verify and auto-align mismatched headers across the repository.

---

## 3. CLI Command Reference

| Action | Command | Description |
|---|---|---|
| **Audit File Sizes** | `gitmap automation guard [dir]` | Audits repository file sizes, common binaries, and large JSONs |
| **Custom Size Limit** | `gitmap automation guard --max-kb 250` | Audits files against a strict 250 KB threshold |
| **Audit Sequences** | `gitmap automation sequence [dir]` | Audits markdown sequence gaps and H1 header alignment |
| **Auto-Fix Headers** | `gitmap automation sequence --fix` | Automatically updates mismatched H1 headers to filename prefix |
| **List Exclusions** | `gitmap automation exclude list` | Displays persistent search and audit exclusions table |
| **Add Exclusion** | `gitmap automation exclude add <pattern> [reason]` | Persists a custom path or glob to the exclusions database |
| **Remove Exclusion** | `gitmap automation exclude rm <pattern>` | Removes an exclusion pattern from the database |
| **Clear Exclusions** | `gitmap automation exclude clear` | Purges all custom exclusions |

---

## 4. Acceptance Criteria & Quality Gates

1. **Zero Bloat Gate:** CI preflight runners must execute `gitmap automation guard` to reject any accidental oversized binary or uncompressed archive.
2. **Search Hygiene:** Multi-core search (`gitmap automation search <pattern>`) must filter out binaries and oversized JSONs by default.
3. **Documentation Hygiene:** Pull requests introducing new documentation or spec files must pass `gitmap automation sequence` without sequence gaps or title mismatches.
