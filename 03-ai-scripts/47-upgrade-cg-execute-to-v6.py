#!/usr/bin/env python3
"""Upgrade all 01-prompts/15-cg-execute prompts and companion skills to V6 architecture."""

from pathlib import Path
import re
import sys

PROMPTS_DIR = Path("01-prompts/15-cg-execute")
AGENTS_SKILLS_DIR = Path(".agents/skills")
CURSOR_SKILLS_DIR = Path(".cursor/skills")

V6_META_TEMPLATE = """
> [!IMPORTANT]
> Prompt Version: 6.0.0
> Runtime: Google Antigravity 2.0 (IDE and CLI)
> Invoke: /{invoke_name} <task>
>
> **Top-Instruction Priority Mandate (Above Precedence / Preamble Precedence):**
> Whatever directives, constraints, checklists, or user instructions are given ABOVE this prompt (including in the user preamble, header blocks, or incoming user request above) are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE. They supersede and strictly override any conflicting general advice, default conventions, or lower-level guidelines below. The agent MUST inspect and follow the instructions above with absolute precedence.
"""

WORKER_BOUNDARIES_OLD = re.compile(
    r"### Boundaries:\n- Read any file in the workspace; edit only your Owned Files: <relative paths>\.\n- Code & Symbol Search: Use GitMap exclusively: `gitmap aum search \"<pattern>\" \[dir\] \[-e <\.ext>\] \[-r\] \[-i\]` or `gitmap search`\. TOTAL BAN on PowerShell `Select-String`, `Get-ChildItem -Recurse`, `git grep`, `grep`, or `findstr`\.\n- After C tool calls, stop and report what you have\.\n- A tool failing twice: reply \"STATUS: BLOCKED\" with exact error and stop\. Never guess paths and never troubleshoot machine\.\n- Workers that find a secret stop and report \"BLOCKED: secret at <file>:<line>\"\. They do not handle it themselves\.\n- Adhere to R1, R2, and R11 by ID\."
)

WORKER_BOUNDARIES_NEW = """### Boundaries & Crash Prevention:
- Read any file in the workspace; edit ONLY your Owned Files: <relative paths>.
- TOTAL BAN ON GIT COMMANDS (LOCK COLLISION PREVENTION): NEVER run ANY git commands (`git add`, `git commit`, `git push`, `git status`, `git diff`, `git checkout`). In shared workspaces, worker git calls create `.git/index.lock` collisions that immediately crash parallel agents. Only the lead orchestrator runs git commands after workers complete.
- TOTAL BAN ON COMMITS: Workers NEVER commit, stage, or push. Committing is exclusively reserved for the Lead Agent at Phase 3 via GitMap (`gitmap cpf "<module> - <summary>"` using hyphen `-`; no colon needed in GitMap cpf as colon is already provided).
- Code & Symbol Search: Use GitMap exclusively: `gitmap aum search "<pattern>" [dir] [-e <.ext>] [-r] [-i]` or `gitmap search`. TOTAL BAN on PowerShell `Select-String`, `Get-ChildItem -Recurse`, `git grep`, `grep`, or `findstr`.
- After C tool calls, stop and report what you have.
- A tool failing twice: reply "STATUS: BLOCKED" with exact error and stop. Never guess paths and never troubleshoot machine.
- Workers that find a secret stop and report "BLOCKED: secret at <file>:<line>". They do not handle it themselves.
- Adhere to R1, R2, and R11 by ID."""

SQLITE_STEP_SNIPPET = """
4. **Concurrency-Safe SQLite Action Logging:**
   - Database location: `.ai-memory/temp-agents/<slug>/agent-task.db`
   - Workers log file changes and intermediate milestones:
     `python 03-ai-scripts/46-agent-sqlite-task-manager.py log-action --db <db-path> --subtask-id <id> --agent "Worker NN" --action "<action>" --file "<file>" --details "<details>"`
   - Workers mark completed subtasks upon meeting acceptance criteria:
     `python 03-ai-scripts/46-agent-sqlite-task-manager.py complete --db <db-path> --subtask-id <id> --agent "Worker NN" --evidence "<proof>"`
"""


