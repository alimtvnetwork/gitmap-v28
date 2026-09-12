# Python Script DRY Architecture, Enums, Pluggable Caching & Fast AI Reading Specification

> **Prompt Version:** 2.1.0
> **Target:** `01-prompts/08-dry-code/01-python-dry-architecture-and-caching.md`
> **Synchronization:** Meta-Repo & AI Scripting Ecosystem

/goal Standardize the architectural design of all Python CI/CD, linting, and fix scripts using strict Python Enum conventions (PascalCase class name ending in Type, UPPER_CASE members, string values mirroring member names), centralized configuration maps, thread-safe lazy regex compilation with pre-initialized None mapping, zero magic strings/numbers, shallow guard clauses (ban nested if pyramids), fast substring pre-filtering, DRY shared engines, multi-folder scoping, customizable extensions, pluggable `tmp/cache/` storage, cross-process atomic file locking, and two-phase incremental `mtime` caching.

## 🎯 Architectural Philosophy

All AI-authored Python scripts in this repository (`03-ai-scripts/` and `.agents/scripts/`) MUST strictly adhere to the following principles:

---

## 1. Python Enum Generation Standard (Strict Standard)

When generating Python `Enum` classes, you must strictly follow these formatting and naming conventions:

1. **Class Name (The Enum Type):**
   - Convention: Use `PascalCase` (UpperCamelCase).
   - Grammar: Use singular nouns ending with `Type` (e.g. `RegexPatternType`, `ScanModeType`, `SeverityType`, `ExitCodeType`, `ArtifactCategoryType`, `CacheKeyType`).
2. **Enum Members (The Variable Names):**
   - Convention: Use `UPPER_CASE` with underscores separating words (`SNAKE_CASE` in all caps).
   - Example: `WINDOWS_BACKSLASH`, `PENDING_APPROVAL`, `SUCCESS`.
3. **Enum Values (The Underlying Data):**
   - Convention: The literal value assigned to the member must be a string.
   - Formatting: The string value must perfectly match the member name in `UPPER_CASE` with underscores.
   - Example: `WINDOWS_BACKSLASH = "WINDOWS_BACKSLASH"`, `PAYMENT_PENDING = "PAYMENT_PENDING"`.

```python
from enum import Enum

class ScanModeType(str, Enum):
    """Enumeration for file scanning modes."""
    CHECK = "CHECK"
    FIX = "FIX"
    STREAM = "STREAM"

class SeverityType(str, Enum):
    """Enumeration for issue severity levels."""
    BLOCKER = "BLOCKER"
    HIGH = "HIGH"
    WARN = "WARN"
    INFO = "INFO"

class ExitCodeType(int, Enum):
    """Enumeration for application exit codes."""
    SUCCESS = 0
    VIOLATIONS_FOUND = 1
    TOOL_ERROR = 2

class ArtifactCategoryType(str, Enum):
    """Enumeration for cleanup artifact categories."""
    PYCACHE = "PYCACHE"
    BINARIES = "BINARIES"
    TEMPORARY = "TEMPORARY"
    CUSTOM = "CUSTOM"
    ALL = "ALL"

class RegexPatternType(str, Enum):
    """Enumeration for cached regex pattern identifiers."""
    WINDOWS_BACKSLASH = "WINDOWS_BACKSLASH"
    LEADING_DOT_SLASH = "LEADING_DOT_SLASH"
    CRLF = "CRLF"
    UNIVERSAL_LINE_ENDING = "UNIVERSAL_LINE_ENDING"
    TRAILING_WHITESPACE = "TRAILING_WHITESPACE"
    SEQ_PREFIX = "SEQ_PREFIX"
    UPPERCASE = "UPPERCASE"
    FILE_URI_WIN = "FILE_URI_WIN"
    DRIVE_ABS_WIN = "DRIVE_ABS_WIN"
    REPO_FILE_URI = "REPO_FILE_URI"
    EXPLICIT_DOUBLE_TRUE = "EXPLICIT_DOUBLE_TRUE"
    EXPLICIT_TRIPLE_TRUE = "EXPLICIT_TRIPLE_TRUE"
    EXPLICIT_PYTHON_TRUE = "EXPLICIT_PYTHON_TRUE"
    COMMENT_PREFIX = "COMMENT_PREFIX"
    COBRA_COMMAND = "COBRA_COMMAND"
    SHORT_DESC = "SHORT_DESC"
    EXAMPLE_USAGE = "EXAMPLE_USAGE"
    CHANGELOG_HEADER = "CHANGELOG_HEADER"
    FILE_NUM_PREFIX = "FILE_NUM_PREFIX"
    H1_HEADER = "H1_HEADER"
    PLACEHOLDER_TOKEN = "PLACEHOLDER_TOKEN"
    NON_ALPHANUMERIC = "NON_ALPHANUMERIC"
```

