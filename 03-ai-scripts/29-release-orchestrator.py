#!/usr/bin/env python3
"""
29-release-orchestrator.py - Standalone Release Orchestrator

Automates the complete release heavy-lifting lifecycle:
  1. Detects and preserves the original starting branch.
  2. Resolves current version and calculates next SemVer (minor default per Rule 0).
  3. Pre-release quality gates verification (06-cicd-local-runner.py --run-tests).
  4. STEP 1: Creates and checks out dedicated release branch: release/vX.Y.Z.
  5. STEP 2: Executes version bump using the repository bump script (37-bump-version.py / bump_versions.py).
  6. STEP 3: Commits version bump changes in the release branch.
  7. STEP 4: Creates the annotated git tag: vX.Y.Z on that release commit.
  8. STEP 5: Merges release branch commit back to main branch, pushes main, release branch, and tag to origin.
  9. Reverts working tree back to the original starting branch.

Usage:
  python 03-ai-scripts/29-release-orchestrator.py
  python 03-ai-scripts/29-release-orchestrator.py --tier patch
  python 03-ai-scripts/29-release-orchestrator.py --tier minor --scope "Feature release"
  python 03-ai-scripts/29-release-orchestrator.py --tier major --scope "Breaking change"
  python 03-ai-scripts/29-release-orchestrator.py --version 6.42.0
  python 03-ai-scripts/29-release-orchestrator.py --dry-run
  python 03-ai-scripts/29-release-orchestrator.py --no-push
"""

import argparse
import datetime
import json
import os
import re
import subprocess
import sys
from pathlib import Path

# Repository root discovery
REPO_ROOT = Path(__file__).resolve().parent.parent

# Canonical version files
VERSION_JSON = REPO_ROOT / "version.json"
PACKAGE_JSON = REPO_ROOT / "package.json"
README_MD = REPO_ROOT / "readme.md"
CHANGELOG_MD = REPO_ROOT / "changelog.md"
CONSTANTS_GO = REPO_ROOT / "cli" / "constants" / "constants.go"
TEMPLATE_VERSION = REPO_ROOT / "prompt-version.template.json"

# Known bump scripts
NODE_BUMP_SCRIPT = REPO_ROOT / "scripts" / "bump-version.mjs"
PYTHON_BUMP_SCRIPT = REPO_ROOT / ".ai-memory" / "release" / "bump_versions.py"
AI_BUMP_SCRIPT = REPO_ROOT / "03-ai-scripts" / "37-bump-version.py"


def run_cmd(cmd, cwd=None, check=True, capture_output=True):
    """Executes a command with cross-platform safety."""
    target_cwd = cwd or str(REPO_ROOT)
    result = subprocess.run(
        cmd,
        cwd=target_cwd,
        shell=False,
        check=check,
        capture_output=capture_output,
        text=True,
    )

    return result


def get_git_output(*args):
    """Executes a git command and returns stripped stdout."""
    res = run_cmd(["git", *args])

    return res.stdout.strip()


def get_current_branch():
    """Detects and returns current git branch name."""
    branch = get_git_output("rev-parse", "--abbrev-ref", "HEAD")
    if not branch or branch == "HEAD":
        raise RuntimeError("Detached HEAD or unable to determine current git branch.")

    return branch


def get_main_branch():
    """Detects whether repository uses 'main' or 'master'."""
    try:
        branches = get_git_output("branch", "--list", "main", "master")
        if "main" in branches:
            return "main"
        if "master" in branches:
            return "master"
    except Exception:
        pass

    return "main"


def read_canonical_version():
    """Reads current SemVer from version.json or package.json."""
    if VERSION_JSON.is_file():
        try:
            with open(VERSION_JSON, "r", encoding="utf-8") as f:
                data = json.load(f)

            raw_ver = data.get("Version") or data.get("version")
            if raw_ver:
                return str(raw_ver).strip()
        except Exception:
            pass

    if PACKAGE_JSON.is_file():
        try:
            with open(PACKAGE_JSON, "r", encoding="utf-8") as f:
                data = json.load(f)

            raw_ver = data.get("version")
            if raw_ver:
                return str(raw_ver).strip()
        except Exception:
            pass

    raise FileNotFoundError("Could not find canonical version in version.json or package.json.")


