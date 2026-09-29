#!/usr/bin/env python3
"""Upgrade v4 prompts (13-plan-audit, 14-execute, 15-cg-execute) and matching skills
to the unified v4 format:
1. Starts directly with [/goal](slashCommand:goal) and [/learn](slashCommand:learn) (no # header at top).
2. Enforces strict no-build and no-test mandate in goal, subagents A=2 H=2, GitMap AUM primary, continuous self-looping.
3. Enforces Bottom-Instruction Priority Mandate (Below Precedence) and '--' bottom border.
4. Injects expanded High-Speed GitMap Acceleration options across below-steps and execute prompts.
5. Synchronizes new CG Execute prompts (29..32) and updates readme.md / .ai-memory/prompts.md.
"""
from __future__ import annotations

from pathlib import Path
import re
import shutil

ROOT = Path(__file__).resolve().parent.parent
V4_DIR = ROOT / "01-prompts" / "v4"
V3_DIR = ROOT / "01-prompts" / "v3"
SKILLS_DIR = ROOT / ".agents" / "skills"

GITMAP_FAST_BLOCK = """#### High-Speed GitMap Acceleration Options (Run Everything Faster)

Always prefer native GitMap commands over slow generic shell pipelines:
1. **Ultra-Fast File & Directory Discovery (AUM Index & Walk):**
   - **Wildcard / Glob Search:** `gitmap find "<wildcard*>" [-ext <ext>]` (alias `gitmap f`)
   - **Exact Filename Search:** `gitmap find-files <name> [-ext <ext>]` (alias `gitmap ff`)
   - **Substring Filename Search:** `gitmap find-files-any "<str>" [-ext <ext>]` (alias `gitmap ffa`)
   - **Prefix / Suffix Search:** `gitmap find-files-startswith <prefix>` (`gitmap ffs`) / `gitmap find-files-endswith <suffix>` (`gitmap ffe`)
   - **List Indexed Repo Files:** `gitmap list-files [pattern] [-ext <ext>]` (alias `gitmap lf`)
   - **Directory Tree & Scaffolding:** `gitmap folder-tree` (alias `gitmap ft`)
   - **Zero-Write File Stream:** `gitmap cat <filepath>`
   - **Instant Multi-Core Regex Search:** `gitmap search "<term>"` or `gitmap aum search "<query>" [dir] --ext <ext>`
2. **Fast Repository Hygiene, Lowercase & Symlink Repair:**
   - **Auto-Lowercase Files (Safe 2-Step `git mv`):** `gitmap lowercase` (alias `gitmap lcf [--dry-run]`)
   - **Lowercase Root Readme:** `gitmap lowercase-readme`
   - **Sync Curated `.gitignore` / `.gitattributes` / `.prettierignore`:** `gitmap commons` (alias `gitmap co` or `gitmap sync all`)
   - **Repair Broken Symlinks:** `gitmap fix-link` (alias `gitmap fixlink`)
   - **Clean Update Temp & Inspect Storage:** `gitmap update-cleanup`, `gitmap storage` (alias `gitmap stor`)
3. **Fast Git State, Execution & Atomic Commits:**
   - **Repo Status & Remote Check:** `gitmap status` (`gitmap st`), `gitmap has-any-updates` (`gitmap hau`), `gitmap latest-branch` (`gitmap lb`)
   - **Fast Cross-Platform Shell Runner:** `gitmap pwsh "<command>"` (`gitmap ps`), `gitmap bash "<command>"` (`gitmap sh`), `gitmap async <cmd>` (`gitmap asyn`)
   - **Semantic Atomic Commit & Push:** `gitmap cpf "<summary>"` (Feature), `gitmap cpb "<summary>"` (Bug), `gitmap cpr "<summary>"` (Release), `gitmap pcp "<summary>"` (Pull-Commit-Push)
   - **Smart CI/CD Pipeline Waiting:** `gitmap pe`, `gitmap pipeline-ai status --json` (`gitmap pl-ai status -t <etaSeconds>`)
"""

