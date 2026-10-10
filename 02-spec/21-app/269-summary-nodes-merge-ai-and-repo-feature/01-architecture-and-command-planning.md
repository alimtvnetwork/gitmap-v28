# Specification: GitMap Summary, Nodes Delegation & Merge-AI Command Suite

**Document ID:** `02-spec/21-app/269-summary-nodes-merge-ai-and-repo-feature/01-architecture-and-command-planning.md`  
**Classification:** Core Application Specification (`02-spec/21-app`)  
**Milestone:** `269-summary-nodes-merge-ai-and-repo-feature`  
**Status:** `ratified`  
**Companion Documents:** [Universal 'Repo Feature' Resolver](repo-feature.md)  

---

## 1. Executive Summary & Problem Statement

Modern multi-repository engineering across local and remote SSH fleets suffers from three acute bottlenecks:
1. **Context Fragmentation:** Developers and AI agents cannot easily determine which repositories have active changes, what recent releases accomplished, or where file churn ("heated files") is concentrated without issuing hundreds of manual `git log` and `git status` calls.
2. **Disconnected CI/CD Telemetry:** Pipeline failures live in remote cloud dashboards (GitHub Actions, GitLab CI), decoupled from local repository summaries. AI agents must navigate browser sessions or raw log dumps to extract failure points.
3. **Redundant Fleet Queries:** In multi-machine fleets (Windows workstation, Linux dev boxes, CI nodes), identical clones of repositories exist across machines. Blindly querying all nodes duplicates network traffic, causes rate limiting, and introduces latency.
4. **Complex Multi-Repo Consolidation:** Merging features from multiple repositories into a unified target is error-prone. Blind rebases cause silent overwrite regressions, and manual staging requires tedious file alignment.

This specification addresses all four challenges through:
* **`gitmap summary` / `gitmap fs`:** Incremental, Split-DB cached release gists, heated file churn tracking, and workspace activity heatmaps.
* **`+pe` / `fspe` / `pe all`:** CI/CD error stack traces fused directly into terminal tree nodes for instant AI diagnosis.
* **`gitmap nodes ...`:** Intelligent fleet deduplication with local-machine precedence.
* **`repo feature`:** Universal target resolution engine for URLs, paths, and new repo slugs.
* **`gitmap merge-ai` (`ma`):** Deterministic multi-repo staging, `01_`/`02_` collision sequencing, `merge-ai-manifest.json` generation, and AI handoff instructions.

---

## 2. Command Architecture & Syntax Grammar

| Category | Canonical Command | Aliases / Variants | Default Parameters | Core Purpose |
| :--- | :--- | :--- | :--- | :--- |
| **Single Repo Summary** | `gitmap summary [$repo] [N]` | — | $N = 8$, target = `.` | Last $N$ releases summarized (≤200 words), heated files, Split-DB cached. |
| **Workspace Full Summary** | `gitmap full summary [N]` | `gitmap full status [N]`, `gitmap fs [N]` | $N = 3$ | TreeView of repos active in last 48h or dirty, with master commit command. |
| **Full Summary + CI/CD Errors** | `gitmap full summary+pe` | `gitmap full status+pe`, `gitmap fs+pe [--json]`, `gitmap fspe` | $N = 3$ | Fuses pipeline error stack traces into the Full Summary tree. |
| **Fleet Pipeline Errors** | `gitmap pe all [--json]` | `gitmap pe all --force-all` | 48h active filter | Fleet-wide CI/CD check. Concise green for passing; detailed errors for failing. |
| **Distributed Fleet Nodes** | `gitmap nodes fs` | `gitmap nodes full status`, `gitmap nodes full summary` | $N = 3$ | Distributed full summary with deduplication (local precedence). |
| **Distributed Single Repo** | `gitmap nodes summary $repo` | — | $N = 8$ | Queries cluster fleet for summary of a specific repo. |
| **Distributed Full + Errors** | `gitmap nodes fs+pe [--json]` | `gitmap nodes fspe [--json]` | $N = 3$ | Distributed summary with pipeline errors across fleet nodes. |
| **Distributed Errors All** | `gitmap nodes pe all [--json]` | `gitmap nodes pipe-error-all` | 48h active filter | Distributed pipeline errors across all nodes with deduplication. |
| **AI Merge Orchestrator** | `gitmap merge-ai <dest> <src...>` | `gitmap ma <dest> <src...>` | Single commit staging | Merges source repos into `<dest>` (repo feature), sequences collisions, emits manifest & instructions. |

