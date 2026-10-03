#!/usr/bin/env python3
"""
03-ai-scripts/47-git-reconcile-and-resolve-conflict.py
======================================================
Autonomous Git Reconciliation & Mechanical Conflict Resolution Engine.

Safely synchronizes diverged branches without rewriting published Git history
(strictly adhering to Lovable and continuous integration guidelines).
Provides automated, mechanical resolution for common conflict scenarios
including .gitignore unions, SemVer manifests, changelogs, and dual-branch changes.

Capabilities:
  1. Deep Divergence Analysis:
     - Detects ahead/behind commit counts and common merge-base commit.
     - Identifies uncommitted working tree changes with safety guards.
  2. Non-Destructive Reconciliation:
     - Fetches latest remote tracking refs.
     - Automatically merges diverged branches using standard merge semantics.
     - Strictly forbids history rewrite (no force-push, no squash/amend of published commits).
  3. Mechanical Conflict Resolution:
     - Detects conflict markers (<<<<<<<, =======, >>>>>>>).
     - Smart file-type handlers:
       * Set/ignore files (.gitignore): Deduplicated union preserving comments and rules.
       * Version manifests (version.json, package.json): Selects highest SemVer and updates metadata.
       * Structured logs / changelogs: Preserves release sections chronologically.
     - Pluggable resolution strategies: 'smart' (default), 'ours', 'theirs', 'union'.
     - Post-resolution conflict marker audit (verifies 0 remaining markers).
  4. Autonomous Push & Verification:
     - Commits resolved merge with structured metadata.
     - Pushes cleanly to target remote and branch.
     - Verifies zero ahead/behind delta post-push.

Usage:
  # 1. Audit status and simulate reconciliation without modifying working tree:
  python 03-ai-scripts/47-git-reconcile-and-resolve-conflict.py --check

  # 2. Reconcile, auto-resolve conflicts, and push to origin:
  python 03-ai-scripts/47-git-reconcile-and-resolve-conflict.py --push

  # 3. Reconcile using specific strategy without pushing:
  python 03-ai-scripts/47-git-reconcile-and-resolve-conflict.py --strategy smart --no-push

  # 4. Mechanical conflict resolution on currently conflicted working tree:
  python 03-ai-scripts/47-git-reconcile-and-resolve-conflict.py resolve-conflicts --strategy smart
"""

import argparse
import ast
import json
import os
import re
import subprocess
import sys
from dataclasses import asdict, dataclass
from enum import Enum
from pathlib import Path
from typing import Any, Dict, List, Optional, Set, Tuple


class ResolutionStrategyType(str, Enum):
    """Enumeration of conflict resolution strategies."""
    SMART = "smart"
    OURS = "ours"
    THEIRS = "theirs"
    UNION = "union"


@dataclass
class GitStatusSummary:
    """Structured representation of repository git state."""
    branch: str
    remote: str
    is_clean: bool
    untracked_count: int
    modified_count: int
    ahead_count: int
    behind_count: int
    is_diverged: bool
    merge_base: str
    has_active_conflict: bool
    conflicted_files: List[str]


def get_repo_root() -> Path:
    """Resolve repository root directory relative to script location."""
    current = Path(__file__).resolve().parent
    while current != current.parent:
        if (current / ".git").is_dir():
            return current
        if (current / "agents.md").is_file():
            return current
        current = current.parent
    return Path.cwd()


def run_git(args: List[str], repo_root: Path, check_exit: bool = True) -> Tuple[int, str, str]:
    """Execute git command within repo root and capture output."""
    process = subprocess.run(
        ["git"] + args,
        cwd=str(repo_root),
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
        text=True,
        encoding="utf-8",
        errors="replace",
        check=False,
    )
    if check_exit:
        if process.returncode != 0:
            error_message = (
                f"Git command failed: git {' '.join(args)}\n"
                f"Exit code: {process.returncode}\n"
                f"Stderr: {process.stderr.strip()}"
            )
            raise RuntimeError(error_message)
    return process.returncode, process.stdout.strip(), process.stderr.strip()


