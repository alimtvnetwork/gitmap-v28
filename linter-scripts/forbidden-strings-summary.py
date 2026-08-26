#!/usr/bin/env python3
"""
Forbidden-Strings Summary Report
=================================

Companion to ``check-forbidden-strings.py``. Produces a focused,
per-rule digest of every current finding plus the **exact** one-line
shell command that fixes them.

Why a separate tool?

``check-forbidden-strings.py`` is a CI gate: it prints every match with
file + line + content so reviewers can audit. It is great for failing
the build, but it is verbose and the generic ``fix_hint`` field in the
TOML is a template, not a runnable command.

This summary tool:

  * Groups findings by rule id (e.g. ``STALE-MODULE-PATH``).
  * Lists each affected file once with a hit-count.
  * Emits a single concrete ``sed`` command per rule that uses the
    canonical ``replacement`` declared in the TOML, scoped to **only
    the files that actually contain the stale text** (not a repo-wide
    sweep).
  * Optionally prints just the command (``--emit-fix-command``) so it
    can be piped to ``bash`` from a CI annotation or a docs snippet.

USAGE
-----

    # Human-readable summary of every rule
    python3 linter-scripts/forbidden-strings-summary.py

    # Single rule
    python3 linter-scripts/forbidden-strings-summary.py --rule STALE-MODULE-PATH

    # Just emit the fix command (one line, runnable)
    python3 linter-scripts/forbidden-strings-summary.py \\
        --rule STALE-MODULE-PATH --emit-fix-command

    # Pipe straight to bash
    python3 linter-scripts/forbidden-strings-summary.py \\
        --rule STALE-MODULE-PATH --emit-fix-command | bash

    # Markdown output (suitable for PR comments / GitHub job summary)
    python3 linter-scripts/forbidden-strings-summary.py --markdown

EXIT CODES
----------
    0  no findings (or --emit-fix-command requested and there is
       nothing to fix — emits a no-op ``: nothing to fix`` line so
       piping to bash is always safe)
    1  one or more findings present (default mode)
    2  invalid usage (unknown rule, missing config, etc.)
"""

from __future__ import annotations

import argparse
import fnmatch
import os
import re
import shlex
import sys
from collections import defaultdict
from pathlib import Path

try:
    import tomllib  # Python 3.11+
except ModuleNotFoundError:
    try:
        import tomli as tomllib  # type: ignore[no-redef]
    except ModuleNotFoundError as exc:
        print(f"Error: {exc}", file=sys.stderr)
        sys.exit("Error: Python 3.11+ required (tomllib), or install 'tomli'.")

REPO_ROOT = Path(__file__).resolve().parent.parent
CONFIG_PATH = REPO_ROOT / "linter-scripts" / "forbidden-strings.toml"
ALWAYS_EXCLUDE_DIRS = {".git", "node_modules", "dist", "build"}


# --------------------------------------------------------------------------- #
# Config + scan (mirrors check-forbidden-strings.py so output is consistent)
# --------------------------------------------------------------------------- #

def load_rules() -> list[dict]:
    if not CONFIG_PATH.exists():
        sys.exit(f"Error: config not found: {CONFIG_PATH}")
    with CONFIG_PATH.open("rb") as fh:
        data = tomllib.load(fh)
    return data.get("rule", [])


def is_excluded_dir(rel_dir: str, extra: list[str]) -> bool:
    parts = rel_dir.replace("\\", "/").split("/")
    excludes = ALWAYS_EXCLUDE_DIRS | set(extra)
    return any(p in excludes for p in parts)


def is_allowlisted(rel: str, allowlist: list[str]) -> bool:
    for pat in allowlist:
        if fnmatch.fnmatch(rel, pat):
            return True
        if rel == pat or rel.startswith(pat + "/"):
            return True
    return False


def scan_rule(rule: dict) -> dict[str, list[tuple[int, str, str]]]:
    """Return ``{rel_path: [(lineno, line_text, exact_match_text), ...]}``."""
    pattern = re.compile(rule["pattern"])
    extra_excl = rule.get("exclude_dirs", [])
    excl_files = rule.get("exclude_files", [])
    allowlist = rule.get("allowlist", [])

    hits: dict[str, list[tuple[int, str, str]]] = defaultdict(list)
    for dirpath, dirnames, filenames in os.walk(REPO_ROOT):
        rel_dir = os.path.relpath(dirpath, REPO_ROOT)
        if rel_dir == ".":
            rel_dir = ""
        if rel_dir and is_excluded_dir(rel_dir, extra_excl):
            dirnames.clear()
            continue
        dirnames[:] = [
            d for d in dirnames
            if d not in ALWAYS_EXCLUDE_DIRS and d not in extra_excl
        ]
        for fname in filenames:
            if any(fnmatch.fnmatch(fname, p) for p in excl_files):
                continue
            full = Path(dirpath) / fname
            rel = full.relative_to(REPO_ROOT).as_posix()
            if is_allowlisted(rel, allowlist):
                continue
            try:
                text = full.read_text(encoding="utf-8", errors="ignore")
            except OSError as exc:
                print(f"Error: {exc}", file=sys.stderr)
                continue
            for lineno, line in enumerate(text.splitlines(), 1):
                for m in pattern.finditer(line):
                    hits[rel].append((lineno, line.rstrip(), m.group(0)))
    return dict(hits)


