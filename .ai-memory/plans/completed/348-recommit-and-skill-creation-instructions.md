# Plan: 348-recommit-and-skill-creation-instructions

## Status: COMPLETED
## Slug: recommit-and-skill-creation-instructions
## Category: Skills & Prompt Engineering

## User Request (Verbatim)
# Recommit and Skill Creation Instructions: high priority instruction, non-negotiable task

I'm just confirming that can recommit, and also the skill is created, time restriction skill inside the coding guideline repository in the skills folder, but you have not given me the skill path so that I can refer into the letterly prompt. So that time restrict part should be added in the letterly prompt. So that when letterly prompt runs, every AI knows that it needs to enable the time restriction using `gitmap`, and a recommit needs to be happening using `gitmap`. And also there could be another part about the AI automation bot and `gitmap` thing. So the AI can add `gitmap` automation. So that could be another skill. So these two files can be referred, for the `gitmap` to get different set of configuration, if the user wants. That would be on call basis. So that means if the prompt says, then and only then it would do that and then work with it. So apparently, these are not, as improved. I think you should be start working on that. What do you think? So do not modify the existing letterly V6 prompt. I mean, execute parent task in N steps. Don't do that, but create a new one. Do you understand? And no letterly or memes needs to have a skill in the skills folder for it that. And also you should tell me that how we can add the skills to Antigravity by the folder or API, guide me that at the end.

## task: Recommit and Skill Creation Instructions
## slug: recommit-and-skill-creation-instructions

# Actionable Items Must Follow Non-Negotiable
1. Git pull BEFORE starting any work — first try `gitmap pull <repo path>` (or run `gitmap pull` in repository folder); if that fails, fall back to regular `git pull` and resolve any merge conflicts so the repository is completely up to date and clean before starting.
2. Write spec under 02-spec/21-app/XX-recommit-and-skill-creation-instructions/ and enqueue plan task in .ai-memory/plans/pending/XX-recommit-and-skill-creation-instructions.md (subtasks in .ai-memory/plans/subtasks/XX-recommit-and-skill-creation-instructions/) first, (xx should be figured out by the folder structure and running AI tools, xx is sequence 01, 02)
3. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
4. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
5. Create a new skill for time restriction in the coding guideline repository under the skills folder.
6. Ensure the letterly prompt includes instructions for enabling time restriction using `gitmap`.
7. Develop a new skill for AI automation bot integration with `gitmap`.
8. Provide guidance on adding skills to Antigravity via folder or API.

## High-Level Execution Phases
- **Phase 1: Planning & Specifications**
  - Task-01: Discovery across time restriction, recommit, automation bot, and Antigravity APIs. (DONE)
  - Task-02: Author comprehensive specs in `02-spec/21-app/348-recommit-and-skill-creation-instructions/` and subtask plans in `.ai-memory/plans/subtasks/348-recommit-and-skill-creation-instructions/`. (DONE)
- **Phase 2: Core Skills & Prompts Creation**
  - Task-03: Create `.agents/skills/gitmap-time-restriction/skill.md`. (DONE)
  - Task-04: Create `.agents/skills/gitmap-automation-bot/skill.md`. (DONE)
  - Task-05: Create `01-prompts/22-letterly/12-execute-recommit-time-restricted-letterly.md`, `01-prompts/14-execute/09-execute-parent-task-with-n-steps-v7.md`, and update indexes. (DONE)
- **Phase 3: Verification, Antigravity Guide & Release**
  - Task-06: Pre-commit checks, author Antigravity skill integration guidance (folder & API), and atomic release commit. (DONE)

## Subtasks Manifest
1. `01-skills-architecture.md`: Specification of time restriction and automation bot skills. (DONE)
2. `02-letterly-and-execution-prompts.md`: Specification of new Letterly prompt and V7 execution prompt. (DONE)
3. `03-antigravity-integration-and-release.md`: Antigravity skill onboarding guide, verification, and release ceremony. (DONE)