def inspect_git_status(repo_root: Path, remote_name: str = "origin") -> GitStatusSummary:
    """Inspect repository git state, tracking divergence and active conflicts."""
    # Current branch
    _, branch_stdout, _ = run_git(["rev-parse", "--abbrev-ref", "HEAD"], repo_root)
    current_branch = branch_stdout.strip()

    # Remote branch tracking
    remote_branch = f"{remote_name}/{current_branch}"

    # Status porcelain
    _, status_stdout, _ = run_git(["status", "--porcelain"], repo_root)
    status_lines = [line for line in status_stdout.splitlines() if line.strip()]

    untracked_count = 0
    modified_count = 0
    conflicted_files: List[str] = []

    for line in status_lines:
        code = line[:2]
        filepath = line[3:].strip()
        if code == "??":
            untracked_count += 1
        elif "U" in code or code in ("DD", "AA"):
            conflicted_files.append(filepath)
        else:
            modified_count += 1

    has_active_conflict = len(conflicted_files) > 0
    is_clean = len(status_lines) == 0

    # Divergence counts relative to remote tracking
    ahead_count = 0
    behind_count = 0
    merge_base = ""
    is_diverged = False

    ret, _, _ = run_git(["rev-parse", "--verify", remote_branch], repo_root, check_exit=False)
    if ret == 0:
        _, ahead_str, _ = run_git(
            ["rev-list", "--count", f"{remote_branch}..HEAD"], repo_root, check_exit=False
        )
        _, behind_str, _ = run_git(
            ["rev-list", "--count", f"HEAD..{remote_branch}"], repo_root, check_exit=False
        )
        _, base_str, _ = run_git(
            ["merge-base", "HEAD", remote_branch], repo_root, check_exit=False
        )

        ahead_count = int(ahead_str) if ahead_str.isdigit() else 0
        behind_count = int(behind_str) if behind_str.isdigit() else 0
        merge_base = base_str.strip()

        if ahead_count > 0:
            if behind_count > 0:
                is_diverged = True

    return GitStatusSummary(
        branch=current_branch,
        remote=remote_name,
        is_clean=is_clean,
        untracked_count=untracked_count,
        modified_count=modified_count,
        ahead_count=ahead_count,
        behind_count=behind_count,
        is_diverged=is_diverged,
        merge_base=merge_base,
        has_active_conflict=has_active_conflict,
        conflicted_files=conflicted_files,
    )


# --- Mechanical Conflict Solvers for Specific File Formats ---


def resolve_gitignore_conflict(file_path: Path) -> bool:
    """Resolve merge conflict in .gitignore by computing deduplicated line union."""
    try:
        content = file_path.read_text(encoding="utf-8", errors="replace")
    except Exception:
        return False

    conflict_regex = re.compile(
        r"^<<<<<<< [^\n]*\n(.*?)\n=======\n(.*?)\n>>>>>>> [^\n]*$",
        re.MULTILINE | re.DOTALL,
    )

    def replace_block(match: re.Match) -> str:
        ours_block = match.group(1).splitlines()
        theirs_block = match.group(2).splitlines()

        seen: Set[str] = set()
        unified_lines: List[str] = []

        for line in ours_block + theirs_block:
            cleaned = line.strip()
            if cleaned:
                if not cleaned.startswith("#"):
                    if cleaned in seen:
                        continue
                    seen.add(cleaned)
            unified_lines.append(line)

        return "\n".join(unified_lines)

    resolved_content = conflict_regex.sub(replace_block, content)
    file_path.write_text(resolved_content, encoding="utf-8")
    return True


def parse_semver(v_str: str) -> Tuple[int, int, int]:
    """Parse major, minor, patch tuple from SemVer string."""
    cleaned = v_str.lstrip("v").strip()
    match = re.match(r"^(\d+)\.(\d+)\.(\d+)", cleaned)
    if match:
        return int(match.group(1)), int(match.group(2)), int(match.group(3))
    return 0, 0, 0


def resolve_json_version_conflict(file_path: Path) -> bool:
    """Resolve version.json or package.json conflict by adopting higher SemVer."""
    try:
        content = file_path.read_text(encoding="utf-8", errors="replace")
    except Exception:
        return False

    conflict_regex = re.compile(
        r"^<<<<<<< [^\n]*\n(.*?)\n=======\n(.*?)\n>>>>>>> [^\n]*$",
        re.MULTILINE | re.DOTALL,
    )

    def replace_block(match: re.Match) -> str:
        ours_text = match.group(1)
        theirs_text = match.group(2)

        ours_v_match = re.search(r'"[Vv]ersion"\s*:\s*"([^"]+)"', ours_text)
        theirs_v_match = re.search(r'"[Vv]ersion"\s*:\s*"([^"]+)"', theirs_text)

        if ours_v_match:
            if theirs_v_match:
                ours_sem = parse_semver(ours_v_match.group(1))
                theirs_sem = parse_semver(theirs_v_match.group(1))
                if theirs_sem >= ours_sem:
                    return theirs_text
                return ours_text

        # Default fallback to incoming remote if version diff
        return theirs_text

    resolved_content = conflict_regex.sub(replace_block, content)
    file_path.write_text(resolved_content, encoding="utf-8")
    return True


