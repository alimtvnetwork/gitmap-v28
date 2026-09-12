---
name: write-memory
description: Persist session decisions, learned conventions, resolved ambiguities, strictly-avoid rules, and issues into .lovable/ memory and indices.
---

# Memory Persistence & Issue Logging

Persist what happened this turn so the next AI knows everything without guessing. Every decision, plan change, unresolved ambiguity, newly discovered pattern, and fixed bug must be written to `.lovable/` before this turn ends.

## Hard Rules (Non-Negotiable)

1. Folder is `.lovable/memory/`, NEVER `memories/`.
2. Every new memory file under `.lovable/memory/` MUST be registered in `.lovable/memory/01-index.md` in the same operation.
3. Every plan added, moved, or completed MUST update `.lovable/plans/01-index.md` in the same operation.
4. Ambiguity files are NEVER duplicated. Open questions go to `.lovable/ambiguous-questions/01-new-ambiguity/xx-<slug>.md`. When answered, the file is MOVED (`mv`) to `.lovable/ambiguous-questions/02-ambiguity-resolved/xx-<slug>.md` with a `## Resolution` block appended.
5. Never overwrite `.lovable/strictly-avoid.md`. Append only. If a rule was already there, do not duplicate it.
6. When updating existing files (especially indexes, `strictly-avoid.md`, `suggestions.md`), preserve all unrelated content. No silent truncation.
7. Strict Relative Git Paths: NEVER write absolute filesystem paths or `file:///` URIs into markdown files, plans, or code comments.
8. Strict Lowercase File Naming: All files, scripts, documentation, and system files MUST use strictly lowercase naming.
9. Root `readme.md` and `.lovable/what-to-read.md` stay in sync. Same file list, same order.
10. Nothing executes this turn beyond writing to `.lovable/`, root `readme.md` lowercase fixing, and `mv`. No application source code changes.

## Memory Routing Protocol

```
New info discovered
├─ Institutional knowledge (pattern / convention / decision)?
│  YES → .lovable/memory/learned/xx-<slug>.md + update .lovable/memory/01-index.md
├─ Must never happen again?
│  YES → append to .lovable/strictly-avoid.md
├─ Idea, not yet approved?
│  YES → .lovable/suggestions.md
├─ Bug / regression?
│  YES → .lovable/issues/xx-<slug>.md (or .lovable/cicd-issues/ if CI/CD)
├─ New or changed plan?
│  YES → .lovable/plans/pending/xx-<slug>.md + update .lovable/plans/01-index.md
├─ Ambiguity / unclear requirement blocking progress?
│  YES → .lovable/ambiguous-questions/01-new-ambiguity/xx-<slug>.md
└─ User answered an open ambiguity?
   YES → mv to .lovable/ambiguous-questions/02-ambiguity-resolved/xx-<slug>.md + append ## Resolution
```

## Completion Confirmation

Reply with this exact markdown block:

```markdown
# Memory Update Complete

- Plans completed this turn: [N]
- Plans created this turn: [N]
- Ambiguities resolved this turn: [N]
- Ambiguities opened this turn: [N]
- Issues logged this turn: [N]
- CI/CD issues logged this turn: [N]
- Memory files written: [N]
- Skills updated: [S]
- Rules updated: [U]
- Suggestions logged: [N]
- Commands logged: [N]
- Root readme lowercase verified: [Yes/No]

## Current State Summary

- Total pending plans: [N]  (from .lovable/plans/01-index.md)
- Total open ambiguities: [N]  (from 01-new-ambiguity/)
- Total CI/CD issues open: [N]  (from cicd-index.md)
- Total institutional memory files: [N]  (from .lovable/memory/01-index.md)

Next turn will read this state cleanly.
```
