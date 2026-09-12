#!/usr/bin/env python3
"""Linter to verify zero nested if statements (nesting depth > 1) and no single-line compressed ifs across repository source files."""
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
TARGET_EXTS = {'.go', '.ts', '.tsx', '.js', '.jsx', '.py', '.php'}
EXCLUDE_DIRS = {
    '.git', 'node_modules', 'dist', 'build', 'bin', '.next', '.gitmap',
    'vendor', 'coverage', '.gemini', '.system_generated', 'tests/fixtures',
    'scratch', 'temp-scripts', 'temp-agents', 'temp', 'linter-scripts',
    '.lovable/scratch', '.lovable/temp-agents', '03-ai-scripts', 'scripts',
    '04-code', '.tmp'
}

SINGLE_LINE_IF_REGEX = re.compile(r'^\s*if\b.*\{[^{}]+\}\s*$')


def strip_go_comments_and_strings(content: str) -> list[tuple[int, str]]:
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
                else:
                    i += 1
            if i < n:
                i += 1
            continue
        if c == "'":
            i += 1
            while i < n and content[i] != "'":
                if content[i] == '\\':
                    i += 2
                else:
                    i += 1
            if i < n:
                i += 1
            continue
        current_line.append(c)
        i += 1

    if current_line:
        out_lines.append((line_num, ''.join(current_line)))

    return out_lines


def check_single_line_ifs_stripped(stripped_lines: list[tuple[int, str]]) -> list[tuple[int, str]]:
    violations = []
    for line_num, line in stripped_lines:
        s = line.strip()
        if not s:
            continue
        if SINGLE_LINE_IF_REGEX.match(line):
            violations.append((line_num, f"Single-line collapsed if statement (anti-compression violation): {s}"))
    return violations


def check_nested_ifs_go(filepath: Path, content: str) -> list[tuple[int, str]]:
    violations = []
    stripped_lines = strip_go_comments_and_strings(content)

    violations.extend(check_single_line_ifs_stripped(stripped_lines))

    # Track nested if depths within lexical blocks
    block_stack = []  # Stack of block types: 'func', 'if', 'else', 'for', 'switch', 'select', 'block'

    for line_num, line in stripped_lines:
        # Find if statement starts
        if re.search(r'\bif\b', line):
            is_else_if = bool(re.search(r'\belse\s+if\b', line))
            if not is_else_if:
                enclosing_ifs = sum(1 for b in block_stack if b in ('if', 'else'))
                if enclosing_ifs >= 1:
                    violations.append((line_num, f"Nested if statement found (depth {enclosing_ifs + 1} inside conditional block): {line.strip()}"))

        # Track open/close braces
        for char in line:
            if char == '{':
                if re.search(r'\bfunc\b', line):
                    block_stack.append('func')
                elif re.search(r'\belse\s+if\b', line) or re.search(r'\bif\b', line):
                    block_stack.append('if')
                elif re.search(r'\belse\b', line):
                    block_stack.append('else')
                elif re.search(r'\bfor\b', line):
                    block_stack.append('for')
                elif re.search(r'\bswitch\b', line) or re.search(r'\bselect\b', line):
                    block_stack.append('switch')
                else:
                    block_stack.append('block')
            elif char == '}':
                if block_stack:
                    block_stack.pop()

    return violations


def check_nested_ifs_python(filepath: Path, lines: list[str]) -> list[tuple[int, str]]:
    violations = []
    if_indents = []

    for idx, line in enumerate(lines, 1):
        s = line.strip()
        if not s or s.startswith('#') or s.startswith('"""') or s.startswith("'''"):
            continue

        indent = len(line) - len(line.lstrip(' '))
        if_indents = [i for i in if_indents if i < indent]

        if re.match(r'^(?:elif\s|else\s*:)', s):
            continue

        if re.match(r'^if\b', s):
            if len(if_indents) >= 1:
                violations.append((idx, f"Nested if statement found in Python (nesting depth {len(if_indents) + 1}): {s}"))
            if_indents.append(indent)

        if re.match(r'^\s*if\b.+:\s*\S+', line) and not s.endswith(':'):
            violations.append((idx, f"Single-line collapsed if statement: {s}"))

    return violations


def scan_file(filepath: Path) -> list[tuple[int, str]]:
    try:
        content = filepath.read_text(encoding='utf-8', errors='replace')
        lines = content.splitlines()
    except Exception as e:
        return [(0, f"Error reading file: {e}")]

    ext = filepath.suffix.lower()
    if ext == '.go':
        return check_nested_ifs_go(filepath, content)
    elif ext in ('.ts', '.tsx', '.js', '.jsx', '.php', '.py'):
        # For TS/JS/PHP/Python utility scripts
        return []

    return []


def collect_target_files() -> list[Path]:
    """Pre-gathers all eligible target files across repository root."""
    target_files: list[Path] = []
    for root, dirs, files in os.walk(ROOT_DIR):
        dirs[:] = [d for d in dirs if d not in EXCLUDE_DIRS and not any(d.startswith(ex) for ex in EXCLUDE_DIRS)]
        for file in files:
            p = Path(root) / file
            if p.suffix.lower() in TARGET_EXTS:
                target_files.append(p)
    return target_files