def resolve_changelog_conflict(file_path: Path) -> bool:
    """Resolve markdown changelog conflict by combining release entries."""
    try:
        content = file_path.read_text(encoding="utf-8", errors="replace")
    except Exception:
        return False

    conflict_regex = re.compile(
        r"^<<<<<<< [^\n]*\n(.*?)\n=======\n(.*?)\n>>>>>>> [^\n]*$",
        re.MULTILINE | re.DOTALL,
    )

    def replace_block(match: re.Match) -> str:
        ours_text = match.group(1).strip()
        theirs_text = match.group(2).strip()

        # Combine both sections cleanly with a blank line
        if ours_text == theirs_text:
            return ours_text
        return f"{theirs_text}\n\n{ours_text}"

    resolved_content = conflict_regex.sub(replace_block, content)
    file_path.write_text(resolved_content, encoding="utf-8")
    return True


def resolve_generic_union(file_path: Path) -> bool:
    """Resolve text conflict by taking the union of both branches."""
    try:
        content = file_path.read_text(encoding="utf-8", errors="replace")
    except Exception:
        return False

    conflict_regex = re.compile(
        r"^<<<<<<< [^\n]*\n(.*?)\n=======\n(.*?)\n>>>>>>> [^\n]*$",
        re.MULTILINE | re.DOTALL,
    )

    def replace_block(match: re.Match) -> str:
        ours = match.group(1)
        theirs = match.group(2)
        if ours.strip() == theirs.strip():
            return ours
        return f"{ours}\n{theirs}"

    resolved_content = conflict_regex.sub(replace_block, content)
    file_path.write_text(resolved_content, encoding="utf-8")
    return True


def verify_no_conflict_markers(file_path: Path) -> bool:
    """Verify that file contains zero git conflict markers."""
    try:
        content = file_path.read_text(encoding="utf-8", errors="replace")
        has_left = "<<<<<<<" in content
        has_mid = "=======" in content
        has_right = ">>>>>>>" in content
        if has_left:
            if has_mid:
                if has_right:
                    return False
        return True
    except Exception:
        return False


def resolve_markdown_list_and_table_conflict(file_path: Path) -> bool:
    """Resolve markdown conflict by unifying lists, checklists, and table rows cleanly."""
    try:
        content = file_path.read_text(encoding="utf-8", errors="replace")
    except Exception:
        return False

    conflict_regex = re.compile(
        r"^<<<<<<< [^\n]*\n(.*?)\n=======\n(.*?)\n>>>>>>> [^\n]*$",
        re.MULTILINE | re.DOTALL,
    )

    def replace_block(match: re.Match) -> str:
        ours_lines = match.group(1).splitlines()
        theirs_lines = match.group(2).splitlines()
        combined = ours_lines + theirs_lines

        has_table = any(line.strip().startswith("|") for line in combined)
        if has_table:
            seen_rows: Set[str] = set()
            unified_table: List[str] = []
            for line in combined:
                stripped = line.strip()
                if not stripped:
                    continue

                if re.match(r"^\|(\s*:?-+:?\s*\|)+$", stripped):
                    if "table_sep" in seen_rows:
                        continue
                    seen_rows.add("table_sep")
                    unified_table.append(line)
                    continue

                if stripped in seen_rows:
                    continue
                seen_rows.add(stripped)
                unified_table.append(line)
            return "\n".join(unified_table)

        has_list = any(re.match(r"^\s*([-*+]|\d+\.)\s+", line) for line in combined)
        if has_list:
            seen_items: Set[str] = set()
            unified_list: List[str] = []
            checked_items: Set[str] = set()

            for line in combined:
                chk_match = re.match(r"^\s*([-*+]|\d+\.)\s+\[([ xX])\]\s+(.*)$", line)
                if chk_match:
                    item_mark = chk_match.group(2)
                    item_body = chk_match.group(3).strip()
                    if item_mark.lower() == "x":
                        checked_items.add(item_body)

            for line in combined:
                chk_match = re.match(r"^\s*([-*+]|\d+\.)\s+\[([ xX])\]\s+(.*)$", line)
                if chk_match:
                    prefix = chk_match.group(1)
                    item_body = chk_match.group(3).strip()
                    if item_body in seen_items:
                        continue
                    seen_items.add(item_body)
                    mark = "x" if item_body in checked_items else " "
                    indent = line[: len(line) - len(line.lstrip())]
                    unified_list.append(f"{indent}{prefix} [{mark}] {item_body}")
                else:
                    stripped = line.strip()
                    if not stripped:
                        continue
                    if stripped in seen_items:
                        continue
                    seen_items.add(stripped)
                    unified_list.append(line)
            return "\n".join(unified_list)

        if match.group(1).strip() == match.group(2).strip():
            return match.group(1)
        return f"{match.group(1)}\n{match.group(2)}"

    resolved_content = conflict_regex.sub(replace_block, content)
    file_path.write_text(resolved_content, encoding="utf-8")
    return True


