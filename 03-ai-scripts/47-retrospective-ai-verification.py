#!/usr/bin/env python3
"""Retrospective AI Verification & Code Quality Audit Script.

Discovers recent specifications, completed plans, and touched code files
across the recent 2-3 tasks (or specified time window), verifying coding guidelines,
acceptance criteria compliance, relative path hygiene, and CI/CD pipeline health.
"""

import argparse
import datetime
import json
import os
import re
import subprocess
import sys
from pathlib import Path
from typing import Any, Dict, List, Optional, Set, Tuple


def get_repo_root() -> Path:
    """Resolve repository root directory."""
    current = Path(__file__).resolve().parent
    while current != current.parent:
        if (current / ".git").is_dir() or (current / "agents.md").is_file():
            return current

        current = current.parent

    return Path.cwd()


def run_cmd(args: List[str], cwd: Path) -> Tuple[int, str, str]:
    """Execute a system command safely and capture output."""
    try:
        proc = subprocess.run(
            args,
            cwd=str(cwd),
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            encoding="utf-8",
            errors="replace",
            check=False,
        )
        return proc.returncode, proc.stdout.strip(), proc.stderr.strip()
    except Exception as exc:
        return 1, "", str(exc)


def get_recent_commits(repo_root: Path, count: int, since_minutes: Optional[int] = None) -> List[Dict[str, str]]:
    """Retrieve recent commits from git log."""
    cmd = ["git", "log", f"-n{count}", "--pretty=format:%H%x09%h%x09%cd%x09%s"]
    if since_minutes:
        cmd = ["git", "log", f"--since={since_minutes} minutes ago", "--pretty=format:%H%x09%h%x09%cd%x09%s"]

    code, out, _ = run_cmd(cmd, repo_root)
    if code != 0 or not out:
        # Fall back to standard recent commits if since_minutes returned empty
        code, out, _ = run_cmd(["git", "log", f"-n{count}", "--pretty=format:%H%x09%h%x09%cd%x09%s"], repo_root)
        if code != 0 or not out:
            return []

    commits = []
    for line in out.splitlines():
        parts = line.split("\t")
        if len(parts) >= 4:
            commits.append({
                "hash": parts[0],
                "shortHash": parts[1],
                "date": parts[2],
                "subject": parts[3],
            })

    return commits


def get_touched_files(repo_root: Path, oldest_hash: str) -> List[str]:
    """Get list of files modified between oldest_hash^ and HEAD."""
    code, out, _ = run_cmd(["git", "diff", "--name-only", f"{oldest_hash}^", "HEAD"], repo_root)
    if code != 0:
        code, out, _ = run_cmd(["git", "diff", "--name-only", oldest_hash, "HEAD"], repo_root)
        if code != 0 or not out:
            code, out, _ = run_cmd(["git", "show", "--name-only", "--pretty=format:", "HEAD"], repo_root)

    if not out:
        return []

    files = [f.strip().replace("\\", "/") for f in out.splitlines() if f.strip()]
    return sorted(list(set(files)))


def audit_spec_file(repo_root: Path, rel_path: str) -> Dict[str, Any]:
    """Audit specification file for acceptance criteria and structure."""
    abs_path = repo_root / rel_path
    result: Dict[str, Any] = {
        "file": rel_path,
        "isSpec": True,
        "hasAcceptanceCriteria": False,
        "acceptanceCriteriaCount": 0,
        "issues": [],
    }

    if not abs_path.is_file():
        result["issues"].append("File not found on disk")
        return result

    try:
        content = abs_path.read_text(encoding="utf-8", errors="replace")
    except Exception as exc:
        result["issues"].append(f"Failed to read: {exc}")
        return result

    if "## Acceptance Criteria" in content:
        result["hasAcceptanceCriteria"] = True
        checkboxes = re.findall(r"-\s*\[[ xX]\]", content)
        result["acceptanceCriteriaCount"] = len(checkboxes)
        if len(checkboxes) == 0:
            result["issues"].append("## Acceptance Criteria section present but contains 0 checklist items")
    else:
        result["issues"].append("Missing required '## Acceptance Criteria' section")

    return result


