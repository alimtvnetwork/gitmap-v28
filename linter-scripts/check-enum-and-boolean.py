#!/usr/bin/env python3
"""
check-enum-and-boolean.py - Linter for Boolean Conventions, Naming, Enums, and Conditionals.

Checks:
1. No explicit boolean comparisons (== true, === true, == false, === false).
2. No inverted success checks (!isSuccess).
3. No single-line if blocks (anti-compression rule).
4. All enum definitions must end in 'Type'.
5. No nested if statements (nesting depth > 1 inside an if/else body in Go source code).
"""

import argparse
from concurrent.futures import ThreadPoolExecutor, as_completed
from importlib import import_module
import json
import os
from pathlib import Path
import re
import subprocess
import sys
import threading
import time
from typing import Any

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8", errors="replace")

sys.path.insert(0, str(Path(__file__).resolve().parent.parent / "03-ai-scripts"))
engine = import_module("02-shared-engine")
chunk_items = engine.chunk_items
WorkerHeartbeatMonitor = engine.WorkerHeartbeatMonitor

ROOT_DIR = Path(__file__).resolve().parent.parent
TARGET_EXTS = {'.go', '.ts', '.tsx', '.php', '.py'}
EXCLUDE_DIRS = {
    '.git', 'node_modules', 'dist', 'build', 'bin', '.next', '.gitmap',
    'vendor', 'coverage', '.gemini', '.system_generated', 'tests/fixtures',
    'scratch', 'temp-scripts', 'temp-agents', 'temp', 'linter-scripts',
    '.lovable/scratch', '.lovable/temp-agents', '03-ai-scripts', 'scripts',
    '04-code', '.tmp'
}

# Regex patterns
EXPLICIT_BOOL_REGEX = re.compile(r'\b(==\s*true|===\s*true|==\s*false|===\s*false)\b')
INVERTED_SUCCESS_REGEX = re.compile(r'!\s*[a-zA-Z0-9_$.->]*\bisSuccess\b')
SINGLE_LINE_IF_REGEX = re.compile(r'^\s*if\b.*\{[^{}]+\}\s*$')
ENUM_MISSING_TYPE_REGEX = re.compile(r'^\s*(?:export\s+)?enum\s+(?!\w+Type\b)\w+\b')

def is_comment_or_doc(line: str, ext: str) -> bool:
    s = line.strip()
    if not s:
        return True
    if ext in ('.go', '.ts', '.tsx', '.php') and (s.startswith('//') or s.startswith('/*') or s.startswith('*')):
        return True
    if ext == '.py' and (s.startswith('#') or s.startswith('"""') or s.startswith("'''")):
        return True
    return False

def strip_go_code(content: str) -> list[tuple[int, str]]:
    """
    Strips comments, string literals, and rune literals from Go code,
    preserving exact line numbers and code structure.
    """
    out_lines = []
    current_line = []
    line_num = 1
    i = 0
    n = len(content)

    while i < n:
        c = content[i]
        if c == '\n':
            out_lines.append((line_num, ''.join(current_line)))
            current_line = []
            line_num += 1
            i += 1
            continue
        if c == '/' and i + 1 < n and content[i + 1] == '/':
            i += 2
            while i < n and content[i] != '\n':
                i += 1
            continue
        if c == '/' and i + 1 < n and content[i + 1] == '*':
            i += 2
            while i < n:
                if content[i] == '\n':
                    out_lines.append((line_num, ''.join(current_line)))
                    current_line = []
                    line_num += 1
                    i += 1
                elif content[i] == '*' and i + 1 < n and content[i + 1] == '/':
                    i += 2
                    break
                else:
                    i += 1
            continue
        if c == '`':
            i += 1
            while i < n and content[i] != '`':
                if content[i] == '\n':
                    out_lines.append((line_num, ''.join(current_line)))
                    current_line = []
                    line_num += 1
                i += 1
            if i < n:
                i += 1
            continue
        if c == '"':
            i += 1
            while i < n and content[i] != '"':
                if content[i] == '\\':
                    i += 2
                elif content[i] == '\n':
                    break
                else:
                    i += 1
            if i < n and content[i] == '"':
                i += 1
            continue
        if c == "'":
            i += 1
            while i < n and content[i] != "'":
                if content[i] == '\\':
                    i += 2
                elif content[i] == '\n':
                    break
                else:
                    i += 1
            if i < n and content[i] == "'":
                i += 1
            continue

        current_line.append(c)
        i += 1

    if current_line or content.endswith('\n'):
        out_lines.append((line_num, ''.join(current_line)))

    return out_lines