# --------------------------------------------------------------------------- #
# Fix command generation
# --------------------------------------------------------------------------- #

def collect_unique_matches(hits: dict[str, list[tuple[int, str, str]]]) -> list[str]:
    """Every distinct exact text the regex matched, across all files.

    A rule like ``STALE-REPO-SLUG`` matches ``coding-guidelines-v24``,
    ``coding-guidelines-v24``, etc. — each needs its own
    ``s|stale|new|g`` substitution because the canonical replacement
    is a single string, not a function of the match.
    """
    seen: set[str] = set()
    for entries in hits.values():
        for _, _, exact in entries:
            seen.add(exact)
    return sorted(seen)


def build_fix_command(
    rule: dict,
    hits: dict[str, list[tuple[int, str, str]]],
) -> str | None:
    """Return a single bash command that, when run, replaces every
    stale match in every affected file with the canonical replacement.

    Returns ``None`` if the rule has no ``replacement`` declared (e.g.
    a guard that has no automatic fix).
    """
    replacement = rule.get("replacement")
    if not replacement:
        return None
    if not hits:
        return ": nothing to fix"

    files = sorted(hits.keys())
    files_arg = " ".join(shlex.quote(p) for p in files)
    sed_exprs = []
    for stale in collect_unique_matches(hits):
        if stale == replacement:
            # Defensive guard: never emit a no-op s|x|x|g substitution
            # (the previous hardcoded fix_hint had this bug).
            continue
        sed_exprs.append(f"-e {shlex.quote(f's|{stale}|{replacement}|g')}")
    if not sed_exprs:
        return ": nothing to fix"
    return f"sed -i {' '.join(sed_exprs)} {files_arg}"


# --------------------------------------------------------------------------- #
# Rendering
# --------------------------------------------------------------------------- #

def render_text(rules: list[dict], scans: dict[str, dict]) -> tuple[str, int]:
    lines: list[str] = []
    total = 0
    lines.append("Forbidden-Strings Summary Report")
    lines.append("=" * 34)
    lines.append("")
    for rule in rules:
        rid = rule["id"]
        hits = scans[rid]
        count = sum(len(v) for v in hits.values())
        total += count
        lines.append(f"[{rid}]")
        lines.append(f"  description : {rule.get('description', '')}")
        lines.append(f"  pattern     : /{rule['pattern']}/")
        if rule.get("replacement"):
            lines.append(f"  replacement : {rule['replacement']}")
        else:
            lines.append("  replacement : (none — manual fix required)")
        if not hits:
            lines.append("  status      : ✅ clean (0 findings)")
            lines.append("")
            continue
        lines.append(
            f"  status      : ❌ {count} finding(s) in {len(hits)} file(s)"
        )
        lines.append("  files       :")
        for path in sorted(hits.keys()):
            entries = hits[path]
            uniq = sorted({m for _, _, m in entries})
            lines.append(
                f"    - {path}  ({len(entries)} hit(s); matches: {', '.join(uniq)})"
            )
        cmd = build_fix_command(rule, hits)
        if cmd is None:
            lines.append("  fix command : (no `replacement` declared — fix manually)")
        else:
            lines.append("  fix command :")
            lines.append(f"    {cmd}")
        lines.append("")
    if total == 0:
        lines.append("All rules clean — no findings.")
    else:
        lines.append(f"Total: {total} finding(s) across {len(rules)} rule(s).")
    return "\n".join(lines) + "\n", total