BOTTOM_BLOCK = """--

## 🚨 Highest Priority Instructions (Appended User Tasks & Instructions Below)

[PASTE USER REQUEST / TASK INSTRUCTIONS HERE — THE AGENT MUST EXECUTE WHATEVER IS WRITTEN BELOW WITH ABSOLUTE PRIORITY AND PRECEDENCE OVER ALL GENERAL GUIDELINES ABOVE]
"""


def extract_goal_line(text: str, fallback_title: str) -> str:
    """Extract existing goal description or synthesize from title."""
    m = re.search(r"\[/goal\]\(slashCommand:goal\)\s*(.+)", text)
    if m:
        raw = m.group(1).strip()
        # Strip trailing period
        raw = raw.rstrip(".")
        if "no-build" not in raw.lower() and "never run build" not in raw.lower():
            raw += (
                " with strict no-build and no-test execution (NEVER run build commands like `go build` "
                "or `npm run build`, and NEVER run test suites like `go test ./...` or `pytest` during "
                "routine execution turns; all compilation and testing are strictly verified later in CI/CD). "
                "Spawn autonomous subagents (A = 2, H = 2) for parallel reading and modular spec generation, "
                "use GitMap high-speed commands as primary, establish a single-agent blueprint in Phase 1 "
                "(first 50% steps budget), and execute disjoint refactoring in Phase 2 (remaining 50% steps budget) "
                "with continuous self-looping until 100% complete and finalized with an atomic push."
            )
        else:
            raw += "."
        return raw

    return (
        f"Autonomously orchestrate and execute {fallback_title} across the codebase in bounded 5-8 file "
        "micro-batches with strict no-build and no-test execution (NEVER run build commands like `go build` "
        "or `npm run build`, and NEVER run test suites like `go test ./...` or `pytest` during routine "
        "execution turns; all compilation and testing are strictly verified later in CI/CD). Spawn autonomous "
        "subagents (A = 2, H = 2) for parallel reading and modular spec generation, use GitMap high-speed "
        "commands as primary, establish a single-agent blueprint in Phase 1 (first 50% steps budget), and "
        "execute disjoint refactoring in Phase 2 (remaining 50% steps budget) with continuous self-looping "
        "until 100% complete and finalized with an atomic push."
    )


def strip_old_header(text: str) -> tuple[str, str]:
    """Return (title, body_starting_from_first_major_section)."""
    lines = text.splitlines()
    title = "the task"
    if lines and lines[0].startswith("# "):
        title = lines[0][2:].split("—")[0].strip()

    # Find the first major section marker (e.g. '## Phase 0', '### Master Task Checklist', '## Variables', '## RULE 0', '## The Unified')
    cut_idx = 0
    for i, line in enumerate(lines):
        if (
            line.startswith("## Phase 0")
            or line.startswith("### Master Task Checklist")
            or line.startswith("## The Unified")
            or line.startswith("## Variables")
            or line.startswith("## RULE 0")
            or line.startswith("## Phase 1")
            or line.startswith("## 1.")
        ):
            cut_idx = i
            break

    if cut_idx > 0:
        # Check if line before cut_idx is '---'
        body = "\n".join(lines[cut_idx:])
    else:
        body = text

    return title, body


