---
name: write-memory
description: Persist session decisions, learned conventions, resolved ambiguities, strictly-avoid rules, and issues into .ai-memory/ memory and indices.
---

# Memory Persistence & Issue Logging

Persist what happened this turn so the next AI knows everything without guessing. Every decision, plan change, unresolved ambiguity, newly discovered pattern, and fixed bug must be written to `.ai-memory/` before this turn ends.

## Hard Rules (Non-Negotiable)

1. Folder is `.ai-memory/memory/`, NEVER `memories/`.
2. Every new memory file under `.ai-memory/memory/` MUST be registered in `.ai-memory/memory/01-index.md` in the same operation.
3. Every plan added, moved, or completed MUST update `.ai-memory/plans/01-index.md` in the same operation.
4. Ambiguity files are NEVER duplicated. Open questions go to `.ai-memory/ambiguous-questions/01-new-ambiguity/xx-<slug>.md`. When answered, the file is MOVED (`mv`) to `.ai-memory/ambiguous-questions/02-ambiguity-resolved/xx-<slug>.md` with a `## Resolution` block appended.
5. Never overwrite `.ai-memory/strictly-avoid.md`. Append only. If a rule was already there, do not duplicate it.
6. When updating existing files (especially indexes, `strictly-avoid.md`, `suggestions.md`), preserve all unrelated content. No silent truncation.
7. Strict Relative Git Paths: NEVER write absolute filesystem paths or `file:///` URIs into markdown files, plans, or code comments.
8. Strict Lowercase File Naming: All files, scripts, documentation, and system files MUST use strictly lowercase naming.
9. Root `readme.md` and `.ai-memory/what-to-read.md` stay in sync. Same file list, same order.
10. Nothing executes this turn beyond writing to `.ai-memory/`, root `readme.md` lowercase fixing, and `mv`. No application source code changes.
11. Mandatory 30-Commit Git History Audit: Prior to authoring or updating memory, execute `git log -n 30 --oneline` to inspect recent commits, summarize progress, and extract what was learned.
12. Recent 20-Task Tracking & Compact Task Register: Maintain and check `.ai-memory/plans/01-index.md` Recent Completed Tasks Register (last 20 tasks) and `.ai-memory/what-to-read.md` to guarantee continuous loop memory.

## Memory Routing Protocol

```
New info discovered
├─ Institutional knowledge (pattern / convention / decision)?
│  YES → .ai-memory/memory/learned/xx-<slug>.md + update .ai-memory/memory/01-index.md
├─ Must never happen again?
│  YES → append to .ai-memory/strictly-avoid.md
├─ Idea, not yet approved?
│  YES → .ai-memory/suggestions.md
├─ Bug / regression?
│  YES → .ai-memory/issues/xx-<slug>.md (or .ai-memory/cicd-issues/ if CI/CD)
├─ New or changed plan?
│  YES → .ai-memory/plans/pending/xx-<slug>.md + update .ai-memory/plans/01-index.md
├─ Ambiguity / unclear requirement blocking progress?
│  YES → .ai-memory/ambiguous-questions/01-new-ambiguity/xx-<slug>.md
└─ User answered an open ambiguity?
   YES → mv to .ai-memory/ambiguous-questions/02-ambiguity-resolved/xx-<slug>.md + append ## Resolution
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

- Total pending plans: [N]  (from .ai-memory/plans/01-index.md)
- Total open ambiguities: [N]  (from 01-new-ambiguity/)
- Total CI/CD issues open: [N]  (from cicd-index.md)
- Total institutional memory files: [N]  (from .ai-memory/memory/01-index.md)

Next turn will read this state cleanly.
```
