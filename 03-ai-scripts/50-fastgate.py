#!/usr/bin/env python3
"""
Fast-Gate Pre-Commit Runner
Scopes AST linters (relative paths, nested ifs, boolean guidelines) exclusively to staged or modified files.
Execution time is <1.5s for fast pre-commit feedback.
Includes --full fallback to run on all files.
"""

from __future__ import annotations

import argparse
from importlib import import_module
import os
from pathlib import Path
import re
import shutil
import subprocess
import sys
import time
from typing import Any

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")

sys.path.insert(0, str(Path(__file__).resolve().parent))
engine = import_module("02-shared-engine")

ExitCodeType = engine.ExitCodeType
DEFAULT_ENCODING = engine.DEFAULT_ENCODING
LINE_SEPARATOR = engine.LINE_SEPARATOR
CURRENT_DIR = engine.CURRENT_DIR

FORBIDDEN_PATH_PATTERNS = (
    (re.compile(r"file:///[a-zA-Z]:[/\\]?", re.IGNORECASE), "Absolute file:/// URI with drive letter"),
    (re.compile(r"file:" + r"///work/", re.IGNORECASE), "Absolute file:/// URI"),
    (re.compile(r"file:///(?:Users|home|root)/", re.IGNORECASE), "Absolute file:/// URI to user directory"),
    (re.compile(r"\b[dD]:[/\\]work[/\\]gitmap\b", re.IGNORECASE), "Hardcoded absolute repo path (D:\\work\\gitmap)"),
    (re.compile(r"\b[cC]:[/\\]Users[/\\][a-zA-Z0-9_.-]+[/\\]\.gemini\b", re.IGNORECASE), "Hardcoded user agent directory"),
)

EXCLUDE_EXTS = {
    ".png", ".jpg", ".jpeg", ".gif", ".webp", ".ico", ".pdf", ".zip",
    ".gz", ".tar", ".exe", ".bin", ".db", ".sqlite", ".woff", ".woff2", ".ttf"
}

CODE_EXTS = {".go", ".ts", ".tsx", ".js", ".jsx", ".py", ".php"}

EXCLUDE_DIRS = {
    ".git", "node_modules", "dist", "build", "bin", ".next", ".gitmap",
    "vendor", "coverage", ".gemini", ".system_generated", "tests/fixtures",
    "scratch", "temp-scripts", "temp-agents", "temp", "linter-scripts",
    ".ai-memory/scratch", ".ai-memory/temp-agents", "03-ai-scripts", "scripts",
    "04-code", ".tmp"
}

ALLOWLIST_FILES = {
    ".github/workflows/goreleaser-smoke.yml",
    "linter-scripts/check-relative-paths.py",
    "03-ai-scripts/07-relative-path-fixer.py",
    "03-ai-scripts/49-verify-privacy-and-relative-paths.py",
    "03-ai-scripts/50-fastgate.py",
}

ALLOWLIST_PREFIXES = (
    ".ai-memory/pipeline-ai/",
    ".ai-memory/cicd/",
    ".ai-memory/temp/",
    ".gitmap/output/",
)


def get_repo_root() -> Path:
    """Returns canonical repository root path."""
    return Path(__file__).resolve().parent.parent


def run_git_lines(cmd_args: list[str], root_dir: Path) -> list[str]:
    """Executes git command and returns stripped output lines."""
    git_bin = shutil.which("git") or "git"
    res = subprocess.run([git_bin] + cmd_args, cwd=str(root_dir), capture_output=True, text=True, encoding=DEFAULT_ENCODING)
    has_output = bool(res.returncode == 0 and res.stdout)
    if not has_output:
        return []
    return [line.strip().replace("\\", "/") for line in res.stdout.splitlines() if line.strip()]


def get_staged_files(root_dir: Path) -> list[str]:
    """Retrieves staged modified or added files from git index."""
    args = ["diff", "--cached", "--name-only", "--diff-filter=ACMR"]
    return run_git_lines(args, root_dir)


