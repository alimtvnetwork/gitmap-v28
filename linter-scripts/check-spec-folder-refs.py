#!/usr/bin/env python3
"""
check-spec-folder-refs.py
=========================

Fail CI if any markdown file references a numbered spec folder
(e.g., `spec/12-cicd-pipeline-workflows/`) that does not exist on
disk and is not allowlisted.

The allowlist (`spec-folder-refs.allowlist`) supports two categories,
selected via `[external]` / `[doc-only]` section headers:

  [external]  Real folders that live in a sibling repository.
              References are valid; we just don't host them here.

  [doc-only]  Illustrative / historical names that do NOT exist
              anywhere. They are prose only and must never become
              live links.

Stale references — i.e. references to numbered folders that are
neither on disk nor in either allowlist section — produce a clearer,
categorized error message that includes:

  - The nearest existing folder (fuzzy match) when one is plausible.
  - A decision tree: "real folder typo? → fix the path" /
    "external sibling? → add under [external]" /
    "documentation-only? → add under [doc-only]".

Reference shapes detected
-------------------------

1. Absolute repo-relative   : `spec/NN-name/...`
2. Relative inside `spec/`  : `./NN-name/...`, `../NN-name/...`, etc.

A "numbered spec folder" matches `^\\d{2}-[a-z0-9-]+$`.

Exit codes
----------
  0 — all references resolve or are allowlisted
  1 — at least one stale folder reference was found
"""

from __future__ import annotations

import difflib
import re
import sys
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent
SPEC_ROOT = REPO_ROOT / "spec"
ALLOWLIST_PATH = Path(__file__).resolve().parent / "spec-folder-refs.allowlist"

NUMBERED_FOLDER_RE = re.compile(r"^\d{2}-[a-z0-9-]+$")
ABSOLUTE_REF_RE = re.compile(
    r"(?<![\w\-])spec/(\d{2}-[a-z0-9-]+)(?=[/)\s\"'#])"
)
RELATIVE_REF_RE = re.compile(
    r"(?<![\w/])(\.{1,2}(?:/\.{2})*)/(\d{2}-[a-z0-9-]+)(?=[/)\s\"'#])"
)

SECTION_EXTERNAL = "external"
SECTION_DOC_ONLY = "doc-only"
VALID_SECTIONS = {SECTION_EXTERNAL, SECTION_DOC_ONLY}

ARCHIVE_LEGACY_DOC_ONLY = {
    "04-generic-cli",
    "05-coding-guidelines",
    "07-generic-release",
    "08-generic-update",
    "09-pipeline",
    "12-consolidated-guidelines",
}


def is_numbered_folder(name: str) -> bool:
    """Return True when name matches the NN-kebab-case pattern."""
    return bool(NUMBERED_FOLDER_RE.match(name))


def list_existing_folders() -> set[str]:
    """Return numbered folder names directly under spec/."""
    if not SPEC_ROOT.is_dir():
        return set()
    return {
        entry.name
        for entry in SPEC_ROOT.iterdir()
        if entry.is_dir() and is_numbered_folder(entry.name)
    }


def is_section_header(line: str) -> bool:
    """Return True when the line is a `[section]` header."""
    return line.startswith("[") and line.endswith("]")


def parse_section_header(line: str) -> str:
    """Return the section name from a `[name]` header line."""
    return line[1:-1].strip().lower()


def load_allowlist() -> dict[str, set[str]]:
    """Parse the allowlist file into category → set-of-names.

    Entries appearing before any section header default to
    [external] for backward compatibility with the legacy flat format.
    Unknown section names are skipped with a warning.
    """
    buckets: dict[str, set[str]] = {
        SECTION_EXTERNAL: set(),
        SECTION_DOC_ONLY: set(),
    }
    if not ALLOWLIST_PATH.is_file():
        return buckets
    current = SECTION_EXTERNAL
    for raw in ALLOWLIST_PATH.read_text(encoding="utf-8").splitlines():
        line = raw.strip()
        if not line or line.startswith("#"):
            continue
        if is_section_header(line):
            section = parse_section_header(line)
            if section not in VALID_SECTIONS:
                print(f"::warning::unknown allowlist section [{section}] ignored")
                current = SECTION_EXTERNAL
                continue
            current = section
            continue
        buckets[current].add(line)
    return buckets


def iter_markdown_files() -> list[Path]:
    """Return every .md file in the repo, excluding noisy directories."""
    skip_dirs = {"node_modules", "dist", ".git", "release-artifacts"}
    files: list[Path] = []
    for path in REPO_ROOT.rglob("*.md"):
        parts = set(path.relative_to(REPO_ROOT).parts)
        if parts & skip_dirs:
            continue
        files.append(path)
    return files