---

## 3. Subsystem Architectural Specifications

### 3.1 Single-Repo Summary Engine (`gitmap summary`)
1. **Target Evaluation:** Resolves `$repo` using standard local path or remote URL. If omitted, defaults to the current working directory.
2. **Release Gist Generator:**
   * Extracts the last $N$ releases (tags sorted chronologically/SemVer, default $N = 8$).
   * Calculates commit log delta between tag boundaries.
   * Produces an executive gist constrained to **≤ 200 words** per release. Raw file lists are strictly suppressed from prose to preserve high-level semantic clarity.
3. **Heated Files Analysis:**
   * Computes churn metrics (modification count + lines added/deleted) per file across the release commit range.
   * Surfaces top 5 heated files and provides concise functional rationale for their modifications.
4. **Split SQLite Cache (`summary.db`):**
   * Stored in `<workspace>/.gitmap/summary.db`.
   * Keyed on `(repo_url, tag_commit_hash, head_commit_hash)`.
   * **Cache Hit:** If current git `HEAD` and release tag commit hashes match SQLite records, returns cached summary in `< 15ms`.
   * **Incremental Delta:** If commits exist beyond the cached tag, calculates only the delta to `HEAD`, updates `summary.db`, and outputs.

---

### 3.2 Workspace Full Summary Engine (`gitmap full summary` / `gitmap fs`)
1. **Activity Heatmap Filter (48-Hour Rule):**
   * Scans all repositories indexed in `gitmap.db`.
   * A repository is included in the output if and only if:
     * It has commit activity within the last **48 hours** (configurable via `Settings.ActivityWindowHours`).
     * **OR** it is in a **dirty state** (unstaged modifications, staged changes, or untracked files).
   * Dormant, clean repositories are excluded to ensure high signal-to-noise ratio.
2. **Tree View Representation:**
   * Emits ANSI TreeView rendering (`├──`, `└──`).
   * Displays the last $N$ releases (default $N = 3$) for qualifying active repositories.
3. **Dirty State Remediation & Master Command:**
   * Identifies dirty files per repository.
   * Prints a suggested commit command for each dirty repository.
   * **Single Master Command:** At the conclusion of the output, prints a single, consolidated shell command allowing the developer to sanitize/commit all pending repositories in one pass:
     ```bash
     gitmap sanitize-all --message "wip: save active progress across dirty repos"
     ```

---

### 3.3 CI/CD Pipeline Telemetry Integration (`+pe` / `fspe` / `pe all`)
1. **Pipeline Extraction:**
   * Interrogates `repodb/pipeline.db` for the most recent CI/CD pipeline run.
   * If passing, renders a clean green checkmark (`✓ CI/CD Passing`).
   * If failing, extracts failed step, job name, exit code, and the last 25 lines of the failure log.
2. **AI Diagnosis Readiness:**
   * Formats stack traces with explicit file paths and line numbers so downstream AI agents can immediately execute Root Cause Analysis (RCA) without leaving the terminal.
3. **`pe all` Optimization:**
   * Evaluates CI/CD state across repositories.
   * Respects the 48-hour activity window by default (avoiding API rate limits across large workspaces).
   * Employs the **Concise Green Rule**: Repositories with passing pipelines output a single green line; only failing pipelines expand with details.
   * `--force-all`: Bypasses the 48-hour window and queries every repository in the workspace.

---

