# Suggestions workflow

Filesystem contract for capturing Lovable suggestions so any AI can continue the backlog reliably.

## Location

- Active suggestions: `.lovable/memory/suggestions/`
- Archived (done) suggestions: `.lovable/memory/suggestions/archive/`
- Legacy freeform log: `.lovable/memory/suggestions/01-suggestions.md` (kept as historical; do not extend)

## File naming

`YYYYMMDD-HHMMSS-suggestion-<slug>.md`

- UTC timestamp at creation.
- `<slug>` is 3 to 6 kebab-case words.
- Example: `20260723-141200-suggestion-fix-spec-01-app-prefix-collisions.md`

## File shape

Frontmatter is required. Body sections are required in the listed order.

```markdown
---
suggestionId: 20260723-141200-fix-spec-01-app-prefix-collisions
createdAt: 2026-07-23T14:12:00Z
source: Lovable
affectedProject: gitmap-cli
status: open   # open | inProgress | done
---

# Fix spec/01-app prefix collisions

## Description
One paragraph explaining the observation.

## Rationale
Why it matters. Link the risk from the reliability report or a concrete failure.

## Proposed change
Specific spec files, code files, or config to change. No implementation code here.

## Acceptance criteria
- Bullet list of verifiable outcomes.

## Completion notes
(Filled in only when status flips to done.)
```

## Lifecycle

1. Create a new file with `status: open`.
2. When work starts, flip to `status: inProgress` and add a `startedAt` field.
3. On completion, flip to `status: done`, add a `completedAt` field, and fill in Completion notes.
4. Move the file to `.lovable/memory/suggestions/archive/` in the same commit that flips to done. Never delete.

## Template

See `_template.md` in this folder. Copy it and rename with a fresh timestamp/slug.

## Rules

- One suggestion per file. No batching.
- Never edit `01-suggestions.md`. It is frozen.
- Never write suggestions into `plan.md`; `plan.md` is the roadmap, this folder is the intake.
- If a suggestion changes materially, do not rewrite history: mark the old one `done` with a note pointing to the new suggestion.
