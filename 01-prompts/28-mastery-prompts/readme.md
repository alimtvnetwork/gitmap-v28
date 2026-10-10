# Mastery Prompts (`28-mastery-prompts`) — Index & Catalog

> [!IMPORTANT]
> Category: `28-mastery-prompts`
> Architecture: Foolproof Standalone Session Bootstrapping Directive
> Scope: Fresh-machine / fresh-chat bootstrap — GitHub, GitMap, repos, memory, chats

---

## Overview

The `28-mastery-prompts` category houses the mastery bootstrap prompt: a single
paste-into-any-chat directive that reproduces the user's full working
environment from zero. It boots a fully operational agent on any machine —
existing or fresh — through eight phases with a check and a fallback at every
step: preflight mode detection, memory intake, GitHub auth, GitMap build and
install, canonical repo-set cloning, coding-guideline intake, the standing
execution protocol, and repo side-chat creation, ending with a brief summary.

Unlike the master onboarding prompt (`27-muse-prompts`), which assumes the
environment, the mastery prompt *builds* the environment. Paste it anywhere;
it handles every situation, including a machine with no memory files, no Go,
no GitHub auth, and no clones yet.

New versions go here as new files. Existing prompts are never modified.

---

## Directory Index & Catalog

| Prompt File | Version | Scope | Key Capabilities |
| :--- | :--- | :--- | :--- |
| [`01-mastery-bootstrap-prompt.md`](01-mastery-bootstrap-prompt.md) | 1.2.0 | Standalone Bootstrap | Turbo-mode operating contract + **GitMap field manual** (pipeline errors: `pe`/`pe -t`/`pe all`/`te all`/`history-ai`; prompt templates: `prompt ls`/`show`; portable repo sets; agent tasks; everyday commands), Phase 0 preflight (fresh/existing machine), full memory intake (MEMORY/USER/AGENTS/SOUL/TOOLS/IDENTITY/alignment/people/groups/repo legend), GitHub auth, Go + GitMap install (quick-installer primary, build-from-source fallback), **GitMap skill-file intake** (mandatory SKILL.md read), 10-repo canonical clone set **via `gitmap clone`** (GITHUB_TOKEN pattern for private repos), guideline intake, Phase 6 execution protocol (breakdown → RUNNING+ETA → A=2/H=2 → 5-min pings → commit+push), repo side-chat creation, full fallback table. |

---

## Directory Tree

```text
01-prompts/28-mastery-prompts/
├── 01-mastery-bootstrap-prompt.md   # the mastery prompt (paste into any fresh chat, any machine)
└── readme.md                        # this index
```
