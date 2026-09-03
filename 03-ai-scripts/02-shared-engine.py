#!/usr/bin/env python3
"""
Shared Core Engine for AI Repository Tooling, CI Fix Scripts & High-Speed Caching
Dual-Platform Engine (100% Native Unix & Windows Support)

Features:
1. Centralized Configuration Maps, Language Manifests, Subsystem Hints, and Top-Level Enums:
   - EncodingType, LanguageType, SubsystemType, ArtifactCategoryType, ScanModeType, SeverityType, ExitCodeType, RegexPatternType, CacheKeyType.
2. Centralized Encodings, Separators, Tokens, Paths, Default Literals (CURRENT_DIR, CACHE_KEY_FILES), and Artifact Clean Presets.
3. Thread-Safe Lazy Regex Registry (Pre-initialized None Map, Dynamic On-The-Fly Auto-Registration with Logging).
4. Dual-Mode Cross-Process Locking:
   - POSIX: Kernel-level `fcntl.flock` (automatic cleanup on process kill/crash).
   - Windows: Atomic `os.O_CREAT | os.O_EXCL` with PID timestamp & stale-lock recovery.
5. Unix Symlink & Cycle Guard: Inode tracking (st_dev, st_ino) preventing infinite recursion.
6. Unix Permission Preservation: Preserves executable bits (chmod +x / st_mode) across atomic writes.
7. Universal Line Ending Normalizer: Aggressively converts CRLF (\\r\\n) and legacy Mac CR (\\r) to clean UNIX LF (\\n).
8. Memory-Safe Chunked Binary Probe: Inspects first 8KB for null-bytes without loading large blobs into RAM.
9. Two-phase incremental mtime-based file streaming (cache-first + parallel scan).
10. Pluggable cache layout in tmp/cache/ (paths, locks, files).
11. Fault-tolerant file reader handling missing/deleted files gracefully (zero crash).
12. Shared formatting and collection joining helpers (format_comma_separated, format_keys).
"""

from collections.abc import Generator, Iterable
from contextlib import contextmanager
from enum import Enum
import json
import os
from pathlib import Path
import re
import sys
import threading
import time
from typing import Any, Callable

# Optional POSIX kernel locking
try:
    import fcntl
    IS_FCNTL_AVAILABLE = True
except ImportError:
    fcntl = None
    IS_FCNTL_AVAILABLE = False

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")

# --- Top-Level Enums Following Standard ---

class EncodingType(str, Enum):
    """Enumeration for standardized character encodings."""
    UTF8 = "utf-8"
    UTF8_SIG = "utf-8-sig"
    UTF16 = "utf-16"
    UTF16_LE = "utf-16le"
    UTF16_BE = "utf-16be"
    ASCII = "ascii"

class LanguageType(str, Enum):
    """Enumeration for detected programming languages."""
    GO = "GO"
    RUST = "RUST"
    PYTHON = "PYTHON"
    TYPESCRIPT = "TYPESCRIPT"
    JAVASCRIPT = "JAVASCRIPT"
    PHP = "PHP"
    CSHARP = "CSHARP"
    SQL = "SQL"
    SHELL = "SHELL"
    MARKDOWN = "MARKDOWN"
    HTML = "HTML"
    CSS = "CSS"
    OTHER = "OTHER"

class SubsystemType(str, Enum):
    """Enumeration for major codebase subsystems."""
    BACKEND = "BACKEND"
    DATABASE = "DATABASE"
    FRONTEND = "FRONTEND"
    CICD = "CICD"
    DOCS = "DOCS"
    CLI = "CLI"
    TESTS = "TESTS"
    UNKNOWN = "UNKNOWN"

class ArtifactCategoryType(str, Enum):
    """Enumeration for cleanup artifact categories."""
    PYCACHE = "PYCACHE"
    BINARIES = "BINARIES"
    TEMPORARY = "TEMPORARY"
    CUSTOM = "CUSTOM"
    ALL = "ALL"

class CacheKeyType(str, Enum):
    """Enumeration for standardized cache and topology JSON keys."""
    FILES = "files"
    TIMESTAMP = "timestamp"
    TOTAL_FILES = "total_files"
    VERSION = "version"
    GENERATED_AT = "generatedAt"
    EXPIRES_AT = "expiresAt"
    TTL_SECONDS = "ttlSeconds"
    SCAN_DURATION_MS = "scanDurationMs"
    ROOT_PATH = "rootPath"
    MANIFESTS = "manifests"
    LANGUAGE_DISTRIBUTION = "languageDistribution"
    LANGUAGE_ROOTS = "languageRoots"
    SUBSYSTEMS = "subsystems"
    ROOTS = "roots"
    ENTRYPOINTS = "entrypoints"
    SCHEMA_FILES = "schemaFiles"
    WORKFLOWS = "workflows"
    SPEC_ROOTS = "specRoots"
    TEST_RUNNERS = "testRunners"

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