---

## 2. Centralized Configuration Maps & Lazy Regex Registry

- **Root Definition:** All configuration parameters (`EXCLUDE_DIRS`, `BINARY_EXTENSIONS`, `DEFAULT_TEXT_EXTENSIONS`, `DEFAULT_CODE_EXTENSIONS`, `DEFAULT_CLI_EXTENSIONS`, `CACHE_BASE_DIR`, `DEFAULT_MAX_FILE_KB`, `ALLOWED_LARGE_FILES`, `CI_JOBS_MATRIX`, `INSTALLER_*`) must be centralized in `02-shared-engine.py`.
- **Zero Magic Strings / Numbers:** Never hardcode raw byte tokens (`b"\x00"`, `b"\xef\xbb\xbf"`), chunk sizes (`8192`), path markers (`"\\\\?\\"`, `"."`), or string literals directly in child functions. Import them from `02-shared-engine.py`.
- **Implicit Boolean Checks:** Never compare booleans against explicit `True` (BAN: `if is_valid == True:` -> MANDATORY: `if is_valid:`).
- **Thread-Safe Lazy Regex Compilation (Zero Import-Time Compilation):**
  - Store raw regex pattern string definitions + compilation flags in a central `REGEX_DEFINITIONS: dict[RegexPatternType, tuple[str, int]]` map in `02-shared-engine.py`.
  - Initialize a mapping `_compiled_patterns: dict[RegexPatternType, re.Pattern | None] = {pt: None for pt in RegexPatternType}` with `None` values on module import.
  - Compile patterns lazily on first demand inside `RegexRegistry.get(pattern_type)` with double-checked `threading.Lock()`.
  - Regex patterns are compiled on first demand, then served in `O(1)` time for all subsequent lookups across threads.

```python
# --- Thread-Safe Lazy Regex Registry (Zero import-time compilation) ---
class RegexRegistry:
    """
    Thread-safe lazy-compiling regex registry.
    Initializes all entries to None mapping on startup.
    Patterns are compiled on-demand upon first get() call and cached.
    """
    _compiled_patterns: dict[RegexPatternType, re.Pattern | None] = {pt: None for pt in RegexPatternType}
    _lock = threading.Lock()

    @classmethod
    def get(cls, pattern_type: RegexPatternType) -> re.Pattern:
        """Lazily compiles on first demand and returns cached immutable re.Pattern."""
        cached = cls._compiled_patterns.get(pattern_type)
        if cached is not None:
            return cached

        with cls._lock:
            if cls._compiled_patterns[pattern_type] is None:
                if pattern_type not in REGEX_DEFINITIONS:
                    raise KeyError(f"Pattern type '{pattern_type}' is not registered in REGEX_DEFINITIONS")
                raw_pattern, flags = REGEX_DEFINITIONS[pattern_type]
                cls._compiled_patterns[pattern_type] = re.compile(raw_pattern, flags)
            return cls._compiled_patterns[pattern_type]
```

---

## 3. Shallow Guard Clauses & Fast Substring Pre-Filtering

- **Total Ban on Deeply Nested If-Blocks:** Avoid 4+ levels of nested indentation pyramids (`if a: if b: if c:`). Replace with flat guard clauses (`if not a: continue`).
- **Fast Substring Pre-Filtering:** Before invoking expensive AST parsers (`ast.parse()`) or running complex multi-line regex engines across hundreds of files, perform a quick substring check (e.g. `if "command" not in content: return []`, `if "file:" not in content and ":\\" not in content: return []`). This yields immediate **10x–50x speedups** on large codebases.

---

## 4. Multi-Folder Scoping & Customizable Extensions

- **Multi-Folder Capability:** All scripts MUST accept a target directory (`<path>` or `--path` / `--dir`) allowing them to run on the full repository, specific submodules, or individual feature folders.
- **Customizable File Extensions:** Supported file extensions must be customizable via `--ext` or function parameters (`extensions=...`), with robust lowercasing and leading dot normalization.
- **Nested Directory Pruning:** When using `os.walk`, prune `dirs[:]` dynamically so nested `.git`, `.gitmap`, and `node_modules` folders inside subprojects are skipped instantly.