def upgrade_prompt_file(file_path: Path) -> bool:
    content = file_path.read_text(encoding="utf-8")
    original = content
    stem = file_path.stem
    slug = stem.split("-", 1)[1]

    # 1. Add V6 Metadata Box if missing
    if "Prompt Version: 6.0.0" not in content:
        # Find closing ``` of parameter block
        match = re.search(r"(```text\s*\nN = 300[\s\S]*?```)", content)
        if match:
            param_block = match.group(1)
            v6_meta = V6_META_TEMPLATE.format(invoke_name=slug)
            content = content.replace(param_block, param_block + "\n" + v6_meta.strip() + "\n")

    # 2. Fix slash command semicolons
    content = content.replace("(slashCommand:goal)", "(slashCommand;goal)")
    content = content.replace("(slashCommand:learn)", "(slashCommand;learn)")
    content = content.replace("(slashCommand:plan)", "(slashCommand;plan)")

    # 3. Upgrade worker boundaries
    if "TOTAL BAN ON GIT COMMANDS (LOCK COLLISION PREVENTION)" not in content:
        if WORKER_BOUNDARIES_OLD.search(content):
            content = WORKER_BOUNDARIES_OLD.sub(WORKER_BOUNDARIES_NEW, content)
        elif "### Boundaries:" in content and "LOCK COLLISION" not in content:
            content = content.replace(
                "### Boundaries:\n- Read any file in the workspace; edit only your Owned Files: <relative paths>.",
                WORKER_BOUNDARIES_NEW,
            )

    # 4. Add SQLite Task DB instructions in Phase 2 if missing
    if "46-agent-sqlite-task-manager.py" not in content:
        target_anchor = "### 7.3 Turn-Yielding & Verification Protocol"
        if target_anchor in content:
            replacement = SQLITE_STEP_SNIPPET.strip() + "\n\n" + target_anchor
            content = content.replace(target_anchor, replacement)

    # 5. Fix GitMap commit hyphen convention in Phase 3
    if 'gitmap cpf "<module> - <summary>"' not in content and "gitmap cpf" in content:
        content = content.replace(
            'Call `gitmap cpf "<module>: <feature summary>"`',
            'Call `gitmap cpf "<module> - <feature summary>"` (using hyphen `-` instead of colon `:`)',
        )

    if content != original:
        file_path.write_text(content, encoding="utf-8", newline="\n")
        return True

    return False


def sync_skills() -> int:
    synced = 0
    prompts = sorted([p for p in PROMPTS_DIR.glob("*.md") if p.name != "readme.md"])

    for p in prompts:
        stem = p.stem
        slug = stem.split("-", 1)[1]

        # Candidate skill names
        skill_names = [slug, f"cg-{slug}"]
        if slug.startswith("cg-"):
            skill_names.append(slug[3:])

        for sname in skill_names:
            agent_skill_file = AGENTS_SKILLS_DIR / sname / "skill.md"
            cursor_skill_file = CURSOR_SKILLS_DIR / sname / "skill.md"

            if agent_skill_file.exists():
                text = agent_skill_file.read_text(encoding="utf-8")
                # Ensure cursor has exact copy
                cursor_skill_file.parent.mkdir(parents=True, exist_ok=True)
                cursor_skill_file.write_text(text, encoding="utf-8", newline="\n")
                synced += 1

    return synced


def main():
    prompts = sorted([p for p in PROMPTS_DIR.glob("*.md") if p.name != "readme.md"])
    upgraded_prompts = 0

    for p in prompts:
        if upgrade_prompt_file(p):
            upgraded_prompts += 1
            print(f"Upgraded prompt: {p.name}")

    synced_skills = sync_skills()
    print(f"Total upgraded prompts: {upgraded_prompts}")
    print(f"Total synchronized skills: {synced_skills}")


if __name__ == "__main__":
    main()
