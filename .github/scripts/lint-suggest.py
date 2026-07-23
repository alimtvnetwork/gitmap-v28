#!/usr/bin/env python3
"""
lint-suggest.py — turn golangci-lint NEW findings into actionable PR
suggestions.

Pairs with lint-diff.py: that script decides which findings are NEW vs
UNCHANGED and gates the build. This script takes the same inputs and
produces a Markdown comment that maps each NEW finding to (a) a
human-readable explanation of the rule, (b) a concrete code-change
template the contributor can apply, and optionally (c) a one-line
`sed`/manual snippet.

Output goes to a file (--out, default GITHUB_STEP_SUMMARY) so the
workflow can both render it inline and post it as a PR comment via
peter-evans/create-or-update-comment.

Why a separate script (and not just lint-diff.py)?
  - lint-diff.py is the gate: it must stay tiny, exit-code driven,
    and free of presentation logic so it stays trivially auditable.
  - This script is presentation: the rule->fix mapping is a living
    table that will grow as the team encounters new linters. Keeping
    it isolated means we can iterate on suggestions without risking
    the gating logic.

Suggestion table coverage matches the linters configured in
.golangci.yml today (gocritic, misspell, nolintlint, unused, gofmt,
revive, errcheck, govet, staticcheck, ineffassign). Unknown linters
fall back to a generic "consult the linter docs" suggestion so the
comment is never blank for a NEW finding.
"""

from __future__ import annotations

import argparse
import json
import os
import re
import sys
from collections import defaultdict
from typing import Callable

Finding = tuple[str, int, str, str]  # (file, line, linter, message)


# Per-linter suggestion builders. Each returns (title, fix_block) where
# fix_block is fenced Markdown (diff/go/text). Builders receive the
# raw message so they can extract identifiers (e.g. the misspelled
# word, the unused symbol) and produce specific guidance.
SuggestionBuilder = Callable[[str], tuple[str, str]]


def main() -> int:
    args = parse_args()
    new_findings = load_new_findings(args.current, args.baseline)

    out_path = args.out or os.environ.get("GITHUB_STEP_SUMMARY", "")
    if not out_path:
        print("::warning::no --out and no GITHUB_STEP_SUMMARY; "
              "writing to stdout", file=sys.stderr)

    markdown = render_markdown(new_findings, repo=args.repo, sha=args.sha)

    if out_path:
        # Append (don't overwrite): GITHUB_STEP_SUMMARY accumulates
        # across steps in the same job, and the PR-comment step reads
        # this same file as its body.
        with open(out_path, "a", encoding="utf-8") as fh:
            fh.write(markdown)
            fh.write("\n")
    else:
        sys.stdout.write(markdown)

    return 0


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--current", required=True,
                        help="Current golangci-lint JSON report")
    parser.add_argument("--baseline", default="",
                        help="Baseline JSON for diffing (optional)")
    parser.add_argument("--out", default="",
                        help="Output Markdown path (default: "
                             "$GITHUB_STEP_SUMMARY)")
    parser.add_argument("--repo", default="",
                        help="owner/repo for permalinks (optional)")
    parser.add_argument("--sha", default="",
                        help="Commit SHA for permalinks (optional)")
    return parser.parse_args()


def load_new_findings(current_path: str,
                      baseline_path: str) -> list[Finding]:
    """Compute the same NEW set lint-diff.py uses, so the comment never
    surfaces a finding that the gate considered pre-existing."""
    current = load_findings(current_path)
    baseline_present = bool(baseline_path) \
        and os.path.exists(baseline_path) \
        and os.path.getsize(baseline_path) > 0
    baseline = load_findings(baseline_path) if baseline_present else set()
    # Without a baseline every finding is "new" from the diff's POV,
    # but we don't want to drown a seeding PR in suggestions. Cap at
    # 50 so the comment stays readable; the rest are still in the
    # uploaded JSON artifact.
    new = sorted(current - baseline)
    if not baseline_present:
        new = new[:50]
    return new


