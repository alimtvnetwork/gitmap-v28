# AI Fix Scripts Index & Tooling Guide

> /goal Master, discover, and execute the repository's suite of ultra-fast Python scripts for linting, path resolution, naming enforcement, version synchronization, local CI verification, polyglot discovery, artifact removal, and plan consolidation.
> /learn Read the script specifications below and run scripts via `python 03-ai-scripts/<script-name>.py`.

---

## 📋 AI Agent Pre-Flight Checklist

Follow this sequence before and during any repository modification task:

- [ ] **/learn** Inspect `02-shared-engine.py` to import centralized constants (`DEFAULT_ENCODING`, `LINE_SEPARATOR`, `TAB_CHAR`, `PATH_SEPARATOR`, `CURRENT_DIR`, `DEVICE_PATH_PREFIX`, `DOT_CHAR`, `UTF8_BOM_BYTES`, `CRLF_BYTES`, `NULL_BYTE`), enums (`RegexPatternType`, `ScanModeType`, `SeverityType`, `ExitCodeType`, `CacheKeyType`, `ArtifactCategoryType`), formatters (`format_comma_separated`, `format_keys`), lazy regex registry, and dual-platform locks.
- [ ] **/goal** Discover repository topology, languages (Go, Rust, Python, TypeScript, PHP, SQL), and subsystem roots using `18-codebase-topology-discoverer.py`.
- [ ] **/goal** Run rapid repo-wide file discovery using `11-fast-file-scanner.py` or instant cache lookup `<1ms`.
- [ ] **/goal** Rapidly read target files or explore folder contents using `17-fast-file-reader.py`.
- [ ] **/goal** Search multi-threaded regex patterns across files using `12-fast-cached-grep.py`.
- [ ] **/learn** Auto-fix whitespace, line endings, and boolean checks using `05-guideline-autofixer.py`.
- [ ] **/goal** Sanitize absolute filesystem paths and `file:///` URIs using `07-relative-path-fixer.py`.
- [ ] **/goal** Safely remove accidental binary blobs, pycache, or test artifacts using `19-artifact-remover.py`.
- [ ] **/goal** Consolidate, archive, and re-sequence Lovable plan files and subtasks using `20-plan-consolidator.py`.
- [ ] **/goal** Audit repository error management and AppError compliance using `python linter-scripts/check-error-management.py`.
- [ ] **/goal** Audit repository for nested if statements (nesting depth > 1) using `python linter-scripts/check-nested-ifs.py`.
- [ ] **/goal** Validate all 18 quality gates in parallel before submitting using `06-cicd-local-runner.py`.

---

## 🛠️ Master Script Catalog & Search Tags

