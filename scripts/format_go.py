#!/usr/bin/env python3
"""scripts/format_go.py — Auto-format Go files with gofmt (cross-platform).

Modes:
  python scripts/format_go.py                # format every .go file under gitmap/
  python scripts/format_go.py --staged       # format only staged .go files (hook mode)
  python scripts/format_go.py path/a.go ...  # format an explicit file list

Exit codes:
  0 — clean or safely re-staged
  1 — partial-staging conflict or missing gofmt
  2 — usage error
"""

from __future__ import annotations

import shutil
import subprocess
import sys
from pathlib import Path


def main() -> int:
    repo_root = Path(__file__).resolve().parent.parent
    gofmt_exe = shutil.which("gofmt")
    if not gofmt_exe:
        print("⚠ gofmt not found — install Go toolchain (https://go.dev/dl/)", file=sys.stderr)
        if "--staged" in sys.argv:
            return 1
        return 0

    mode = "all"
    explicit_files: list[str] = []

    if len(sys.argv) > 1:
        arg = sys.argv[1]
        if arg == "--staged":
            mode = "staged"
        elif arg in ("-h", "--help"):
            print(__doc__)
            return 0
        elif arg.startswith("-"):
            print(f"✗ unknown flag: {arg}", file=sys.stderr)
            return 2
        else:
            mode = "explicit"
            explicit_files = sys.argv[1:]

    files: list[Path] = []
    if mode == "staged":
        try:
            out = subprocess.check_output(
                ["git", "diff", "--cached", "--name-only", "--diff-filter=ACM", "--", "*.go"],
                cwd=str(repo_root),
                text=True,
                stderr=subprocess.DEVNULL,
            )
            files = [repo_root / line.strip() for line in out.splitlines() if line.strip()]
        except subprocess.SubprocessError:
            files = []
    elif mode == "all":
        gitmap_dir = repo_root / "gitmap"
        if gitmap_dir.is_dir():
            files = list(gitmap_dir.rglob("*.go"))
    else:
        files = [Path(f) for f in explicit_files]

    files = [f for f in files if f.is_file()]
    if not files:
        print("  (no .go files to format)")
        return 0

    # Find dirty files
    dirty: list[Path] = []
    for f in files:
        try:
            out = subprocess.check_output([gofmt_exe, "-l", str(f)], text=True, stderr=subprocess.DEVNULL).strip()
            if out:
                dirty.append(f)
        except subprocess.SubprocessError:
            pass

    if not dirty:
        print("  ✓ all Go files already gofmt-clean")
        return 0

    print(f"Formatting {len(dirty)} file(s):")
    for f in dirty:
        print(f"    {f.relative_to(repo_root) if f.is_relative_to(repo_root) else f}")
        subprocess.run([gofmt_exe, "-w", str(f)], check=False)

    if mode == "staged":
        conflicts: list[str] = []
        for f in dirty:
            rel = str(f.relative_to(repo_root))
            try:
                staged_blob = subprocess.check_output(["git", "show", f":{rel}"], cwd=str(repo_root))
                # run gofmt on staged content
                p = subprocess.Popen([gofmt_exe], stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.DEVNULL)
                staged_fmt, _ = p.communicate(input=staged_blob)
                with open(f, "rb") as fh:
                    worktree_content = fh.read()
                if staged_fmt == worktree_content:
                    subprocess.run(["git", "add", "--", rel], cwd=str(repo_root), check=False)
                else:
                    conflicts.append(rel)
            except Exception:
                conflicts.append(rel)

        if conflicts:
            print("\n✗ Partial-staging detected — these files have unstaged edits and", file=sys.stderr)
            print("  were formatted but NOT re-staged:", file=sys.stderr)
            for c in conflicts:
                print(f"    {c}", file=sys.stderr)
            print("\n  Resolve with:  git add <file>   (commits the formatted version)", file=sys.stderr)
            return 1
        print("  ✓ formatted files re-staged")

    return 0


if __name__ == "__main__":
    sys.exit(main())
