---
name: lazy-regex-and-pattern-matching
description: "Executes the Lazy Regex And Pattern Matching prompt. Autonomously scan, plan, refactor, and verify all regular expression usage across the codebase, eliminating raw regexp.MustCompile and regexp.Compile outside package lazyregex, replacing inline compilation and blind. Use when the user asks to run lazy-regex-and-pattern-matching, or the task is about coding-guideline execution for this rule family."
---

# Lazy Regex And Pattern Matching

Source prompt: `01-prompts/15-cg-execute/21-lazy-regex-and-pattern-matching.md`

## Instructions

1. Read `01-prompts/15-cg-execute/21-lazy-regex-and-pattern-matching.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

Autonomously scan, plan, refactor, and verify all regular expression usage across the codebase, eliminating raw regexp.MustCompile and regexp.Compile outside package lazyregex, replacing inline compilation and blind.