def parse_semver(ver_str):
    """Parses X.Y.Z into a tuple of ints (major, minor, patch)."""
    clean_ver = ver_str.lstrip("v")
    match = re.match(r"^(\d+)\.(\d+)\.(\d+)$", clean_ver)
    if not match:
        raise ValueError(f"Invalid SemVer format: '{ver_str}' (expected X.Y.Z)")

    return int(match.group(1)), int(match.group(2)), int(match.group(3))


def calculate_next_version(current_ver, tier):
    """Calculates next SemVer based on tier (Rule 0: default minor, patch resets to 0)."""
    major, minor, patch = parse_semver(current_ver)

    if tier == "patch":
        patch += 1
    elif tier == "minor":
        minor += 1
        patch = 0
    elif tier == "major":
        major += 1
        minor = 0
        patch = 0
    else:
        raise ValueError(f"Unknown bump tier: '{tier}'. Expected patch, minor, or major.")

    return f"{major}.{minor}.{patch}"


def bootstrap_bump_script_if_needed():
    """Creates a basic bump script if none exists in the repository."""
    if AI_BUMP_SCRIPT.is_file() or PYTHON_BUMP_SCRIPT.is_file() or NODE_BUMP_SCRIPT.is_file():
        return

    scripts_dir = REPO_ROOT / "03-ai-scripts"
    scripts_dir.mkdir(parents=True, exist_ok=True)

    bootstrap_content = '''#!/usr/bin/env python3
import argparse
import datetime
import json
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent

parser = argparse.ArgumentParser()
parser.add_argument("--version", "-v", required=True)
parser.add_argument("--scope", "-s", default="Routine release")
args = parser.parse_args()

v_json = ROOT / "version.json"
if v_json.is_file():
    data = json.loads(v_json.read_text(encoding="utf-8"))
    data["version"] = args.version
    data["releaseDate"] = datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d")
    v_json.write_text(json.dumps(data, indent=2) + "\\n", encoding="utf-8")

p_json = ROOT / "package.json"
if p_json.is_file():
    data = json.loads(p_json.read_text(encoding="utf-8"))
    data["version"] = args.version
    p_json.write_text(json.dumps(data, indent=2) + "\\n", encoding="utf-8")

print(f"Successfully bumped to {args.version}")
'''
    with open(AI_BUMP_SCRIPT, "w", encoding="utf-8") as f:
        f.write(bootstrap_content)

    print(f"[*] Bootstrapped missing bump script: {AI_BUMP_SCRIPT.relative_to(REPO_ROOT)}")


def create_and_checkout_release_branch(next_version, dry_run=False):
    """Step 1: Creates and switches to a dedicated release branch."""
    branch_name = f"release/v{next_version}"
    if dry_run:
        print(f"[DRY RUN] Would create and checkout release branch: '{branch_name}'")
        return branch_name

    print(f"[*] Step 1: Creating and switching to release branch: '{branch_name}'...")
    try:
        run_cmd(["git", "checkout", "-b", branch_name])
    except Exception:
        run_cmd(["git", "checkout", branch_name])
    current = get_current_branch()
    print(f"[*] Active branch is now: '{current}'")

    return branch_name