def load_findings(path: str) -> set[Finding]:
    """Mirror of lint-diff.py's loader. Kept duplicated (rather than
    imported) so this script can run from any working directory and
    has no Python-side dependency on its sibling — they're independent
    deployables that happen to read the same JSON shape."""
    if not path or not os.path.exists(path):
        return set()
    if os.path.getsize(path) == 0:
        return set()
    try:
        with open(path, encoding="utf-8") as fh:
            data = json.load(fh)
    except (json.JSONDecodeError, OSError) as err:
        print(f"::warning::could not parse {path}: {err}", file=sys.stderr)
        return set()

    out: set[Finding] = set()
    for issue in data.get("Issues") or []:
        pos = issue.get("Pos") or {}
        file = pos.get("Filename", "")
        line = int(pos.get("Line", 0) or 0)
        linter = issue.get("FromLinter", "")
        message = (issue.get("Text") or "").strip()
        if file and linter:
            out.add((file, line, linter, message))
    return out


# ---------------------------------------------------------------------------
# Suggestion builders
# ---------------------------------------------------------------------------

def suggest_misspell(message: str) -> tuple[str, str]:
    # Message shape: "`colour` is a misspelling of `color`"
    match = re.search(r"`([^`]+)`\s+is a misspelling of\s+`([^`]+)`",
                      message)
    if not match:
        return ("Fix spelling per US English convention.",
                "```text\nReplace the flagged word with the suggested "
                "spelling.\n```")
    wrong, right = match.group(1), match.group(2)
    fix = (f"```diff\n- {wrong}\n+ {right}\n```\n"
           f"_All occurrences in the same file/comment block should be "
           f"updated together._")
    return (f"Replace `{wrong}` with `{right}` (US English).", fix)


def suggest_gocritic(message: str) -> tuple[str, str]:
    if "paramTypeCombine" in message:
        return ("Combine consecutive params of the same type.",
                "```diff\n- func f(a string, b string) {}\n"
                "+ func f(a, b string) {}\n```")
    if "sprintfQuotedString" in message:
        return ("Use `%q` instead of `\"%s\"` for quoted strings.",
                "```diff\n- fmt.Sprintf(`KEY=\"%s\"`, val)\n"
                "+ fmt.Sprintf(`KEY=%q`, val)\n```")
    if "filepathJoin" in message:
        return ("Don't pass strings containing path separators to "
                "`filepath.Join`.",
                "```diff\n- filepath.Join(\"/a/b\", \"c\")\n"
                "+ filepath.Join(filepath.FromSlash(\"/a/b\"), \"c\")\n```")
    if "ifElseChain" in message:
        return ("Rewrite long if/else chains as `switch`.",
                "```go\nswitch {\ncase cond1: ...\ncase cond2: ...\n"
                "default: ...\n}\n```")
    return ("Apply the gocritic checker's suggested rewrite.",
            "```text\nSee https://go-critic.com/overview for the named "
            "checker.\n```")


def suggest_unused(message: str) -> tuple[str, str]:
    # "func `pointEnvAt` is unused"  /  "var `x` is unused"
    match = re.search(r"`([^`]+)`\s+is unused", message)
    name = match.group(1) if match else "<symbol>"
    return (f"Remove the unused symbol `{name}`.",
            f"```diff\n- // entire declaration of {name}\n```\n"
            "_If it's part of a future API, add a `//nolint:unused` "
            "directive with a short justification._")


def suggest_nolintlint(message: str) -> tuple[str, str]:
    return ("Remove the unused `//nolint` directive.",
            "```diff\n- foo() //nolint:gosec // no longer needed\n"
            "+ foo()\n```\n"
            "_The underlying linter no longer flags this line; the "
            "directive became dead weight._")


def suggest_gofmt(_message: str) -> tuple[str, str]:
    return ("Run `gofmt -w` on the file.",
            "```bash\ngofmt -w path/to/file.go\n```")


def suggest_errcheck(_message: str) -> tuple[str, str]:
    return ("Handle the returned error explicitly.",
            "```diff\n- doThing()\n+ if err := doThing(); err != nil "
            "{\n+     return fmt.Errorf(\"do thing: %w\", err)\n+ }\n```\n"
            "_Per project policy: never swallow errors. Log to "
            "`os.Stderr` if return is impossible._")