### 3.4 Distributed Fleet Protocol (`gitmap nodes ...`)
1. **Deduplication Problem:** In a multi-node SSH cluster (e.g. Workstation, Node-1 Ubuntu, Node-2 Mac), repositories are frequently cloned redundantly across nodes. Querying every node for every repo leads to duplicate network overhead.
2. **Two-Phase Handshake with Local Precedence:**
   * **Phase 1 (Local Catalog):** Local master catalogs all locally hosted repository URLs and paths.
   * **Phase 2 (Remote Discovery):** Local master issues a lightweight SSH handshake to fleet nodes: `gitmap scan export --lean-manifest`, receiving a minimal JSON array of `{url, path, machine}`.
   * **Phase 3 (Partitioning):**
     * Local master retains all repositories it already hosts. Local status is evaluated locally without remote SSH calls.
     * Only repositories that do *not* exist on the local master are delegated to the specific remote nodes hosting them.
   * **Phase 4 (Unified Aggregation):** Merges local and remote outputs into a single TreeView or JSON payload.

---

### 3.5 AI Merge Orchestrator (`gitmap merge-ai` / `gitmap ma`)
1. **Target Resolution via "Repo Feature":**
   * Evaluates `<dest-target>` according to [repo-feature.md](repo-feature.md).
2. **Input Sources:**
   * Space-separated URLs: `url1 url2 url3`
   * Comma-separated URLs: `url1, url2, url3`
   * Text file: `file.txt` (one URL per line)
   * Config file: `config.json`
3. **Single-Commit Staging Strategy:**
   * GitMap prepares the destination staging workspace for a **single atomic commit**.
   * Multi-commit stacking/rebasing is intentionally disabled to prevent silent merge regressions.
4. **Collision Sequencing (`01_`, `02_`):**
   * Source repositories are inspected in parallel to determine chronological order (earliest to latest).
   * Base repo files are placed directly in the destination directory tree.
   * For subsequent repos:
     * Non-colliding files are placed directly into their destination paths.
     * Conflicting files are renamed with sequential prefixes:
       `01_<filename>` (Base/existing version)
       `02_<filename>` (Source 2 version)
       `03_<filename>` (Source 3 version)
5. **Artifact Generation:**
   * **`merge-ai-manifest.json`:** Comprehensive JSON manifest recording all source repos, commit hashes, branches, and the complete collision map.
   * **`instruction.md`:** Authoritative guidance document instructing the downstream AI agent to consolidate the prefixed files, preserve all business logic, flatten file names back to canonical names, and verify the build.

---

## 4. Split SQLite Database Schema (`summary.db`)

Located at `<workspace>/.gitmap/summary.db`:

```sql
PRAGMA journal_mode = WAL;
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS Repositories (
    repo_id INTEGER PRIMARY KEY AUTOINCREMENT,
    repo_url TEXT UNIQUE NOT NULL,
    canonical_slug TEXT NOT NULL,
    local_path TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ReleaseSummaries (
    summary_id INTEGER PRIMARY KEY AUTOINCREMENT,
    repo_id INTEGER NOT NULL REFERENCES Repositories(repo_id) ON DELETE CASCADE,
    tag_name TEXT NOT NULL,
    tag_commit_hash TEXT NOT NULL,
    release_date DATETIME NOT NULL,
    summary_gist TEXT NOT NULL,
    heated_files_json TEXT NOT NULL,
    word_count INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(repo_id, tag_name, tag_commit_hash)
);

CREATE TABLE IF NOT EXISTS RepoHeadStates (
    repo_id INTEGER PRIMARY KEY REFERENCES Repositories(repo_id) ON DELETE CASCADE,
    head_commit_hash TEXT NOT NULL,
    is_dirty BOOLEAN NOT NULL DEFAULT 0,
    dirty_files_count INTEGER NOT NULL DEFAULT 0,
    last_activity_at DATETIME NOT NULL,
    last_scanned_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE VIEW IF NOT EXISTS ViewRepoReleaseSummaries AS
SELECT 
    r.repo_url,
    r.canonical_slug,
    r.local_path,
    h.head_commit_hash,
    h.is_dirty,
    h.dirty_files_count,
    h.last_activity_at,
    s.tag_name,
    s.tag_commit_hash,
    s.release_date,
    s.summary_gist,
    s.heated_files_json
FROM Repositories r
JOIN RepoHeadStates h ON r.repo_id = h.repo_id
LEFT JOIN ReleaseSummaries s ON r.repo_id = s.repo_id
ORDER BY r.canonical_slug ASC, s.release_date DESC;
```