def execute_version_bump(next_version, scope, dry_run=False):
    """Step 2: Executes the version bump on the release branch via scripts or standalone fallback."""
    if dry_run:
        print(f"[DRY RUN] Would bump version to {next_version} (scope: {scope})")
        return

    print(f"[*] Step 2: Executing version bump to v{next_version}...")

    # Check 1: AI Scripts bump script (pure Python)
    if AI_BUMP_SCRIPT.is_file():
        print(f"[*] Invoking AI Python bump script: {AI_BUMP_SCRIPT.relative_to(REPO_ROOT)}")
        res = run_cmd([sys.executable, str(AI_BUMP_SCRIPT), "--version", next_version, "--scope", scope], check=False)
        if res.returncode == 0:
            return

    # Check 2: .ai-memory Python bump script
    if PYTHON_BUMP_SCRIPT.is_file():
        print(f"[*] Invoking Python bump script: {PYTHON_BUMP_SCRIPT.relative_to(REPO_ROOT)}")
        res = run_cmd([sys.executable, str(PYTHON_BUMP_SCRIPT), "--version", next_version, "--scope", scope], check=False)
        if res.returncode == 0:
            return

    # Check 3: Node bump script
    if NODE_BUMP_SCRIPT.is_file():
        print(f"[*] Invoking Node bump script: {NODE_BUMP_SCRIPT.relative_to(REPO_ROOT)}")
        res = run_cmd(["node", str(NODE_BUMP_SCRIPT), "--version", next_version, "--scope", scope], check=False)
        if res.returncode != 0:
            run_cmd(["node", str(NODE_BUMP_SCRIPT), next_version, scope], check=False)
        return

    today_str = datetime.datetime.now(datetime.timezone.utc).strftime("%Y-%m-%d")

    # Fallback in-place updates: version.json
    if VERSION_JSON.is_file():
        with open(VERSION_JSON, "r", encoding="utf-8") as f:
            v_data = json.load(f)
        v_data["version"] = next_version
        v_data["releaseDate"] = today_str
        with open(VERSION_JSON, "w", encoding="utf-8") as f:
            json.dump(v_data, f, indent=2)
            f.write("\n")

    # Fallback: package.json
    if PACKAGE_JSON.is_file():
        with open(PACKAGE_JSON, "r", encoding="utf-8") as f:
            p_data = json.load(f)
        p_data["version"] = next_version
        with open(PACKAGE_JSON, "w", encoding="utf-8") as f:
            json.dump(p_data, f, indent=2)
            f.write("\n")

    # Fallback: changelog.md
    if CHANGELOG_MD.is_file():
        with open(CHANGELOG_MD, "r", encoding="utf-8") as f:
            cl_content = f.read()

        entry_header = f"## [v{next_version}] - {today_str}\n\n### Added\n- {scope}\n\n"
        if f"[v{next_version}]" not in cl_content:
            if "# Changelog\n" in cl_content:
                cl_content = cl_content.replace("# Changelog\n", f"# Changelog\n\n{entry_header}", 1)
            else:
                cl_content = f"# Changelog\n\n{entry_header}{cl_content}"

            with open(CHANGELOG_MD, "w", encoding="utf-8") as f:
                f.write(cl_content)


def stage_and_commit_release(next_version, scope, dry_run=False):
    """Step 3: Stages release files and commits on the release branch."""
    commit_msg = f"release: v{next_version} {scope}"

    if dry_run:
        print(f"[DRY RUN] Would stage changes and commit on release branch: '{commit_msg}'")
        return "dryrun_commit_sha"

    # Stage release-specific and sync-regenerated files
    release_candidates = [
        VERSION_JSON,
        PACKAGE_JSON,
        CHANGELOG_MD,
        README_MD,
        CONSTANTS_GO,
        TEMPLATE_VERSION,
        NODE_BUMP_SCRIPT,
        PYTHON_BUMP_SCRIPT,
        AI_BUMP_SCRIPT,
        REPO_ROOT / ".gitmap" / "release",
        REPO_ROOT / "public" / "health-score.json",
        REPO_ROOT / "src" / "data" / "specTree.json",
        REPO_ROOT / "02-spec" / "19-main-worker-service" / "98-changelog.md",
        REPO_ROOT / "reports" / "spec-verification" / "coverage.md",
    ]
    for vf in release_candidates:
        if vf.exists():
            run_cmd(["git", "add", str(vf)], check=False)

    # Commit
    run_cmd(["git", "commit", "-m", commit_msg])
    commit_sha = get_git_output("rev-parse", "HEAD")
    print(f"[*] Step 3: Committed release changes on release branch: {commit_sha[:8]} ('{commit_msg}')")

    return commit_sha


