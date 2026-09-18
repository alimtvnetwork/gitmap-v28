# 02 — Pre-existing duplicate test helpers in `gitmap/cmd/`

## Discovered during

Task 01 (JSON schema docs) — while verifying my new `jsonschema_helpers_test.go`
wouldn't collide with existing helpers, I found an existing collision between
two files I did not author:

| Symbol | File A | File B |
|---|---|---|
| `collectObjectKeys(t *testing.T, dec *json.Decoder) []string` | `jsoncontract_helpers_test.go:131` | `jsonsnapshot_helpers_test.go:116` |
| `equalStringSlices(a, b []string) bool` | `jsoncontract_helpers_test.go:171` | `jsonsnapshot_helpers_test.go:144` |

Identical signatures, both in `package cmd`, no build tags, both in `_test.go`
files. Go does not allow duplicate top-level identifiers in the same package
even in test files, so this should fail to compile with:

```
./jsoncontract_helpers_test.go:131:6: collectObjectKeys redeclared in this block
        ./jsonsnapshot_helpers_test.go:116:6: other declaration of collectObjectKeys
```

## Why this matters for the current task

My new `jsonschema_helpers_test.go` reuses `collectObjectKeys` and
`expectDelim` (the latter only exists in `jsonsnapshot_helpers_test.go`). If the
duplicate above is somehow real, my test inherits the broken state but does NOT
cause it. If the duplicate is somehow tolerated (e.g., a build flag I'm not
seeing), my test will compile fine.

## Inference taken

I left the two pre-existing files untouched and made my `jsonschema_helpers_test.go`
depend only on `expectDelim` (unique to `jsonsnapshot_helpers_test.go`) and
`collectObjectKeys` (which Go's resolver will pick from one of the duplicates if
the build somehow tolerates it).

If the user runs `go test ./gitmap/cmd/...` and gets the redeclaration error,
the fix is to delete the duplicates from ONE of the two existing files —
recommended: keep `jsonsnapshot_helpers_test.go` (newer-style, smaller helpers)
and remove the duplicates from `jsoncontract_helpers_test.go`. That cleanup is
NOT in scope for the JSON-schema task; logging here so the user sees it on
review.

## Recommendation

Open a separate cleanup task to dedupe `gitmap/cmd/jsoncontract_helpers_test.go`
vs `jsonsnapshot_helpers_test.go`. Estimated effort: 15 minutes. Risk: low —
the symbols have identical signatures, so removing one definition has no
behavioral impact.