def collect_absolute_refs(text: str) -> set[str]:
    """Collect `spec/NN-name` references from text."""
    return set(ABSOLUTE_REF_RE.findall(text))


def collect_relative_refs(file_path: Path, text: str) -> set[str]:
    """Resolve `./NN-name` / `../NN-name` refs to top-level folder names.

    Only references that resolve back into spec/ are returned.
    """
    rel_targets: set[str] = set()
    if "spec" not in file_path.parts:
        return rel_targets
    base = file_path.parent
    for dots, name in RELATIVE_REF_RE.findall(text):
        target = (base / dots / name).resolve()
        if SPEC_ROOT in target.parents and target.parent == SPEC_ROOT:
            rel_targets.add(name)
    return rel_targets


def find_nearest_folder(name: str, candidates: set[str]) -> str | None:
    """Return the closest matching folder name, or None if none plausible."""
    matches = difflib.get_close_matches(name, candidates, n=1, cutoff=0.6)
    return matches[0] if matches else None


def find_stale_refs(
    existing: set[str], allow: dict[str, set[str]]
) -> list[tuple[Path, str]]:
    """Return list of (file, missing-folder) tuples not allowlisted."""
    allowed = allow[SECTION_EXTERNAL] | allow[SECTION_DOC_ONLY]
    stale: list[tuple[Path, str]] = []
    for md_file in iter_markdown_files():
        try:
            text = md_file.read_text(encoding="utf-8", errors="replace")
        except OSError as exc:
            import sys as _sys; print(f"Error: {exc}", file=_sys.stderr)
            continue
        refs = collect_absolute_refs(text) | collect_relative_refs(md_file, text)
        for folder in sorted(refs):
            rel_parts = md_file.relative_to(REPO_ROOT).parts
            is_archive_doc = rel_parts[:2] == ("spec", "26-spec-outsides")
            if is_archive_doc and folder in ARCHIVE_LEGACY_DOC_ONLY:
                continue
            if folder in existing or folder in allowed:
                continue
            stale.append((md_file, folder))
    return stale


def render_guidance(folder: str, suggestion: str | None) -> list[str]:
    """Build the multi-line guidance block shown for one stale ref."""
    lines: list[str] = []
    if suggestion:
        lines.append(f"       💡 Did you mean spec/{suggestion}/ ?")
    lines.append("       Resolve by choosing one:")
    lines.append("         (a) Typo / rename     → fix the path in the markdown.")
    lines.append("         (b) Sibling-repo ref  → add under [external] in")
    lines.append("                                 spec-folder-refs.allowlist.")
    lines.append("         (c) Documentation-only → add under [doc-only] in")
    lines.append("                                 spec-folder-refs.allowlist.")
    return lines


def print_report(
    stale: list[tuple[Path, str]],
    existing: set[str],
    allow: dict[str, set[str]],
) -> None:
    """Render a CI-friendly categorized report."""
    print("=" * 70)
    print("  SPEC FOLDER REFERENCE CHECK")
    print("=" * 70)
    print(f"  Existing numbered spec folders : {len(existing)}")
    print(f"  Allowlisted [external] folders : {len(allow[SECTION_EXTERNAL])}")
    print(f"  Allowlisted [doc-only] folders : {len(allow[SECTION_DOC_ONLY])}")
    print(f"  Stale references found         : {len(stale)}")
    print("=" * 70)
    if not stale:
        print("  ✅ All spec/NN-name references resolve or are allowlisted.")
        return
    for md_file, folder in stale:
        rel = md_file.relative_to(REPO_ROOT)
        suggestion = find_nearest_folder(folder, existing)
        print(f"  🔴 {rel}")
        print(f"       → spec/{folder}/  (folder does not exist)")
        for line in render_guidance(folder, suggestion):
            print(line)
    print("=" * 70)


def main() -> int:
    """Entry point."""
    import sys
    if hasattr(sys.stdout, 'reconfigure'):
        sys.stdout.reconfigure(encoding='utf-8')
    if not SPEC_ROOT.is_dir():
        print(f"::error::spec/ directory not found at {SPEC_ROOT}")
        return 1
    existing = list_existing_folders()
    allow = load_allowlist()
    stale = find_stale_refs(existing, allow)
    print_report(stale, existing, allow)
    return 1 if stale else 0


if __name__ == "__main__":
    sys.exit(main())