---

## 5. Decomposed Pure Functions (< 25 Lines Each)

- Monolithic scripts are strictly forbidden.
- Scripts MUST be broken down into small, composable, single-responsibility functions.
- Separate file I/O, cache state verification, violation analysis, and CLI output formatting into distinct functions.

---

## 6. Pluggable `tmp/cache/` Structure & Cross-Process Locking

All cache data is organized in structured, pluggable subdirectories under `tmp/cache/`:

1. **`tmp/cache/paths/`**: Stores repository file path listings, topology discoveries, and metadata indexes.
2. **`tmp/cache/locks/`**: Cross-process file locks (`repo-cache.lock`) preventing corruption when multiple agents or subagents operate simultaneously.
3. **`tmp/cache/files/`**: Cached tokenized contents or AST data.

### Dual-Platform Locking & Stale Lock Recovery
- **POSIX (Linux/macOS):** Native kernel `fcntl.flock(LOCK_EX | LOCK_NB)` automatically cleaned up on crash or process kill.
- **Windows (NTFS):** Atomic `os.O_CREAT | os.O_EXCL` with PID timestamp and automatic eviction for stale locks (>15s).

---

## 7. Two-Phase Incremental Caching & Missing File Tolerance

To achieve sub-15ms repository-wide execution without redundant disk I/O:

1. **Phase 1 (Cache-First Processing):**
   - Read pre-computed file metadata (`mtime`, `size`) from `tmp/cache/repo-file-cache.json`.
   - Deliver cached files immediately (<0.1ms) so the consumer starts processing file #1 without waiting.
2. **Phase 2 (Streaming Discovery for New / Modified Files):**
   - Concurrently stream directory entries via `os.walk` / `os.scandir`.
   - On Unix, track `(st_dev, st_ino)` to prevent symlink recursion cycles.
3. **Fault-Tolerant File Reading (Zero Crash on Deletions):**
   - If a file is deleted or missing during traversal, `read_file_safe()` returns `None` instead of raising an unhandled exception.
   - Deleted files are automatically evicted from the cache without treating the deletion as a failure.

---

## 8. High-Speed File Reading & Exploration for AI Agents

AI agents and subagents should avoid slow, recursive shell commands (`Get-ChildItem -Recurse`, `dir /s`, or brute-force glob searches). Instead, use the optimized Python toolchain:

| Task | Recommended AI Command | Speed |
|---|---|:---:|
| **List Folder Files** | `python 03-ai-scripts/17-fast-file-reader.py --list-folder <dir>` | **<1ms** |
| **Fast Safe File Read** | `python 03-ai-scripts/17-fast-file-reader.py --read-file <path>` | **<1ms** |
| **Search File Paths** | `python 03-ai-scripts/17-fast-file-reader.py --search-pattern "<term>"` | **<2ms** |
| **Full Repo File Index** | `python 03-ai-scripts/11-fast-file-scanner.py --lang ts,go --path spec/` | **~14ms** |
| **Parallel Content Grep** | `python 03-ai-scripts/12-fast-cached-grep.py --pattern "<text>"` | **~12ms** |
| **File Manipulation CLI** | `python 03-ai-scripts/03-file-manipulator.py <cmd> <dir>` | **~15ms** |
| **Codebase Topology Routing**| `python 03-ai-scripts/18-codebase-topology-discoverer.py --query <db\|backend\|go>` | **<1ms** |
| **Safe Artifact Removal**| `python 03-ai-scripts/19-artifact-remover.py --clean-pycache --dry-run` | **~10ms** |
| **Plan Consolidation**| `python 03-ai-scripts/20-plan-consolidator.py --dry-run` | **~12ms** |

> [!IMPORTANT]
> **If any fast reader script (e.g. `17-fast-file-reader.py`) is missing:** You must immediately recreate it using Python standard libraries (`os`, `sys`, `pathlib`, `re`, `argparse`, `json`). It MUST support flags: `--list-folder <path> [--ext .md,.ts]`, `--read-file <path> [--max-bytes N]`, and `--search-pattern "<regex>" [--path <dir>]`. Ensure strict UTF-8 output (`sys.stdout.reconfigure(encoding="utf-8")`) and implement local caching.