def create_release_tag(next_version, commit_sha, dry_run=False):
    """Step 4: Creates annotated git tag on the release commit."""
    tag_name = f"v{next_version}"

    if dry_run:
        print(f"[DRY RUN] Would create annotated tag '{tag_name}' at {commit_sha}")
        return tag_name

    # Create annotated tag
    run_cmd(["git", "tag", "-a", tag_name, "-m", f"Release {tag_name}", commit_sha])
    print(f"[*] Step 4: Created annotated tag: {tag_name} -> {commit_sha[:8]}")

    return tag_name


def merge_release_to_main(release_branch, main_branch="main", dry_run=False):
    """Step 5a: Puts release commit back to the main branch via merge."""
    if dry_run:
        print(f"[DRY RUN] Would checkout '{main_branch}' and merge '{release_branch}'")
        return

    print(f"[*] Step 5a: Checking out '{main_branch}' and merging '{release_branch}'...")
    run_cmd(["git", "checkout", main_branch])
    run_cmd(["git", "merge", release_branch])
    print(f"[OK] Merged release branch '{release_branch}' into '{main_branch}'.")


def push_release(release_branch, tag_name, main_branch="main", dry_run=False):
    """Step 5b: Pushes main branch, release branch, and tag to remote repository."""
    if dry_run:
        print(f"[DRY RUN] Would push main '{main_branch}', release branch '{release_branch}', and tag '{tag_name}' to origin")
        return

    print(f"[*] Step 5b: Pushing '{main_branch}' to origin...")
    run_cmd(["git", "push", "origin", main_branch])

    print(f"[*] Pushing release branch '{release_branch}' to origin...")
    run_cmd(["git", "push", "origin", release_branch])

    print(f"[*] Pushing tag '{tag_name}' to origin...")
    run_cmd(["git", "push", "origin", tag_name])


def revert_to_original_branch(original_branch, dry_run=False):
    """Switches git working tree back to the starting branch."""
    if dry_run:
        print(f"[DRY RUN] Would revert back to original branch: '{original_branch}'")
        return

    current = get_current_branch()
    if current != original_branch:
        print(f"[*] Reverting back to original branch: '{original_branch}' (from '{current}')...")
        run_cmd(["git", "checkout", original_branch])

    restored = get_current_branch()
    if restored != original_branch:
        raise RuntimeError(
            f"Failed to restore original branch! Current branch is '{restored}', expected '{original_branch}'"
        )

    print(f"[OK] Working tree successfully restored to original branch: '{restored}'")


def verify_pre_release_quality_gates(dry_run=False, skip_tests=False):
    """Executes full unit test suites and CI quality gates prior to release."""
    if skip_tests:
        print("[!] Warning: Pre-release test execution skipped via --skip-tests flag.")
        return

    if dry_run:
        print("[DRY RUN] Would execute full unit test suites and CI quality gates: python 03-ai-scripts/06-cicd-local-runner.py --run-tests")
        return

    print("[*] Running full pre-release unit test suites and CI quality gates (python 03-ai-scripts/06-cicd-local-runner.py --run-tests)...")
    runner_script = REPO_ROOT / "03-ai-scripts" / "06-cicd-local-runner.py"
    res = run_cmd([sys.executable, str(runner_script), "--run-tests"], capture_output=False)
    if res.returncode != 0:
        raise RuntimeError("Pre-release quality gates / unit tests failed! Releases are forbidden on failing tests.")


