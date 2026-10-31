---
name: string-normalization-and-equalfoldany
description: Executes the String Normalization & EqualFoldAny prompt. Autonomously scan, audit, plan, and refactor repetitive string comparisons, ToLower/TrimSpace calls, and chained equality checks to strutil.EqualFoldAnyTrim.
---

# String Normalization & EqualFoldAny

Source prompt: `01-prompts/15-cg-execute/37-string-normalization-and-equalfoldany.md`
Companion skill: `.agents/skills/cg-string-normalization-and-equalfoldany/skill.md`

## Instructions

1. Read `01-prompts/15-cg-execute/37-string-normalization-and-equalfoldany.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. Replace manual trimming + case-folding + chained `||` checks with centralized `strutil.EqualFoldAnyTrim(target, ...candidates)`.
4. Check canonical helper locations (`pkg/strutil/strutil.go` or repo string utility) before creating duplicate utilities.
5. Do not shorten, paraphrase, or skip checklist items in the source prompt.
