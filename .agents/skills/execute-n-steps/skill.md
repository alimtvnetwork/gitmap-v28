---
name: execute-n-steps
description: Formats input text verbatim into high-priority instructions, action items, additional learning and planning directives, and execute-parent-task-with-n-steps-v6 skill invocation suffix.
---

# Letterly Execute N Steps Mode (Alias)

Source prompt: `01-prompts/22-letterly/03-execute-n-steps.md`

## Instructions

1. Read `01-prompts/22-letterly/03-execute-n-steps.md` in full before doing the task.
2. Follow `letterly-execute-n-steps` specifications.
3. Output `# High Priority Instruction` followed by `${Input Text Verbatim}`.
4. Output `# Actionable Items Must Follow Non-Negotiable` with numbered steps.
5. Append `Must follow and spawn agent using` and `[06-execute-parent-task-with-n-steps-v6.md](file;.agents/skills/execute-parent-task-with-n-steps-v6)`.
6. Include `learn [/learn](slashCommand;learn) if you have to learn something and [/plan](slashCommand;plan) stuff before working please.`