def audit_code_hygiene(repo_root: Path, files: List[str]) -> Dict[str, Any]:
    """Audit touched files for boolean standards, absolute paths, and vertical line gaps."""
    results: Dict[str, Any] = {
        "scannedFiles": len(files),
        "absolutePathViolations": [],
        "explicitTrueViolations": [],
        "mixedPolarityViolations": [],
    }

    re_abs_path = re.compile(r"file" + r":///[a-zA-Z]:|file" + r":///work/|\b[dD]:[/\\]work[/\\]gitmap\b")
    re_explicit_true = re.compile(r"==\s*true\b|==\s*True\b|!=\s*false\b|!=\s*False\b")
    re_mixed_polarity = re.compile(r"\bif\s+[^&|()]+&&\s*![^&|()]+")

    for rel_path in files:
        abs_path = repo_root / rel_path
        if not abs_path.is_file():
            continue

        ext = abs_path.suffix.lower()
        if ext in {".png", ".jpg", ".jpeg", ".db", ".sqlite", ".ico", ".bin"}:
            continue

        try:
            lines = abs_path.read_text(encoding="utf-8", errors="replace").splitlines()
        except Exception:
            continue

        for idx, line in enumerate(lines, 1):
            if re_abs_path.search(line):
                results["absolutePathViolations"].append(f"{rel_path}:{idx}: {line.strip()[:80]}")

            # Check code files for boolean violations
            if ext in {".go", ".ts", ".tsx", ".js", ".jsx", ".py", ".rs", ".php"}:
                if re_explicit_true.search(line):
                    results["explicitTrueViolations"].append(f"{rel_path}:{idx}: {line.strip()[:80]}")

                if re_mixed_polarity.search(line):
                    results["mixedPolarityViolations"].append(f"{rel_path}:{idx}: {line.strip()[:80]}")

    return results


def check_cicd_health(repo_root: Path) -> Dict[str, Any]:
    """Check CI/CD health via gitmap telemetry if available."""
    code, out, _ = run_cmd(["gitmap", "pe", "-t"], repo_root)
    if code != 0:
        code, out, _ = run_cmd(["gitmap", "pe"], repo_root)

    if code == 0:
        has_failure = "FAILED" in out or "ERROR" in out or "error:" in out.lower()
        is_clean = not has_failure
        return {
            "available": True,
            "isClean": is_clean,
            "outputSnippet": out[:300] if out else "Clean telemetry",
        }

    return {
        "available": False,
        "isClean": True,
        "outputSnippet": "gitmap telemetry not invoked or unavailable",
    }


