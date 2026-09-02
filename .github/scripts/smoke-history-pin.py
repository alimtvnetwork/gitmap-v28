#!/usr/bin/env python3
"""smoke-history-pin.py — Smoke test: gitmap history-pin collapses revisions (cross-platform).

Usage:
  python .github/scripts/smoke-history-pin.py <path-to-gitmap-binary>

Spec: spec/04-generic-cli/16-history-rewrite.md
"""

from __future__ import annotations

import hashlib
import re
import shutil
import subprocess
import sys
import tempfile
from pathlib import Path


def git_cmd(cwd: Path, *args: str) -> str:
    return subprocess.check_output(["git"] + list(args), cwd=str(cwd), text=True, stderr=subprocess.DEVNULL).strip()


def sha256_text(text: str) -> str:
    return hashlib.sha256(text.encode("utf-8")).hexdigest()


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

        # 3 commits with distinct content in X
        for n in (1, 2, 3):
            (repo / "X").write_text(f"version {n} contents of X\n", encoding="utf-8")
            subprocess.check_call(["git", "add", "X"], cwd=str(repo))
            subprocess.check_call(
                ["git", "-c", "user.name=ci", "-c", "user.email=ci@x", "commit", "-m", f"X v{n}"],
                cwd=str(repo),
                stdout=subprocess.DEVNULL,
            )
        subprocess.check_call(["git", "push", "origin", "main"], cwd=str(repo), stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)

        # Check 3 distinct hashes pre-rewrite
        shas = git_cmd(repo, "log", "--all", "--pretty=format:%H", "--", "X").splitlines()
        contents = set()
        for s in shas:
            contents.add(git_cmd(repo, "show", f"{s}:X"))
        if len(contents) != 3:
            print(f"Pre-check failed: expected 3 distinct contents, got {len(contents)}", file=sys.stderr)
            return 1

        # Run history-pin
        res = subprocess.run(
            [str(gitmap_bin), "history-pin", "X", "--no-push", "--keep-sandbox"],
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

        # Check that X in sandbox hashes identically across all revisions
        sb_shas = git_cmd(sandbox, "log", "--all", "--pretty=format:%H", "--", "X").splitlines()
        sb_contents = {git_cmd(sandbox, "show", f"{s}:X") for s in sb_shas}
        if len(sb_contents) != 1:
            print(f"FAIL: history-pin left {len(sb_contents)} distinct content hashes (expected 1)", file=sys.stderr)
            return 1

        print("PASS: history-pin collapsed X to a single content hash across all history")

        # Scenario B: multi-path pin + --message scoping
        for n in (1, 2):
            (repo / "Y").write_text(f"Y v{n}\n", encoding="utf-8")
            (repo / "Z").write_text(f"Z v{n}\n", encoding="utf-8")
            subprocess.check_call(["git", "add", "Y", "Z"], cwd=str(repo))
            subprocess.check_call(
                ["git", "-c", "user.name=ci", "-c", "user.email=ci@x", "commit", "-m", f"Y/Z v{n}"],
                cwd=str(repo),
                stdout=subprocess.DEVNULL,
            )

        (repo / "Z").write_text("Z final\n", encoding="utf-8")
        subprocess.check_call(["git", "add", "Z"], cwd=str(repo))
        subprocess.check_call(
            ["git", "-c", "user.name=ci", "-c", "user.email=ci@x", "commit", "-m", "Z only"],
            cwd=str(repo),
            stdout=subprocess.DEVNULL,
        )
        subprocess.check_call(["git", "push", "origin", "main"], cwd=str(repo), stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)

        shutil.rmtree(sandbox, ignore_errors=True)

        res_b = subprocess.run(
            [str(gitmap_bin), "history-pin", "X", "Y", "--no-push", "--keep-sandbox", "--message", "pinned by ci"],
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
        for p in ("X", "Y"):
            p_shas = git_cmd(sandbox_b, "log", "--all", "--pretty=format:%H", "--", p).splitlines()
            p_contents = {git_cmd(sandbox_b, "show", f"{s}:{p}") for s in p_shas}
            if len(p_contents) != 1:
                print(f"FAIL: multi-path pin: {p} has {len(p_contents)} distinct hashes (expected 1)", file=sys.stderr)
                return 1

        z_shas = git_cmd(sandbox_b, "log", "--all", "--pretty=format:%H", "--", "Z").splitlines()
        z_contents = {git_cmd(sandbox_b, "show", f"{s}:Z") for s in z_shas}
        if len(z_contents) < 2:
            print(f"FAIL: history-pin leaked into unrelated path Z (collapsed to {len(z_contents)})", file=sys.stderr)
            return 1

        zonly = git_cmd(sandbox_b, "log", "--all", "--pretty=format:%s").splitlines()
        if zonly.count("Z only") != 1:
            print("FAIL: --message leaked into untouched 'Z only' commit", file=sys.stderr)
            return 1

        print("PASS: multi-path pin (X, Y) + --message scoping; Z untouched")

    return 0


if __name__ == "__main__":
    sys.exit(main())
