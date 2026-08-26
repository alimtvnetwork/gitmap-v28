# Spec Cross-Link Checker

`check-spec-cross-links.py` walks `spec/`, parses every markdown link, and
fails CI if any internal link points to a missing file or a non-existent
heading anchor.

## Local run

```bash
python3 linter-scripts/check-spec-cross-links.py --root spec

# Optional: JSON report

python3 linter-scripts/check-spec-cross-links.py --root spec --json

# Optional: GitHub annotations

python3 linter-scripts/check-spec-cross-links.py --root spec --github
```

## Exit codes

| Code | Meaning |
|------|---------|
| 0    | All internal links resolve |
| 1    | One or more broken links / missing sections |
| 2    | Invocation error |

## What is checked

- Markdown links of the form `[text](path)` and `[text](path#anchor)`.
- Path must resolve to an existing file (relative to source `.md`, or absolute from repo root).
- If `#anchor` is present, it must match an existing H1–H6 heading slug in the target file.
- Links inside fenced code blocks (```` ``` ```` or `~~~`) are ignored — they are examples, not real references.
- Links inside `<spec-placeholder>...</spec-placeholder>` blocks are ignored — this is the preferred placeholder format for cross-references that authors will activate later. See `spec/_template.md` §"Placeholder cross-references" for the copy-paste snippet.
- Regular HTML comments (`<!-- ... -->`) are **not** stripped: a broken link inside an ordinary comment is still flagged so license headers / TODOs that drift can't hide bugs. To opt a block out of validation, wrap it in `<spec-placeholder>` instead.
- External URLs (`http://`, `https://`, `mailto:`, etc.) and project schemes (`mem://`, `user-uploads://`, `knowledge://`) are skipped.

## Companion check: placeholder block formatting

Because comment-wrapped links are skipped here, malformed placeholder
blocks (missing `-->`, stray prose, external links pretending to be
placeholders) would otherwise rot silently. Run the companion linter to
validate placeholder syntax before commit:

```bash
python3 linter-scripts/check-placeholder-comments.py --root spec --repo-root .
```

It enforces the snippet shape documented in `spec/_template.md` §Placeholder cross-references and runs automatically in the pre-commit hook alongside the cross-link checker.

Rule summary (full docs in the script header):

| Code | Check |
|------|-------|
| P-001 | Intent text (`TODO:` / `FIXME:` / `reason="…"`) must start with an allowlisted imperative verb and end with a period |
| P-002 | Body lines must be `- [text](link)` bullets |
| P-003 | Links must be relative `.md` paths (no `http(s)://`, no anchor-only) |
| P-004 | Block must contain ≥1 valid bullet |
| P-005 | No blank lines inside the block |
| P-006 | Every opening marker must have a matching closer |
| P-007 | No duplicate placeholder targets — within a file or across files (anchor ignored, paths resolved) |

Extend the P-001 verb allowlist for a single run with `--allow-verb <verb>` (repeatable). The default set (`activate`, `add`, `link`, `replace`, `wire`, `update`, `write`, `create`, `document`, `cross-reference`) covers the activation-language used in the spec template.

## Allowlist (waivers)

Known-broken links live in `linter-scripts/spec-cross-links.allowlist`,
one waiver per line:

```

# Comments start with `#` (only at line start; anchor `#` inside entries is preserved)

spec/path/to/file.md:42:./missing-target.md
spec/other.md:99:./file.md#missing-section
```

Format: `<relpath-from-repo-root>:<line>:<exact-target-as-written>`.
Remove a waiver as soon as the underlying link is fixed.

## CI

Runs as the `cross-links` job in `.github/workflows/ci.yml` on every push
and pull request to `main`.

---

## Auto-Fix Suggester (companion)

`suggest-spec-cross-link-fixes.py` consumes the same allowlist and the
same broken-link set as the checker, then proposes the closest match for
each failure using `difflib.SequenceMatcher`:

- `missing-file` → fuzzy-match against every `*.md` under `spec/` (basename
  match wins ties).
- `missing-section` → fuzzy-match the requested anchor against the slug of
  every heading in the resolved target file.

### Modes

| Flag | Behavior |
|------|----------|
| _(default)_ | Report-only. Always exits 0. Used by CI as advisory annotations. |
| `--apply` | Rewrite files in place when `confidence >= --min-confidence` (default `0.82`). Exits 1 if any low-confidence breakage remains. |
| `--github` | Emit GitHub Actions annotations (`::warning` for auto-fixable, `::notice` for manual-review). |
| `--json`   | Machine-readable report on stdout. |

### Local usage

```bash

# Preview what would change

python3 linter-scripts/suggest-spec-cross-link-fixes.py --root spec

# Apply high-confidence fixes only

python3 linter-scripts/suggest-spec-cross-link-fixes.py --root spec --apply

# Stricter threshold

python3 linter-scripts/suggest-spec-cross-link-fixes.py --root spec --apply --min-confidence 0.9
```

### CI integration

The `cross-links` job runs the suggester unconditionally after the
blocking checker. Suggestions are uploaded as the
`spec-cross-links-suggestions` artifact and surfaced as PR annotations.
The suggester never fails the build — it is purely advisory so reviewers
can opt in to applying fixes locally.
