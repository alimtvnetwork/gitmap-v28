#!/usr/bin/env python3
"""
Forbidden-Strings Allowlist CLI
================================

Adds repo-relative file paths to the ``allowlist`` of a specific
``[[rule]]`` block in ``linter-scripts/forbidden-strings.toml``.

Use this when a forbidden-strings violation is **legitimate** — for
example, a historical session log that must literally cite the old
slug, or a changelog entry quoting a deprecated module path.

Guard rails (every one is intentional — do NOT bypass):

1. ``--reason`` is mandatory and non-empty. The reason is written into
   the TOML as a one-line comment above the inserted entries so future
   readers know why the waiver exists.
2. Each path must exist on disk.
3. Each path must actually contain a match for the rule's regex.
   Allowlisting a path that does not match is meaningless and usually
   indicates a typo — refuse it loudly.
4. Duplicates are skipped silently (idempotent).
5. The TOML is edited in place as text — we deliberately do **not**
   round-trip through a TOML writer so handcrafted comments,
   grouping, and trailing-comma style are preserved.

USAGE
-----

    # Explicit paths
    python3 linter-scripts/allowlist-forbidden-string.py \\
        --rule STALE-REPO-SLUG \\
        --reason "Audit trail of the v14 -> v17 rebrand." \\
        .lovable/memory/sessions/2026-04-24-batch-cleanup-and-rebrand.md

    # Auto mode: scan the repo and waive every current finding under one reason
    python3 linter-scripts/allowlist-forbidden-string.py \\
        --rule STALE-REPO-SLUG \\
        --reason "Historical migration documentation." \\
        --auto

    # Dry-run (preview only, no write)
    python3 ... --rule STALE-REPO-SLUG --reason "..." --auto --dry-run

    # Re-verify after applying
    python3 ... --rule STALE-REPO-SLUG --reason "..." --auto --check

EXIT CODES
----------
    0  success (paths added, or already present, or dry-run preview)
    1  generic failure (unknown rule, IO error, etc.)
    2  invalid usage (missing/empty reason, no paths supplied)
    3  one or more paths failed validation (missing on disk, or no
       regex match — see stderr for the per-path reason)
    4  --check requested and the linter still reports violations after
       the write (means new findings appeared OR a path was added that
       the scanner doesn't actually consider waived)
"""

from __future__ import annotations

import argparse
import os
import re
import subprocess
import sys
from pathlib import Path

try:
    import tomllib  # Python 3.11+
except ModuleNotFoundError:
    try:
        import tomli as tomllib  # type: ignore[no-redef]
    except ModuleNotFoundError as exc:
        import sys as _sys; print(f"Error: {exc}", file=_sys.stderr)
        sys.exit("Error: Python 3.11+ required (tomllib), or install 'tomli'.")

REPO_ROOT = Path(__file__).resolve().parent.parent
TOML_PATH = REPO_ROOT / "linter-scripts" / "forbidden-strings.toml"
CHECKER_PATH = REPO_ROOT / "linter-scripts" / "check-forbidden-strings.py"

# Directories the scanner skips unconditionally — mirror
# ``check-forbidden-strings.py::ALWAYS_EXCLUDE_DIRS`` so --auto produces
# the same finding set as CI.
ALWAYS_EXCLUDE_DIRS = {".git", "node_modules", "dist", "build"}


# --------------------------------------------------------------------------- #
# Config loading
# --------------------------------------------------------------------------- #

def load_rules() -> list[dict]:
    with TOML_PATH.open("rb") as fh:
        return tomllib.load(fh).get("rule", [])


def find_rule(rules: list[dict], rule_id: str) -> dict:
    for rule in rules:
        if rule.get("id") == rule_id:
            return rule
    available = ", ".join(r.get("id", "?") for r in rules) or "(none)"
    raise SystemExit(
        f"Error: rule id '{rule_id}' not found in {TOML_PATH.name}. "
        f"Available: {available}"
    )


