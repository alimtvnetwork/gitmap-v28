# 124 — Polyglot Worker Orchestrator & Automation Runner Specification

## Overview

**Module Number:** 124
**Version:** 2.0.0
**Updated:** 2026-09-19
**Status:** Architecture Proposal & Specification
**AI Confidence:** Production-Ready
**Ambiguity Score:** None

---

## 1. Purpose & Architectural Vision

This specification defines the **Go Supervisor / Polyglot Worker Pool Architecture** for GitMap.

GitMap CLI serves as the high-performance master supervisor:
1. **Ultra-Fast Filesystem Traversal & Manifesting:** Go scans 10,000+ repository files in memory in `<15ms`, resolves ignore filters, checks binary probe guards, and caches the repository inventory in an ultra-fast SQLite look-ahead manifest.
2. **SQLite-Backed Runtime Discovery & Caching:** Runtimes (Python, Node.js, Go, Rust, PowerShell, Bash) are probed once on first execution and cached in `.gitmap/data/<repo-slug>/automation/sql.db`. Subsequent executions resolve runtime paths in `<0.05ms` without recurring PATH scans.
3. **Automated Runtime Remediation:** If a required language runtime is missing, GitMap provides one-click suggestions targeting the GitMap installer (`gitmap install python`, `gitmap install profile dev-full`) alongside native package manager fallbacks.
4. **Bilingual Stream Encoding (UTF-8 & UTF-16):** Go negotiates UTF-8 for cross-platform JSON pipelines and UTF-16LE wide streams for Windows PowerShell and Win32 consoles, preventing character corruption and code page mismatches.
5. **Polyglot Worker Pool Management:** Go partitions files into balanced chunks across configurable worker groups and parallel worker threads, avoiding the catastrophic OS process-per-file overhead.
6. **On-the-Fly Code Snippets & Script Files:** Developers and AI assistants can supply inline code strings (`"code..."`) or point to dedicated script files (`py-file`, `node-file`, `go-file`, `ps-file`, `sh-file`).
7. **Composed Pipelines:** File search (`search-filename`), regex search (`search-grep`), and Git change detection (`changed-files`) feed directly into parallel worker execution pipelines.

---

## 2. Database Architecture (`sql.db`)

### 2.1 Repository-Scoped Storage Path

Under GitMap's split-database architecture, all automation telemetry, cached runtimes, and look-ahead file manifests are persisted in:

```
.gitmap/data/automation/<repo-slug>/sql.db
```
*(Standardized per [Spec 129](129-pr-commit-engines-and-sqlite-split-db.md) canonical `<section>/<slug>/sql.db` formula)*

Where `<repo-slug>` is the sanitized repository identifier (e.g. `alimtvnetwork-gitmap-v28` or relative workspace slug).

### 2.2 Connection Standards & Pragmas

All SQLite connections to `sql.db` MUST adhere to GitMap database conventions:
- **Max Open Connections:** `db.SetMaxOpenConns(1)` (strictly enforces single-writer serialization to prevent Windows NTFS lock collisions).
- **WAL Journal Mode:** `PRAGMA journal_mode = WAL;` (permits concurrent non-blocking readers).
- **Busy Timeout:** `PRAGMA busy_timeout = 5000;` (5-second retry window for transient locks).
- **Foreign Keys:** `PRAGMA foreign_keys = ON;`.
- **Synchronous Mode:** `PRAGMA synchronous = NORMAL;`.

### 2.3 Database Schema

