# {Title of the Document}

**Version:** 1.0.0
**Updated:** YYYY-MM-DD
**AI Confidence:** Draft
**Ambiguity:** None

---

## Keywords

`keyword-one` · `keyword-two` · `keyword-three`

---

## Scoring

| Criterion | Status |
|-----------|--------|
| `00-overview.md` present in module | ✅ / ❌ |
| AI Confidence assigned | ✅ / ❌ |
| Ambiguity assigned | ✅ / ❌ |
| Keywords present | ✅ / ❌ |
| Scoring table present | ✅ / ❌ |

---

## Purpose

One paragraph: what this document specifies, who reads it, and what they should be able to do after reading it. No project marketing — just scope.

---

## Document Inventory

(For `00-overview.md` files only — list every sibling file in this module.)

| # | File | Purpose |
|---|------|---------|
| 01 | `01-…md` | … |
| 99 | `99-consistency-report.md` | Structural health check |

---

## Specification

The actual content. Use H2 (`##`) for top-level sections, H3 (`###`) for sub-sections, H4 (`####`) only when truly needed. Avoid H5+.

### Examples

Use fenced code blocks with the language hint set:

```ts
// TypeScript example
```

```go
// Go example
```

Tables for rule comparisons:

| ❌ Don't | ✅ Do | Why |
|---------|------|-----|
| … | … | … |

---

## Cross-References

Add file-relative links here. Always include `.md`. Examples (replace before committing):

```
- [Related module](../NN-related-module/00-overview.md)
- [Strictly-avoid quick reference](../17-consolidated-guidelines/00-strictly-avoid-quickref.md)
```

### Placeholder cross-references (copy-paste snippet)

Use one of two formats when you need to reserve a link before its
target file exists. The cross-link checker ignores both, so neither
will fail CI.

**Preferred — custom tag** (`<spec-placeholder>`). Self-documenting,
greppable, and the only block format the cross-link checker actively
ignores:

```markdown
<spec-placeholder reason="Activate when target is created.">
- [Target Title](../NN-module-name/00-overview.md)
- [Target Title](../NN-module-name/01-file-name.md#section-anchor)
</spec-placeholder>
```

**Legacy — HTML comment**. Still accepted by the placeholder linter,
but **not** stripped by the cross-link checker — the only reason links
inside it pass today is that they happen to resolve, or they live
inside a fenced code block. New placeholders should use the tag form
above.

```markdown
<!-- TODO: Activate when target is created.
- [Target Title](../NN-module-name/00-overview.md)
- [Target Title](../NN-module-name/01-file-name.md#section-anchor)
-->
```

Guidelines for placeholders (both formats):
- Keep the comment block contiguous (no blank lines inside).
- Replace `NN-module-name` and `01-file-name.md` with real paths before removing the comment markers.
- Remove the `<spec-placeholder>` / `</spec-placeholder>` wrappers (or `<!--`/`-->` for the legacy form) once the target exists.
- If the anchor (`#section-anchor`) is unknown, omit it and add it later.
- Prefer `<spec-placeholder>` for new authoring — only it is recognised by the cross-link checker's selective ignore.
- Don't reserve the same target file in more than one placeholder — the linter's P-007 rule flags duplicates within a file *and* across files (anchor differences are collapsed). If two specs really need to link to the same future doc, activate the placeholder once it lands and let the live link be referenced from both places.
- The wording after `TODO:` (or inside `reason="…"`) must start with an imperative verb from the P-001 allowlist (`activate`, `add`, `link`, `replace`, `wire`, `update`, `write`, `create`, `document`, `cross-reference`) and end with a period. Extend the list with `--allow-verb <verb>` when invoking the linter if you need a different verb.

### How to activate placeholders

When the target file finally lands, "activate" the placeholder by
promoting the bullets out of the comment so the cross-link checker
starts validating them. Three steps, in order:

1. Confirm the target file exists and the anchor (if any) matches a real heading.
2. Delete the opening `<!-- TODO: ...` line and the closing `-->` line.
3. Run `python3 linter-scripts/check-spec-cross-links.py --root spec --repo-root .`
   to verify the now-live links resolve.

#### Example 1 — single placeholder, target now exists

**Before** (placeholder, ignored by the checker):

```markdown
<!-- TODO: Activate when target is created.
- [Database conventions](../04-database-conventions/00-overview.md)
-->
```

**After** (live link, validated by the checker):

```markdown
- [Database conventions](../04-database-conventions/00-overview.md)
```

#### Example 2 — partial activation, one target still pending

Split the block: promote the resolved bullet out, keep the unresolved
one wrapped. Do **not** leave a half-commented block — the linter
(`check-placeholder-comments.py`) will reject mixed prose inside `<!-- -->`.

**Before:**

```markdown
<!-- TODO: Activate when targets are created.
- [Naming conventions](../04-database-conventions/01-naming-conventions.md)
- [Schema design](../04-database-conventions/02-schema-design.md)
-->
```

**After** (first target shipped, second still pending):

```markdown
- [Naming conventions](../04-database-conventions/01-naming-conventions.md)

<!-- TODO: Activate when target is created.
- [Schema design](../04-database-conventions/02-schema-design.md)
-->
```

#### Example 3 — anchor added after the fact

If you originally omitted the `#section-anchor`, add it during
activation rather than leaving a stale link to the file root.

**Before:**

```markdown
<!-- TODO: Activate when target is created.
- [Free-text columns](../04-database-conventions/02-schema-design.md)
-->
```

**After** (anchor confirmed against the target's heading):

```markdown
- [Free-text columns](../04-database-conventions/02-schema-design.md#free-text-columns)
```

> **Do not** simply delete the `<!--` / `-->` markers without re-running
> the cross-link checker — a mistyped path or stale anchor will only
> surface once the link is live.

---

## Pre-flight checklist (delete before committing)

- [ ] H1 title matches the file's purpose
- [ ] Version + Updated date filled in (ISO `YYYY-MM-DD`)
- [ ] AI Confidence + Ambiguity assigned
- [ ] Keywords list non-empty
- [ ] Scoring table filled in
- [ ] Cross-references resolve (run `python linter-scripts/check-spec-cross-links.py --root spec --repo-root .`)
- [ ] Validator passes (`python linter-scripts/validate-guidelines.py`)
- [ ] Sync scripts run in order — see `17-consolidated-guidelines/01-spec-authoring.md` §X.2

---

*Spec template — see `17-consolidated-guidelines/01-spec-authoring.md` for full authoring conventions.*