# --- Centralized Constants for Encodings, Separators, Tokens & Default Literals ---
DEFAULT_ENCODING = EncodingType.UTF8.value
UTF8_SIG_ENCODING = EncodingType.UTF8_SIG.value
UTF16_ENCODING = EncodingType.UTF16.value
UTF16_LE_ENCODING = EncodingType.UTF16_LE.value
UTF16_BE_ENCODING = EncodingType.UTF16_BE.value

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

# Raw Byte Constants
UTF8_BOM_BYTES = b"\xef\xbb\xbf"
CRLF_BYTES = b"\r\n"
NULL_BYTE = b"\x00"

# Binary Probe & Chunk Constants
BINARY_PROBE_CHUNK_SIZE = 8192
DEFAULT_MAX_WORKERS = 4

# Standard Installer Script File Names & Exclusions
INSTALLER_BASH_NAME = "install.sh"
INSTALLER_PWSH_NAME = "install.ps1"
INSTALLER_EXCLUDE_PARTS = ("node_modules", ".git", "dist", "build")

# Standard Git Command String Constants
GIT_EXECUTABLE = "git"
GIT_CMD_LS_FILES = "ls-files"
GIT_CMD_RM = "rm"
GIT_FLAG_FORCE = "-f"
GIT_FLAG_RECURSIVE = "-r"
GIT_FLAG_ERROR_UNMATCH = "--error-unmatch"

# Standard Status Badges
STATUS_GIT_TRACKED = " [Git Tracked]"
STATUS_UNTRACKED = " [Untracked]"

# Cache Dictionary Keys
CACHE_KEY_FILES = CacheKeyType.FILES.value
CACHE_KEY_TOTAL_FILES = CacheKeyType.TOTAL_FILES.value
CACHE_KEY_TIMESTAMP = CacheKeyType.TIMESTAMP.value
CACHE_KEY_VERSION = CacheKeyType.VERSION.value
CACHE_KEY_GENERATED_AT = CacheKeyType.GENERATED_AT.value
CACHE_KEY_EXPIRES_AT = CacheKeyType.EXPIRES_AT.value
CACHE_KEY_TTL_SECONDS = CacheKeyType.TTL_SECONDS.value
CACHE_KEY_SCAN_DURATION_MS = CacheKeyType.SCAN_DURATION_MS.value
CACHE_KEY_ROOT_PATH = CacheKeyType.ROOT_PATH.value
CACHE_KEY_MANIFESTS = CacheKeyType.MANIFESTS.value
CACHE_KEY_LANGUAGE_DISTRIBUTION = CacheKeyType.LANGUAGE_DISTRIBUTION.value
CACHE_KEY_LANGUAGE_ROOTS = CacheKeyType.LANGUAGE_ROOTS.value
CACHE_KEY_SUBSYSTEMS = CacheKeyType.SUBSYSTEMS.value
CACHE_KEY_ROOTS = CacheKeyType.ROOTS.value
CACHE_KEY_ENTRYPOINTS = CacheKeyType.ENTRYPOINTS.value
CACHE_KEY_SCHEMA_FILES = CacheKeyType.SCHEMA_FILES.value
CACHE_KEY_WORKFLOWS = CacheKeyType.WORKFLOWS.value
CACHE_KEY_SPEC_ROOTS = CacheKeyType.SPEC_ROOTS.value
CACHE_KEY_TEST_RUNNERS = CacheKeyType.TEST_RUNNERS.value

# Centralized Pycache & Temp Artifact Cleanup Presets
PYCACHE_DIR_NAMES: tuple[str, ...] = (
    "__pycache__", ".pytest_cache", ".coverage", ".mypy_cache", ".ruff_cache"
)
PYCACHE_FILE_EXTENSIONS: tuple[str, ...] = (
    ".pyc", ".pyo", ".pyd"
)
TEMP_ARTIFACT_EXTENSIONS: tuple[str, ...] = (
    ".tmp", ".log", ".swp", ".bak", ".orig", ".swo", ".swn"
)
TEMP_ARTIFACT_FILENAMES: tuple[str, ...] = (
    ".DS_Store", "Thumbs.db", "desktop.ini", ".directory"
)

