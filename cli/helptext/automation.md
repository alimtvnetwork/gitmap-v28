# gitmap automation

High-performance native Go automation subsystem, multi-core search, coding guideline enforcement, database migrations, CI/CD preflight gates, release orchestration, and git hygiene.

## Usage

```bash
gitmap automation <subcommand> [flags]
gitmap aum <subcommand> [flags]
gitmap auto <subcommand> [flags]
gitmap scripts <subcommand> [flags]
```

## Description

`gitmap automation` (canonical alias: `gitmap aum`) provides high-speed compiled Go tooling replacing legacy external scripts. It implements thread-safe lazy regex compilation, multi-core streaming directory traversal, universal CRLF-to-LF and trailing whitespace normalization, sub-millisecond in-memory caching (<0.05ms) with zero disk bloat, coding guideline linters, SQLite schema code generators and migration runners, parallel CI preflight checkers, automated release bumping, and safe git hygiene purgers.

## Subcommands Overview

### Phase 1: Code Quality, Naming & Coding Guidelines

- `search <pattern> [dir]` (aliases: `grep`, `find-text`): Multi-core streaming search across repository files with lazy regex or literal match.
- `newlines [paths...]` (aliases: `fix-newlines`, `lf`): Polyglot CRLF to LF and trailing whitespace normalizer.
- `cache [status|read|warm|clear]` (alias: `mem-cache`): Manage sub-millisecond in-memory file cache with zero disk pollution.
- `benchmark [target]` (aliases: `bench`, `compare`): Run side-by-side Go vs Python execution benchmarks.
- `relative-paths [dir]` (aliases: `rel-paths`, `paths`, `links`): Audit and sanitize forbidden absolute filesystem paths in documentation and code.
- `naming [dir]` (aliases: `bool-naming`, `conventions`): Audit boolean comparison anti-patterns and affirmative naming conventions.
- `result-wrapper [dir]` (aliases: `res-wrapper`, `result-audit`): Audit Go functions returning multi-value map/slice error tuples.
- `params [dir]` (aliases: `param-audit`, `arity`): Audit function signatures for argument reduction and parameter structs.
- `enums [dir]` (aliases: `enum-audit`, `types-enum`): Audit enums for `*Type` suffixes and ban raw numeric rune casts.
- `guard [dir]` (aliases: `file-guard`, `size-guard`, `blob-guard`): Audit repository file sizes, common binaries, and large JSONs.
- `sequence [dir]` (aliases: `seq`, `titles`, `seq-audit`): Audit markdown sequence numbering gaps and H1 header alignment.
- `exclude [list|add|rm|clear]` (aliases: `waiver`, `exclusions`): Manage persistent search and audit exclusion patterns.

### Phase 2: Topology, Schema & Database Generation

- `topology [dir]` (aliases: `topo`, `codebase-topology`): Discover and cache codebase topology (subsystems, schemas, workflows, languages).
- `db-generate [db-path]` (aliases: `db-gen`, `gen-db`): Generate Go structs, TypeScript interfaces, and column enums from SQLite tables.
- `db-migrate [db-path] [migrations-dir]` (aliases: `migrate`, `db-up`): Execute ordered SQL migrations with rollback protection and tracking.
- `schema-audit [db-path]` (aliases: `audit-schema`, `db-audit`): Validate SQLite split-db schemas against naming conventions (PascalCase, affirmative booleans, PKs).

### Phase 3: CI/CD, Local Pre-Flight & Multi-Core Checkers

- `preflight [dir]` (aliases: `check-all`, `ci-local`, `pre-commit`): Parallel multi-core CI preflight checks (gofmt, relpaths, nested-if, naming).
- `test-inventory [dir]` (aliases: `tests-inv`, `inventory`): Discover unit and integration tests, output inventory manifest with durations.
- `purge-actions` (aliases: `purge-artifacts`, `clean-actions`): Query GitHub Actions API and purge obsolete artifacts to maintain 0.0 GB footprint.
- `smoke-test [dir]` (aliases: `installer-smoke`, `smoke`): Cross-platform dry-run validator for installed tools and installer scripts.

