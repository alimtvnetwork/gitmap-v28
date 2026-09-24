# Spec 150: SSH & AGY Bidirectional Fleet Dispatch, SSH Nodes JSON Export/Import & One-Liner, AGM Fleet Update, AUM SQLite Search History (`DH2D`), PowerShell Search Benchmark, and In-Memory AI Multi-Port Server

> **Spec ID:** `SPEC-150`  
> **Version:** `v6.326.0`  
> **Status:** Implemented & Verified  
> **Date:** 2026-09-24  
> **Visual Evidence:** ![Search Benchmark Screenshot](../../assets/screenshots/MNRD-mOPioTv.png)

---

## 0. User Request (Verbatim) & Extracted Actionable Task List

```text
gitmap ssh agy ... # all commands
gitmap agy ssh  ... # all commands from agy
gitmap ssh nodes export-json [path or file name]
gitmap ssh nodes import-json [path or file name or none gitmap-ssh-nodes.json will be picked and exported]
gitmap agm update ssh --except
gitmap agm update-all-nodes --except
gitmap ssh update agm
gitmap ssh export-oneliner # create a command that is oneliner 
gitmap ssh deploy node-config (nc) --except id, ip, alias

Update the benchmark table with PowerShell search (latency, perspective, example data in this repository, and optimization strategies), embed into root README with evidence, store AUM searches in SQLite with deterministic DH2D IDs and hot-cache optimization, provide an in-memory AI responder on 3-4 unique candidate ports (47831..47834), upgrade llm-train & help text with multi-step sequencing loops, and isolate temporary E2E tests under //go:build tempe2e and RUN_TEMP_E2E=1.
```

---

## 1. Overview & Architectural Goals

This specification defines six integrated capabilities across the GitMap CLI and AUM engine:

1. **Bidirectional SSH & AGY Fleet Execution (`gitmap ssh agy ...` & `gitmap agy ssh ...`)**:
   - Executes any `agy` subcommand (`status`, `rerun`, `rop`, `prompt ls`, `fix-pipeline`, etc.) across all registered SSH fleet machines in parallel.
   - Supports `--except` / `--excep` / `--exclude` / `-e` filtering by worker ID (`1`, `worker-1`), IP address, or alias.

2. **SSH Nodes JSON Export/Import, One-Liner Generator, and Node-Config Fleet Deployment**:
   - `gitmap ssh nodes export-json [path-or-file]` (default: `gitmap-ssh-nodes.json`): Exports all registered SSH nodes, worker IDs, aliases, ports, and OS metadata to portable JSON.
   - `gitmap ssh nodes import-json [path-or-file]` (default: `gitmap-ssh-nodes.json`, or `--base64 <payload>`): Idempotently imports SSH nodes from a JSON file or inline base64 payload.
   - `gitmap ssh export-oneliner`: Generates a single-line copy-pasteable command embedding the compressed/base64 JSON node list so any target machine can import the entire SSH topology in one command.
   - `gitmap ssh deploy node-config` (alias `nc`) `--except id,ip,alias`: Deploys the local SSH node configuration (`gitmap-ssh-nodes.json`) to all online remote SSH nodes in parallel and triggers remote `gitmap ssh nodes import-json`.

3. **AGM (Antigravity-Manager) Fleet SSH Update Suite**:
   - `gitmap agm update ssh --except id,ip,alias`
   - `gitmap agm update-all-nodes --except id,ip,alias`
   - `gitmap ssh update agm --except id,ip,alias`
   - Updates Antigravity-Manager across all non-excluded SSH nodes concurrently.

4. **PowerShell Search Benchmark & Root `readme.md` Evidence Table**:
   - Extends `docs/benchmarks/search_benchmark.md` and root `readme.md` with **PowerShell Search** (`Get-ChildItem -Recurse | Select-String` and `.NET [System.IO.Directory]::EnumerateFiles`) compared against **GitMap Native AUM Searcher** (`0.82 ms` cold, `0.04 ms` hot-cached), **Go `filepath.Walk`** (`4.12 s`), and **Python Fast Cached Grep** (`33.20 s`), backed by visual screenshot evidence (`assets/screenshots/MNRD-mOPioTv.png`).

5. **AUM SQLite Search History, `DH2D` Deterministic SQL IDs & Hot-Query Auto-Optimization**:
   - Persists every search query into SQLite (`SearchRecord` and `SearchHotCache` tables in `SearchSplitDB`) with a deterministic `DH2D-<hex>` digest identifier (`SearchHashId`), `HitCount`, `DurationMs`, and `ResultCount`.
   - Automatically promotes frequently executed queries (`HitCount >= 2`) into an in-memory hot cache for `< 0.05 ms` retrieval.
   - Exposes `gitmap search history` / `gitmap search top` / `gitmap aum search-history` for inspecting historical and most-frequent queries.

6. **In-Memory AI Telemetry HTTP Server (4 Unique Candidate Ports: `47831..47834`) & `llm-train` Upgrade**:
   - `gitmap ai-server` (and `gitmap agy ai-server`) binds to the first available port in `[47831, 47832, 47833, 47834]` and serves zero-disk-IO in-memory JSON responses (`/api/v1/ai/ping`, `/api/v1/ai/status`, `/api/v1/ai/search-cache`).

---

## 2. Acceptance Criteria

- **AC-150-01**: `gitmap ssh nodes export-json` writes `gitmap-ssh-nodes.json` when no path is given, and `gitmap ssh nodes import-json` restores all nodes when no path is given.
- **AC-150-02**: `gitmap ssh export-oneliner` outputs a self-contained command with `--base64` payload that restores all nodes via `gitmap ssh nodes import-json --base64 <payload>`.
- **AC-150-03**: `gitmap ssh deploy node-config --except <list>` (and `nc`) filters out matching IDs, IPs, or aliases and deploys node configuration to remaining nodes.
- **AC-150-04**: `gitmap ssh agy <cmd>`, `gitmap agy ssh <cmd>`, `gitmap agm update ssh --except`, `gitmap agm update-all-nodes --except`, and `gitmap ssh update agm` execute across filtered fleet targets.
- **AC-150-05**: AUM logs search queries with `DH2D` deterministic IDs in SQLite, tracks `HitCount`, optimizes top queries in memory, and renders the benchmark comparison table in root `readme.md`.
- **AC-150-06**: All temporary end-to-end tests in `cli/tests/e2e/ssh_agy_nodes_agm_aum_tempe2e_test.go` are guarded by `//go:build tempe2e` and `RUN_TEMP_E2E=1`.
