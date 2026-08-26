#!/usr/bin/env python3
"""
check-prompts-loaded.py — Verify that the prompt-loading contract from
``.lovable/coding-guidelines/coding-guidelines.md`` actually holds on
disk: the index file ``.lovable/prompts.md`` must exist, and every
prompt file under ``.lovable/prompts/`` must be referenced by it.

This is the on-disk equivalent of the AI-side "read all prompts before
generating code" rule. If the index drifts from the prompt directory
the AI cannot honour the rule, so we fail the build instead.

Failure modes detected:

  * The index file is missing.
  * The prompts directory is missing or empty.
  * A prompt file exists on disk but is not referenced by the index
    ("orphan prompt").
  * The index references a prompt filename that does not exist on disk
    ("dangling reference").

Exit codes:
  0  PASS — index exists and matches the prompt directory exactly.
  1  FAIL — at least one drift detected.
  2  ERROR — paths missing or unreadable.

Usage::

    python3 linter-scripts/check-prompts-loaded.py
    python3 linter-scripts/check-prompts-loaded.py \\
        --index .lovable/prompts.md \\
        --prompts-dir .lovable/prompts
"""
from __future__ import annotations

import argparse
import re
import sys
from pathlib import Path

EXIT_PASS = 0
EXIT_FAIL = 1
EXIT_ERROR = 2

DEFAULT_INDEX = ".lovable/prompts.md"
DEFAULT_PROMPTS_DIR = ".lovable/prompts"


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Verify .lovable/prompts.md indexes every prompt file.",
    )
    parser.add_argument("--index", default=DEFAULT_INDEX,
                        help=f"Path to the prompt index (default: {DEFAULT_INDEX})")
    parser.add_argument("--prompts-dir", default=DEFAULT_PROMPTS_DIR,
                        help=f"Path to the prompts directory (default: {DEFAULT_PROMPTS_DIR})")
    return parser.parse_args()


def list_prompt_files(prompts_dir: Path) -> list[str]:
    """Return sorted prompt filenames (basenames) in ``prompts_dir``."""
    files = [p.name for p in prompts_dir.iterdir() if p.is_file()]
    return sorted(files)


def find_orphans(prompt_files: list[str], index_body: str) -> list[str]:
    """Prompts on disk that the index never names."""
    return [f for f in prompt_files if f not in index_body]


PROMPT_REF_RE = re.compile(r"\b(\d{2,}-[A-Za-z0-9._-]+\.[A-Za-z0-9]+)\b")


def find_dangling(prompt_files: list[str], index_body: str) -> list[str]:
    """Filenames the index mentions that don't exist on disk."""
    on_disk = set(prompt_files)
    mentioned = {m.group(1) for m in PROMPT_REF_RE.finditer(index_body)}
    return sorted(m for m in mentioned if m not in on_disk)


def report(prompt_files: list[str], orphans: list[str],
           dangling: list[str]) -> bool:
    """Print a human-readable summary and return ``True`` on PASS."""
    print(f"🔎 Prompts on disk: {len(prompt_files)}")
    for f in prompt_files:
        print(f"   • {f}")
    print("")
    has_orphans = bool(orphans)
    has_dangling = bool(dangling)
    if has_orphans:
        print(f"❌ Orphan prompts (not referenced by the index) ({len(orphans)}):")
        for f in orphans:
            print(f"   • {f}")
        print("")
    if has_dangling:
        print(f"❌ Dangling references (in index, missing on disk) ({len(dangling)}):")
        for f in dangling:
            print(f"   • {f}")
        print("")
    is_pass = not has_orphans and not has_dangling
    if is_pass:
        print("✅ Prompt index is in sync with .lovable/prompts/.")
    return is_pass


def main() -> int:
    import sys
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    args = parse_args()
    index_path = Path(args.index)
    prompts_dir = Path(args.prompts_dir)
    if not index_path.is_file():
        sys.stderr.write(
            f"ERROR: prompt index not found at {index_path}.\n"
            "       Create it per .lovable/coding-guidelines/coding-guidelines.md.\n"
        )
        return EXIT_ERROR
    if not prompts_dir.is_dir():
        sys.stderr.write(f"ERROR: prompts directory not found at {prompts_dir}.\n")
        return EXIT_ERROR
    prompt_files = list_prompt_files(prompts_dir)
    if not prompt_files:
        sys.stderr.write(f"ERROR: no prompt files found under {prompts_dir}.\n")
        return EXIT_ERROR
    index_body = index_path.read_text(encoding="utf-8")
    orphans = find_orphans(prompt_files, index_body)
    dangling = find_dangling(prompt_files, index_body)
    is_pass = report(prompt_files, orphans, dangling)
    return EXIT_PASS if is_pass else EXIT_FAIL


if __name__ == "__main__":
    sys.exit(main())
