#!/usr/bin/env python3
"""Targeted pre-gate that audits changed Go files for fix-repo format cleanliness."""

import os
import re
import subprocess
import sys


def main():
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    if hasattr(sys.stderr, "reconfigure"):
        sys.stderr.reconfigure(encoding="utf-8")

    repo_root = os.path.abspath(os.path.join(os.path.dirname(__file__), "..", ".."))
    gitmap_dir = os.path.join(repo_root, "gitmap")

    base = os.environ.get("GH_PR_BASE_SHA", "").strip() or "HEAD~1"

    try:
        diff_cmd = ["git", "diff", "--name-only", "--diff-filter=AM", base, "--", "gitmap/**/*.go"]
        res = subprocess.run(diff_cmd, cwd=repo_root, capture_output=True, text=True, encoding="utf-8", errors="replace")
        changed_files = [line.strip() for line in (res.stdout or "").splitlines() if line.strip()]
    except Exception:
        changed_files = []

    if not changed_files:
        print(f"fix-repo audit: no .go files changed vs {base} — skipping.")
        sys.exit(0)

    token_regex = re.compile(r"[A-Za-z0-9._-]+-v[0-9]+")
    affected = []

    for rel_path in changed_files:
        full_path = os.path.join(repo_root, rel_path)
        if not os.path.isfile(full_path):
            continue

        try:
            with open(full_path, "r", encoding="utf-8", errors="replace") as f:
                content = f.read()
                if token_regex.search(content):
                    affected.append(os.path.relpath(full_path, gitmap_dir))
        except Exception:
            continue

    if not affected:
        print("fix-repo audit: no fix-repo-shaped tokens in changed .go files — skipping.")
        sys.exit(0)

    print(f"fix-repo audit: scanning {len(affected)} affected .go file(s)…")

    try:
        gofmt_res = subprocess.run(
            ["gofmt", "-l"] + affected,
            cwd=gitmap_dir,
            capture_output=True,
            text=True,
            encoding="utf-8",
            errors="replace",
        )
    except FileNotFoundError:
        print("ERROR: gofmt not on PATH", file=sys.stderr)
        sys.exit(2)

    dirty = [line.strip() for line in (gofmt_res.stdout or "").splitlines() if line.strip()]
    if dirty:
        print(
            "::error title=fix-repo gofmt regression::These .go files contain fix-repo-style version tokens AND are not gofmt-clean — the v4.8.0+ post-rewrite gofmt step did not run on them.",
            file=sys.stderr,
        )
        for f in dirty:
            print(f"  {f}", file=sys.stderr)

        print("\n----- gofmt -d (diff) -----", file=sys.stderr)
        diff_res = subprocess.run(
            ["gofmt", "-d"] + dirty,
            cwd=gitmap_dir,
            capture_output=True,
            text=True,
            encoding="utf-8",
            errors="replace",
        )
        print(diff_res.stdout or "", file=sys.stderr)
        print("---------------------------", file=sys.stderr)
        print(
            "\nRoot cause: the fix-repo invocation that produced these changes did NOT auto-run gofmt.\nFix locally with:  cd gitmap && gofmt -w .\n",
            file=sys.stderr,
        )
        sys.exit(1)

    print("✅ All fix-repo touched .go files are gofmt-clean.")
    sys.exit(0)


if __name__ == "__main__":
    main()
