#!/usr/bin/env python3
"""
check-function-lengths.py
=========================

🟡🔴 Tiered enforcement for function body length in scripts/**/*.{sh,ps1}.

Tiers (body lines, brace/closer-exclusive):
  ≤  8 lines  → BEST PRACTICE  (silent OK)
  9–15 lines  → WARN           (printed, exit 0 unless --strict)
 16–25 lines  → CODE RED FAIL  (printed, exit 1) — unless waived
   > 25 lines → HARD FAIL      (printed, exit 1) — only allowed with
                                a `framework=true` waiver, e.g. a
                                language/framework-imposed signature
                                or large `switch`/`match` block.

Waiver syntax (place the comment on the line ABOVE the function header
or as the first body line):

  Bash:        # lint-allow: function-length reason="framework signature" framework=true
  PowerShell:  # lint-allow: function-length reason="big switch" max=25

  - `reason=...`     required (free text, must be quoted)
  - `max=N`          required for the 16–25 tier (N must be ≤ 25)
  - `framework=true` required for the >25 tier (caps at 60)
  - For the warn tier (9–15 lines, within hard max), `reason=...` alone suffices.

Discovery:
  - Walks `scripts/` (override with --root).
  - Includes top-level installer/runner files: fix-repo.{sh,ps1},
    visibility-change.{sh,ps1}, run.{sh,ps1}.

Detection:
  - Bash:   `name() {` ... matching `}` at column 0 (or a single `}` line).
  - PowerShell: `function Verb-Noun {` ... matching `}` via brace counting.

Flags:
  --strict        Treat WARN tier (9–15) as a failure too.
  --max N         Backwards-compat: override CODE-RED ceiling (default 15).
  --hard-max N    Override the absolute hard ceiling (default 25).
  --verbose       Print every checked function with its tier.

Exit codes:
  0  no failures (warnings allowed unless --strict)
  1  one or more CODE-RED or HARD-FAIL violations
  2  usage / discovery error
"""
from __future__ import annotations
import argparse
import re
import sys
from pathlib import Path

BEST_PRACTICE_LINES = 8
MAX_LINES = 15           # CODE RED ceiling (legacy --max default)
HARD_MAX_LINES = 25      # absolute ceiling without framework waiver
FRAMEWORK_HARD_CAP = 60  # even framework waivers cannot exceed this

TOP_LEVEL_FILES = (
    "fix-repo.sh", "fix-repo.ps1",
    "visibility-change.sh", "visibility-change.ps1",
    "run.sh", "run.ps1",
)

BASH_FN_RE = re.compile(r"^\s*([A-Za-z_][A-Za-z0-9_]*)\s*\(\)\s*\{?\s*$")
PS_FN_RE = re.compile(r"^\s*function\s+([A-Za-z_][A-Za-z0-9_-]*)\s*\{?\s*", re.IGNORECASE)
WAIVER_RE = re.compile(
    r"#\s*lint-allow:\s*function-length"
    r"(?:\s+reason=\"([^\"]+)\")?"
    r"(?:\s+max=(\d+))?"
    r"(?:\s+framework=(true|false))?",
    re.IGNORECASE,
)


def is_bash_file(path: Path) -> bool:
    return path.suffix == ".sh"


def is_ps_file(path: Path) -> bool:
    return path.suffix == ".ps1"


def measure_braced_block(lines: list[str], start_idx: int) -> int:
    """Return body line count by brace counting from `{` on/after start_idx."""
    depth = 0
    seen_open = False
    body = 0
    for i in range(start_idx, len(lines)):
        line = lines[i]
        opens = line.count("{")
        closes = line.count("}")
        if not seen_open and opens == 0:
            continue
        seen_open = True
        depth += opens - closes
        if i > start_idx:
            body += 1
        if depth <= 0 and seen_open:
            return max(body - 1, 0)
    return body


def parse_waiver(lines: list[str], header_idx: int) -> dict | None:
    """Look for a `# lint-allow: function-length ...` comment on the line above
    the function header, or as the first body line. Returns parsed fields or None."""
    candidates = []
    if header_idx > 0:
        candidates.append(lines[header_idx - 1])
    if header_idx + 1 < len(lines):
        candidates.append(lines[header_idx + 1])
    for raw in candidates:
        match = WAIVER_RE.search(raw)
        if not match:
            continue
        return {
            "reason": match.group(1) or "",
            "max": int(match.group(2)) if match.group(2) else None,
            "framework": (match.group(3) or "").lower() == "true",
        }
    return None


def classify(length: int, ceiling: int, hard_ceiling: int) -> str:
    """Return one of: 'ok', 'warn', 'fail', 'hard-fail'."""
    if length <= BEST_PRACTICE_LINES:
        return "ok"
    if length <= ceiling:
        return "warn"
    if length <= hard_ceiling:
        return "fail"
    return "hard-fail"