def orchestrate_release(tier="minor", explicit_version=None, scope=None, dry_run=False, push=True, skip_tests=False):
    """Executes the complete 5-step release orchestration flow."""
    # 1. Capture starting branch
    original_branch = get_current_branch()
    main_branch = get_main_branch()
    print(f"[*] Starting release orchestration on branch: '{original_branch}' (main branch: '{main_branch}')")

    # 2. Resolve versions
    current_ver = read_canonical_version()
    if explicit_version:
        next_ver = explicit_version.lstrip("v")
    else:
        next_ver = calculate_next_version(current_ver, tier)

    default_scope = scope or f"Release v{next_ver}"
    print(f"[*] Version Plan: {current_ver} -> {next_ver} (Tier: {tier})")

    # Pre-release quality gates and full unit test execution
    verify_pre_release_quality_gates(dry_run=dry_run, skip_tests=skip_tests)

    release_branch = f"release/v{next_ver}"
    tag_name = f"v{next_ver}"

    try:
        # STEP 1: Create and checkout release branch FIRST
        create_and_checkout_release_branch(next_ver, dry_run=dry_run)

        # STEP 2: Bump version on the release branch using script
        execute_version_bump(next_ver, default_scope, dry_run=dry_run)

        # STEP 3: Commit bump changes in the release branch
        commit_sha = stage_and_commit_release(next_ver, default_scope, dry_run=dry_run)

        # STEP 4: Create annotated tag on that release commit
        create_release_tag(next_ver, commit_sha, dry_run=dry_run)

        # STEP 5: Put that commit back to the main branch (and push)
        merge_release_to_main(release_branch, main_branch=main_branch, dry_run=dry_run)

        is_push_enabled = push and not dry_run
        if is_push_enabled:
            push_release(release_branch, tag_name, main_branch=main_branch, dry_run=dry_run)

    finally:
        # Restore original starting branch if different from current
        revert_to_original_branch(original_branch, dry_run=dry_run)

    print("\n" + "=" * 60)
    print("[OK] RELEASE ORCHESTRATION COMPLETE")
    print(f"  - Starting Branch:  {original_branch}")
    print(f"  - Previous Version: {current_ver}")
    print(f"  - Released Version: {next_ver}")
    print(f"  - Release Branch:   {release_branch}")
    print(f"  - Release Tag:      {tag_name}")
    print(f"  - Merged To Main:   {main_branch}")
    print(f"  - Active Branch:    {get_current_branch()} [Preserved]")
    print("=" * 60 + "\n")


def parse_arguments():
    """Configures CLI argument parser."""
    parser = argparse.ArgumentParser(
        description="29-release-orchestrator: Autonomous release lifecycle with 5-step release branching."
    )
    parser.add_argument(
        "-t",
        "--tier",
        choices=["patch", "minor", "major"],
        default="minor",
        help="SemVer bump tier (default: minor per Rule 0)",
    )
    parser.add_argument(
        "-v",
        "--version",
        dest="explicit_version",
        default=None,
        help="Explicit SemVer string (overrides --tier)",
    )
    parser.add_argument(
        "-s",
        "--scope",
        default=None,
        help="One-line description/scope of the release",
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Simulate the release workflow without modifying files or git",
    )
    parser.add_argument(
        "--no-push",
        action="store_true",
        help="Do not push release branch and tag to remote repository",
    )
    parser.add_argument(
        "--skip-tests",
        action="store_true",
        help="Skip pre-release unit test and quality gate verification (emergency use only)",
    )

    return parser.parse_args()


def main():
    """Main CLI entrypoint."""
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    if hasattr(sys.stderr, "reconfigure"):
        sys.stderr.reconfigure(encoding="utf-8")

    args = parse_arguments()
    should_push = not args.no_push

    orchestrate_release(
        tier=args.tier,
        explicit_version=args.explicit_version,
        scope=args.scope,
        dry_run=args.dry_run,
        push=should_push,
        skip_tests=args.skip_tests,
    )


if __name__ == "__main__":
    main()
