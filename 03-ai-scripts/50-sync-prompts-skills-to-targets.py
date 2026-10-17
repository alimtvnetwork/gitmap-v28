#!/usr/bin/env python3
"""Prompts + Skills sync to target repos (repos that do not need spec).

WHAT THIS SCRIPT SYNCS (from coding-guidelines-v24 into each target):
  1. 01-prompts/          -> <target>/01-prompts/
  2. .agents/skills/      -> <target>/.agents/skills/
  3. .cursor/skills/      -> <target>/.cursor/skills/

WHAT THIS SCRIPT NEVER SYNCS (explicit fail-safe, default-deny):
  - 02-spec/              NEVER synced. Not shared specs, not app specs, nothing.
                          The allowlist contains ONLY the three roots above; any path
                          outside them is denied by default. The is_denied() guard
                          below is checked on EVERY copy and delete path.
  - .ai-memory/           Target repos own their memory/plans; never touched.
  - 03-ai-scripts/        Sync scripts ship via the fleet sync (38) or are copied
                          deliberately by hand; this script stays scoped to
                          prompts + skills only.

HOW TO ADD A TARGET:
  Option A (registry file): add an entry to scripts/prompt-sync-targets.json:
      {"name": "my-repo", "path": "/absolute/path/to/my-repo"}
  Option B (CLI): --target /path/to/my-repo  (repeatable; bypasses the registry)
  The registry ships with story-telling-v1 as the example target.

RE-RUNNABLE + IDEMPOTENT:
  - Files are compared by content hash; identical files are skipped (no-op).
  - Canonical skills/prompts are updated in place; target-local skills (folders
    that exist only in the target) are KEPT, never deleted (see 0c5fad76).
  - Running twice in a row: second run reports 0 changes. Safe to re-run anytime.

USAGE:
  gitmap py 03-ai-scripts/50-sync-prompts-skills-to-targets.py --dry-run
  gitmap py 03-ai-scripts/50-sync-prompts-skills-to-targets.py --target ~/workspace/custom-skills
  gitmap py 03-ai-scripts/50-sync-prompts-skills-to-targets.py --target story-telling-v1 --no-commit
"""
from __future__ import annotations

import argparse
import hashlib
import json
import os
from pathlib import Path
import shutil
import sys

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")

SOURCE_ROOT = Path(__file__).resolve().parent.parent
DEFAULT_CONFIG = SOURCE_ROOT / "scripts" / "prompt-sync-targets.json"

# --- Allowlist: the ONLY trees this script may ever write into a target. ---
SYNC_ROOTS = (
    "01-prompts",
    ".agents/skills",
    ".cursor/skills",
)

# --- Fail-safe default-deny: paths that must NEVER be written, walked, or deleted. ---
DENIED_ROOTS = (
    "02-spec",          # ALL specs: shared (01-20) and app (21+) alike. Never sync spec.
    ".ai-memory",       # target-owned memory/plans.
    "03-ai-scripts",    # sync scripts ship via fleet sync (38) or by hand.
    ".git",
    "__pycache__",
    ".pytest_cache",
)

DENIED_EXTS = {".pyc", ".pyo", ".tmp"}


def is_denied(rel_path: str) -> bool:
    """Fail-safe default-deny gate. rel_path is relative to the target root.

    Any path whose first segment (or any segment) matches a denied root is
    rejected. 02-spec/ is denied unconditionally — this is the program's
    hard boundary: repos that do not need spec must never receive spec.
    """
    parts = Path(rel_path.replace("\\", "/")).parts
    lowered = [p.lower() for p in parts]

    for denied in DENIED_ROOTS:
        if denied.lower() in lowered:
            return True

    name = parts[-1] if parts else ""
    if name in (".DS_Store",):
        return True

    if Path(name).suffix.lower() in DENIED_EXTS:
        return True

    return False


def is_sync_root(rel_path: str) -> bool:
    """Default-allow ONLY the three sync roots; deny everything else."""
    norm = rel_path.replace("\\", "/")
    return any(norm == root or norm.startswith(root + "/") for root in SYNC_ROOTS)


def file_hash(path: Path) -> str:
    """SHA-256 of file bytes; empty string when unreadable."""
    try:
        digest = hashlib.sha256()
        with path.open("rb") as handle:
            for chunk in iter(lambda: handle.read(65536), b""):
                digest.update(chunk)
        return digest.hexdigest()
    except OSError:
        return ""


