# Plan 94: ORM Code Generator Integration, ID Naming & Error Standards, and Fast Cache Engine

**Title:** Database ORM Integration, Generator Upgrades, Tesla ID Naming Standards, Universal AppError Wrapping, Fast Cache Engine & Codebase Release  
**Status:** Pending  
**Parent Task Workflow:** Parent Task N-Step Continuous Loop & Multi-Agent Orchestration (Prompt v2.1.0)  
**Budget (N):** 150 steps  
**Target Directory:** `.lovable/plans/pending/`  

---

## 1. Root Cause & Problem Statement

1. **Disconnected Database ORM & Code Generation**:
   - `gitmap/dbengine` provides a rich, multi-dialect SQL query compiler and repository engine (`DbWrapper`, `Repository[T, F]`, `QueryBuilder[T, F]`, safe scan coercers).
   - However, currently only 4 models in `gitmap/pipelinedb` are connected. Tables in `gitmap/repodb` (`RepoFile`, `SearchCache`, `FileSequence`, `SequenceHistory`, `RepoScanLog`), `gitmap/db` (`ClusterNode`, `ClusterRun`, etc.), and key store entities are completely disconnected from the ORM, relying on ad-hoc raw SQL strings and handwritten scanners.
   - The generator `03-ai-scripts/30-db-struct-enum-generator.py` only generates read queries, does not parse `db:` tags, and contains hardcoded hacks.
2. **Tesla ID Naming & Schema Convention Violations**:
   - Spec `spec/04-database-conventions/01-naming-conventions.md` and `spec/01-naming` mandate PascalCase acronyms (`Id`, `Db`, `Url`, `Api`) and forbid all-caps clusters (`ID`, `DB`, `URL`, `API`).
   - Multiple structs in `gitmap/model/`, `gitmap/repodb/`, `gitmap/store/`, and `gitmap/db/` contain legacy `ID`, bare `Id` columns, or snake_case table/column names.
   - Lookup tables and "code values" (`TaskType`, `ProjectType`, cluster enums) violate Spec Rule 13's canonical `(Id, Code, Label, Description)` standard.
3. **Swallowed Errors & Missing AppError Envelopes**:
   - Over 15 instances of swallowed errors (`_ = `) exist in `pipelinedb/pipeline_split_ops.go` (dropping scan errors in stats and recent runs), `dbengine/wrapper.go` (swallowing transaction rollback errors), and `repodb/repo_db.go`.
   - `repodb` and `db` functions return standard Go `error` instead of domain-specific `*apperror.AppError`.
4. **Fast Cache & Search Bottlenecks (Go vs Python)**:
   - Python's `03-ai-scripts/` engine caches files in `tmp/cache/repo-file-cache.json` (<0.1ms startup), sniffs binary files with an 8KB null-byte probe, and prunes 20+ directories.
   - Go's `gitmap/indexer/walker.go` issues a synchronous SQLite `SELECT` query for *every individual file* inside `filepath.WalkDir`, lacks a binary sniffer (only checking file size), and only prunes `.git` and `node_modules`.
5. **Release Alignment**:
   - `03-ai-scripts/29-release-bumper.py` synchronizes manifests (`version.json`, `package.json`, `constants.go`, `changelog.md`), and `gitmap release` handles git tags, branches, and cross-compilation. They must be validated through `14-version-sync-checker.py` and local CI runner.

---

## 2. Task-Specific Rules & Constraints

1. **Strict PascalCase Acronym Standards**:
   - All acronyms must be PascalCase/camelCase: `Id`, `Db`, `Url`, `Api`. Zero all-caps `ID`, `DB`, `URL`, `API` in struct fields, method signatures, parameters, or database columns.
   - Primary keys must follow `<Entity>Id` format; foreign keys must follow `<TargetEntity>Id`.
2. **Zero Swallowed Errors & Universal AppError**:
   - Zero `_ = ` discards on database operations, row scans, or transaction rollbacks.
   - Every database method must return `(*T, *apperror.AppError)`, `([]T, *apperror.AppError)`, or `*apperror.AppError`.
3. **ORM Model Code Generation**:
   - Use `03-ai-scripts/30-db-struct-enum-generator.py` with `db:` tag support to generate `<Model>FieldType` enums, `Scan<Model>` row scanners, and `<Model>DbRepo` typed repositories wrapping `dbengine.Repository`.
4. **High-Performance File Indexing & Binary Sniffer**:
   - Batch-load file indexing timestamps in 1 query upfront before walking disk.
   - Add 8KB null-byte sniffer (`bytes.IndexByte(buf, 0) != -1`) before saving file contents into SQLite.
   - Port `EXCLUDE_DIRS` from `02-shared-engine.py` into Go indexer and searcher.
5. **Coding Guidelines Enforcement**:
   - Functions $\le 15$ lines, zero nested ifs, affirmative booleans only (`is*`, `has*`), Unix LF line endings.

---

## 3. Subtasks Ledger

- **Subtask 94.01**: Generator Upgrades & Tag-Driven Code Generation (`03-ai-scripts/30-db-struct-enum-generator.py`).
- **Subtask 94.02**: `repodb` & `db` ORM Connection, Entity Generation, and PascalCase `Id` Normalization (`gitmap/repodb/`, `gitmap/db/`, `gitmap/model/`).
- **Subtask 94.03**: Zero-Swallowed Errors & Universal `*apperror.AppError` Wrapping in Database Layers (`gitmap/pipelinedb/`, `gitmap/dbengine/`, `gitmap/repodb/`).
- **Subtask 94.04**: Go Fast Cache Optimization, Upfront Batch Indexing, 8KB Binary Sniffer & Directory Pruning (`gitmap/indexer/`, `gitmap/searcher/`, `gitmap/fsutil/`).
- **Subtask 94.05**: Codebase Release Synchronization, Pre-flight Verification & CI Quality Gates (`03-ai-scripts/`, `gitmap/release/`, local CI runner).