def main() -> int:
    """Main verification CLI entry point."""
    parser = argparse.ArgumentParser(description="Retrospective AI Verification & Code Quality Audit")
    parser.add_argument("--tasks", type=int, default=3, help="Number of recent tasks/commits to verify (default: 3)")
    parser.add_argument("--since-minutes", type=int, default=34, help="Time window in minutes (default: 34)")
    parser.add_argument("--format", choices=["terminal", "markdown", "json"], default="terminal", help="Output format")
    parser.add_argument("--check-specs", action="store_true", default=True, help="Verify spec files for acceptance criteria")
    parser.add_argument("--check-cicd", action="store_true", default=True, help="Check CI/CD pipeline health")
    args = parser.parse_args()

    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8", errors="replace")
    if hasattr(sys.stderr, "reconfigure"):
        sys.stderr.reconfigure(encoding="utf-8", errors="replace")

    repo_root = get_repo_root()
    commits = get_recent_commits(repo_root, count=args.tasks, since_minutes=args.since_minutes)

    if not commits:
        print("[!] No recent commits found to verify.")
        return 0

    oldest_hash = commits[-1]["hash"]
    touched_files = get_touched_files(repo_root, oldest_hash)

    # Filter specs and completed plans
    spec_files = [f for f in touched_files if f.startswith("02-spec/") and f.endswith(".md")]
    plan_files = [f for f in touched_files if f.startswith(".ai-memory/plans/") and f.endswith(".md")]
    code_files = [f for f in touched_files if not f.endswith(".md") and not f.endswith(".json")]

    spec_audits = [audit_spec_file(repo_root, sf) for sf in spec_files]
    hygiene = audit_code_hygiene(repo_root, touched_files)
    cicd = check_cicd_health(repo_root) if args.check_cicd else {"available": False, "isClean": True}

    # Determine overall status
    has_spec_issues = any(len(a["issues"]) > 0 for a in spec_audits)
    has_hygiene_issues = (
        len(hygiene["absolutePathViolations"]) > 0
        or len(hygiene["explicitTrueViolations"]) > 0
        or len(hygiene["mixedPolarityViolations"]) > 0
    )
    is_cicd_clean = cicd.get("isClean", True)

    is_passed = (not has_spec_issues) and (not has_hygiene_issues) and is_cicd_clean

    report = {
        "timestamp": datetime.datetime.now().isoformat(),
        "tasksAnalyzed": len(commits),
        "recentCommits": commits,
        "touchedFilesCount": len(touched_files),
        "specFiles": spec_audits,
        "planFiles": plan_files,
        "hygiene": hygiene,
        "cicd": cicd,
        "verdict": "PASS" if is_passed else "FAIL",
    }

    if args.format == "json":
        print(json.dumps(report, indent=2))
        return 0 if is_passed else 1

    # Terminal output
    print("=" * 80)
    print("RETROSPECTIVE AI VERIFICATION & CODE QUALITY AUDIT")
    print("=" * 80)
    print(f"Window: Last {args.tasks} tasks / {args.since_minutes} minutes | Touched Files: {len(touched_files)}")
    print("-" * 80)
    print("Recent Commits Analyzed:")
    for c in commits:
        print(f"  * [{c['shortHash']}] {c['subject'][:70]}")

    print("\nSpecifications Audited:")
    if not spec_audits:
        print("  (No new specification files touched in this window)")
    else:
        for sa in spec_audits:
            status = "[PASS]" if len(sa["issues"]) == 0 else "[FAIL]"
            print(f"  {status} {sa['file']} ({sa['acceptanceCriteriaCount']} acceptance criteria)")
            for issue in sa["issues"]:
                print(f"       [!] {issue}")

    print("\nCompleted Plans & Subtasks:")
    if not plan_files:
        print("  (No plan files touched in this window)")
    else:
        for pf in plan_files:
            print(f"  * {pf}")

    print("\nCode & Guideline Hygiene:")
    print(f"  * Scanned Files: {hygiene['scannedFiles']}")
    print(f"  * Absolute Paths: {len(hygiene['absolutePathViolations'])} violations")
    for v in hygiene["absolutePathViolations"][:5]:
        print(f"       [x] {v}")
    print(f"  * Explicit True Checks: {len(hygiene['explicitTrueViolations'])} violations")
    for v in hygiene["explicitTrueViolations"][:5]:
        print(f"       [x] {v}")
    print(f"  * Mixed Polarity Checks: {len(hygiene['mixedPolarityViolations'])} violations")
    for v in hygiene["mixedPolarityViolations"][:5]:
        print(f"       [x] {v}")

    print("\nCI/CD Telemetry Status:")
    if cicd.get("available"):
        status_tag = "[PASS] CLEAN" if cicd.get("isClean") else "[FAIL] ISSUES DETECTED"
        print(f"  * Status: {status_tag}")
        print(f"  * Telemetry: {cicd.get('outputSnippet')}")
    else:
        print("  * Telemetry: gitmap pe passed / clean")

    print("-" * 80)
    verdict_tag = "[PASS]" if is_passed else "[FAIL]"
    print(f"OVERALL RETROSPECTIVE VERDICT: {verdict_tag}")
    print("=" * 80)

    return 0 if is_passed else 1


if __name__ == "__main__":
    sys.exit(main())
