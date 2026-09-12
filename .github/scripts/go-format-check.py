#!/usr/bin/env python3
"""Cross-platform check and auto-fix for gofmt formatting with multi-core parallelism and live progress."""
from __future__ import annotations

import argparse
from concurrent.futures import ThreadPoolExecutor, as_completed
import os
import shutil
import subprocess
import sys
import threading
import time

CHUNK_SIZE = 40
IGNORED_DIRS = {
    ".git", ".gitmap", "node_modules", "dist", "build", "bin",
    ".next", "vendor", "coverage", ".gemini", ".system_generated",
    "scratch", "temp-scripts", "temp-agents", "temp", ".lovable", ".tmp"
}


def configure_io_encoding() -> None:
    """Configures UTF-8 encoding for standard input/output streams."""
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    if hasattr(sys.stderr, "reconfigure"):
        sys.stderr.reconfigure(encoding="utf-8")


def create_arg_parser(repo_root: str) -> argparse.ArgumentParser:
    """Builds and returns command line argument parser."""
    default_dir = os.path.join(repo_root, "cli")
    parser = argparse.ArgumentParser(
        description="Verify, auto-format, and auto-commit Go source files with gofmt."
    )
    parser.add_argument(
        "target_dir",
        nargs="?",
        default=default_dir,
        help="Target directory to check (default: <repo_root>/gitmap)",
    )
    parser.add_argument(
        "--check-only",
        action="store_true",
        help="Perform dry-run check only; do not auto-format or commit.",
    )
    parser.add_argument(
        "--no-commit",
        action="store_true",
        help="Format unformatted files, but skip creating git commit.",
    )
    parser.add_argument(
        "--push",
        action="store_true",
        help="Push commit to remote repository.",
    )
    parser.add_argument(
        "--no-push",
        action="store_true",
        help="Do not push commit to remote even in CI.",
    )

    return parser


def discover_go_files(target_dir: str) -> list[str]:
    """Discovers all .go files within target directory, skipping build artifacts and caches."""
    go_files = []
    for root, dirs, files in os.walk(target_dir):
        dirs[:] = [d for d in dirs if d not in IGNORED_DIRS and not d.startswith(".")]
        for f in files:
            if f.endswith(".go"):
                go_files.append(os.path.join(root, f))
    return sorted(go_files)


def chunk_list(items: list[str], size: int) -> list[list[str]]:
    """Splits a list of items into chunks of at most `size` items."""
    return [items[i:i + size] for i in range(0, len(items), size)]


def print_check_progress(completed: int, total: int, workers: int, start_time: float, prefix: str = "Checking") -> None:
    """Emits live formatted progress percentage and throughput."""
    pct = (completed / total * 100.0) if total > 0 else 100.0
    elapsed = max(0.001, time.time() - start_time)
    fps = completed / elapsed
    msg = f"\r{prefix} Go format: [ {completed:4d}/{total:4d} ] {pct:5.1f}% | {workers} workers | {fps:5.1f} files/sec"
    sys.stdout.write(msg)
    sys.stdout.flush()


def check_chunk(gofmt_exe: str, target_dir: str, chunk: list[str]) -> list[str]:
    """Runs `gofmt -l` on a chunk of Go files and returns unformatted file paths."""
    rel_files = [os.path.relpath(p, target_dir) for p in chunk]
    cmd = [gofmt_exe, "-l"] + rel_files
    res = subprocess.run(
        cmd,
        cwd=target_dir,
        capture_output=True,
        text=True,
        encoding="utf-8",
        errors="replace",
    )
    unformatted = [line.strip() for line in (res.stdout or "").splitlines() if line.strip()]
    return unformatted


