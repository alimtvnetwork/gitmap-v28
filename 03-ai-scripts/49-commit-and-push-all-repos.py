#!/usr/bin/env python3
"""
03-ai-scripts/49-commit-and-push-all-repos.py
=============================================
Autonomous Multi-Repository Commit & Push Orchestrator.

Discovers all Git repositories across a target workspace directory (e.g. D:\\work,
~/git-work, or custom root), audits working tree status and divergence, stages
uncommitted changes, creates atomic conventional commits, and safely pushes to
upstream remote tracking branches.

Capabilities:
  1. Workspace Root Auto-Discovery:
     - Automatically locates workspace directory across Windows (D:\\work, C:\\work),
       POSIX/Linux (~/git-work, ~/work), environment variables (WORK_DIR), or custom --dir flag.
  2. Multi-Repository Topology Discovery:
     - Recursively discovers genuine top-level Git repositories.
     - Automatically skips build artifacts, dependency directories (node_modules, target, vendor, dist, .cache).
  3. Working Tree & Divergence Audit:
     - Detects dirty files, untracked changes, ahead/behind commits, and detached HEADs.
  4. Atomic Commit Automation:
     - Stages modifications cleanly (git add -A).
     - Generates meaningful conventional commits (chore(sync): ...) or accepts custom messages.
  5. Safe Remote Push & Upstream Binding:
     - Automatically configures upstream tracking (`git push -u <remote> <branch>`) when missing.
     - Gracefully reports permission or network issues without aborting the entire fleet.
  6. Dry-Run & JSON Reporting:
     - Full simulation via `--dry-run` and machine-readable output via `--json`.

Usage:
  # 1. Audit status across workspace (read-only):
  python 03-ai-scripts/49-commit-and-push-all-repos.py --check

  # 2. Commit and push all dirty and ahead repositories in default workspace:
  python 03-ai-scripts/49-commit-and-push-all-repos.py

  # 3. Target specific workspace directory (e.g. D:\\work or ~/git-work):
  python 03-ai-scripts/49-commit-and-push-all-repos.py --dir "d:/work"

  # 4. Preview actions without mutating disk or remotes:
  python 03-ai-scripts/49-commit-and-push-all-repos.py --dry-run

  # 5. Push unpushed commits only (skip committing dirty trees):
  python 03-ai-scripts/49-commit-and-push-all-repos.py --push-only

  # 6. Run internal self-test suite:
  python 03-ai-scripts/49-commit-and-push-all-repos.py --self-test
"""
from __future__ import annotations

import argparse
from dataclasses import asdict, dataclass
from enum import Enum
import json
import os
from pathlib import Path
import subprocess
import sys
import tempfile
import time
from typing import Any, Dict, List, Optional, Tuple

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")
if hasattr(sys.stderr, "reconfigure"):
    sys.stderr.reconfigure(encoding="utf-8")


class RepoStatusType(str, Enum):
    """Categorized status of a Git repository."""
    CLEAN = "clean"
    DIRTY = "dirty"
    AHEAD = "ahead"
    BEHIND = "behind"
    DIVERGED = "diverged"
    DETACHED = "detached"
    NO_COMMITS = "no_commits"
    ERROR = "error"


class OperationStatusType(str, Enum):
    """Result status of commit/push operation."""
    SUCCESS = "success"
    SKIPPED = "skipped"
    SIMULATED = "simulated"
    FAILED = "failed"


SKIP_PARENT_PATTERNS = {
    "node_modules",
    "target",
    "dist",
    "build",
    "vendor",
    ".cache",
    ".tmp",
    "tmp",
    "__pycache__",
    ".venv",
    "venv",
    ".gitmap",
}


@dataclass
class RepoAuditResult:
    """Detailed audit metrics for a single repository."""
    name: str
    relative_path: str
    absolute_path: str
    branch: str
    primary_remote: str
    status: RepoStatusType
    is_dirty: bool
    dirty_file_count: int
    untracked_count: int
    ahead_count: int
    behind_count: int
    has_upstream: bool
    is_detached: bool