# Centralized 18 CI Quality Gate Job Definitions
CI_JOBS_MATRIX: dict[str, list[str]] = {
    "Relative Path Check": [sys.executable, "linter-scripts/check-relative-paths.py"],
    "Prompts Loaded Check": [sys.executable, "linter-scripts/check-prompts-loaded.py"],
    "Readme Install Section Check": [sys.executable, "linter-scripts/check-readme-install-section.py"],
    "Forbidden Strings Check": [sys.executable, "linter-scripts/check-forbidden-strings.py"],
    "Newline Styling Check": [sys.executable, "linter-scripts/check-newline-styling.py"],
    "Fast File Scanner Cache": [sys.executable, "03-ai-scripts/11-fast-file-scanner.py", "--check"],
    "File Size Guard": [sys.executable, "03-ai-scripts/13-file-size-guard.py"],
    "Version Sync Check": [sys.executable, "03-ai-scripts/14-version-sync-checker.py"],
    "Bundle Installer Generation": ["node", "scripts/generate-bundle-installers.mjs"],
    "Spec Tree Sync": ["node", "scripts/sync-spec-tree.mjs"],
    "Codegen Determinism Check": [sys.executable, "linters-cicd/codegen/scripts/verify_codegen_determinism.py"],
    "Spec Verification Coverage": ["node", "scripts/spec-verification/generate-coverage-report.mjs", "--strict", "--out", "reports/spec-verification/coverage.md"],
    "Validate Version JSON": ["node", "scripts/validate-version-json.mjs"],
    "Doc Links Check": ["node", "scripts/docs/check-doc-links.mjs", "readme.md"],
    "Check File Sizes Baseline": [sys.executable, "linter-scripts/check-file-sizes.py", "--check"],
    "Newline Styling MJS Check": ["node", "linter-scripts/check-newline-styling.mjs"],
    "Spec Folder References Check": [sys.executable, "linter-scripts/check-spec-folder-refs.py"],
    "Sequence Integrity Check": [sys.executable, "linter-scripts/check-sequence-integrity.py"],
    "Prompt & Spec Path Integrity Check": [sys.executable, "linter-scripts/check-prompt-and-spec-paths.py"],
    "Linters CI/CD Test Suite": [sys.executable, "linters-cicd/tests/run.py"],
    "Go Base Test Suite": ["go", "test", "-C", "04-code/golang", "./..."],
}

# --- Module-Level Directory & File Constants ---
CACHE_BASE_DIR = Path("tmp/cache")
CACHE_PATHS_DIR = CACHE_BASE_DIR / "paths"
CACHE_LOCKS_DIR = CACHE_BASE_DIR / "locks"
CACHE_FILES_DIR = CACHE_BASE_DIR / "files"
LEGACY_CACHE_FILE = Path("tmp/repo-file-cache.json")
PRIMARY_CACHE_FILE = CACHE_BASE_DIR / "repo-file-cache.json"

PRIMARY_TOPOLOGY_CACHE_FILE = CACHE_PATHS_DIR / "codebase-topology-cache.json"
LEGACY_TOPOLOGY_CACHE_FILE = CACHE_BASE_DIR / "codebase-topology-cache.json"
DEFAULT_TTL_SECONDS = 1800  # 30 Minutes

DEFAULT_MAX_FILE_KB = 2048
LOCK_TIMEOUT_SECONDS = 5.0
STALE_LOCK_SECONDS = 15.0
MAX_READ_SIZE_BYTES = 20 * 1024 * 1024  # 20MB memory safety cap

EXCLUDE_DIRS = {
    ".git", ".gitmap", "gitmap", ".git-map",
    "node_modules", "dist", "build", ".venv", "venv",
    ".gemini", "tmp", ".system_generated", "vendor", ".cache",
    ".next", "bin", "obj", "coverage", "__pycache__",
    ".vs", ".idea", ".agent", "release-artifacts", "release-assets",
    ".turbo", ".parcel-cache",
}

BINARY_EXTENSIONS = {
    ".png", ".jpg", ".jpeg", ".gif", ".webp", ".ico", ".svg",
    ".pdf", ".zip", ".tar", ".gz", ".7z", ".rar", ".bz2",
    ".exe", ".dll", ".so", ".dylib", ".bin", ".o", ".a",
    ".db", ".sqlite", ".sqlite3",
    ".woff", ".woff2", ".ttf", ".eot", ".otf",
    ".mp3", ".mp4", ".wav", ".avi", ".mov",
    ".pyc", ".pyo", ".pyd", ".class",
}

DEFAULT_TEXT_EXTENSIONS = (
    ".md", ".markdown", ".py", ".ts", ".tsx", ".js", ".jsx",
    ".json", ".yaml", ".yml", ".go", ".php", ".cs", ".sh", ".ps1"
)

DEFAULT_CODE_EXTENSIONS = (
    ".ts", ".tsx", ".js", ".jsx", ".go", ".py", ".php", ".cs"
)

DEFAULT_CLI_EXTENSIONS = (
    ".go", ".ts", ".tsx", ".py", ".php"
)

ALLOWED_LARGE_FILES = {
    "src/data/specTree.json",
    "src\\data\\specTree.json",
    "slides-app/dist.zip",
    "slides-app\\dist.zip",
}