| # | Script | Primary Purpose | Speed | Discovery Tags |
|:---:|---|---|:---:|---|
| **01** | `01-index.md` | Master index, script catalog, AI instructions, and tag registry | — | `docs`, `index`, `ai-instructions`, `catalog` |
| **02** | `02-shared-engine.py` | Shared engine: constants, regex registry with on-the-fly registration, locks, cache | ~2ms | `core`, `engine`, `constants`, `regex`, `locking`, `cache`, `enums`, `helpers` |
| **03** | `03-file-manipulator.py` | Mass lowercasing, sequence fixing, and UTF-8 LF normalization CLI | ~15ms | `rename`, `lowercase`, `sequence`, `encoding`, `cli` |
| **04** | `04-newline-fixer.py` | Fixes trailing whitespace and missing final newlines across folders | ~15ms | `newlines`, `whitespace`, `crlf`, `lf`, `formatting` |
| **05** | `05-guideline-autofixer.py` | Composite runner combining newline fixing and boolean naming checks | ~25ms | `autofix`, `composite`, `guidelines`, `booleans` |
| **06** | `06-cicd-local-runner.py` | Runs all 18 CI quality checks locally via `ThreadPoolExecutor` | ~35ms | `ci-cd`, `runner`, `parallel`, `quality-gates`, `test` |
| **07** | `07-relative-path-fixer.py` | Detects and fixes absolute paths / `file:///` URIs across folders | ~30ms | `paths`, `relative-paths`, `absolute-paths`, `sanitizer` |
| **08** | `08-naming-autofixer.py` | Enforces lowercase filenames, boolean conventions, and condition rules | ~20ms | `naming`, `booleans`, `is-prefix`, `has-prefix`, `linter` |
| **09** | `09-cli-help-auditor.py` | Validates CLI `--help` examples against actual implementations | ~25ms | `cli`, `help`, `cobra`, `commander`, `docstrings` |
| **10** | `10-encoding-normalizer.py` | Normalizes all files to strict UTF-8 with UNIX LF line endings | ~35ms | `encoding`, `utf-8`, `bom-stripping`, `unix-lf` |
| **11** | `11-fast-file-scanner.py` | High-speed repo file scanner (<15ms full scan, <1ms cache query) | ~14ms | `scanner`, `cache`, `indexing`, `file-list`, `discovery` |
| **12** | `12-fast-cached-grep.py` | Parallel regex matcher leveraging pre-warmed file cache | ~12ms | `grep`, `search`, `regex`, `parallel`, `content-search` |
| **13** | `13-file-size-guard.py` | Audits repository files for oversized binary blobs (>2MB) | ~10ms | `file-size`, `blob-guard`, `security`, `binary-check` |
| **14** | `14-version-sync-checker.py` | Verifies synchronization of `version.json`, `package.json`, `changelog.md` | ~5ms | `version`, `sync`, `changelog`, `package-json`, `release` |
| **15** | `15-sequence-and-title-auditor.py` | Audits and aligns numeric file sequence prefixes and `# H1` titles | ~20ms | `sequence`, `title`, `h1-headers`, `markdown-audit` |
| **16** | `16-installer-smoke-tester.py` | Generic installer smoke test validating script placeholders & hashes | ~8ms | `installer`, `smoke-test`, `install-sh`, `install-ps1` |
| **17** | `17-fast-file-reader.py` | AI agent fast file reader and folder explorer using `tmp/cache/` | <1ms | `reader`, `explorer`, `instant-read`, `ai-tool` |
| **18** | `18-codebase-topology-discoverer.py` | Universal polyglot codebase & topology discovery with TTL cache | ~15ms | `topology`, `discovery`, `polyglot`, `routing`, `cache-ttl`, `ai-tool` |
| **19** | `19-artifact-remover.py` | Safe interactive artifact remover with git index untracking (`git rm`) | ~10ms | `artifact-remover`, `cleanup`, `git-rm`, `pycache`, `safety-guard` |
| **20** | `20-plan-consolidator.py` | Fast Lovable plans & subtasks consolidator and index synchronizer | ~12ms | `plans`, `consolidator`, `subtasks`, `resequence`, `plan-cleanup` |
| **21** | `21-sequence-integrity-linter.py` | Verifies numeric sequences and headers across plans and subtasks | ~10ms | `sequence`, `linter`, `plans`, `integrity` |
| **22** | `22-doc-path-linter.py` | Lints markdown references and verifies documentation paths | ~15ms | `paths`, `docs`, `linter`, `relative-paths` |
| **23** | `23-coding-guideline-path-consolidator.py` | Consolidates coding guideline references to canonical specs | ~18ms | `guidelines`, `paths`, `consolidator` |
| **24** | `24-spec-path-migrator.py` | Migrates legacy spec references to updated paths | ~15ms | `spec`, `migration`, `paths` |
| **25** | `25-repo-migrator.py` | Repository-wide asset and structural migration utility | ~25ms | `migrator`, `repo`, `assets` |
| **26** | `26-go-code-formatter.py` | Cross-platform Go code formatter via gofmt with staged support | ~20ms | `go`, `gofmt`, `formatter`, `staged` |
| **27** | `27-misspell-auditor.py` | Audits and auto-fixes British to American English spelling | ~15ms | `spelling`, `misspell`, `us-english`, `autofix` |
| **28** | `28-go-preflight-ci.py` | Runs local Go test and golangci-lint preflight verification | ~35ms | `go`, `test`, `lint`, `preflight`, `ci-cd` |
| **29** | `29-release-orchestrator.py` | Automated release orchestrator, branch lifecycle manager, and tag pusher | ~10ms | `release`, `orchestrator`, `bump`, `branch`, `tag`, `lifecycle` |
| **30** | `30-db-struct-enum-generator.py` | Auto-generates type-safe column name enums from Go model structs | ~12ms | `db`, `models`, `enums`, `generator`, `codegen`, `columns` |
| **31** | `31-db-migration-runner.py` | Standalone external SQLite and schema migration runner with introspection | ~15ms | `db`, `migrations`, `sqlite`, `runner`, `schema` |
| **32** | `32-deep-consolidator.py` | Deep plans and subtasks consolidator preserving 100% technical detail, subtask ledgers, and verified outcomes | ~15ms | `plans`, `consolidator`, `milestone`, `subtasks`, `safety-backup` |
| **33** | `33-git-history-tracer-and-purger.py` | Traces deleted files in Git, pre-flight inspection with selective exclusion, workspace restoration, and deep history purging | ~20ms | `git`, `history`, `tracer`, `restore`, `purge`, `filter-repo`, `deleted-files` |
| **33** | `33-test-inventory-generator.py` | Generates `.lovable/test-inventory.json` and manages atomic file change tracking under lock | ~15ms | `test-inventory`, `tracking`, `locking`, `changes` |
| **34** | `34-schema-scanner.py` | Scans SQL table definitions for PascalCase, PK conventions, and affirmative booleans | ~15ms | `schema`, `database`, `sqlite`, `scanner`, `linter`, `erd` |
| **35** | `35-result-wrapper-auditor.py` | Audits Go functions returning multi-value map/slice error tuples for ResultMap and AppError compliance | ~15ms | `result`, `result-wrapper`, `result-map`, `apperror`, `linter`, `audit` |