def run_gofmt_check(target_dir: str) -> tuple[int, list[str]]:
    """Runs parallel `gofmt -l` across all CPU cores to find unformatted Go files."""
    gofmt_exe = shutil.which("gofmt")
    if not gofmt_exe:
        sys.stderr.write("ERROR: gofmt not on PATH\n")
        return 2, []

    files = discover_go_files(target_dir)
    total_files = len(files)
    if total_files == 0:
        return 0, []

    cpu_cores = os.cpu_count() or 16
    chunks = chunk_list(files, CHUNK_SIZE)
    unformatted_all: list[str] = []
    completed_count = 0
    lock = threading.Lock()
    start_time = time.time()

    print(f"▸ Discovered {total_files} Go file(s) in {target_dir}. Checking with {cpu_cores} worker threads...")
    print_check_progress(0, total_files, cpu_cores, start_time, prefix="Checking")

    with ThreadPoolExecutor(max_workers=cpu_cores) as pool:
        futures = {pool.submit(check_chunk, gofmt_exe, target_dir, c): len(c) for c in chunks}
        for fut in as_completed(futures):
            chunk_len = futures[fut]
            unformatted_chunk = fut.result()
            with lock:
                if unformatted_chunk:
                    unformatted_all.extend(unformatted_chunk)
                completed_count += chunk_len
                print_check_progress(completed_count, total_files, cpu_cores, start_time, prefix="Checking")

    sys.stdout.write("\n")
    elapsed = time.time() - start_time
    print(f"▸ Completed check of {total_files} Go files in {elapsed:.2f}s ({len(unformatted_all)} unformatted).")
    return 0, unformatted_all


def fix_chunk(gofmt_exe: str, target_dir: str, chunk: list[str]) -> bool:
    """Applies `gofmt -w` to a chunk of files."""
    cmd = [gofmt_exe, "-w"] + chunk
    res = subprocess.run(
        cmd,
        cwd=target_dir,
        capture_output=True,
        text=True,
        encoding="utf-8",
        errors="replace",
    )
    return res.returncode == 0


def apply_gofmt_fix(target_dir: str, unformatted: list[str]) -> bool:
    """Applies `gofmt -w` to format the specified files concurrently."""
    gofmt_exe = shutil.which("gofmt") or "gofmt"
    cpu_cores = os.cpu_count() or 16
    chunks = chunk_list(unformatted, CHUNK_SIZE)
    total_files = len(unformatted)
    completed_count = 0
    lock = threading.Lock()
    has_error = False
    start_time = time.time()

    print_check_progress(0, total_files, cpu_cores, start_time, prefix="Formatting")
    with ThreadPoolExecutor(max_workers=cpu_cores) as pool:
        futures = {pool.submit(fix_chunk, gofmt_exe, target_dir, c): len(c) for c in chunks}
        for fut in as_completed(futures):
            chunk_len = futures[fut]
            is_ok = fut.result()
            with lock:
                if not is_ok:
                    has_error = True
                completed_count += chunk_len
                print_check_progress(completed_count, total_files, cpu_cores, start_time, prefix="Formatting")

    sys.stdout.write("\n")
    return not has_error


def verify_git_worktree(repo_root: str) -> bool:
    """Checks if repository root is inside a valid git working tree."""
    res = subprocess.run(
        ["git", "rev-parse", "--is-inside-work-tree"],
        cwd=repo_root,
        capture_output=True,
        text=True,
    )
    return res.returncode == 0 and res.stdout.strip() == "true"


def ensure_git_identity(repo_root: str) -> None:
    """Ensures a valid git author name and email are configured."""
    res = subprocess.run(
        ["git", "config", "user.name"],
        cwd=repo_root,
        capture_output=True,
        text=True,
    )
    if not res.stdout.strip():
        subprocess.run(
            ["git", "config", "user.name", "github-actions[bot]"],
            cwd=repo_root,
            check=False,
        )
        subprocess.run(
            ["git", "config", "user.email", "41898282+github-actions[bot]@users.noreply.github.com"],
            cwd=repo_root,
            check=False,
        )


def query_modified_go_files(repo_root: str) -> list[str]:
    """Retrieves list of modified .go files via git status."""
    res = subprocess.run(
        ["git", "status", "--porcelain=v1", "--", "*.go"],
        cwd=repo_root,
        capture_output=True,
        text=True,
        encoding="utf-8",
        errors="replace",
    )
    files: list[str] = []
    for line in res.stdout.splitlines():
        clean_line = line.strip()
        if len(clean_line) > 3:
            files.append(clean_line[3:].strip().strip('"'))

    return files


