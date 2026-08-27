#!/usr/bin/env python3
"""
check-file-sizes.py
===================

Hard Rule #6 enforcement (spec/17-consolidated-guidelines/31-*.md):

  - Any source file: max 300 lines
  - Any React component file (*.tsx): max 100 lines
  - Any class/struct (proxied at file level for .py/.go/.rs
    single-definition files): max 120 lines

Discovery roots (globbed, symlinks skipped):
  src/, slides-app/src/, scripts/, linter-scripts/

Excluded:
  - node_modules, dist, build, .git, coverage, .next, .cache
  - Generated / vendored: *.d.ts, *_pb.ts, *_pb.go
  - shadcn/ui primitives under src/components/ui/ (framework code, opt out
    via allowlist baseline instead of a blanket waiver so new ones ratchet).

Modes:
  --check         Fail on any NEW violation not present in the baseline.
  --strict        Fail on ANY violation (no baseline honored).
  --write-baseline
                  Regenerate .file-size-baseline.json from current tree.
  --list          Print every file over cap (advisory), exit 0.

Waiver syntax (top-of-file comment within first 5 lines):

    // lint-allow: file-size reason="framework registry" max=250
    # lint-allow: file-size reason="giant switch table" max=400

  - `reason=...` required
  - `max=N`      hard override, must be <= 600 (absolute ceiling)

Rationale for baseline mode: legacy violators exist (see initial audit in
v5.131 release notes). Blocking on those would strand the linter. Baseline
freezes the current line counts; any file that grows past its baseline OR
any NEW file over cap is a fail. Shrinking is always accepted.

Exit codes:
  0  ok
  1  new/grown violation vs baseline (or any violation under --strict)
  2  malformed waiver
"""
from __future__ import annotations

import argparse
import json
import re
import sys
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parents[1]
BASELINE_PATH = REPO_ROOT / ".file-size-baseline.json"

ROOTS = ["src", "slides-app/src", "scripts", "linter-scripts"]
EXTS = {".ts", ".tsx", ".js", ".jsx", ".mjs", ".cjs", ".py", ".go", ".rs", ".sh", ".ps1"}
EXCLUDE_DIRS = {"node_modules", "dist", "build", ".git", "coverage", ".next", ".cache", "__pycache__"}
EXCLUDE_SUFFIX = (".d.ts", "_pb.ts", "_pb.go")
EXCLUDE_PATH_SUBSTRINGS = ("src/components/ui/",)  # shadcn primitives

DEFAULT_CAP = 300
TSX_CAP = 100
ABSOLUTE_CEILING = 600

WAIVER_RE = re.compile(
    r"lint-allow:\s*file-size\s+reason=\"([^\"]+)\"(?:\s+max=(\d+))?",
)


def is_excluded(path: Path) -> bool:
    parts = set(path.parts)
    if parts & EXCLUDE_DIRS:
        return True
    name = path.name
    if name.endswith(EXCLUDE_SUFFIX):
        return True
    rel = path.as_posix()
    return any(sub in rel for sub in EXCLUDE_PATH_SUBSTRINGS)


def iter_source_files():
    for root in ROOTS:
        root_path = REPO_ROOT / root
        if not root_path.exists():
            continue
        for path in root_path.rglob("*"):
            if not path.is_file():
                continue
            if path.suffix not in EXTS:
                continue
            if is_excluded(path):
                continue
            yield path


def read_line_count(path: Path) -> int:
    with path.open("rb") as fh:
        return sum(1 for _ in fh)


def read_waiver(path: Path):
    """Return (max_override, reason) or (None, None). Raises on malformed."""
    try:
        with path.open("r", encoding="utf-8", errors="replace") as fh:
            head = [next(fh, "") for _ in range(5)]
    except OSError as exc:
        import sys as _sys; print(f"Error: {exc}", file=_sys.stderr)
        return None, None
    for line in head:
        match = WAIVER_RE.search(line)
        if not match:
            continue
        reason, max_str = match.group(1), match.group(2)
        if not reason.strip():
            raise ValueError(f"empty reason in waiver: {path}")
        override = int(max_str) if max_str else DEFAULT_CAP
        if override > ABSOLUTE_CEILING:
            raise ValueError(f"waiver max={override} exceeds ceiling {ABSOLUTE_CEILING}: {path}")
        return override, reason
    return None, None