def check_nested_ifs_in_go(content: str, filepath: Path) -> list[tuple[int, str]]:
    if "gitmap" in filepath.parts and "constants" in filepath.parts:
        return []

    stripped = strip_go_code(content)
    violations = []
    if_stack = []
    brace_depth = 0

    for line_num, code in stripped:
        trimmed = code.strip()
        if not trimmed:
            continue

        is_else_if = 'else if' in trimmed
        is_if = bool(re.search(r'\bif\b', trimmed)) and '{' in trimmed

        closes_at_start = len(re.match(r'^\}+', trimmed).group(0)) if re.match(r'^\}+', trimmed) else 0
        current_depth = brace_depth - closes_at_start

        while if_stack and current_depth <= if_stack[-1]:
            if_stack.pop()

        if is_if:
            if not is_else_if:
                if if_stack:
                    violations.append((line_num, f"Nested 'if' detected (depth {len(if_stack) + 1}): '{trimmed}'"))
                if_stack.append(brace_depth)

        opens = code.count('{')
        closes = code.count('}')
        brace_depth += (opens - closes)

        while if_stack and brace_depth <= if_stack[-1]:
            if_stack.pop()

    return violations

def check_file(filepath: Path) -> list[str]:
    ext = filepath.suffix
    if ext not in TARGET_EXTS:
        return []
    if filepath.name.endswith('_test.go') or '.test.' in filepath.name or '.spec.' in filepath.name:
        return []

    try:
        content = filepath.read_text(encoding='utf-8', errors='replace')
    except Exception:
        return []

    lines = content.splitlines()
    violations = []
    is_constants_file = ("gitmap" in filepath.parts and "constants" in filepath.parts)

    for idx, line in enumerate(lines, start=1):
        if is_comment_or_doc(line, ext):
            continue

        # 1. Explicit boolean comparison
        if EXPLICIT_BOOL_REGEX.search(line):
            violations.append(f"{filepath}:{idx}: Explicit boolean comparison: '{line.strip()}'")

        # 2. Inverted success check
        if INVERTED_SUCCESS_REGEX.search(line):
            violations.append(f"{filepath}:{idx}: Inverted success check (!isSuccess): '{line.strip()}'")

        # 3. Single-line if statement
        if not is_constants_file and SINGLE_LINE_IF_REGEX.search(line) and not line.strip().startswith('//'):
            violations.append(f"{filepath}:{idx}: Single-line 'if' block (must expand to multiple lines): '{line.strip()}'")

        # 4. Enums without Type suffix
        if ext in ('.ts', '.tsx', '.php') and ENUM_MISSING_TYPE_REGEX.search(line):
            violations.append(f"{filepath}:{idx}: Enum missing 'Type' suffix: '{line.strip()}'")

    # 5. Nested ifs in Go files
    if ext == '.go':
        nested = check_nested_ifs_in_go(content, filepath)
        for line_num, msg in nested:
            violations.append(f"{filepath}:{line_num}: {msg}")

    return violations

def collect_target_files() -> list[Path]:
    """Pre-gathers all eligible target files across repository root."""
    target_files: list[Path] = []
    for root, dirs, files in os.walk(ROOT_DIR):
        dirs[:] = [d for d in dirs if d not in EXCLUDE_DIRS and not any(ex in d for ex in EXCLUDE_DIRS)]
        for file in files:
            p = Path(root) / file
            if p.suffix in TARGET_EXTS and not p.name.endswith('_test.go'):
                target_files.append(p)
    return target_files


def print_scan_progress(completed: int, total: int, workers: int, start_time: float) -> None:
    """Emits live scan percentage and throughput."""
    pct = (completed / total * 100.0) if total > 0 else 100.0
    elapsed = max(0.001, time.time() - start_time)
    fps = completed / elapsed
    msg = f"\rChecking boolean & enum compliance: [ {completed:4d}/{total:4d} ] {pct:5.1f}% | {workers} workers | {fps:5.1f} files/sec"
    sys.stdout.write(msg)
    sys.stdout.flush()


