---
name: letterly-mobile
description: Formats input text verbatim into mobile single-paragraph goal and learn execution directives followed by execute-parent-task-with-n-steps-v6 skill suffix without extra action items or conversational prefixes.
---

# Letterly Mobile Mode

> **[/goal](slashCommand;goal)** Formats input text verbatim into mobile single-paragraph goal and learn execution directives followed by `execute-parent-task-with-n-steps-v6` skill suffix without extra action items or conversational prefixes.
> **[/learn](slashCommand;learn)** Whatever is given as an input, do not add filler or conversational conversational prefixes. Follow the single-paragraph mobile format and append the mandatory V6 skill reference.

**Source prompt:** `01-prompts/22-letterly/01-mobile.md`

---

## 1. When to Use

Activate this skill when:
- Processing mobile-formatted prompt inputs for immediate autonomous execution.
- Compacting multi-paragraph requests into a single continuous paragraph directive.
- Appending the mandatory `[06-execute-parent-task-with-n-steps-v6.md](file;.agents/skills/execute-parent-task-with-n-steps-v6)` suffix.

---

## 2. Formatting Rules

1. **Zero Conversational Framing:** Never output "Certainly! Here's the output based on your instructions:" or similar filler phrases.
2. **Single Paragraph Output:** Combine all input parts into a single continuous line/paragraph without newlines.
3. **No Intermediate Action Items:** Do not generate extra markdown lists or actionable checklists before the directive.
4. **Clean Verbatim Input:** Strip spoken filler words (`um`, `ah`, `wh`) while preserving the core technical directive verbatim.
5. **Mandatory Suffix:** Append `- must follow the skill [06-execute-parent-task-with-n-steps-v6.md](file;.agents/skills/execute-parent-task-with-n-steps-v6)` at the end of the line.

---

## 3. Output Format

```
[/goal](slashCommand;goal) [/learn](slashCommand;learn) ${Input Text Verbatim} - must follow the skill [06-execute-parent-task-with-n-steps-v6.md](file;.agents/skills/execute-parent-task-with-n-steps-v6)
```