### Phase 4: Release, SemVer & Version Synchronization

- `version-sync [dir]` (aliases: `sync-version`, `ver-sync`): Audit and synchronize repository version manifests.
- `release-bump [major|minor|patch]` (aliases: `bump`, `semver-bump`): Compute next semantic version and update repository manifests.
- `milestones [milestone-id]` (aliases: `milestone`, `consolidate-milestones`): Query GitHub milestone issues and format structured release notes.

### Phase 5: Documentation, Spec Migration & Memory Consolidation

- `help-audit [dir]` (aliases: `help-check`, `doc-audit`): Audit CLI commands for help descriptions and documentation parity.
- `plan-consolidate [dir]` (aliases: `consolidate-plans`, `consolidate`): Cluster completed plans and subtasks into milestone summaries.
- `doc-links [dir]` (aliases: `check-links`, `doc-paths`): Validate markdown relative links across documentation and memory.
- `spec-migrate [dir]` (aliases: `resequence-spec`, `migrate-spec`): Re-sequence specification file prefixes and update cross-links.

### Phase 6: Git Hygiene, Cleanup & History Purging

- `clean-artifacts` (aliases: `clean-build`, `rm-artifacts`): Safely delete build binaries, test dumps, pycache, and temporary files.
- `changed-files` (aliases: `git-changes`, `diff-files`): Discover modified, staged, and untracked files for targeted linter passes.
- `purge-history` (aliases: `trace-history`, `clean-history`): Trace and identify large historical git blobs without rewriting recent commits.
- `format-go [dir]` (aliases: `gofmt-ast`, `fmt-go`): AST-aware Go code formatter organizing imports and enforcing UTF-8 LF.

## Command Reference & Flags

### Phase 1: Code Quality, Naming & Coding Guidelines

#### Search Flags (`gitmap aum search <pattern> [dir]`)

Multi-core streaming search across repository files with lazy regex or literal match fast paths.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--regex` | `-r` | `false` | Enable regular expression matching |
| `--ignore-case` | `-i` | `false` | Perform case-insensitive matching |
| `--ext` | `-e` | `nil` | Filter search by file extensions (e.g. `.go, .ts`) |
| `--workers` | `-w` | `0` | Number of concurrent worker threads (default: CPU count) |
| `--max-json-kb` | | `500` | Maximum JSON size in KB before auto-exclusion |
| `--include-binaries` | | `false` | Include binary files in search |
| `--include-large-json` | | `false` | Include oversized JSON files in search |

#### Newlines Flags (`gitmap aum newlines [paths...]`)

Polyglot CRLF to LF and trailing whitespace normalizer.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--fix` | `-f` | `false` | Apply normalized line endings and whitespace to disk |
| `--dry-run` | | `false` | Preview modified file counts without writing to disk |

#### Cache Commands (`gitmap aum cache [action] [path]`)

Manage sub-millisecond in-memory file cache (<0.05ms) with zero disk pollution.

| Action | Description |
|--------|-------------|
| `status` | Display current in-memory cache capacity, file count, and hit/miss statistics |
| `read <path>` | Retrieve a file directly from memory cache or load and warm into cache |
| `warm [dir]` | Pre-populate the in-memory cache for a directory |
| `clear` / `purge` / `rm` | Purge the in-memory cache and reclaim all memory |

#### Benchmark Commands (`gitmap aum benchmark [target]`)

Run side-by-side Go vs Python execution benchmarks for performance verification.

| Target | Description |
|--------|-------------|
| `all` | Benchmark all subsystems (search, newlines, cache) |
| `search` | Benchmark multi-core Go streaming search against legacy Python grep |
| `newlines` | Benchmark Go LF normalizer against legacy Python newline normalizer |