# --------------------------------------------------------------------------- #
# Repo scan (mirrors check-forbidden-strings.py logic, minus exclude_files
# globbing — we only need *which paths match*, not full reporting)
# --------------------------------------------------------------------------- #

def is_excluded_dir(rel_dir: str, extra_excludes: list[str]) -> bool:
    parts = rel_dir.replace("\\", "/").split("/")
    excludes = ALWAYS_EXCLUDE_DIRS | set(extra_excludes)
    return any(p in excludes for p in parts)


def find_matching_paths(rule: dict) -> list[str]:
    """Return repo-relative paths that contain at least one match for
    ``rule['pattern']`` and are not already covered by an exclude/allow."""
    import fnmatch

    pattern = re.compile(rule["pattern"])
    extra_excl_dirs = rule.get("exclude_dirs", [])
    excl_files = rule.get("exclude_files", [])
    allowlist = rule.get("allowlist", [])

    def already_waived(rel: str) -> bool:
        for pat in allowlist:
            if fnmatch.fnmatch(rel, pat):
                return True
            if rel == pat or rel.startswith(pat + "/"):
                return True
        return False

    def file_excluded(name: str) -> bool:
        return any(fnmatch.fnmatch(name, pat) for pat in excl_files)

    hits: list[str] = []
    for dirpath, dirnames, filenames in os.walk(REPO_ROOT):
        rel_dir = os.path.relpath(dirpath, REPO_ROOT)
        if rel_dir == ".":
            rel_dir = ""
        if rel_dir and is_excluded_dir(rel_dir, extra_excl_dirs):
            dirnames.clear()
            continue
        dirnames[:] = [
            d for d in dirnames
            if d not in ALWAYS_EXCLUDE_DIRS and d not in extra_excl_dirs
        ]
        for fname in filenames:
            if file_excluded(fname):
                continue
            full = Path(dirpath) / fname
            rel = full.relative_to(REPO_ROOT).as_posix()
            if already_waived(rel):
                continue
            try:
                text = full.read_text(encoding="utf-8", errors="ignore")
            except OSError as exc:
                import sys as _sys; print(f"Error: {exc}", file=_sys.stderr)
                continue
            if pattern.search(text):
                hits.append(rel)
    return sorted(set(hits))


# --------------------------------------------------------------------------- #
# Path validation
# --------------------------------------------------------------------------- #

def validate_paths(rule: dict, paths: list[str]) -> tuple[list[str], list[str]]:
    """Split ``paths`` into ``(valid, errors)``.

    A path is valid iff it exists AND its contents match the rule regex.
    """
    pattern = re.compile(rule["pattern"])
    valid: list[str] = []
    errors: list[str] = []
    for raw in paths:
        rel = raw.replace("\\", "/").lstrip("./")
        full = REPO_ROOT / rel
        if not full.exists():
            errors.append(f"{rel}: does not exist on disk")
            continue
        if not full.is_file():
            errors.append(f"{rel}: not a regular file")
            continue
        try:
            text = full.read_text(encoding="utf-8", errors="ignore")
        except OSError as exc:
            import sys as _sys; print(f"Error: {exc}", file=_sys.stderr)
            errors.append(f"{rel}: read error: {exc}")
            continue
        if not pattern.search(text):
            errors.append(
                f"{rel}: does not contain a match for /{rule['pattern']}/ "
                f"(refusing — allowlisting this would be a no-op)"
            )
            continue
        valid.append(rel)
    return valid, errors


# --------------------------------------------------------------------------- #
# TOML in-place edit
# --------------------------------------------------------------------------- #

# Match the start of the target [[rule]] block: a line ``[[rule]]`` followed
# (within the block, before the next ``[[rule]]`` or EOF) by ``id = "..."``.
RULE_HEADER_RE = re.compile(r"^\[\[rule\]\]\s*$", re.MULTILINE)
ID_LINE_RE = re.compile(r'^id\s*=\s*"([^"]+)"\s*$', re.MULTILINE)
ALLOWLIST_OPEN_RE = re.compile(r"^allowlist\s*=\s*\[", re.MULTILINE)


