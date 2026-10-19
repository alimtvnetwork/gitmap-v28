#!/usr/bin/env python3
"""Upgrade all 38 prompts in 01-prompts/15-cg-execute/ and 01-prompts/v4/15-cg-execute/.

Enforces:
- 6-Part Task Verification Schema in /goal and Phase 1A narrative
- Task Comprehension & Model Runtime Invariants in Phase 1A Pillars
- Understood / AI Comprehension field injection in Task-01 and Task-02
- Azure AI / Model Runtime Context field injection in Task-01 and Task-02
- Preservation of Prompt 31 bottom-precedence mandate
- 1:1 Mirroring to 01-prompts/v4/15-cg-execute/ with byte verification
- Strict relative Git paths, lowercase naming, and LF newline hygiene
"""

from __future__ import annotations

import argparse
from pathlib import Path
import sys

PROMPTS_DIR = Path("01-prompts/15-cg-execute")
V4_PROMPTS_DIR = Path("01-prompts/v4/15-cg-execute")

GOAL_SEARCH = "with the 5 Mandatory Pillars and 4-part verification schema"
GOAL_REPLACE = "with the 5 Mandatory Pillars and 6-part verification schema"

PHASE_1A_SEARCH = "using the 4-part verification schema:"
PHASE_1A_REPLACE = (
    "using the 6-part verification schema (mandating explicit "
    "**Understood / AI Comprehension** and **Azure AI / Model Runtime Context** "
    "fields for every task):"
)

PILLAR_5_SEARCH = (
    "   * **Pillar 5: Multi-Agent Concurrency Topology:** Detail swarm "
    "structure (`A = 3`, `H = 2`, Capacity = 6), disjoint file assignments, "
    "and immediate dynamic respawn."
)

INVARIANTS_INSERT = (
    "   * **Task Comprehension & Model Runtime Invariants:** In the task breakdown, "
    "every single task MUST explicitly state `- **Understood / AI Comprehension:** "
    "[YES] — [1–2 concise sentences proving deep understanding of intent, scope, "
    "and verified constraints; answering did the AI understand the requirements or not]` "
    "and `- **Azure AI / Model Runtime Context:** [Explicitly document Azure AI, "
    "Azure OpenAI, or local model endpoint context if applicable, or state "
    "'N/A — Local CLI/GitMap']`."
)

TASK_01_SEARCH = (
    "1. **Task-01: [Descriptive Task Title]**\n"
    "   - **State:** `[IN PROGRESS — EXECUTING IMMEDIATELY]`"
)

TASK_01_REPLACE = (
    "1. **Task-01: [Descriptive Task Title]**\n"
    "   - **State:** `[IN PROGRESS — EXECUTING IMMEDIATELY]`\n"
    "   - **Understood / AI Comprehension:** [YES] — [1–2 concise sentences proving deep understanding of intent, scope, and verified constraints; answering did the AI understand the requirements or not]\n"
    "   - **Azure AI / Model Runtime Context:** [Explicitly document Azure AI, Azure OpenAI, or local model endpoint context if applicable, or state 'N/A — Local CLI/GitMap']"
)

TASK_02_SEARCH = (
    "2. **Task-02: [Descriptive Task Title]**\n"
    "   - **State:** `[QUEUED — EXECUTING NOW WITHOUT USER PROMPT]`"
)

TASK_02_REPLACE = (
    "2. **Task-02: [Descriptive Task Title]**\n"
    "   - **State:** `[QUEUED — EXECUTING NOW WITHOUT USER PROMPT]`\n"
    "   - **Understood / AI Comprehension:** [YES] — [1–2 concise sentences proving deep understanding of intent, scope, and verified constraints; answering did the AI understand the requirements or not]\n"
    "   - **Azure AI / Model Runtime Context:** [Explicitly document Azure AI, Azure OpenAI, or local model endpoint context if applicable, or state 'N/A — Local CLI/GitMap']"
)

PROMPT_31_FILENAME = "31-cg-execute-in-below-steps.md"
NORMALIZED_PHASE_1A_HEADER = "### Phase 1A: Verbatim Capture, Task Extraction & Chat Output Gate (Step 0)"
RAW_PHASE_1A_HEADER = "### ### ### Phase 1A: Verbatim Capture, Task Extraction & Chat Output Gate (Step 0)"


