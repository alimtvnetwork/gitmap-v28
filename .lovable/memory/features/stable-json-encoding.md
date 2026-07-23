---
name: stable-json-encoding
description: gitmap/stablejson package — guaranteed-stable field order for consumer-facing JSON outputs, no reflection on struct shape
type: feature
---

# Stable JSON Encoding (`gitmap/stablejson`)

Tiny package providing `Field{Key, Value}` + `WriteArray(w, [][]Field)`.
Encodes JSON arrays of objects with field order pinned by the
caller's slice — never by struct field iteration.

## Why

`encoding/json` emits struct fields in declaration order today, but:
1. Go 2 / encoding/json/v2 has discussed alphabetical key ordering
2. Refactors / IDE tools silently reorder struct fields
3. Reflection-based field walks change with embedding/omitempty/tags

stablejson sidesteps all three by NEVER reflecting on a struct.
Only individual VALUES go through `json.Marshal`.

## Output contract (byte-compat with `json.Encoder.SetIndent("", "  ")`)

- 2-space indentation
- empty array → `[]\n` (NOT `null`)
- trailing `\n` (matches Encoder.Encode)
- preserves caller key order verbatim

Verified by `stablejson_test.go::TestWriteArray_ByteCompatWithEncoder`.
Migrating an existing caller does NOT require regenerating goldens.

## When to use

Any new `--format=json` CLI surface that downstream scripts will
parse. First adopter (v3.152.0): `gitmap startup-list --format=json`
via `encodeStartupListJSON` in `gitmap/cmd/startuplistrender.go`.

## When NOT to use

- Internal/debug JSON dumps where order doesn't matter
- Wide structs (40+ fields) like `model.ScanRecord` — listing every
  key by hand is more error-prone than the encoding/json risk it
  guards against. Use only for narrow, consumer-facing schemas.

## Field-name source of truth

Callers should declare on-the-wire names as package-level constants
(see `startupListJSONKeyName` / `_KeyPath` / `_KeyExec`) so any
rename is a single grep-friendly diff.