@dataclass
class RepoExecutionResult:
    """Execution outcome for a single repository."""
    name: str
    relative_path: str
    branch: str
    commit_status: OperationStatusType
    push_status: OperationStatusType
    commit_message: Optional[str] = None
    commit_sha: Optional[str] = None
    error_message: Optional[str] = None


def run_cmd(args: List[str], cwd: Path, timeout_seconds: int = 45) -> Tuple[int, str, str]:
    """Execute a shell command with utf-8 decoding and timeout protection."""
    try:
        proc = subprocess.run(
            args,
            cwd=str(cwd),
            capture_output=True,
            text=True,
            encoding="utf-8",
            errors="replace",
            timeout=timeout_seconds,
        )
        return proc.returncode, proc.stdout.strip(), proc.stderr.strip()
    except subprocess.TimeoutExpired:
        return 124, "", f"Command timed out after {timeout_seconds}s: {' '.join(args)}"
    except Exception as exc:
        return 1, "", f"Execution error: {exc}"


def detect_workspace_root(explicit_dir: Optional[str] = None) -> Path:
    """Autonomously detect target workspace root directory across operating systems."""
    if explicit_dir and explicit_dir.strip():
        target = Path(explicit_dir.strip()).expanduser().resolve()
        if target.is_dir():
            return target
        raise ValueError(f"Specified workspace directory does not exist: {explicit_dir}")

    env_dir = os.environ.get("WORK_DIR", "").strip()
    if env_dir:
        p_env = Path(env_dir).expanduser().resolve()
        if p_env.is_dir():
            return p_env

    # 1. Check parent directory if named work or git-work
    cwd = Path.cwd().resolve()
    for candidate in [cwd, *cwd.parents]:
        if candidate.name.lower() in ("work", "git-work"):
            return candidate

    # 2. Check standard Windows drive roots
    for drive in ("D:", "C:", "E:"):
        for sub in ("work", "git-work"):
            candidate = Path(f"{drive}{os.sep}{sub}")
            if candidate.is_dir():
                return candidate

    # 3. Check home directory paths
    home = Path.home()
    for candidate_name in ("git-work", "work"):
        home_work = home / candidate_name
        if home_work.is_dir():
            return home_work

    # Default fallback to parent of current working directory
    return cwd.parent


def discover_repositories(root: Path) -> List[Path]:
    """Discover all authentic Git repositories directly or nested under root directory."""
    discovered: List[Path] = []
    root_resolved = root.resolve()

    # If root itself is a git repository
    if (root_resolved / ".git").is_dir():
        discovered.append(root_resolved)

    for git_dir in root_resolved.glob("**/.git"):
        repo_dir = git_dir.parent
        # Skip if within build / cache artifacts
        parts = repo_dir.relative_to(root_resolved).parts
        if any(skip_part in parts for skip_part in SKIP_PARENT_PATTERNS):
            continue
        if repo_dir not in discovered:
            discovered.append(repo_dir)

    discovered.sort(key=lambda p: str(p.relative_to(root_resolved)).lower())
    return discovered


