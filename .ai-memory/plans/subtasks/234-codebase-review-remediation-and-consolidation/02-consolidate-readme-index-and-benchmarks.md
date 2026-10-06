# Subtask 02: Consolidate README Index and Elevate Benchmarks

> **Parent Plan:** [234-codebase-review-remediation-and-consolidation.md](../../234-codebase-review-remediation-and-consolidation.md)  
> **Spec Reference:** [02-spec/21-app/234-codebase-review-remediation-and-consolidation/01-architecture-spec.md](../../../../02-spec/21-app/234-codebase-review-remediation-and-consolidation/01-architecture-spec.md)  
> **Master Ledger:** [02-spec/21-app/234-codebase-review-remediation-and-consolidation/00-master-audit-ledger.md](../../../../02-spec/21-app/234-codebase-review-remediation-and-consolidation/00-master-audit-ledger.md)  
> **Status:** `DONE`  
> **Target Areas:**  
> - Root README: `readme.md`  
> - Reading Guide: `what-to-read.md`  
> - Benchmark Data: `benchmark.md` -> `docs/benchmarks/benchmark.md`  
> - Benchmark Documentation Directory: `docs/benchmarks/`  

---

## 1. Technical Objective

Eliminate the byte-identical duplication between root `readme.md` and `what-to-read.md` (currently 2,522 lines / 169 KB each). Shrink root `readme.md` into an executive navigational index (~150–250 lines) that links out to the 15 command category documents in `docs/commands/`, engineering specifications in `02-spec/`, and performance benchmarks. Elevate root `benchmark.md` to `docs/benchmarks/benchmark.md`, link to it from `readme.md`, and safely remove `benchmark.md` from the root directory.

---

## 2. Architectural Context & Ingestion

1. **Current State:**
   - `readme.md` and `what-to-read.md` share identical SHA-256 hashes and line counts (2,522 lines each).
   - Any update to one file causes documentation drift unless both files are manually duplicated.
   - `benchmark.md` is orphaned at the repository root alongside existing benchmarks in `docs/benchmarks/search_benchmark.md`.
2. **Target State:**
   - `readme.md`: Authoritative front door for humans and AI agents. Outlines GitMap architecture, quickstart, core capabilities, and references sub-documents.
   - `what-to-read.md`: Navigational reading guide for autonomous agents, defining reading order: `.ai-memory/overview.md` -> `readme.md` -> `02-spec/` -> `docs/commands/`.
   - `docs/benchmarks/benchmark.md`: Authoritative polyglot performance comparison matrix (GitMap vs PowerShell vs Python) promoted into the canonical `docs/benchmarks/` directory.

---

## 3. Implementation Steps

### Step 1: Elevate `benchmark.md` to `docs/benchmarks/benchmark.md`
- Copy `benchmark.md` content into `docs/benchmarks/benchmark.md`.
- Verify synergy with `docs/benchmarks/search_benchmark.md`.
- Ensure all image and script references use relative paths (e.g., `03-ai-scripts/43-run-search-benchmarks.py`).
- Remove root `benchmark.md` from the root working tree and git tracking.

### Step 2: Author Executive Index `readme.md`
- Restructure `readme.md` to a concise index (~150–250 lines):
  - **Header & Badges:** GitMap title, version `v6.497.0`, author attribution (MD ALIM UL KARIM / RISEUP ASIA LLC).
  - **Core Value Proposition:** Fast multi-repo discovery, manifest preservation (`gitmap.json`), parallel cloner, split-DB indexing (`DH2D`), and AI-first workflows.
  - **Quickstart Guide:** Standard installation and basic commands (`gitmap scan`, `gitmap find`, `gitmap clone`).
  - **Command Categories Matrix:** Table linking to the 15 command category guides in `docs/commands/`:
    - Core Workflows (`docs/commands/01-core.md`)
    - Clone & Manifests (`docs/commands/02-clone.md`)
    - Search & Discovery (`docs/commands/03-search.md`)
    - Fleet & SSH Delegation (`docs/commands/04-fleet.md`)
    - Pipeline & Diagnostics (`docs/commands/05-pipeline.md`)
    - Database & State (`docs/commands/06-database.md`)
    - (and remaining category guides)
  - **Performance Benchmarks Section:** Summary table with direct link to `docs/benchmarks/benchmark.md`.
  - **Specifications Index:** Direct link to `02-spec/` and canonical coding guidelines.

### Step 3: Refactor `what-to-read.md`
- Replace duplicate 2,522-line file with a focused reading sequence guide:
  - Phase 1: Identity & Architecture (`version.json`, `readme.md`).
  - Phase 2: Memory & Context (`.ai-memory/overview.md`, `.ai-memory/plans/`).
  - Phase 3: Specifications (`02-spec/01-spec-authoring-guide/`, `02-spec/02-coding-guidelines/`).
  - Phase 4: Command Reference (`docs/commands/`).
- Clarify that root `readme.md` is the primary entry point, eliminating the requirement to maintain identical copies.

### Step 4: Verify Cross-Links and Paths
- Validate all markdown links are strictly relative (no `file:///` or Windows drive letters).
- Ensure positive boolean naming and clean formatting throughout.

---

## 4. Verification Protocol

```bash
# 1. Verify benchmark.md elevation
test -f docs/benchmarks/benchmark.md && echo "Elevated benchmark exists"
test ! -f benchmark.md && echo "Root benchmark removed"

# 2. Verify line count reduction on readme.md
wc -l readme.md

# 3. Verify what-to-read.md is no longer identical to readme.md
diff -q readme.md what-to-read.md || echo "Files successfully de-duplicated"

# 4. Confirm strictly relative links in readme.md
grep -i "D:" readme.md docs/benchmarks/benchmark.md || echo "No absolute paths found"
```

---

## 5. Acceptance Criteria

- [x] `docs/benchmarks/benchmark.md` created with complete benchmark data from root `benchmark.md`.
- [x] Root `benchmark.md` removed from root directory.
- [x] `readme.md` reduced to an executive index linking to `docs/commands/`, `02-spec/`, and `docs/benchmarks/benchmark.md`.
- [x] `what-to-read.md` converted into a concise agent reading sequence guide.
- [x] De-duplication achieved (SHA-256 hashes of `readme.md` and `what-to-read.md` are distinct).
- [x] All markdown links adhere strictly to relative git paths.
- [x] Zero absolute paths (`D:\...`) present in documentation files.
