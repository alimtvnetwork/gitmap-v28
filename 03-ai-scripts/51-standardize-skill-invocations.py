#!/usr/bin/env python3
"""Standardizes skill invocations across .agents/skills and .cursor/skills.

Ensures every skill declares its canonical invocation format:
  > **Invocation:** `@[skill-name](file;.agents/skills/skill-name)`
for Antigravity, and:
  > **Invocation:** `@[skill-name](file;.cursor/skills/skill-name)`
for Cursor.
"""

from __future__ import annotations

import os
import sys

REPO_ROOT = os.path.abspath(os.path.join(os.path.dirname(__file__), ".."))


def update_skill_invocation(file_path: str, skill_name: str, is_cursor: bool) -> bool:
    skill_dir = ".cursor/skills" if is_cursor else ".agents/skills"
    inv_line = f"> **Invocation:** `@[{skill_name}](file;{skill_dir}/{skill_name})`"

    with open(file_path, "r", encoding="utf-8", errors="ignore") as f:
        content = f.read()

    # If invocation already exists with exact format, skip
    if f"> **Invocation:** `@[{skill_name}](file;{skill_dir}/{skill_name})`" in content:
        return False

    # If invocation exists with bare format without backticks, update it
    bare_pattern = f"> **Invocation:** @[{skill_name}](file;{skill_dir}/{skill_name})"
    if bare_pattern in content:
        content = content.replace(bare_pattern, inv_line)
        with open(file_path, "w", encoding="utf-8") as f:
            f.write(content)
        return True

    if not content.startswith("---"):
        return False

    parts = content.split("---", 2)
    if len(parts) < 3:
        return False

    fm = parts[1]
    body = parts[2]
    lines = body.splitlines(keepends=True)

    # Check if body has > [!IMPORTANT] early (within first 35 lines)
    important_idx = -1
    for idx, line in enumerate(lines[:35]):
        if "> [!IMPORTANT]" in line:
            important_idx = idx
            break

    if important_idx != -1:
        insert_idx = important_idx + 1
        for j in range(important_idx + 1, min(len(lines), important_idx + 12)):
            line_str = lines[j]
            if (
                "> **" in line_str
                or "> Prompt Version:" in line_str
                or "> Synchronization:" in line_str
                or "> Runtime:" in line_str
                or "> Canonical Prompt Specification:" in line_str
                or "> Toolchain Discovery & Primacy:" in line_str
                or "> Invoke:" in line_str
            ):
                insert_idx = j + 1
            elif line_str.strip() == ">" or not line_str.startswith(">"):
                break
        lines.insert(insert_idx, f"{inv_line}  \n")
    else:
        header_idx = -1
        for idx, line in enumerate(lines):
            if line.startswith("# "):
                header_idx = idx
                break
        if header_idx != -1:
            lines.insert(header_idx + 1, f"\n{inv_line}\n")
        else:
            lines.insert(0, f"\n{inv_line}\n")

    new_content = "---" + fm + "---" + "".join(lines)
    with open(file_path, "w", encoding="utf-8") as f:
        f.write(new_content)
    return True


def run_standardization() -> int:
    targets = [
        (os.path.join(REPO_ROOT, ".agents", "skills"), False),
        (os.path.join(REPO_ROOT, ".cursor", "skills"), True),
    ]

    total_updated = 0
    total_scanned = 0

    for dir_path, is_cursor in targets:
        if not os.path.exists(dir_path):
            continue
        for entry in sorted(os.listdir(dir_path)):
            skill_dir = os.path.join(dir_path, entry)
            if not os.path.isdir(skill_dir):
                continue
            skill_file = os.path.join(skill_dir, "skill.md")
            if not os.path.exists(skill_file):
                continue
            total_scanned += 1
            if update_skill_invocation(skill_file, entry, is_cursor):
                total_updated += 1

    print(f"Scanned {total_scanned} skills across .agents and .cursor. Updated: {total_updated}")
    return 0


if __name__ == "__main__":
    sys.exit(run_standardization())