# Language Extension Mapping
LANG_EXT_MAP: dict[str, list[str]] = {
    "go": [".go"],
    "golang": [".go"],
    "ts": [".ts", ".tsx", ".mts", ".cts"],
    "typescript": [".ts", ".tsx", ".mts", ".cts"],
    "tsx": [".tsx"],
    "js": [".js", ".jsx", ".mjs", ".cjs"],
    "javascript": [".js", ".jsx", ".mjs", ".cjs"],
    "jsx": [".jsx"],
    "py": [".py", ".pyi"],
    "python": [".py", ".pyi"],
    "php": [".php", ".phtml"],
    "cs": [".cs"],
    "csharp": [".cs"],
    "rust": [".rs"],
    "rs": [".rs"],
    "md": [".md", ".markdown"],
    "markdown": [".md", ".markdown"],
    "json": [".json"],
    "yaml": [".yaml", ".yml"],
    "yml": [".yaml", ".yml"],
    "sh": [".sh", ".bash"],
    "bash": [".sh", ".bash"],
    "ps1": [".ps1", ".psm1", ".psd1"],
    "powershell": [".ps1", ".psm1", ".psd1"],
    "sql": [".sql"],
    "html": [".html", ".htm"],
    "css": [".css", ".scss", ".sass", ".less"],
    "c": [".c", ".h"],
    "cpp": [".cpp", ".hpp", ".cc", ".cxx"],
}

# Manifest File Signatures by Language
LANGUAGE_MANIFESTS: dict[LanguageType, tuple[str, ...]] = {
    LanguageType.GO: ("go.mod", "go.sum", "go.work"),
    LanguageType.RUST: ("Cargo.toml", "Cargo.lock"),
    LanguageType.PYTHON: ("pyproject.toml", "setup.py", "requirements.txt", "Pipfile", "poetry.lock", "uv.lock"),
    LanguageType.TYPESCRIPT: ("tsconfig.json", "tsconfig.base.json"),
    LanguageType.JAVASCRIPT: ("package.json", "jsconfig.json"),
    LanguageType.PHP: ("composer.json", "composer.lock", "artisan"),
    LanguageType.CSHARP: ("*.csproj", "*.sln", "NuGet.Config"),
    LanguageType.SQL: ("schema.sql", "migrations", "prisma/schema.prisma"),
}

# Subsystem Directory Indicators
SUBSYSTEM_DIR_HINTS: dict[SubsystemType, tuple[str, ...]] = {
    SubsystemType.DATABASE: ("db", "database", "migrations", "migration", "sql", "schemas", "schema", "prisma", "drizzle"),
    SubsystemType.BACKEND: ("cmd", "internal", "pkg", "api", "routes", "controllers", "handlers", "server", "services", "backend"),
    SubsystemType.FRONTEND: ("components", "views", "pages", "ui", "web", "frontend", "client", "app", "slides-app"),
    SubsystemType.CICD: (".github", "workflows", "scripts", "linter-scripts", "linters-cicd", "ci", ".lovable"),
    SubsystemType.DOCS: ("02-spec", "docs", "doc", "documentation", "prompts", ".lovable/prompts"),
    SubsystemType.TESTS: ("tests", "test", "02-spec", "__tests__", "testing", "fixtures"),
    SubsystemType.CLI: ("cli", "cmd", "commands", "bin"),
}

# Known Subsystem Entrypoint Files
SUBSYSTEM_ENTRYPOINTS: tuple[str, ...] = (
    "main.go", "main.py", "main.rs", "server.ts", "server.js",
    "app.py", "app.go", "index.ts", "index.js"
)

# Subsystem & Language Aliases for Fast CLI Routing
QUERY_ALIASES: dict[str, LanguageType | SubsystemType] = {
    "db": SubsystemType.DATABASE,
    "database": SubsystemType.DATABASE,
    "sql": SubsystemType.DATABASE,
    "migrations": SubsystemType.DATABASE,
    "schema": SubsystemType.DATABASE,
    "backend": SubsystemType.BACKEND,
    "server": SubsystemType.BACKEND,
    "api": SubsystemType.BACKEND,
    "frontend": SubsystemType.FRONTEND,
    "ui": SubsystemType.FRONTEND,
    "web": SubsystemType.FRONTEND,
    "client": SubsystemType.FRONTEND,
    "app": SubsystemType.FRONTEND,
    "ci": SubsystemType.CICD,
    "cicd": SubsystemType.CICD,
    "workflow": SubsystemType.CICD,
    "actions": SubsystemType.CICD,
    "docs": SubsystemType.DOCS,
    "doc": SubsystemType.DOCS,
    "02-spec": SubsystemType.DOCS,
    "prompt": SubsystemType.DOCS,
    "prompts": SubsystemType.DOCS,
    "cli": SubsystemType.CLI,
    "commands": SubsystemType.CLI,
    "tests": SubsystemType.TESTS,
    "test": SubsystemType.TESTS,
    "qa": SubsystemType.TESTS,
    "go": LanguageType.GO,
    "golang": LanguageType.GO,
    "rs": LanguageType.RUST,
    "rust": LanguageType.RUST,
    "py": LanguageType.PYTHON,
    "python": LanguageType.PYTHON,
    "ts": LanguageType.TYPESCRIPT,
    "typescript": LanguageType.TYPESCRIPT,
    "js": LanguageType.JAVASCRIPT,
    "javascript": LanguageType.JAVASCRIPT,
    "php": LanguageType.PHP,
    "cs": LanguageType.CSHARP,
    "csharp": LanguageType.CSHARP,
    "sh": LanguageType.SHELL,
    "shell": LanguageType.SHELL,
    "bash": LanguageType.SHELL,
    "ps1": LanguageType.SHELL,
    "powershell": LanguageType.SHELL,
    "md": LanguageType.MARKDOWN,
    "markdown": LanguageType.MARKDOWN,
}