#### Relative-Paths Flags (`gitmap aum relative-paths [dir]`)

Audit and sanitize forbidden absolute filesystem paths (`D:\...`, `file:///`) in documentation and code.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--fix` | `-f` | `false` | Auto-fix recognized path patterns |
| `--ext` | `-e` | `nil` | Filter by file extensions (e.g. `.md, .ts`) |
| `--json` | | `false` | Output results as machine-readable JSON |

#### Naming Flags (`gitmap aum naming [dir]`)

Audit boolean comparison anti-patterns (`== true`, `=== true`) and enforce affirmative boolean naming (`is*`, `has*`).

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--ext` | `-e` | `nil` | Filter by file extensions |
| `--json` | | `false` | Output results as machine-readable JSON |

#### Result-Wrapper Flags (`gitmap aum result-wrapper [dir]`)

Audit Go functions returning multi-value `(map, error)` or `([]T, error)` tuples and enforce `ResultMap`/`ResultSlice` monads.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--json` | | `false` | Output results as machine-readable JSON |

#### Params Flags (`gitmap aum params [dir]`)

Audit function signatures exceeding 3–4 parameters and recommend dedicated `*Params` structs in `types.go`.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--json` | | `false` | Output results as machine-readable JSON |

#### Enums Flags (`gitmap aum enums [dir]`)

Audit enum definitions for mandatory `*Type` suffixes and ban raw numeric `rune(10)` casts.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--json` | | `false` | Output results as machine-readable JSON |

#### Guard Flags (`gitmap aum guard [dir]`)

Audit repository file sizes, common binaries, and large JSONs with interactive exclusion prompts.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--max-kb` | `-k` | `500` | Maximum allowed file size in KB |
| `--max-json-kb` | | `500` | Maximum allowed JSON size before auto-exclusion |
| `--interactive` | `-i` | `true` | Prompt interactively when binaries are found |
| `--auto-exclude` | | `false` | Automatically exclude detected binaries without prompting |
| `--json` | | `false` | Output results as machine-readable JSON |

#### Sequence Flags (`gitmap aum sequence [dir]`)

Audit markdown sequence numbering gaps and H1 header alignment across documentation and specifications.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--fix` | `-f` | `false` | Automatically fix mismatched H1 headers to file prefix numbers |
| `--json` | | `false` | Output results as machine-readable JSON |

#### Exclude Commands (`gitmap aum exclude [action] [pattern] [reason]`)

Manage persistent search and audit exclusion patterns stored in SQLite (`sql.db`).

| Action | Description |
|--------|-------------|
| `list` | Display all configured search and audit exclusion patterns |
| `add <pattern> [reason]` | Register a persistent exclusion pattern with an optional reason |
| `rm <pattern>` | Remove an existing exclusion pattern |
| `clear` / `purge` | Clear all custom exclusion patterns |

---

### Phase 2: Topology, Schema & Database Generation

#### Topology Flags (`gitmap aum topology [dir]`)

Discover and cache codebase topology (subsystems, schemas, workflows, languages) with routing query support.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--refresh` | `-r` | `false` | Force refresh topology discovery cache |
| `--json` | `-j` | `false` | Output results as raw JSON |
| `--query` | `-q` | `""` | Search specific subsystem or language routing |
| `--ttl` | | `1800` | Cache TTL in seconds |

#### Db-Generate Flags (`gitmap aum db-generate [db-path]`)

Generate Go structs, TypeScript interfaces, and column enums from SQLite tables.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--lang` | | `all` | Target language: `go`, `ts`, or `all` |
| `--out` | | `""` | Output directory for generated files |
| `--dry-run` | | `false` | Preview generated code without writing |
| `--struct-dir` | | `""` | Optional directory of model structs |

#### Db-Migrate Flags (`gitmap aum db-migrate [db-path] [migrations-dir]`)

Execute ordered SQL migrations with transaction rollback protection and migration history tracking.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--dry-run` | | `false` | Preview migrations without writing |
| `--status` | | `false` | Display migration history and current version |
| `--rollback` | | `0` | Rollback last N migrations |
| `--sql` | | `""` | Execute inline SQL migration statement |

