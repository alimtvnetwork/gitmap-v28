---
name: variadic-and-spread-parameters
description: Executes the Variadic & Spread Parameters prompt. Autonomously scan, plan, refactor, and fix all rigid slice and array parameters across Go, TypeScript, and Rust to variadic and spread patterns.
---

# Variadic & Spread Parameters, Rest Elements & Slices

Source prompt: `01-prompts/15-cg-execute/36-variadic-and-spread-parameters.md`
Companion skill: `.agents/skills/cg-variadic-and-spread-parameters/skill.md`

## Instructions

1. Read `01-prompts/15-cg-execute/36-variadic-and-spread-parameters.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. Replace rigid slice/array parameters with variadic patterns (`...T` in Go, `...items: readonly T[]` in TypeScript, `&[T]` in Rust).
4. Eliminate artificial single-item slice/array literals at call sites.
5. Support native slice/array unpacking via spread syntax at multi-item call sites.
6. Handle empty/nil collections safely without errors.
7. Do not shorten, paraphrase, or skip checklist items in the source prompt.