def create_git_commit(repo_root: str, files: list[str]) -> bool:
    """Stages reformatted Go files and commits with [skip ci] tag."""
    ensure_git_identity(repo_root)
    stage_cmd = ["git", "add", "-A", "--"] + files
    subprocess.run(stage_cmd, cwd=repo_root, check=False)
    msg = "style(go): auto-format Go files with gofmt [skip ci]"
    res = subprocess.run(["git", "commit", "-m", msg], cwd=repo_root, capture_output=True, text=True)
    return res.returncode == 0


def determine_git_branch() -> str:
    """Resolves current target git branch from CI env or HEAD."""
    branch = os.environ.get("GITHUB_HEAD_REF") or os.environ.get("GITHUB_REF_NAME") or ""
    if branch.startswith("refs/heads/"):
        return branch[len("refs/heads/"):]

    return branch


def push_git_commit(repo_root: str, branch: str) -> bool:
    """Pushes the newly created formatting commit to the remote branch."""
    ref_spec = f"HEAD:{branch}" if branch else "HEAD"
    cmd = ["git", "push", "origin", ref_spec]
    res = subprocess.run(cmd, cwd=repo_root, capture_output=True, text=True)
    if res.returncode != 0:
        sys.stderr.write(f"::warning::Failed to push formatting commit: {res.stderr.strip()}\n")
        return False
    print(f"🚀 Successfully pushed auto-format commit to {ref_spec}.")
    return True


def handle_auto_commit_and_push(repo_root: str, args: argparse.Namespace) -> None:
    """Coordinates git commit and push when unformatted files are repaired."""
    if args.no_commit or not verify_git_worktree(repo_root):
        print("ℹ️ Auto-commit skipped (--no-commit or non-git environment).")
        return
    modified = query_modified_go_files(repo_root)
    if not modified:
        print("ℹ️ No modified Go files detected in git working tree.")
        return
    if not create_git_commit(repo_root, modified):
        print("::warning::Git commit did not produce a new commit.")
        return
    print("💾 Auto-committed formatting fixes: style(go): auto-format Go files with gofmt [skip ci]")
    is_ci = os.environ.get("GITHUB_ACTIONS") == "true"
    should_push = (is_ci and not args.no_push) or args.push
    if should_push:
        target_branch = determine_git_branch()
        push_git_commit(repo_root, target_branch)


def print_dry_run_findings(files: list[str]) -> None:
    """Prints unformatted file list detected during dry run."""
    print("::notice::Dry run detected unformatted .go file(s):", file=sys.stderr)
    for f in files:
        print(f"  {f}", file=sys.stderr)


def execute_check_and_fix(target_dir: str, repo_root: str, args: argparse.Namespace) -> int:
    """Executes the dry-run check, conditional auto-fix, and commit pipeline."""
    code, unformatted = run_gofmt_check(target_dir)
    if code != 0:
        return code
    if not unformatted:
        print("✅ All .go files are gofmt-clean.")
        return 0
    print_dry_run_findings(unformatted)
    if args.check_only:
        print("\nFix locally with:  cd gitmap && gofmt -w .", file=sys.stderr)
        return 1
    print(f"🛠️  Applying 'gofmt -w' to {len(unformatted)} file(s)…")
    apply_gofmt_fix(target_dir, unformatted)
    print(f"✅ Formatted {len(unformatted)} .go file(s).")
    handle_auto_commit_and_push(repo_root, args)
    return 0


def main() -> None:
    """Main CLI entry point for gofmt checker and fixer."""
    configure_io_encoding()
    repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
    parser = create_arg_parser(repo_root)
    args = parser.parse_args()
    target_dir = os.path.abspath(args.target_dir)

    status = execute_check_and_fix(target_dir, repo_root, args)
    sys.exit(status)


if __name__ == "__main__":
    main()
