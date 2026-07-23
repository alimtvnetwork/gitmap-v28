# What to Read — AI Onboarding Map

> Purpose: tell a fresh AI session **exactly which files to open** to understand
> the project, the JSON output contracts, and how to add new schemas / outputs
> (called "slides" loosely — any new typed JSON surface).
>
> Read order matters. Top to bottom.

---

## 1. Project Orientation (read first, in order)

| Order | File | Why |
|---|---|---|
| 1 | `README.md` (root) | High-level product, install, command surface, project structure tree |
| 2 | `.lovable/memory/index.md` | Master memory index — Core rules + topic links |
| 3 | `.lovable/overview.md` | Project overview & invariants |
| 4 | `.lovable/strictly-avoid.md` | Hard prohibitions (do NOT violate) |
| 5 | `.lovable/coding-guidelines.md` | Code style + size limits |
| 6 | `.lovable/plan.md` | Active roadmap |
| 7 | `spec/03-general/10-strictly-prohibited.md` | Spec-level forbidden actions |

---

## 2. Folder Map (where things live)

```
/                                      repo root
├── README.md                          product README + pinned version
├── CHANGELOG.md                       version-by-version changes
├── gitmap/                            Go CLI source (the product)
│   ├── cmd/                           command handlers + contract tests
│   ├── constants/                     ALL string constants (no magic strings)
│   ├── model/                         shared structs (JSON record shapes)
│   ├── stablejson/                    deterministic JSON writer (key-by-key)
│   ├── store/                         SQLite layer + migrations
│   ├── completion/                    shell completion generators
│   ├── release/                       release workflow + semver
│   ├── formatter/                     human/JSON/CSV output formatters
│   ├── helptext/                      embedded markdown help files
│   └── scripts/                       embedded install/uninstall scripts
├── gitmap-updater/                    standalone updater binary
├── spec/                              specifications (source of truth)
│   ├── 03-general/                    cross-cutting rules (logging, build, prohibited)
│   ├── 04-generic-cli/                per-command specs
│   ├── 05-coding-guidelines/          Go style, naming, errors
│   ├── 08-json-schemas/               ★ JSON Schema definitions for every JSON output
│   └── 12-consolidated-guidelines/    merged style guide
├── scripts/                           helper shell/PowerShell scripts
├── .github/workflows/                 CI pipelines
├── .gitmap/                           runtime artifacts (release/, backup/, etc.) — DO NOT hand-edit
├── .lovable/                          AI memory + plans + prompts
│   ├── memory/                        durable knowledge (index, features, tech, project)
│   ├── plans/                         active multi-step plans
│   ├── prompts/                       reusable prompts (write-memory, read-memory)
│   ├── pending-issues/                open issues
│   ├── solved-issues/                 closed issues with root cause
│   └── cicd-issues/                   CI/CD-specific issues
└── src/                               React docs site (Vite + shadcn)
```

---

## 3. How JSON Outputs ("slides") Are Structured

Every JSON-emitting command in gitmap follows the same triangle:

```
spec/08-json-schemas/<name>.schema.json    ← contract (JSON Schema draft-07)
        │
        ├── gitmap/model/<name>.go         ← Go struct mirroring the schema
        │
        ├── gitmap/cmd/<name>render.go     ← stablejson encoder (key order = wire contract)
        │
        └── gitmap/cmd/<name>_jsonschema_contract_test.go  ← drift guard:
                                              schema ↔ encoder ↔ golden bytes
```

### Canonical example — `amend audit`

| Layer | File |
|---|---|
| Schema | `spec/08-json-schemas/amend-audit.schema.json` |
| Model  | `gitmap/model/amendment.go` (`AmendmentRecord`, `AmendAuthor`) |
| Encoder | `gitmap/cmd/amendauditrender.go` (`encodeAmendAuditJSON` + `amendAuditKey*` constants) |
| Writer | `gitmap/cmd/amendaudit.go` (`writeAmendAudit` → file on disk) |
| Schema contract test | `gitmap/cmd/amendaudit_jsonschema_contract_test.go` |
| Golden-bytes contract test | `gitmap/cmd/amendauditjson_contract_test.go` |
| Golden fixture | `gitmap/cmd/testdata/amend_audit_canonical.json` |

