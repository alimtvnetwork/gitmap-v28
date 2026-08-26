# `rename_intake` — diff-mode audit JSON schema

The `--list-changed-files` audit emitted by
[`check-placeholder-comments.py`](./check-placeholder-comments.py)
documents *which* repo-relative paths the diff-mode intake considered
and *why* each one was kept or skipped. With `--with-similarity` the
same payload is enriched with rename/copy provenance recovered from
`git diff --name-status -M -C` (or the equivalent
`--changed-files` payload).

This file is the canonical schema reference. It documents the JSON
shape, the score/unscored distinction, and the STDOUT/STDERR contract
so dashboards and downstream tools can ingest the audit without
scraping the human text table.

## Stream contract — STDERR vs STDOUT

The audit is **always written to STDERR**, regardless of `--json` or
the verbosity flags. STDOUT is reserved for the linter's primary
payload:

| Stream | Content | Format |
|---|---|---|
| **STDOUT** | The human violation summary, **or** with `--json` a single-document JSON object containing the violations array. | Always one well-formed document. |
| **STDERR** | The `--list-changed-files` audit table (text columns) **or** with `--json` a JSON array of `rename_intake` records. Plus the usual progress / dedupe-footer lines. | Text rows or one JSON array. |

This split lets a CI job pipe STDOUT into `jq` for violation triage
while a separate sink (artifact upload, dashboard ingest, log scraper)
consumes the audit on STDERR. **No interleaving.** A row never lands on
STDOUT and a violation never lands on STDERR.

The audit array is a single JSON document — it is emitted in one write
after all rows are collected, so partial / streaming readers see either
the complete array or nothing.

## Record schema

Each row is a JSON object with the following fields:

| Field | Type | Required | Description |
|---|---|---|---|
| `path` | `string` | yes | Repo-relative POSIX path of the considered file (the post-rename / post-copy "new" side for R/C rows). |
| `status` | `string` | yes | One of the closed audit vocabulary: `matched`, `ignored-extension`, `ignored-out-of-root`, `ignored-missing`, `ignored-deleted`. |
| `reason` | `string` | yes | Human-readable explanation for the `status`. Stable enough for log search but not part of the machine contract — do not parse. |
| `similarity` | `object \| null` | only with `--with-similarity` | Rename/copy provenance, see below. Absent from the object entirely when the flag is not set (not `null`). |

Without `--with-similarity` the `similarity` key is **explicitly
removed** from every record so dashboards parsing the legacy schema
continue to work unchanged. The flag is purely additive.

## Status reference

Every `--list-changed-files` row carries a `status` drawn from a
**closed five-value vocabulary**. This section is the single
canonical reference: it lists each status, the exact `reason`
wording the linter emits, the conditions that trigger it, and what
the row looks like in all three output surfaces (JSON audit, text
table, CSV export). Use it when wiring a dashboard, writing a
`grep` rule for CI logs, or sanity-checking a new test fixture.

