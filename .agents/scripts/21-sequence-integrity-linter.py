#!/usr/bin/env python3
"""
21-sequence-integrity-linter.py
===============================
Unified AI tool to verify sequential file integrity and internal relative
links across prompts, execution plans, and agent skills.

Exits:
    0 - All sequence references and relative links resolve cleanly.
    1 - Missing or broken sequential file references detected.
"""

import os
import re
import sys
from pathlib import Path

# Tokens indicating generic template examples / documentation placeholders
PLACEHOLDER_TOKENS = (
    "<slug>", "<subslug>", "<work-slug>", "<date>", "<version>",
    "<category>", "<area>", "<name>", "<ext>", "<owner>", "<repo>",
    "XX-", "NN-", "01-<", "00-<", "{", "}", "*", "path/to/", "...",
    "<version-slug>", "<work_slug>", "<subtask_slug>", "<target>", "<module>",
    "/XX/", "xx-", "vX.Y.Z", "vX.", "/<", "XX"
)

# Target directories to audit for sequence integrity
AUDIT_DIRS = (
    ".lovable/prompts",
    ".lovable/plans",
    ".agents/skills",
    ".lovable/coding-guidelines",
    ".lovable/memory/standards"
)

# Directories/files explicitly exempt (e.g., historical archives and migration transaction logs)
EXEMPT_PATHS = {
    ".lovable/memory/transactions/spec-migration-transaction-log.md",
}

EXEMPT_DIR_PARTS = {
    "/done/", "/_archive/", "/completed/", "/sessions/"
}

# Regex to match markdown links: [text](target)
MD_LINK_RE = re.compile(r'\[([^\]]*)\]\(([^)\s]+)(?:\s+"[^"]*")?\)')

# Regex to match single-backtick path references: `spec/...` or `.lovable/...`
PATH_BACKTICK_RE = re.compile(r'`([^`\n]+)`')


def is_placeholder(target: str) -> bool:
    """Return True if the target contains placeholder tokens."""
    for token in PLACEHOLDER_TOKENS:
        if token in target:
            return True
    return False


def is_external_or_special(target: str) -> bool:
    """Return True if target is a web URL, anchor, email, or special URI."""
    clean = target.strip()
    return clean.startswith((
        "http://", "https://", "mailto:", "tel:", "ftp://",
        "conversation://", "file:///", "#", "javascript:"
    ))


def is_valid_file_extension(target: str) -> bool:
    """Return True if string ends with a standard tracked file extension."""
    valid_exts = (
        ".md", ".py", ".sh", ".ps1", ".mjs", ".js", ".ts", ".tsx",
        ".go", ".json", ".yaml", ".yml", ".toml", ".png", ".svg", ".gif"
    )
    return any(target.endswith(ext) for ext in valid_exts)


def is_exempt(rel_posix_path: str) -> bool:
    """Check if file is in the exemption set."""
    if rel_posix_path in EXEMPT_PATHS:
        return True
    return any(part in rel_posix_path for part in EXEMPT_DIR_PARTS)


def resolve_reference(citing_file: Path, raw_target: str, repo_root: Path) -> Path | None:
    """Attempt to resolve a target path against repo root and local directory."""
    clean = raw_target.split("#")[0].split("?")[0].strip().replace("\\", "/")
    if not clean or is_placeholder(clean) or is_external_or_special(clean):
        return None

    # Check relative to repository root
    root_resolved = (repo_root / clean).resolve()
    if root_resolved.exists():
        return root_resolved

    # Check relative to citing file's parent directory
    local_resolved = (citing_file.parent / clean).resolve()
    if local_resolved.exists():
        return local_resolved

    # Check if target is a sibling file in citing directory
    sibling_resolved = (citing_file.parent / Path(clean).name).resolve()
    if sibling_resolved.exists():
        return sibling_resolved

    return None


def audit_file(file_path: Path, repo_root: Path) -> list[tuple[int, str, str]]:
    """Audit single markdown file for broken sequence/link references."""
    rel_posix = file_path.resolve().relative_to(repo_root).as_posix()
    if is_exempt(rel_posix):
        return []

    try:
        content = file_path.read_text(encoding="utf-8", errors="replace")
    except Exception as read_err:
        return [(1, f"Unreadable file: {read_err}", rel_posix)]

    violations = []
    lines = content.splitlines()

    # Filter out fenced code blocks to avoid flagging illustrative code examples
    in_code_fence = False
    for line_idx, line in enumerate(lines, start=1):
        stripped = line.strip()
        if stripped.startswith("```") or stripped.startswith("~~~"):
            in_code_fence = not in_code_fence
            continue

        if in_code_fence:
            continue

        # 1. Check explicit markdown links
        for m in MD_LINK_RE.finditer(line):
            raw_target = m.group(2)
            if is_external_or_special(raw_target) or is_placeholder(raw_target):
                continue
            resolved = resolve_reference(file_path, raw_target, repo_root)
            if resolved is None:
                clean_target = raw_target.split("#")[0].strip()
                if clean_target and not is_placeholder(clean_target):
                    violations.append((line_idx, line.strip(), clean_target))

        # 2. Check backticked repository paths
        for m in PATH_BACKTICK_RE.finditer(line):
            candidate = m.group(1).strip()
            # Only test strings that look like actual file paths in tracked folders
            if candidate.startswith(("spec/", ".lovable/", ".agents/", "linter-scripts/", "scripts/")):
                if is_valid_file_extension(candidate) and not is_placeholder(candidate):
                    resolved = resolve_reference(file_path, candidate, repo_root)
                    if resolved is None:
                        violations.append((line_idx, line.strip(), candidate))

    return violations


def main() -> int:
    repo_root = Path(__file__).resolve().parent.parent
    total_files = 0
    total_violations = 0

    print("[INFO] Auditing sequential file integrity and internal links across prompts, plans, and skills...")

    for audit_dir_name in AUDIT_DIRS:
        target_dir = repo_root / audit_dir_name
        if not target_dir.exists():
            continue

        for root_str, _, file_names in os.walk(target_dir):
            for file_name in file_names:
                if not file_name.endswith(".md"):
                    continue

                full_path = Path(root_str) / file_name
                rel_path = full_path.resolve().relative_to(repo_root).as_posix()
                total_files += 1

                file_violations = audit_file(full_path, repo_root)
                if file_violations:
                    total_violations += len(file_violations)
                    print(f"\n[FAIL] {rel_path} ({len(file_violations)} broken reference(s)):")
                    for line_num, _, target in file_violations[:5]:
                        print(f"  Line {line_num}: '{target}' -> not found in repository")
                    if len(file_violations) > 5:
                        print(f"  ... and {len(file_violations) - 5} more.")

    print(f"\n[SUMMARY] Audited {total_files} documents across {len(AUDIT_DIRS)} target directories.")
    if total_violations == 0:
        print("[PASS] All sequential file references and document links resolved successfully.")
        return 0

    print(f"[ERROR] Detected {total_violations} broken sequential reference(s).")
    return 1


if __name__ == "__main__":
    sys.exit(main())
