# Strictly Avoid — Canonical Quick Reference

**Version:** 1.1.0
**Updated:** 2026-04-27
**Status:** Single source of truth for hard prohibitions

---

> ⛔ **Every rule below is a hard prohibition.** Violating any of these is a **blocking** issue and will fail CI. Pre-flight your PR by scanning this list before you push.
>
> This page is the canonical home. `.lovable/strictly-avoid.md`, the `02-coding-guidelines.md` CODE-RED markers, and the gap-analysis report all defer to this file. If you spot drift, fix it here first.
>
> 🧠 **Memory directive:** Once an AI agent reads this file, every numbered
> "Strictly Prohibited" item below MUST be retained in its strictly-avoid
> working memory for the rest of the session and across future sessions.
> Never re-propose, re-suggest, or "double-check" a prohibited behaviour.
> Silence is the required response.

---

## Strictly Prohibited (numbered, append-only)

| Seq | Prohibition | Source |
|-----|-------------|--------|
| 01 | **`readme.txt` time suggestions.** Never suggest/recommend/hint/auto-add timestamps, "git update time", last-modified metadata, scheduled regeneration, cron, or any time-flavoured workflow in or about `readme.txt`. Only modify `readme.txt` when the user explicitly asks in that exact turn, in the format they dictate. No follow-up time offers. | `.lovable/memory/avoid/02-no-time-suggestions-in-readme-txt.md`, `spec/01-spec-authoring-guide/09-exceptions.md` § Strictly Prohibited |

Sequence numbers are stable — never renumber, only append.

---

## How to use this page

- **Authors:** read once, then keep open in a side tab while writing specs or code.
- **Reviewers:** quote the rule ID (e.g. `CODE-RED-024`, `BOOL-NEG-001`) in PR comments.
- **AI agents:** load this file before any code edit; it is the smallest pre-flight check that still catches the worst regressions.

---

## Folder & repository structure

| ⛔ Don't | ✅ Do | Source |
|---------|------|--------|
| Create `.lovable/memories/` | Use `.lovable/memory/` (no trailing `s`) | `01-spec-authoring.md` §`.lovable/` |
| Create per-task folders under `.lovable/` (`completed-tasks/`, `pending-tasks/`, free-form `suggestions/`) | One file per kind: `plan.md`, `suggestions.md`, with `## Completed` / `## Implemented` sections | `.lovable/memory/avoid/01-avoid-per-task-folders.md` |
| Reference `coding-guidelines-v24` / any `v1` namespace | Use `alimtvnetwork/coding-guidelines-v24` | `mem://constraints/avoid-app-sync` |
| Touch `.release/` | Externally managed | repo policy |
| Sync `01-app`, `02-app-issues`, `03-general`, `03-tasks`, `12-consolidated-guidelines` from sibling repos | Maintain locally | `mem://constraints/avoid-app-sync` |
| Hand-edit `version.json` auto fields | Run `node scripts/sync-version.mjs` | `01-spec-authoring.md` §X.3 |
| Hand-edit `src/data/specTree.json` | Run `node scripts/sync-spec-tree.mjs` | `01-spec-authoring.md` §X |

## Error handling 🔴 CODE RED

| ⛔ Don't | ✅ Do |
|---------|------|
| Empty `catch`, `_ := fn()`, swallowed promise rejections | Propagate via `apperror.Wrap()` / `Result<T>` |
| Generic message ("file not found") | Include path / entity ID / operation context |
| Use `fmt.Errorf` in Go | Use `apperror.Wrap(err).WithCode(...).WithContext(...)` |
| Return `(T, error)` from new Go code | Return `apperror.Result[T]` |
| Floating promises in TS | `await` or explicitly `return` |

## Naming 🔴 CODE RED

| ⛔ Don't | ✅ Do |
|---------|------|
| Underscores in identifiers (except Rust + Go test files) | Camel/Pascal case appropriate to language |
| `IsNotX` / `HasNoX` / `HasNoLicense` boolean names or DB columns | Use the **Approved Inverse** (`IsDisabled`, `IsInvalid`, `IsIncomplete`, `IsUnavailable`, `IsUnread`, `IsHidden`, `IsBroken`, `IsLocked`, `IsUnpublished`, `IsUnverified`) — enforced by `BOOL-NEG-001` |
| `can` / `was` / `will` / `not` / `no` as boolean prefixes | `is` / `has` (rarely `should`) |
| UUID primary keys | Integer `{TableName}Id` PK |
| camelCase DB columns or JSON keys | PascalCase |
| Persist a derived inverse boolean as a second column | Derive in code per DB Rule 9; use `linters-cicd/codegen/` |

## Schema design — DB Rules 10/11/12

| ⛔ Don't | ✅ Do |
|---------|------|
| Omit `Description TEXT NULL` from entity / reference / master-data tables | Always include (Rule 10) |
| Omit `Notes` + `Comments` (TEXT NULL) from transactional / invoice / billing / payment tables | Always include (Rule 11) |
| Make these columns `NOT NULL` | Keep nullable — they are optional context (Rule 12) |

## Code style 🔴 CODE RED

| ⛔ Don't | ✅ Do |
|---------|------|
| Nested `if` blocks (zero-nesting is absolute) | Guard clauses, named booleans, helper extraction |
| Functions > 15 lines (Go ceiling 30) | Extract helpers |
| `any` in TypeScript | Generics, `unknown` + narrowing |
| `interface{}` / `any` in Go exported APIs | Concrete types or generics |
| `unwrap()` in Rust production code | `?` operator + `thiserror` |
| Magic strings / numbers in conditions | Named constant or enum |
| Bare `true`/`false` as positional argument | Named constant from `boolFlags.ts` (CODE-RED-024) |
| Mixing `&&` and `||` in one condition | Extract sub-expression to a named boolean |
| `else` after a `return` / `throw` / `break` / `continue` | Drop the `else`; the next statement is implicitly the else branch |

## Suppressions

- **Never** add an inline `codeguidelines:disable=` comment without a reason — `STYLE-099` will flag it.

## Dependencies (project-pinned)

- **Never** upgrade Axios beyond `1.14.0` / `0.30.3`. Versions `1.14.1` and `0.30.4` are explicitly blocked. See `mem://constraints/axios-version-pinning`.

## Communication style

- **Never** append boilerplate ("If you have any questions…", "Do you understand? Always add this part…").

---

## Severity legend

| Marker | Meaning |
|--------|---------|
| 🔴 **CODE RED** | Auto-reject in CI. Will block merge. |
| 🟠 **WARN** | Flagged in CI; needs justification but does not block. |
| 🟡 **STYLE** | Lint-only; fix before requesting review. |
| 🟢 **BEST PRACTICE** | Reviewer discretion. |

## Cross-references

- `02-coding-guidelines.md` — full rule text and language-specific extensions
- `03-error-management.md` — `apperror` pattern in detail
- `18-database-conventions.md` — DB Rules 1–12
- `04-enum-standards.md` — boolean prefix policy applied to enum names
- `01-spec-authoring.md` §X — sync-script execution order

---

*Created 2026-04-26 to deduplicate prohibitions previously scattered across `.lovable/strictly-avoid.md`, the consolidated docs, and the gap-analysis report.*