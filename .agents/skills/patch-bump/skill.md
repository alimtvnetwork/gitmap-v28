---
name: patch-bump
description: "Executes the Patch Version Bump - Release Management prompt. Trigger Keywords & Aliases: patch bump, patch release, bump patch version, hotfix bump, bugfix release. Use when the user asks to run patch-bump, or the task is about a version bump or a release."
---

# Patch Version Bump - Release Management

Source prompt: `01-prompts/17-release-management/03-patch-bump.md`

## Instructions

1. Read `01-prompts/17-release-management/03-patch-bump.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## GitMap Release Non-Negotiables
- **Pre-Flight Pull:** Always run `gitmap pull` (fallback `git pull`) before touching any files.
- **Dedicated Release Verb:** Commit release changes on the release branch using `gitmap cpr "RELEASE - vX.Y.Z <scope>"`.
- **TOTAL BAN ON CPF FOR RELEASES:** Never use `gitmap cpf` for releases (it adds a conflicting `Feature:` prefix).
- **TOTAL BAN ON NESTED PREFIXES:** Never write `chore(release):` or `release:` inside the message argument.
- **TOTAL BAN ON THIRD BRACKETS & COLONS:** Use `"RELEASE - vX.Y.Z <scope>"` with hyphen and no square brackets.

## Goal

Trigger Keywords & Aliases: patch bump, patch release, bump patch version, hotfix bump, bugfix release.
