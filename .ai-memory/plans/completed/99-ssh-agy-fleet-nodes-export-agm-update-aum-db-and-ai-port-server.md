# Plan 99: SSH & AGY Bidirectional Fleet Dispatch, SSH Nodes JSON Export/Import & One-Liner, AGM Fleet Update, AUM SQLite Search History (`DH2D`), PowerShell Search Benchmark, and In-Memory AI Multi-Port Server

> **Status:** `COMPLETED`  
> **Release Version:** `v6.325.0`  
> **Spec Reference:** [`02-spec/21-app/150-ssh-agy-fleet-nodes-export-agm-update-aum-db-and-ai-port-server.md`](../../../02-spec/21-app/150-ssh-agy-fleet-nodes-export-agm-update-aum-db-and-ai-port-server.md)

---

## 1. Completed Subtasks

- [x] **SUBTASK-99-01**: Synchronized workspace (`git pull --rebase origin main`), persisted screenshot evidence to `assets/screenshots/MNRD-mOPioTv.png`, and authored Spec 150.
- [x] **SUBTASK-99-02**: Added **PowerShell Search** (`Get-ChildItem -Recurse | Select-String` and `.NET EnumerateFiles`) into `docs/benchmarks/search_benchmark.md` and root `readme.md` with repository sample timings and visual evidence. Implemented `cli/searcher/search_history_db.go` storing deterministic `DH2D-<HEX>` SQL IDs (`SearchHotCache`), `HitCount`, and auto-promoting frequent searches (`HitCount >= 2`) into an in-memory `<0.04ms` hot cache (`gitmap search history`).
- [x] **SUBTASK-99-03**: Wired bidirectional `gitmap ssh agy <subcmd>` and `gitmap agy ssh <subcmd>` with `--except` / `--excep` filtering by numeric ID (`1`), worker ID (`worker-1`), IP, or alias.
- [x] **SUBTASK-99-04**: Implemented `gitmap ssh nodes export-json [path]` (default `gitmap-ssh-nodes.json`), `gitmap ssh nodes import-json [path]` (default `gitmap-ssh-nodes.json` or `--base64`), `gitmap ssh export-oneliner`, and `gitmap ssh deploy node-config (nc) --except id,ip,alias` in `cli/cmdssh/ssh_nodes_export_import.go`.
- [x] **SUBTASK-99-05**: Implemented `gitmap agm update ssh --except`, `gitmap agm update-all-nodes --except`, and `gitmap ssh update agm --except`.
- [x] **SUBTASK-99-06**: Implemented `cli/cmd/ai_memory_server.go` (`gitmap ai-server`) binding to the first available port in `[47831, 47832, 47833, 47834]` and serving zero-disk-IO in-memory JSON responses (`/api/v1/ai/ping`, `/api/v1/ai/status`, `/api/v1/ai/search-cache`). Enhanced `llm-train` / `llm-docs` and CLI help text.
- [x] **SUBTASK-99-07**: Verified isolated temporary E2E suite `cli/tests/e2e/ssh_agy_nodes_agm_aum_tempe2e_test.go` (`//go:build tempe2e`, `RUN_TEMP_E2E=1`) and released `v6.325.0`.
