#!/usr/bin/env python3
"""Cross-platform single linter baseline-diff runner and gate.

Runs one focused golangci-lint analyzer, normalizes the resulting JSON report,
diffs against an optional baseline JSON, and emits GitHub Actions annotations.
Works natively across Windows, Linux, and macOS without bash or jq.
"""
import argparse
import json
import os
import re
import subprocess
import sys


def parse_args():
    parser = argparse.ArgumentParser(description="Single linter baseline-diff checker")
    parser.add_argument("lint_dir", nargs="?", default=os.environ.get("LINT_DIR", "gitmap"),
                        help="Directory to lint (default: gitmap)")
    parser.add_argument("--linter", default=os.environ.get("LINTER", ""),
                        help="The single golangci-lint analyzer to enable")
    parser.add_argument("--baseline", default=os.environ.get("BASELINE", ""),
                        help="Path to baseline JSON report (optional)")
    parser.add_argument("--current-out", default=os.environ.get("CURRENT_OUT", ""),
                        help="Path to write current JSON report (required)")
    parser.add_argument("--text-filter", default=os.environ.get("TEXT_FILTER", ""),
                        help="Optional regex filter on issue text")
    parser.add_argument("--label", default=os.environ.get("LABEL", ""),
                        help="Display label for annotations (default: LINTER)")
    parser.add_argument("--skip-run", action="store_true",
                        default=bool(os.environ.get("SKIP_LINT_RUN")),
                        help="Skip running golangci-lint (use existing current-out file)")
    return parser.parse_args()


def run_linter(lint_dir, linter, current_out):
    os.makedirs(os.path.dirname(os.path.abspath(current_out)), exist_ok=True)
    cmd = [
        "golangci-lint", "run",
        "--no-config",
        "--disable-all",
        f"--enable={linter}",
        "--timeout=5m",
        "--issues-exit-code=0",
        "--out-format=json",
        "./..."
    ]
    try:
        res = subprocess.run(
            cmd,
            cwd=lint_dir,
            capture_output=True,
            text=True,
            encoding="utf-8",
            errors="replace"
        )
    except FileNotFoundError:
        print("ERROR: golangci-lint not on PATH", file=sys.stderr)
        sys.exit(2)

    with open(current_out, "w", encoding="utf-8") as fh:
        fh.write(res.stdout or '{"Issues": []}')


def normalize_report(path, linter, text_filter_re):
    if not path or not os.path.isfile(path) or os.path.getsize(path) == 0:
        return {}

    try:
        with open(path, "r", encoding="utf-8", errors="replace") as fh:
            data = json.load(fh)
    except Exception as exc:
        print(f"WARN: Could not parse JSON from {path}: {exc}", file=sys.stderr)
        return {}

    issues = data.get("Issues") or []
    normalized = {}
    for issue in issues:
        if issue.get("FromLinter") != linter:
            continue
        pos = issue.get("Pos") or {}
        filename = (pos.get("Filename") or "").replace("\\", "/")
        if not filename:
            continue
        line = int(pos.get("Line") or 1)
        col = int(pos.get("Column") or 1)
        raw_text = issue.get("Text") or ""
        clean_text = " ".join(raw_text.split())

        if text_filter_re and not text_filter_re.search(clean_text):
            continue

        key = f"{filename}|{line}|{clean_text}"
        if key not in normalized:
            normalized[key] = (filename, line, col, clean_text)

    return normalized


def main():
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    if hasattr(sys.stderr, "reconfigure"):
        sys.stderr.reconfigure(encoding="utf-8")

    args = parse_args()

    if not args.linter:
        print("ERROR: LINTER env var or --linter argument is required (e.g. LINTER=gocritic)", file=sys.stderr)
        sys.exit(2)

    if not args.current_out:
        print("ERROR: CURRENT_OUT env var or --current-out argument is required", file=sys.stderr)
        sys.exit(2)

    label = args.label if args.label else args.linter
    text_filter_re = re.compile(args.text_filter) if args.text_filter else None

    if not args.skip_run:
        run_linter(args.lint_dir, args.linter, args.current_out)

    current_issues = normalize_report(args.current_out, args.linter, text_filter_re)
    seeding = not (args.baseline and os.path.isfile(args.baseline) and os.path.getsize(args.baseline) > 0)
    baseline_issues = normalize_report(args.baseline, args.linter, text_filter_re) if not seeding else {}

    current_keys = set(current_issues.keys())
    baseline_keys = set(baseline_issues.keys())
    new_keys = sorted(current_keys - baseline_keys)
    new_count = len(new_keys)

    print("========================================================================")
    print(f"  {label.upper()} DIFF (baseline-diff, full-path only)")
    print("========================================================================")
    print(f"  current  : {args.current_out}")
    print(f"  baseline : {args.baseline if not seeding else '<none — seeding mode>'}")
    print(f"  + NEW    : {new_count}")
    print("========================================================================")

    if new_count == 0:
        print(f"OK: no new {label} findings.")
        sys.exit(0)

    for key in new_keys:
        filename, line, col, text = current_issues[key]
        if seeding:
            print(f"::warning file={filename},line={line},col={col}::[{label}] {text} (seeding baseline)")
        else:
            print(f"::error file={filename},line={line},col={col}::[{label}] {text} (NEW vs baseline)")

    if seeding:
        print("Seeding mode — not failing the build.", file=sys.stderr)
        sys.exit(0)

    print(f"\nFAIL: {new_count} new {label} finding(s). Fix the issues above.", file=sys.stderr)
    sys.exit(1)


if __name__ == "__main__":
    main()
