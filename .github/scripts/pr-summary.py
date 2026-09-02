#!/usr/bin/env python3
"""pr-summary.py — Build a concise PR comment summarizing full-suite-guard (cross-platform).

Usage:
  python .github/scripts/pr-summary.py <artifact_dir> <out_path>
"""

from __future__ import annotations

import os
import re
import sys
from pathlib import Path


def gate_row(label: str, env_name: str, detail: str) -> str:
    raw = os.getenv(env_name, "")
    deduped = os.getenv("SHA_DEDUPED", "false") == "true"
    if raw == "success" and deduped:
        icon = "♻ deduped"
    elif raw == "success":
        icon = "✅ pass"
    elif raw == "failure":
        icon = "❌ fail"
    elif raw == "canceled":
        icon = "🛑 canceled"
    elif raw == "skipped":
        icon = "⏭ skipped"
    elif not raw:
        icon = "⏭ not run"
    else:
        icon = f"⚠ {raw}"
    return f"| {label} | {icon} | {detail} |"


def main() -> int:
    if len(sys.argv) < 3:
        print("Usage: pr-summary.py <artifact_dir> <out_path>", file=sys.stderr)
        return 1

    artifact_dir = Path(sys.argv[1])
    out_path = Path(sys.argv[2])

    test_out = artifact_dir / "test-output.txt"
    lint_out = artifact_dir / "lint-output.txt"

    repo = os.getenv("GITHUB_REPO", "unknown/unknown")
    run_id = os.getenv("GITHUB_RUN_ID", "0")
    sha = os.getenv("GITHUB_SHA_SHORT", "HEAD")
    result = os.getenv("FULL_SUITE_RESULT", "unknown")

    run_url = f"https://github.com/{repo}/actions/runs/{run_id}"
    artifacts_url = f"{run_url}#artifacts"

    tests_passed = 0
    tests_failed = 0
    has_test_output = test_out.is_file()
    first_failures: list[str] = []

    if has_test_output:
        try:
            with open(test_out, "r", encoding="utf-8", errors="replace") as f:
                for line in f:
                    if re.match(r"^ok\s", line):
                        tests_passed += 1
                    elif re.match(r"^FAIL\s", line):
                        tests_failed += 1
                    if re.match(r"^--- FAIL:|^FAIL\s", line) and len(first_failures) < 10:
                        first_failures.append(line.rstrip())
        except OSError:
            pass

    lint_issues = 0
    has_lint_output = lint_out.is_file()
    first_lint_lines: list[str] = []

    if has_lint_output:
        try:
            with open(lint_out, "r", encoding="utf-8", errors="replace") as f:
                for line in f:
                    if re.match(r"^[^\s].+:\d+:\d+:", line):
                        lint_issues += 1
                        if len(first_lint_lines) < 10:
                            first_lint_lines.append(line.rstrip())
        except OSError:
            pass

    tests_status = "✅ pass" if tests_failed == 0 else "❌ fail"
    if not has_test_output:
        tests_status = "⚠ no output"

    lint_status = "✅ clean" if lint_issues == 0 else f"❌ {lint_issues} issue(s)"
    if not has_lint_output:
        lint_status = "⚠ no output"

    overall_emoji = "✅" if result == "success" else "❌"

    lines: list[str] = [
        "<!-- gitmap-ci-summary -->",
        f"## {overall_emoji} CI Summary — `{sha[:10]}`",
        "",
        f"Full Suite Guard: **{result}**",
        "",
        "### Gates",
        "",
        "| Gate | Status | Detail |",
        "| --- | --- | --- |",
        gate_row("Spell check (misspell)", "GATE_SPELL_CHECK", "US locale, whole repo"),
        gate_row("Lint (golangci-lint + vet + gofmt + goimports)", "GATE_LINT", "v1.64.8 strict, --issues-exit-code=1"),
        gate_row("Lint script unit tests", "GATE_LINT_SCRIPT_TESTS", "bash + jq"),
        gate_row("Lint baseline diff", "GATE_LINT_BASELINE_DIFF", "soft-gate, new-issues only"),
        gate_row("Repo-policy checks", "GATE_REPO_POLICY", "naming, legacy-refs, generate-drift, layout"),
        gate_row("Vulnerability scan (govulncheck)", "GATE_VULNCHECK", "v1.1.4, third-party = fail"),
        gate_row("Tests (go test ./... -count=1)", "GATE_TEST", f"{tests_passed} pkg(s) ok, {tests_failed} failed"),
        gate_row("JSON snapshot smoke", "GATE_JSON_SNAPSHOT_FAST", "fast subset"),
        gate_row("Installer smoke (Linux/macOS)", "GATE_INSTALLER_SMOKE", "install.sh contract"),
        gate_row("Installer smoke (Windows)", "GATE_INSTALLER_SMOKE_WINDOWS", "install.ps1 contract"),
        gate_row("Full-suite guard", "GATE_FULL_SUITE_GUARD", "test + lint replay, SARIF upload"),
        gate_row("Cross-compile build", "GATE_BUILD", "6 GOOS/GOARCH targets"),
        "",
        "### Detail",
        "",
        "| Stage | Result | Detail |",
        "| --- | --- | --- |",
        f"| `go test ./...` | {tests_status} | {tests_passed} package(s) ok, {tests_failed} failed |",
        f"| `golangci-lint` (strict) | {lint_status} | v1.64.8, --max-issues=0 |",
        "",
    ]

    if first_failures:
        lines.extend(["### First failing tests", "", "```"] + first_failures + ["```", ""])

    if first_lint_lines:
        lines.extend(["### First lint findings", "", "```"] + first_lint_lines + ["```", ""])

    lines.extend([
        "### Logs & artifacts",
        "",
        f"- [Full run logs]({run_url})",
        f"- [Download artifacts]({artifacts_url}) — `full-suite-outputs` contains the raw `test-output.txt` and `lint-output.txt`",
        "",
        "<sub>Reproduce locally: `./scripts/preflight-ci.sh`</sub>",
    ])

    out_path.parent.mkdir(parents=True, exist_ok=True)
    with open(out_path, "w", encoding="utf-8") as f:
        f.write("\n".join(lines) + "\n")

    print(f"✓ pr-summary: wrote {out_path}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