def sync_tree(src: Path, dst: Path, is_dry_run: bool) -> dict:
    """Mirror src into dst. Returns counts. Never deletes (additive-safe).

    - Canonical files: updated when content differs, skipped when identical.
    - New files: copied.
    - Target-local files (not in src): KEPT, counted as 'kept_local'.
    - Anything hitting is_denied(): skipped and counted as 'skipped_denied'
      (this is how 02-spec/ is provably excluded).
    """
    stats = {"copied": 0, "updated": 0, "skipped": 0, "kept_local": 0,
             "skipped_denied": 0, "errors": []}

    if not src.exists():
        stats["errors"].append(f"source tree missing: {src}")
        return stats

    if not is_dry_run:
        dst.mkdir(parents=True, exist_ok=True)

    # 1. Walk source: copy new / update changed.
    for root, dirs, files in os.walk(src):
        dirs[:] = [d for d in sorted(dirs)
                   if not is_denied((Path(root).relative_to(src) / d).as_posix())]

        rel_root = Path(root).relative_to(src)
        target_dir = dst / rel_root
        rel_dir_str = rel_root.as_posix()

        if is_denied(rel_dir_str):
            stats["skipped_denied"] += 1
            dirs[:] = []
            continue

        if not is_dry_run:
            target_dir.mkdir(parents=True, exist_ok=True)

        for name in sorted(files):
            src_file = Path(root) / name
            rel_file = (rel_root / name).as_posix()

            if is_denied(rel_file):
                stats["skipped_denied"] += 1
                continue

            dst_file = target_dir / name

            if not dst_file.exists():
                if not is_dry_run:
                    shutil.copy2(src_file, dst_file)
                stats["copied"] += 1
            elif file_hash(src_file) != file_hash(dst_file):
                if not is_dry_run:
                    shutil.copy2(src_file, dst_file)
                stats["updated"] += 1
            else:
                stats["skipped"] += 1

    # 2. Walk target: report target-local extras (kept, never deleted).
    if dst.exists():
        for root, dirs, files in os.walk(dst):
            dirs[:] = [d for d in sorted(dirs)
                       if not is_denied((Path(root).relative_to(dst) / d).as_posix())]

            rel_root = Path(root).relative_to(dst)

            if is_denied(rel_root.as_posix()):
                dirs[:] = []
                continue

            for name in sorted(files):
                rel_file = (rel_root / name).as_posix()

                if is_denied(rel_file):
                    continue

                if not (src / rel_file).exists():
                    stats["kept_local"] += 1

    return stats


def load_registry(config_path: Path) -> list[dict]:
    """Load the target registry file. Missing file -> empty list (CLI only)."""
    if not config_path.exists():
        return []

    try:
        data = json.loads(config_path.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError) as exc:
        raise SystemExit(f"ERROR: cannot read registry {config_path}: {exc}")

    targets = data.get("targets", [])
    if not isinstance(targets, list):
        raise SystemExit(f"ERROR: 'targets' must be a list in {config_path}")

    return targets


def resolve_target(entry: dict | str) -> tuple[str, Path]:
    """Return (name, path) for a registry entry or CLI path string."""
    if isinstance(entry, str):
        name, raw = Path(entry).name or entry, entry
    else:
        name = str(entry.get("name", "")).strip() or "unnamed"
        raw = str(entry.get("path", "")).strip()

    if not raw:
        raise SystemExit(f"ERROR: target '{name}' has no path.")

    path = Path(os.path.expandvars(os.path.expanduser(raw)))

    if not path.is_absolute():
        path = (SOURCE_ROOT.parent / path).resolve()

    return name, path


def sync_target(name: str, target: Path, is_dry_run: bool) -> dict:
    """Sync prompts+skills into one target. Returns the per-target report dict."""
    report: dict = {"name": name, "path": str(target), "status": "OK",
                    "error": None, "stats": {}, "denied_check": None}

    if not target.exists():
        report["status"] = "SKIP"
        report["error"] = "target directory does not exist"
        return report

    if not target.is_dir():
        report["status"] = "SKIP"
        report["error"] = "target is not a directory"
        return report

    # Pre-flight fail-safe proof: 02-spec/ must not land in the target via us.
    report["denied_check"] = "02-spec/ denied by is_denied() gate (fail-safe default-deny)"

    totals = {"copied": 0, "updated": 0, "skipped": 0, "kept_local": 0,
              "skipped_denied": 0, "errors": []}

    for root in SYNC_ROOTS:
        src = SOURCE_ROOT / root
        dst = target / root
        stats = sync_tree(src, dst, is_dry_run=is_dry_run)

        for key in ("copied", "updated", "skipped", "kept_local", "skipped_denied"):
            totals[key] += stats[key]

        totals["errors"].extend(stats["errors"])

    report["stats"] = totals

    if totals["errors"]:
        report["status"] = "FAIL"
        report["error"] = "; ".join(totals["errors"])

    report["changed"] = (totals["copied"] + totals["updated"]) > 0
    return report