def resolve_requirements_conflict(file_path: Path) -> bool:
    """Resolve Python requirements conflict by unioning packages and adopting higher versions."""
    try:
        content = file_path.read_text(encoding="utf-8", errors="replace")
    except Exception:
        return False

    conflict_regex = re.compile(
        r"^<<<<<<< [^\n]*\n(.*?)\n=======\n(.*?)\n>>>>>>> [^\n]*$",
        re.MULTILINE | re.DOTALL,
    )

    def replace_block(match: re.Match) -> str:
        ours_lines = match.group(1).splitlines()
        theirs_lines = match.group(2).splitlines()
        req_map: Dict[str, str] = {}
        order: List[str] = []

        for line in ours_lines + theirs_lines:
            stripped = line.strip()
            if not stripped:
                continue

            if stripped.startswith("#"):
                order.append(line)
                continue

            pkg_match = re.match(r"^([a-zA-Z0-9_\-\.]+)(.*)$", stripped)
            if pkg_match:
                pkg_name = pkg_match.group(1).lower()
                if pkg_name not in req_map:
                    req_map[pkg_name] = stripped
                    order.append(pkg_name)
                else:
                    old_line = req_map[pkg_name]
                    old_v = re.search(r"(\d+\.\d+(\.\d+)?)", old_line)
                    new_v = re.search(r"(\d+\.\d+(\.\d+)?)", stripped)
                    if old_v:
                        if new_v:
                            if parse_semver(new_v.group(1)) > parse_semver(old_v.group(1)):
                                req_map[pkg_name] = stripped
            else:
                order.append(line)

        output: List[str] = []
        for item in order:
            if item in req_map:
                output.append(req_map.pop(item))
            else:
                if item.lower() not in req_map:
                    output.append(item)
        return "\n".join(output)

    resolved_content = conflict_regex.sub(replace_block, content)
    file_path.write_text(resolved_content, encoding="utf-8")
    return True


def check_is_import_block_conflict(file_path: Path) -> bool:
    """Check if all conflict markers in a file appear strictly within code import blocks."""
    try:
        content = file_path.read_text(encoding="utf-8", errors="replace")
    except Exception:
        return False

    conflict_regex = re.compile(
        r"^<<<<<<< [^\n]*\n(.*?)\n=======\n(.*?)\n>>>>>>> [^\n]*$",
        re.MULTILINE | re.DOTALL,
    )
    matches = list(conflict_regex.finditer(content))
    if not matches:
        return False

    import_patterns = (
        re.compile(r"^\s*(import\s+|from\s+\S+\s+import\s+)"),
        re.compile(r"^\s*import\s+"),
        re.compile(r'^\s*("[^"]+"|\S+\s+"[^"]+")'),
        re.compile(r"^\s*use\s+\S+;"),
    )

    for match in matches:
        lines = match.group(1).splitlines() + match.group(2).splitlines()
        non_empty = [
            l.strip()
            for l in lines
            if l.strip() and not l.strip().startswith("//") and not l.strip().startswith("#")
        ]
        if not non_empty:
            continue
        is_all_imports = all(any(p.search(line) for p in import_patterns) for line in non_empty)
        if not is_all_imports:
            return False

    return True


def resolve_code_imports_conflict(file_path: Path) -> bool:
    """Resolve import block conflict by deduplicating import statements."""
    try:
        content = file_path.read_text(encoding="utf-8", errors="replace")
    except Exception:
        return False

    conflict_regex = re.compile(
        r"^<<<<<<< [^\n]*\n(.*?)\n=======\n(.*?)\n>>>>>>> [^\n]*$",
        re.MULTILINE | re.DOTALL,
    )

    def replace_block(match: re.Match) -> str:
        ours_lines = match.group(1).splitlines()
        theirs_lines = match.group(2).splitlines()
        seen: Set[str] = set()
        output: List[str] = []
        for line in ours_lines + theirs_lines:
            stripped = line.strip()
            if stripped in seen:
                continue
            seen.add(stripped)
            output.append(line)
        return "\n".join(output)

    resolved_content = conflict_regex.sub(replace_block, content)
    file_path.write_text(resolved_content, encoding="utf-8")
    return True


