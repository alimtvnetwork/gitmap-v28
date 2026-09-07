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
    """Executes a git command in the target directory and returns stripped stdout."""
    return subprocess.check_output(
        ["git"] + list(args),
        cwd=str(cwd),
        text=True,
        stderr=subprocess.DEVNULL,
    ).strip()


def sha256_text(text: str) -> str:
    """Returns SHA256 hex digest for UTF-8 text."""
    return hashlib.sha256(text.encode("utf-8")).hexdigest()


def init_test_repos(work: Path) -> tuple[Path, Path]:
    """Initializes a bare origin repository and clones a local working repo."""
    origin = work / "origin.git"
    origin.mkdir()
    subprocess.check_call(["git", "init", "--bare", str(origin)], stdout=subprocess.DEVNULL)

    repo = work / "work"
    subprocess.check_call(
        ["git", "clone", str(origin), str(repo)],
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
    )

    return origin, repo


def commit_files(repo: Path, files: dict[str, str], message: str) -> None:
    """Writes files to repo, stages them, and commits with standard CI credentials."""
    for rel_path, content in files.items():
        target = repo / rel_path
        target.write_text(content, encoding="utf-8")
        subprocess.check_call(["git", "add", rel_path], cwd=str(repo))

    subprocess.check_call(
        ["git", "-c", "user.name=ci", "-c", "user.email=ci@x", "commit", "-m", message],
        cwd=str(repo),
        stdout=subprocess.DEVNULL,
    )


def push_main(repo: Path) -> None:
    """Pushes local main branch to origin."""
    subprocess.check_call(
        ["git", "push", "origin", "main"],
        cwd=str(repo),
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
    )


def get_path_contents(repo: Path, path: str) -> set[str]:
    """Returns the set of distinct blob contents for a given path across all commits."""
    shas = git_cmd(repo, "log", "--all", "--pretty=format:%H", "--", path).splitlines()

    return {git_cmd(repo, "show", f"{s}:{path}") for s in shas if s.strip()}


def parse_sandbox(output: str, label: str = "output") -> Path | None:
    """Parses and validates the sandbox directory path emitted by gitmap."""
    match = re.search(r"sandbox kept at ([^\r\n]+)", output)
    if not match:
        print(f"FAIL: could not parse sandbox path from {label}", file=sys.stderr)
        return None

    sandbox = Path(match.group(1).strip())
    if not sandbox.is_dir():
        print(f"FAIL: sandbox dir {sandbox} does not exist", file=sys.stderr)
        return None

    return sandbox


def run_scenario_a(gitmap_bin: Path, repo: Path) -> Path | None:
    """Executes Scenario A: single path pinning."""
    for n in (1, 2, 3):
        commit_files(repo, {"X": f"version {n} contents of X\n"}, f"X v{n}")

    push_main(repo)

    contents = get_path_contents(repo, "X")
    if len(contents) != 3:
        print(f"Pre-check failed: expected 3 distinct contents, got {len(contents)}", file=sys.stderr)
        return None

    res = subprocess.run(
        [str(gitmap_bin), "history-pin", "X", "--no-push", "--keep-sandbox"],
        cwd=str(repo),
        capture_output=True,
        text=True,
    )
    sys.stderr.write(res.stdout + res.stderr)

    sandbox = parse_sandbox(res.stdout + res.stderr)
    if not sandbox:
        return None

    sb_contents = get_path_contents(sandbox, "X")
    if len(sb_contents) != 1:
        print(f"FAIL: history-pin left {len(sb_contents)} distinct content hashes (expected 1)", file=sys.stderr)
        return None

    print("PASS: history-pin collapsed X to a single content hash across all history")

    return sandbox


def run_scenario_b(gitmap_bin: Path, repo: Path, prev_sandbox: Path) -> int:
    """Executes Scenario B: multi-path pin + --message scoping."""
    for n in (1, 2):
        commit_files(repo, {"Y": f"Y v{n}\n", "Z": f"Z v{n}\n"}, f"Y/Z v{n}")

    commit_files(repo, {"Z": "Z final\n"}, "Z only")
    push_main(repo)

    shutil.rmtree(prev_sandbox, ignore_errors=True)

    res_b = subprocess.run(
        [str(gitmap_bin), "history-pin", "X", "Y", "--no-push", "--keep-sandbox", "--message", "pinned by ci"],
        cwd=str(repo),
        capture_output=True,
        text=True,
    )
    sys.stderr.write(res_b.stdout + res_b.stderr)

    sandbox_b = parse_sandbox(res_b.stdout + res_b.stderr, "scenario B output")
    if not sandbox_b:
        return 1

    for p in ("X", "Y"):
        p_contents = get_path_contents(sandbox_b, p)
        if len(p_contents) != 1:
            print(f"FAIL: multi-path pin: {p} has {len(p_contents)} distinct hashes (expected 1)", file=sys.stderr)
            return 1

    z_contents = get_path_contents(sandbox_b, "Z")
    if len(z_contents) < 2:
        print(f"FAIL: history-pin leaked into unrelated path Z (collapsed to {len(z_contents)})", file=sys.stderr)
        return 1

    zonly = git_cmd(sandbox_b, "log", "--all", "--pretty=format:%s").splitlines()
    if zonly.count("Z only") != 1:
        print("FAIL: --message leaked into untouched 'Z only' commit", file=sys.stderr)
        return 1

    print("PASS: multi-path pin (X, Y) + --message scoping; Z untouched")

    return 0


def main() -> int:
    if len(sys.argv) < 2:
        print("path to gitmap binary required", file=sys.stderr)
        return 1

    gitmap_bin = Path(sys.argv[1]).resolve()
    if not gitmap_bin.is_file():
        print(f"gitmap binary not found at {gitmap_bin}", file=sys.stderr)
        return 1

    with tempfile.TemporaryDirectory() as work_dir:
        _, repo = init_test_repos(Path(work_dir))

        sandbox = run_scenario_a(gitmap_bin, repo)
        if not sandbox:
            return 1

        return run_scenario_b(gitmap_bin, repo, sandbox)


if __name__ == "__main__":
    sys.exit(main())