#### Schema-Audit Flags (`gitmap aum schema-audit [db-path]`)

Validate SQLite split-db schemas against naming conventions (PascalCase, affirmative booleans, PKs).

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--strict` | | `false` | Enforce strict column naming and PK rules |
| `--dir` | | `""` | Target directory containing SQL or Go schemas |
| `--json` | | `false` | Output results as machine-readable JSON |

---

### Phase 3: CI/CD, Local Pre-Flight & Multi-Core Checkers

#### Preflight Flags (`gitmap aum preflight [dir]`)

Parallel multi-core CI preflight checks across CPU cores (gofmt, relpaths, nested-if, naming).

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--fail-fast` | | `false` | Halt on first check failure |
| `--workers` | `-w` | `0` | Number of worker threads (default: CPU cores) |
| `--phase` | `-p` | `all` | Execution phase (`all`, `test`, `lint`) |
| `--filter` | `-k` | `""` | Filter checks by name substring |
| `--json` | | `false` | Output results as machine-readable JSON |

#### Test-Inventory Flags (`gitmap aum test-inventory [dir]`)

Discover unit and integration tests, output inventory manifest with duration baselines.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--out` | | `""` | Output manifest path (default: `.ai-memory/test-inventory.json`) |
| `--refresh` | | `false` | Force refresh existing inventory manifest |
| `--slow-threshold` | | `4.0` | Threshold in seconds to categorize slow tests |
| `--force-run-all` | | `false` | Mark all discovered tests as needing run |
| `--record` | | `nil` | Record modified file paths safely |
| `--query-recent` | | `false` | Query tests associated with recent changes |
| `--clear` | | `false` | Clear recent changes log |
| `--json` | | `false` | Output results as machine-readable JSON |

#### Purge-Actions Flags (`gitmap aum purge-actions`)

Query GitHub Actions API and purge obsolete artifacts to maintain 0.0 GB storage footprint.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--repo` | | `""` | Target repository slug (`owner/repo`) |
| `--older-than` | | `0` | Purge artifacts older than N days |
| `--dry-run` | | `false` | Preview artifacts without deleting |
| `--workers` | `-w` | `12` | Number of concurrent deletion threads |
| `--json` | | `false` | Output results as machine-readable JSON |

#### Smoke-Test Flags (`gitmap aum smoke-test [dir]`)

Cross-platform dry-run validator for installed tools and installer scripts.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--tools` | | `nil` | List of tool binaries to verify in PATH |
| `--workers` | `-w` | `0` | Number of worker threads (default: CPU cores) |
| `--filter` | `-k` | `""` | Filter targets by name substring |
| `--json` | | `false` | Output results as machine-readable JSON |

---

### Phase 4: Release, SemVer & Version Synchronization

#### Version-Sync Flags (`gitmap aum version-sync [dir]`)

Audit and synchronize repository version manifests across Go constants, package files, and documentation.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--fix` | `-f` | `false` | Automatically harmonize mismatched manifests |
| `--target-version` | `-t` | `""` | Override canonical version to harmonize |
| `--json` | | `false` | Output results as machine-readable JSON |

#### Release-Bump Flags (`gitmap aum release-bump [major|minor|patch]`)

