# Subtask Plan 01: Scan Engine --rc, --separate, Auto-Merge & Sequential Allocator

- **Parent Task Slug:** `gitmap-repo-cache-scan-clone-and-ls-ecosystem`
- **Subtask ID:** `Task-01`
- **Owner:** Worker 01
- **Status:** PENDING
- **Owned Files:**
  - `cli/cmdscan/flags.go`
  - `cli/cmdscan/scan.go`
  - `cli/cmdscan/scan_rc_export.go`
  - `cli/cmdscan/scan_rc_merge.go`
  - `cli/cmdscan/scan_rc_test.go`

---

## 1. Technical Objective
Implement `--rc` / `--repo-cache` and `--separate` / `--seprate` / `--sep` / `--s` in `gitmap scan`:
1. Parse flags in `cli/cmdscan/flags.go`.
2. Intercept scan execution in `cli/cmdscan/scan.go` after record generation.
3. In `cli/cmdscan/scan_rc_export.go`:
   - Resolve `repo-cache` companion root via `store.OpenSpecialReposSplitDB()`.
   - If not `--separate`: dispatch to `MergeScanRecordsWithExistingManifest()`.
   - If `--separate`: scan `repo-cache/` for existing `XX-gitmap.json`, compute `nextSeq = max(Prefix) + 1`, and format `fmt.Sprintf("%02d-gitmap.json", nextSeq)`.
4. In `cli/cmdscan/scan_rc_merge.go`:
   - Target `repo-cache/01-gitmap/gitmap.json`.
   - If missing: write initial batch as indented JSON.
   - If existing: deserialize, normalize URLs, match existing repos, update metadata, append new repos, and atomically re-serialize.
5. In `cli/cmdscan/scan_rc_test.go`:
   - Test flag parsing (`--rc`, `--separate`, `--sep`, `--s`).
   - Test merge deduplication and field updating.
   - Test sequential allocation from `01-` to `02-`, `03-`.