---

## 🏛️ Core Shared Engine Architecture (`02-shared-engine.py`)

`02-shared-engine.py` is the single source of truth for all repository automation scripts.

### Centralized Constants & Configurations

```python
DEFAULT_ENCODING = "utf-8"
CURRENT_DIR = "."
DOT_CHAR = "."
EMPTY_STRING = ""
LINE_SEPARATOR = "\n"
CARRIAGE_RETURN = "\r"
CRLF_SEPARATOR = "\r\n"
TAB_CHAR = "\t"
PATH_SEPARATOR = "/"
WINDOWS_PATH_SEPARATOR = "\\"
DEVICE_PATH_PREFIX = "\\\\?\\"
COMMA_SPACE_SEPARATOR = ", "
UTF8_BOM_BYTES = b"\xef\xbb\xbf"
CRLF_BYTES = b"\r\n"
NULL_BYTE = b"\x00"
BINARY_PROBE_CHUNK_SIZE = 8192
DEFAULT_MAX_WORKERS = 4
```

### Key Architectural Components

1. **Lazy Regex Compilation with On-the-Fly Dynamic Registration:** Regex definitions are mapped to `None` on module load. When `RegexRegistry.get(pattern_type)` is called:
   - If present in the cache, returns the compiled regex.
   - If not yet compiled, compiles it on-demand with double-checked thread locking.
   - If an unregistered pattern is requested, it logs the event, auto-registers it in both dictionaries on-the-fly, compiles it, and returns the immutable compiled regex.
