#!/usr/bin/env python3
"""Rewire spec cross-links broken by the a2477488 stub purge.

Every file deleted in a2477488 was a "Compatibility reference retained for
existing cross-links" stub pointing at a canonical document. This script
rewrites links that targeted those stubs to point at the canonical targets
directly, then re-runs the cross-link linter to verify.
"""
import os
import re
import subprocess
import sys

REPO = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
os.chdir(REPO)

PURGE = "a2477488"


def git(*args):
    return subprocess.run(
        ["git", *args], capture_output=True, text=True, check=True
    ).stdout


def main():
    # 1. Map deleted stub -> canonical target (relative to repo root).
    deleted = git("diff", "--name-only", "--diff-filter=D", f"{PURGE}~1", PURGE, "--", "02-spec/").split()
    redirect = {}
    for stub in deleted:
        try:
            content = git("show", f"{PURGE}~1:{stub}")
        except subprocess.CalledProcessError:
            continue
        m = re.search(r"\[Open (?:canonical overview|current document)\]\(([^)]+)\)", content)
        if not m:
            print(f"WARN: no redirect found in stub {stub}")
            continue
        target = m.group(1)
        stub_dir = os.path.dirname(stub)
        canonical = os.path.normpath(os.path.join(stub_dir, target))
        redirect[stub] = canonical

    print(f"Mapped {len(redirect)} deleted stubs to canonical targets")

    # 2. Special case: root changelog is CHANGELOG.md (uppercase).
    # 3. For each broken ref from the linter, rewrite the link.
    out = subprocess.run(
        [sys.executable, "linter-scripts/check-spec-cross-links.py"],
        capture_output=True, text=True,
    ).stdout

    refs = []
    cur_source = None
    cur_line = None
    cur_target = None
    for line in out.splitlines():
        m = re.match(r"\s*\[missing-file\]\s+(\S+):(\d+)", line)
        if m:
            cur_source, cur_line = m.group(1), int(m.group(2))
            continue
        m = re.match(r"\s*target:\s+(\S+)", line)
        if m and cur_source:
            refs.append((cur_source, cur_line, m.group(1)))
            cur_source = None

    print(f"Found {len(refs)} broken refs")

    fixed = 0
    skipped = []
    # Group by source file for efficient editing.
    by_file = {}
    for source, line, target in refs:
        by_file.setdefault(source, []).append((line, target))

    for source, items in sorted(by_file.items()):
        with open(source, encoding="utf-8") as f:
            lines = f.readlines()
        changed = False
        for line_no, target in items:
            idx = line_no - 1
            if idx < 0 or idx >= len(lines):
                skipped.append((source, line_no, target, "line out of range"))
                continue
            line = lines[idx]
            # Resolve the link target to a repo-relative path.
            if target == "../../../changelog.md":
                new_target_rel = "CHANGELOG.md"
                new_link = target.replace("changelog.md", "CHANGELOG.md")
            else:
                src_dir = os.path.dirname(source)
                resolved = os.path.normpath(os.path.join(src_dir, target))
                if resolved not in redirect:
                    skipped.append((source, line_no, target, "not a purged stub"))
                    continue
                canonical = redirect[resolved]
                new_link = os.path.relpath(canonical, src_dir).replace(os.sep, "/")
                if not new_link.startswith("."):
                    new_link = "./" + new_link
            # Replace only the exact link target occurrences on this line.
            # Match markdown link destinations: ](target) possibly with #anchor.
            pattern = re.compile(r"\(" + re.escape(target) + r"(#?[^)]*)\)")
            new_line, n = pattern.subn(lambda m: "(" + new_link + m.group(1) + ")", line)
            if n == 0:
                # Try reference-style or bare occurrence.
                if target in line:
                    new_line = line.replace(target, new_link)
                    n = 1
            if n > 0:
                lines[idx] = new_line
                changed = True
                fixed += 1
            else:
                skipped.append((source, line_no, target, "target text not on line"))
        if changed:
            with open(source, "w", encoding="utf-8") as f:
                f.writelines(lines)

    print(f"Rewired {fixed} links")
    for s in skipped:
        print(f"SKIP: {s}")

    # 4. Re-run the linter.
    out = subprocess.run(
        [sys.executable, "linter-scripts/check-spec-cross-links.py"],
        capture_output=True, text=True,
    ).stdout
    remaining = out.count("[missing-file]")
    print(f"Remaining broken refs: {remaining}")
    return 0 if remaining == 0 else 1


if __name__ == "__main__":
    sys.exit(main())