def print_report(reports: list[dict], is_dry_run: bool) -> None:
    """Clear per-target report + summary table."""
    mode = "DRY-RUN (no writes)" if is_dry_run else "LIVE"
    print(f"\nMode: {mode}")
    print("=" * 100)

    for rep in reports:
        print(f"\nTarget : {rep['name']}  ({rep['path']})")
        print(f"Status : {rep['status']}")

        if rep["error"]:
            print(f"Error  : {rep['error']}")
            continue

        stats = rep["stats"]
        print(f"Fail-safe : {rep['denied_check']}")
        print(f"Copied new        : {stats['copied']}")
        print(f"Updated (changed) : {stats['updated']}")
        print(f"Skipped (identical): {stats['skipped']}")
        print(f"Kept target-local : {stats['kept_local']}")
        print(f"Denied/skipped    : {stats['skipped_denied']}  (02-spec/ and friends never sync)")

    print("\n" + "=" * 100)
    print("SUMMARY")
    print("=" * 100)
    print(f"{'Target':<28} | {'Status':<6} | {'New':>5} | {'Upd':>5} | {'Same':>6} | {'Local':>5} | {'Denied':>6} | Note")
    print("-" * 100)

    for rep in reports:
        name = str(rep["name"])[:28]
        status = str(rep["status"])

        if rep["error"]:
            print(f"{name:<28} | {status:<6} | {'-':>5} | {'-':>5} | {'-':>6} | {'-':>5} | {'-':>6} | {rep['error']}")
            continue

        stats = rep["stats"]
        note = "changed" if rep["changed"] else "clean (idempotent no-op)"

        print(f"{name:<28} | {status:<6} | {stats['copied']:>5} | {stats['updated']:>5} | "
              f"{stats['skipped']:>6} | {stats['kept_local']:>5} | {stats['skipped_denied']:>6} | {note}")


def main() -> None:
    parser = argparse.ArgumentParser(
        description="Sync canonical 01-prompts/ + .agents/skills/ + .cursor/skills/ "
                    "into target repos. NEVER syncs 02-spec/ (fail-safe default-deny)."
    )
    parser.add_argument("--target", "-t", action="append", default=[],
                        help="Target repo path (or registry name). Repeatable. "
                             "If omitted, the registry file is used.")
    parser.add_argument("--config", "-c", default=str(DEFAULT_CONFIG),
                        help=f"Target registry JSON (default: {DEFAULT_CONFIG})")
    parser.add_argument("--dry-run", action="store_true",
                        help="Preview changes without writing anything.")
    parser.add_argument("--list", action="store_true",
                        help="List registry targets and exit.")
    args = parser.parse_args()

    registry = load_registry(Path(args.config))

    if args.list:
        print(f"Registry: {args.config}")
        for entry in registry:
            name, path = resolve_target(entry)
            print(f"  - {name}: {path}")
        return

    # CLI --target wins; otherwise use the registry.
    if args.target:
        entries: list[dict | str] = list(args.target)
    else:
        entries = registry

    if not entries:
        print("No targets: registry is empty and --target was not given.")
        print(f"Registry file: {args.config}")
        return

    print(f"Source : coding-guidelines-v24 ({SOURCE_ROOT})")
    print(f"Syncs  : {', '.join(SYNC_ROOTS)}")
    print(f"NEVER  : 02-spec/ (fail-safe default-deny), .ai-memory/, 03-ai-scripts/")

    reports = []
    for entry in entries:
        name, path = resolve_target(entry)
        reports.append(sync_target(name, path, is_dry_run=args.dry_run))

    print_report(reports, is_dry_run=args.dry_run)

    failures = [r for r in reports if r["status"] == "FAIL"]

    if failures:
        raise SystemExit(1)


if __name__ == "__main__":
    main()
