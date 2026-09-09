#!/usr/bin/env python3
"""Linter to verify that no absolute filesystem paths or file:/// URIs exist in repository files."""
from __future__ import annotations

import argparse
from concurrent.futures import ThreadPoolExecutor, as_completed
from importlib import import_module
import json
import os
from pathlib import Path
import re
import subprocess
import sys
import time
from typing import Any

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")

sys.path.insert(0, str(Path(__file__).resolve().parent.parent / "03-ai-scripts"))
engine = import_module("02-shared-engine")

chunk_items = engine.chunk_items
WorkerHeartbeatMonitor = engine.WorkerHeartbeatMonitor

FORBIDDEN_PATTERNS = [
    (re.compile(r"file:///[a-zA-Z]:[/\\]?", re.IGNORECASE), "Absolute file:/// URI with drive letter"),
    (re.compile(r"file:///work/", re.IGNORECASE), "Absolute file:/// URI"),
    (re.compile(r"file:///(?:Users|home|root)/", re.IGNORECASE), "Absolute file:/// URI to user directory"),
    (re.compile(r"\b[dD]:[/\\]work[/\\]gitmap\b", re.IGNORECASE), "Hardcoded absolute repo path (D:\\work\\gitmap)"),
    (re.compile(r"\b[cC]:[/\\]Users[/\\][a-zA-Z0-9_.-]+[/\\]\.gemini\b", re.IGNORECASE), "Hardcoded user agent directory"),
]

EXCLUDE_EXTS = {
    ".png", ".jpg", ".jpeg", ".gif", ".webp", ".ico", ".pdf", ".zip",
    ".gz", ".tar", ".exe", ".bin", ".db", ".sqlite", ".woff", ".woff2", ".ttf"
}

ALLOWLIST_FILES = {
    ".github/workflows/goreleaser-smoke.yml",
    "linter-scripts/check-relative-paths.py",
    ".lovable/ai-fix-scripts/07-relative-path-fixer.py",
}


def parse_arguments() -> argparse.Namespace:
    """Parses CLI flags for parallel relative paths checker."""
    parser = argparse.ArgumentParser(description="Check for absolute filesystem paths in repo.")
    parser.add_argument("--changed-only", "-c", action="store_true", help="Check only files changed in recent commits")
    parser.add_argument("--commits", "-n", type=int, default=20, help="Commit window for changed files (default: 20)")
    parser.add_argument("--chunk-size", type=int, default=8, help="Number of files per worker chunk (default: 8)")
    parser.add_argument("--workers", "-w", type=int, default=10, help="Worker concurrency (default: 10)")
    parser.add_argument("--quiet", "-q", action="store_true", help="Quiet output (no banner)")

    return parser.parse_args()


def ensure_changed_files_manifest(repo_root: Path, commits: int) -> Path:
    """Ensures git-changed-files.json is generated in .lovable/temp/."""
    manifest = repo_root / ".lovable/temp/git-changed-files.json"
    is_fresh = manifest.is_file() and (time.time() - manifest.stat().st_mtime < 120.0)
    if not is_fresh:
        extractor = repo_root / "03-ai-scripts/27-git-changed-files.py"
        subprocess.run([sys.executable, str(extractor), "--commits", str(commits), "--quiet"], cwd=str(repo_root), check=True)

    return manifest


def load_changed_files(repo_root: Path, commits: int) -> list[str]:
    """Loads changed files list from git-changed-files.json."""
    manifest = ensure_changed_files_manifest(repo_root, commits)
    data = json.loads(manifest.read_text(encoding="utf-8"))
    files_list = data.get("files", [])
    raw_files = [item["path"] if isinstance(item, dict) else str(item) for item in files_list]

    return raw_files


def collect_all_tracked_files(repo_root: Path) -> list[str]:
    """Retrieves all tracked files using git ls-files."""
    res = subprocess.run(["git", "ls-files"], cwd=str(repo_root), capture_output=True, text=True, encoding="utf-8")
    if res.returncode != 0:
        return []
    lines = [f.strip() for f in res.stdout.splitlines() if f.strip()]

    return lines


def deduplicate_candidate_files(raw_files: list[str], repo_root: Path) -> list[str]:
    """Deduplicates file paths using dictionary lookup and filters excluded extensions."""
    seen: dict[str, str] = {}
    for f in raw_files:
        norm = f.replace("\\", "/").strip()
        ext = os.path.splitext(norm)[1].lower()
        if ext in EXCLUDE_EXTS or norm in ALLOWLIST_FILES:
            continue
        if (repo_root / norm).is_file():
            seen[norm] = norm
    deduped = sorted(seen.keys())

    return deduped