# Centralized Raw Regex Definitions: Enum -> (Pattern String, Flags)
REGEX_DEFINITIONS: dict[RegexPatternType | str, tuple[str, int]] = {
    RegexPatternType.WINDOWS_BACKSLASH: (r"\\", 0),
    RegexPatternType.LEADING_DOT_SLASH: (r"^\./", 0),
    RegexPatternType.CRLF: (r"\r\n", 0),
    RegexPatternType.UNIVERSAL_LINE_ENDING: (r"\r\n|\r", 0),
    RegexPatternType.TRAILING_WHITESPACE: (r"[ \t]+$", re.MULTILINE),
    RegexPatternType.SEQ_PREFIX: (r"^([0-9]+)-(.*)$", 0),
    RegexPatternType.UPPERCASE: (r"[A-Z]", 0),
    RegexPatternType.FILE_URI_WIN: (r"file:///[A-Za-z]:/[^\s\)\]\"'>]+", 0),
    RegexPatternType.DRIVE_ABS_WIN: (r"(?<![A-Za-z0-9_])[A-Za-z]:\\[A-Za-z0-9_\\.-]+", 0),
    RegexPatternType.REPO_FILE_URI: (r"file:///[A-Za-z]:/[^/]+/coding-guidelines/([^\s\)\]\"'>]+)", 0),
    RegexPatternType.EXPLICIT_DOUBLE_TRUE: (r"==\s*true\b", re.IGNORECASE),
    RegexPatternType.EXPLICIT_TRIPLE_TRUE: (r"===\s*true\b", re.IGNORECASE),
    RegexPatternType.EXPLICIT_PYTHON_TRUE: (r"==\s*True\b", 0),
    RegexPatternType.COMMENT_PREFIX: (r"^\s*(//|#|\*|/\*)", 0),
    RegexPatternType.COBRA_COMMAND: (r"var\s+(\w+Cmd)\s*=\s*&cobra\.Command\s*\{([^}]+)\}", re.DOTALL),
    RegexPatternType.SHORT_DESC: (r"Short:\s*\"[^\"]+\"", 0),
    RegexPatternType.EXAMPLE_USAGE: (r"Example:\s*\"[^\"]+\"", 0),
    RegexPatternType.CHANGELOG_HEADER: (r"##\s+\[v?([0-9]+\.[0-9]+\.[0-9]+[^\]]*)\]", 0),
    RegexPatternType.FILE_NUM_PREFIX: (r"^([0-9]+)-(.*)\.md$", 0),
    RegexPatternType.H1_HEADER: (r"^(#\s+)([0-9]+)(\s*[-—:]\s*)(.*)$", re.MULTILINE),
    RegexPatternType.PLACEHOLDER_TOKEN: (r"[A-Z0-9_]*PLACEHOLDER[A-Z0-9_]*", 0),
    RegexPatternType.NON_ALPHANUMERIC: (r"[^a-zA-Z0-9_-]+", 0),
}

# --- Thread-Safe Lazy Regex Registry (Zero import-time compilation, on-the-fly dynamic auto-registration) ---
class RegexRegistry:
    """
    Thread-safe lazy-compiling regex registry.
    Initializes all known entries to None mapping on startup.
    If an unknown pattern or key is requested, it logs the event, registers it on-the-fly in both maps, and compiles.
    """
    _compiled_patterns: dict[RegexPatternType | str, re.Pattern | None] = {pt: None for pt in RegexPatternType}
    _lock = threading.Lock()

    @classmethod
    def get(
        cls,
        pattern_type: RegexPatternType | str,
        default_pattern: str | None = None,
        flags: int = 0
    ) -> re.Pattern:
        """
        Lazily compiles on first demand and returns cached immutable re.Pattern.
        Auto-registers dynamic regex strings on-the-fly if missing.
        """
        cached = cls._compiled_patterns.get(pattern_type)
        if cached is not None:
            return cached

        with cls._lock:
            # Check again within lock
            cached = cls._compiled_patterns.get(pattern_type)
            if cached is not None:
                return cached

            # Auto-register on the fly if pattern is not present
            if pattern_type not in REGEX_DEFINITIONS:
                raw_pattern = default_pattern or str(pattern_type)
                print(f"ℹ️ [RegexRegistry] Dynamic on-the-fly registration: '{pattern_type}' -> '{raw_pattern}'")
                REGEX_DEFINITIONS[pattern_type] = (raw_pattern, flags)

            raw_pattern, compile_flags = REGEX_DEFINITIONS[pattern_type]
            compiled = re.compile(raw_pattern, compile_flags)
            cls._compiled_patterns[pattern_type] = compiled
            return compiled

    @classmethod
    def get_group(cls, *pattern_types: RegexPatternType | str) -> tuple[re.Pattern, ...]:
        """Lazily retrieves a tuple of compiled re.Pattern objects."""
        return tuple(cls.get(pt) for pt in pattern_types)