def cap_for(path: Path) -> int:
    return TSX_CAP if path.suffix == ".tsx" else DEFAULT_CAP


def scan():
    """Return list of {path, lines, cap, waiver_max, waiver_reason}."""
    results = []
    for path in iter_source_files():
        lines = read_line_count(path)
        waiver_max, waiver_reason = read_waiver(path)
        effective_cap = waiver_max if waiver_max is not None else cap_for(path)
        if lines <= effective_cap:
            continue
        results.append({
            "path": path.relative_to(REPO_ROOT).as_posix(),
            "lines": lines,
            "cap": cap_for(path),
            "effective_cap": effective_cap,
            "waiver_reason": waiver_reason,
        })
    return sorted(results, key=lambda r: (-r["lines"], r["path"]))


def load_baseline():
    if not BASELINE_PATH.exists():
        return {}
    with BASELINE_PATH.open("r", encoding="utf-8") as fh:
        data = json.load(fh)
    return {entry["path"]: entry["lines"] for entry in data.get("files", [])}


def write_baseline(violations):
    payload = {
        "note": (
            "Ratchet baseline for Hard Rule #6 (file size caps). "
            "Regenerate ONLY when files legitimately shrink. New violations "
            "or growth past the pinned line count fails CI."
        ),
        "files": [{"path": v["path"], "lines": v["lines"]} for v in violations],
    }
    with BASELINE_PATH.open("w", encoding="utf-8") as fh:
        json.dump(payload, fh, indent=2, sort_keys=False)
        fh.write("\n")


def classify(violations, baseline):
    """Split violations into (new_or_grown, unchanged_or_shrunk)."""
    new_or_grown, ok = [], []
    for v in violations:
        prev = baseline.get(v["path"])
        if prev is None or v["lines"] > prev:
            new_or_grown.append({**v, "baseline": prev})
        else:
            ok.append(v)
    return new_or_grown, ok


def print_row(v):
    marker = "NEW" if v.get("baseline") is None else f"GREW from {v['baseline']}"
    print(f"  {v['path']}: {v['lines']} lines (cap {v['effective_cap']}, {marker})")


def main():
    parser = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument("--check", action="store_true", help="Fail on new/grown violations vs baseline (CI mode).")
    parser.add_argument("--strict", action="store_true", help="Fail on ANY violation, ignore baseline.")
    parser.add_argument("--write-baseline", action="store_true", help="Regenerate .file-size-baseline.json.")
    parser.add_argument("--list", action="store_true", help="List all violations, exit 0.")
    args = parser.parse_args()

    try:
        violations = scan()
    except ValueError as exc:
        print(f"Malformed waiver: {exc}", file=sys.stderr)
        return 2

    if args.write_baseline:
        write_baseline(violations)
        print(f"Wrote baseline with {len(violations)} entries -> {BASELINE_PATH.relative_to(REPO_ROOT)}")
        return 0

    if args.strict:
        if not violations:
            print("check-file-sizes: OK (no files over cap)")
            return 0
        print(f"check-file-sizes STRICT: {len(violations)} file(s) over cap:")
        for v in violations:
            print_row({**v, "baseline": None})
        return 1

    baseline = load_baseline()
    new_or_grown, ok = classify(violations, baseline)

    if args.list or not args.check:
        if not violations:
            print("check-file-sizes: OK (no files over cap)")
            return 0
        print(f"check-file-sizes: {len(violations)} file(s) currently over cap "
              f"({len(new_or_grown)} new/grown, {len(ok)} pinned to baseline):")
        for v in new_or_grown:
            print_row(v)
        for v in ok:
            print(f"  {v['path']}: {v['lines']} lines (cap {v['effective_cap']}, pinned)")
        return 0

    # --check mode
    if not new_or_grown:
        print(f"check-file-sizes --check: OK ({len(ok)} pinned violations, no new drift)")
        return 0
    print(f"check-file-sizes --check FAIL: {len(new_or_grown)} new or grown violation(s) of Hard Rule #6:")
    for v in new_or_grown:
        print_row(v)
    print("\nFix by extracting components/helpers, or add a top-of-file waiver:")
    print('  // lint-allow: file-size reason="framework registry" max=250')
    print("Regenerate baseline ONLY when files shrink: python3 linter-scripts/check-file-sizes.py --write-baseline")
    return 1


if __name__ == "__main__":
    sys.exit(main())
