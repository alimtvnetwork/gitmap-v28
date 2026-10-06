# Subtask 03: Fix version.json Identity & Author Cloning Architecture Spec

> **Parent Plan:** [234-codebase-review-remediation-and-consolidation.md](../../234-codebase-review-remediation-and-consolidation.md)  
> **Spec Reference:** [02-spec/21-app/234-codebase-review-remediation-and-consolidation/02-component-and-cli-spec.md](../../../../02-spec/21-app/234-codebase-review-remediation-and-consolidation/02-component-and-cli-spec.md)  
> **Master Ledger:** [02-spec/21-app/234-codebase-review-remediation-and-consolidation/00-master-audit-ledger.md](../../../../02-spec/21-app/234-codebase-review-remediation-and-consolidation/00-master-audit-ledger.md)  
> **Status:** `DONE`  
> **Target Areas:**  
> - Repository canonical identity: `version.json`  
> - Cloning architecture documentation: `docs/commands/cloning-architecture.md`  
> - Clone package mapping: `cli/cloner/`, `cli/clonefrom/`, `cli/clonenow/`, `cli/clonepick/`, `cli/clonenext/`  

---

## 1. Technical Objective

1. **Restore Canonical Repository Identity:** Update `version.json` to properly reflect the GitMap repository (`Title`: "GitMap", `RepoSlug`: "gitmap-v28", `RepoUrl`: "https://github.com/alimtvnetwork/gitmap-v28", `name`: "gitmap", `description`: "Autonomous developer companion and CLI for ultra-fast repository scanning, polyglot automation, cluster/SSH delegation, pipeline self-healing, and multi-tier AI agent task orchestration") instead of the stale "Coding Guidelines" placeholder. Preserve existing author metadata, folder specifications, prompts array, and canonical version `6.497.0`.
2. **Author Cloning Architecture Document:** Create `docs/commands/cloning-architecture.md` mapping all five clone subsystems (`cli/cloner`, `cli/clonefrom`, `cli/clonenow`, `cli/clonepick`, and `cli/clonenext`). Clearly explain each package's distinct data model (`model.ScanRecord`, `clonefrom.Plan`, `clonenow.Plan`, `clonepick.Plan`, `clonenext.ParsedRepo`), entrypoints, and specialized execution semantics, formally defending against monolithic single-file collapse.

---

## 2. Target File Inventory

### 2.1 Configuration Files to Update
1. `version.json` — Root single-source-of-truth version and repository identity descriptor.

### 2.2 Documentation Files to Author
1. `docs/commands/cloning-architecture.md` — Unified cloning architecture reference, mapping guide, and strategy manual.

### 2.3 Clone Packages to Map (Read-Only Architectural Analysis)
1. `cli/cloner/` — In-memory, scan-driven pipeline cloner consuming `model.ScanRecord`.
2. `cli/clonefrom/` — Plan-driven batch cloner consuming `clonefrom.Plan`.
3. `cli/clonenow/` — Round-trip manifest cloner consuming `clonenow.Plan`.
4. `cli/clonepick/` — Sparse-checkout and path-filtering cloner consuming `clonepick.Plan`.
5. `cli/clonenext/` — Next-version bumping and multi-repo remote cloner consuming `clonenext.ParsedRepo`.

---

## 3. Implementation Steps

### Step 1: Ingest & Inspect `version.json` Fields
- Review `version.json` at root to identify all mismatched identity attributes.
- Ensure `Version` is strictly maintained at `6.497.0` (zero unauthorized version increments).
- Verify author attribution (`MD ALIM UL KARIM`, `RISEUP ASIA LLC`) remains intact.

### Step 2: Update `version.json` Identity Attributes
- Modify the following top-level fields in `version.json`:
  ```json
  "Title": "GitMap",
  "RepoSlug": "gitmap-v28",
  "RepoUrl": "https://github.com/alimtvnetwork/gitmap-v28",
  "name": "gitmap",
  "Description": "Autonomous developer companion and CLI for ultra-fast repository scanning, polyglot automation, cluster/SSH delegation, pipeline self-healing, and multi-tier AI agent task orchestration."
  ```
- Retain all `folders`, `prompts`, and structural stats blocks.

### Step 3: Author `docs/commands/cloning-architecture.md`
- Detail the rationale for five separate modular packages:
  - `cli/cloner`: Scan-driven pipeline cloner (consumes `model.ScanRecord`, SQLite DB integration, multi-repo workspace tree discovery).
  - `cli/clonefrom`: User manifest-driven batch cloner (reads user JSON/CSV files, parses `clonefrom.Plan`, supports custom destination/branch/depth).
  - `cli/clonenow`: Scan artifact round-trip cloner (consumes `gitmap scan` output JSON/CSV/text, parses `clonenow.Plan`, preserves exact recorded folder paths, handles SSH/HTTPS mode toggle).
  - `cli/clonepick`: Interactive sparse-checkout / partial clone engine (parses `clonepick.Plan`, materializes specific paths, persists selections in SQLite DB for `--replay`).
  - `cli/clonenext`: Next-version bumping and multi-repo remote cloner (parses `clonenext.ParsedRepo`, resolves `vN` / `v++` / `v+1`, checks GitHub API, clones or updates repos).
- Document data model fields, CLI flag parity, progress rendering, and error handling for each engine.
- Include a Mermaid architecture diagram illustrating data flows from CLI input to disk.

### Step 4: Verification of JSON Validity & Documentation Integrity
- Validate `version.json` using a JSON linter/parser.
- Verify strictly relative paths throughout `docs/commands/cloning-architecture.md` (no absolute paths or absolute URI schemes).

---

## 4. Verification Protocol

```bash
# 1. Verify JSON syntax of version.json
python3 -c "import json; data=json.load(open('version.json')); assert data['Title'] == 'GitMap'; assert data['RepoSlug'] == 'gitmap-v28'; assert data['Version'] == '6.497.0'; print('version.json identity check passed!')"

# 2. Verify creation of cloning architecture document
test -f docs/commands/cloning-architecture.md && echo "cloning-architecture.md exists!"

# 3. Verify presence of all 5 clone engines in cloning-architecture.md
python3 -c "
content = open('docs/commands/cloning-architecture.md').read()
for pkg in ['cli/cloner', 'cli/clonefrom', 'cli/clonenow', 'cli/clonepick', 'cli/clonenext']:
    assert pkg in content, f'Missing package reference: {pkg}'
for model in ['model.ScanRecord', 'clonefrom.Plan', 'clonenow.Plan', 'clonepick.Plan', 'clonenext.ParsedRepo']:
    assert model in content, f'Missing model reference: {model}'
print('All 5 clone engines and models verified in cloning-architecture.md!')
"
```

---

## 5. Acceptance Criteria

- [x] `version.json` `Title` updated to `"GitMap"`.
- [x] `version.json` `RepoSlug` updated to `"gitmap-v28"`.
- [x] `version.json` `RepoUrl` updated to `"https://github.com/alimtvnetwork/gitmap-v28"`.
- [x] `version.json` `name` updated to `"gitmap"`.
- [x] `version.json` `Version` strictly maintained at `"6.497.0"`.
- [x] `docs/commands/cloning-architecture.md` authored with complete documentation of all 5 clone engines.
- [x] Data models (`model.ScanRecord`, `clonefrom.Plan`, `clonenow.Plan`, `clonepick.Plan`, `clonenext.ParsedRepo`) documented with field breakdowns.
- [x] Rationale against monolithic single-file collapse fully articulated.
- [x] Zero absolute paths used in documentation.