def audit_repository(repo_path: Path, root: Path) -> RepoAuditResult:
    """Inspect and categorize the working tree and remote divergence state of a repository."""
    rel_path = str(repo_path.relative_to(root)) if repo_path != root else repo_path.name
    repo_name = repo_path.name

    # Check branch / HEAD
    code, branch_out, _ = run_cmd(["git", "rev-parse", "--abbrev-ref", "HEAD"], repo_path)
    branch = branch_out if code == 0 else "unknown"
    is_detached = (branch == "HEAD")

    # Check primary remote
    code, remotes_out, _ = run_cmd(["git", "remote"], repo_path)
    remotes = remotes_out.splitlines() if code == 0 and remotes_out else []
    primary_remote = "origin" if "origin" in remotes else (remotes[0] if remotes else "")

    # Check status (porcelain)
    code, status_out, _ = run_cmd(["git", "status", "--porcelain"], repo_path)
    status_lines = [line for line in status_out.splitlines() if line.strip()]
    is_dirty = len(status_lines) > 0
    dirty_file_count = len(status_lines)
    untracked_count = sum(1 for line in status_lines if line.startswith("??"))

    # Check commits count (handle empty repositories)
    code, rev_count_out, _ = run_cmd(["git", "rev-list", "-n", "1", "HEAD"], repo_path)
    has_commits = (code == 0 and len(rev_count_out.strip()) > 0)

    if not has_commits:
        return RepoAuditResult(
            name=repo_name,
            relative_path=rel_path,
            absolute_path=str(repo_path),
            branch=branch,
            primary_remote=primary_remote,
            status=RepoStatusType.NO_COMMITS,
            is_dirty=is_dirty,
            dirty_file_count=dirty_file_count,
            untracked_count=untracked_count,
            ahead_count=0,
            behind_count=0,
            has_upstream=False,
            is_detached=is_detached,
        )

    # Check upstream tracking branch and divergence
    code, upstream_out, _ = run_cmd(["git", "rev-parse", "--abbrev-ref", "@{u}"], repo_path)
    has_upstream = (code == 0 and len(upstream_out.strip()) > 0)

    ahead_count = 0
    behind_count = 0

    if has_upstream:
        code_a, ahead_out, _ = run_cmd(["git", "rev-list", "@{u}..HEAD", "--count"], repo_path)
        ahead_count = int(ahead_out) if code_a == 0 and ahead_out.isdigit() else 0

        code_b, behind_out, _ = run_cmd(["git", "rev-list", "HEAD..@{u}", "--count"], repo_path)
        behind_count = int(behind_out) if code_b == 0 and behind_out.isdigit() else 0

    # Determine overall status
    if is_detached:
        overall_status = RepoStatusType.DETACHED
    elif is_dirty:
        overall_status = RepoStatusType.DIRTY
    elif ahead_count > 0 and behind_count > 0:
        overall_status = RepoStatusType.DIVERGED
    elif ahead_count > 0:
        overall_status = RepoStatusType.AHEAD
    elif behind_count > 0:
        overall_status = RepoStatusType.BEHIND
    else:
        overall_status = RepoStatusType.CLEAN

    return RepoAuditResult(
        name=repo_name,
        relative_path=rel_path,
        absolute_path=str(repo_path),
        branch=branch,
        primary_remote=primary_remote,
        status=overall_status,
        is_dirty=is_dirty,
        dirty_file_count=dirty_file_count,
        untracked_count=untracked_count,
        ahead_count=ahead_count,
        behind_count=behind_count,
        has_upstream=has_upstream,
        is_detached=is_detached,
    )


