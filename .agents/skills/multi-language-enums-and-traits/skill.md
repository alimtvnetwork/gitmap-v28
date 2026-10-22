---
name: multi-language-enums-and-traits
description: "Executes the Multi Language Enums And Traits prompt. Autonomously scan, plan, refactor, and fix all Enum, Trait, and Pattern Matching architectures across Go, TypeScript, PHP, Rust, and Python codebases, enforcing string-backed enums, MIME and media type enums (MimeType for Blob/Clipboard/Headers), HasEnumHelpers traits, exhaustive. Use when the user asks to run multi-language-enums-and-traits, or the task is about coding-guideline execution for this rule family."
---

# Multi Language Enums And Traits

Source prompt: `01-prompts/15-cg-execute/16-multi-language-enums-and-traits.md`

## Instructions

1. Read `01-prompts/15-cg-execute/16-multi-language-enums-and-traits.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Spec 04 Addendum — Repository Pattern & Enum Actions (normative)

- Every action is an enum member; dispatch by exhaustive switch/match only; `default` raises a typed error — never a silent no-op, never a wildcard arm on a finite domain enum.
- Check `IsValid` before dispatch; external input enters only through `Parse`/`TryParse`.
- Util-methods checkup for every enum, in every language: `String()`/label, `Parse`/`TryParse`, `IsValid`, `All`/`Values` (PHP: `HasEnumHelpers` `values()`/`names()` count toward `All`/`Values` only). Missing util = add it in the same change.
- `fromOrDefault()` MUST NOT be used for action dispatch or external input — unknown values return a typed error.

## Goal

Autonomously scan, plan, refactor, and fix all Enum, Trait, and Pattern Matching architectures across Go, TypeScript, PHP, Rust, and Python codebases, enforcing string-backed enums, MIME and media type enums (MimeType for Blob/Clipboard/Headers), HasEnumHelpers traits, exhaustive.
