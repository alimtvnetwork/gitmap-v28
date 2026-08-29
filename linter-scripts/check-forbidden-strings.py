#!/usr/bin/env python3
"""
Forbidden Strings Linter
=========================
Generic, TOML-driven scanner that fails CI when forbidden patterns
appear outside of allowlisted paths.

Config:  linter-scripts/forbidden-strings.toml
Usage:   python3 linter-scripts/check-forbidden-strings.py
CI:      wired in .github/workflows/ci.yml (lint job).
"""

from __future__ import annotations

import fnmatch
import os
import re
import sys

try:
    import tomllib  # Python 3.11+
except ModuleNotFoundError:
    try:
        import tomli as tomllib  # type: ignore[no-redef]
    except ModuleNotFoundError as exc:
        import sys as _sys; print(f"Error: {exc}", file=_sys.stderr)
        sys.exit("Error: Python 3.11+ required (tomllib), or install 'tomli'.")

ALWAYS_EXCLUDE_DIRS = {".git", "node_modules", "dist", "build"}
CONFIG_PATH = os.path.join(os.path.dirname(__file__), "forbidden-strings.toml")


def load_config(path: str) -> list[dict]:
    with open(path, "rb") as fh:
        data = tomllib.load(fh)
    rules = data.get("rule", [])
    if not rules:
        print("⚠️  No [[rule]] entries in config — nothing to check.")
    return rules


def is_excluded_dir(dirpath: str, extra_excludes: list[str]) -> bool:
    parts = dirpath.replace("\\", "/").split("/")
    all_excludes = ALWAYS_EXCLUDE_DIRS | set(extra_excludes)
    return any(p in all_excludes for p in parts)


def is_allowlisted(rel_path: str, allowlist: list[str]) -> bool:
    for pattern in allowlist:
        if fnmatch.fnmatch(rel_path, pattern):
            return True
        if rel_path == pattern or rel_path.startswith(pattern + "/"):
            return True
    return False


def is_excluded_file(filename: str, exclude_files: list[str]) -> bool:
    return any(fnmatch.fnmatch(filename, pat) for pat in exclude_files)


def scan_rule(rule: dict, root: str) -> list[str]:
    rule_id = rule["id"]
    pattern = re.compile(rule["pattern"])
    exclude_dirs = rule.get("exclude_dirs", [])
    exclude_files = rule.get("exclude_files", [])
    allowlist = rule.get("allowlist", [])
    findings: list[str] = []

    for dirpath, dirnames, filenames in os.walk(root):
        rel_dir = os.path.relpath(dirpath, root)
        if rel_dir == ".":
            rel_dir = ""

        if is_excluded_dir(rel_dir, exclude_dirs):
            dirnames.clear()
            continue

        dirnames[:] = [
            d for d in dirnames
            if d not in ALWAYS_EXCLUDE_DIRS and d not in exclude_dirs
        ]

        for fname in filenames:
            if is_excluded_file(fname, exclude_files):
                continue

            full = os.path.join(dirpath, fname)
            rel = os.path.relpath(full, root)

            if is_allowlisted(rel, allowlist):
                continue

            try:
                with open(full, "r", encoding="utf-8", errors="ignore") as fh:
                    for lineno, line in enumerate(fh, 1):
                        match = pattern.search(line)
                        if not match:
                            continue
                        findings.append({
                            "path": rel,
                            "line": lineno,
                            "matched": match.group(0),
                            "content": line.rstrip(),
                        })
            except (OSError, UnicodeDecodeError) as exc:
                import sys as _sys; print(f"Error: {exc}", file=_sys.stderr)
                continue

    return findings


def main() -> int:
    import sys
    if hasattr(sys.stdout, 'reconfigure'):
        sys.stdout.reconfigure(encoding='utf-8')
    config = load_config(CONFIG_PATH)
    root = os.getcwd()
    total_failures = 0

    for rule in config:
        rule_id = rule.get("id", "UNKNOWN")
        desc = rule.get("description", "")
        fix_hint = rule.get("fix_hint", "")

        print(f"🔍 [{rule_id}] {desc}")
        findings = scan_rule(rule, root)

        if findings:
            total_failures += len(findings)
            print(f"❌ FAIL [{rule_id}]: {len(findings)} forbidden phrase(s) blocked:")
            print("")
            for f in findings:
                # GitHub Actions error annotation — surfaces in PR checks UI.
                print(
                    f"::error file={f['path']},line={f['line']},title=Forbidden phrase blocked [{rule_id}]"
                    f"::[{rule_id}] matched forbidden phrase \"{f['matched']}\" "
                    f"at {f['path']}:{f['line']}"
                )
                # Human-readable line for local terminal use.
                print(f"  ⛔ [{rule_id}] {f['path']}:{f['line']}")
                print(f"     matched phrase : \"{f['matched']}\"")
                print(f"     line content   : {f['content']}")
                print("")
            if fix_hint:
                print(f"   Fix hint: {fix_hint}")
            print("")
        else:
            print(f"✅ PASS: no forbidden strings found.")
        print("")

    if total_failures:
        print(f"💥 {total_failures} total finding(s) across all rules.")
        return 1

    print("🎉 All forbidden-string rules passed.")
    return 0


if __name__ == "__main__":
    sys.exit(main())