def locate_rule_block(toml_text: str, rule_id: str) -> tuple[int, int]:
    """Return ``(block_start, block_end)`` byte offsets for the ``[[rule]]``
    block whose id matches ``rule_id``."""
    starts = [m.start() for m in RULE_HEADER_RE.finditer(toml_text)]
    if not starts:
        raise SystemExit("Error: no [[rule]] blocks found in TOML.")
    boundaries = starts + [len(toml_text)]
    for i, begin in enumerate(starts):
        end = boundaries[i + 1]
        block = toml_text[begin:end]
        m = ID_LINE_RE.search(block)
        if m and m.group(1) == rule_id:
            return begin, end
    raise SystemExit(f"Error: rule id '{rule_id}' not found in TOML text.")


def insert_into_allowlist(
    toml_text: str,
    rule_id: str,
    new_paths: list[str],
    reason: str,
) -> tuple[str, list[str]]:
    """Insert ``new_paths`` into the ``allowlist = [...]`` of ``rule_id``.

    Returns ``(new_text, actually_added)``. Skips paths already present.
    Creates the ``allowlist`` key if the rule doesn't have one yet
    (placed just before the closing of the block).
    """
    block_start, block_end = locate_rule_block(toml_text, rule_id)
    block = toml_text[block_start:block_end]

    open_match = ALLOWLIST_OPEN_RE.search(block)
    if open_match is None:
        # No allowlist key yet — append one at the end of the block.
        existing: list[str] = []
    else:
        # Find matching ``]`` for the ``[`` after ``allowlist =``. Bracket
        # nesting is unlikely inside this list, but count just in case.
        i = open_match.end()  # position right after ``[``
        depth = 1
        while i < len(block) and depth > 0:
            ch = block[i]
            if ch == "[":
                depth += 1
            elif ch == "]":
                depth -= 1
                if depth == 0:
                    break
            i += 1
        if depth != 0:
            raise SystemExit(
                f"Error: malformed allowlist for rule '{rule_id}' "
                f"(missing closing ']')."
            )
        existing_block = block[open_match.end():i]
        existing = re.findall(r'"([^"]+)"', existing_block)

    existing_set = set(existing)
    to_add = [p for p in new_paths if p not in existing_set]
    if not to_add:
        return toml_text, []

    safe_reason = reason.replace("\n", " ").strip()
    indent = "  "
    new_entries = [f"{indent}# {safe_reason}"]
    new_entries += [f'{indent}"{p}",' for p in to_add]
    new_block_chunk = "\n" + "\n".join(new_entries)

    if open_match is None:
        # Append a brand-new allowlist key at the end of the block.
        # Strip trailing whitespace from the block, append the key, then
        # restore one trailing newline so the next [[rule]] stays separated.
        stripped = block.rstrip()
        new_key = f"\nallowlist = [{new_block_chunk}\n]\n"
        new_block = stripped + new_key
        # Preserve the original trailing whitespace style (one blank line
        # before the next block, if any).
        trailing = block[len(block.rstrip()):]
        if trailing and not new_block.endswith(trailing):
            new_block += trailing.lstrip("\n")
    else:
        # Insert before the closing ``]`` of the existing allowlist.
        # ``i`` points at the ``]`` character.
        before = block[:i].rstrip()
        after = block[i:]  # starts with ``]``
        # Drop a trailing comma-less entry's lack-of-comma situation by
        # always emitting our entries with trailing commas. The previous
        # last entry may or may not end with ``,`` — TOML accepts both, so
        # we don't force-rewrite it.
        new_block = before + new_block_chunk + "\n" + after

    new_text = toml_text[:block_start] + new_block + toml_text[block_end:]
    return new_text, to_add


# --------------------------------------------------------------------------- #
# CLI
# --------------------------------------------------------------------------- #

