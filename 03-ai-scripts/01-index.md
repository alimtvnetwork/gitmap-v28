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
| **29** | `29-release-orchestrator.py` | Autonomous release lifecycle: bump, commit, release branch, tag, and original branch revert | ~20ms | `release`, `orchestrator`, `git-branch`, `tag`, `semver`, `revert-branch` |
| **30** | `30-enum-generator.py` | Multi-file Go enum scaffolder (variant, vars, test, readme) with byte/int/string & CLI/JSON modes | ~10ms | `go`, `enum`, `scaffolder`, `base-enum`, `generator` |
| **31** | `31-md-gap-fixer.py` | Fixes multiple consecutive empty lines (3+ newlines) into exactly 1 empty line in markdown files | ~8ms | `markdown`, `formatting`, `newlines`, `gaps`, `cleanup` |
| **32** | `32-git-history-file-tracer.py` | Traces, pre-flights, restores, or permanently purges removed files from Git history | ~40ms | `git`, `history`, `tracer`, `restore`, `purge`, `filter-repo`, `preflight` |
| **33** | `33-test-inventory-generator.py` | Centralized test cataloging, duration tracking, and test resolution for releases | ~15ms | `test`, `inventory`, `cache`, `duration`, `release` |
| **34** | `34-purge-github-actions-artifacts.py` | Purges stored GitHub Actions artifacts to enforce Zero-Storage Mandate | ~25ms | `ci-cd`, `artifacts`, `github-actions`, `quota`, `cleanup` |
| **35** | `35-db-struct-enum-generator.py` | Inspects Go models, generates column enums and strongly-typed repository builders | ~20ms | `db`, `generator`, `scaffolder`, `repository`, `enums` |
| **36** | `36-multi-repo-folder-migrator.py` | Migrates child repositories from `.lovable`/`spec` to `.ai-memory`/`02-spec` | ~50ms | `migrator`, `multi-repo`, `structure`, `ai-memory` |
| **37** | `37-bump-version.py` | Repository-aware SemVer version bumper & manifest synchronizer | ~15ms | `version`, `bump`, `semver`, `sync`, `changelog` |

---

## 🔍 Detailed Script Specifications & CLI Reference

Every automation script in `03-ai-scripts/` is documented below with a collapsible `<details>` block detailing its architectural motivation (Why It Exists), runtime capabilities (What It Does), and command-line usage examples for check and fix/execution modes.

<details>
<summary><strong>01 — <code>01-index.md</code>: Script Catalog and AI Pre-Flight Guide</strong></summary>

#### Why It Exists

Serves as the central directory, documentation manifest, and AI discovery entry point for all automation and linting tools in `03-ai-scripts/`. It eliminates confusion around script availability, defines the mandatory pre-flight checklist for AI agents, and catalogs discovery tags for rapid tooling lookup.

#### What It Does

- Catalogs all 35 scripts in sequential numeric order with descriptions, execution speeds, and discovery tags.
- Outlines the mandatory AI pre-flight checklist guiding agents from topology discovery to local CI gate verification.
- Houses detailed collapsible `<details>` specifications, architectural summaries of `02-shared-engine.py`, and repository-wide operational coding rules.

#### CLI Usage & Examples

```bash
# View documentation via terminal
cat 03-ai-scripts/01-index.md

# Inspect using AI fast reader tool
python 03-ai-scripts/17-fast-file-reader.py --read-file 03-ai-scripts/01-index.md

# Audit markdown formatting and gap compliance
python 03-ai-scripts/31-md-gap-fixer.py 03-ai-scripts/01-index.md

# Auto-fix markdown gaps in-place
python 03-ai-scripts/31-md-gap-fixer.py 03-ai-scripts/01-index.md --fix
```

</details>

<details>
<summary><strong>02 — <code>02-shared-engine.py</code>: Centralized Constants, Regex Registry, Dual-Platform Locks, and Caching</strong></summary>

#### Why It Exists

Eliminates code duplication and architectural drift across repository tooling by centralizing foundational constants, shared types/enums, thread-safe lazy regex compilation, POSIX/Windows cross-process file locks, high-speed cached directory crawling, and atomic file I/O operations.

#### What It Does

- Centralizes repository constants (`DEFAULT_ENCODING`, `LINE_SEPARATOR`, `DEVICE_PATH_PREFIX`, `UTF8_BOM_BYTES`, `CRLF_BYTES`, etc.).
- Defines shared enums: `RegexPatternType`, `ScanModeType`, `SeverityType`, `ExitCodeType`, `CacheKeyType`, `ArtifactCategoryType`, `EncodingType`, `LanguageType`, `SubsystemType`.
- Implements thread-safe `RegexRegistry` with lazy compilation and on-the-fly dynamic auto-registration.
- Implements `CrossProcessLock` supporting POSIX `fcntl.flock` and Windows atomic `os.O_CREAT | os.O_EXCL` with stale lock eviction.
- Implements two-phase file walking: cache-first `stream_cached_files()` and disk-crawling `stream_directory_files()` with inode cycle detection.
- Provides atomic file writers (`write_file_lf`) preserving POSIX permissions and line-ending normalizers.

#### CLI Usage & Examples

```bash
# Verify module execution and self-test
python 03-ai-scripts/02-shared-engine.py

# Import engine utilities in Python tools
python -c "from importlib import import_module; engine = import_module('02-shared-engine'); print('Engine loaded:', engine.DEFAULT_ENCODING)"
```

</details>

<details>
<summary><strong>03 — <code>03-file-manipulator.py</code>: Mass Renaming, Lowercasing, and Sequence Normalization CLI</strong></summary>

#### Why It Exists

Enforces the strict lowercase filename mandate (Rule 4) and maintains clean sequential numbering across repository documentation, prompt files, and scripts without breaking Git history or corrupting file encodings.

#### What It Does

- Subcommand `lowercase`: Recursively discovers mixed-case or uppercase files and folders, safely renaming them via git-aware operations with case-insensitive collision guards.
- Subcommand `sequence`: Re-numbers prefixed files (e.g. `01-`, `02-`) monotonically, eliminating numeric gaps while preserving pinned mappings.
- Subcommand `encoding`: Standardizes encodings and line endings across target directory trees.

