# gitmap help

Show the gitmap command reference. Without arguments, prints the full
grouped help screen. With a topic, prints per-command help.

## Usage

    gitmap help                    # full grouped help
    gitmap help <command>          # per-command help (e.g. gitmap help clone)
    gitmap help --compact          # one-line-per-command dense layout
    gitmap help --groups           # show only the group banners
    gitmap help --filter <query>   # show only rows matching <query>
    gitmap help -f <query>         # short alias for --filter
    gitmap help --json             # machine-readable JSON of the help registry

## Flags

- `--compact`         Dense single-column layout (no group banners).
- `--groups`          Print only the intent-based group headings.
- `--filter <query>`  Case-insensitive substring filter over command + description.
                      Matches are highlighted; zero-hit queries print fuzzy suggestions.
- `-f <query>`        Short form of `--filter`.
- `--json`            Emit the help registry as JSON. ANSI color is stripped.
                      Composes with `--filter` to scope the output.

## Examples

### Filter to a single area

    gitmap help --filter clone
    gitmap help -f release

### List only the group banners

    gitmap help --groups

### Pipe machine-readable help to jq

    gitmap help --json | jq '.groups[] | {group, count: (.lines|length)}'
    gitmap help --json --filter ssh | jq '.count'

## JSON schema

The `--json` payload conforms to [`spec/08-json-schemas/help-json.schema.json`](../../spec/08-json-schemas/help-json.schema.json)
(JSON Schema draft 2020-12). Contract test `helpjson_jsonschema_contract_test.go`
validates runtime output against the schema on every build to prevent drift.

## Notes

- Per-command help (e.g. `gitmap help clone`) accepts a `--pretty` flag
  to toggle the styled vs plain renderer.
- Glyph rendering follows `--glyphs auto|rich|safe` (see `gitmap help glyphs`).

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter help
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