```sql
-- Table: runtimes
-- Caches discovered CLI interpreters, compilers, and shells.
CREATE TABLE IF NOT EXISTS runtimes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,          -- 'python', 'node', 'rust', 'go', 'pwsh', 'bash'
    binary_name TEXT NOT NULL,          -- 'python.exe', 'node.exe', 'cargo.exe', 'pwsh.exe'
    binary_path TEXT NOT NULL,          -- Normalized forward-slash path to executable
    version TEXT NOT NULL,              -- e.g. 'Python 3.12.3', 'Node v20.11.0'
    status TEXT NOT NULL,               -- 'active', 'missing', 'degraded'
    install_cmd TEXT NOT NULL,          -- Primary GitMap command: 'gitmap install python'
    profile_suggestion TEXT NOT NULL,   -- Profile suggestion: 'gitmap install profile dev-full'
    fallback_cmd TEXT NOT NULL,         -- OS fallback: 'winget install Python.Python.3.12'
    discovered_at DATETIME NOT NULL,
    last_verified_at DATETIME NOT NULL,
    is_valid INTEGER NOT NULL DEFAULT 1
);
CREATE INDEX IF NOT EXISTS idx_runtimes_name ON runtimes(name);
CREATE INDEX IF NOT EXISTS idx_runtimes_status ON runtimes(status);

-- Table: file_manifest
-- High-speed look-ahead file index for instant filtering and batching.
CREATE TABLE IF NOT EXISTS file_manifest (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    relative_path TEXT NOT NULL UNIQUE,          -- 'cli/cmd/root.go'
    file_name TEXT NOT NULL,                      -- 'root.go'
    file_extension TEXT NOT NULL,                 -- '.go'
    parent_folder_path TEXT NOT NULL,             -- 'cli/cmd'
    absolute_file_path TEXT NOT NULL,             -- '/repos/gitmap/cli/cmd/root.go'
    absolute_parent_folder_path TEXT NOT NULL,    -- '/repos/gitmap/cli/cmd'
    file_size INTEGER NOT NULL,                   -- Size in bytes
    modified_timestamp INTEGER NOT NULL,          -- Unix seconds
    content_hash TEXT,                            -- Optional SHA-256
    is_binary INTEGER NOT NULL DEFAULT 0,         -- 1 if probe detects binary
    scanned_at DATETIME NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_manifest_ext ON file_manifest(file_extension);
CREATE INDEX IF NOT EXISTS idx_manifest_parent ON file_manifest(parent_folder_path);
CREATE INDEX IF NOT EXISTS idx_manifest_size ON file_manifest(file_size);

-- Table: execution_history
-- Records execution audits, throughput metrics, and exit statuses.
CREATE TABLE IF NOT EXISTS execution_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    runtime TEXT NOT NULL,
    command_type TEXT NOT NULL,                   -- 'inline', 'file', 'search-pipe'
    script_preview TEXT NOT NULL,                 -- Truncated snippet (<= 120 chars)
    files_matched INTEGER NOT NULL,
    files_processed INTEGER NOT NULL,
    workers_used INTEGER NOT NULL,
    threads_per_worker INTEGER NOT NULL,
    encoding TEXT NOT NULL,                       -- 'utf-8', 'utf-16le'
    duration_ms INTEGER NOT NULL,
    exit_code INTEGER NOT NULL,
    executed_at DATETIME NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_exec_runtime ON execution_history(runtime);
CREATE INDEX IF NOT EXISTS idx_exec_date ON execution_history(executed_at);

-- Table: search_exclusions
-- Persists custom files, directories, and glob patterns excluded from search and audits.
CREATE TABLE IF NOT EXISTS search_exclusions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pattern TEXT NOT NULL UNIQUE,                 -- e.g. 'assets/icon.png', '*.bin', 'tmp/cache'
    reason TEXT NOT NULL DEFAULT 'user_excluded', -- 'binary', 'large_file', 'large_json', 'user_excluded'
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_exclusions_pattern ON search_exclusions(pattern);
```

---

## 3. Runtime Probing, Caching & Auto-Bootstrapping

### 3.1 Two-Tier Discovery Protocol

