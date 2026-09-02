#!/usr/bin/env python3
"""linter-scripts/check-runner-dispatch-antipatterns.py — Guard against dispatch anti-patterns in runner scripts (cross-platform).

Spec: spec/15-distribution-and-runner/06-fix-repo-forwarding.md
"""

from __future__ import annotations

import datetime
import os
import re
import sys
from pathlib import Path


def extract_region_sh(content: str) -> tuple[int, str] | None:
    for i, line in enumerate(content.splitlines(), start=1):
        if re.match(r"^\s*fix-repo\)", line):
            return i, line
    return None


def extract_region_ps(content: str) -> tuple[int, str] | None:
    for i, line in enumerate(content.splitlines(), start=1):
        if re.match(r'^\s*"fix-repo"\s*\{', line):
            return i, line
    return None


def main() -> int:
    if sys.platform == "win32":
        try:
            sys.stdout.reconfigure(encoding="utf-8", errors="replace")
            sys.stderr.reconfigure(encoding="utf-8", errors="replace")
        except Exception:
            pass

    repo_root = Path(__file__).resolve().parent.parent
    sh_path = repo_root / "run.sh"
    ps_path = repo_root / "run.ps1"

    if not sh_path.is_file() or not ps_path.is_file():
        print("::error::run.sh or run.ps1 missing", file=sys.stderr)
        return 2

    with open(sh_path, "r", encoding="utf-8", errors="replace") as f:
        sh_text = f.read()
    with open(ps_path, "r", encoding="utf-8", errors="replace") as f:
        ps_text = f.read()

    sh_region = extract_region_sh(sh_text)
    ps_region = extract_region_ps(ps_text)

    if not sh_region or not ps_region:
        print("::error::fix-repo dispatch arms not found in runners", file=sys.stderr)
        return 2

    findings: list[tuple[str, int, str, str, str, str]] = []

    def check_forbid(file_name: str, line_no: int, text: str, pattern: str, reason: str) -> None:
        m = re.search(pattern, text)
        if m:
            findings.append((file_name, line_no, "FORBIDDEN", reason, pattern, m.group(0)))

    def check_require(file_name: str, line_no: int, text: str, pattern: str, reason: str) -> None:
        if not re.search(pattern, text):
            findings.append((file_name, line_no, "MISSING", reason, pattern, "(no match)"))

    # Check run.sh
    sh_line, sh_line_text = sh_region
    check_forbid("run.sh", sh_line, sh_line_text, r'"\$\*"', '"$*" collapses argv into one string; use "$@"')
    check_forbid("run.sh", sh_line, sh_line_text, r"\beval\b", "eval-based dispatch breaks argv quoting")
    check_forbid("run.sh", sh_line, sh_line_text, r"\bbash\s+-c\b", "`bash -c` wrapper loses argv boundaries")
    check_forbid("run.sh", sh_line, sh_line_text, r"\bsh\s+-c\b", "`sh -c` wrapper loses argv boundaries")
    check_forbid("run.sh", sh_line, sh_line_text, r"\bxargs\b", "xargs reformats argv via stdin and loses quoting")
    check_require("run.sh", sh_line, sh_line_text, r'\bexec\b[^#]*fix[-_]repo\.(sh|py)[^#]*"\$@"', 'must use `exec ... fix-repo "$@"`')

    # Check run.ps1
    ps_line, ps_line_text = ps_region
    check_forbid("run.ps1", ps_line, ps_line_text, r"\$args\s*-join", "$args -join collapses argv into a single string")
    check_forbid("run.ps1", ps_line, ps_line_text, r"-join\s+\$args", "-join $args collapses argv into a single string")
    check_forbid("run.ps1", ps_line, ps_line_text, r'"\$args"', '"$args" interpolation flattens argv; use @args splatting')
    check_forbid("run.ps1", ps_line, ps_line_text, r"\[string\]::Join", "[string]::Join on $args collapses argv into one string")
    check_forbid("run.ps1", ps_line, ps_line_text, r"\bInvoke-Expression\b", "Invoke-Expression on argv breaks quoting and is unsafe")
    check_forbid("run.ps1", ps_line, ps_line_text, r"\biex\b", "`iex` on argv is unsafe")
    check_forbid("run.ps1", ps_line, ps_line_text, r"Start-Process\b", "Start-Process detaches child")
    check_forbid("run.ps1", ps_line, ps_line_text, r"Start-Job\b", "Start-Job detaches child")
    check_forbid("run.ps1", ps_line, ps_line_text, r"\bcmd\s+/c\b", "`cmd /c` wrapper loses quoting")
    check_require("run.ps1", ps_line, ps_line_text, r"@args", "must invoke inner with @args splatting")
    check_require("run.ps1", ps_line, ps_line_text, r"exit\s+\$LASTEXITCODE", "must end with `exit $LASTEXITCODE`")

    if not findings:
        print("\n✅ runner dispatch guard: no anti-patterns found")
        return 0

    print(f"\n════════════════════ violation summary ({len(findings)}) ════════════════════", file=sys.stderr)
    for i, (f, l, kind, reason, pat, match) in enumerate(findings, start=1):
        print(f"  {i:<3}  {f}:{l}  {kind:<10}  {reason}", file=sys.stderr)
        print(f"       │ regex: {pat}\n       │ match: <<{match}>>", file=sys.stderr)
        print(f"::error file={f},line={l}::[{kind}] {reason}", file=sys.stderr)
    print("═══════════════════════════════════════════════════════════", file=sys.stderr)
    return 1


if __name__ == "__main__":
    sys.exit(main())
