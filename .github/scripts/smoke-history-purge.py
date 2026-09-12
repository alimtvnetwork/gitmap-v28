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
    """Executes a git command in the target directory and returns stripped stdout."""
    return subprocess.check_output(
        ["git"] + list(args),
        cwd=str(cwd),
        text=True,
        stderr=subprocess.DEVNULL,
    ).strip()


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
    """Executes Scenario A: removes secret.env from all history."""
    readme_lines: list[str] = []
    for n in range(1, 6):
        readme_lines.append(f"line {n}\n")
        files = {"README.md": "".join(readme_lines)}
        if n <= 3:
            files["secret.env"] = f"API_KEY=leaked-{n}\n"
        commit_files(repo, files, f"commit {n}")

    push_main(repo)

    pre = git_cmd(repo, "log", "--all", "--oneline", "--", "secret.env").splitlines()
    if len(pre) != 3:
        print(f"Pre-check failed: expected 3 commits with secret.env, got {len(pre)}", file=sys.stderr)
        return None

    res = subprocess.run(
        [str(gitmap_bin), "history-purge", "secret.env", "--no-push", "--keep-sandbox"],
        cwd=str(repo),
        capture_output=True,
        text=True,
    )
    sys.stderr.write(res.stdout + res.stderr)

    sandbox = parse_sandbox(res.stdout + res.stderr)
    if not sandbox:
        return None

    remaining = git_cmd(sandbox, "log", "--all", "--oneline", "--", "secret.env").splitlines()
    if remaining:
        print(f"FAIL: history-purge left {len(remaining)} commits referencing secret.env", file=sys.stderr)
        return None

    readme_commits = git_cmd(sandbox, "log", "--all", "--oneline", "--", "README.md").splitlines()
    if len(readme_commits) < 5:
        print(f"FAIL: README commits collapsed to {len(readme_commits)} (expected >=5)", file=sys.stderr)
        return None

    print("PASS: history-purge removed secret.env from all history; README.md untouched")

    return sandbox


def run_scenario_b(gitmap_bin: Path, repo: Path, prev_sandbox: Path) -> int:
    """Executes Scenario B: --message scoping verification."""
    shutil.rmtree(prev_sandbox, ignore_errors=True)
    res_b = subprocess.run(
        [str(gitmap_bin), "history-purge", "secret.env", "--no-push", "--keep-sandbox", "--message", "scrubbed by ci"],
        cwd=str(repo),
        capture_output=True,
        text=True,
    )
    sys.stderr.write(res_b.stdout + res_b.stderr)

    sandbox_b = parse_sandbox(res_b.stdout + res_b.stderr, "scenario B output")
    if not sandbox_b:
        return 1

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


def main() -> int:
    if len(sys.argv) < 2:
        print("path to gitmap binary required", file=sys.stderr)
        return 1

    gitmap_bin = Path(sys.argv[1]).resolve()
    if not gitmap_bin.is_file():
        print(f"gitmap binary not found at {gitmap_bin}", file=sys.stderr)
        return 1

    test_base = os.path.join(tempfile.gettempdir(), "gitmap", "test")
    os.makedirs(test_base, exist_ok=True)
    with tempfile.TemporaryDirectory(dir=test_base) as work_dir:
        _, repo = init_test_repos(Path(work_dir))

        sandbox = run_scenario_a(gitmap_bin, repo)
        if not sandbox:
            return 1

        return run_scenario_b(gitmap_bin, repo, sandbox)


if __name__ == "__main__":
    sys.exit(main())