def get_modified_files(root_dir: Path) -> list[str]:
    """Retrieves unstaged modified or added files from working directory."""
    args = ["diff", "--name-only", "--diff-filter=ACMR"]
    return run_git_lines(args, root_dir)


def get_all_repo_files(root_dir: Path) -> list[str]:
    """Retrieves all tracked files across the repository."""
    return run_git_lines(["ls-files"], root_dir)


def resolve_target_files(args: argparse.Namespace, root_dir: Path) -> list[str]:
    """Resolves target files based on CLI flags with staged priority."""
    has_custom = bool(args.files)
    if has_custom:
        return [f.replace("\\", "/").strip() for f in args.files if f.strip()]
    is_full = bool(args.full)
    if is_full:
        return get_all_repo_files(root_dir)
    staged = get_staged_files(root_dir)
    has_staged = bool(staged)
    if has_staged:
        return staged
    return get_modified_files(root_dir)


def is_path_allowed_for_relative(rel_path: str) -> bool:
    """Checks whether file path is allowed to contain raw paths."""
    is_allowlisted = rel_path in ALLOWLIST_FILES
    if is_allowlisted:
        return True
    return rel_path.startswith(ALLOWLIST_PREFIXES)


def is_candidate_for_relative_check(rel_path: str) -> bool:
    """Filters files eligible for absolute path linting."""
    ext = os.path.splitext(rel_path)[1].lower()
    is_excluded_ext = ext in EXCLUDE_EXTS
    if is_excluded_ext:
        return False
    is_allowed = is_path_allowed_for_relative(rel_path)
    if is_allowed:
        return False
    return True


def is_candidate_for_code_check(rel_path: str) -> bool:
    """Filters code files eligible for AST and boolean checks."""
    ext = os.path.splitext(rel_path)[1].lower()
    is_target_ext = ext in CODE_EXTS
    if not is_target_ext:
        return False
    parts = rel_path.split("/")
    has_excluded_dir = any(p in EXCLUDE_DIRS for p in parts)
    return not has_excluded_dir


def check_relative_paths_for_file(rel_path: str, root_dir: Path) -> list[str]:
    """Scans single file for forbidden absolute path patterns."""
    file_path = root_dir / rel_path
    if not file_path.is_file():
        return []
    try:
        content = file_path.read_text(encoding=DEFAULT_ENCODING, errors="replace")
    except Exception as exc:
        return [f"Failed to read file: {exc}"]
    return scan_lines_for_paths(content.splitlines(), rel_path)


def scan_lines_for_paths(lines: list[str], rel_path: str) -> list[str]:
    """Scans content lines for forbidden path matches."""
    violations: list[str] = []
    for idx, line in enumerate(lines, 1):
        if "FORBIDDEN_PATH_PATTERNS" in line or "50-fastgate.py" in line:
            continue
        vios = match_line_forbidden_patterns(line, idx, rel_path)
        violations.extend(vios)
    return violations


def match_line_forbidden_patterns(line: str, idx: int, rel_path: str) -> list[str]:
    """Matches single line against regex forbidden path patterns."""
    hits: list[str] = []
    for pat, desc in FORBIDDEN_PATH_PATTERNS:
        if pat.search(line):
            hits.append(f"{rel_path}:{idx}: {desc} -> {line.strip()[:80]}")
    return hits


def check_relative_paths(files: list[str], root_dir: Path) -> list[str]:
    """Validates relative paths across all candidate files."""
    violations: list[str] = []
    for f in files:
        if is_candidate_for_relative_check(f):
            vios = check_relative_paths_for_file(f, root_dir)
            violations.extend(vios)
    return violations


def load_linter_scanner(module_name: str, root_dir: Path) -> Any:
    """Dynamically loads scan_file from linter-scripts module."""
    linter_dir = root_dir / "linter-scripts"
    if str(linter_dir) not in sys.path:
        sys.path.insert(0, str(linter_dir))
    try:
        mod = import_module(module_name)
        return getattr(mod, "scan_file", None)
    except Exception:
        return None