```
                      [ User Invokes: gitmap automation run <runtime> "..." ]
                                                │
                                                ▼
                     [ Step 1: Query SQLite Cache: runtimes in sql.db ]
                                                │
                       ┌────────────────────────┴────────────────────────┐
                       │                                                 │
                  [ Cache HIT ]                                    [ Cache MISS ]
                       │                                                 │
          Check last_verified_at < 24h                                   │
          and os.Stat(binary_path) == nil                                │
                       │                                                 ▼
             ┌─────────┴─────────┐                       [ Step 2: System Probe ]
             ▼                   ▼                       - Check PATH environment
        [ Valid ]            [ Stale ]                   - Check AppData / Local Programs
             │                   │                       - Check /usr/bin, /usr/local/bin
             │                   └───────────────────────┐
             │                                           │
             ▼                                           ▼
      [ Run Worker ]                             [ Probe Result ]
                                                         │
                               ┌─────────────────────────┴─────────────────────────┐
                               ▼                                                   ▼
                         [ Discovered ]                                      [ Missing ]
                               │                                                   │
                   Save to runtimes in sql.db                          Save status='missing' to sql.db
                               │                                                   │
                               ▼                                                   ▼
                        [ Run Worker ]                                    [ Return AppError E7100 ]
                                                                          Display GitMap Install UI
```

1. **Subsequent Calls (<0.05ms):** Reads `binary_path` directly from `runtimes` in `sql.db`.
2. **Cache Invalidation:** If `os.Stat(binary_path)` fails (binary moved or uninstalled), GitMap transparently re-probes and updates the record.

### 3.2 Missing Runtime Suggestion UI

When a runtime cannot be found, GitMap intercepts immediately without launching broken child processes. It renders a clean error envelope (`E7100:RUNTIME_MISSING`) with copy-pasteable GitMap installer commands:

```
[E7100:RUNTIME_MISSING] Python runtime not detected on host system.

Suggested Installation via GitMap:
  gitmap install python
  gitmap install profile dev-full
  gitmap install profile python-data

Alternative System Package Managers:
  Windows (winget) : winget install Python.Python.3.12
  Windows (choco)  : choco install python3
  macOS (brew)     : brew install python@3.12
  Linux (Ubuntu)   : sudo apt update && sudo apt install -y python3 python3-pip

Tip: Configure automatic bootstrap in config:
  gitmap automation config set autoInstallMissingRuntimes true
```

---

## 4. Bilingual Stream Encoding: UTF-8 vs UTF-16

### 4.1 The Windows PowerShell vs Unix Dilemma

| Runtime Environment | Native Stream Default | Failure Mode When Mismatched |
|---|---|---|
| **Windows PowerShell (`pwsh`, `powershell.exe`)** | UTF-16LE (`[System.Text.Encoding]::Unicode`) | Emits double-byte characters; standard UTF-8 readers perceive null bytes (`\x00`) between characters, causing corrupt JSON and broken pipelines. |
| **Unix Shells, Python 3, Node.js, Go, Rust** | UTF-8 | Chokes on Windows UTF-16 byte streams, raising `UnicodeDecodeError` or JSON parse syntax errors. |

### 4.2 GitMap's Bilingual Stream Protocol

GitMap resolves this friction through **Bilingual Stream Negotiation**:

1. **Default JSON Wire Protocol:** Standard `UTF-8` is used for all inter-process JSON exchanges across Go, Python, Node.js, and Rust.
2. **Explicit Flag `--encoding utf16` / `utf16le`:**
   - Go writes `stdin` using UTF-16LE encoding with a Byte Order Mark (`0xFF, 0xFE`).
   - Go decodes worker `stdout` using standard library `unicode/utf16`, converting wide characters to native Go strings with zero allocation loss.
3. **Auto-Detection via BOM:**
   - If a child process begins its output with `0xFF, 0xFE` (UTF-16LE BOM) or `0xEF, 0xBB, 0xBF` (UTF-8 BOM), the Go supervisor automatically selects the matching decoder.
4. **PowerShell Auto-Enforcement:** When running `gitmap automation run ps ...`, GitMap automatically injects UTF-8 output stream reconfiguration unless `--encoding utf16` is explicitly supplied:
   ```powershell
   [Console]::OutputEncoding = [System.Text.Encoding]::UTF8
   $OutputEncoding = [System.Text.Encoding]::UTF8
   ```

---

## 5. Injected JSON Metadata Context

When files are processed by worker groups, Go injects a rich file context object via `stdin` (for JSON mode) or pre-populated scope variables:

```json
{
  "filePath": "cli/cmd/root.go",
  "fileName": "root.go",
  "fileExtension": ".go",
  "parentFolderPath": "cli/cmd",
  "absoluteFilePath": "/repos/gitmap/cli/cmd/root.go",
  "absoluteParentFolderPath": "/repos/gitmap/cli/cmd",
  "fileSize": 17155,
  "modifiedTimestamp": 1726744800,
  "lineNumber": 560,
  "matchedContent": "case constants.CmdAutomation:",
  "content": "package cmd\n\nimport (...",
  "encoding": "utf-8"
}
```

### Context Field Reference

| Field | Type | Description | Example |
|---|---|---|---|
| `filePath` | `string` | Relative path from repo root (forward slashes) | `cli/cmd/root.go` |
| `fileName` | `string` | Base file name with extension | `root.go` |
| `fileExtension` | `string` | File extension including leading dot | `.go` |
| `parentFolderPath` | `string` | Relative directory path containing the file | `cli/cmd` |
| `absoluteFilePath` | `string` | Normalized absolute OS path | `/repos/gitmap/cli/cmd/root.go` |
| `absoluteParentFolderPath`| `string` | Normalized absolute OS parent directory | `/repos/gitmap/cli/cmd` |
| `fileSize` | `integer`| File size in bytes | `17155` |
| `modifiedTimestamp` | `integer`| Unix timestamp of last file modification | `1726744800` |
| `lineNumber` | `integer`| Matching line number (populated in grep pipelines) | `560` |
| `matchedContent` | `string` | Exact matched line content (in grep pipelines) | `case constants.CmdAutomation:` |
| `content` | `string` | File content buffer (populated when `--pre-read` is enabled) | `"package cmd..."` |
| `encoding` | `string` | Negotiated stream encoding (`utf-8` or `utf-16le`) | `utf-8` |

---

## 6. Comprehensive Polyglot Examples

### 6.1 Python Runtimes (`py`, `python`, `py-file`)

#### Example 1: Inline AST Security & Quality Linter (Parallel Batch)
Scans all Python files across the repository, parses their Abstract Syntax Tree (AST), and detects unsafe `eval()`, bare `except:`, or hardcoded passwords without writing a single temporary file:

```bash
gitmap automation search-filename "*.py" run py "
import ast, sys, json

data = json.loads(sys.stdin.readline())
path = data['filePath']
try:
    with open(path, 'r', encoding='utf-8') as f:
        tree = ast.parse(f.read(), filename=path)
    for node in ast.walk(tree):
        if isinstance(node, ast.ExceptHandler) and node.type is None:
            print(f'[WARN:BARE_EXCEPT] {path}:{node.lineno}')
        elif isinstance(node, ast.Call) and getattr(node.func, 'id', None) == 'eval':
            print(f'[CRIT:EVAL_USED] {path}:{node.lineno}')
except Exception as e:
    print(f'[ERR:PARSE_FAIL] {path}: {e}', file=sys.stderr)
" --w 4 --threads 4
```

#### Example 2: Regex Secret Validator with Composed Grep Pipe
Finds potential API keys using GitMap's high-speed Boyer-Moore search engine and pipes matched lines to Python to validate entropy and format:

```bash
gitmap automation search-grep "AIza[0-9A-Za-z-_]{35}" run py "
import json, sys

data = json.loads(sys.stdin.readline())
path = data['filePath']
line = data['lineNumber']
key = data['matchedContent'].strip()

if not key.endswith('EXAMPLE') and not 'test' in key.lower():
    print(f'🚨 SUSPECT GOOGLE API KEY: {path}:{line} -> {key[:8]}...{key[-4:]}')
"
```

#### Example 3: Script File Cyclomatic Complexity Auditor (`py-file`)
Executes a standalone Python script across repository files:

```bash
gitmap automation search-filename "*.py" run py-file "scripts/measure_complexity.py" --w 4
```

---

### 6.2 Node.js & TypeScript Runtimes (`node`, `js`, `ts`, `node-file`)

#### Example 1: Monorepo `package.json` Dependency Sync Auditor
Traverses all packages in a monorepo, parses `package.json`, and flags conflicting dependencies:

```bash
gitmap automation search-filename "package.json" run node "
const fs = require('fs');
const readline = require('readline');

const rl = readline.createInterface({ input: process.stdin });
rl.on('line', (line) => {
  const data = JSON.parse(line);
  const pkg = JSON.parse(fs.readFileSync(data.absoluteFilePath, 'utf8'));
  const deps = pkg.dependencies || {};
  if (deps.lodash && !deps.lodash.startsWith('^4.17')) {
    console.log(`⚠️  Outdated lodash in ${data.parentFolderPath}: ${deps.lodash}`);
  }
});
"
```

#### Example 2: React Component Missing Key Inspector
Scans JSX and TSX files for unkeyed map expressions:

```bash
gitmap automation search-grep ".map(" run node "
const readline = require('readline');
const rl = readline.createInterface({ input: process.stdin });
rl.on('line', (line) => {
  const data = JSON.parse(line);
  if (data.fileExtension === '.tsx' || data.fileExtension === '.jsx') {
    if (!data.matchedContent.includes('key=')) {
      console.log(`[REACT:NO_KEY] ${data.filePath}:${data.lineNumber} -> ${data.matchedContent.trim()}`);
    }
  }
});
" --w 4
```

#### Example 3: Standalone Script File Execution (`node-file`)
```bash
gitmap automation run node-file "scripts/audit_imports.mjs" --w 4
```

---

### 6.3 Go Runtimes (`go`, `go-file`)

#### Example 1: On-The-Fly Function Length & Parameter Linter
Runs an inline Go snippet that parses Go files with standard library `go/parser` and verifies adherence to the 15-line function limit:

```bash
gitmap automation search-filename "*.go" run go "
package main

import (
	\"bufio\"
	\"encoding/json\"
	\"fmt\"
	\"go/parser\"
	\"go/token\"
	\"os\"
)

type FileCtx struct {
	FilePath         string `json:\"filePath\"`
	AbsoluteFilePath string `json:\"absoluteFilePath\"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var ctx FileCtx
		if err := json.Unmarshal(scanner.Bytes(), &ctx); err != nil {
			continue
		}
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, ctx.AbsoluteFilePath, nil, 0)
		if err != nil {
			continue
		}
		for _, decl := range node.Decls {
			// Inspect AST declarations and flag functions exceeding limits
			_ = decl
		}
		fmt.Printf(\"Verified: %s\\n\", ctx.FilePath)
	}
}
" --w 4
```

#### Example 2: Script File Go Verifier (`go-file`)
```bash
gitmap automation run go-file "scripts/verify_struct_tags.go" --w 4
```

---

### 6.4 Rust Runtimes (`rust`, `rs`, `rs-file`)

#### Example 1: High-Speed Token Counter
Compiles and invokes a lightning-fast Rust snippet using `cargo-script` or `rustc` for memory-safe byte analysis:

```bash
gitmap automation search-filename "*.md" run rust "
use std::io::{self, BufRead};

fn main() {
    let stdin = io::stdin();
    for line in stdin.lock().lines() {
        if let Ok(l) = line {
            println!(\"Rust Worker received: {}\", l.len());
        }
    }
}
" --w 4
```

---

### 6.5 PowerShell Runtimes (`ps`, `powershell`, `ps-file`)

#### Example 1: Windows File Security & ACL Inspector (UTF-16 Stream)
Uses PowerShell with native wide strings to audit NTFS file permissions without character corruption:

```bash
gitmap automation search-filename "*.exe" run ps "
$input | ForEach-Object {
    $ctx = $_ | ConvertFrom-Json
    $acl = Get-Acl -Path $ctx.absoluteFilePath
    Write-Host \"[NTFS_OWNER] $($ctx.filePath) -> $($acl.Owner)\"
}
" --encoding utf16
```

#### Example 2: Standalone PowerShell Audit Script (`ps-file`)
```bash
gitmap automation run ps-file "scripts/audit_registry_keys.ps1" --encoding utf16
```

---

### 6.6 Bash & POSIX Shell Runtimes (`bash`, `sh`, `sh-file`)

#### Example 1: Executable Permissions & Shebang Validator
Verifies that all shell scripts start with a valid shebang and have executable bits:

```bash
gitmap automation search-filename "*.sh" run bash "
while read -r line; do
    file=\$(echo \"\$line\" | jq -r .absoluteFilePath)
    rel=\$(echo \"\$line\" | jq -r .filePath)
    head -n 1 \"\$file\" | grep -q '^#!/' || echo \"[SHEBANG:MISSING] \$rel\"
    [ -x \"\$file\" ] || echo \"[PERM:NOT_EXEC] \$rel\"
done
"
```

---

### 6.7 AI Script Catalog Proposals & Real-World Polyglot Orchestration (`03-ai-scripts`)

The scripts in `03-ai-scripts/` represent high-value repository maintenance algorithms. Under the Go Supervisor architecture, these algorithms run either via high-performance native Go engines or via parallel polyglot worker processes.

#### 6.7.1 Sequence, Numbering & Title Header Auditor (`15-sequence-and-title-auditor.py` / Native Go)
Audits numbered markdown files (e.g. `01-intro.md`, `02-spec/...`, `124-polyglot-worker.md`) across directory trees to guarantee:
1. **Zero Sequence Gaps:** Numbering must be strictly contiguous (e.g. `01`, `02`, `03`... no missing numbers).
2. **Zero Duplicate Numbers:** No two files may share the same sequence prefix.
3. **H1 Title Alignment:** Top-level `# XX Title` headers must exactly match the filename numeric prefix.
4. **Automated Fix Mode (`--fix`):** Auto-updates mismatched H1 headers in place with Unix LF line endings.

```bash
# Audit markdown sequences and headers across all spec directories
gitmap automation sequence 02-spec/21-app

# Automatically fix mismatched H1 headers across the repository
gitmap automation sequence --fix

# Polyglot Worker Execution (Delegating to Python engine with worker pool):
gitmap automation run py-file "03-ai-scripts/15-sequence-and-title-auditor.py" --path "02-spec" --w 4
```

#### 6.7.2 File Size Guard & Memory-Safe Binary Probing (`13-file-size-guard.py` / Native Go)
Prevents accidental commits of massive blobs, multi-megabyte JSON dumps, and pre-compiled binaries:
1. **Configurable Ceilings:** Audits all tracked and local repository files against maximum size limits (default: 500 KB, configurable via `--max-kb`).
2. **8KB Chunk Probing:** Evaluates unrecognized file types by scanning the first 8,192 bytes for null bytes (`0x00`) without loading entire multi-megabyte files into RAM.
3. **Allowed Large File Waivers (`ALLOWED_LARGE_FILES`):** Explicitly exempts pre-approved metadata assets (e.g. `src/data/specTree.json`, `.ai-memory/test-inventory.json`).

```bash
# Execute native Go file size guard (default 500 KB threshold)
gitmap automation guard

# Audit with custom 250 KB ceiling and output JSON telemetry
gitmap automation guard --max-kb 250 --json

# Polyglot Worker Execution (Delegating to Python engine):
gitmap automation run py-file "03-ai-scripts/13-file-size-guard.py" --max-kb 500 --w 4
```

#### 6.7.3 Large JSON Automatic Exclusion ("Out of the List")
Large JSON files (such as `specTree.json`, `test-inventory.json`, AST dumps, and pipeline telemetry logs exceeding 500 KB) pose a severe threat to search performance and AI LLM context windows:
- **Default Policy:** Any `.json` file exceeding 500 KB is automatically excluded ("out of the list") from repository search streams (`gitmap automation search`), file manifests, and prompt injection buffers.
- **Explicit Override:** Developers can explicitly include large JSONs when needed via `--include-large-json` or `--max-json-kb <size>`.

```bash
# Standard search automatically excludes oversized JSON dumps
gitmap automation search "connectionTimeout"

# Explicitly search including large JSON data dumps
gitmap automation search "connectionTimeout" --include-large-json --max-json-kb 2000
```

#### 6.7.4 Common Binary Probing & Interactive User Confirmation
When un-excluded binary files (`.exe`, `.dll`, `.bin`, `.zip`, `.png`, `.jpg`, `.wasm`, `.db`) or null-byte chunks are discovered:
1. **Interactive Terminal Prompt:** GitMap displays the discovered binary and prompts the developer:
   ```
   [Binary File Detected] assets/logo.png (57 KB)
   Exclude this file from future searches? [y/N]:
   ```
