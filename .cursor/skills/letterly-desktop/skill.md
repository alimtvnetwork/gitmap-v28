---
name: letterly-desktop
description: Formats input text verbatim into desktop high-priority instructions, non-negotiable action items, and execute-parent-task-with-n-steps-v6 skill invocation suffix without conversational boilerplate.
---

# Letterly Desktop Mode

> **[/goal](slashCommand;goal)** Formats input text verbatim into desktop high-priority instructions, non-negotiable action items, and `execute-parent-task-with-n-steps-v6` skill invocation suffix without conversational boilerplate.
> **[/learn](slashCommand;learn)** Ingest user input verbatim, structure into High Priority Instruction followed by Actionable Items, and append the V6 execution skill reference as a suffix.

**Source prompt:** `01-prompts/22-letterly/02-desktop.md`

---

## 1. When to Use

Activate this skill when:
- Formatting structured desktop execution prompts from transcribed or raw user requests.
- Preparing high-priority tasks requiring explicit numbered action items.
- Appending the mandatory `[06-execute-parent-task-with-n-steps-v6.md](file;.agents/skills/execute-parent-task-with-n-steps-v6)` skill invocation.

---

## 2. Formatting Rules

1. **Zero Conversational Framing:** Never output conversational greetings or introductions like "Certainly! Here's the output:".
2. **High Priority Section:** Lead with `# High Priority Instruction` followed by `[/goal](slashCommand;goal) [/learn](slashCommand;learn) ${Input Text Verbatim}`.
3. **Actionable Items Section:** Follow with `# Actionable Items Must Follow Non-Negotiable` listing concrete steps derived from the instructions.
4. **Mandatory Suffix:** Append the subagent spawn directive pointing to `[06-execute-parent-task-with-n-steps-v6.md](file;.agents/skills/execute-parent-task-with-n-steps-v6)` as a suffix.

---

## 3. Output Format

```markdown
# High Priority Instruction

[/goal](slashCommand;goal) [/learn](slashCommand;learn) ${Input Text Verbatim}

# Actionable Items Must Follow Non-Negotiable

1. Write a plan and spec first
2. ...

Must follow and spawn an agent using the following skill

[06-execute-parent-task-with-n-steps-v6.md](file;.agents/skills/execute-parent-task-with-n-steps-v6)
```