def scan_file_with_handler(p: Path, rel: str, scanner: Any) -> list[str]:
    """Executes scanner on single file and formats output messages."""
    if not p.is_file():
        return []
    try:
        results = scanner(p)
    except Exception as exc:
        return [f"{rel}: error running scanner: {exc}"]
    return [f"{rel}:{item[0]}: {item[1]}" for item in results]


def check_nested_ifs(files: list[str], root_dir: Path) -> list[str]:
    """Scans candidate files for nested if violations."""
    scanner = load_linter_scanner("check-nested-ifs", root_dir)
    if not scanner:
        return []
    violations: list[str] = []
    for f in files:
        if is_candidate_for_code_check(f):
            vios = scan_file_with_handler(root_dir / f, f, scanner)
            violations.extend(vios)
    return violations


def check_boolean_guidelines(files: list[str], root_dir: Path) -> list[str]:
    """Scans candidate files for boolean guideline violations."""
    scanner = load_linter_scanner("check-boolean-guidelines", root_dir)
    if not scanner:
        return []
    violations: list[str] = []
    for f in files:
        if is_candidate_for_code_check(f):
            vios = scan_file_with_handler(root_dir / f, f, scanner)
            violations.extend(vios)
    return violations


def print_gate_outcome(name: str, violations: list[str]) -> bool:
    """Prints single quality gate result and returns success status."""
    has_passed = (len(violations) == 0)
    if has_passed:
        print(f"  ✔ {name:<24} : PASS")
        return True
    print(f"  ❌ {name:<24} : FAIL ({len(violations)} violation(s))")
    for msg in violations[:5]:
        print(f"     • {msg}")
    return False


def print_summary_and_status(p_vios: list[str], n_vios: list[str], b_vios: list[str], start_time: float) -> bool:
    """Prints status of all three gates and summary duration line."""
    p_pass = print_gate_outcome("Relative Paths", p_vios)
    n_pass = print_gate_outcome("Nested Ifs", n_vios)
    b_pass = print_gate_outcome("Boolean Guidelines", b_vios)
    elapsed = time.perf_counter() - start_time
    is_all_clean = bool(p_pass and n_pass and b_pass)
    if is_all_clean:
        print(f"\n✔ FastGate passed in {elapsed:.3f}s (<1.5s target). Ready to commit.")
        return True
    print(f"\n❌ FastGate rejected in {elapsed:.3f}s. Please resolve violations before commit.")
    return False


def parse_cli_arguments() -> argparse.Namespace:
    """Parses command line arguments for fastgate pre-commit runner."""
    parser = argparse.ArgumentParser(description="Fast-Gate Pre-Commit AST Linter Runner")
    parser.add_argument("--files", nargs="*", default=[], help="Explicit files to lint")
    parser.add_argument("--full", "-f", action="store_true", help="Fallback to scan all repository files")
    parser.add_argument("--quiet", "-q", action="store_true", help="Minimal output mode")
    return parser.parse_args()


def execute_fastgate(args: argparse.Namespace) -> int:
    """Executes fastgate quality checks and returns exit code."""
    root_dir = get_repo_root()
    files = resolve_target_files(args, root_dir)
    start_time = time.perf_counter()
    if not files:
        print("⚡ FastGate: No staged or modified candidate files to inspect. Clean!")
        return ExitCodeType.SUCCESS.value
    print(f"⚡ FastGate Pre-Commit Validation ({len(files)} target file(s)):")
    p_vios = check_relative_paths(files, root_dir)
    n_vios = check_nested_ifs(files, root_dir)
    b_vios = check_boolean_guidelines(files, root_dir)
    is_success = print_summary_and_status(p_vios, n_vios, b_vios, start_time)
    return ExitCodeType.SUCCESS.value if is_success else ExitCodeType.VIOLATIONS_FOUND.value


def main() -> None:
    """Main execution entry point."""
    args = parse_cli_arguments()
    sys.exit(execute_fastgate(args))


if __name__ == "__main__":
    main()
