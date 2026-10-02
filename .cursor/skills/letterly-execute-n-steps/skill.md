---
name: letterly-execute-n-steps
description: Formats input text verbatim into high-priority instructions, action items, additional learning and planning directives, and execute-parent-task-with-n-steps-v6 skill invocation suffix.
---

# Letterly Execute N Steps Mode

> **[/goal](slashCommand;goal)** Formats input text verbatim into high-priority instructions, action items, additional learning and planning directives, and `execute-parent-task-with-n-steps-v6` skill invocation suffix.
> **[/learn](slashCommand;learn)** Ingest user input verbatim, structure into High Priority Instruction and Actionable Items, append skill invocation suffix, and include `/learn` and `/plan` pre-execution directives.

**Source prompt:** `01-prompts/22-letterly/03-execute-n-steps.md`

---

## 1. When to Use

Activate this skill when:
- Structuring multi-step task directives with explicit planning phases before execution.
- Creating standardized N-step task prompts with V6 subagent delegation.
- Enforcing pre-flight spec reading and planning protocols.

---

## 2. Formatting Rules

1. **Zero Conversational Framing:** Omit all chatty filler, acknowledgments, and explanations.
2. **High Priority Section:** `# High Priority Instruction` followed by `${Input Text Verbatim}`.
3. **Actionable Items Section:** `# Actionable Items Must Follow Non-Negotiable` with numbered steps.
4. **Mandatory Suffix:** `Must follow and spawn agent using` followed by `[06-execute-parent-task-with-n-steps-v6.md](file;.agents/skills/execute-parent-task-with-n-steps-v6)`.
5. **Additional Instructions:** Include `learn [/learn](slashCommand;learn) if you have to learn something and [/plan](slashCommand;plan) stuff before working please.`

---

## 3. Output Format

```markdown
# High Priority Instruction

${Input Text Verbatim}

# Actionable Items Must Follow Non-Negotiable

1. Write a plan and spec first
2. ...

Must follow and spawn agent using

[06-execute-parent-task-with-n-steps-v6.md](file;.agents/skills/execute-parent-task-with-n-steps-v6)

## Additional Instructions

learn [/learn](slashCommand;learn) if you have to learn something and [/plan](slashCommand;plan) stuff before working please.
```