def is_waiver_valid(tier: str, length: int, waiver: dict | None) -> bool:
    if waiver is None:
        return False
    if not waiver["reason"]:
        return False
    if tier == "warn":
        # 9–15 tier (best-practice exceedance, still within max=15) only needs reason.
        return True
    if tier == "fail":
        # 16–25 tier requires explicit max >= length
        return waiver["max"] is not None and length <= waiver["max"] <= HARD_MAX_LINES
    if tier == "hard-fail":
        # >25 requires framework=true and stays under absolute cap
        return waiver["framework"] and length <= FRAMEWORK_HARD_CAP
    return True


def scan_file(path: Path) -> list[tuple[str, int, int, dict | None]]:
    """Return list of (function_name, line_number, body_length, waiver_or_None)."""
    try:
        lines = path.read_text(encoding="utf-8", errors="replace").splitlines()
    except OSError as exc:
        print(f"::error::cannot read {path}: {exc}", file=sys.stderr)
        return []
    fn_re = BASH_FN_RE if is_bash_file(path) else PS_FN_RE
    findings: list[tuple[str, int, int, dict | None]] = []
    for idx, line in enumerate(lines):
        match = fn_re.match(line)
        if not match:
            continue
        name = match.group(1)
        length = measure_braced_block(lines, idx)
        waiver = parse_waiver(lines, idx)
        findings.append((name, idx + 1, length, waiver))
    return findings


def discover_targets(root: Path, repo_root: Path) -> list[Path]:
    targets: list[Path] = []
    if root.exists():
        for ext in ("*.sh", "*.ps1"):
            targets.extend(root.rglob(ext))
    for name in TOP_LEVEL_FILES:
        candidate = repo_root / name
        if candidate.exists():
            targets.append(candidate)
    return sorted(set(targets))


def report(rel: Path, lineno: int, name: str, length: int, tier: str, ceiling: int) -> None:
    if tier == "warn":
        print(f"::warning file={rel},line={lineno}::"
              f"function '{name}' has {length} body lines "
              f"(best ≤{BEST_PRACTICE_LINES}, max ≤{ceiling}; warn tier)")
    elif tier == "fail":
        print(f"::error file={rel},line={lineno}::"
              f"function '{name}' has {length} body lines "
              f"(CODE RED — exceeds {ceiling}; allowed up to {HARD_MAX_LINES} only with "
              f"`# lint-allow: function-length reason=\"...\" max=N`)")
    elif tier == "hard-fail":
        print(f"::error file={rel},line={lineno}::"
              f"function '{name}' has {length} body lines "
              f"(HARD FAIL — exceeds absolute {HARD_MAX_LINES}; only allowed with "
              f"`# lint-allow: function-length reason=\"...\" framework=true`)")


def main() -> int:
    import sys
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--root", default="scripts", help="directory to scan")
    parser.add_argument("--repo-root", default=".", help="repo root for top-level files")
    parser.add_argument("--max", type=int, default=MAX_LINES,
                        help="CODE-RED ceiling (default 15)")
    parser.add_argument("--hard-max", type=int, default=HARD_MAX_LINES,
                        help="absolute hard ceiling (default 25)")
    parser.add_argument("--strict", action="store_true",
                        help="treat WARN tier (9–15 lines) as a failure")
    parser.add_argument("--verbose", action="store_true")
    args = parser.parse_args()

    repo_root = Path(args.repo_root).resolve()
    root = (repo_root / args.root).resolve()
    targets = discover_targets(root, repo_root)
    if not targets:
        print(f"::error::no .sh/.ps1 files found under {root}", file=sys.stderr)
        return 2

    failures = 0
    warnings = 0
    waived = 0
    checked = 0
    for path in targets:
        for name, lineno, length, waiver in scan_file(path):
            checked += 1
            tier = classify(length, args.max, args.hard_max)
            rel = path.relative_to(repo_root)
            if tier == "ok":
                if args.verbose:
                    print(f"OK   {rel}:{lineno} {name} ({length})")
                continue
            if tier in ("warn", "fail", "hard-fail") and is_waiver_valid(tier, length, waiver):
                waived += 1
                if args.verbose:
                    print(f"WAIV {rel}:{lineno} {name} ({length}, tier={tier})")
                continue
            report(rel, lineno, name, length, tier, args.max)
            if tier == "warn":
                warnings += 1
                if args.strict:
                    failures += 1
            else:
                failures += 1

    print(f"\nChecked {checked} functions across {len(targets)} files; "
          f"{failures} failure(s), {warnings} warning(s), {waived} waived.",
          file=sys.stderr)
    return 1 if failures else 0


if __name__ == "__main__":
    sys.exit(main())