def get_compiled_regex(
    pattern_type: RegexPatternType | str,
    default_pattern: str | None = None,
    flags: int = 0
) -> re.Pattern:
    """Convenience functional accessor for RegexRegistry.get."""
    return RegexRegistry.get(pattern_type, default_pattern=default_pattern, flags=flags)

def get_compiled_regex_group(*pattern_types: RegexPatternType | str) -> tuple[re.Pattern, ...]:
    """Convenience functional accessor for RegexRegistry.get_group."""
    return RegexRegistry.get_group(*pattern_types)

# --- Shared String Formatting & Collection Utilities ---

def format_comma_separated(items: Iterable[Any]) -> str:
    """Combines an iterable of items into a clean comma-separated string."""
    return COMMA_SPACE_SEPARATOR.join(str(item) for item in items)

def format_keys(mapping: Any, separator: str = COMMA_SPACE_SEPARATOR) -> str:
    """Extracts and formats keys from a dictionary or iterable with separator."""
    if hasattr(mapping, "keys"):
        return separator.join(str(k) for k in mapping.keys())
    return separator.join(str(k) for k in mapping)

# --- Path & File Utility Functions ---

def is_ignored_directory(dir_name: str, custom_excludes: set[str] | None = None) -> bool:
    """Checks if directory name is in the global or custom exclusion list."""
    excludes = EXCLUDE_DIRS if custom_excludes is None else EXCLUDE_DIRS | custom_excludes
    return dir_name.lower() in {d.lower() for d in excludes}

def is_ignored_path(path: str | Path, custom_excludes: set[str] | None = None) -> bool:
    """Checks if any segment of the path matches an excluded directory."""
    excludes = EXCLUDE_DIRS if custom_excludes is None else EXCLUDE_DIRS | custom_excludes
    excludes_lower = {d.lower() for d in excludes}
    parts = Path(path).parts
    return any(p.lower() in excludes_lower for p in parts)

def is_binary_file(file_path: Path) -> bool:
    """
    Checks if file is binary by extension or memory-safe 8KB chunk probing for null bytes.
    Avoids loading full large files into RAM.
    """
    is_binary_ext = file_path.suffix.lower() in BINARY_EXTENSIONS
    if is_binary_ext:
        return True

    try:
        if not file_path.is_file():
            return False

        with open(file_path, "rb") as f:
            chunk = f.read(BINARY_PROBE_CHUNK_SIZE)
            return NULL_BYTE in chunk
    except Exception:
        return False

def is_allowed_large_file(file_path: str | Path) -> bool:
    """Checks if file is on the explicit waiver list for large generated assets."""
    norm = normalize_rel_path(file_path).lstrip(f"{CURRENT_DIR}{PATH_SEPARATOR}")
    return norm in {normalize_rel_path(f).lstrip(f"{CURRENT_DIR}{PATH_SEPARATOR}") for f in ALLOWED_LARGE_FILES}

def normalize_rel_path(path: str | Path) -> str:
    """Converts a path into a canonical relative POSIX path using PATH_SEPARATOR."""
    re_slash = get_compiled_regex(RegexPatternType.WINDOWS_BACKSLASH)
    re_lead = get_compiled_regex(RegexPatternType.LEADING_DOT_SLASH)
    p_str = re_slash.sub(PATH_SEPARATOR, str(path))
    return re_lead.sub(EMPTY_STRING, p_str)

def normalize_extensions(extensions: tuple | set | list | str | None) -> set[str] | None:
    """Normalizes custom extensions into a lowercased set with leading dots."""
    if not extensions:
        return None

    if isinstance(extensions, str):
        raw_items = [e.strip() for e in extensions.split(",") if e.strip()]
    else:
        raw_items = [str(e).strip() for e in extensions if str(e).strip()]

    normalized = set()
    for item in raw_items:
        clean = item.lower()
        if not clean.startswith(DOT_CHAR):
            clean = f"{DOT_CHAR}{clean}"
        normalized.add(clean)

    return normalized if normalized else None

def read_file_safe(
    path: str | Path,
    max_bytes: int = MAX_READ_SIZE_BYTES,
    encoding: str = DEFAULT_ENCODING
) -> str | None:
    """
    Memory-safe and fault-tolerant file reader.
    Handles missing or deleted files gracefully with zero crashes.
    Normalizes both CRLF (\\r\\n) and legacy Mac CR (\\r) to strict UNIX LF (\\n).
    """
    p = Path(path)
    if not p.is_file():
        return None

    try:
        re_univ_nl = get_compiled_regex(RegexPatternType.UNIVERSAL_LINE_ENDING)
        with open(p, "r", encoding=encoding, errors="replace") as f:
            raw_text = f.read(max_bytes)
            return re_univ_nl.sub(LINE_SEPARATOR, raw_text)
    except (FileNotFoundError, PermissionError, OSError):
        return None