def upgrade_prompt_file(path: Path) -> None:
    """Transform a prompt file in v4 into the clean v4 goal-first, bottom-priority structure."""
    if path.name == "readme.md":
        return

    text = path.read_text(encoding="utf-8")

    # Skip rewriting the header of 29, 30, 31, 32 in 15-cg-execute as they were freshly authored in v4 format
    if path.name in {
        "29-code-dryness-and-library-extraction.md",
        "30-clean-repo-build-and-caches.md",
        "31-cg-execute-in-below-steps.md",
        "32-cg-follow-other-prompts.md",
    }:
        return

    # If it is 09-parent-task-in-below-steps.md, just ensure GitMap fast block is upgraded
    if path.name == "09-parent-task-in-below-steps.md":
        if "High-Speed GitMap Acceleration Options" not in text:
            text = text.replace(
                "#### GitMap High-Efficiency Commands & Tooling (Search, Script Runner, Commit)",
                GITMAP_FAST_BLOCK + "\n#### GitMap High-Efficiency Commands & Tooling (Search, Script Runner, Commit)",
            )
            path.write_text(text, encoding="utf-8")
        return

    title, body = strip_old_header(text)
    goal_desc = extract_goal_line(text, title)

    header = f"""[/goal](slashCommand:goal) {goal_desc}

[/learn](slashCommand:learn) Enforce the Bottom-Instruction Priority Mandate: whatever directives, custom requirements, checklists, or user instructions are appended BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and strictly NON-NEGOTIABLE. Ingest and execute the bottom instructions ahead of all default conventions above, internalize the single-agent architectural blueprint, and persist all progress into `.ai-memory/plans/` and memory logs.

> **Bottom-Instruction Priority Mandate (Below Precedence / Suffix Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or guidelines above. The agent MUST inspect and follow the instructions below with absolute precedence.

```text
N = 200 (Total self-loop steps budget)
A = 2   (Number of spawned autonomous subagents running concurrently, default: 2)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)

System Concurrency Capacity = A × H = 2 agents × 2 hands = 4 concurrent subtask operations
```

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Planning, Detailed Spec, and Lean Subtask Generation)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Parallel Execution, Self-Looping, Targeted Quality Linting)
```

N, A, H, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

{GITMAP_FAST_BLOCK}
---

"""

    # Replace Top-Instruction references in body with Bottom-Instruction
    body = body.replace(
        "Top-Instruction Priority Mandate (Preamble Precedence)",
        "Bottom-Instruction Priority Mandate (Below Precedence)",
    )
    body = body.replace(
        "TOP-INSTRUCTION PRIORITY MANDATE (PREAMBLE PRECEDENCE)",
        "BOTTOM-INSTRUCTION PRIORITY MANDATE (BELOW PRECEDENCE)",
    )

    # Ensure bottom block exists
    if "\n--\n" not in body and not body.rstrip().endswith("[PASTE USER REQUEST / TASK INSTRUCTIONS HERE — THE AGENT MUST EXECUTE WHATEVER IS WRITTEN BELOW WITH ABSOLUTE PRIORITY AND PRECEDENCE OVER ALL GENERAL GUIDELINES ABOVE]"):
        body = body.rstrip() + "\n\n" + BOTTOM_BLOCK

    new_content = header + body
    path.write_text(new_content, encoding="utf-8")


def main() -> None:
    # 1. Upgrade all prompts in v4/13-plan-audit, v4/14-execute, v4/15-cg-execute
    target_folders = ["13-plan-audit", "14-execute", "15-cg-execute"]
    upgraded_count = 0
    for folder in target_folders:
        fdir = V4_DIR / folder
        if not fdir.exists():
            continue
        for md_file in sorted(fdir.glob("*.md")):
            if md_file.name == "readme.md":
                continue
            upgrade_prompt_file(md_file)
            upgraded_count += 1

    # 2. Also update v3/14-execute/09-parent-task-in-below-steps.md with expanded GitMap options
    v3_below = V3_DIR / "14-execute" / "09-parent-task-in-below-steps.md"
    v4_below = V4_DIR / "14-execute" / "09-parent-task-in-below-steps.md"
    if v4_below.exists():
        shutil.copy2(v4_below, v3_below)

    # 3. Copy new CG execute prompts (29..32) into v3/15-cg-execute and root 01-prompts/15-cg-execute
    for fname in [
        "29-code-dryness-and-library-extraction.md",
        "30-clean-repo-build-and-caches.md",
        "31-cg-execute-in-below-steps.md",
        "32-cg-follow-other-prompts.md",
    ]:
        src_f = V4_DIR / "15-cg-execute" / fname
        for dest_dir in [V3_DIR / "15-cg-execute", ROOT / "01-prompts" / "15-cg-execute"]:
            if dest_dir.exists():
                shutil.copy2(src_f, dest_dir / fname)

    print(f"Upgraded {upgraded_count} prompt files in 01-prompts/v4/ and synced new CG execute prompts.")


if __name__ == "__main__":
    main()
