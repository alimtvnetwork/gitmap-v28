---
name: Fixture Stamp SHA Hash
description: Optional sha=<12-hex> field on `// fixture-stamp:` markers detects content drift; ValidateBody fails if body hash diverges; autobump refreshes both generation and sha
type: feature
---
# Fixture Stamp SHA Hash

The `// fixture-stamp:` marker accepts an optional `sha=<12-hex>` field — the first 12 hex chars of SHA-256 over the fixture body **with all fixture-stamp lines stripped** (so the hash is stable across marker edits).

## Public surface (`gitmap/fixtureversion/hash.go`)

- `BodyHashExcludingMarker(body) string` — full 64-hex SHA-256 of body minus stamp lines.
- `ShortHash(full) string` — first `HashShortLen` (12) chars.
- `HashMatches(body, recordedShortHash) bool` — empty recorded hash = opt-out (returns true).
- `RewriteOrAppendSHA(head, newSHA) string` — replaces `sha=` if present, appends to marker line if missing.

## Validation

- `Stamp.SHA string` — optional field; empty means "do not enforce".
- `ValidateBody(body, stamp, want) error` — runs `Validate` then checks the hash. Failure message: `fixture %q body hash mismatch: marker records sha=X but actual body sha=Y`.
- `MustValidateBody` now calls `ValidateBody` (was `Validate`), so any test using stamped fixtures gets drift detection automatically.

## Autobump integration

- `BumpRequest.NewSHA` field added.
- `MustValidateBodyWithAutobump` computes `ShortHash(BodyHashExcludingMarker(body))` and passes it via `BumpRequest`. Under `GITMAP_FIXTURE_AUTOBUMP=1` both generation AND sha are refreshed in a single in-place rewrite of the source file.

## Migration

Existing markers without `sha=` keep validating (opt-out). To enable drift detection on a fixture, add `sha=<12-hex>` to its marker — the first `make fixtures-bump RUN=...` after that will keep it fresh.

First consumer: `fixRepoV9ToV12FixtureBody` in `gitmap/cmd/fixrepo_rewrite_v9tov12_test.go` (sha=7e1463d1eae6).

## Why hash-with-marker-stripped

If the hash were computed over the body including the marker, every successful bump would invalidate the hash it just wrote — chicken-and-egg. Stripping the stamp line(s) before hashing makes the value stable across any number of marker rewrites.