| `status` | When it fires | Exact `reason` text |
|---|---|---|
| `matched` | Path is under `--root`, has an allow-listed extension, and the post-state file exists on disk. Counted in the linted set. | `under --root, extension allowed, file present on disk` |
| `ignored-extension` | Path is under `--root` but its suffix isn't in the allowlist (e.g. a `.txt` change while linting `.md`). | `extension '<.ext or (none)>' not in allowlist <sorted list>` |
| `ignored-out-of-root` | Path resolves outside `--root` (e.g. a top-level `readme.md` change while linting `spec/`). | `path is outside --root <resolved-abs-path>` |
| `ignored-missing` | Post-state path no longer exists on disk — typically reverted in a later commit of the same push, or filtered out by `.gitignore` on checkout. | `post-state path is not on disk (reverted later in the push, or filtered by .gitignore on checkout)` |
| `ignored-deleted` | A `D`-status row, **or** the OLD side of a rename/copy: there is no post-state file to scan. `reason` is per-source — see [next sub-section](#ignored-deleted-reason-format). | (per-source — see below) |

The five values are **the only ones the renderer will ever emit**;
adding a new status is a major-version contract change announced in
`changelog.md`. Field-order, casing, and dash hyphenation are part
of the contract — `Matched`, `ignored_deleted`, etc. will never be
emitted.

### `ignored-deleted` reason format

`ignored-deleted` is the only status whose `reason` text varies by
intake source. Each delete is tagged at parse time with one of the
provenance keys below, and the audit emitter looks the key up in a
per-source reason table (centralised at `_DELETED_REASON` in
`check-placeholder-comments.py`). All four surfaces (text table,
JSON payload, CSV export, dedupe footer) use the same wording for
a given source.

| Source tag | Triggered by | Full `reason` text | Stable substring (CI-grep safe) |
|---|---|---|---|
| `diff-D` | A real `D`-status row from `git diff --name-status` | `git diff reported D (deleted): file removed in the diff range, no post-state to lint` | `git diff reported D (deleted)` |
| `changed-files-D` | An authored `--changed-files` payload row shaped exactly `` `D\tpath` `` | `` --changed-files payload row shaped `D\tpath`: explicit delete marker, no post-state to lint `` | `` --changed-files payload row shaped `` |
| *(unknown future tag)* | Captured as a delete by a parser this map hasn't caught up with | `path captured as a delete by the diff intake but provenance is unknown — treated as ignored-deleted for safety` | `provenance is unknown` |

The full sentence may be re-worded for clarity in a future minor
release; the **substring column** is what's safe to grep for in CI
log scrapers. To switch on the machine-readable tag instead of the
wording, pass `--list-changed-files-verbose` and read the `source`
field — see [Verbose mode](#verbose-mode---list-changed-files-verbose)
below.

### Examples in every output surface

The samples below come from one invocation against a fixture that
exercises all five statuses. The `--changed-files` payload is:

```text
M	spec/keep.md
D	spec/gone.md
R090	spec/old.md	spec/new.md
A	spec/notes.txt
A	outside/x.md
A	spec/missing.md
```

Disk state: `spec/keep.md` exists, the rename's `spec/new.md` and
the new `spec/missing.md` do **not** (simulating a revert-on-push),
and the linter is invoked with `--root spec` from a working
directory whose absolute path is `/repo` (the linter resolves
`--root` to an absolute path before printing it in the
`ignored-out-of-root` reason — your `<resolved-abs-path>` will of
course differ).

#### JSON (`--json --with-similarity`)

```json
[
  {
    "path": "spec/keep.md",
    "status": "matched",
    "reason": "under --root, extension allowed, file present on disk",
    "similarity": null
  },
  {
    "path": "spec/new.md",
    "status": "ignored-missing",
    "reason": "post-state path is not on disk (reverted later in the push, or filtered by .gitignore on checkout)",
    "similarity": {"kind": "R", "score": 90, "old_path": "spec/old.md"}
  },
  {
    "path": "spec/notes.txt",
    "status": "ignored-extension",
    "reason": "extension '.txt' not in allowlist ['.md']",
    "similarity": null
  },
  {
    "path": "outside/x.md",
    "status": "ignored-out-of-root",
    "reason": "path is outside --root /repo/spec",
    "similarity": null
  },
  {
    "path": "spec/missing.md",
    "status": "ignored-missing",
    "reason": "post-state path is not on disk (reverted later in the push, or filtered by .gitignore on checkout)",
    "similarity": null
  },
  {
    "path": "spec/gone.md",
    "status": "ignored-deleted",
    "reason": "--changed-files payload row shaped `D\\tpath`: explicit delete marker, no post-state to lint",
    "similarity": null
  }
]
```

#### Text table (default `--list-changed-files`, `--with-similarity`)

```text
── placeholder-comments: changed-file audit (6 row(s); +similarity columns) ──
  status               path             kind  score  old          reason
  ----------------------------------------------------------------------
  matched              spec/keep.md     -     -      -            under --root, extension allowed, file present on disk
  ignored-missing      spec/new.md      R     90     spec/old.md  post-state path is not on disk (reverted later in the push, or filtered by .gitignore on checkout)
  ignored-extension    spec/notes.txt   -     -      -            extension '.txt' not in allowlist ['.md']
  ignored-out-of-root  outside/x.md     -     -      -            path is outside --root /repo/spec
  ignored-missing      spec/missing.md  -     -      -            post-state path is not on disk (reverted later in the push, or filtered by .gitignore on checkout)
  ignored-deleted      spec/gone.md     -     -      -            --changed-files payload row shaped `D\tpath`: explicit delete marker, no post-state to lint
  totals: matched=1  ignored-extension=1  ignored-out-of-root=1  ignored-missing=2  ignored-deleted=1
```

With `--list-changed-files-verbose` the same run grows a `source`
column that holds the raw provenance tag for every
`ignored-deleted` row (and `-` everywhere else):

```text
  status               path             source           reason
  -------------------------------------------------------------
  ignored-deleted      spec/gone.md     changed-files-D  --changed-files payload row shaped `D\tpath`: explicit delete marker, no post-state to lint
```

#### CSV export (`--similarity-csv -`)

```text
path,status,reason,kind,score,old_path
spec/keep.md,matched,"under --root, extension allowed, file present on disk",,,
spec/new.md,ignored-missing,"post-state path is not on disk (reverted later in the push, or filtered by .gitignore on checkout)",R,90,spec/old.md
spec/notes.txt,ignored-extension,extension '.txt' not in allowlist ['.md'],,,
outside/x.md,ignored-out-of-root,path is outside --root /repo/spec,,,
spec/missing.md,ignored-missing,"post-state path is not on disk (reverted later in the push, or filtered by .gitignore on checkout)",,,
spec/gone.md,ignored-deleted,"--changed-files payload row shaped `D\tpath`: explicit delete marker, no post-state to lint",,,
```

Cross-surface invariants — true for every status, every output:

- **Row count parity**: the JSON array, the text table body, and
  the CSV body always have the **same number of rows** for the same
  invocation (after dedupe / `--only-changed-status` filtering).
- **Reason wording parity**: the `reason` cell is identical
  byte-for-byte across JSON, text, and CSV — including the
  per-source variation on `ignored-deleted`.
- **Similarity columns** stay empty for any row whose `similarity`
  is `null` (plain rows AND every `ignored-deleted` row); the
  `kind`/`score`/`old_path` cells are only populated on R/C rows.
- **`--with-similarity` is purely additive** for non-deleted
  statuses — toggling it doesn't change which rows appear or what
  `status`/`reason` they carry.

The reason-substring guarantees in the table above are exercised
by `linter-scripts/tests/test_rename_intake_emission_gating.py` and
`linter-scripts/tests/test_ignored_deleted_*.py`, so a re-wording
would break CI before it could surprise a downstream log scraper.

### `reason` for `ignored-deleted` rows

> **Canonical reference:** see [`ignored-deleted` reason
> format](#ignored-deleted-reason-format) under
> [Status reference](#status-reference) for the full per-source
> table, the `--list-changed-files-verbose` source-tag contract,
> and worked examples in JSON / text / CSV. The summary below is
> kept for backwards-compatible deep-links.

The `reason` field on `ignored-deleted` rows is per-source so a
reviewer can see *why* a path was classified as deleted. Two source
vocabularies are emitted today:

| Source | Triggered by | `reason` substring (stable for log-grep) |
|---|---|---|
| `diff-D` | A real `D`-status row from `git diff --name-status` | `git diff reported D (deleted)` |
| `changed-files-D` | An authored `--changed-files` payload row shaped exactly `D\tpath` | `--changed-files payload row shaped` |

The full `reason` text is intentionally not part of the machine
contract (it may be re-worded for clarity), but the substrings above
are stable and safe to grep in CI logs. New provenance tags will be
added alongside their parser changes; an unknown tag falls back to
a clearly labelled "provenance unknown" reason rather than crashing.

### Verbose mode (`--list-changed-files-verbose`)

Pass `--list-changed-files-verbose` (alongside `--list-changed-files`)
to expose the raw provenance tag for every `ignored-deleted` row.
Two surfaces, two contracts:

- **JSON** — adds a `"source": str|null` key to **every** row
  (`str` on `ignored-deleted` rows, `null` everywhere else) so the
  schema stays regular for downstream consumers. Off (the default),
  the key is **omitted entirely** — legacy 3-key schema preserved
  byte-for-byte.
- **Text** — appends a `source` column at the end of the fixed-
  width cells (just before the variable-width `reason`). Non-
  deleted rows render `-` to match the surrounding blank-cell
  convention.

Verbose mode also promises the `reason` wording is **machine-stable**
for `ignored-deleted` rows — no future re-wording — so CI scripts
can match on the full string instead of the substring guarantee
above. Composes cleanly with `--with-similarity` (similarity
columns first, then `source` last), `--dedupe-changed-files`
(first-seen `source` wins), and `--only-changed-status` (filter
runs after source attachment).

### `similarity` sub-object

When `--with-similarity` is on, every record carries a `similarity`
key whose value is either `null` (plain `A` / `M` / `D` rows, no
rename/copy observed) or an object:

| Field | Type | Description |
|---|---|---|
| `kind` | `string` | `"R"` for rename, `"C"` for copy. No other values. |
| `score` | `integer \| null` | Git's similarity percentage (0–100), or `null` when the row is **unscored** — see next section. |
| `old_path` | `string` | The OLD-side repo-relative path the rename/copy is from. Always present on R/C rows. |

## `score` — scored vs unscored

Git emits the similarity percentage on `--name-status -M -C` rows
(e.g. `R092\told\tnew` → 92 % similar). The audit preserves that
integer verbatim in `score`. Two cases produce a `null` score, which
the text renderer prints as a single dash (`-`):

1. **Authored `--changed-files` payloads** that omit the percentage.
   Both shapes are accepted:
   - tab form `R\told\tnew` (no leading score) → `score: null`
   - arrow form `old => new` → `score: null`
   These are common when the changed-file list is hand-built or
   concatenated from tools that don't preserve git's score column.
2. **Plain rows** (`A` / `M` / `D`) where `similarity` itself is
   `null` — there is no rename/copy to score.

A score of `0` is **not** the same as a missing score. `0` means git
observed the rename pair and rated them entirely dissimilar; `null`
means the score was never recorded. Keep them distinct in dashboards.

The text renderer reserves the marker `?` for a future "unknown but
expected" state; it is **not currently emitted**. Today only `-` (for
`null`) and the integer string appear in the `score` column.

## Worked example

`git diff --name-status -M -C main...HEAD` against a branch that
renamed one file, copied one, modified one, and deleted one yields:

```text
R092    docs/old-name.md    docs/new-name.md
C075    spec/template.md    spec/copy.md
M       readme.md
D       legacy.md
```

With `python3 linter-scripts/check-placeholder-comments.py
--diff-base main --list-changed-files --with-similarity --json` the
STDERR audit is:

```json
[
  {
    "path": "docs/new-name.md",
    "status": "matched",
    "reason": "in --root and extension allowed",
    "similarity": {"kind": "R", "score": 92, "old_path": "docs/old-name.md"}
  },
  {
    "path": "spec/copy.md",
    "status": "matched",
    "reason": "in --root and extension allowed",
    "similarity": {"kind": "C", "score": 75, "old_path": "spec/template.md"}
  },
  {
    "path": "readme.md",
    "status": "matched",
    "reason": "in --root and extension allowed",
    "similarity": null
  },
  {
    "path": "legacy.md",
    "status": "ignored-deleted",
    "reason": "file no longer exists on disk (D row)",
    "similarity": null
  }
]
```

Drop `--with-similarity` and the four `similarity` keys disappear
entirely (legacy schema).

## Ready-to-copy example: every `status` × scored / unscored

The worked example above shows the four most common cases. The
payload below is the **full matrix**: every `status` value in the
closed audit vocabulary paired with each shape `similarity` can
take (`{kind, score: int, old_path}`, `{kind, score: null,
old_path}`, and `null`). It's a verified `--with-similarity --json`
output, ready to drop straight into a fixture, a doc snippet, or
a downstream consumer's mock.

The same payload is shipped as a checked-in artifact at
[`linter-scripts/examples/rename-intake-audit.json`](examples/rename-intake-audit.json)
so test code can `json.load()` it without copy-paste drift.

| `status` | `similarity` shape | When it happens |
|---|---|---|
| `matched` | `{kind, score: int, old_path}` | Git's `R<nn>` / `C<nn>` row, post-state under `--root` with an allowed extension. |
| `matched` | `{kind, score: null, old_path}` | Authored `--changed-files` payload using `R\told\tnew` or `old => new` (no leading score). |
| `matched` | `null` | Plain `A` / `M` row — no rename/copy provenance to attach. |
| `ignored-extension` | `{kind, score: int, old_path}` | Renamed to / from an allow-listed extension; rename signal preserved so callers can audit it. |
| `ignored-extension` | `{kind, score: null, old_path}` | Same, but the source payload was unscored. |
| `ignored-extension` | `null` | Plain row whose extension isn't allow-listed. |
| `ignored-out-of-root` | `{kind, score: int, old_path}` | Rename whose post-state path falls outside `--root`. |
| `ignored-out-of-root` | `{kind, score: null, old_path}` | Same, unscored intake. |
| `ignored-out-of-root` | `null` | Plain row outside `--root`. |
| `ignored-missing` | `{kind, score: int, old_path}` | Rename whose post-state file vanished from disk (revert-during-push, .gitignore on checkout). |
| `ignored-missing` | `null` | Plain row whose post-state file is missing. |
| `ignored-deleted` | `null` | `D` row, or the **pre-state** half of a rename — there is no post-state to score, so `similarity` is **always** `null` for this status. |

> **Note** — `ignored-deleted` rows always carry `similarity: null`
> because the row represents a path that no longer exists; the
> rename's *new* side (with the score) is recorded as a separate
> `matched` / `ignored-*` row in the same audit.

```json
[
  {
    "path": "docs/new-name.md",
    "status": "matched",
    "reason": "under --root, extension allowed, file present on disk",
    "similarity": {"kind": "R", "score": 92, "old_path": "docs/old-name.md"}
  },
  {
    "path": "spec/copy.md",
    "status": "matched",
    "reason": "under --root, extension allowed, file present on disk",
    "similarity": {"kind": "C", "score": 75, "old_path": "spec/template.md"}
  },
  {
    "path": "spec/renamed-no-score.md",
    "status": "matched",
    "reason": "under --root, extension allowed, file present on disk",
    "similarity": {"kind": "R", "score": null, "old_path": "spec/old.md"}
  },
  {
    "path": "spec/copy-no-score.md",
    "status": "matched",
    "reason": "under --root, extension allowed, file present on disk",
    "similarity": {"kind": "C", "score": null, "old_path": "spec/src.md"}
  },
  {
    "path": "readme.md",
    "status": "matched",
    "reason": "under --root, extension allowed, file present on disk",
    "similarity": null
  },
  {
    "path": "spec/notes.txt",
    "status": "ignored-extension",
    "reason": "extension '.txt' not in allowlist ['.md']",
    "similarity": {"kind": "R", "score": 88, "old_path": "spec/legacy.txt"}
  },
  {
    "path": "spec/draft.txt",
    "status": "ignored-extension",
    "reason": "extension '.txt' not in allowlist ['.md']",
    "similarity": {"kind": "R", "score": null, "old_path": "spec/sketch.txt"}
  },
  {
    "path": "spec/scratch.txt",
    "status": "ignored-extension",
    "reason": "extension '.txt' not in allowlist ['.md']",
    "similarity": null
  },
  {
    "path": "tools/moved-here.md",
    "status": "ignored-out-of-root",
    "reason": "path is outside --root spec",
    "similarity": {"kind": "R", "score": 95, "old_path": "spec/moved-from-here.md"}
  },
  {
    "path": "tools/cloned-here.md",
    "status": "ignored-out-of-root",
    "reason": "path is outside --root spec",
    "similarity": {"kind": "C", "score": null, "old_path": "spec/cloned-from-here.md"}
  },
  {
    "path": "docs/outside.md",
    "status": "ignored-out-of-root",
    "reason": "path is outside --root spec",
    "similarity": null
  },
  {
    "path": "spec/missing-rename.md",
    "status": "ignored-missing",
    "reason": "post-state path is not on disk (reverted later in the push, or filtered by .gitignore on checkout)",
    "similarity": {"kind": "R", "score": 81, "old_path": "spec/old-missing-name.md"}
  },
  {
    "path": "spec/missing.md",
    "status": "ignored-missing",
    "reason": "post-state path is not on disk (reverted later in the push, or filtered by .gitignore on checkout)",
    "similarity": null
  },
  {
    "path": "spec/deleted.md",
    "status": "ignored-deleted",
    "reason": "git diff reported D (deleted)",
    "similarity": null
  }
]
```

Reading the payload: scan a row's `similarity`. `null` ⇒ no
rename/copy info (a plain row, or any `ignored-deleted` row).
An object whose `score` is an integer ⇒ git observed and rated
the pair. An object whose `score` is `null` ⇒ the row IS a
rename/copy (`kind` is meaningful) but no percentage was ever
recorded. Treat `score: 0` as "git observed and rated 0% similar"
— it is **not** the same as `score: null`.

## Compatibility notes

- **Field order is not part of the contract.** Parse by key, not by
  position. The current implementation emits the keys in the order
  shown above, but that may change.
- **Unknown statuses** are never emitted today — the closed vocabulary
  is enforced at the renderer. New statuses, if added, will land in a
  major version bump and be announced in `changelog.md`.
- **`--dedupe-changed-files`** runs before serialisation; the JSON
  array contains at most one record per `path` when that flag is set,
  with first-seen-wins semantics applied to the `similarity` record
  too.
- **`--only-changed-status`** runs after dedupe and after similarity
  attachment, so filtered-out rows are absent from the array but their
  totals still appear in the human footer on STDERR.

## Validating `rename_intake` output in CI

The repo ships a stdlib-only validator at
[`linter-scripts/validate-rename-intake.py`](./validate-rename-intake.py)
that enforces every rule documented above (key set, closed `status`
vocabulary, scored/unscored/`null` similarity trichotomy, plus the
cross-field invariant that `ignored-deleted` rows always carry
`similarity: null`). It has **zero external dependencies** so CI
doesn't need an extra `pip install` step — Python 3.10+ is
enough.

### Quick reference

```text
usage: validate-rename-intake.py [INPUT] [--with-similarity]
                                 [--with-labels] [--allow-empty]
                                 [--print-schema] [--quiet]
```

| Flag | When to use it |
|---|---|
| `--with-similarity` | The upstream linter was invoked with `--with-similarity`. Validates the enriched 4-key schema (per-record `similarity` object/null required). |
| `--with-labels` | The upstream linter was invoked with `--similarity-labels`. Implies `--with-similarity`. Requires `score_kind` on every non-null `similarity` object. |
| `--allow-empty` | Treat `[]` as valid. Off by default because an empty array almost always means a misconfigured invocation. |
| `--print-schema` | Print the formal JSON Schema (Draft 2020-12) for the requested mode and exit. Pipe into `check-jsonschema` / `ajv` if you want a richer validator. |
| `--quiet` | Suppress the success line; failures still print. |

`INPUT` is a file path, or `-` (the default) to read from stdin.

Exit codes: **`0`** on success, **`1`** on schema violations (details on STDERR), **`2`** on invalid JSON or a CLI usage error. CI gates can rely on these without parsing output.

### Capturing the audit on STDERR

`rename_intake` JSON is emitted on **STDERR** (see the
[stream contract](#stream-contract--stderr-vs-stdout)), so a CI step
typically redirects STDERR into a file and STDOUT into `/dev/null`
(or a separate sink for the violation summary):

```bash
python3 linter-scripts/check-placeholder-comments.py \
  --diff-base origin/main \
  --list-changed-files --with-similarity --json \
  > /dev/null 2> rename-intake.json

python3 linter-scripts/validate-rename-intake.py \
  rename-intake.json --with-similarity
```

If you don't want a temp file, pipe STDERR straight into stdin:

```bash
python3 linter-scripts/check-placeholder-comments.py \
  --diff-base origin/main \
  --list-changed-files --with-similarity --json \
  2>&1 >/dev/null | \
python3 linter-scripts/validate-rename-intake.py - --with-similarity
```

(The `2>&1 >/dev/null` order is deliberate: it sends STDERR to the
pipe and discards STDOUT.)

### Drop-in CI snippets

**GitHub Actions:**

```yaml
- name: Validate rename_intake JSON
  run: |
    python3 linter-scripts/check-placeholder-comments.py \
      --diff-base origin/${{ github.base_ref || 'main' }} \
      --list-changed-files --with-similarity --json \
      > /dev/null 2> rename-intake.json
    python3 linter-scripts/validate-rename-intake.py \
      rename-intake.json --with-similarity

- name: Upload audit artifact
  if: always()
  uses: actions/upload-artifact@v4
  with:
    name: rename-intake-${{ github.run_id }}
    path: rename-intake.json
```

**GitLab CI:**

```yaml
validate_rename_intake:
  image: python:3.12-slim
  script:
    - python3 linter-scripts/check-placeholder-comments.py
        --diff-base "$CI_MERGE_REQUEST_DIFF_BASE_SHA"
        --list-changed-files --with-similarity --json
        > /dev/null 2> rename-intake.json
    - python3 linter-scripts/validate-rename-intake.py
        rename-intake.json --with-similarity
  artifacts:
    when: always
    paths: [rename-intake.json]
```

**Pre-commit hook (local sanity check):**

```yaml
- repo: local
  hooks:
    - id: validate-rename-intake
      name: Validate rename_intake JSON
      language: system
      pass_filenames: false
      entry: bash -c '
        python3 linter-scripts/check-placeholder-comments.py
          --diff-base HEAD --list-changed-files --with-similarity
          --json > /dev/null 2> /tmp/rename-intake.json &&
        python3 linter-scripts/validate-rename-intake.py
          /tmp/rename-intake.json --with-similarity'
```

### Using the formal JSON Schema with external tooling

If your CI already runs a richer JSON-Schema validator (for cross-
language schema gating, registry-backed schemas, etc.) you can
export the equivalent Draft 2020-12 schema and feed it in:

```bash
python3 linter-scripts/validate-rename-intake.py \
  --print-schema --with-similarity > rename-intake.schema.json

# Example: check-jsonschema (pipx install check-jsonschema)

check-jsonschema --schemafile rename-intake.schema.json \
  rename-intake.json
```

The bundled validator and the printed schema are kept in lock-step
by the test suite (`test_validate_rename_intake.py`), so either gate
can be used interchangeably.

## See also

- [`check-placeholder-comments.py --help`](./check-placeholder-comments.py)
  — full flag reference for `--list-changed-files`,
  `--with-similarity`, `--dedupe-changed-files`,
  `--only-changed-status`, `--similarity-csv`.
- [`validate-rename-intake.py --help`](./validate-rename-intake.py)
  — stdlib-only schema validator for `rename_intake` JSON, suitable
  for CI gating; see *"Validating `rename_intake` output in CI"*
  above for ready-to-paste pipeline snippets.
- `linter-scripts/tests/test_with_similarity_flag.py` — executable
  examples of every shape documented above.
- `linter-scripts/tests/test_similarity_csv_export.py` — schema
  examples for the `--similarity-csv` export, including the
  scored-vs-unscored distinction.

## CSV export (`--similarity-csv PATH`)

For spreadsheet review (Excel, Numbers, LibreOffice, `csvkit`,
pandas), the same audit can be exported as RFC 4180 CSV to a file
path or to STDOUT (`-`). The header is **always**
`path,status,reason,kind,score,old_path` regardless of whether
`--with-similarity` was passed; the four similarity columns simply
stay empty when no `_RenameSimilarity` is attached.

| `score` cell | Meaning |
|---|---|
| Empty (`""`) | **Unscored** — plain A/M/D row, OR an authored `--changed-files` payload that omitted the percentage. |
| `"0"` | Git observed the rename/copy and rated the pair entirely dissimilar. |
| `"1"` … `"100"` | Git's similarity percentage. |

`ISBLANK(E2)` and `E2=0` are **not** the same condition — keep them
distinct when filtering. Dedupe and `--only-changed-status` run
BEFORE the export, so the CSV mirrors what you saw on STDERR.

### Field-separator dialect — `--similarity-csv-format {csv,tsv}`

By default the export uses RFC 4180 CSV (comma-separated, double-
quote quoting). When commas are inconvenient — paths or reasons
contain commas, your spreadsheet imports cleaner from tabs, or
you're piping into `cut -f` / `awk -F'\t'` / `column -t -s$'\t'` /
`q -t` — pass `--similarity-csv-format tsv` to switch to the
stdlib `csv.excel_tab` dialect:

```bash
python3 linter-scripts/check-placeholder-comments.py \
  --list-changed-files --similarity-csv audit.tsv \
  --similarity-csv-format tsv
```

Contract:

| | `csv` (default) | `tsv` |
|---|---|---|
| Separator | `,` | `\t` |
| Quoting | RFC 4180 (double-quote) | Same rules; kicks in only for cells with embedded tabs / newlines / quotes |
| Header row | Identical column names + order | Identical column names + order |
| Score-cell convention | Empty = unscored, `"0"` = observed dissimilar | Same |
| `--similarity-labels` append | 7th `score_kind` column | Same |

Only the separator changes — the column order, header row, and
empty-vs-`0` score convention are identical across both dialects.
The output file extension is **not** auto-rewritten: pass an
explicit `.tsv` path when you want one. Cells containing tabs,
newlines, or quotes still round-trip safely under TSV — the
stdlib dialect quotes them — so the export is lossless either
way.

## Per-kind score labels (`--similarity-labels`)

`score` is an integer, but its *meaning* depends on the row's kind:
for an `R` row it's how alike the two paths are (100 = byte-identical
move); for a `C` row it's how much of the source survived in the
copy (100 = verbatim duplicate). The two are not directly comparable
— a 90 % rename and a 90 % copy describe different observations.

`--similarity-labels` (opt-in, requires `--with-similarity`) makes
that distinction explicit by attaching a canonical `score_kind`
discriminator to every R/C row. The vocabulary is closed:

| `score_kind` | Applies to | Score semantics |
|---|---|---|
| `rename-similarity` | `kind=R`, integer score | How alike the two paths are. |
| `copy-similarity` | `kind=C`, integer score | How much of the source survived. |
| `unscored` | `kind=R` or `kind=C`, `score=null` | Kind is meaningful, magnitude isn't. |
| *(absent / empty)* | plain A/M/D rows | No rename/copy provenance to label. |

Surfaces touched:

- **JSON audit** — adds `score_kind` to the nested `similarity`
  object. Plain rows whose `similarity` is `null` get no label
  (absence already means "no provenance"). The legacy schema is
  preserved byte-for-byte when the flag is off.
- **Text table** — appends a `meaning` column after `old`. Plain
  rows render `-` for visual consistency with the other similarity
  cells.
- **CSV export** — appends a 7th `score_kind` column. The first six
  columns are unchanged so positional readers that hard-code indices
  0–5 keep working without modification; opt into index 6 only when
  you care about the discriminator.

Score-of-`0` is classified by kind, **not** as `unscored` —
`unscored` is reserved for `score=null` (authored payloads without a
percentage). A `C` row with `score=0` is `copy-similarity` (git
observed the pair and rated them entirely dissimilar), not
`unscored`.