def parse_args() -> argparse.Namespace:
    p = argparse.ArgumentParser(
        prog="allowlist-forbidden-string",
        description=(
            "Add file paths to the allowlist of a [[rule]] in "
            "linter-scripts/forbidden-strings.toml. A non-empty "
            "--reason is mandatory."
        ),
    )
    p.add_argument(
        "--rule", required=True,
        help="Rule id (e.g. STALE-REPO-SLUG). Must exist in the TOML.",
    )
    p.add_argument(
        "--reason", required=True,
        help=(
            "Justification for the waiver. Written into the TOML as a "
            "one-line comment above the new entries. Required."
        ),
    )
    p.add_argument(
        "paths", nargs="*",
        help="Repo-relative paths to allowlist.",
    )
    p.add_argument(
        "--auto", action="store_true",
        help=(
            "Scan the repo for current findings of --rule and waive "
            "all of them under --reason. Combine with --dry-run to "
            "preview."
        ),
    )
    p.add_argument(
        "--dry-run", action="store_true",
        help="Show what would change but do not write the TOML.",
    )
    p.add_argument(
        "--check", action="store_true",
        help=(
            "After writing, re-run check-forbidden-strings.py and exit "
            "with code 4 if any findings remain."
        ),
    )
    return p.parse_args()


def main() -> int:
    args = parse_args()

    if not args.reason.strip():
        print("Error: --reason must be a non-empty string.", file=sys.stderr)
        return 2

    rules = load_rules()
    rule = find_rule(rules, args.rule)

    if args.auto:
        if args.paths:
            print(
                "Error: --auto cannot be combined with explicit paths "
                "(remove one or the other).",
                file=sys.stderr,
            )
            return 2
        candidate_paths = find_matching_paths(rule)
        if not candidate_paths:
            print(
                f"No outstanding findings for rule '{args.rule}' — "
                f"nothing to allowlist."
            )
            return 0
        print(f"Auto mode: {len(candidate_paths)} candidate path(s) found.")
    else:
        if not args.paths:
            print(
                "Error: provide at least one path, or pass --auto to "
                "scan the repo.",
                file=sys.stderr,
            )
            return 2
        candidate_paths = args.paths

    valid, errors = validate_paths(rule, candidate_paths)

    if errors:
        print("Path validation errors:", file=sys.stderr)
        for err in errors:
            print(f"  - {err}", file=sys.stderr)
        # Hard fail unless every path was valid. We refuse to write a
        # partial allowlist when the user clearly intended the whole set.
        if not valid:
            return 3
        print(
            f"\n{len(valid)} valid path(s) will still be applied; "
            f"{len(errors)} rejected.",
            file=sys.stderr,
        )

    if not valid:
        print("Nothing to do.")
        return 3 if errors else 0

    toml_text = TOML_PATH.read_text(encoding="utf-8")
    new_text, added = insert_into_allowlist(
        toml_text, args.rule, valid, args.reason,
    )

    if not added:
        print(
            f"All {len(valid)} path(s) already present in "
            f"allowlist for '{args.rule}'. No changes."
        )
        return 0

    print(f"Adding {len(added)} path(s) to allowlist for '{args.rule}':")
    for p in added:
        print(f"  + {p}")
    print(f"Reason: {args.reason.strip()}")

    if args.dry_run:
        print("\n[dry-run] TOML not written.")
        return 0

    TOML_PATH.write_text(new_text, encoding="utf-8")
    print(f"\nWrote {TOML_PATH.relative_to(REPO_ROOT)}.")

    if args.check:
        print("\nRe-running check-forbidden-strings.py ...")
        result = subprocess.run(
            [sys.executable, str(CHECKER_PATH)],
            cwd=REPO_ROOT,
        )
        if result.returncode != 0:
            print(
                "Linter still reports findings — exiting 4.",
                file=sys.stderr,
            )
            return 4
        print("Linter clean.")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
