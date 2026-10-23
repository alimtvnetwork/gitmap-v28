# Subtask Execution Plan: Coding Guidelines, Skill Synchronization & Release Ceremony (ST-07)

## Document Metadata
- **Subtask ID**: ST-07 (`07-subtask-coding-guidelines-and-skills.md`)
- **Parent Plan**: `.ai-memory/plans/pending/314-gitmap-heatmap-and-url-handling-instructions.md`
- **Specification Reference**:
  - `02-spec/21-app/314-gitmap-heatmap-and-url-handling-instructions/00-overview.md`
  - `02-spec/21-app/314-gitmap-heatmap-and-url-handling-instructions/01-architecture-spec.md`
  - `02-spec/21-app/314-gitmap-heatmap-and-url-handling-instructions/02-component-spec.md`
  - `02-spec/21-app/314-gitmap-heatmap-and-url-handling-instructions/99-consistency-report.md`
- **Status**: `QUEUED / READY FOR PHASE 2 WAVE 3 & PHASE 3`
- **Assigned Worker**: Swarm Worker 1 (Wave 3) & Lead Orchestrator
- **Disjoint File Ownership**:
  - `02-spec/02-coding-guidelines/06-ai-optimization/10-anti-pattern-replacements.md` (new)
  - `02-spec/02-coding-guidelines/06-ai-optimization/11-heatmap-telemetry-and-url-caching.md` (new)
  - `.agents/skills/gitmap/skill.md` (update)
  - `cli/cmd/llm/llm_urls.go` (registry update)
  - Root changelog and version metadata

---

## 1. Objectives & Architectural Intent

To ensure all autonomous AI agents (Antigravity, Cline, Claude, Copilot, Muse) learn, respect, and invoke the newly introduced GitMap heatmap discovery and URL caching capabilities, repository documentation must be formally expanded:
1. **Author Anti-Pattern Replacements Guideline (`10-anti-pattern-replacements.md`)**:
   - Replaces manual directory grepping and raw curl calls with high-speed GitMap native equivalents (`gitmap aum search`, `gitmap url read`, `gitmap heatmap`).
   - Fixes references in `cli/cmd/llm/llm_urls.go` pointing to missing spec.
2. **Author Heatmap Telemetry & URL Caching Guideline (`11-heatmap-telemetry-and-url-caching.md`)**:
   - Establishes rules for telemetry logging, multi-AI database resolution, segmented downloading ($\ge 10$ splits, 5KB chunks, `aria2c` fallback), and Split-DB caching.
3. **Synchronize GitMap Skill (`.agents/skills/gitmap/skill.md`)**:
   - Updates CLI cheatsheets and autonomous agent instructions with `gitmap heatmap`, `gitmap rc heatmap`, and `gitmap url read`.
4. **SemVer Minor Release Ceremony**:
   - Validates relative paths across all files, bumps minor version (`gitmap release` / `minor-bump`), and performs an atomic GitMap commit (`gitmap cpf`).

---

## 2. Specification Deliverables Details

### 2.1 `10-anti-pattern-replacements.md`
- **Path**: `02-spec/02-coding-guidelines/06-ai-optimization/10-anti-pattern-replacements.md`
- **Core Topics**:
  - Anti-pattern 1: Using raw `curl`, `wget`, or custom HTTP scripts without local caching -> Replace with `gitmap url read <url>`.
  - Anti-pattern 2: Single-stream downloads for files $> 1\text{--}2\text{ MB}$ -> Replace with automatic 10-split segmented downloading or `aria2c`.
  - Anti-pattern 3: Grepping filesystem for past agent edits -> Replace with `gitmap heatmap --repo <alias>`.
  - Anti-pattern 4: Hardcoding absolute paths in scripts -> Replace with relative Git-root paths.

### 2.2 `11-heatmap-telemetry-and-url-caching.md`
- **Path**: `02-spec/02-coding-guidelines/06-ai-optimization/11-heatmap-telemetry-and-url-caching.md`
- **Core Topics**:
  - Multi-tier SQLite heatmap database layout:
    - Root DB: `EffectiveDataDir()/heatmap.db`
    - Per-repo DB: `repodb/<repo>.db`
    - Multi-AI companion DBs: `repo-cache/.heatmap/<seq>-<ai>/<repo>/heatmap.db`
  - URL Split-DB architecture:
    - Path: `repodb/url_cache.db` with WAL mode and single-writer sync mutex.
    - Sub-millisecond unbuffered stdout streaming via `output.Raw()`.
  - High-speed segmented download threshold:
    - Any payload $\ge 1\text{--}2\text{ MB}$ must use 10 parallel segments or 5KB chunks.
    - Check `aria2c` on host; fallback to native Go Range workers.
  - LLM training integration:
    - Ingesting heatmap files and live endpoints via `--heatmap <path|url>`.

### 2.3 `.agents/skills/gitmap/skill.md` Updates
- Add command documentation for:
  - `gitmap heatmap`: Discovery catalog, database inventory, alias mappings, recent tasks, and suggestion pills.
  - `gitmap rc heatmap`: Repository-filtered heatmap inspection.
  - `gitmap url read <url>`: Cached URL streaming with `--refresh`, `--no-cache`, and `--json`.
  - `gitmap llm train --heatmap <path|url>`: AI decision trace ingestion.

---

## 3. Implementation Steps & Micro-Batches

### Step 3.1: Author Specification Files
- Create `02-spec/02-coding-guidelines/06-ai-optimization/10-anti-pattern-replacements.md` with complete rule definitions, before/after code blocks, and AI directives.
- Create `02-spec/02-coding-guidelines/06-ai-optimization/11-heatmap-telemetry-and-url-caching.md` with comprehensive architectural diagrams and invariants.
- Update `cli/cmd/llm/llm_urls.go` to include the verified relative URL and path for guideline 11.

### Step 3.2: Update GitMap Skill Documentation
- Edit `.agents/skills/gitmap/skill.md` to append the new command reference blocks and guidelines into the cheatsheet sections.

### Step 3.3: Linter & Hygiene Validation
- Run relative path hygiene checks using `03-ai-scripts/07-relative-path-fixer.py` to confirm zero absolute paths or drive letters.
- Run `python 03-ai-scripts/check-enum-and-boolean.py` to confirm positive boolean naming.

### Step 3.4: SemVer Minor Release Ceremony (Lead Orchestrator)
- Execute minor version bump using repository release tool (`minor-bump` skill / `gitmap release`).
- Synchronize changelogs and update root version tags.
- Commit all changes using atomic GitMap command:
  ```bash
  gitmap cpf "heatmap-url-cache - implement heatmap discovery, url split-db caching, segmented downloader, and llm training integration"
  ```

---

## 4. Coding Guidelines Compliance Checklist

- [x] **Positive Booleans**: `hasCached`, `hasAria2`, `isAvailable`, `isExpired`.
- [x] **Strict Relative Paths**: Every markdown file link and path reference is relative.
- [x] **Worker Git Ban**: Lead orchestrator exclusively performs git commits via `gitmap cpf`.
- [x] **Zero Builds in Subagents**: No build or full test execution during subtask run.

---

## 5. Verification & Acceptance Criteria

1. **Specs Authored**: Both guideline files exist in `02-spec/02-coding-guidelines/06-ai-optimization/` and contain valid markdown with zero absolute paths.
2. **Skill Updated**: `.agents/skills/gitmap/skill.md` reflects `gitmap heatmap` and `gitmap url read`.
3. **Linter Pass**: `03-ai-scripts/01-check-relative-paths.py` reports zero violations.
4. **Clean Release**: Version incremented properly according to SemVer minor bump rules.
