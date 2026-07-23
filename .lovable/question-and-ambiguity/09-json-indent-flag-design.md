# 09 — `--json-indent` flag design

## Original task

> Add a `--json-indent` flag to let me choose between minified and pretty-printed stable JSON output while preserving byte-identical ordering.

## Ambiguities

1. **Flag type.** Boolean (`--minified`)? Enum (`compact|pretty`)? Integer (spaces per level)?
2. **Default value.** 0 (preserve nothing) vs. 2 (preserve current pretty default).
3. **Behavior on non-json formats.** Error / warn / silently ignore for `--format=jsonl|csv|table`?
4. **Backward compatibility.** Does `encodeStartupListJSON(w, entries)` keep its 2-arg signature (12 test call sites) or get a new third parameter?
5. **Indent character.** Spaces only, or also tabs (`--json-indent=tab`)?

## Decisions taken

| # | Choice | Rationale |
|---|--------|-----------|
| 1 | **Integer 0..8** | Matches `python -m json.tool --indent`, `jq --indent`. 0 = minified is the conventional "compact" sentinel. Bool would force a separate flag for indent width. Enum would lock out 4-space teams. |
| 2 | **Default = 2** | Preserves the long-standing pretty output exactly, so all 12+ existing JSON golden fixtures pass without regeneration. |
| 3 | **Silently ignore for non-json formats** | Lets shell scripts pass `--json-indent=N --format=$F` for any `$F` without conditionals. Documented in helptext + flag description. |
| 4 | **Two-function split**: `encodeStartupListJSON(w, entries)` (legacy, 2-arg, default indent) + new `encodeStartupListJSONIndent(w, entries, n)`. Dispatcher uses the indent variant. | Zero churn in existing tests; clean migration path. |
| 5 | **Spaces only for now** | Tabs would need a `string`-typed flag; YAGNI. `indentSpaces(n)` helper centralizes the conversion so a future `--json-indent=tab` lands in one place. |

## Implementation summary

- **`gitmap/stablejson/stablejson.go`**: new public `WriteArrayIndent(w, items, indent string)`; `WriteArray` becomes a 1-line wrapper at `indent="  "` (preserves byte-compat with `json.Encoder.SetIndent("", "  ")`).
- **`gitmap/stablejson/writers.go`** (new, 124 lines): unexported `writeArrayPretty`, `writeArrayMinified`, `writeObject(_, _, indent)`, `writeCompactObject`, `writeKeyValue` — split out to keep `stablejson.go` under the 200-line budget.
- **`gitmap/cmd/startuplistrender.go`**: added `encodeStartupListJSONIndent` + `indentSpaces` helper; legacy `encodeStartupListJSON` becomes a thin default-indent wrapper. Dispatcher signature gained `jsonIndent int`.
- **`gitmap/cmd/startup.go`**: `parseStartupListFlags` now returns `(format string, jsonIndent int, err error)`. Flag validated 0..8 even when format ignores it.
- **`gitmap/constants/constants_startup.go`**: new `FlagStartupListJSONIndent`, `FlagDescStartupListJSONIndent`, `StartupListJSONIndentDefault=2`, `StartupListJSONIndentMax=8`, `ErrStartupListBadJSONIndent`.
- **`gitmap/cmd/startuplistjson_indent_contract_test.go`** (new, 190 lines): five tests covering key-order-stable across {0,1,2,4,8}, indent=0 byte-exact, indent=2 vs legacy byte-equal, empty-list `[]\n` invariant, indent-width-is-N-spaces, plus flag-parsing accept/reject boundaries.
- **`gitmap/helptext/startup-list.md`** (118 lines, under 120 cap): documented flag, examples, ignore-for-non-json behavior.
- **Version**: 3.167.0 → **3.168.0** (minor — new user-facing flag).

## Verification

- `GOOS=windows go build ./...` → clean.
- `go test ./stablejson/` on linux → all pre-existing tests pass (byte-compat-with-Encoder included).
- Pre-existing `startup/win{backend,shortcut}.go` linux-build leak and `collectObjectKeys` redeclaration in cmd test helpers are still present — both pre-date this change (logged in entries 02 and 04 respectively).

## Counter

Task 09 of 40.