def execute_commit_and_push(
    repo_path: Path,
    audit: RepoAuditResult,
    custom_message: Optional[str] = None,
    is_dry_run: bool = False,
    is_commit_only: bool = False,
    is_push_only: bool = False,
) -> RepoExecutionResult:
    """Stage, commit, and push changes for a single repository."""
    res = RepoExecutionResult(
        name=audit.name,
        relative_path=audit.relative_path,
        branch=audit.branch,
        commit_status=OperationStatusType.SKIPPED,
        push_status=OperationStatusType.SKIPPED,
    )

    if audit.status == RepoStatusType.NO_COMMITS:
        return res

    if audit.is_detached:
        res.error_message = "Repository is in detached HEAD state; skipping automatic commit/push"
        return res

    # 1. Commit Phase
    should_commit = audit.is_dirty and not is_push_only
    if should_commit:
        commit_msg = (
            custom_message.strip()
            if (custom_message and custom_message.strip())
            else f"chore(sync): automated repository sync and working tree commit ({audit.dirty_file_count} files)"
        )
        res.commit_message = commit_msg

        if is_dry_run:
            res.commit_status = OperationStatusType.SIMULATED
        else:
            code_add, _, err_add = run_cmd(["git", "add", "-A"], repo_path)
            if code_add != 0:
                res.commit_status = OperationStatusType.FAILED
                res.error_message = f"git add failed: {err_add}"
                return res

            code_ci, ci_out, err_ci = run_cmd(["git", "commit", "-m", commit_msg], repo_path)
            if code_ci == 0:
                res.commit_status = OperationStatusType.SUCCESS
                code_sha, sha_out, _ = run_cmd(["git", "rev-parse", "--short", "HEAD"], repo_path)
                res.commit_sha = sha_out if code_sha == 0 else ""
            else:
                res.commit_status = OperationStatusType.FAILED
                res.error_message = f"git commit failed: {err_ci}"
                return res

    # 2. Push Phase
    has_pending_commits = (audit.ahead_count > 0) or (res.commit_status == OperationStatusType.SUCCESS)
    should_push = has_pending_commits and not is_commit_only and bool(audit.primary_remote)

    if should_push:
        remote = audit.primary_remote or "origin"
        branch = audit.branch

        if is_dry_run:
            res.push_status = OperationStatusType.SIMULATED
        else:
            push_cmd = ["git", "push"]
            if not audit.has_upstream:
                push_cmd.extend(["-u", remote, branch])
            else:
                push_cmd.extend([remote, branch])

            code_push, _, err_push = run_cmd(push_cmd, repo_path)
            if code_push == 0:
                res.push_status = OperationStatusType.SUCCESS
            else:
                err_clean = err_push.replace("\r", " ").replace("\n", " ")
                if any(kw in err_push.lower() for kw in ("permission to", "denied to", "403", "forbidden")):
                    res.push_status = OperationStatusType.SKIPPED
                    res.error_message = "External upstream (read-only / 403 denied)"
                else:
                    res.push_status = OperationStatusType.FAILED
                    res.error_message = f"git push failed ({err_clean[:180]})"

    return res


