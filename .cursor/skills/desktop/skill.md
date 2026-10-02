---
name: desktop
description: Formats input text verbatim into desktop high-priority instructions, non-negotiable action items, and execute-parent-task-with-n-steps-v6 skill invocation suffix without conversational boilerplate.
---

# Letterly Desktop Mode (Alias)

Source prompt: `01-prompts/22-letterly/02-desktop.md`

## Instructions

1. Read `01-prompts/22-letterly/02-desktop.md` in full before doing the task.
2. Follow `letterly-desktop` specifications.
3. Output `# High Priority Instruction` followed by `[/goal](slashCommand;goal) [/learn](slashCommand;learn) ${Input Text Verbatim}`.
4. Output `# Actionable Items Must Follow Non-Negotiable` with numbered steps.
5. Append `Must follow and spawn an agent using the following skill` and `[06-execute-parent-task-with-n-steps-v6.md](file;.agents/skills/execute-parent-task-with-n-steps-v6)`.
