---
name: major-bump
description: "Executes the Major Version Bump - Release Management prompt. Trigger Keywords & Aliases: major bump, major release, bump major version, breaking bump, breaking release. Use when the user asks to run major-bump, or the task is about a version bump or a release."
---

# Major Version Bump - Release Management

Source prompt: `01-prompts/17-release-management/01-major-bump.md`

## Instructions

1. Read `01-prompts/17-release-management/01-major-bump.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## GitMap Release Non-Negotiables
- **Pre-Flight Pull:** Always run `gitmap pull` (fallback `git pull`) before touching any files.
- **Dedicated Release Verb:** Commit release changes on the release branch using `gitmap cpr "vX.0.0 <scope>"`.
- **TOTAL BAN ON CPF FOR RELEASES:** Never use `gitmap cpf` for releases (it adds a conflicting `Feature:` prefix).
- **TOTAL BAN ON NESTED & REDUNDANT PREFIXES:** Never write `chore(release):`, `release:`, `RELEASE - `, or `VERSION - ` inside the message argument for `gitmap cpr`. Format directly as `"vX.0.0 <scope>"`.
- **TOTAL BAN ON THIRD BRACKETS & COLONS:** Use `"vX.0.0 <scope>"` with no square brackets or colons.

## Goal

Trigger Keywords & Aliases: major bump, major release, bump major version, breaking bump, breaking release.
