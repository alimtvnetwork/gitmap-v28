# Subtask 02: Canonical Asset Ingestion & version.json Manifest Generation

- **Parent Plan:** [363-gitmap-synchronization-and-version-control](../../pending/363-gitmap-synchronization-and-version-control.md)
- **Subtask ID:** `TASK-02`
- **Slug:** `gitmap-synchronization-and-version-control - Task 02`
- **Target Package:** `cli/cmdcg/`, `cli/cmdprompt/`, `pkg/manifest/`
- **Status:** `pending`
- **Assigned Worker:** Ingestion Subagent

---

## 1. Goal

Implement the canonical asset discovery and manifest generator engine. This subsystem crawls authoritative source repositories (`coding-guidelines`, `gitmap`), computes cross-platform stable SHA-256 hash digests for all managed assets, and generates a standardized, deterministic `syncManifest` embedded within `version.json`.

---

## 2. In Scope

1. **Asset Discovery Crawler (`asset_crawler.go`):**
   - Recursive scanning of designated asset trees:
     - `01-prompts/*.md` (Canonical AI prompts)
     - `.agents/skills/**` (Antigravity skill bundles and scripts)
     - `02-spec/**` (Core architectural specifications)
     - `03-ai-scripts/**` (Automated workflow and sync scripts)
   - Ignore patterns: temporary files (`.tmp`, `.bak`), build caches, and OS metadata files (`.DS_Store`, `Thumbs.db`).

2. **Cross-Platform Hash Computation (`asset_hasher.go`):**
   - Read file contents with line-ending normalization (stripping `\r` before hashing) to ensure identical SHA-256 hashes across Windows (`CRLF`) and Linux/macOS (`LF`).
   - Compute individual asset SHA-256 digests and root section Merkle hashes (`promptsHash`, `skillsHash`, `specsHash`, `scriptsHash`).

3. **`syncManifest` Generation & Serialization (`manifest_writer.go`):**
   - Construct typed `SyncManifest` struct:
     - `manifestVersion`: e.g. `"1.0.0"`
     - `generatedAt`: ISO 8601 UTC timestamp
     - `assetCount`: total managed assets
     - `sectionHashes`: map of section root digests
     - `assets`: map of relative path to SHA-256 hex digest
   - Inject or update the `"syncManifest"` key within `version.json` while preserving all existing root fields (`Version`, `Title`, `Authors`, etc.).
   - Sort keys deterministically to eliminate spurious git diffs.

4. **Manifest Reading & Validation (`manifest_reader.go`):**
   - Parse `version.json` from any target repository.
   - Validate schema version and verify root integrity.
   - Return strongly-typed `Result[SyncManifest]` containers.

---

## 3. Out of Scope

- File synchronization network transfer or disk overwrite (deferred to Task 03 and Task 04).
- Terminal CLI routing and flag handling (deferred to Task 04).
- Full builds or test suite execution.

---

## 4. Architectural Invariants

- **Zero Nested Ifs:** All validation routines must use early return guard clauses.
- **Positive Booleans:** `hasManifest`, `isValid`, `isNormalized`.
- **Function Limit:** Helper functions <= 15 lines max.
- **Leaf Error Sovereignty:** Return `*apperror.AppError` for file read or JSON marshal failures.

---

## 5. Verification Checklist

- [ ] Asset crawler recursively discovers all files in `01-prompts/`, `.agents/skills/`, `02-spec/`, and `03-ai-scripts/`.
- [ ] SHA-256 computation produces identical digests regardless of Windows (`CRLF`) or Linux (`LF`) line endings.
- [ ] `syncManifest` is cleanly written to `version.json` without modifying unrelated fields.
- [ ] Serialized JSON keys are deterministically ordered.
- [ ] Manifest reader correctly parses existing and missing manifest blocks.
