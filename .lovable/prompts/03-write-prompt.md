# Write Memory

> **Purpose:** After completing work or at the end of a session, the AI must persist everything it learned, did, and left undone — so the next AI session can pick up seamlessly with zero context loss.
>
> **When to run:** At the end of every session, after completing a task batch, or when explicitly asked to "update memory", "write memory", or "end memory".

---

## Core Principle

> **The memory system is the project's brain.** If you did something and didn't write it down, it didn't happen. If something is pending and you didn't record it, it will be lost. Write memory as if the next AI has amnesia — because it does.

---

## Phase 1 — Audit Current State

Before writing anything, take inventory:

- **Done this session:** every task completed, every file created/modified/deleted, every decision made.
- **Still pending:** started-but-unfinished, discussed-but-unstarted, blocked items.
- **Learned:** new patterns/conventions, gotchas, user preferences (explicit or implicit).
- **Went wrong:** bugs + root causes, failed approaches, things never to repeat.

---

## Phase 2 — Update Memory Files

**Target:** `.lovable/memory/`

1. Read `mem://index.md` (or the equivalent index file) — understand what already exists, do not duplicate.
2. Update existing memory files affected by this session. Mark completed items with `✅`. **Never truncate or overwrite unrelated entries.**
3. If new knowledge doesn't fit any existing file, create `XX-descriptive-name.md` in the appropriate subfolder and **immediately update the index** in the same operation.
4. Update workflow state in `.lovable/memory/workflow/`:

| Status | Marker |
|--------|--------|
| Done | `✅ Done` |
| In Progress | `🔄 In Progress` |
| Pending | `⏳ Pending` |
| Blocked | `🚫 Blocked — [reason]` |
| Avoid / Skip | `🚫 Blocked — [avoid]` |

---

## Phase 3 — Update Plans & Suggestions

### 3A — Plans (`.lovable/plan.md`)

- Update task statuses (done / in progress / pending).
- Add tasks discovered this session.
- Fully complete plan items move to a `## Completed` section at the bottom of the same file (do not delete).

### 3B — Suggestions (`.lovable/suggestions.md`)

Single file. Structure:

```markdown
## Active Suggestions
### [Title]
- **Status:** Pending | In Review | Approved | Rejected
- **Priority:** High | Medium | Low
- **Description:** What and why
- **Added:** [date / session ref]

## Implemented Suggestions
### [Title]
- **Implemented:** [date / session ref]
- **Notes:** Implementation details
```

When a suggestion is implemented: move it from Active → Implemented, add notes, reference the commit/file.

---

## Phase 4 — Update Issues

### 4A — Pending (`.lovable/pending-issues/XX-short-description.md`)

```markdown
# [Title]
## Description
## Root Cause   (or "Under investigation.")
## Steps to Reproduce
## Attempted Solutions
- [ ] Approach 1 — [result]
## Priority   High | Medium | Low
## Blocked By (if applicable)
```

### 4B — Solved (`.lovable/solved-issues/XX-short-description.md`)

When resolved, **move** the file from `pending-issues/` → `solved-issues/` and append:

```markdown
## Solution
## Iteration Count
## Learning
## What NOT to Repeat
```

### 4C — Strictly Avoided (`.lovable/strictly-avoid.md`)

If a solved issue revealed a forbidden pattern:

```markdown
- **[Pattern Name]:** [Why forbidden]. See: `.lovable/solved-issues/XX-filename.md`
```

Anything the user says to skip or avoid → add here too.

### 4D — CI/CD Issues (`.lovable/cicd-issues/XX-name.md` + `.lovable/cicd-index.md`)

Every CI/CD pipeline failure (linters, vulncheck, build, release, test workflows) gets logged here.

- **File:** `.lovable/cicd-issues/XX-short-name.md` (sequence starts at `01`, lowercase-hyphenated).
- **Index:** `.lovable/cicd-index.md` — a summary table of all CI/CD issues with status. Update in the same operation as creating/editing an issue file.
- **No duplicates** — if the same root cause recurs, append a `## Recurrence` section to the existing file instead of creating a new one.
- **Required structure** for each issue: Pipeline (tool/command/runner), Symptom (raw error), Root Cause, Fix, Verification, Status, Prevention.

---

## Phase 5 — Consistency Validation

1. **Index integrity** — every file in `.lovable/memory/` (incl. subfolders) is listed in the index. If not, add it.
2. **Cross-reference** — every `✅ Done` in `plan.md` has evidence (memory entry, solved issue, code change). Every actionable `pending-issues/` item is reflected in `plan.md` or `suggestions.md`. No file in **both** pending and solved.
3. **Orphan check** — no memory file without an index entry, no "Implemented" suggestion without code evidence, no solved issue without a `## Solution` section.
4. **Final confirmation** respond with:

```
✅ Memory update complete.
Session Summary:
- Tasks completed: [X]
- Tasks pending: [Y]
- New memory files created: [Z]
- Issues resolved: [N]
- Issues opened: [M]
- Suggestions added: [S]
- Suggestions implemented: [T]

Files modified: [list]
Inconsistencies found and fixed: [list, or "None"]
The next AI session can pick up from: [current state + next logical step]
```

---

## File Naming & Structure Rules

| Rule | Example |
|------|---------|
| Numeric prefix | `01-auth-flow.md` |
| Lowercase, hyphenated | `03-error-handling.md` ✅ / `03_Error_Handling.md` ❌ |
| Plans | single `.lovable/plan.md` |
| Suggestions | single `.lovable/suggestions.md` |
| Pending issues | one file per issue under `pending-issues/` |
| Solved issues | one file per issue under `solved-issues/` |
| Memory | grouped by topic under `.lovable/memory/{workflow,features,tech,style,project,...}/` |
| Completed plans/suggestions | `## Completed` section in the **same** file (no separate `completed/` folders) |

### Folder reference

```
.lovable/
├── overview.md
├── strictly-avoid.md
├── user-preferences
├── plan.md
├── suggestions.md
├── prompt.md
├── prompts/
│   ├── 01-read-prompt.md
│   └── 02-write-prompt.md
├── memory/
│   ├── index.md
│   ├── workflow/
│   ├── features/
│   ├── tech/
│   ├── style/
│   ├── project/
│   └── suggestions/
├── pending-issues/
├── solved-issues/
├── cicd-index.md
└── cicd-issues/
```

> ⚠️ **NEVER** create `.lovable/memories/` (with trailing `s`). The correct path is `.lovable/memory/`.

---

## Anti-Corruption Rules

1. **Never delete history** — mark done, move to completed sections, never remove.
2. **Never overwrite blindly** — read before writing, preserve existing content.
3. **Never leave orphans** — every file indexed, every reference resolves.
4. **Never split what should be unified** — plans and suggestions live in ONE file each.
5. **Never mix states** — an issue is either pending or solved, never both.
6. **Never skip the index update** — creating a memory file requires updating `index.md` in the same operation.
7. **Never assume the next AI knows anything** — write as if explaining to a stranger with only the files to go on.

---

*Version 1.0. Must stay in sync with [01-read-prompt.md](./01-read-prompt.md).*