def upgrade_prompt_content(filename: str, content: str) -> str:
    """Apply the Understood and Azure AI upgrades to a single prompt."""
    updated = content

    # Normalize Prompt 31 header if repetitive hashes exist
    if filename == PROMPT_31_FILENAME and RAW_PHASE_1A_HEADER in updated:
        updated = updated.replace(RAW_PHASE_1A_HEADER, NORMALIZED_PHASE_1A_HEADER)

    # 1. Update /goal line
    if GOAL_SEARCH in updated:
        updated = updated.replace(GOAL_SEARCH, GOAL_REPLACE)
    else:
        print(f"[WARN] Goal search string not found in {filename}")

    # 2. Update Phase 1A Step 2 intro
    if PHASE_1A_SEARCH in updated:
        updated = updated.replace(PHASE_1A_SEARCH, PHASE_1A_REPLACE)
    else:
        print(f"[WARN] Phase 1A search string not found in {filename}")

    # 3. Add Task Comprehension & Model Runtime Invariants after Pillar 5
    if PILLAR_5_SEARCH in updated:
        replacement_block = f"{PILLAR_5_SEARCH}\n{INVARIANTS_INSERT}"
        updated = updated.replace(PILLAR_5_SEARCH, replacement_block)
    else:
        print(f"[WARN] Pillar 5 search string not found in {filename}")

    # 4. Inject into Task-01
    if TASK_01_SEARCH in updated:
        updated = updated.replace(TASK_01_SEARCH, TASK_01_REPLACE)
    else:
        print(f"[WARN] Task-01 search string not found in {filename}")

    # 5. Inject into Task-02 (Group 1 prompts only)
    if filename != PROMPT_31_FILENAME:
        if TASK_02_SEARCH in updated:
            updated = updated.replace(TASK_02_SEARCH, TASK_02_REPLACE)
        else:
            print(f"[WARN] Task-02 search string not found in {filename}")

    return updated


def main() -> int:
    parser = argparse.ArgumentParser(description="Upgrade CG Execute prompts with Understood & Azure AI fields.")
    parser.add_argument("--dry-run", action="store_true", help="Perform checks without writing files")
    args = parser.parse_args()

    if not PROMPTS_DIR.exists():
        print(f"[ERROR] Prompts directory not found: {PROMPTS_DIR}")
        return 1

    prompt_files = sorted([p for p in PROMPTS_DIR.glob("*.md") if p.name != "readme.md"])
    if len(prompt_files) != 38:
        print(f"[ERROR] Expected 38 prompt files, found {len(prompt_files)}")
        return 1

    print(f"Discovered {len(prompt_files)} prompt files in {PROMPTS_DIR}")

    V4_PROMPTS_DIR.mkdir(parents=True, exist_ok=True)

    upgraded_count = 0
    verification_errors = 0

    for prompt_path in prompt_files:
        filename = prompt_path.name
        original_text = prompt_path.read_text(encoding="utf-8")
        upgraded_text = upgrade_prompt_content(filename, original_text)

        # Quick validation of upgraded text
        if "Understood / AI Comprehension" not in upgraded_text:
            print(f"[ERROR] Validation failed: 'Understood / AI Comprehension' missing in {filename}")
            verification_errors += 1
        if "Azure AI / Model Runtime Context" not in upgraded_text:
            print(f"[ERROR] Validation failed: 'Azure AI / Model Runtime Context' missing in {filename}")
            verification_errors += 1
        if "6-part verification schema" not in upgraded_text:
            print(f"[ERROR] Validation failed: '6-part verification schema' missing in {filename}")
            verification_errors += 1
        if "4-part verification schema" in upgraded_text:
            print(f"[ERROR] Validation failed: Legacy '4-part verification schema' remains in {filename}")
            verification_errors += 1

        if not args.dry_run:
            # Write upgraded file to 01-prompts/15-cg-execute/
            with open(prompt_path, "w", encoding="utf-8", newline="\n") as f:
                f.write(upgraded_text)

            # Mirror to 01-prompts/v4/15-cg-execute/
            v4_path = V4_PROMPTS_DIR / filename
            with open(v4_path, "w", encoding="utf-8", newline="\n") as f:
                f.write(upgraded_text)

            # Byte verification
            orig_bytes = prompt_path.read_bytes()
            v4_bytes = v4_path.read_bytes()
            if orig_bytes != v4_bytes:
                print(f"[ERROR] Byte verification failed for {filename} and {v4_path}")
                verification_errors += 1
            else:
                upgraded_count += 1
        else:
            upgraded_count += 1

    if verification_errors > 0:
        print(f"[FAILED] Completed with {verification_errors} verification error(s).")
        return 1

    mode_label = "Dry-run checked" if args.dry_run else "Successfully upgraded and mirrored"
    print(f"[SUCCESS] {mode_label} all {upgraded_count} prompts across:")
    print(f"  - {PROMPTS_DIR} (38 files)")
    print(f"  - {V4_PROMPTS_DIR} (38 files)")
    return 0


if __name__ == "__main__":
    sys.exit(main())