def read_file_lf(path: str | Path, encoding: str = DEFAULT_ENCODING) -> str:
    """Reads a text file ensuring strict UNIX LF. Returns empty string if file does not exist."""
    content = read_file_safe(path, encoding=encoding)
    return content if content is not None else EMPTY_STRING

def write_file_lf(
    path: str | Path,
    content: str,
    encoding: str = DEFAULT_ENCODING
) -> bool:
    """
    Atomic write ensuring strict UNIX LF line endings and Unix permission preservation.
    Encodes file according to the specified encoding constant (default: utf-8 without BOM).
    Preserves original executable bits (chmod +x / st_mode) on Linux/macOS.
    """
    p = Path(path)
    p.parent.mkdir(parents=True, exist_ok=True)
    temp_path = p.with_name(f"{p.name}.tmp_{os.getpid()}_{int(time.time()*1000)}")

    original_mode = None
    if p.exists():
        try:
            original_mode = p.stat().st_mode
        except Exception:
            pass

    try:
        re_univ_nl = get_compiled_regex(RegexPatternType.UNIVERSAL_LINE_ENDING)
        lf_content = re_univ_nl.sub(LINE_SEPARATOR, content)
        with open(temp_path, "wb") as f:
            f.write(lf_content.encode(encoding))

        if original_mode is not None:
            try:
                os.chmod(temp_path, original_mode)
            except Exception:
                pass

        try:
            temp_path.replace(p)
        except PermissionError:
            with open(p, "wb") as f:
                f.write(lf_content.encode(encoding))
            if temp_path.exists():
                try:
                    temp_path.unlink()
                except Exception:
                    pass

        return True
    except Exception:
        if temp_path.exists():
            try:
                temp_path.unlink()
            except Exception:
                pass
        return False

# --- Dual-Platform Cross-Process Locking Mechanism ---

@contextmanager
def atomic_cache_lock(lock_name: str = "repo-cache.lock", timeout: float = LOCK_TIMEOUT_SECONDS):
    """
    Dual-platform cross-process lock:
    - On Unix: Uses kernel-level `fcntl.flock` (automatic cleanup on crash/SIGKILL).
    - On Windows: Uses `os.O_CREAT | os.O_EXCL` with PID timestamp & stale lock eviction (>15s).
    """
    CACHE_LOCKS_DIR.mkdir(parents=True, exist_ok=True)
    lock_file = CACHE_LOCKS_DIR / lock_name
    start_time = time.time()
    is_acquired = False
    lock_fd = None

    if IS_FCNTL_AVAILABLE:
        try:
            lock_fd = open(lock_file, "w")
            while time.time() - start_time < timeout:
                try:
                    fcntl.flock(lock_fd, fcntl.LOCK_EX | fcntl.LOCK_NB)
                    is_acquired = True
                    break
                except (BlockingIOError, OSError):
                    time.sleep(0.02)
        except Exception:
            is_acquired = False

        try:
            yield is_acquired
        finally:
            if lock_fd is not None:
                try:
                    fcntl.flock(lock_fd, fcntl.LOCK_UN)
                    lock_fd.close()
                except Exception:
                    pass
    else:
        while time.time() - start_time < timeout:
            try:
                if lock_file.exists():
                    lock_age = time.time() - lock_file.stat().st_mtime
                    if lock_age > STALE_LOCK_SECONDS:
                        try:
                            lock_file.unlink()
                        except Exception:
                            pass

                fd = os.open(str(lock_file), os.O_CREAT | os.O_EXCL | os.O_RDWR)
                os.write(fd, f"pid={os.getpid()}{LINE_SEPARATOR}time={time.time()}".encode(DEFAULT_ENCODING))
                os.close(fd)
                is_acquired = True
                break
            except FileExistsError:
                time.sleep(0.02)
            except Exception:
                break

        try:
            yield is_acquired
        finally:
            if is_acquired:
                if lock_file.exists():
                    try:
                        lock_file.unlink()
                    except Exception:
                        pass

# --- Pluggable Cache Management ---

def load_repo_cache() -> dict[str, Any]:
    """Loads pre-computed repository file cache from tmp/cache/ or legacy tmp/."""
    for target in (PRIMARY_CACHE_FILE, LEGACY_CACHE_FILE):
        if target.exists():
            try:
                with open(target, "r", encoding=DEFAULT_ENCODING) as f:
                    data = json.load(f)
                    if isinstance(data, dict):
                        if CACHE_KEY_FILES in data:
                            return data
            except Exception:
                pass
    return {}