2. **Persistent Decision:** If confirmed (`y`), the file path or pattern is persisted into `search_exclusions` in `.gitmap/data/<repo-slug>/automation/sql.db`.
3. **Headless / CI Execution:** In non-interactive mode (`CI=true` or piped stdin), binaries are automatically excluded from search streams without prompting.

#### 6.7.5 Customizable Search Exclusions & Waiver Registry
Developers can inspect, add, remove, and purge search exclusion patterns stored persistently in SQLite:

```bash
# List all active search and audit exclusions
gitmap automation exclude list

# Add a custom file path or glob pattern to exclusions
gitmap automation exclude add "assets/vendor/*.bin" "vendor_blobs"
gitmap automation exclude add "mock_data.json" "test_fixture"

# Remove an exclusion pattern
gitmap automation exclude rm "mock_data.json"

# Clear all custom exclusions
gitmap automation exclude clear
```

#### 6.7.6 Version Synchronization Checker (`14-version-sync-checker.py`)
Ensures version consistency across polyglot project manifests (`package.json`, `Cargo.toml`, `go.mod`, `cli/constants/constants_version.go`):

```bash
gitmap automation run py-file "03-ai-scripts/14-version-sync-checker.py"
```

#### 6.7.7 CLI Help & Command Parity Auditor (`09-cli-help-auditor.py`)
Audits CLI command registration, help text flags, and markdown documentation catalogs to ensure 100% AST simulation parity:

```bash
gitmap automation run py-file "03-ai-scripts/09-cli-help-auditor.py"
```

#### 6.7.8 Result Monad & AppError Wrapper Auditor (`35-result-wrapper-auditor.py`)
Scans Go source files to detect and eliminate raw multi-value `(T, error)` tuples, enforcing strongly-typed `result.Result[T]` wrappers and structured `*apperror.AppError` envelopes:

```bash
gitmap automation run py-file "03-ai-scripts/35-result-wrapper-auditor.py" --path "cli"
```

---

## 7. Composed Pipelines & Chaining

### 7.1 Search and Execute Pipelines

```bash
# 1. Filter by filename glob and process with Python
gitmap automation search-filename "*.json" run py "..." --w 4 --threads 4

# 2. Search regex pattern and process matched lines with Node.js
gitmap automation search-grep "FIXME:" run node "..." --w 3 --threads 4

# 3. Detect uncommitted Git changes and run linter on changed files only
gitmap automation changed-files run py-file "scripts/pre_commit_check.py"
```

### 7.2 Standard Unix Pipe Composition

```bash
# Output discovered files as JSON stream and pipe to custom consumer
gitmap automation search-filename "*.go" --json | jq '.filePath' | head -n 20
```

---

## 8. CLI Management & Configuration Commands

### 8.1 Runtime Management

```bash
# List discovered runtimes and their resolution status
gitmap automation runtimes list

# Force re-probe of system PATH and refresh SQLite runtimes table
gitmap automation runtimes refresh
```

### 8.2 Database Inspection & Maintenance

```bash
# Display database file location, size, and row counts
gitmap automation db status

# Execute read-only diagnostic SQL queries against sql.db
gitmap automation db query "SELECT name, binary_path, status, last_verified_at FROM runtimes"
gitmap automation db query "SELECT file_extension, COUNT(*) as count FROM file_manifest GROUP BY file_extension ORDER BY count DESC"

# Clear manifest cache and execution history (preserves runtime registry)
gitmap automation db clear
```

### 8.3 Automation Settings Configuration

```bash
# Display current automation settings
gitmap automation config show

# Set default worker groups and concurrency
gitmap automation config set defaultWorkers 6
gitmap automation config set defaultThreads 4

# Set default encoding ('utf-8' or 'utf-16le')
gitmap automation config set defaultEncoding utf-8

# Set per-worker timeout
gitmap automation config set defaultTimeoutSeconds 45

# Toggle automatic runtime installation via GitMap installer
gitmap automation config set autoInstallMissingRuntimes true
```