2. **Dual-Platform Cross-Process Locking:** POSIX kernel `fcntl.flock` on Linux/macOS (auto-cleans on SIGKILL/crash) and atomic `os.O_CREAT | os.O_EXCL` with 15s stale eviction on Windows.
3. **Two-Phase Caching Pipeline:** `stream_cached_files()` yields indexed files in <0.1ms; `stream_directory_files()` discovers new files while guarding against symlink recursion on Unix via inode tracking `(st_dev, st_ino)`.
4. **Fault-Tolerant I/O:** `read_file_safe()` and `write_file_lf()` preserve Unix file execution permissions (`st_mode`), normalize line endings, and prevent crash loops on concurrently deleted files.
5. **Shared Formatters:** `format_comma_separated(items)` and `format_keys(mapping)` for clean console reporting.

---

## 🎯 AI Operational & Style Rules

1. **Blank Line Before `if` Statements:** Always insert a blank line before an `if` condition when preceded by assignments, definitions, or statements (ensures visual separation between setup and branching).
2. **No Absolute Paths:** Always use relative paths from the git root. Never output `C:\...` or `file:///...`.
3. **Implicit Booleans:** Always evaluate positive booleans implicitly (`if is_valid:`, never `if is_valid == True:`).
4. **Prefix Boolean Variables & Functions:** Use `is_` or `has_` prefix for all boolean variables and return functions (`is_ready`, `has_match`, `is_success`, `has_failures`).
5. **Enums Format:** Python enums MUST use `PascalCase` class name ending in `Type`, `UPPER_CASE` members, and string values mirroring the member names.
6. **Quality Gates:** Before completing standard work sessions, execute `python 03-ai-scripts/06-cicd-local-runner.py --no-tests` and verify all quality gates pass (test execution is disabled unless explicitly commanded by repository owner). For release tasks or CI/CD fix tasks (`ci-cd-fix`), execute `python 03-ai-scripts/06-cicd-local-runner.py` (or `--run-tests`) to run all unit tests, linters, and quality gates properly.
7. **Atomic Change Tracking:** Append all modified files under lock via `python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`.

- [29-release-orchestrator.py](29-release-orchestrator.py): Orchestrates git release, bump version, branches, tags, and commits.
- [32-deep-consolidator.py](32-deep-consolidator.py): Deep plans and subtasks consolidator preserving 100% technical detail, subtask ledgers, and verified outcomes.
- [33-git-history-tracer-and-purger.py](33-git-history-tracer-and-purger.py): Traces deleted and active files in Git, selective pre-flight exclusion, workspace restoration, Recycle Bin file deletion (`--delete`), and deep history purging (`--purge`). Features automatic OS temp directory backups (`%TEMP%` / `/tmp`) with instant rollback instructions. Built-in presets: `--spec-25-audit` (`spec/21-app/25-app-spec-audit`), `--spec-audit` (`spec/19-main-worker-service/audit`), `--audit` (repo-wide `*audit*`), `--lovable`, `--lovable-subtasks`, `--lovable-md`, `--spec`, `--spec-md`, and positional root/folder scans.
- [33-test-inventory-generator.py](33-test-inventory-generator.py): Centralized test inventory manifest generator (`.lovable/test-inventory.json`) and atomic file change tracker with cross-platform mutex (`.lovable/temp/recent-file-changes.lock`).
- [35-result-wrapper-auditor.py](35-result-wrapper-auditor.py): Audits Go functions returning multi-value map/slice error tuples for ResultMap and AppError compliance.
- [36-param-struct-auditor.py](36-param-struct-auditor.py): Audits Go function signatures for argument reduction, dedicated parameter structs/DTOs, types.go declarations, and affirmative boolean naming.

---

## 🔍 Linters & AST Checks

- [linter-scripts/check-error-management.py](../linter-scripts/check-error-management.py): Linter enforcing zero bare panics, zero bare exits, zero exitWith aliases, zero naked apperror statements, zero disguised HandleError(nil), and zero swallowed database errors.
- [linter-scripts/check-nested-ifs.py](../linter-scripts/check-nested-ifs.py): AST linter enforcing zero nested if statements (nesting depth > 1), inverted early guard returns, and elimination of single-line collapsed if blocks across Go, TS/JS, Python, and PHP codebases.
