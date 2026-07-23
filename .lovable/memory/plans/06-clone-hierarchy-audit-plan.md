# Plan 06 — Clone Hierarchy Audit & Hardening

> **Status:** drafted (awaiting Phase 0 spec sign-off)
> **Owner:** AI
> **Scope:** `gitmap clone <file>` for `.json` and `.csv` manifests
> **Out of scope:** direct-URL clones, `.txt` (`git clone …`) manifests, `clone-next` family

## Goal

Confirm — and then *prove with regression tests* — that
`gitmap clone manifest.json` and `gitmap clone manifest.csv` recreate the
exact folder hierarchy described by each record's `RelativePath`, under
`--target-dir`, no matter how the manifest was authored (mixed slashes,
nested groups, leading `./`, Windows paths, missing fields, etc.).

The current code path *should* do this already (see
`gitmap/cloner/cloner.go::cloneOne` joining `targetDir` with
`rec.RelativePath`), but there is no end-to-end test that asserts the
on-disk tree, and the parsers do no path normalisation. This plan closes
both gaps without introducing a new sub-command.

## Why now

User report: "implement a clone command that reads JSON and CSV and
preserves the original folder hierarchy." Investigation showed the
feature already exists at the code level but is under-specified and
under-tested — easy to silently regress. Audit + harden is the right
response.

---

## Phase 0 — Spec & open questions

Files to add/update:

- `spec/01-app/110-clone-hierarchy-guarantee.md` (new) — codify the
  guarantee, the normalisation rules, and the failure modes.
- Cross-link from `spec/01-app/05-cloner.md` (one-line "see also" footer).
- Memory: `mem://features/clone-hierarchy-guarantee.md` summarising the
  contract for future sessions.

### Open questions (need user defaults)

| # | Question | Proposed default |
|---|----------|------------------|
| Q1 | Path separator normalisation: should `\` in CSV/JSON `RelativePath` be rewritten to `/` on Unix and vice versa on Windows? | **Yes**, normalise to `filepath.FromSlash` at parse time. |
| Q2 | Reject vs. sanitise absolute paths (`/foo`, `C:\foo`) and parent escapes (`../../etc`) inside `RelativePath`? | **Reject** with a per-row error; the rest of the manifest still runs. Mirrors G305 hardening already in place for zip extraction. |
| Q3 | Empty `RelativePath` — fall back to `RepoName`, or treat as a manifest error? | **Fall back to `RepoName`**; if both are empty, error that row. Preserves current behaviour for hand-written manifests. |
| Q4 | Should the audit add a `--dry-run` that prints the planned tree without cloning? | **Yes** (cheap, very useful for hierarchy debugging). Behind a flag, default off. |

Plan moves to Phase 1 only after the user confirms or overrides Q1–Q4.

---

## Phase 1 — Parser hardening (`gitmap/formatter/`)

1. Add `normaliseRelativePath(string) (string, error)` in a new
   `formatter/relpath.go`:
   - `filepath.FromSlash` → OS-native separators.
   - `filepath.Clean` to collapse `.` and `//`.
   - Reject if result is absolute (`filepath.IsAbs`) or starts with `..`.
   - Reject Windows drive letters on any OS (`C:\…`).
2. Apply it in `ParseCSV` (`rowToRecord`) and `ParseJSON` (post-decode
   loop). Per-row failures bubble up as a `ManifestError` with row index
   and field — added to `model.CloneSummary.Errors` rather than aborting
   the run.
3. Tests in `gitmap/formatter/relpath_test.go` covering each Q2 case +
   mixed slashes + non-ASCII path segments.

## Phase 2 — Cloner contract tests (`gitmap/cloner/`)

End-to-end tree assertion using a stub `git` binary on `$PATH` that just
`mkdir`s the destination + writes a sentinel file. No network.

- `cloner/hierarchy_test.go` with table-driven cases:
  - Flat (`a`, `b`, `c`).
  - Nested two-deep (`group/repo`).
  - Mixed depths in one manifest.
  - Windows-style separators in input on Linux runner.
  - Empty `RelativePath` → falls back to `RepoName` per Q3.
  - Malicious `../../escape` row → recorded as failure, sibling rows
    succeed.
- After each run, walk `targetDir` and compare to the expected tree
  (sorted). Fail loudly on any extra/missing directory.

## Phase 3 — `--dry-run` (gated on Q4 = yes)

- New flag in `parseCloneFlags`. When set:
  - Skip `git` invocation entirely.
  - Print the planned tree as a sorted, indented listing.
  - Still validate paths through Phase 1 helpers — surfaces manifest
    errors without touching disk.
- Test in `cmd/clone_dryrun_test.go`: golden file of the tree printout.

## Phase 4 — Docs, helptext, changelog

- Update `gitmap/helptext/clone.md` with the hierarchy guarantee, the
  normalisation rules, and the new `--dry-run` flag.
- Append a changelog entry in both `CHANGELOG.md` and
  `src/data/changelog.ts` (Plan-06-Clone-Hierarchy-Hardening).
- Update spec `05-cloner.md` "Behavior → File-based clone" with explicit
  reference to spec 110.

## Phase 5 — QA pass

1. `go build ./...`
2. `go test ./gitmap/formatter/... ./gitmap/cloner/... ./gitmap/cmd/...`
3. `bash .github/scripts/check-constants-naming.sh`
4. `golangci-lint run` (pinned v1.64.x as per `mem://tech/static-analysis-security`).
5. Manual smoke:
   ```
   gitmap clone testdata/hier.json --target-dir /tmp/h --dry-run
   gitmap clone testdata/hier.csv  --target-dir /tmp/h
   tree /tmp/h | diff - testdata/hier.tree.golden
   ```

---

## Definition of done

- All tests in Phases 1–3 green.
- Spec 110 + memory file in place; cross-links updated.
- `gitmap clone <file>` documented as a hierarchy-preserving operation
  with explicit normalisation and rejection rules.
- A future session reading `mem://features/clone-hierarchy-guarantee`
  can describe the contract without opening the source.

## Risks / non-goals

- We are **not** redesigning the manifest schema — `ScanRecord` stays
  as-is. If new fields are wanted, that's a separate plan.
- We are **not** changing direct-URL or text-manifest behaviour.
- Stubbing `git` via `$PATH` shim is the only portable way to assert the
  on-disk tree without network; CI must allow `PATH` mutation in tests
  (already the case).
