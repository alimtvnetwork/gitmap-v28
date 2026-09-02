#!/usr/bin/env python3
"""scripts/build_stamp.py — pre-build provenance stamp for stale-checkout detection (cross-platform).

Usage:
  python scripts/build_stamp.py
  python scripts/build_stamp.py --strict
"""

from __future__ import annotations

import argparse
import hashlib
import re
import subprocess
import sys
from pathlib import Path

STAMP_SCRIPT_VERSION = "1.0.0"


def probe_git(repo_root: Path, args: list[str], strict: bool) -> str:
    try:
        out = subprocess.check_output(
            ["git", "-C", str(repo_root)] + args,
            stderr=subprocess.DEVNULL,
            text=True,
        ).strip()
        return out if out else "(unknown)"
    except Exception:
        if strict:
            print("build-stamp: git command failed in strict mode", file=sys.stderr)
            sys.exit(1)
        return "(unknown)"


def probe_constants_version(constants_file: Path) -> str:
    if not constants_file.is_file():
        return "(unknown — constants.go missing)"
    try:
        with open(constants_file, "r", encoding="utf-8") as f:
            for line in f:
                m = re.match(r'^const\s+Version\s*=\s*"([^"]+)"', line.strip())
                if m:
                    return m.group(1)
    except OSError:
        pass
    return "(unknown — pattern miss)"


def fingerprint_file(repo_root: Path, label: str, path: Path) -> str:
    rel = path.relative_to(repo_root).as_posix()
    if not path.is_file():
        return f"  {label:<22} (missing — {rel})"

    try:
        with open(path, "rb") as f:
            data = f.read()
        sha12 = hashlib.sha256(data).hexdigest()[:12]
        lines = len(data.splitlines())
        return f"  {label:<22} sha256:{sha12}  lines:{lines}  {rel}"
    except OSError:
        return f"  {label:<22} (unreadable — {rel})"


def detect_redecl_risk(repo_root: Path, updaterepo_file: Path, updatedebug_file: Path, strict: bool) -> str:
    if not updaterepo_file.is_file() or not updatedebug_file.is_file():
        return "  redecl-risk-check       skipped (one or both source files missing)"

    repo_has = 0
    debug_has = 0
    pattern = re.compile(r"^func\s+(fileExists|fileExistsLoose)\(")

    try:
        with open(updaterepo_file, "r", encoding="utf-8") as f:
            repo_has = sum(1 for line in f if pattern.match(line.strip()))
    except OSError:
        pass

    try:
        with open(updatedebug_file, "r", encoding="utf-8") as f:
            debug_has = sum(1 for line in f if pattern.match(line.strip()))
    except OSError:
        pass

    if repo_has > 0 and debug_has > 0:
        msg = (
            "  redecl-risk-check       ⚠ FAIL — fileExists/fileExistsLoose declared in both files\n"
            "                           (this checkout predates the v3.113.0 fsutil migration)\n"
            "                           expected fix: git pull origin main"
        )
        if strict:
            print("build-stamp: redeclaration risk detected in strict mode", file=sys.stderr)
            sys.exit(1)
        return msg
    return "  redecl-risk-check       ok (no local fileExists* in cmd/ — fsutil migration present)"


def main() -> int:
    parser = argparse.ArgumentParser(description="Pre-build provenance stamp.")
    parser.add_argument("--strict", action="store_true", help="Exit 1 if git or checks fail")
    args = parser.parse_args()

    repo_root = Path(__file__).resolve().parent.parent
    constants_file = repo_root / "gitmap" / "constants" / "constants.go"
    updaterepo_file = repo_root / "gitmap" / "cmd" / "updaterepo.go"
    updatedebug_file = repo_root / "gitmap" / "cmd" / "updatedebugwindows.go"

    commit = probe_git(repo_root, ["rev-parse", "HEAD"], args.strict)
    short = probe_git(repo_root, ["rev-parse", "--short=10", "HEAD"], args.strict)
    branch = probe_git(repo_root, ["rev-parse", "--abbrev-ref", "HEAD"], args.strict)
    describe = probe_git(repo_root, ["describe", "--tags", "--always", "--dirty"], args.strict)
    cdate = probe_git(repo_root, ["log", "-1", "--format=%cI"], args.strict)
    csubj = probe_git(repo_root, ["log", "-1", "--format=%s"], args.strict)
    decl_ver = probe_constants_version(constants_file)

    fp_const = fingerprint_file(repo_root, "constants.go", constants_file)
    fp_update = fingerprint_file(repo_root, "updaterepo.go", updaterepo_file)
    fp_debug = fingerprint_file(repo_root, "updatedebugwindows.go", updatedebug_file)
    redecl = detect_redecl_risk(repo_root, updaterepo_file, updatedebug_file, args.strict)

    print(f"=== gitmap build-stamp v{STAMP_SCRIPT_VERSION} ====================================")
    print("Provenance for stale-checkout detection. Compare these against the SHA")
    print("and version you expected to build — if they don't match, run")
    print("'git pull origin main' before debugging the build error.\n")
    print("git")
    print(f"  commit                  {commit}")
    print(f"  short                   {short}")
    print(f"  branch                  {branch}")
    print(f"  describe                {describe}")
    print(f"  commit-date             {cdate}")
    print(f"  commit-subject          {csubj}\n")
    print("source")
    print(f"  declared-version        {decl_ver}")
    print(fp_const)
    print(fp_update)
    print(fp_debug)
    print("\nguards")
    print(redecl)
    print("=====================================================================")
    return 0


if __name__ == "__main__":
    sys.exit(main())