def parse_cli_args() -> argparse.Namespace:
    """Parses command line arguments for boolean and enum linter."""
    parser = argparse.ArgumentParser(description="Check boolean and enum conventions.")
    parser.add_argument("--changed-only", "-c", action="store_true", help="Check only changed files")
    parser.add_argument("--commits", "-n", type=int, default=20, help="Commit window (default: 20)")
    parser.add_argument("--chunk-size", type=int, default=8, help="Files per chunk (default: 8)")
    parser.add_argument("--workers", "-w", type=int, default=10, help="Concurrency (default: 10)")

    return parser.parse_args()


def load_changed_targets(commits: int) -> list[Path]:
    """Loads target files from git-changed-files.json."""
    manifest = ROOT_DIR / ".lovable/temp/git-changed-files.json"
    is_fresh = manifest.is_file() and (time.time() - manifest.stat().st_mtime < 120.0)
    if not is_fresh:
        extractor = ROOT_DIR / "03-ai-scripts/27-git-changed-files.py"
        subprocess.run([sys.executable, str(extractor), "--commits", str(commits), "--quiet"], cwd=str(ROOT_DIR), check=True)
    data = json.loads(manifest.read_text(encoding="utf-8"))
    targets: list[Path] = []
    for item in data.get("files", []):
        path_str = item["path"] if isinstance(item, dict) else str(item)
        p = ROOT_DIR / path_str
        if p.is_file() and p.suffix.lower() in TARGET_EXTS and not p.name.endswith('_test.go'):
            rel = p.relative_to(ROOT_DIR).as_posix()
            if not any(ex in rel.split("/") for ex in EXCLUDE_DIRS):
                targets.append(p)

    return targets


def resolve_candidate_files(args: argparse.Namespace) -> list[Path]:
    """Resolves eligible files according to changed-only flag."""
    if args.changed_only:
        return load_changed_targets(args.commits)

    return collect_target_files()


def check_file_chunk(chunk: list[Path]) -> list[str]:
    """Checks a chunk of files for compliance."""
    chunk_violations: list[str] = []
    for p in chunk:
        v = check_file(p)
        if v:
            chunk_violations.extend(v)

    return chunk_violations


def execute_chunked_checks(
    chunks: list[list[Path]],
    total_files: int,
    cpu_cores: int,
    monitor: Any,
) -> list[str]:
    """Executes parallel chunked check using ThreadPoolExecutor."""
    all_violations: list[str] = []
    completed_count = 0
    start_time = time.time()
    with ThreadPoolExecutor(max_workers=cpu_cores) as pool:
        futures = {pool.submit(check_file_chunk, chunk): len(chunk) for chunk in chunks}
        for fut in as_completed(futures):
            chunk_len = futures[fut]
            chunk_violations = fut.result()
            all_violations.extend(chunk_violations)
            completed_count += chunk_len
            monitor.increment_processed(chunk_len)
            print_scan_progress(completed_count, total_files, cpu_cores, start_time)

    return all_violations


def report_violations_and_exit(all_violations: list[str], total_files: int, elapsed: float) -> int:
    """Formats violations report and returns exit code."""
    sys.stdout.write("\n")
    print(f"Scanned {total_files} source files for boolean, enum, and conditional compliance in {elapsed:.2f}s.\n")
    if all_violations:
        print(f"❌ FAILED: Found {len(all_violations)} violation(s):")
        for v in all_violations[:50]:
            print(f"  - {v}")
        return 1
    print("✅ PASS: All boolean, naming, enum, and conditional checks passed (zero explicit booleans, inverted success, or nested ifs).")

    return 0


def main() -> int:
    """Main execution entrypoint."""
    args = parse_cli_args()
    mode_label = f" (changed only, last {args.commits} commits)" if args.changed_only else ""
    print(f"=== Running Boolean & Enum Linter (check-enum-and-boolean.py){mode_label} ===")
    target_files = resolve_candidate_files(args)
    total_files = len(target_files)
    chunks = chunk_items(target_files, args.chunk_size)
    cpu_cores = min(args.workers, max(1, len(chunks)))
    print(f"▸ Discovered {total_files} source file(s). Checking compliance with {cpu_cores} worker threads...")
    monitor = WorkerHeartbeatMonitor(total_files, "files", 5.0, cpu_cores)
    monitor.start()
    start_time = time.time()
    all_violations = execute_chunked_checks(chunks, total_files, cpu_cores, monitor)
    monitor.stop()
    elapsed = time.time() - start_time
    exit_code = report_violations_and_exit(all_violations, total_files, elapsed)

    return exit_code


if __name__ == '__main__':
    sys.exit(main())