#### CLI Usage & Examples

```bash
# Preview lowercase renaming in target directory (dry-run check)
python 03-ai-scripts/03-file-manipulator.py lowercase 01-prompts --dry-run

# Apply lowercase renaming to all files under target directory
python 03-ai-scripts/03-file-manipulator.py lowercase 01-prompts

# Re-sequence numbered files while maintaining old order
python 03-ai-scripts/03-file-manipulator.py sequence 01-prompts --keep-old-order

# Re-sequence with pinned filenames
python 03-ai-scripts/03-file-manipulator.py sequence 01-prompts --pin "readme=00,intro=01"
```

</details>

<details>
<summary><strong>04 — <code>04-newline-fixer.py</code>: Trailing Whitespace and POSIX Final Newline Normalizer</strong></summary>

#### Why It Exists

Prevents noisy Git diffs, lint failures, and POSIX compiler warnings caused by trailing whitespace on code lines, carriage returns (`\r\n`), or missing final newlines at the end of files.

#### What It Does

- Recursively audits repository text files against whitespace conventions.
- Strips trailing spaces and tabs (`[ \t]+$`) from each line.
- Converts CRLF (`\r\n`) and CR (`\r`) to standard UNIX LF (`\n`).
- Guarantees exactly one trailing newline at EOF (no missing newline, no excessive trailing blank lines).

#### CLI Usage & Examples

```bash
# Audit whitespace and newlines across repository (check mode, exits 1 on violations)
python 03-ai-scripts/04-newline-fixer.py .

# Audit specific directory with custom extension filter
python 03-ai-scripts/04-newline-fixer.py 02-spec --ext .md,.json

# Auto-fix whitespace and newlines across repository in-place (fix mode)
python 03-ai-scripts/04-newline-fixer.py . --fix

# Auto-fix specific folder
python 03-ai-scripts/04-newline-fixer.py 01-prompts --fix
```

</details>

<details>
<summary><strong>05 — <code>05-guideline-autofixer.py</code>: Composite Autofixer for Newlines, Boolean Naming, and Guidelines</strong></summary>

#### Why It Exists

Provides a high-speed composite pre-commit runner that chains multiple guideline fixes (newline normalization and boolean convention auditing) into a single sub-25ms command.

#### What It Does

- Orchestrates `04-newline-fixer.py` to clean trailing whitespace and ensure EOF newlines.
- Invokes `08-naming-autofixer.py` to audit boolean conventions (explicit `== True` comparisons and banned boolean naming patterns).
- Reports consolidated diagnostics with file paths, line numbers, and error summaries.

#### CLI Usage & Examples

```bash
# Run composite audit across all repository files without modifying disk (check mode)
python 03-ai-scripts/05-guideline-autofixer.py . --check-only

# Check specific directory without modifying files
python 03-ai-scripts/05-guideline-autofixer.py 02-spec --check-only

# Automatically apply newline fixes and report boolean rule warnings (fix mode)
python 03-ai-scripts/05-guideline-autofixer.py .

# Target specific directory with custom extension filter
python 03-ai-scripts/05-guideline-autofixer.py 04-code --ext .go,.ts,.py
```

</details>

<details>
<summary><strong>06 — <code>06-cicd-local-runner.py</code>: High-Concurrency Local Test Suite Runner with Selective Failure Logging</strong></summary>

#### Why It Exists

Enables developers and AI agents to run all 31+ local CI quality gates concurrently in under 6 seconds before pushing code, preventing broken builds and wasted review cycles on remote GitHub Actions.

#### What It Does

- Concurrently executes all repository quality gates across multiple languages using `ThreadPoolExecutor`.
- Features intelligent log filtering: emits quiet green checkmarks on success; prints full stack traces and diffs on failure.
- Supports synchronous sequential mode (`--sync`) for deterministic step-by-step debugging.
- Supports job filtering (`--filter`) and machine-readable JSON exports (`--json`).

#### CLI Usage & Examples

```bash
# Run all quality gates concurrently (check mode, quiet on success, full logs on failure)
python 03-ai-scripts/06-cicd-local-runner.py

# Show detailed output for all gates (both passed and failed)
python 03-ai-scripts/06-cicd-local-runner.py --all

# Filter specific quality gate by name substring
python 03-ai-scripts/06-cicd-local-runner.py --filter "Go Base"

# Run sequentially in 1 worker thread for debugging
python 03-ai-scripts/06-cicd-local-runner.py --sync

# Output machine-readable JSON report
python 03-ai-scripts/06-cicd-local-runner.py --json --output tmp/ci-report.json
```

</details>

<details>
<summary><strong>07 — <code>07-relative-path-fixer.py</code>: Absolute Path and file:/// URI Detector and Sanitizer</strong></summary>

#### Why It Exists

Enforces Rule 5 (Strict Relative Git Paths Mandate and TOTAL BAN on absolute paths / `file:///` URIs) across all repository documentation, prompts, specs, and code comments.

#### What It Does

