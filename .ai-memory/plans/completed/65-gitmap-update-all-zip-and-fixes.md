# Completed Plan: 65-gitmap-update-all-zip-and-fixes

## Request Summary
Enhance GitMap fleet management, update distribution, caching lifecycle, prompts search primacy, pull output deduplication, and gitignore policy:
1. **Nodes Table Formatting**: Invert table display sequence in `gitmap nodes` so capabilities table is rendered first (TOP) and registered node status table is at the bottom. Prune redundant `ROLE` column from the capabilities table.
2. **Update All & Zip Distribution**: Implement `gitmap update all --include-others`, `gitmap update all zip --include-others`, and `gitmap update-all-zip --include-others` (alias: `uaz`). Support Anti-Gravity Manager (AGM) packaging, local zip packaging, SSH streaming (SCP equivalent) to remote temp directory, and triggering remote installer scripts (PowerShell for Windows, Shell for Linux/macOS).
3. **AUM Search Cache Lifecycle**: Add `CleanSearchHotCache()` to reset in-memory cache, purge SQLite hot cache (`.gitmap/data/search/sql.db`), and run vacuum on `gitmap search clean`. Add hermetic synthetic E2E tests strictly without proprietary brand names.
4. **Prompt Primacy**: Update V6 prompt and skill to ban raw unindexed search tools (`rg`, `ripgrep`, `grep`, `Select-String`, `findstr`), require GitMap search commands with concrete examples, require relative path disambiguation for repeated packages, and mandate structured failure sub-nodes.
5. **Repository Deduplication & Dual Solutions**: Deduplicate tracked repository records by canonical clean path. Disambiguate duplicate repository names with distinct paths. Provide structured sub-nodes with up to two alternative solutions (Option 1 vs Option 2).
6. **Gitignore Policy**: Ensure `.gitmap/` is NOT ignored globally. Enforce only `.gitmap/backup/` exclusion by default, allowing projects to legitimately track files inside `.gitmap/`. Normalize ignore patterns to eliminate duplicate entries.

## Verified Deliverables

### Task-01: Nodes Table Reordering and Role Column Pruning
- **Files Modified**: `cli/cmd/nodes_cmd.go`, `cli/cmd/nodes_cmd_test.go`
- **Evidence**:
  - `renderUnifiedNodesTable` invokes `renderCommandsMatrixTable` first and `renderPrimaryNodesTable` second.
  - `renderCommandsMatrixHeader` and `renderCommandsMatrixRow` omit the `ROLE` column.
  - `TestUnifiedNodes_RenderTable` asserts matrix table appears before primary table and does not contain `ROLE` in capabilities header.

### Task-02: Multi-Node Update Engine & Zip SCP Distribution
- **Files Modified**: `cli/cmd/rootutility.go`, `cli/cmdupdate/update_fleet.go`, `cli/cmdupdate/update_fleet_test.go`
- **Evidence**:
  - Registered `uaz`, `update-all-zip`, `updateallzip` commands and recognized `--zip`, `zip`, `--include-others`.
  - Added `IsZip` and `IncludeOthers` to `FleetUpdateOptions`.
  - Expanded fleet target resolution to all unified cluster and hosts nodes when `--include-others` is specified.
  - Implemented local zip packaging, SSH streaming (`StreamFileToRemote`) to remote temp directory, and remote script execution.
  - Unit tests in `update_fleet_test.go` verify flag parsing, routing, and target expansion.

### Task-03: AUM Search Cache Lifecycle & Safe Test Fixtures
- **Files Modified**: `cli/searcher/search_history_db.go`, `cli/cmd/search.go`, `cli/tests/e2e/ssh_agy_nodes_agm_aum_tempe2e_test.go`
- **Evidence**:
  - `CleanSearchHotCache` resets `aumHotCache`, deletes `SearchHotCache` and `SearchQueryLog` records from SQLite, and vacuums the database.
  - `gitmap search clean` executes `CleanSearchHotCache` with green confirmation output.
  - Added hermetic synthetic E2E test `TestTempE2E_AUMSearchCacheLifecycleAndClean` using generic mock names (`generic_alpha_repo`, `sample_service_node`, `test_fixture_beta`).

### Task-04: Canonical Prompts: Replace Raw rg with GitMap Commands
- **Files Modified**: `01-prompts/14-execute/13-execute-parent-task-with-n-steps-v6.md`, `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`
- **Evidence**:
  - Added explicit total ban on `rg` and `ripgrep` across Section 3, Section 12, Section 14, and Worker Brief.
  - Added concrete GitMap search command examples (`gitmap aum search`, `gitmap search`, `gitmap find`, `gitmap lf`, `gitmap cat`).
  - Added explicit instruction to disambiguate repeated packages by path and render structured failure sub-nodes with Option 1 vs Option 2.

### Task-05: Repository Deduplication & Sub-Node Solutions
- **Files Modified**: `cli/cmdpull/pull_efficient.go`, `cli/cmdpull/pull_efficient_render.go`, `cli/cmdpull/pull_remediation_hint.go`
- **Evidence**:
  - `deduplicateTrackedRecords` normalizes and deduplicates records by `filepath.Clean(strings.ToLower(r.AbsolutePath))`.
  - Disambiguated duplicate repository records by displaying distinct paths.
  - Implemented `ResolveStructuredRemediation` and updated `renderFailedGroup` and `renderDirtyGroup` to render dual options (Option 1 vs Option 2).

### Task-06: .gitignore Policy: Retain .gitmap, Exclude Only .gitmap/backup
- **Files Modified**: `cli/cmdpull/pull.go`, `cli/cmdignore/ignore_groups.go`, `cli/cmdignore/fix_ignore.go`
- **Evidence**:
  - `buildDefaultTrackedPathsToCheck` inspects exclusively `.gitmap/backup/` and resume JSON files.
  - `ImmutableDefaultRules` and `EnsureImmutableRules` enforce exclusively `.gitmap/backup/`.
  - `fix_ignore.go` deduplicates ignore patterns using normalized string trimming (`strings.Trim(trimmed, "/ \t\r\n")`).
