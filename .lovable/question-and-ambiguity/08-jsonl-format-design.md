# 08 — `--format=jsonl` design choices

## Original task

> Add a `--format=jsonl` mode to output one startup entry per line while keeping stable key order within each object.

## Ambiguity

Three points needed a call:

1. **Empty-list framing.** Three plausible options for an empty list:
   (a) zero bytes; (b) a single `\n`; (c) `[]\n` (mirroring `--format=json`).
2. **Trailing newline on the last record.** Some JSONL producers omit it (so `cat a.jsonl b.jsonl` produces invalid output), others include it (so concatenation stays valid).
3. **Inter-token whitespace.** "Compact" can mean strictly no spaces (`{"k":v,"k2":v2}`) or "no newlines but spaces after `:` and `,`". The existing `--format=json` is pretty-printed with 2-space indent; jsonl needs to pick its own style.

A fourth, smaller ambiguity: where should the JSONL builder live — duplicate the Field-slice construction inside the renderer, or factor a shared helper between `--format=json` and `--format=jsonl`?

## Options considered

| # | Question | Option A | Option B | Chosen |
|---|----------|----------|----------|--------|
| 1 | Empty list framing | `\n` (one blank line) | Zero bytes | **B**: matches jq `--compact-output` behavior; `wc -l` equals record count |
| 2 | Trailing newline | Omit on last record | Always include | **Always include**: `cat a.jsonl b.jsonl` stays valid |
| 3 | Inter-token spaces | `{"k": v, "k2": v2}` | `{"k":v,"k2":v2}` | **No spaces**: maximally compact, matches `jq -c` exactly |
| 4 | Field construction | Duplicate per format | Shared helper | **Shared helper** (`buildStartupListJSONItems`): one diff for any future column add/rename |

## Recommendation & decision taken

Implemented all "B" / shared-helper options. Concrete deliverables:

1. **`gitmap/stablejson/stablejson.go`**: added `WriteJSONLines(w, items)` and `writeCompactObject` helper. Empty input writes zero bytes; non-empty writes `{...}\n` per record with no inter-token whitespace.
2. **`gitmap/cmd/startuplistrender.go`**: extracted `buildStartupListJSONItems` as the shared source of (field name, field order, value); both `encodeStartupListJSON` and the new `encodeStartupListJSONL` route through it.
3. **`gitmap/cmd/startup.go`**: added `constants.StartupListFormatJSONL` to the `parseStartupListFlags` switch.
4. **`gitmap/constants/constants_startup.go`**: added `StartupListFormatJSONL = "jsonl"`, updated `FlagDescStartupListFormat` and `ErrStartupListBadFormat` to mention the new format.
5. **`gitmap/cmd/startuplistjsonl_contract_test.go`** (new, 190 lines, under the 200-line budget): five tests covering empty-emits-nothing, single-entry byte-exact, multi-entry line count, key-order-stable (via `json.Decoder.Token()`), and special-chars-match-json (cross-format value parity).
6. **`gitmap/helptext/startup-list.md`**: added a `--format=jsonl` section with concrete sample output (113 lines, under the 120-line cap).
7. **Version bump**: `3.166.0` → `3.167.0` (minor — new user-facing flag value).

## Verification

- `go build` for `GOOS=windows ./...` → clean (full build including the JSONL test file via the `cmd` package).
- `go test ./stablejson/` on linux → pass.
- Linux `./...` build fails ONLY in pre-existing `startup/winbackend.go` and `startup/winshortcut.go` which lack `//go:build windows` tags (unrelated to this change — flagged separately for a future cleanup task).
- `gofmt` reformatted the constants block alignment; no semantic changes.
- No existing test pins `ErrStartupListBadFormat` or `FlagDescStartupListFormat` text — updating both was safe.

## Counter

Task 08 of 40. One follow-up worth flagging (not done in this turn): the `startup/win{backend,shortcut}.go` files leak Windows-only symbols into linux/darwin builds because they're filename-suffixed but not `//go:build`-tagged. Pre-existing, unrelated to JSONL, but blocks any non-Windows `go test ./cmd/` run.