def run_orchestration(
    target_dir: Optional[str] = None,
    is_dry_run: bool = False,
    is_check_only: bool = False,
    is_commit_only: bool = False,
    is_push_only: bool = False,
    custom_message: Optional[str] = None,
    as_json: bool = False,
) -> Dict[str, Any]:
    """Execute complete discovery, audit, and push orchestration across target workspace."""
    start_time = time.time()
    workspace_root = detect_workspace_root(target_dir)
    repositories = discover_repositories(workspace_root)

    audits: List[RepoAuditResult] = []
    executions: List[RepoExecutionResult] = []

    for repo_path in repositories:
        audit = audit_repository(repo_path, workspace_root)
        audits.append(audit)

        if not is_check_only:
            exec_res = execute_commit_and_push(
                repo_path=repo_path,
                audit=audit,
                custom_message=custom_message,
                is_dry_run=is_dry_run,
                is_commit_only=is_commit_only,
                is_push_only=is_push_only,
            )
            executions.append(exec_res)

    duration = round(time.time() - start_time, 2)

    total_repos = len(audits)
    dirty_repos = sum(1 for a in audits if a.is_dirty)
    ahead_repos = sum(1 for a in audits if a.ahead_count > 0)
    clean_repos = sum(1 for a in audits if a.status == RepoStatusType.CLEAN)

    committed_count = sum(1 for e in executions if e.commit_status == OperationStatusType.SUCCESS)
    pushed_count = sum(1 for e in executions if e.push_status == OperationStatusType.SUCCESS)
    failed_count = sum(
        1 for e in executions
        if e.commit_status == OperationStatusType.FAILED or e.push_status == OperationStatusType.FAILED
    )

    report_payload = {
        "workspace_root": str(workspace_root),
        "total_repositories": total_repos,
        "clean_repositories": clean_repos,
        "dirty_repositories": dirty_repos,
        "ahead_repositories": ahead_repos,
        "committed_count": committed_count,
        "pushed_count": pushed_count,
        "failed_count": failed_count,
        "is_dry_run": is_dry_run,
        "duration_seconds": duration,
        "audits": [asdict(a) for a in audits],
        "executions": [asdict(e) for e in executions],
    }

    if as_json:
        print(json.dumps(report_payload, indent=2))
        return report_payload

    # Human-Readable Formatted Console Table
    print(f"\n{'=' * 80}")
    print(f"  MULTI-REPOSITORY ORCHESTRATION REPORT")
    print(f"{'=' * 80}")
    print(f"  Workspace Root : {workspace_root}")
    print(f"  Total Repos    : {total_repos}")
    print(f"  Clean Repos    : {clean_repos}")
    print(f"  Dirty Repos    : {dirty_repos}")
    print(f"  Ahead Repos    : {ahead_repos}")
    if not is_check_only:
        print(f"  Committed      : {committed_count}")
        print(f"  Pushed         : {pushed_count}")
        print(f"  Failed         : {failed_count}")
        print(f"  Mode           : {'SIMULATION (DRY-RUN)' if is_dry_run else 'LIVE EXECUTION'}")
    print(f"  Duration       : {duration}s")
    print(f"{'-' * 80}")
    print(f"{'REPOSITORY':<35} {'BRANCH':<12} {'STATUS':<10} {'DIRTY':<7} {'AHEAD':<7} {'OUTCOME':<15}")
    print(f"{'-' * 80}")

    for a in audits:
        rel = a.relative_path if len(a.relative_path) <= 34 else a.relative_path[:31] + "..."
        outcome_parts = []
        matching_exec = next((e for e in executions if e.relative_path == a.relative_path), None)
        if matching_exec:
            if matching_exec.commit_status in (OperationStatusType.SUCCESS, OperationStatusType.SIMULATED):
                outcome_parts.append(f"CI:{matching_exec.commit_status.value}")
            if matching_exec.push_status in (OperationStatusType.SUCCESS, OperationStatusType.SIMULATED):
                outcome_parts.append(f"PU:{matching_exec.push_status.value}")
            elif matching_exec.push_status == OperationStatusType.SKIPPED and matching_exec.error_message:
                outcome_parts.append("SKIP:external")
            elif matching_exec.commit_status == OperationStatusType.FAILED or matching_exec.push_status == OperationStatusType.FAILED:
                outcome_parts.append("ERR")
        outcome_str = " | ".join(outcome_parts) if outcome_parts else ("AUDIT" if is_check_only else "IDLE")

        dirty_str = str(a.dirty_file_count) if a.is_dirty else "-"
        ahead_str = str(a.ahead_count) if a.ahead_count > 0 else "-"
        print(f"{rel:<35} {a.branch:<12} {a.status.value:<10} {dirty_str:<7} {ahead_str:<7} {outcome_str:<15}")

    print(f"{'=' * 80}\n")
    return report_payload