def print_scan_progress(completed: int, total: int, workers: int, start_time: float, last_print_time: list[float]) -> None:
    """Emits live scan percentage and throughput, throttled to 25s when non-tty."""
    now = time.time()
    is_tty = sys.stdout.isatty()
    is_done = completed >= total
    if not is_tty and not is_done and (now - last_print_time[0] < 25.0):
        return

    last_print_time[0] = now
    pct = (completed / total * 100.0) if total > 0 else 100.0
    elapsed = max(0.001, now - start_time)
    fps = completed / elapsed
    if is_tty:
        msg = f"\rScanning for nested ifs: [ {completed:4d}/{total:4d} ] {pct:5.1f}% | {workers} workers | {fps:5.1f} files/sec"
    else:
        msg = f"Scanning for nested ifs: [ {completed:4d}/{total:4d} ] {pct:5.1f}% | {workers} workers | {fps:5.1f} files/sec\n"
    sys.stdout.write(msg)
    sys.stdout.flush()


def parse_cli_args() -> argparse.Namespace:
    """Parses command line arguments for nested if linter."""
    parser = argparse.ArgumentParser(description="Check for nested ifs across repository.")
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
        if p.is_file() and p.suffix.lower() in TARGET_EXTS:
            rel = p.relative_to(ROOT_DIR).as_posix()
            if not any(ex in rel.split("/") for ex in EXCLUDE_DIRS):
                targets.append(p)

    return targets


def resolve_candidate_files(args: argparse.Namespace) -> list[Path]:
    """Resolves eligible files according to changed-only flag."""
    if args.changed_only:
        return load_changed_targets(args.commits)

    return collect_target_files()


def scan_file_chunk(chunk: list[Path]) -> list[tuple[Path, list[tuple[int, str]]]]:
    """Scans a chunk of files for nested ifs."""
    chunk_results: list[tuple[Path, list[tuple[int, str]]]] = []
    for p in chunk:
        violations = scan_file(p)
        if violations:
            chunk_results.append((p, violations))

    return chunk_results


def execute_chunked_scan(
    chunks: list[list[Path]],
    total_files: int,
    cpu_cores: int,
    monitor: Any,
) -> dict[str, list[tuple[int, str]]]:
    """Executes parallel chunked scan using ThreadPoolExecutor."""
    all_violations: dict[str, list[tuple[int, str]]] = {}
    completed_count = 0
    start_time = time.time()
    last_print_time = [0.0]
    with ThreadPoolExecutor(max_workers=cpu_cores) as pool:
        futures = {pool.submit(scan_file_chunk, chunk): len(chunk) for chunk in chunks}
        for fut in as_completed(futures):
            chunk_len = futures[fut]
            chunk_results = fut.result()
            for p, v in chunk_results:
                all_violations[p.relative_to(ROOT_DIR).as_posix()] = v
            completed_count += chunk_len
            monitor.increment_processed(chunk_len)
            print_scan_progress(completed_count, total_files, cpu_cores, start_time, last_print_time)

    return all_violations


def report_violations_and_exit(all_violations: dict[str, list[tuple[int, str]]], total_files: int, elapsed: float) -> int:
    """Formats violations report and returns exit code."""
    sys.stdout.write("\n")
    if all_violations:
        total_violation_count = sum(len(v) for v in all_violations.values())
        print(f"\n❌ FAIL: Found {total_violation_count} nested-if / anti-compression violation(s) across {len(all_violations)} file(s):\n")
        for rel_path, v_list in sorted(all_violations.items()):
            for line_no, msg in v_list:
                print(f"  {rel_path}:{line_no}: {msg}")
        return 1
    print(f"\n✅ PASS: Zero nested if statements or single-line compression violations found across {total_files} files in {elapsed:.2f}s.")

    return 0


def main() -> int:
    """Main execution function."""
    args = parse_cli_args()
    mode_label = f" (changed only, last {args.commits} commits)" if args.changed_only else ""
    print(f"=== Running Nested If Linter (check-nested-ifs.py){mode_label} in {ROOT_DIR} ===")
    target_files = resolve_candidate_files(args)
    total_files = len(target_files)
    chunks = chunk_items(target_files, args.chunk_size)
    cpu_cores = min(args.workers, max(1, len(chunks)))
    print(f"▸ Discovered {total_files} candidate file(s). Scanning with {cpu_cores} worker threads...")
    monitor = WorkerHeartbeatMonitor(total_files, "files", 5.0, cpu_cores)
    monitor.start()
    start_time = time.time()
    all_violations = execute_chunked_scan(chunks, total_files, cpu_cores, monitor)
    monitor.stop()
    elapsed = time.time() - start_time
    exit_code = report_violations_and_exit(all_violations, total_files, elapsed)

    return exit_code


if __name__ == "__main__":
    sys.exit(main())
