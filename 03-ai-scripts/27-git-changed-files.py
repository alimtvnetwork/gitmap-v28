#!/usr/bin/env python3
"""
27-git-changed-files.py — Fast Git Commit History & Working-Tree Changed Files Extractor
Extracts files changed in the last N commits or incremental commit range and current working tree.
Deduplicates with dictionary lookup, saves last_commit_hash checkpoint, and exports to JSON, YAML, and TXT.
"""
from __future__ import annotations

import argparse
from importlib import import_module
import json
import os
from pathlib import Path
import sys
import time
from typing import Any

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")

sys.path.insert(0, str(Path(__file__).parent))
engine = import_module("02-shared-engine")

extract_git_changed_files = engine.extract_git_changed_files
DEFAULT_COMMITS = 20
DEFAULT_OUT_DIR = ".lovable/temp"


def parse_cli_args() -> argparse.Namespace:
    """Parses command line arguments for git changed files extractor."""
    parser = argparse.ArgumentParser(description="Extract changed files from git commit window.")
    parser.add_argument("--commits", "-n", type=int, default=DEFAULT_COMMITS, help="Number of commits (default: 20)")
    parser.add_argument("--output-dir", "-o", type=str, default=DEFAULT_OUT_DIR, help="Output folder")
    parser.add_argument("--format", type=str, choices=["all", "json", "yaml", "txt"], default="all", help="Output format")
    parser.add_argument("--no-incremental", "--full-window", action="store_true", help="Bypass checkpoint and force full window")
    parser.add_argument("--checkpoint", type=str, default=None, help="Custom checkpoint JSON file path")
    parser.add_argument("--since-commit", type=str, default=None, help="Diff specifically from given commit hash to HEAD")
    parser.add_argument("--quiet", "-q", action="store_true", help="Quiet mode (no banner)")
    parser.add_argument("--verify", action="store_true", help="Verify files exist on filesystem")

    return parser.parse_args()


def serialize_yaml(records: list[dict[str, Any]], metadata: dict[str, Any]) -> str:
    """Serializes changed file records and checkpoint metadata into standard YAML format."""
    lines = [
        f"last_commit_hash: \"{metadata.get('last_commit_hash', '')}\"",
        f"commit_range: \"{metadata.get('commit_range', '')}\"",
        f"is_incremental: {str(metadata.get('is_incremental', False)).lower()}",
        f"total_files: {len(records)}",
        "changed_files:",
    ]
    for rec in records:
        lines.append(f"  - path: \"{rec['path']}\"")
        lines.append(f"    status: \"{rec.get('status', 'unknown')}\"")
        lines.append(f"    extension: \"{rec.get('extension', '')}\"")
    output = "\n".join(lines) + "\n"

    return output


def serialize_txt(records: list[dict[str, Any]], metadata: dict[str, Any]) -> str:
    """Serializes file paths line-by-line with header comments for fast iteration."""
    lines = [
        f"# last_commit_hash: {metadata.get('last_commit_hash', '')}",
        f"# commit_range: {metadata.get('commit_range', '')}",
        f"# total_files: {len(records)}",
    ]
    lines.extend(rec["path"] for rec in records)
    output = "\n".join(lines) + "\n"

    return output


def serialize_json(records: list[dict[str, Any]], metadata: dict[str, Any]) -> str:
    """Serializes changed files and commit checkpoint into machine-readable JSON."""
    payload = dict(metadata)
    payload["total_files"] = len(records)
    payload["files"] = records
    output = json.dumps(payload, indent=2, ensure_ascii=False) + "\n"

    return output


def write_export_file(target_path: Path, content: str) -> None:
    """Writes content to target path creating parent directories if needed."""
    target_path.parent.mkdir(parents=True, exist_ok=True)
    target_path.write_text(content, encoding="utf-8")


def export_files_bundle(
    records: list[dict[str, Any]],
    metadata: dict[str, Any],
    out_dir: Path,
    fmt: str,
) -> dict[str, Path]:
    """Exports files in selected formats into destination folder."""
    written: dict[str, Path] = {}
    if fmt in ("all", "json"):
        p_json = out_dir / "git-changed-files.json"
        write_export_file(p_json, serialize_json(records, metadata))
        written["json"] = p_json
    if fmt in ("all", "yaml"):
        p_yaml = out_dir / "git-changed-files.yaml"
        write_export_file(p_yaml, serialize_yaml(records, metadata))
        written["yaml"] = p_yaml
    if fmt in ("all", "txt"):
        p_txt = out_dir / "git-changed-files.txt"
        write_export_file(p_txt, serialize_txt(records, metadata))
        written["txt"] = p_txt

    return written


def filter_existing_files(records: list[dict[str, Any]], repo_root: Path) -> list[dict[str, Any]]:
    """Filters records to only include files currently present on disk."""
    existing = [r for r in records if (repo_root / r["path"]).is_file()]

    return existing


def print_summary(
    records: list[dict[str, Any]],
    metadata: dict[str, Any],
    written: dict[str, Path],
    elapsed_sec: float,
) -> None:
    """Prints execution summary banner and checkpoint details."""
    mode_str = "INCREMENTAL (DELTA RANGE)" if metadata.get("is_incremental") else "FULL WINDOW"
    head_sha = metadata.get("last_commit_hash", "")[:12]
    prev_sha = (metadata.get("previous_checkpoint_hash") or "none")[:12]
    print("=" * 64)
    print("       GIT CHANGED FILES EXTRACTOR (CHECKPOINTED DELTA ENGINE)")
    print("=" * 64)
    print(f"🎯 Current HEAD Commit     : {head_sha}")
    print(f"📌 Previous Checkpoint SHA : {prev_sha}")
    print(f"🔄 Extraction Mode         : {mode_str}")
    print(f"📐 Commit Range Evaluated  : {metadata.get('commit_range', 'unknown')}")
    print(f"📦 Total Deduplicated Files: {len(records)}")
    print(f"⏱  Extraction Duration     : {elapsed_sec:.3f}s")
    for fmt_name, path in written.items():
        print(f"📄 [{fmt_name.upper():4s}] Exported to: {path}")
    print("=" * 64)


def resolve_checkpoint_file(args: argparse.Namespace, out_dir: Path) -> Path:
    """Resolves checkpoint JSON file path from args or default location."""
    if args.checkpoint:
        return Path(args.checkpoint)

    return out_dir / "git-changed-files.json"


def main() -> int:
    """Main CLI entrypoint."""
    args = parse_cli_args()
    start_time = time.perf_counter()
    repo_root = Path(__file__).resolve().parent.parent
    out_dir = repo_root / args.output_dir
    chk_file = resolve_checkpoint_file(args, out_dir)
    records, metadata = extract_git_changed_files(
        repo_root,
        commit_count=args.commits,
        checkpoint_file=chk_file,
        force_full=args.no_incremental,
        since_commit=args.since_commit,
    )
    if args.verify:
        records = filter_existing_files(records, repo_root)
    written = export_files_bundle(records, metadata, out_dir, args.format)
    elapsed = time.perf_counter() - start_time
    if not args.quiet:
        print_summary(records, metadata, written, elapsed)

    return 0


if __name__ == "__main__":
    sys.exit(main())