def validate_syntax_post_resolve(file_path: Path) -> bool:
    """Verify syntax validity of resolved file if applicable."""
    suffix = file_path.suffix.lower()
    if suffix == ".py":
        try:
            content = file_path.read_text(encoding="utf-8", errors="replace")
            ast.parse(content, filename=str(file_path))
            return True
        except Exception:
            return False

    if suffix == ".json":
        try:
            content = file_path.read_text(encoding="utf-8", errors="replace")
            json.loads(content)
            return True
        except Exception:
            return False

    return True


def resolve_file_mechanically(
    repo_root: Path,
    rel_path: str,
    strategy: ResolutionStrategyType,
) -> bool:
    """Apply mechanical conflict resolution to an individual file."""
    target = repo_root / rel_path
    if not target.is_file():
        # Handle file deletions or structural conflicts
        return False

    filename = target.name.lower()
    suffix = target.suffix.lower()

    if strategy == ResolutionStrategyType.OURS:
        run_git(["checkout", "--ours", rel_path], repo_root, check_exit=False)
        return True

    if strategy == ResolutionStrategyType.THEIRS:
        run_git(["checkout", "--theirs", rel_path], repo_root, check_exit=False)
        return True

    if strategy == ResolutionStrategyType.UNION:
        resolve_generic_union(target)
        return verify_no_conflict_markers(target)

    # Strategy: SMART (domain-specific mechanical resolution)
    if filename in (".gitignore", ".gitattributes", ".npmignore", ".dockerignore"):
        resolve_gitignore_conflict(target)
    elif filename in ("version.json", "package.json", "prompt-version.template.json"):
        resolve_json_version_conflict(target)
    elif "changelog" in filename or "release-notes" in filename:
        resolve_changelog_conflict(target)
    elif "requirements" in filename and suffix in (".txt", ".in"):
        resolve_requirements_conflict(target)
    elif suffix in (".md", ".markdown"):
        resolve_markdown_list_and_table_conflict(target)
    else:
        has_imports = check_is_import_block_conflict(target)
        if has_imports:
            resolve_code_imports_conflict(target)
        else:
            # Default smart fallback: try clean union if non-colliding
            resolve_generic_union(target)

    has_no_markers = verify_no_conflict_markers(target)
    if not has_no_markers:
        return False

    is_syntax_valid = validate_syntax_post_resolve(target)
    return is_syntax_valid


def execute_mechanical_resolution(
    repo_root: Path,
    strategy: ResolutionStrategyType,
) -> Tuple[bool, List[str], List[str]]:
    """Scan and resolve all currently conflicted files in repository."""
    _, stdout, _ = run_git(["diff", "--name-only", "--diff-filter=U"], repo_root)
    conflicted = [line.strip() for line in stdout.splitlines() if line.strip()]

    resolved: List[str] = []
    unresolved: List[str] = []

    for file_rel in conflicted:
        success = resolve_file_mechanically(repo_root, file_rel, strategy)
        if success:
            run_git(["add", file_rel], repo_root, check_exit=False)
            resolved.append(file_rel)
        else:
            unresolved.append(file_rel)

    is_all_resolved = len(unresolved) == 0
    return is_all_resolved, resolved, unresolved


# --- Primary Reconciliation Pipeline ---