Compute next semantic version and update repository manifests across SSoT files.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--version` | `-v` | `""` | Explicit version string to bump to (e.g. `6.263.0`) |
| `--dry-run` | `-d` | `false` | Preview version bump without modifying files |
| `--tag` | | `false` | Create git release tag (`vX.Y.Z`) |
| `--push` | | `false` | Push release branch and tag to remote |
| `--skip-tests` | | `false` | Skip pre-release quality gate checks |
| `--scope` | `-s` | `"Automated release orchestration"` | Scope description for release commit |
| `--bullet` | `-b` | `nil` | Changelog bullet points (repeatable) |
| `--json` | | `false` | Output results as machine-readable JSON |

#### Milestones Flags (`gitmap aum milestones [milestone-id]`)

Query GitHub milestone issues and format structured release notes.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--out` | `-o` | `""` | Output file path for generated release notes |
| `--json` | | `false` | Output results as machine-readable JSON |

---

### Phase 5: Documentation, Spec Migration & Memory Consolidation

#### Help-Audit Flags (`gitmap aum help-audit [dir]`)

Audit CLI commands for help descriptions and documentation parity across `cli/helptext/`.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--strict` | `-s` | `false` | Fail with exit code 1 if violations exist |
| `--ext` | `-e` | `nil` | Filter by file extensions |
| `--json` | | `false` | Output results as machine-readable JSON |

#### Plan-Consolidate Flags (`gitmap aum plan-consolidate [dir]`)

Cluster completed plans and subtasks in `.ai-memory/plans/` into milestone summaries.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--threshold` | `-t` | `5` | Minimum completed plans threshold to cluster |
| `--dry-run` | `-d` | `false` | Preview consolidation without modifying files |
| `--force` | | `false` | Bypass confirmation prompts |
| `--json` | | `false` | Output results as machine-readable JSON |

#### Doc-Links Flags (`gitmap aum doc-links [dir]`)

Validate markdown relative links across documentation and memory.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--fix` | `-f` | `false` | Autofix known outdated path references |
| `--json` | | `false` | Output results as machine-readable JSON |

#### Spec-Migrate Flags (`gitmap aum spec-migrate [dir]`)

Re-sequence specification file prefixes and update cross-links across `02-spec/`.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--from` | | `0` | Source spec number to migrate from |
| `--to` | | `0` | Target spec number to migrate to |
| `--dry-run` | `-d` | `false` | Preview migration without disk changes |
| `--json` | | `false` | Output results as machine-readable JSON |

---

### Phase 6: Git Hygiene, Cleanup & History Purging

#### Clean-Artifacts Flags (`gitmap aum clean-artifacts`)

Safely delete build binaries, test dumps, pycache, and temporary files.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--dir` | | `.` | Root directory to scan for artifacts |
| `--dry-run` | `-n` | `false` | Preview matching items without deleting |
| `--all` | `-a` | `false` | Clean all preset categories (pycache, temp, binaries) |
| `--verbose` | `-v` | `false` | Show all discovered artifact paths |
| `--clean-pycache` | | `false` | Remove python bytecode and test cache |
| `--clean-temp` | | `false` | Remove temporary files (`.tmp, .log, .swp`) |
| `--clean-binaries` | | `false` | Remove compiled binaries (`.exe, .syso, etc.`) |
| `--json` | | `false` | Output results as machine-readable JSON |

#### Changed-Files Flags (`gitmap aum changed-files`)

Discover modified, staged, and untracked files for targeted linter passes.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--dir` | | `.` | Target repository directory |
| `--base` | | `""` | Git base reference or commit (e.g. `origin/main, HEAD~1`) |
| `--commits` | `-n` | `20` | Number of recent commits to evaluate |
| `--staged-only` | | `false` | Discover staged files only |
| `--verify` | | `false` | Verify that files currently exist on disk |
| `--json` | | `false` | Output results as machine-readable JSON |

#### Purge-History Flags (`gitmap aum purge-history`)