def suggest_revive(_message: str) -> tuple[str, str]:
    return ("Apply the revive style rule.",
            "```text\nSee revive rule docs: https://revive.run/r\n```")


def suggest_govet(_message: str) -> tuple[str, str]:
    return ("Address the `go vet` warning.",
            "```text\nRun `go vet ./...` locally for the full context.\n```")


def suggest_staticcheck(_message: str) -> tuple[str, str]:
    return ("Apply the staticcheck recommendation.",
            "```text\nLook up the SA/ST/S code at "
            "https://staticcheck.dev/docs/checks/\n```")


def suggest_ineffassign(_message: str) -> tuple[str, str]:
    return ("Remove or use the ineffective assignment.",
            "```diff\n- x := compute()\n- x = override()\n"
            "+ x := override()\n```")


def suggest_generic(_message: str) -> tuple[str, str]:
    return ("Consult the linter's documentation.",
            "```text\nNo template available for this linter yet.\n```")


SUGGESTERS: dict[str, SuggestionBuilder] = {
    "misspell": suggest_misspell,
    "gocritic": suggest_gocritic,
    "unused": suggest_unused,
    "nolintlint": suggest_nolintlint,
    "gofmt": suggest_gofmt,
    "errcheck": suggest_errcheck,
    "revive": suggest_revive,
    "govet": suggest_govet,
    "staticcheck": suggest_staticcheck,
    "ineffassign": suggest_ineffassign,
}


def build_suggestion(linter: str, message: str) -> tuple[str, str]:
    return SUGGESTERS.get(linter, suggest_generic)(message)


# ---------------------------------------------------------------------------
# Markdown rendering
# ---------------------------------------------------------------------------

# Sentinel marker so the PR-comment action can find and replace the
# previous comment in-place instead of stacking duplicates on every push.
COMMENT_MARKER = "<!-- gitmap-lint-suggestions -->"


def render_markdown(findings: list[Finding], repo: str,
                    sha: str) -> str:
    if not findings:
        return (f"{COMMENT_MARKER}\n"
                "### golangci-lint — no new findings\n\n"
                "All lint warnings on this PR were already present on "
                "`main`. Nothing to fix here. ✅\n")

    lines: list[str] = [
        COMMENT_MARKER,
        "### golangci-lint — actionable suggestions",
        "",
        f"Found **{len(findings)}** new finding(s) introduced by this "
        "change. Each entry below maps the warning to a concrete fix "
        "you can apply.",
        "",
    ]

    grouped: dict[str, list[Finding]] = defaultdict(list)
    for f in findings:
        grouped[f[0]].append(f)

    for file in sorted(grouped):
        lines.append(f"#### `{file}`")
        lines.append("")
        for _, line, linter, message in grouped[file]:
            title, fix = build_suggestion(linter, message)
            location = format_location(file, line, repo, sha)
            lines.extend([
                f"- **L{line}** · `[{linter}]` — {escape_md(message)}",
                f"  - 📍 {location}",
                f"  - 💡 **Suggested fix:** {title}",
                "",
                indent(fix, "    "),
                "",
            ])
        lines.append("")

    lines.append("---")
    lines.append("_Suggestions are templates, not patches — review before "
                 "applying. The full JSON report is attached as the "
                 "`golangci-lint-report` artifact._")
    return "\n".join(lines)


def format_location(file: str, line: int, repo: str, sha: str) -> str:
    if repo and sha:
        return (f"[`{file}:{line}`]"
                f"(https://github.com/{repo}/blob/{sha}/{file}#L{line})")
    return f"`{file}:{line}`"


def escape_md(text: str) -> str:
    # Minimal escape: only neutralize backticks-inside-message that
    # would break our inline-code formatting. Pipes don't matter
    # because we render bullets, not tables.
    return text.replace("|", "\\|")


def indent(text: str, prefix: str) -> str:
    return "\n".join(prefix + line if line else line
                     for line in text.splitlines())


if __name__ == "__main__":
    sys.exit(main())