---

## 5. Master Implementation & AI Compliance Checklist

- [ ] **Phase 1: Split-DB Foundation (`summary.db`)**
  - [ ] Implement database migration applying tables `Repositories`, `ReleaseSummaries`, `RepoHeadStates`.
  - [ ] Implement `ViewRepoReleaseSummaries` view.
  - [ ] Implement cache hit validator (hash comparison returns in < 15ms).
  - [ ] Implement incremental delta calculation for newly added commits.
- [ ] **Phase 2: Single-Repo Summary (`gitmap summary`)**
  - [ ] Implement Cobra command `gitmap summary [path|url] [N]`.
  - [ ] Enforce default $N = 8$.
  - [ ] Implement release gist generator with ≤ 200 words constraint.
  - [ ] Implement heated files churn analysis (top 5 files by modification frequency).
  - [ ] Format terminal output with green checkmarks (`✓`).
- [ ] **Phase 3: Full Summary & Activity Heatmap (`gitmap fs`)**
  - [ ] Implement `gitmap full summary [N]`, `full status [N]`, and `fs [N]`.
  - [ ] Enforce default $N = 3$.
  - [ ] Implement 48-hour activity window filter.
  - [ ] Implement dirty working tree detector.
  - [ ] Render ANSI TreeView representation.
  - [ ] Output per-repo commit hints and a single master sanitize command at bottom.
- [ ] **Phase 4: CI/CD Pipeline Telemetry Integration (`+pe` / `fspe` / `pe all`)**
  - [ ] Implement `full summary+pe`, `fs+pe`, and `fspe`.
  - [ ] Connect pipeline diagnostics engine to `repodb/pipeline.db`.
  - [ ] Truncate failure stack traces to 25 lines with line numbers.
  - [ ] Implement `gitmap pe all` with 48h active filter and concise green status.
  - [ ] Implement `--force-all` and `--json` support.
- [ ] **Phase 5: Distributed Fleet Nodes Protocol (`gitmap nodes ...`)**
  - [ ] Implement Phase 1 local cataloging.
  - [ ] Implement Phase 2 remote discovery handshake.
  - [ ] Implement Phase 3 deduplication with local precedence (skip remote calls for locally hosted repos).
  - [ ] Implement parallel SSH delegation for unique remote repos.
  - [ ] Merge and render unified TreeView or JSON.
- [ ] **Phase 6: Universal "Repo Feature" Resolver**
  - [ ] Implement target resolution engine per `repo-feature.md`.
  - [ ] Support Remote URLs, local paths with/without `.git`, and bare slugs.
  - [ ] Leverage authenticated GitMap credentials for GitHub repo creation.
- [ ] **Phase 7: AI Merge Orchestrator (`gitmap merge-ai` / `ma`)**
  - [ ] Support space-separated URLs, comma-separated URLs, text files, and config JSON.
  - [ ] Implement parallel clone and chronological inspection.
  - [ ] Stage destination workspace for a single atomic commit.
  - [ ] Sequence colliding files into `01_<filename>`, `02_<filename>`.
  - [ ] Emit `merge-ai-manifest.json` and root `instruction.md` with action checklist.
- [ ] **Phase 8: End-to-End & Regression Verification**
  - [ ] Validate `< 25ms` cache hit latency on `gitmap summary`.
  - [ ] Verify zero redundant SSH queries on multi-node fleets with identical repo sets.
  - [ ] Verify that `merge-ai` creates sequenced files without silent file loss.