def run_self_tests() -> bool:
    """Run internal test suite validating discovery, audit, and commit/push logic."""
    print("Running self-tests for 49-commit-and-push-all-repos.py...")
    with tempfile.TemporaryDirectory() as tmp_dir:
        tmp_path = Path(tmp_dir)

        # 1. Test Workspace Root Auto-Detection with explicit path
        detected = detect_workspace_root(str(tmp_path))
        assert detected == tmp_path, f"Expected {tmp_path}, got {detected}"

        # 2. Test Repository Creation & Discovery
        repo1 = tmp_path / "repo-alpha"
        repo1.mkdir()
        run_cmd(["git", "init", "-b", "main"], repo1)
        run_cmd(["git", "config", "user.name", "Tester"], repo1)
        run_cmd(["git", "config", "user.email", "tester@example.com"], repo1)

        # Create .gitignore ignoring target
        gitignore_file = repo1 / ".gitignore"
        gitignore_file.write_text("target/\n", encoding="utf-8")
        run_cmd(["git", "add", ".gitignore"], repo1)
        run_cmd(["git", "commit", "-m", "add gitignore"], repo1)

        # Create build folder that must be skipped
        fake_build = tmp_path / "repo-alpha" / "target" / "nested-git"
        fake_build.mkdir(parents=True)
        run_cmd(["git", "init"], fake_build)

        discovered = discover_repositories(tmp_path)
        assert repo1 in discovered, "repo-alpha must be discovered"
        assert fake_build not in discovered, "nested git in target must be skipped"

        # 4. Commit test file
        test_file = repo1 / "readme.md"
        test_file.write_text("# Test Repo\n", encoding="utf-8")
        run_cmd(["git", "add", "readme.md"], repo1)
        run_cmd(["git", "commit", "-m", "add readme"], repo1)

        audit_clean = audit_repository(repo1, tmp_path)
        assert audit_clean.status == RepoStatusType.CLEAN, f"Expected clean, got {audit_clean.status}"

        # 5. Modify file to create dirty state
        test_file.write_text("# Test Repo Modified\n", encoding="utf-8")
        audit_dirty = audit_repository(repo1, tmp_path)
        assert audit_dirty.is_dirty is True
        assert audit_dirty.status == RepoStatusType.DIRTY

        # 6. Test Dry-Run Execution
        exec_dry = execute_commit_and_push(repo1, audit_dirty, is_dry_run=True)
        assert exec_dry.commit_status == OperationStatusType.SIMULATED

        # 7. Test Real Commit Execution
        exec_real = execute_commit_and_push(repo1, audit_dirty, is_dry_run=False, is_commit_only=True)
        assert exec_real.commit_status == OperationStatusType.SUCCESS
        assert exec_real.commit_sha is not None and len(exec_real.commit_sha) > 0

        audit_post = audit_repository(repo1, tmp_path)
        assert audit_post.is_dirty is False

    print("All self-tests PASSED successfully!")
    return True


def main() -> None:
    parser = argparse.ArgumentParser(
        description="Autonomous Multi-Repository Commit & Push Orchestrator across target workspace (e.g. D:\\work)."
    )
    parser.add_argument(
        "--dir",
        type=str,
        default=None,
        help="Workspace root directory (default: auto-detect D:\\work, ~/git-work, or WORK_DIR).",
    )
    parser.add_argument(
        "--check",
        action="store_true",
        help="Audit mode only. Display repository status without modifying working trees or pushing.",
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Simulate staging, committing, and pushing without performing disk mutations or remote pushes.",
    )
    parser.add_argument(
        "--commit-only",
        action="store_true",
        help="Stage and commit dirty repositories without pushing to remotes.",
    )
    parser.add_argument(
        "--push-only",
        action="store_true",
        help="Push ahead repositories without creating commits for dirty working trees.",
    )
    parser.add_argument(
        "--message",
        "-m",
        type=str,
        default=None,
        help="Custom commit message for dirty repositories (defaults to conventional chore(sync): ...).",
    )
    parser.add_argument(
        "--json",
        action="store_true",
        help="Output results in JSON format for automated pipelines.",
    )
    parser.add_argument(
        "--self-test",
        action="store_true",
        help="Run internal unit and integration tests.",
    )

    args = parser.parse_args()

    if args.self_test:
        success = run_self_tests()
        sys.exit(0 if success else 1)

    try:
        report = run_orchestration(
            target_dir=args.dir,
            is_dry_run=args.dry_run,
            is_check_only=args.check,
            is_commit_only=args.commit_only,
            is_push_only=args.push_only,
            custom_message=args.message,
            as_json=args.json,
        )
        if report.get("failed_count", 0) > 0:
            sys.exit(1)
    except Exception as exc:
        print(f"Error: {exc}", file=sys.stderr)
        sys.exit(1)


if __name__ == "__main__":
    main()