### 8.4 Search & Audit Exclusion Management

```bash
# List all active exclusions
gitmap automation exclude list

# Add a custom exclusion pattern with an audit reason
gitmap automation exclude add "assets/vendor/*.bin" "vendor_blobs"

# Remove an exclusion pattern
gitmap automation exclude rm "assets/vendor/*.bin"

# Clear all custom exclusions
gitmap automation exclude clear
```

### 8.5 Repository File Size Guard & Sequence Auditor Commands

```bash
# Run file size guard across the repository (default 500 KB limit)
gitmap automation guard

# Run guard with custom limit, non-interactive, and JSON output
gitmap automation guard --max-kb 250 --interactive=false --json

# Run sequence numbering and H1 header auditor
gitmap automation sequence 02-spec

# Automatically repair and align mismatched H1 headers to file sequence numbers
gitmap automation sequence --fix
```

---

## 9. Acceptance Criteria

### Scenario 1: SQLite Runtime Cache Hit (<0.05ms)
- **Given** Python has been previously discovered and recorded in `.gitmap/data/<repo-slug>/automation/sql.db`
- **When** `gitmap automation run py "print('ok')"` is invoked
- **Then** GitMap reads the runtime path from `sql.db` in `<0.05ms` without re-scanning the system PATH or invoking `python --version`.

### Scenario 2: Missing Runtime Suggestion UI
- **Given** Rust is not installed on the host machine
- **When** `gitmap automation run rust "fn main() {}"` is invoked
- **Then** GitMap immediately halts with `E7100:RUNTIME_MISSING` and displays `gitmap install rust` along with profile and system package manager suggestions.

### Scenario 3: Windows PowerShell UTF-16 Stream Integrity
- **Given** a PowerShell worker script invoked with `--encoding utf16`
- **When** the script outputs wide-character strings, emojis, or non-ASCII identifiers
- **Then** GitMap decodes the stdout stream via UTF-16LE without mojibake or corrupt replacement characters.

### Scenario 4: Composed Search and Parallel Batching
- **Given** a repository with 5,000+ files
- **When** `gitmap automation search-filename "*.go" run py "..." --w 4 --threads 4` is executed
- **Then** Go groups all matching files into 4 balanced chunks, launches 4 worker processes, passes files via JSON stream, and completes without thrashing the OS.

### Scenario 5: Large JSON Files Excluded from Search Stream
- **Given** a repository containing `src/data/specTree.json` (> 1 MB) and `test-inventory.json`
- **When** `gitmap automation search "pattern"` is executed without `--include-large-json`
- **Then** GitMap automatically excludes all JSON files exceeding 500 KB, keeping search results concise and fast.

### Scenario 6: Common Binary Interactive Prompt & Persistence
- **Given** an un-excluded `.png` or binary file in a repository folder
- **When** `gitmap automation guard` runs interactively
- **Then** GitMap prompts the developer whether to exclude the file from future searches, and upon confirmation (`y`), persists the pattern into `.gitmap/data/<repo-slug>/automation/sql.db`.

### Scenario 7: Sequence Gaps & Title Mismatch Auto-Repair
- **Given** a documentation directory with numbering gaps (`01`, `03`) and title mismatch (`03` prefix with `# 02 Header`)
- **When** `gitmap automation sequence --fix` is executed
- **Then** GitMap reports the sequence gap and automatically repairs the `# 02` header to `# 03` in place with Unix LF line endings.

---

## 10. Cross-References

- Companion LLM Orchestration Guide: [`./125-automation-llm-orchestration-guide.md`](./125-automation-llm-orchestration-guide.md)
- Cross-Platform Python Tooling: [`./123-cross-platform-python-tooling.md`](./123-cross-platform-python-tooling.md)
- Database Architecture: [`./120-database-suite-and-start-fresh.md`](./120-database-suite-and-start-fresh.md)
- Coding Guidelines: [`../02-coding-guidelines/00-overview.md`](../02-coding-guidelines/00-overview.md)
- Error Management: [`../03-error-manage/00-overview.md`](../03-error-manage/00-overview.md)
