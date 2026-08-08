# JSON Schema docs for `gitmap-v28` CLI outputs

This directory holds **machine-readable JSON Schema documents** for every CLI
command that emits stable JSON. Downstream consumers (jq pipelines, GitHub
Actions, monitoring agents, internal tools) can use these schemas to:

1. **Validate** that a `gitmap-v28 <cmd> --json` output conforms to the documented
   shape before parsing it.
2. **Pin** the expected key set so a future field rename is caught at the
   consumer's CI step instead of in production.
3. **Reproduce** key ordering — important when consumers do byte-level diffing
   (e.g. golden-file tests in downstream repos).

## Dialect

All schemas in this directory use **JSON Schema draft 2020-12**
(`"$schema": "https://json-schema.org/draft/2020-12/schema"`).

## The `propertyOrder` extension

JSON Schema deliberately does NOT standardize key ordering — JSON objects are
unordered per RFC 8259. However, `gitmap-v28`'s stablejson outputs DO emit keys in
a contractually fixed order (see `gitmap-v28/stablejson/stablejson.go`), and some
consumers rely on that for byte-identical diffs.

To express that contract in the schema we use a non-standard
`"propertyOrder": [...]` array on each object schema. Consumers that care about
order can read this array; standards-compliant validators ignore it (it's
treated as an unknown annotation per §6.4 of the 2020-12 spec).

## Coverage status

| Command | Schema | Contract | Notes |
|---|---|---|---|
| `gitmap-v28 startup-list --json` | [`startup-list.schema.json`](startup-list.schema.json) | **strict** | Backed by `gitmap-v28/stablejson`; key order is contractual |
| `gitmap-v28 list-releases --json` | [`list-releases.schema.json`](list-releases.schema.json) | **strict** | Per-repo view; lowerCamel keys mirror `model.ReleaseRecord` `json:` tags; pinned by `gitmap-v28/cmd/listreleases_jsonschema_contract_test.go` |
| `gitmap-v28 list-releases --all-repos --json` | [`list-releases-all-repos.schema.json`](list-releases-all-repos.schema.json) | **strict** | Joined Release+Repo view; PascalCase keys preserve the legacy `json.MarshalIndent` surface (`store.ReleaseAcrossRepos` has no `json:` tags) |
| `gitmap-v28 scan` (JSON + CSV)   | [`scan-output.schema.json`](scan-output.schema.json) + [`scan-output.README.md`](scan-output.README.md) | **field-name stable** | Field names/keys are permanent; new fields are appended only. Key ordering is NOT byte-contractual (uses `encoding/json`). |
| _(others)_ | — | — | See [`_TODO.md`](_TODO.md) — currently emit via `json.MarshalIndent` so ordering is NOT contractual until they migrate to `stablejson` |

## How a downstream consumer uses these

```bash
# Fetch the schema
curl -O https://raw.githubusercontent.com/.../spec/08-json-schemas/startup-list.schema.json

# Validate a real output (using ajv-cli as one example)
gitmap-v28 startup-list --json | ajv validate -s startup-list.schema.json -d -
```

## Schema authorship rules

1. **One schema file per CLI subcommand** that emits JSON. Filename:
   `<subcommand>.schema.json` (kebab-case, matches the CLI verb).
2. **Every object schema MUST list `"required": [...]`** for keys the consumer
   can rely on. Optional keys go in `"properties"` only.
3. **Every object schema SHOULD include `"propertyOrder": [...]`** matching the
   exact emit order. If ordering is not contractual, omit `propertyOrder` and
   add a `"$comment"` explaining why.
4. **`"additionalProperties": false`** by default. Use `true` only when the
   command may grow new top-level fields without a major version bump.
5. **Pair every schema with a contract test** in `gitmap-v28/cmd/` named
   `<subcommand>_jsonschema_contract_test.go` that:
   - embeds the schema,
   - emits a fixture via the actual encoder function,
   - asserts the output validates AND key order matches `propertyOrder`.

## Why not auto-generate from Go structs?

Auto-generation via reflection is tempting but lies about the
non-stablejson outputs — `encoding/json` does not guarantee field order, so a
generated schema with `propertyOrder` would document a contract that does not
exist. Hand-written schemas force us to be honest: only the stablejson outputs
get `propertyOrder`. Migration to stablejson is tracked in `_TODO.md`.