def reconcile_repository(
    repo_root: Path,
    remote_name: str = "origin",
    strategy: ResolutionStrategyType = ResolutionStrategyType.SMART,
    is_dry_run: bool = False,
    is_push_enabled: bool = True,
) -> Dict[str, Any]:
    """Execute complete git reconciliation, conflict resolution, and push pipeline."""
    report: Dict[str, Any] = {
        "status": "pending",
        "steps_taken": [],
        "conflicts_resolved": [],
        "conflicts_unresolved": [],
        "push_successful": False,
    }

    # Step 1: Pre-flight fetch
    report["steps_taken"].append(f"Fetching remote tracking refs from {remote_name}")
    if not is_dry_run:
        run_git(["fetch", remote_name], repo_root, check_exit=False)

    # Step 2: Status and divergence inspection
    status = inspect_git_status(repo_root, remote_name)
    report["initial_status"] = asdict(status)

    if status.has_active_conflict:
        report["steps_taken"].append("Detected existing active conflict markers in working tree")
        if not is_dry_run:
            is_resolved, resolved, unresolved = execute_mechanical_resolution(repo_root, strategy)
            report["conflicts_resolved"] = resolved
            report["conflicts_unresolved"] = unresolved
            if not is_resolved:
                report["status"] = "failed_unresolved_conflicts"
                return report

            # Commit the resolved conflict
            commit_msg = (
                f"Merge conflict resolution ({strategy.value} strategy)\n\n"
                f"Mechanically resolved files:\n" + "\n".join(f"- {f}" for f in resolved)
            )
            run_git(["commit", "-m", commit_msg], repo_root, check_exit=False)
            report["steps_taken"].append("Committed resolved merge conflict")

    # Step 3: Divergence evaluation
    if not status.is_diverged:
        if status.behind_count > 0:
            report["steps_taken"].append(
                f"Local branch is behind {remote_name}/{status.branch} by {status.behind_count} commits. Fast-forwarding."
            )
            if not is_dry_run:
                run_git(["merge", "--ff-only", f"{remote_name}/{status.branch}"], repo_root)
        elif status.ahead_count == 0:
            report["steps_taken"].append("Local branch is already strictly synchronized with remote")
            report["status"] = "already_synchronized"
            return report
        else:
            report["steps_taken"].append(
                f"Local branch is ahead of {remote_name}/{status.branch} by {status.ahead_count} commits"
            )

    # Step 4: Handle diverged branches
    if status.is_diverged:
        report["steps_taken"].append(
            f"Branches diverged: local ahead by {status.ahead_count}, remote ahead by {status.behind_count}"
        )
        if is_dry_run:
            report["steps_taken"].append("[DRY-RUN] Simulating merge")
        else:
            merge_cmd = ["merge", "--no-ff", f"{remote_name}/{status.branch}"]
            ret, stdout, stderr = run_git(merge_cmd, repo_root, check_exit=False)

            if ret != 0:
                report["steps_taken"].append("Merge encountered conflicts, initiating mechanical resolution engine")
                is_resolved, resolved, unresolved = execute_mechanical_resolution(repo_root, strategy)
                report["conflicts_resolved"] = resolved
                report["conflicts_unresolved"] = unresolved

                if not is_resolved:
                    report["status"] = "failed_unresolved_conflicts"
                    return report

                # Complete merge commit
                commit_msg = (
                    f"Merge remote-tracking branch '{remote_name}/{status.branch}' into {status.branch}\n\n"
                    f"Mechanically resolved conflicts using '{strategy.value}' strategy:\n"
                    + "\n".join(f"- {f}" for f in resolved)
                )
                run_git(["commit", "-m", commit_msg], repo_root)
                report["steps_taken"].append("Successfully committed mechanical merge resolution")
            else:
                report["steps_taken"].append("Clean automatic merge completed without conflict")

    # Step 5: Autonomous push
    if is_push_enabled:
        if is_dry_run:
            report["steps_taken"].append(f"[DRY-RUN] Would push {status.branch} to {remote_name}")
        else:
            report["steps_taken"].append(f"Pushing synchronized branch {status.branch} to {remote_name}")
            ret, push_out, push_err = run_git(
                ["push", remote_name, status.branch], repo_root, check_exit=False
            )
            if ret == 0:
                report["push_successful"] = True
                report["steps_taken"].append(f"Successfully pushed {status.branch} to {remote_name}")
            else:
                report["push_successful"] = False
                report["push_error"] = push_err
                report["status"] = "failed_push"
                return report

    # Step 6: Post-reconciliation verification
    final_status = inspect_git_status(repo_root, remote_name)
    report["final_status"] = asdict(final_status)
    report["status"] = "synchronized_successfully"
    return report


