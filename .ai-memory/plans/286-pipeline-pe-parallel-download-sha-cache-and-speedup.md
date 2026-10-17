# Parent Plan: 286-pipeline-pe-parallel-download-sha-cache-and-speedup

## Description
Accelerate `gitmap pe` pipeline log downloads, parallelize failure processing, eliminate redundant GitHub queries on cache hits, fix stale TTL cache leaks across commit boundaries, and implement a persistent commit SHA cache (`PipelineShaCache`) in `repodb/pipeline.db` for sub-10ms response times.

## Status
- **Status:** COMPLETED
- **Created:** 2026-10-18
- **Phase:** 3 (Verification & Release)

## Architectural Decomposition
- **Subtask 01:** Parallel Log Download & Fast ANSI/Regex Extraction
  - Concurrently fetch failing job logs with bounded worker pool in `pipeline_query_logs.go`.
  - Bypass ANSI regex on plain lines and optimize failure scanning in `pipeline_extract_scan.go` and `pipeline_extract_clean.go`.
- **Subtask 02:** SQLite PipelineShaCache & TTL Leak Remediation
  - Add `PipelineShaCache` DDL and queries in `cli/pipelinedb/`.
  - Implement instant SHA cache hit rendering and cache serialization in `pipeline_sha_cache.go`.
  - Fix stale TTL cache leak when local HEAD differs from cached SHA in `pipeline_cache_eval.go`.
  - Add in-memory TTL caching for release tags and PR counts in `pipeline_logs_repo.go`.

## Subtasks
- [.ai-memory/plans/subtasks/286-pipeline-pe-parallel-download-sha-cache-and-speedup/subtask-01.md](.ai-memory/plans/subtasks/286-pipeline-pe-parallel-download-sha-cache-and-speedup/subtask-01.md)
- [.ai-memory/plans/subtasks/286-pipeline-pe-parallel-download-sha-cache-and-speedup/subtask-02.md](.ai-memory/plans/subtasks/286-pipeline-pe-parallel-download-sha-cache-and-speedup/subtask-02.md)
