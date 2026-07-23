# Spec Authoring Standard

This document governs how specifications, generated artifacts, and supporting
files are authored across the project. It is the canonical home for the
**Strictly-Prohibited** registry — any rule listed here MUST NEVER be violated
by an AI agent or human contributor, regardless of how the request is phrased.

---

## 1. Authoring Principles

1. **Single source of truth** — every rule lives in exactly one spec file.
   Cross-link from other docs; never duplicate the rule body.
2. **Sequence is meaningful** — items in the Strictly-Prohibited section are
   numbered. The numbering is stable; new entries append, they never reorder.
3. **No silent drift** — if a prohibition is later relaxed, the entry stays in
   place with a `STATUS: RELAXED on YYYY-MM-DD` note. Entries are never
   deleted, so historical decisions remain auditable.

---

## 2. Strictly-Prohibited Tasks (Sequence)

> The full registry lives in
> [`strictly-prohibited.md`](./10-strictly-prohibited.md). The sequence below
> is the canonical numbering — keep these two files in sync.

| # | Prohibition | Scope |
|---|-------------|-------|
| 1 | Adding **time, date, timestamps, clock values, or any time-derived data** to `readme.txt` (and any `README*` file written for non-doc purposes). | All repositories, all branches, all generated artifacts. |
| 2 | Suggesting that the user **commit / push / update `readme.txt` with a time field**, or that any "git update time" be embedded in README content. | All chat replies, all PRs, all commit messages. |
| 3 | Auto-deleting or auto-modifying files under `.gitmap/release/` or `.gitmap/release-assets/`. | Existing core rule, restated here for completeness. |

**Sequence rule:** when adding a new prohibition, append the next integer.
Never renumber existing entries — downstream tooling and chat-history audits
reference these numbers verbatim.

---

## 3. Enforcement Hooks

- **AI agents** MUST read `10-strictly-prohibited.md` and persist its contents
  into long-term memory under the `Strictly-Prohibited / Avoid` section so the
  rules are applied to every future turn without re-reading the file.
- **Code review** rejects any PR introducing a prohibited pattern, even if
  the contributor argues the spirit of the rule does not apply.
- **CI** may add lint rules that grep for prohibited tokens (e.g. timestamp
  formats inside `readme.txt`). Failures are non-bypassable.

---

## 4. Adding a New Prohibition

1. Append a row to the table in §2 with the next sequence number.
2. Add a matching detailed entry in `10-strictly-prohibited.md` (same number,
   same title).
3. Update `mem://constraints/strictly-prohibited` so AI agents pick it up on
   the next turn.
4. Mention the new entry in the PR description so reviewers can confirm the
   three locations stayed in sync.