def save_repo_cache(cache_data: dict[str, Any]) -> None:
    """Saves repository cache safely with atomic locking and dual-path sync."""
    CACHE_PATHS_DIR.mkdir(parents=True, exist_ok=True)
    with atomic_cache_lock("repo-cache-write.lock"):
        temp_primary = PRIMARY_CACHE_FILE.with_suffix(".json.tmp")
        try:
            with open(temp_primary, "w", encoding=DEFAULT_ENCODING) as f:
                json.dump(cache_data, f, indent=2)
            temp_primary.replace(PRIMARY_CACHE_FILE)
        except Exception:
            pass

        try:
            temp_legacy = LEGACY_CACHE_FILE.with_suffix(".json.tmp")
            with open(temp_legacy, "w", encoding=DEFAULT_ENCODING) as f:
                json.dump(cache_data, f, indent=2)
            temp_legacy.replace(LEGACY_CACHE_FILE)
        except Exception:
            pass

# --- Two-Phase Streaming Engine with Inode Cycle Protection ---

def stream_cached_files(
    cache_data: dict[str, Any],
    root_dir: str = CURRENT_DIR,
    extensions: set[str] | tuple | None = None,
    custom_excludes: set[str] | None = None
) -> Generator[Path, None, None]:
    """Phase 1: Streams valid files from cache first. Automatically skips missing/deleted/excluded files."""
    file_list = cache_data.get(CACHE_KEY_FILES, [])
    norm_root = normalize_rel_path(root_dir).rstrip(PATH_SEPARATOR)
    ext_set = normalize_extensions(extensions)

    for rel_path in file_list:
        norm_p = normalize_rel_path(rel_path)
        if norm_root:
            if norm_root != CURRENT_DIR:
                if not norm_p.startswith(norm_root + PATH_SEPARATOR):
                    if norm_p != norm_root:
                        continue

        if is_ignored_path(norm_p, custom_excludes=custom_excludes):
            continue

        p = Path(norm_p)
        if not p.exists():
            continue

        if is_binary_file(p):
            continue

        if ext_set:
            if p.suffix.lower() not in ext_set:
                continue

        yield p

def stream_directory_files(
    root_dir: str = CURRENT_DIR,
    extensions: set[str] | tuple | None = None,
    custom_excludes: set[str] | None = None
) -> Generator[Path, None, None]:
    """
    Phase 2: Walks filesystem pruning ignored folders (including nested .git, .gitmap, node_modules).
    Guarded with visited inode tracking (st_dev, st_ino) on Unix to prevent symlink recursion cycles.
    """
    ext_set = normalize_extensions(extensions)
    visited_inodes: set[tuple[int, int]] = set()

    for root, dirs, files in os.walk(root_dir, followlinks=False):
        try:
            st = os.stat(root)
            inode_key = (st.st_dev, st.st_ino)
            if inode_key in visited_inodes:
                dirs[:] = []
                continue
            visited_inodes.add(inode_key)
        except Exception:
            pass

        dirs[:] = [d for d in dirs if not is_ignored_directory(d, custom_excludes=custom_excludes)]
        for f in files:
            p = Path(os.path.join(root, f))
            if is_binary_file(p):
                continue

            if ext_set:
                if p.suffix.lower() not in ext_set:
                    continue

            yield p

def process_repository_files(
    processor_fn: Callable[[Path], Any],
    root_dir: str = CURRENT_DIR,
    extensions: set[str] | tuple | list | str | None = None,
    is_use_cache: bool = True,
    custom_excludes: set[str] | None = None
) -> dict[str, Any]:
    """
    Two-Phase Universal Pipeline:
    1. Starts immediately with cached files if available (<0.1ms).
    2. Streams and discovers new / modified files on disk.
    3. Gracefully skips missing or removed files during processing.
    4. Executes processor_fn on each unique file and aggregates statistics.
    """
    start_time = time.perf_counter()
    cache_data = load_repo_cache() if is_use_cache else {}
    processed_paths: set[str] = set()
    results = []
    norm_exts = normalize_extensions(extensions)

    if cache_data:
        if CACHE_KEY_FILES in cache_data:
            for p in stream_cached_files(cache_data, root_dir=root_dir, extensions=norm_exts, custom_excludes=custom_excludes):
                norm_p = normalize_rel_path(p)
                if norm_p not in processed_paths:
                    processed_paths.add(norm_p)
                    res = processor_fn(p)
                    if res is not None:
                        results.append(res)

    for p in stream_directory_files(root_dir=root_dir, extensions=norm_exts, custom_excludes=custom_excludes):
        norm_p = normalize_rel_path(p)
        if norm_p not in processed_paths:
            processed_paths.add(norm_p)
            res = processor_fn(p)
            if res is not None:
                results.append(res)

    elapsed_ms = (time.perf_counter() - start_time) * 1000

    return {
        "total_files": len(processed_paths),
        "results": results,
        "elapsed_ms": elapsed_ms,
    }