def scan_single_line(line: str, line_idx: int, rel_path: str) -> list[tuple[str, int, str]]:
    """Checks a single line against all forbidden patterns."""
    if "FORBIDDEN_PATTERNS" in line or "check-relative-paths.py" in line:
        return []
    violations = []
    for pattern, desc in FORBIDDEN_PATTERNS:
        if pattern.search(line):
            violations.append((rel_path, line_idx, f"{desc}: {line.strip()[:100]}"))

    return violations


def scan_file_lines(lines: list[str], rel_path: str) -> list[tuple[str, int, str]]:
    """Iterates through file lines and aggregates pattern violations."""
    violations: list[tuple[str, int, str]] = []
    for idx, line in enumerate(lines, 1):
        line_violations = scan_single_line(line, idx, rel_path)
        violations.extend(line_violations)

    return violations


def scan_file_path(rel_path: str, repo_root: Path) -> list[tuple[str, int, str]]:
    """Reads file content and returns detected path violations."""
    full_path = repo_root / rel_path
    try:
        content = full_path.read_text(encoding="utf-8", errors="replace")
        lines = content.splitlines()
    except Exception as exc:
        return [(rel_path, 0, f"Failed to read file: {exc}")]
    violations = scan_file_lines(lines, rel_path)

    return violations


def process_file_chunk(
    worker_id: str,
    chunk_idx: int,
    chunk: list[str],
    repo_root: Path,
    monitor: Any,
    log_picks: bool,
) -> list[tuple[str, int, str]]:
    """Processes a chunk of 8 files within worker thread."""
    if log_picks:
        print(f"[{worker_id}] Picked chunk {chunk_idx + 1} ({len(chunk)} files)...")
    chunk_violations: list[tuple[str, int, str]] = []
    for rel_path in chunk:
        monitor.update_worker(worker_id, f"Scanning {Path(rel_path).name}")
        violations = scan_file_path(rel_path, repo_root)
        chunk_violations.extend(violations)
        monitor.increment_processed(1)
    monitor.update_worker(worker_id, "Idle")

    return chunk_violations


def execute_parallel_scan(
    chunks: list[list[str]],
    workers: int,
    repo_root: Path,
    monitor: Any,
    log_picks: bool,
) -> list[tuple[str, int, str]]:
    """Executes parallel chunked scanning across ThreadPoolExecutor workers."""
    all_violations: list[tuple[str, int, str]] = []
    with ThreadPoolExecutor(max_workers=workers) as pool:
        futures = {
            pool.submit(process_file_chunk, f"Worker-{idx % workers + 1}", idx, chunk, repo_root, monitor, log_picks): chunk
            for idx, chunk in enumerate(chunks)
        }
        for fut in as_completed(futures):
            v_list = fut.result()
            all_violations.extend(v_list)

    return all_violations


def report_violations_and_exit(
    violations: list[tuple[str, int, str]],
    total_files: int,
    elapsed_sec: float,
) -> int:
    """Prints scan outcome report and returns exit code."""
    if violations:
        print(f"\n❌ FAIL: Found {len(violations)} absolute path / URI violation(s):\n")
        for file_path, line_no, msg in violations:
            print(f"  {file_path}:{line_no}: {msg}")
        return 1
    print(f"\n✅ PASS: No absolute filesystem paths or file:/// URIs found across {total_files} files ({elapsed_sec:.2f}s).")

    return 0


def resolve_candidate_files(args: argparse.Namespace, repo_root: Path) -> list[str]:
    """Resolves and deduplicates candidate files based on CLI flags."""
    raw = load_changed_files(repo_root, args.commits) if args.changed_only else collect_all_tracked_files(repo_root)
    files = deduplicate_candidate_files(raw, repo_root)

    return files


def print_mode_header(is_changed_only: bool, commits: int) -> None:
    """Prints scan header mode."""
    mode_desc = f"Changed files only (last {commits} commits)" if is_changed_only else "All repository files"
    print(f"=== Checking for Absolute Paths and file:/// URIs [{mode_desc}] ===")


def main() -> int:
    """Main execution function."""
    args = parse_arguments()
    print_mode_header(args.changed_only, args.commits)
    repo_root = Path(__file__).resolve().parent.parent
    files = resolve_candidate_files(args, repo_root)
    chunks = chunk_items(files, args.chunk_size)
    workers = min(args.workers, max(1, len(chunks)))
    monitor = WorkerHeartbeatMonitor(len(files), "files", 5.0, workers)
    monitor.start()
    start_time = time.perf_counter()
    log_picks = bool(len(chunks) <= 20 or args.changed_only)
    violations = execute_parallel_scan(chunks, workers, repo_root, monitor, log_picks)
    monitor.stop()
    elapsed = time.perf_counter() - start_time
    exit_code = report_violations_and_exit(violations, len(files), elapsed)

    return exit_code


if __name__ == "__main__":
    sys.exit(main())
