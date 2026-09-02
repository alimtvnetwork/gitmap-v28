#!/usr/bin/env python3
"""smoke-history-purge.py — Smoke test: gitmap history-purge removes a file from all history (cross-platform).

Usage:
  python .github/scripts/smoke-history-purge.py <path-to-gitmap-binary>

Spec: spec/04-generic-cli/16-history-rewrite.md
"""

from __future__ import annotations

import re
import shutil
import subprocess
import sys
import tempfile
from pathlib import Path


def git_cmd(cwd: Path, *args: str) -> str:
    return subprocess.check_output(["git"] + list(args), cwd=str(cwd), text=True, stderr=subprocess.DEVNULL).strip()


def main() -> int:
    if len(sys.argv) < 2:
        print("path to gitmap binary required", file=sys.stderr)
        return 1

    gitmap_bin = Path(sys.argv[1]).resolve()
    if not gitmap_bin.is_file():
        print(f"gitmap binary not found at {gitmap_bin}", file=sys.stderr)
        return 1

    with tempfile.TemporaryDirectory() as work_dir:
        work = Path(work_dir)
        origin = work / "origin.git"
        origin.mkdir()
        subprocess.check_call(["git", "init", "--bare", str(origin)], stdout=subprocess.DEVNULL)

        repo = work / "work"
        subprocess.check_call(["git", "clone", str(origin), str(repo)], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)

        # 5 commits with secret.env in first 3
        readme = repo / "README.md"
        secret = repo / "secret.env"
        for n in range(1, 6):
            with open(readme, "a", encoding="utf-8") as f:
                f.write(f"line {n}\n")
            if n <= 3:
                secret.write_text(f"API_KEY=leaked-{n}\n", encoding="utf-8")
                subprocess.check_call(["git", "add", "secret.env"], cwd=str(repo))
            subprocess.check_call(["git", "add", "README.md"], cwd=str(repo))
            subprocess.check_call(
                ["git", "-c", "user.name=ci", "-c", "user.email=ci@x", "commit", "-m", f"commit {n}"],
                cwd=str(repo),
                stdout=subprocess.DEVNULL,
            )
        subprocess.check_call(["git", "push", "origin", "main"], cwd=str(repo), stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)

        # Pre-check
        pre = git_cmd(repo, "log", "--all", "--oneline", "--", "secret.env").splitlines()
        if len(pre) != 3:
            print(f"Pre-check failed: expected 3 commits with secret.env, got {len(pre)}", file=sys.stderr)
            return 1

        # Run history-purge
        res = subprocess.run(
            [str(gitmap_bin), "history-purge", "secret.env", "--no-push", "--keep-sandbox"],
            cwd=str(repo),
            capture_output=True,
            text=True,
        )
        sys.stderr.write(res.stdout + res.stderr)

        sandbox_match = re.search(r"sandbox kept at ([^\r\n]+)", res.stdout + res.stderr)
        if not sandbox_match:
            print("FAIL: could not parse sandbox path from output", file=sys.stderr)
            return 1

        sandbox = Path(sandbox_match.group(1).strip())
        if not sandbox.is_dir():
            print(f"FAIL: sandbox dir {sandbox} does not exist", file=sys.stderr)
            return 1

        # Verify secret.env gone from sandbox
        remaining = git_cmd(sandbox, "log", "--all", "--oneline", "--", "secret.env").splitlines()
        if remaining:
            print(f"FAIL: history-purge left {len(remaining)} commits referencing secret.env", file=sys.stderr)
            return 1

        # Verify README survived
        readme_commits = git_cmd(sandbox, "log", "--all", "--oneline", "--", "README.md").splitlines()
        if len(readme_commits) < 5:
            print(f"FAIL: README commits collapsed to {len(readme_commits)} (expected >=5)", file=sys.stderr)
            return 1

        print("PASS: history-purge removed secret.env from all history; README.md untouched")

        # Scenario B: --message scoping
        shutil.rmtree(sandbox, ignore_errors=True)
        res_b = subprocess.run(
            [str(gitmap_bin), "history-purge", "secret.env", "--no-push", "--keep-sandbox", "--message", "scrubbed by ci"],
            cwd=str(repo),
            capture_output=True,
            text=True,
        )
        sys.stderr.write(res_b.stdout + res_b.stderr)

        sb_b_match = re.search(r"sandbox kept at ([^\r\n]+)", res_b.stdout + res_b.stderr)
        if not sb_b_match:
            print("FAIL: could not parse sandbox path from scenario B output", file=sys.stderr)
            return 1

        sandbox_b = Path(sb_b_match.group(1).strip())
        messages = git_cmd(sandbox_b, "log", "--all", "--pretty=format:%s").splitlines()
        scrubbed = messages.count("scrubbed by ci")
        untouched = sum(1 for m in messages if m in ("commit 4", "commit 5"))

        if scrubbed < 1:
            print(f"FAIL: --message did not rewrite any touched commit (got {scrubbed})", file=sys.stderr)
            return 1
        if untouched != 2:
            print(f"FAIL: --message leaked into untouched commits (expected 2 originals, got {untouched})", file=sys.stderr)
            return 1

        print(f"PASS: --message scoped to touched commits only ({scrubbed} rewritten, {untouched} originals kept)")

    return 0


if __name__ == "__main__":
    sys.exit(main())
