#!/usr/bin/env python3
"""Cross-platform check for gofmt formatting."""
import os
import subprocess
import sys


def main():
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    if hasattr(sys.stderr, "reconfigure"):
        sys.stderr.reconfigure(encoding="utf-8")

    repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
    gitmap_dir = sys.argv[1] if len(sys.argv) > 1 else os.path.join(repo_root, "gitmap")

    try:
        res = subprocess.run(
            ["gofmt", "-l", "."],
            cwd=gitmap_dir,
            capture_output=True,
            text=True,
            encoding="utf-8",
            errors="replace"
        )
    except FileNotFoundError:
        print("ERROR: gofmt not on PATH", file=sys.stderr)
        sys.exit(2)

    unformatted = [line.strip() for line in (res.stdout or "").splitlines() if line.strip()]
    if unformatted:
        print("::error::The following .go files are not gofmt-clean:", file=sys.stderr)
        for f in unformatted:
            print(f"  {f}", file=sys.stderr)
        print("\nFix locally with:  cd gitmap && gofmt -w .", file=sys.stderr)
        sys.exit(1)

    print("✅ All .go files are gofmt-clean.")
    sys.exit(0)


if __name__ == "__main__":
    main()