def render_markdown(rules: list[dict], scans: dict[str, dict]) -> tuple[str, int]:
    """Markdown report.

    The persisted summary report is itself scanned by the linter, so we
    deliberately avoid embedding any literal stale string (forbidden brand
    names, legacy module paths, legacy CDN domains, etc.). When a rule
    declares ``summary_description`` / ``summary_pattern`` /
    ``summary_replacement``, those redacted variants are used instead of
    the raw fields. The exact stale matches per file are masked the same
    way, and the fix command is replaced with a pointer to the CLI which
    can emit the runnable command on demand.
    """
    lines: list[str] = ["# Forbidden-Strings Summary Report", ""]
    total = 0
    for rule in rules:
        rid = rule["id"]
        hits = scans[rid]
        count = sum(len(v) for v in hits.values())
        total += count
        desc = rule.get("summary_description") or rule.get("description", "")
        pat = rule.get("summary_pattern") or rule["pattern"]
        repl = rule.get("summary_replacement") or rule.get("replacement")
        lines.append(f"## `{rid}`")
        lines.append("")
        lines.append(f"- **Description:** {desc}")
        lines.append(f"- **Pattern:** `{pat}`")
        if repl:
            lines.append(f"- **Canonical replacement:** `{repl}`")
        if not hits:
            lines.append("- **Status:** ✅ clean (0 findings)")
            lines.append("")
            continue
        lines.append(
            f"- **Status:** ❌ {count} finding(s) in {len(hits)} file(s)"
        )
        lines.append("")
        lines.append("| File | Hits | Stale matches |")
        lines.append("|------|-----:|---------------|")
        masked = "`<redacted>`"
        for path in sorted(hits.keys()):
            entries = hits[path]
            uniq_count = len({m for _, _, m in entries})
            mask_cells = ", ".join([masked] * uniq_count)
            lines.append(
                f"| `{path}` | {len(entries)} | {mask_cells} |"
            )
        lines.append("")
        if rule.get("replacement"):
            lines.append(
                "**Fix command:** run "
                f"`python3 linter-scripts/forbidden-strings-summary.py "
                f"--rule {rid} --emit-fix-command | bash`"
            )
        else:
            lines.append("_No `replacement` declared — fix manually._")
        lines.append("")
    if total == 0:
        lines.append("✅ All rules clean — no findings.")
    else:
        lines.append(f"💥 Total: **{total}** finding(s) across {len(rules)} rule(s).")
    return "\n".join(lines) + "\n", total


# --------------------------------------------------------------------------- #
# CLI
# --------------------------------------------------------------------------- #

def parse_args() -> argparse.Namespace:
    p = argparse.ArgumentParser(
        prog="forbidden-strings-summary",
        description=(
            "Per-rule digest of forbidden-strings findings with an "
            "exact, runnable replacement command."
        ),
    )
    p.add_argument("--rule", help="Limit report to a single rule id.")
    p.add_argument(
        "--markdown", action="store_true",
        help="Render Markdown (suitable for PR comments / job summary).",
    )
    p.add_argument(
        "--emit-fix-command", action="store_true",
        help=(
            "Emit only the bash fix command (no headers). Combine with "
            "--rule to scope to one rule. Pipes safely to `| bash`."
        ),
    )
    return p.parse_args()


def select_rules(all_rules: list[dict], rule_id: str | None) -> list[dict]:
    if rule_id is None:
        return all_rules
    for r in all_rules:
        if r.get("id") == rule_id:
            return [r]
    available = ", ".join(r.get("id", "?") for r in all_rules) or "(none)"
    sys.exit(
        f"Error: rule '{rule_id}' not found in TOML. Available: {available}"
    )


def main() -> int:
    import sys
    if hasattr(sys.stdout, 'reconfigure'):
        sys.stdout.reconfigure(encoding='utf-8')

    args = parse_args()
    all_rules = load_rules()
    if not all_rules:
        sys.exit("Error: no [[rule]] entries in TOML.")

    rules = select_rules(all_rules, args.rule)
    scans: dict[str, dict] = {r["id"]: scan_rule(r) for r in rules}

    if args.emit_fix_command:
        if args.rule is None:
            sys.exit(
                "Error: --emit-fix-command requires --rule (one command "
                "per rule)."
            )
        rule = rules[0]
        cmd = build_fix_command(rule, scans[rule["id"]])
        if cmd is None:
            sys.stdout.write(
                f": no `replacement` declared for {rule['id']} — fix manually\n"
            )
            return 0
        sys.stdout.write(cmd + "\n")
        return 0

    renderer = render_markdown if args.markdown else render_text
    output, total = renderer(rules, scans)
    sys.stdout.write(output)

    # If running inside GitHub Actions, also append to the job summary
    # so reviewers see the report inline on the PR without opening logs.
    summary_path = os.environ.get("GITHUB_STEP_SUMMARY")
    if summary_path:
        try:
            with open(summary_path, "a", encoding="utf-8") as fh:
                md, _ = render_markdown(rules, scans)
                fh.write("\n")
                fh.write(md)
        except OSError as exc:
            print(f"Error: {exc}", file=sys.stderr)
            pass  # non-fatal

    return 1 if total > 0 else 0


if __name__ == "__main__":
    sys.exit(main())