Trace and identify large historical git blobs without rewriting recent commits.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--dir` | | `.` | Target repository directory |
| `--min-size-mb` | | `1.0` | Minimum blob size in MB to identify |
| `--path` | | `""` | Target file or path pattern to trace |
| `--dry-run` | `-n` | `true` | Preview large blobs without rewriting history |
| `--confirm` | `-y` | `false` | Bypass confirmation prompt |
| `--json` | | `false` | Output results as machine-readable JSON |

#### Format-Go Flags (`gitmap aum format-go [dir]`)

AST-aware Go code formatter organizing imports and enforcing UTF-8 LF.

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--dir` | | `.` | Directory or file to format |
| `--write` | `-w` | `false` | Write formatted changes back to disk |
| `--check` | `-c` | `false` | Check and return violations if files need formatting |
| `--staged` | | `false` | Format only staged Go files |
| `--json` | | `false` | Output results as machine-readable JSON |

## Examples

### Phase 1: Code Quality, Naming & Coding Guidelines

```bash
# 1. Multi-core regex search across the codebase with lazy regex compilation
gitmap aum search "func Run" cli --ext .go

# 2. Case-insensitive literal search bypassing the regex engine
gitmap aum search "connection reset" --ignore-case

# 3. Polyglot newline normalization in dry-run mode
gitmap aum newlines --dry-run

# 4. Atomically normalize CRLF to LF and trim whitespace across repository
gitmap aum newlines --fix

# 5. Pre-warm and inspect sub-millisecond in-memory cache
gitmap aum cache warm cli
gitmap aum cache status

# 6. Read file directly from memory cache (<0.05ms)
gitmap aum cache read cli/helptext/automation.md

# 7. Run side-by-side Go vs Python performance benchmark
gitmap aum benchmark all
gitmap aum benchmark search
gitmap aum benchmark newlines

# 8. Audit and auto-fix forbidden absolute filesystem paths in documentation
gitmap aum relative-paths docs --fix

# 9. Audit boolean comparison anti-patterns and affirmative naming
gitmap aum naming cli --ext .go

# 10. Audit functions returning multi-value error tuples instead of Result monads
gitmap aum result-wrapper cli

# 11. Audit function signatures exceeding 4 parameters
gitmap aum params cli

# 12. Audit enums for mandatory *Type suffixes and ban numeric rune casts
gitmap aum enums cli

# 13. Audit file sizes, oversized JSON files, and binary files
gitmap aum guard . --max-kb 500 --interactive

# 14. Audit markdown sequence numbering gaps and auto-fix H1 titles
gitmap aum sequence 02-spec --fix

# 15. Manage persistent search and audit exclusion patterns
gitmap aum exclude list
gitmap aum exclude add "build/bin/*" "compiled binaries"
gitmap aum exclude rm "build/bin/*"
```

### Phase 2: Topology, Schema & Database Generation

```bash
# 1. Discover codebase topology and cache subsystem graph
gitmap aum topology .

# 2. Query routing information for a specific subsystem or language
gitmap aum topology --query cli
gitmap aum topology --query go

# 3. Generate Go structs, TypeScript interfaces, and enums from SQLite tables
gitmap aum db-generate .gitmap/data/automation/sql.db --lang all --out pkg/models

# 4. Preview database code generation in dry-run mode
gitmap aum db-generate .gitmap/data/automation/sql.db --dry-run

# 5. Execute ordered SQL database migrations
gitmap aum db-migrate .gitmap/data/automation/sql.db migrations/

# 6. Check database migration status and history
gitmap aum db-migrate .gitmap/data/automation/sql.db --status

# 7. Rollback last migration
gitmap aum db-migrate .gitmap/data/automation/sql.db --rollback 1

# 8. Audit SQLite split-db schemas against repository naming conventions
gitmap aum schema-audit .gitmap/data/automation/sql.db --strict
```

### Phase 3: CI/CD, Local Pre-Flight & Multi-Core Checkers

