#!/usr/bin/env python3
"""
45-antigravity-run-audit.py - Audit an Antigravity CLI (agy) run against the V4 execute rules.

Reads the NDJSON stream written by `agy -p ... --output-format stream-json` and the git
state of the workspace the run used, then prints a PASS/FAIL table and a Markdown report.

Usage:
  python 03-ai-scripts/45-antigravity-run-audit.py --log <run.ndjson> --repo <workspace> --base <commit>
  python 03-ai-scripts/45-antigravity-run-audit.py --log <run.ndjson> --repo <workspace> --base <commit> --report <out.md>
"""

import argparse
import collections
import json
from pathlib import Path
import re
import subprocess
import sys

if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")

ENCODING = "utf-8"
BOM = "\ufeff"
EVENT_STEP = "step_update"
EVENT_RESULT = "result"
TOOL_STEP_TYPES = frozenset({"tool", "subagent"})
STATE_DONE = "DONE"
STATUS_SUCCESS = "SUCCESS"
SUBAGENT_TOOL_NAME = "invoke_subagent"
AGENT_STATE_MARKER = "\\.gemini\\antigravity"
FILE_URI_PREFIX = "file:///"
STATUS_ERROR = "ERROR"
LEDGER_GLOB = ".ai-memory/temp-agents/*/ledger.md"
ABSOLUTE_PATH_PATTERN = re.compile(r"[A-Za-z]:\\[^\"'\s]+")
BANNED_COMMANDS = {
    "whole-tree staging": re.compile(r"git\s+add\s+(-A|--all|\.(\s|$))"),
    "push": re.compile(r"git\s+push|gitmap\s+cp[fbr]"),
    "build": re.compile(r"npm\s+run\s+build|go\s+build|vite\s+build|bun\s+run\s+build"),
    "test suite": re.compile(r"go\s+test|pytest|vitest|npm\s+(run\s+)?test|--run-tests"),
    "go generate": re.compile(r"go\s+generate"),
}
PASS_LABEL = "PASS"
FAIL_LABEL = "FAIL"


def read_events(log_path: Path) -> list:
    """Returns every JSON event in the stream log, skipping non-JSON lines."""
    lines = log_path.read_text(encoding=ENCODING, errors="replace").splitlines()
    cleaned = [line.lstrip(BOM) for line in lines]

    return [json.loads(line) for line in cleaned if line.startswith("{")]


def done_tool_steps(events: list) -> list:
    """Returns finished tool steps from the stream."""
    steps = [event[EVENT_STEP] for event in events if EVENT_STEP in event]
    is_done_tool = lambda step: step.get("step_type") in TOOL_STEP_TYPES and step.get("state") == STATE_DONE

    return [step for step in steps if is_done_tool(step)]


def final_result(events: list) -> dict:
    """Returns the final result event payload, or an empty dict when the run did not finish."""
    results = [event[EVENT_RESULT] for event in events if EVENT_RESULT in event]

    return results[-1] if results else {}


def parameter_text(step: dict) -> str:
    """Flattens a tool step's parameters into one searchable string."""
    parameters = step.get("tool_info", {}).get("parameters", {})

    return json.dumps(parameters, ensure_ascii=False)


def banned_command_hits(steps: list) -> dict:
    """Maps each banned command category to the parameter strings that matched it."""
    texts = [parameter_text(step) for step in steps]
    hits = {name: [text for text in texts if pattern.search(text)] for name, pattern in BANNED_COMMANDS.items()}

    return {name: matches for name, matches in hits.items() if matches}


def outside_workspace_paths(steps: list, workspace: Path) -> list:
    """Returns absolute paths in tool parameters outside the workspace and Antigravity's own state folder."""
    root = str(workspace.resolve()).lower()
    paths = [path for step in steps for path in ABSOLUTE_PATH_PATTERN.findall(parameter_text(step))]
    normalized = [path.replace("\\\\", "\\").rstrip("\\") for path in paths]
    is_allowed = lambda path: path.lower().startswith(root) or AGENT_STATE_MARKER in path.lower()

    return sorted({path for path in normalized if is_allowed(path) is False})


def worker_logs(events: list) -> list:
    """Returns the transcript paths of every subagent the run launched."""
    infos = [event[EVENT_STEP].get("subagent_info") or {} for event in events if EVENT_STEP in event]
    uris = {agent.get("log_uri", "") for info in infos for agent in info.get("subagents", [])}

    return sorted(Path(uri.removeprefix(FILE_URI_PREFIX)) for uri in uris if uri.startswith(FILE_URI_PREFIX))


def worker_errors(events: list) -> dict:
    """Maps each subagent transcript name to its failed steps' error messages."""
    logs = [path for path in worker_logs(events) if path.is_file()]
    records = {path.parent.parent.parent.name: read_events(path) for path in logs}

    return {name: [row.get("error", "") for row in rows if row.get("status") == STATUS_ERROR] for name, rows in records.items()}


def git_lines(repo: Path, *args: str) -> list:
    """Runs a git command in the workspace and returns its non-empty output lines."""
    result = subprocess.run(["git", "-C", str(repo), *args], capture_output=True, text=True, encoding=ENCODING)

    return [line for line in result.stdout.splitlines() if line.strip()]