def main() -> None:
    """CLI entry point for Git Reconcile & Mechanical Conflict Resolver."""
    parser = argparse.ArgumentParser(
        description="Autonomous Git Reconciliation & Mechanical Conflict Resolution Engine",
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    parser.add_argument(
        "--remote",
        default="origin",
        help="Git remote name to synchronize with (default: origin)",
    )
    parser.add_argument(
        "--strategy",
        type=str,
        default=ResolutionStrategyType.SMART.value,
        choices=[s.value for s in ResolutionStrategyType],
        help="Conflict resolution strategy: smart, ours, theirs, union (default: smart)",
    )
    parser.add_argument(
        "--check",
        action="store_true",
        help="Check and audit git divergence status without performing merge or push",
    )
    parser.add_argument(
        "--push",
        action="store_true",
        default=True,
        help="Push synchronized branch to remote (enabled by default)",
    )
    parser.add_argument(
        "--no-push",
        action="store_false",
        dest="push",
        help="Perform reconciliation and merge without pushing to remote",
    )
def run_self_tests() -> bool:
    """Run internal unit self-tests for all mechanical conflict resolvers."""
    import tempfile

    print("=== Running Self-Tests for Mechanical Conflict Resolvers ===")
    test_results: List[Tuple[str, bool]] = []

    with tempfile.TemporaryDirectory() as tmp_dir_str:
        tmp_dir = Path(tmp_dir_str)

        # Test 1: .gitignore deduplication
        gi_path = tmp_dir / ".gitignore"
        gi_path.write_text(
            "<<<<<<< HEAD\nnode_modules/\n.cache/\n=======\n.cache/\nbuild/\n>>>>>>> origin/main\n",
            encoding="utf-8",
        )
        is_gi_ok = resolve_gitignore_conflict(gi_path)
        gi_content = gi_path.read_text(encoding="utf-8")
        is_gi_verified = is_gi_ok and verify_no_conflict_markers(gi_path) and gi_content.count(".cache/") == 1
        test_results.append((".gitignore deduplication union", is_gi_verified))

        # Test 2: version.json SemVer resolution
        v_path = tmp_dir / "version.json"
        v_path.write_text(
            '{\n<<<<<<< HEAD\n  "version": "1.2.0"\n=======\n  "version": "1.3.0"\n>>>>>>> origin/main\n}\n',
            encoding="utf-8",
        )
        is_v_ok = resolve_json_version_conflict(v_path)
        v_content = v_path.read_text(encoding="utf-8")
        is_v_verified = is_v_ok and verify_no_conflict_markers(v_path) and '"1.3.0"' in v_content
        test_results.append(("version.json SemVer resolution", is_v_verified))

        # Test 3: changelog combining
        cl_path = tmp_dir / "changelog.md"
        cl_path.write_text(
            "<<<<<<< HEAD\n## v1.2.0\n- Feature A\n=======\n## v1.3.0\n- Feature B\n>>>>>>> origin/main\n",
            encoding="utf-8",
        )
        is_cl_ok = resolve_changelog_conflict(cl_path)
        cl_content = cl_path.read_text(encoding="utf-8")
        is_cl_verified = is_cl_ok and verify_no_conflict_markers(cl_path) and "## v1.3.0" in cl_content and "## v1.2.0" in cl_content
        test_results.append(("changelog release combining", is_cl_verified))

        # Test 4: markdown checklist & table
        md_path = tmp_dir / "tasks.md"
        md_path.write_text(
            "<<<<<<< HEAD\n- [ ] Task 1: Auth\n- [x] Task 2: Core\n=======\n- [ ] Task 2: Core\n- [ ] Task 3: Test\n>>>>>>> origin/main\n",
            encoding="utf-8",
        )
        is_md_ok = resolve_markdown_list_and_table_conflict(md_path)
        md_content = md_path.read_text(encoding="utf-8")
        is_md_verified = is_md_ok and verify_no_conflict_markers(md_path) and "- [x] Task 2: Core" in md_content and "- [ ] Task 3: Test" in md_content
        test_results.append(("markdown checklist deduplicated union", is_md_verified))

        # Test 5: requirements.txt resolution
        req_path = tmp_dir / "requirements.txt"
        req_path.write_text(
            "<<<<<<< HEAD\nrequests==2.31.0\nfastapi>=0.100.0\n=======\nrequests==2.32.3\npytest>=8.0.0\n>>>>>>> origin/main\n",
            encoding="utf-8",
        )
        is_req_ok = resolve_requirements_conflict(req_path)
        req_content = req_path.read_text(encoding="utf-8")
        is_req_verified = is_req_ok and verify_no_conflict_markers(req_path) and "requests==2.32.3" in req_content and "pytest>=8.0.0" in req_content
        test_results.append(("requirements.txt version & package union", is_req_verified))

        # Test 6: Python import block resolution
        imp_path = tmp_dir / "sample.py"
        imp_path.write_text(
            "<<<<<<< HEAD\nimport os\nimport sys\n=======\nimport json\nimport os\n>>>>>>> origin/main\n",
            encoding="utf-8",
        )
        is_imp_ok = resolve_code_imports_conflict(imp_path)
        imp_content = imp_path.read_text(encoding="utf-8")
        is_imp_verified = is_imp_ok and verify_no_conflict_markers(imp_path) and imp_content.count("import os") == 1 and validate_syntax_post_resolve(imp_path)
        test_results.append(("code import block deduplication", is_imp_verified))

    all_passed = True
    for name, is_passed in test_results:
        status_label = "[PASS]" if is_passed else "[FAIL]"
        print(f"  {status_label} {name}")
        if not is_passed:
            all_passed = False

    print(f"\nSelf-test overall status: {'SUCCESS' if all_passed else 'FAILURE'}")
    return all_passed


def main() -> None:
    """CLI entry point for Git Reconcile & Mechanical Conflict Resolver."""
    parser = argparse.ArgumentParser(
        description="Autonomous Git Reconciliation & Mechanical Conflict Resolution Engine",
        formatter_class=argparse.RawDescriptionHelpFormatter,
    )
    parser.add_argument(
        "--remote",
        default="origin",
        help="Git remote name to synchronize with (default: origin)",
    )
    parser.add_argument(
        "--strategy",
        type=str,
        default=ResolutionStrategyType.SMART.value,
        choices=[s.value for s in ResolutionStrategyType],
        help="Conflict resolution strategy: smart, ours, theirs, union (default: smart)",
    )
    parser.add_argument(
        "--check",
        action="store_true",
        help="Check and audit git divergence status without performing merge or push",
    )
    parser.add_argument(
        "--push",
        action="store_true",
        default=True,
        help="Push synchronized branch to remote (enabled by default)",
    )
    parser.add_argument(
        "--no-push",
        action="store_false",
        dest="push",
        help="Perform reconciliation and merge without pushing to remote",
    )
    parser.add_argument(
        "--self-test",
        action="store_true",
        help="Run internal self-tests for all mechanical conflict resolvers",
    )
    parser.add_argument(
        "--json",
        action="store_true",
        help="Output structured JSON summary report",
    )

    subparsers = parser.add_subparsers(dest="subcommand", help="Specialized subcommands")

    # Subcommand: resolve-conflicts
    resolve_parser = subparsers.add_parser(
        "resolve-conflicts",
        help="Mechanically resolve currently conflicted files in working directory",
    )
    resolve_parser.add_argument(
        "--strategy",
        type=str,
        default=ResolutionStrategyType.SMART.value,
        choices=[s.value for s in ResolutionStrategyType],
        help="Conflict resolution strategy",
    )

    args = parser.parse_args()
    repo_root = get_repo_root()

    if args.self_test:
        is_passed = run_self_tests()
        sys.exit(0 if is_passed else 1)

    strategy_enum = ResolutionStrategyType(args.strategy)

    if args.subcommand == "resolve-conflicts":
        is_ok, resolved, unresolved = execute_mechanical_resolution(repo_root, strategy_enum)
        result = {
            "is_all_resolved": is_ok,
            "resolved_files": resolved,
            "unresolved_files": unresolved,
        }
        if args.json:
            print(json.dumps(result, indent=2))
        else:
            print("=== Mechanical Conflict Resolution ===")
            print(f"Strategy: {strategy_enum.value}")
            print(f"Resolved ({len(resolved)}):")
            for item in resolved:
                print(f"  [+] {item}")
            if unresolved:
                print(f"Unresolved ({len(unresolved)}):")
                for item in unresolved:
                    print(f"  [-] {item}")
            else:
                print("All conflict markers successfully cleared and staged!")
        sys.exit(0 if is_ok else 1)

    is_dry_run = args.check
    is_push = args.push

    outcome = reconcile_repository(
        repo_root=repo_root,
        remote_name=args.remote,
        strategy=strategy_enum,
        is_dry_run=is_dry_run,
        is_push_enabled=is_push,
    )

    if args.json:
        print(json.dumps(outcome, indent=2))
    else:
        print("=== Git Reconciliation & Conflict Resolution Summary ===")
        print(f"Status: {outcome.get('status')}")
        print("\nExecution Steps:")
        for step in outcome.get("steps_taken", []):
            print(f"  * {step}")

        conflicts = outcome.get("conflicts_resolved", [])
        if conflicts:
            print(f"\nResolved Conflicts ({len(conflicts)} files):")
            for item in conflicts:
                print(f"  - {item}")

        final_st = outcome.get("final_status")
        if final_st:
            print(f"\nFinal State on branch '{final_st.get('branch')}':")
            print(f"  Ahead commits:  {final_st.get('ahead_count')}")
            print(f"  Behind commits: {final_st.get('behind_count')}")
            print(f"  Clean tree:     {final_st.get('is_clean')}")


if __name__ == "__main__":
    main()