### Existing JSON surfaces (all live in `spec/08-json-schemas/`)

`amend-audit`, `amend-list`, `bookmark-list`, `diff-profiles`, `export`,
`find-next`, `help-json`, `history`, `latest-branch`, `list-releases`,
`list-releases-all-repos`, `list-versions`, `llm-docs`, `probe-report`,
`project-repos`, `scan-output`, `scan-project`, `ssh-list`, `startup-list`,
`stats`, `temp-release-list`, `version-history`, `watch`.

Each `.schema.json` declares `type`, `required[]`, and `properties{}` —
the contract tests assert the live encoder emits exactly those keys in
the order declared by the matching `<name>Key*` constants.

---

## 4. How to Add a New JSON Output (Recipe)

1. **Author the schema** → `spec/08-json-schemas/<new-name>.schema.json`
   (draft-07, list `required` alphabetically, declare every `properties` key).
2. **Add the Go struct** → `gitmap/model/<new-name>.go`.
3. **Add wire-key constants** → top of `gitmap/cmd/<new-name>render.go`
   (`<newName>Key<Field> = "<jsonKey>"`). Order of these constants IS the wire order.
4. **Write the encoder** using `gitmap/stablejson.WriteObject` with one
   `stablejson.Field{Key, Value}` per constant. Pre-render nested objects /
   arrays via `WriteObjectIndent` / `WriteArrayIndent` and embed as
   `json.RawMessage` — never call `json.MarshalIndent`.
5. **Add a JSON Schema contract test** → `gitmap/cmd/<newname>_jsonschema_contract_test.go`
   modelled on `amendaudit_jsonschema_contract_test.go`:
   - assert `type == "object"`
   - assert sorted `required` matches a hard-coded slice
   - run the encoder and assert every emitted key is in `properties`.
6. **Add a golden-bytes test** → `gitmap/cmd/<newname>json_contract_test.go`
   modelled on `amendauditjson_contract_test.go`, using
   `assertGoldenBytesDeterministic` + `assertSchemaKeysFirstObject`.
7. **Generate the golden fixture**:
   ```
   GITMAP_UPDATE_GOLDEN=1 GITMAP_ALLOW_GOLDEN_UPDATE=1 \
     go test ./gitmap/cmd -run <NewName>JSONContract
   ```
8. **Wire the command** → add CLI handler in `gitmap/cmd/`, register in
   `dispatch`, add help in `gitmap/helptext/`, add constants ID in
   `gitmap/constants/constants_cli.go`.
9. **Bump version** (minor for new feature) → `gitmap/constants/constants.go`
   `Version`, `src/constants/index.ts`, `CHANGELOG.md`, pin in `README.md`.

---

## 5. Constants & Magic-String Discipline

- **Never** inline a user-visible string. All go in `gitmap/constants/constants_*.go`.
- CLI IDs (command names + aliases) live **only** in `constants_cli.go`.
- Domain-specific bundles: `constants_cd.go`, `constants_clone.go`, etc.
- Error format strings end in `Fmt` or start with `Err`. See existing files.

---

## 6. Tests You Must Run Before Saying "Done"

```
nix run nixpkgs#go_1_24 -- test ./gitmap/... -count=1
nix run nixpkgs#go_1_24 -- vet ./gitmap/...
```

For JSON contract changes also run the targeted contract suite:
```
nix run nixpkgs#go_1_24 -- test ./gitmap/cmd -run 'JSONContract|JSONSchema'
```

---

## 7. When Stuck — Memory Cross-References

- Encoder pattern & key-order rule → `mem://features/stablejson-usage` (if absent, see `amendauditrender.go`).
- Contract-test pattern → `gitmap/cmd/amendaudit_jsonschema_contract_test.go`.
- Version bump procedure → `.lovable/memory/project/version-bump-procedure.md`.
- Strictly prohibited actions → `.lovable/memory/constraints/strictly-prohibited.md`.
- Code style limits → `.lovable/memory/style/code-constraints.md`.

---

## Changelog

- v1 — initial map. Created in response to "what to read" request so any fresh
  AI session can locate the JSON-output contract triangle and add new outputs
  without spelunking.
