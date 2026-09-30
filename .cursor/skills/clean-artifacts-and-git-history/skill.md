---
name: clean-artifacts-and-git-history
description: "Executes the Artifact Sanitization & Git History Preservation prompt. Ensure that NO assets, zip files from artifacts, test data, temporary scratch scripts, or extraneous generated code are accidentally committed to or retained in the Git repository. Use when the user asks to run clean-artifacts-and-git-history, or the task is about commits, boolean cleanups, or multi-agent code fixes."
---

# Artifact Sanitization & Git History Preservation

Source prompt: `01-prompts/09-commit-and-multi-agent-code-fix/07-clean-artifacts-and-git-history.md`

## Instructions

1. Read `01-prompts/09-commit-and-multi-agent-code-fix/07-clean-artifacts-and-git-history.md` in full before doing the task.
2. Execute that prompt verbatim. It is the source of truth for this workflow.
3. A direct instruction in the current user message overrides the prompt when they conflict.
4. Do not shorten, paraphrase, or skip checklist items in the source prompt.

## Goal

Ensure that NO assets, zip files from artifacts, test data, temporary scratch scripts, or extraneous generated code are accidentally committed to or retained in the Git repository.