```bash
# 1. Run all local preflight multi-core quality checks
gitmap aum preflight .

# 2. Run preflight with fail-fast enabled on 8 worker threads
gitmap aum preflight . --fail-fast --workers 8

# 3. Filter preflight checks by keyword
gitmap aum preflight . --filter relpaths

# 4. Generate test inventory manifest with duration baselines
gitmap aum test-inventory . --out .ai-memory/test-inventory.json

# 5. Refresh test inventory and categorize slow tests exceeding 3 seconds
gitmap aum test-inventory . --refresh --slow-threshold 3.0

# 6. Purge obsolete GitHub Actions artifacts to maintain 0.0 GB storage
gitmap aum purge-actions --repo alimtvnetwork/gitmap-v28 --older-than 7

# 7. Preview Actions artifact purge in dry-run mode
gitmap aum purge-actions --dry-run

# 8. Cross-platform dry-run smoke test for installer tools
gitmap aum smoke-test . --tools gitmap,agy,go
```

### Phase 4: Release, SemVer & Version Synchronization

```bash
# 1. Audit version synchronization across all repository manifests
gitmap aum version-sync .

# 2. Automatically synchronize and harmonize mismatched version manifests
gitmap aum version-sync . --fix --target-version 6.263.0

# 3. Preview semantic version patch bump in dry-run mode
gitmap aum release-bump patch --dry-run

# 4. Execute minor version bump with changelog bullets and git tag
gitmap aum release-bump minor --tag --bullet "Added AUM automation suite" --bullet "Expanded helptext documentation"

# 5. Query GitHub milestone issues and format structured release notes
gitmap aum milestones 37 --out docs/release-notes-v6.263.0.md
```

### Phase 5: Documentation, Spec Migration & Memory Consolidation

```bash
# 1. Audit CLI commands for help text parity and AST completeness
gitmap aum help-audit cli

# 2. Enforce strict CLI help audit in CI
gitmap aum help-audit cli --strict

# 3. Cluster completed plans in .ai-memory/plans/ into milestone summaries
gitmap aum plan-consolidate .ai-memory/plans --threshold 5

# 4. Preview plan memory consolidation in dry-run mode
gitmap aum plan-consolidate .ai-memory/plans --dry-run

# 5. Validate markdown relative link integrity across documentation
gitmap aum doc-links 02-spec

# 6. Autofix known outdated markdown documentation links
gitmap aum doc-links 02-spec --fix

# 7. Re-sequence specification number prefixes and update cross-links
gitmap aum spec-migrate 02-spec --from 120 --to 121 --dry-run
```

### Phase 6: Git Hygiene, Cleanup & History Purging

```bash
# 1. Preview artifact cleanup across build binaries, pycache, and temp dumps
gitmap aum clean-artifacts --all --dry-run

# 2. Safely clean python bytecode and test caches
gitmap aum clean-artifacts --clean-pycache

# 3. Clean temporary files (.tmp, .log, .swp) with verbose file listing
gitmap aum clean-artifacts --clean-temp --verbose

# 4. Discover modified, staged, and untracked files for targeted linter passes
gitmap aum changed-files --base origin/main

# 5. Discover staged files only
gitmap aum changed-files --staged-only

# 6. Trace and identify large historical git blobs (>2 MB) without rewriting history
gitmap aum purge-history --min-size-mb 2.0 --dry-run

# 7. Check Go source files for AST formatting and import order violations
gitmap aum format-go cli --check

# 8. Format Go source files with AST-aware import grouping and UTF-8 LF
gitmap aum format-go cli --write

# 9. Format only staged Go source files
gitmap aum format-go --staged --write
```

## Cross-References

- AUM Specification & Roadmap: `02-spec/21-app/128-aum-automation-suite-and-roadmap.md`
- Polyglot Worker Orchestrator: `02-spec/21-app/124-polyglot-worker-orchestrator-and-automation-runner.md`
- Automation LLM Guide: `02-spec/21-app/125-automation-llm-orchestration-guide.md`
- Coding Guidelines: `02-spec/02-coding-guidelines/01-cross-language/`
- See also: `gitmap find-files`, `gitmap find-regex-read`, `gitmap os`