- Utilizes ultra-fast substring pre-filtering for path markers (`C:\`, `D:\`, `/home/`, `/Users/`, `file:///`, and Windows drive letters).
- Detects prohibited absolute filesystem paths and file URIs.
- When `--fix` is passed, automatically converts detected local filesystem paths into canonical relative Git repository paths.

#### CLI Usage & Examples

```bash
# Audit repository for absolute paths and file:/// URIs (check mode, exits 1 on violations)
python 03-ai-scripts/07-relative-path-fixer.py .

# Audit specific directory with custom extension filter
python 03-ai-scripts/07-relative-path-fixer.py 02-spec --ext .md,.json

# Automatically sanitize and replace absolute paths in-place (fix mode)
python 03-ai-scripts/07-relative-path-fixer.py . --fix

# Auto-fix specific documentation directory
python 03-ai-scripts/07-relative-path-fixer.py 01-prompts --fix
```

</details>

<details>
<summary><strong>08 — <code>08-naming-autofixer.py</code>: Lowercase Naming and Boolean Prefix Enforcement Linter</strong></summary>

#### Why It Exists

Enforces Rule 1 (Boolean Principles: TOTAL BAN on explicit `== True` checks) and repository naming standards across Go, TypeScript, Python, and other languages.

#### What It Does

- Performs high-speed substring pre-filtering for `true` / `True` before running regular expressions.
- Flags anti-patterns such as `if condition == true`, `if condition === True`, and negative boolean variable prefixes (`isNotReady`, `hasNoAccess`).
- Audits lowercase filename compliance across the repository.

#### CLI Usage & Examples

```bash
# Audit boolean comparisons and naming conventions across repository (check mode)
python 03-ai-scripts/08-naming-autofixer.py .

# Audit specific code directory with language extension filters
python 03-ai-scripts/08-naming-autofixer.py 04-code --ext .go,.ts,.py

# Audit target source folder
python 03-ai-scripts/08-naming-autofixer.py src/ --ext .ts,.tsx
```

</details>

<details>
<summary><strong>09 — <code>09-cli-help-auditor.py</code>: Command-Line Documentation and Help Verification Tool</strong></summary>

#### Why It Exists

Ensures that CLI `--help` examples, flags, and command descriptions match actual code implementations across Go (Cobra), TypeScript (Commander), Python (Argparse/Click), and PHP (Symfony), preventing stale or misleading CLI documentation.

#### What It Does

- Parses CLI command definitions, subcommand trees, flag declarations, and help annotations.
- Detects undocumented subcommands, broken usage examples, or missing help descriptions.
- Supports parallel multi-worker execution, strict failure mode (`--strict`), and structured JSON reports.

#### CLI Usage & Examples

```bash
# Audit CLI commands and help text across repository (check mode)
python 03-ai-scripts/09-cli-help-auditor.py .

# Strict audit failing with exit code 1 on warnings
python 03-ai-scripts/09-cli-help-auditor.py . --strict

# Audit specific CLI directory with custom worker count
python 03-ai-scripts/09-cli-help-auditor.py 04-code/golang/cmd --ext .go -w 4

# Export audit summary to JSON file
python 03-ai-scripts/09-cli-help-auditor.py . --json --output tmp/cli-audit.json
```

</details>

<details>
<summary><strong>10 — <code>10-encoding-normalizer.py</code>: UTF-8 and LF Line Ending Normalization Utility</strong></summary>

#### Why It Exists

Eliminates UTF-8 Byte Order Marks (BOM) and Windows CRLF line endings that cause compilation errors, bash script execution failures (`\r: command not found`), and noisy Git diffs.

#### What It Does

- Audits text files for UTF-8 BOM (`\xef\xbb\xbf`) and carriage returns (`\r`).
- When `--fix` is passed, strips BOM markers and normalizes all line breaks to strict UNIX LF (`\n`) in-place while preserving file permissions.

#### CLI Usage & Examples

```bash
# Audit repository files for BOM markers and CRLF line endings (check mode)
python 03-ai-scripts/10-encoding-normalizer.py .

# Audit specific directory with extension filter
python 03-ai-scripts/10-encoding-normalizer.py 01-prompts --ext .md

# Strip BOM and normalize line endings to LF across repository (fix mode)
python 03-ai-scripts/10-encoding-normalizer.py . --fix

# Auto-fix specific directory
python 03-ai-scripts/10-encoding-normalizer.py 02-spec --ext .md,.json --fix
```

</details>

<details>
<summary><strong>11 — <code>11-fast-file-scanner.py</code>: High-Performance Filesystem Scanner with Memory Cache</strong></summary>

#### Why It Exists

Provides instantaneous (<1ms from cache, <15ms cold scan) file discovery and inventory across the repository, replacing slow recursive shell walks for AI agents and CI scripts.

#### What It Does

- Recursively indexes repository files while respecting ignore patterns (`.git`, `node_modules`, `tmp/`).
- Maintains a pre-warmed file index cache in `tmp/cache/repo-file-cache.json`.
- Supports filtering by programming language alias (`--lang go,ts,py`), file extension (`--ext`), or substring (`--search`).
- Supports CI verification mode (`--check`) to validate index integrity.

#### CLI Usage & Examples

```bash
# Scan and index repository files to cache
python 03-ai-scripts/11-fast-file-scanner.py

# CI validation check verifying index integrity (check mode)
python 03-ai-scripts/11-fast-file-scanner.py --check

# Query existing cached index without walking the filesystem (<1ms)
python 03-ai-scripts/11-fast-file-scanner.py --query-cache --lang go,ts

# Filter files matching substring with stats breakdown
python 03-ai-scripts/11-fast-file-scanner.py --search "guideline" --stats
```

</details>

<details>
<summary><strong>12 — <code>12-fast-cached-grep.py</code>: Multi-Threaded Pattern Matcher Over Cached Repository Files</strong></summary>

#### Why It Exists

Delivers sub-15ms multi-threaded regex and pattern matching across the repository by leveraging the pre-warmed file cache from `11-fast-file-scanner.py`.

#### What It Does

- Reads indexed file paths from `tmp/cache/repo-file-cache.json` (or falls back to live disk scan if cache is missing).
- Uses `ThreadPoolExecutor` to search files in parallel without shell spawning overhead.
- Supports regex patterns (`--regex`), case sensitivity (`--case-sensitive`), language aliases, and extension filters.

#### CLI Usage & Examples

```bash
# Fast search for pattern across all repository files
python 03-ai-scripts/12-fast-cached-grep.py --pattern "appfault.AppError"

# Search with regular expression in Go files only
python 03-ai-scripts/12-fast-cached-grep.py --pattern "func (.*) Validate" --regex --lang go

# Case-sensitive search limited to markdown documentation
python 03-ai-scripts/12-fast-cached-grep.py --pattern "CODE RED" --case-sensitive --ext .md

# Bounded search with max result limit
python 03-ai-scripts/12-fast-cached-grep.py --pattern "TODO" --max 20
```

</details>

<details>
<summary><strong>13 — <code>13-file-size-guard.py</code>: Binary Blob and Oversized File Gatekeeper</strong></summary>

#### Why It Exists

Acts as a repository gatekeeper preventing accidental commits of oversized binary assets, database dumps, build artifacts, or huge log files that bloat Git history.

#### What It Does

- Scans tracked repository files against a configurable size limit (default 2048 KB / 2MB).
- Prunes ignored paths (`.git`, `node_modules`, `tmp/`).
- Reports files exceeding the threshold with exact sizes, exiting with code 1 if violations are found.

#### CLI Usage & Examples

```bash
# Audit repository for files exceeding 2MB default threshold (check mode)
python 03-ai-scripts/13-file-size-guard.py

# Audit with custom file size threshold (e.g. 500 KB)
python 03-ai-scripts/13-file-size-guard.py --max-kb 500

# Audit specific subfolder with custom threshold
python 03-ai-scripts/13-file-size-guard.py --path 02-spec --max-kb 1000
```

</details>

<details>
<summary><strong>14 — <code>14-version-sync-checker.py</code>: Semantic Version Synchronization Auditor</strong></summary>

#### Why It Exists

Prevents release desynchronization by verifying that version numbers match 100% identically across `version.json`, `package.json`, `changelog.md`, and any language manifest files.

#### What It Does

- Reads semantic version strings from root metadata files (`version.json`, `package.json`).
- Checks `changelog.md` to ensure the latest release header matches the current version.
- Returns exit code 0 when versions are in sync, or exit code 1 with mismatch details.

#### CLI Usage & Examples

```bash
# Validate version synchronization across root metadata files (check mode)
python 03-ai-scripts/14-version-sync-checker.py

# Validate version files in specific project path
python 03-ai-scripts/14-version-sync-checker.py .
```

</details>

<details>
<summary><strong>15 — <code>15-sequence-and-title-auditor.py</code>: Numeric File Sequence and H1 Title Consistency Checker</strong></summary>

#### Why It Exists

Ensures that numbered markdown files follow clean, continuous numeric sequences without gaps, and verifies that the primary `# H1` header inside each file matches its numeric prefix.

#### What It Does

- Inspects filename numeric prefixes (`01-`, `02-`, ...) in target documentation directories.
- Verifies that file numbers are strictly monotonic without missing steps or duplicates.
- Parses the `# H1` heading and checks consistency against the file prefix.
- When `--fix` is passed, automatically synchronizes the H1 title numbers with the filename prefix.

#### CLI Usage & Examples

```bash
# Audit sequence numbers and H1 titles across spec directory (check mode)
python 03-ai-scripts/15-sequence-and-title-auditor.py 02-spec

# Audit prompts directory
python 03-ai-scripts/15-sequence-and-title-auditor.py 01-prompts

# Auto-fix H1 title numbers in markdown files to match filename prefixes (fix mode)
python 03-ai-scripts/15-sequence-and-title-auditor.py 02-spec --fix
```

</details>

<details>
<summary><strong>16 — <code>16-installer-smoke-tester.py</code>: Bash and PowerShell Installer Sandbox Smoke Test Runner</strong></summary>

#### Why It Exists

Validates installer scripts (`install.sh`, `install.ps1`) in a sandbox before release, guaranteeing that scripts have no unreplaced placeholder tokens, include valid SHA-256 verification, and handle non-destructive rollbacks.

#### What It Does

- Scans shell and PowerShell installer files for unresolved tokens (e.g. `PLACEHOLDER`, `<version>`).
- Checks for checksum verification logic and safe update/fallback mechanisms.
- Verifies UNIX LF line endings for bash installers.
- Supports parallel testing with `--workers` and `--filter`.

#### CLI Usage & Examples

```bash
# Run installer smoke tests across repository (check mode)
python 03-ai-scripts/16-installer-smoke-tester.py

# Show detailed output for all installer scripts
python 03-ai-scripts/16-installer-smoke-tester.py --all

# Filter specific installer script by name
python 03-ai-scripts/16-installer-smoke-tester.py --filter "install.sh"

# Run sequentially and save JSON report
python 03-ai-scripts/16-installer-smoke-tester.py --sync --json --output tmp/installer-report.json
```

</details>

<details>
<summary><strong>17 — <code>17-fast-file-reader.py</code>: Instantaneous File and Directory Reader</strong></summary>

#### Why It Exists

Tailored specifically for AI agents to instantly inspect files, read directory trees, and search patterns in `<1ms` without spawning heavy subshells or consuming excessive context tokens.

#### What It Does

- Explores directories (`--list-folder`) with optional extension filters.
- Reads text files safely (`--read-file`) with bounded byte limits (`--max-bytes`, default 100KB) to avoid token context overflow.
- Performs fast substring or regex pattern searching (`--search-pattern`).

#### CLI Usage & Examples

```bash
# Read target file contents safely
python 03-ai-scripts/17-fast-file-reader.py --read-file 03-ai-scripts/01-index.md

# List folder contents with markdown extension filter
python 03-ai-scripts/17-fast-file-reader.py --list-folder 02-spec --ext .md

# Search pattern across files in directory
python 03-ai-scripts/17-fast-file-reader.py --search-pattern "AppError" --path 04-code --ext .go
```

</details>

<details>
<summary><strong>18 — <code>18-codebase-topology-discoverer.py</code>: Polyglot Architecture Mapper with TTL Caching</strong></summary>

#### Why It Exists

Discovers the architectural layout and language composition of any repository, routing AI agents and tools to the correct subsystem roots (Go backend, TypeScript frontend, SQL schemas, docs) without guessing.

#### What It Does

- Analyzes repository root and subdirectories to detect active languages (Go, Rust, Python, TypeScript, PHP, SQL, C#).
- Classifies subsystems (Backend, Database, Frontend, CI/CD, Documentation).
- Caches discovery results in `tmp/cache/paths/codebase-topology-cache.json` with a 30-minute TTL.
- Supports interactive queries (`--query`), summary breakdowns (`--summary`), and JSON exports (`--json`).

#### CLI Usage & Examples

```bash
# Print topology summary breakdown
python 03-ai-scripts/18-codebase-topology-discoverer.py --summary

# Force refresh topology cache
python 03-ai-scripts/18-codebase-topology-discoverer.py --refresh

# Query backend or language routing
python 03-ai-scripts/18-codebase-topology-discoverer.py --query go

# Output raw JSON topology map
python 03-ai-scripts/18-codebase-topology-discoverer.py --json
```

</details>

<details>
<summary><strong>19 — <code>19-artifact-remover.py</code>: Safe Test Artifact and Cache Cleaner with Git Index Hygiene</strong></summary>

#### Why It Exists

Safely cleans build artifacts, compiler binaries, test logs, and temporary caches while preserving Git hygiene and untracking accidentally committed files using `git rm`.

#### What It Does

- Categorizes artifacts by preset: `--clean-pycache`, `--clean-temp`, `--clean-binaries`, or `--clean-all`.
- Supports interactive confirmation, safety backup to a local Trash Bin, and `--permanent` unlinking.
- Implements Plan Mode First (`--plan` / `--dry-run`) to preview matching files before taking destructive action.

#### CLI Usage & Examples

```bash
# Preview all cleanup candidates without deleting anything (check / plan mode)
python 03-ai-scripts/19-artifact-remover.py --clean-all --plan

# Preview pycache cleanup in specific folder
python 03-ai-scripts/19-artifact-remover.py --clean-pycache --dir 03-ai-scripts --plan

# Clean temporary files and pycache safely (fix mode)
python 03-ai-scripts/19-artifact-remover.py --clean-temp --clean-pycache --force

# Clean all unneeded artifacts permanently
python 03-ai-scripts/19-artifact-remover.py --clean-all --force --permanent
```

</details>

<details>
<summary><strong>20 — <code>20-plan-consolidator.py</code>: Lovable Plan Index Synchronizer and Archiver</strong></summary>

#### Why It Exists

Automates the lifecycle management of Lovable execution plans and subtasks (`.ai-memory/plans/pending/`, `completed/`, `subtasks/`), keeping planning directories lean, archived, and sequentially ordered.

#### What It Does

- Automatically creates a timestamped safety backup git branch before any operations.
- Archives completed plans from `pending/` to `completed/`.
- Cleans up obsolete subtask files and untracks them from the Git index.
- Re-sequences remaining pending plans monotonically without number gaps.
- Supports `--dry-run` to preview all plan movements.

#### CLI Usage & Examples

```bash
# Preview plan consolidation and re-sequencing without changes (check / dry-run mode)
python 03-ai-scripts/20-plan-consolidator.py --dry-run

# Archive completed plan, re-sequence pending plans, and create safety branch (fix mode)
python 03-ai-scripts/20-plan-consolidator.py --archive 01-plan-name.md --resequence --backup --force

# Clean up obsolete subtask directories
python 03-ai-scripts/20-plan-consolidator.py --clean-subtasks --force
```

</details>

<details>
<summary><strong>21 — <code>21-sequence-integrity-linter.py</code>: Plan and Subtask Sequence Integrity Validator</strong></summary>

#### Why It Exists

Enforces integrity across all sequential documentation, prompts, execution plans, and agent skills, ensuring that numeric references and internal relative markdown links resolve to real files on disk.

#### What It Does

- Audits documents in `.ai-memory/prompts`, `.ai-memory/plans`, `.agents/skills`, `.ai-memory/coding-guidelines`, etc.
- Parses markdown links `[text](target)` and inline backtick paths.
- Resolves target paths against the repository root while ignoring template placeholders (`<slug>`, `vX.Y.Z`, `XX-`).
- Returns exit code 0 when all links resolve cleanly, or exit code 1 with broken link details.

#### CLI Usage & Examples

```bash
# Audit sequence integrity and relative links across repository plans and skills (check mode)
python 03-ai-scripts/21-sequence-integrity-linter.py
```

</details>

<details>
<summary><strong>22 — <code>22-doc-path-linter.py</code>: Documentation Relative Path Reference Checker</strong></summary>

#### Why It Exists

Autonomously verifies all markdown links and path references across `01-prompts/`, `02-spec/`, `.ai-memory/`, and `.agents/`, eliminating dead references and preventing AI agent hallucinations.

#### What It Does

- Executes comprehensive markdown path verification against all tracked documentation directories.
- Validates that every referenced file and anchor exists on disk.
- Reports exact missing files and line numbers, exiting with code 1 on dead links.

#### CLI Usage & Examples

```bash
# Verify all documentation and prompt path references exist on disk (check mode)
python 03-ai-scripts/22-doc-path-linter.py
```

</details>

<details>
<summary><strong>23 — <code>23-coding-guideline-path-consolidator.py</code>: Coding Guideline Cross-Link Canonicalizer</strong></summary>

#### Why It Exists

Autonomously consolidates references from the legacy nested path (`.ai-memory/coding-guidelines/coding-guidelines.md`) to the canonical path (`.ai-memory/coding-guidelines.md`) across all documentation, code, specs, and linters.

#### What It Does

- Recursively scans markdown, Python, JavaScript, TypeScript, Go, shell, and config files.
- Replaces old paths with the canonical relative path `.ai-memory/coding-guidelines.md`.
- Preserves file formatting and skips ignored directories (`.git`, `node_modules`, `tmp`).

#### CLI Usage & Examples

```bash
# Consolidate coding guideline paths across all repository files (fix mode)
python 03-ai-scripts/23-coding-guideline-path-consolidator.py
```

</details>

<details>
<summary><strong>24 — <code>24-spec-path-migrator.py</code>: Specification Directory Path Migrator</strong></summary>

#### Why It Exists

Migrates references from the old un-prefixed specification directory (`02-spec/`) to the standardized numbered directory (`02-spec/`) across all documentation, code, tests, and configuration files.

#### What It Does

- Replaces slash-based (`02-spec/...`) and backslash-based (`spec\\...`) path patterns with `02-spec/...`.
- Migrates specific known configuration and index references.
- Reports all modified files and the count of updated paths.

#### CLI Usage & Examples

```bash
# Migrate spec/ path references to 02-spec/ across the repository (fix mode)
python 03-ai-scripts/24-spec-path-migrator.py
```

</details>

<details>
<summary><strong>25 — <code>25-repo-migrator.py</code>: Transactional Repository Asset Migrator with Rollback</strong></summary>

#### Why It Exists

Provides a fully transactional, reversible engine to restructure repository layouts (such as migrating legacy `02-spec/` to `02-spec/`, `.ai-memory/prompts/` to `01-prompts/`, and `.ai-memory/ai-fix-scripts/` to `03-ai-scripts/`) with 100% undo/redo safety.

#### What It Does

- Logs all structural modifications to an atomic SQLite transaction journal (`tmp/migrations.db`).
- Supports Plan Mode First (`--plan` / `--dry-run`) to preview planned movements without touching disk.
- Provides full `undo` and `redo` subcommands using file content snapshots.
- Lists historical migration transactions (`ls` / `history`).

#### CLI Usage & Examples

```bash
# Preview layout migration actions without modifying files (check / plan mode)
python 03-ai-scripts/25-repo-migrator.py migrate --plan

# Inspect past migration transactions
python 03-ai-scripts/25-repo-migrator.py history

# Execute layout migration with confirmation bypass (fix mode)
python 03-ai-scripts/25-repo-migrator.py migrate --force

# Roll back the latest applied migration transaction (undo mode)
python 03-ai-scripts/25-repo-migrator.py undo

# Re-apply a rolled-back transaction (redo mode)
python 03-ai-scripts/25-repo-migrator.py redo
```

</details>

<details>
<summary><strong>26 — <code>26-go-code-formatter.py</code>: Cross-Platform gofmt Code Formatter with Staged File Support</strong></summary>

#### Why It Exists

Formats Go source files concurrently using native `gofmt` across all CPU cores, supporting pre-commit staged file filtering to keep Git commits clean and formatted.

#### What It Does

- Discovers all `.go` files across target directories or reads staged files via `git diff --cached`.
- Runs `gofmt -w` concurrently using a multi-worker thread pool.
- Supports sequential fallback (`--sync`), custom worker counts (`--workers`), and JSON output reports.

#### CLI Usage & Examples

```bash
# Format only staged Go files before commit (check / staged fix mode)
python 03-ai-scripts/26-go-code-formatter.py --staged

# Format all Go files in repository concurrently
python 03-ai-scripts/26-go-code-formatter.py

# Format specific directory with 4 workers
python 03-ai-scripts/26-go-code-formatter.py 04-code/golang --workers 4

# Run sequentially in 1 worker thread
python 03-ai-scripts/26-go-code-formatter.py --sync
```

</details>

<details>
<summary><strong>27 — <code>27-misspell-auditor.py</code>: American English Spelling Linter and Autofixer</strong></summary>

#### Why It Exists

Enforces American English spelling across repository code and documentation, preventing CODE RED British English spelling regressions (e.g. `behavior`, `color`, `initialize`).

#### What It Does

- Scans files using a high-performance regex dictionary of common British-to-American spelling variants (or delegates to the native `misspell` binary if installed in `PATH`).
- Supports `--staged` mode to check only git-staged changes before commit.
- When `--fix` is passed, automatically replaces British spellings with American English equivalents.

#### CLI Usage & Examples

```bash
# Audit repository files for British English spelling variants (check mode)
python 03-ai-scripts/27-misspell-auditor.py

# Audit only staged files in Git index
python 03-ai-scripts/27-misspell-auditor.py --staged

# Audit specific files or directories
python 03-ai-scripts/27-misspell-auditor.py 02-spec/ 01-prompts/

# Auto-fix British English spellings to US English in-place (fix mode)
python 03-ai-scripts/27-misspell-auditor.py --fix

# Auto-fix only staged files before commit
python 03-ai-scripts/27-misspell-auditor.py --staged --fix
```

</details>

<details>
<summary><strong>28 — <code>28-go-preflight-ci.py</code>: Standalone Go Unit Test and Linter Runner</strong></summary>

#### Why It Exists

Runs standalone Go unit tests and `golangci-lint` verification concurrently across all Go packages before pushing code, avoiding long CI test feedback loops.

#### What It Does

- Discovers all Go modules and packages in the repository.
- Concurrently executes `go test -v -race` and `golangci-lint run`.
- Supports phase selection: `all`, `test`, or `lint`.
- Features smart log filtering, sequential execution mode (`--sync`), and JSON output.

#### CLI Usage & Examples

```bash
# Run Go tests and linters concurrently (check mode, quiet on success, logs on failure)
python 03-ai-scripts/28-go-preflight-ci.py

# Run unit tests only
python 03-ai-scripts/28-go-preflight-ci.py test

# Run linters only
python 03-ai-scripts/28-go-preflight-ci.py lint

# Run with full banners and progress output
python 03-ai-scripts/28-go-preflight-ci.py --all-paths

# Filter specific package or test name
python 03-ai-scripts/28-go-preflight-ci.py --filter "appfault"

# Run sequentially (1 worker) and output JSON
python 03-ai-scripts/28-go-preflight-ci.py --sync --json --output tmp/go-preflight.json
```

</details>

<details>
<summary><strong>29 — <code>29-release-orchestrator.py</code>: Semantic Release Automation Orchestrator</strong></summary>

#### Why It Exists

Eliminates manual release ceremonies and branch errors by orchestrating the entire release lifecycle: SemVer bumping, changelog updates, release branch creation, Git tagging, and automatic return to the developer's original branch.

#### What It Does

- Detects and remembers the starting Git branch.
- Calculates the next SemVer version (`minor` by default, or `--tier patch|major`, or explicit `-v`).
- Executes version synchronization and commits bump changes.
- Creates the release branch `release/vX.Y.Z` and annotated Git tag `vX.Y.Z`.
- Switches back to the starting branch to preserve the developer's working environment.
- Supports `--dry-run` to preview all release operations safely.

#### CLI Usage & Examples

```bash
# Simulate a minor version release without modifying files (check / dry-run mode)
python 03-ai-scripts/29-release-orchestrator.py --dry-run

# Simulate a patch release with custom scope
python 03-ai-scripts/29-release-orchestrator.py -t patch -s "Fix boolean conditions" --dry-run

# Execute minor version release lifecycle (execution mode)
python 03-ai-scripts/29-release-orchestrator.py -t minor -s "Add AI scripts documentation"

# Execute explicit version release without pushing to remote
python 03-ai-scripts/29-release-orchestrator.py -v 1.5.0 --no-push
```

</details>

<details>
<summary><strong>30 — <code>30-enum-generator.py</code>: Multi-File Go Enum Generator with BasicEnum Scaffolding</strong></summary>

#### Why It Exists

Eliminates repetitive boilerplate and enforces type-safety standards by automatically scaffolding multi-file Go enum packages (`variant.go`, `vars.go`, `variant_test.go`, `readme.md`) adhering to repository conventions.

#### What It Does

- Generates 4 standardized files in `04-code/golang/pkg/enum/{name}type/`:
  - `variant.go`: Enum type definition, typed constants, predicates (`Is...`), and zero-value check.
  - `vars.go`: Labels lookup slice, map compiler, `All()`, `Values()`, and `Parse()` with `*appfault.AppError`.
  - `variant_test.go`: Complete unit tests verifying invariants, interfaces, and JSON round-trips.
  - `readme.md`: Package documentation.
- Supports underlying types: `byte`, `uint8`, `uint16`, `int`, `string`.
- Accepts CLI arguments or JSON configuration files.

#### CLI Usage & Examples

```bash
# Preview generated enum files without writing to disk (check / dry-run mode)
python 03-ai-scripts/30-enum-generator.py --name DeliveryMode --type byte --items Fast,Standard,Slow --dry-run

# Scaffold new Go enum package (generation mode)
python 03-ai-scripts/30-enum-generator.py --name DeliveryMode --type byte --items Fast,Standard,Slow

# Scaffold enum with custom package name and zero-value
python 03-ai-scripts/30-enum-generator.py --name ProcessState --type string --items Running,Paused,Stopped --zero-value Unknown --package state

# Generate from JSON configuration file
python 03-ai-scripts/30-enum-generator.py --config config/enums/delivery-mode.json
```

</details>

<details>
<summary><strong>31 — <code>31-md-gap-fixer.py</code>: Markdown Consecutive Empty Line Normalizer</strong></summary>

#### Why It Exists

Enforces clean markdown formatting by collapsing excessive blank lines (3 or more consecutive newlines / 2+ empty lines) into exactly one single empty line, preventing sprawling markdown files and visual bloat.

#### What It Does

- Scans markdown files (`.md`) across the repository or specified target path.
- Detects consecutive newline sequences `\n{3,}` and compresses them to `\n\n`.
- Supports check mode (audits and reports files with gaps, exiting 1 if gaps exist) and `--fix` mode (applies updates in-place).

#### CLI Usage & Examples

```bash
# Audit markdown files in target directory for excessive gaps (check mode, exits 1 on gaps)
python 03-ai-scripts/31-md-gap-fixer.py 03-ai-scripts/

# Audit single markdown file
python 03-ai-scripts/31-md-gap-fixer.py 03-ai-scripts/01-index.md

# Auto-fix excessive gaps in-place across directory (fix mode)
python 03-ai-scripts/31-md-gap-fixer.py 03-ai-scripts/ --fix

# Auto-fix single markdown file in-place
python 03-ai-scripts/31-md-gap-fixer.py 03-ai-scripts/01-index.md --fix
```

</details>

<details>
<summary><strong>32 — <code>32-git-history-file-tracer.py</code>: Git History Removed File Tracer, Restorer & Purger</strong></summary>

#### Why It Exists

Provides a high-performance auditing and recovery tool to trace files deleted throughout Git history. It prevents accidental knowledge loss during plan compaction, enables fast restoration of removed files to working tree or staging folders, and provides an administrative purge mechanism to permanently scrub deleted files from all Git commits, branches, and tags.

#### What It Does

- Deeply inspects Git deletion history (`git log --diff-filter=D`) across the repository or targeted subdirectories.
- Provides built-in fast presets:
  - `--preset-audit` (or `--audit`, `--spec-audit`): Targets spec audit reports directory (`02-spec/25-app-spec-audit/`).
  - `--preset-lovable` (or `--lovable`): Targets `.ai-memory/` folder.
  - `--preset-plans` (or `--plans`): Targets `.ai-memory/plans/`.
  - `--preset-subtasks` (or `--subtasks`): Targets `.ai-memory/plans/subtasks/`.
  - `--preset-spec` (or `--spec`): Targets `02-spec/`.
- Displays a structured, numbered pre-flight report with deletion commit SHA, date, author, commit message, pre-deletion file size, and current on-disk presence.
- Supports selective filtering via numbered ranges/lists (`--include 1-5`, `--exclude 2,4`) and glob patterns (`--exclude-pattern`).
- Restores removed files from pre-deletion commits either in-place (`--restore`) or into a staging directory (`--restore-to <dir>`).
- Safely moves working tree files to the OS Recycle Bin / Trash (`--delete`, `--remove`, `--trash`) after creating a full physical backup in the OS temporary directory.
- Safely purges deleted files from entire Git history across all refs/tags (`--purge --confirm-purge`) with automatic pre-purge safety backup branches and OS temp directory backups.

#### CLI Usage & Examples

```bash
# Pre-flight inspection of removed audit files in 02-spec/25-app-spec-audit/
python 03-ai-scripts/32-git-history-file-tracer.py --preset-audit

# Pre-flight inspection of removed markdown files in .ai-memory/
python 03-ai-scripts/32-git-history-file-tracer.py --preset-lovable

# Pre-flight inspection of removed subtasks in .ai-memory/plans/subtasks/
python 03-ai-scripts/32-git-history-file-tracer.py --preset-subtasks

# Custom path and extension scan from repository root
python 03-ai-scripts/32-git-history-file-tracer.py --path 04-code/ --ext .go,.ts

# Filter by numbering: include items 1 to 5 and exclude item 3
python 03-ai-scripts/32-git-history-file-tracer.py --preset-lovable --include 1-5 --exclude 3

# Restore selected files to a staging directory
python 03-ai-scripts/32-git-history-file-tracer.py --preset-lovable --include 1,2 --restore-to tmp/recovery/

# Restore selected files in-place into repository working tree
python 03-ai-scripts/32-git-history-file-tracer.py --preset-lovable --include 1 --restore

# Safely move working tree files to OS Recycle Bin (with OS temp backup)
python 03-ai-scripts/32-git-history-file-tracer.py --preset-audit --delete

# Permanently purge selected removed files from all Git history (all branches and tags)
python 03-ai-scripts/32-git-history-file-tracer.py --preset-lovable --include 1,2 --purge --confirm-purge
```

</details>

<details>
<summary><strong>33 — <code>33-test-inventory-generator.py</code>: Test Inventory Generator & Atomic Change Recorder</strong></summary>

#### Why It Exists

Maintains a single source of truth for repository test coverage at `.ai-memory/test-inventory.json` (similar to gitmap architecture) and provides a concurrency-safe atomic change tracking mechanism (`.ai-memory/temp/recent-file-changes.json`) with cross-platform file locking so multi-agent tasks can record file modifications without race conditions or crashes.

#### What It Does

- Discovers and parses unit tests across Go (`*_test.go`), TypeScript/JavaScript (`*.test.ts`, `*.spec.ts`), and Python (`test_*.py`).
- Maps tests to target source files, extracts function names, and generates SHA-256 code hashes.
- Generates and maintains `.ai-memory/test-inventory.json` with test counts, status tracking, and dirty flags.
- Under file lock (`.ai-memory/temp/recent-file-changes.lock`), records distinct repository-relative file paths to `.ai-memory/temp/recent-file-changes.json` and automatically resolves all associated tests that must be executed when release verification or test fixes are explicitly requested.

#### CLI Usage & Examples

```bash
# Scan repository and generate / update .ai-memory/test-inventory.json
python 03-ai-scripts/33-test-inventory-generator.py

# Check test inventory cache freshness (<= 5 days old and profiled)
python 03-ai-scripts/33-test-inventory-generator.py --check-age --max-age-days 5

# Safely record modified files under lock and resolve associated tests
python 03-ai-scripts/33-test-inventory-generator.py --record "04-code/golang/pkg/appfault/appfault.go"

# Query currently recorded modified files and associated test list
python 03-ai-scripts/33-test-inventory-generator.py --query-recent
```

</details>

<details>
<summary><strong>34 — <code>34-purge-github-actions-artifacts.py</code>: GitHub Actions Artifact Purger</strong></summary>

#### Why It Exists

Autonomously discovers and deletes stored GitHub Actions artifacts via the GitHub API across target repositories to enforce the Zero-Storage Actions Mandate and keep account storage usage at 0.0 GB (staying well within the 0.5 GB free quota).

#### What It Does

- Queries the GitHub Actions API for uploaded artifacts across repositories (`alimtvnetwork/coding-guidelines-v24`, `alimtvnetwork/gitmap-v28`).
- Concurrently purges old artifacts using thread pools to avoid quota exhaustion.
- Enforces zero artifact retention policies without breaking commit statuses.

#### CLI Usage & Examples

```bash
# Purge artifacts across default repositories
python 03-ai-scripts/34-purge-github-actions-artifacts.py

# Purge artifacts for a specific repository
python 03-ai-scripts/34-purge-github-actions-artifacts.py --repo alimtvnetwork/coding-guidelines-v24
```

</details>

<details>
<summary><strong>35 — <code>35-db-struct-enum-generator.py</code>: Database Model Struct & Type-Safe Column Enum Generator</strong></summary>

#### Why It Exists

Inspects Go model structs and auto-generates type-safe column enums, `enums/consts.go`, parent package aliases, and strongly-typed repository query builders utilizing `coding-guidelines/common/pkg/dbengine` and `coding-guidelines/common/pkg/appfault`.

#### What It Does

- Parses Go struct definitions, extracting public fields and `db:"column_name"` tags.
- Automatically resolves Go module root and computes relative package import paths.
- Generates `enums/<model>.go` with strict O(1) map validation, JSON serialization/deserialization with `appfault.AppError`, and type-safe `.Is<Field>()` predicates.
- Generates generic `Repository` constructors, row scanners, and CRUD query builder helpers.
- Formats all generated Go files using `gofmt`.

#### CLI Usage & Examples

```bash
# Preview generation across target model directory (dry run)
python 03-ai-scripts/35-db-struct-enum-generator.py --dir 04-code/golang/pkg/dbengine --dry-run

# Generate enums and repositories from specific Go file
python 03-ai-scripts/35-db-struct-enum-generator.py --file 04-code/golang/pkg/models/item.go

# Generate enums into explicit output directory
python 03-ai-scripts/35-db-struct-enum-generator.py --file 04-code/golang/pkg/models/item.go --out-dir 04-code/golang/pkg/generated/item
```

</details>

<details>
<summary><strong>37 — <code>37-bump-version.py</code>: Repository-Aware SemVer Version Bumper & Synchronizer</strong></summary>

#### Why It Exists

Automates accurate SemVer version bumping and synchronization across project manifests (`version.json`, `package.json`), documentation (`readme.md`), and changelogs (`changelog.md`), preventing manual version drift and ensuring release pipelines have a reliable cross-platform Python tool.

#### What It Does

- Discovers canonical version from `version.json` or `package.json`.
- Calculates next SemVer version (`minor` by default per Rule 0, or `patch`/`major`).
- Updates `version.json`, `package.json`, `prompt-version.template.json`, and `readme.md` version badges and install snippets.
- Prepends formatted changelog entries to `changelog.md`.
- Triggers `npm run sync` if defined to regenerate spec trees and health scores.
- Supports `--dry-run` to preview all file modifications.

#### CLI Usage & Examples

```bash
# Preview minor version bump without modifying files (check / dry-run mode)
python 03-ai-scripts/37-bump-version.py --dry-run

# Bump minor version (default) with custom scope
python 03-ai-scripts/37-bump-version.py --tier minor --scope "Add release branching lifecycle"

# Bump patch version
python 03-ai-scripts/37-bump-version.py --tier patch --scope "Fix release orchestrator flow"

# Explicit version bump
python 03-ai-scripts/37-bump-version.py --version 6.42.0
```

</details>

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
6. **Quality Gates & Test Execution Policy:**
   - For standard tasks: Execute `python 03-ai-scripts/06-cicd-local-runner.py --no-tests` (test execution is strictly disabled unless explicitly commanded by the repository owner).
   - For release tasks: Execute `python 03-ai-scripts/06-cicd-local-runner.py --run-tests` to ensure 100% test passage before release.
   - For CI/CD fix tasks (`ci-cd-fix`): Execute `python 03-ai-scripts/06-cicd-local-runner.py` (or `--run-tests`) to run all quality gates, linters, and unit test suites properly to catch and repair all pipeline failures.