def collect_facts(events: list, repo: Path, base: str) -> dict:
    """Gathers every fact the checks and the report need."""
    steps = done_tool_steps(events)
    result = final_result(events)
    tool_counts = collections.Counter(step.get("tool_name", "unknown") for step in steps)

    return {
        "result": result,
        "tool_counts": tool_counts,
        "subagent_calls": tool_counts[SUBAGENT_TOOL_NAME],
        "worker_errors": worker_errors(events),
        "banned": banned_command_hits(steps),
        "outside_paths": outside_workspace_paths(steps, repo),
        "commits": git_lines(repo, "log", "--oneline", f"{base}..HEAD"),
        "changed": git_lines(repo, "diff", "--name-status", f"{base}..HEAD"),
        "uncommitted": git_lines(repo, "status", "--porcelain"),
        "ledgers": sorted(str(path.relative_to(repo)) for path in repo.glob(LEDGER_GLOB)),
    }


def build_checks(facts: dict) -> list:
    """Returns (label, check name, detail) rows for the V4 rules."""
    result = facts["result"]
    is_success = result.get("status") == STATUS_SUCCESS
    commit_count = len(facts["commits"])
    error_count = sum(len(errors) for errors in facts["worker_errors"].values())

    return [
        (PASS_LABEL if is_success else FAIL_LABEL, "Run finished", result.get("status", "no result event")),
        (PASS_LABEL if facts["subagent_calls"] else FAIL_LABEL, "Subagents used (R5)", f"{facts['subagent_calls']} calls"),
        (FAIL_LABEL if error_count else PASS_LABEL, "Worker steps without errors", f"{error_count} failed steps in {len(facts['worker_errors'])} workers"),
        (FAIL_LABEL if facts["banned"] else PASS_LABEL, "No banned commands (R1, R8, R10, R15)", ", ".join(facts["banned"]) or "none"),
        (FAIL_LABEL if facts["outside_paths"] else PASS_LABEL, "No paths outside the workspace", f"{len(facts['outside_paths'])} found"),
        (PASS_LABEL if facts["ledgers"] else FAIL_LABEL, "Ledger written", ", ".join(facts["ledgers"]) or "none"),
        (PASS_LABEL if commit_count == 1 else FAIL_LABEL, "Exactly one commit (R9)", f"{commit_count} commits"),
        (FAIL_LABEL if facts["uncommitted"] else PASS_LABEL, "Nothing left uncommitted", f"{len(facts['uncommitted'])} entries"),
    ]


def render_report(facts: dict, checks: list) -> str:
    """Renders the Markdown report."""
    usage = facts["result"].get("usage", {})
    lines = ["# Antigravity Run Audit", "", "| Result | Check | Detail |", "|---|---|---|"]
    lines += [f"| {label} | {name} | {detail} |" for label, name, detail in checks]
    lines += ["", f"- Duration: {facts['result'].get('duration_seconds', 'n/a')} s", f"- Tokens: {usage.get('total_tokens', 'n/a')}"]
    lines += ["", "## Tool calls", ""] + [f"- `{name}`: {count}" for name, count in facts["tool_counts"].most_common()]
    lines += ["", "## Commits", ""] + [f"- {line}" for line in facts["commits"] or ["none"]]
    lines += ["", "## Files changed", ""] + [f"- `{line}`" for line in facts["changed"] or ["none"]]
    lines += ["", "## Uncommitted", ""] + [f"- `{line}`" for line in facts["uncommitted"] or ["none"]]
    lines += ["", "## Banned command hits", ""] + [f"- {name}: `{text[:200]}`" for name, texts in facts["banned"].items() for text in texts]
    lines += ["", "## Paths outside the workspace", ""] + [f"- `{path}`" for path in facts["outside_paths"] or ["none"]]
    lines += ["", "## Worker errors", ""] + [f"- `{name}`: {len(errors)} x `{sorted(set(errors))[:3]}`" for name, errors in facts["worker_errors"].items() if errors]

    return "\n".join(lines) + "\n"


def parse_arguments() -> argparse.Namespace:
    """Configures the command-line interface."""
    parser = argparse.ArgumentParser(description="Audit an agy stream-json run against the V4 execute rules.")
    parser.add_argument("--log", required=True, type=Path, help="NDJSON file written by agy --output-format stream-json")
    parser.add_argument("--repo", required=True, type=Path, help="Workspace the run used")
    parser.add_argument("--base", required=True, help="Commit the workspace was at before the run")
    parser.add_argument("--report", type=Path, help="Optional path for the Markdown report")

    return parser.parse_args()


def main() -> int:
    """Runs the audit and returns 0 when every check passes, 1 otherwise."""
    args = parse_arguments()
    facts = collect_facts(read_events(args.log), args.repo, args.base)
    checks = build_checks(facts)
    report = render_report(facts, checks)
    print(report)

    if args.report:
        args.report.write_text(report, encoding=ENCODING)

    has_failure = any(label == FAIL_LABEL for label, _, _ in checks)

    return 1 if has_failure else 0


if __name__ == "__main__":
    sys.exit(main())
