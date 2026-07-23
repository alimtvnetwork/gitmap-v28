# Go Build Failure: copyFile Redeclared + Unused Import

## Ticket

After adding settings sync tools (OBS, VS Code, Windows Terminal), the
Go build fails with two errors:

```
cmd/update.go:204:6: copyFile redeclared in this block
    cmd/installobs.go:281:6: other declaration of copyFile
cmd/installvscode.go:11:2: "github.com/alimtvnetwork/gitmap-v27/gitmap-v27/constants" imported and not used
```

## Symptoms

1. `go build` fails with exit code 1.
2. CI release pipeline fails at the build step.

## Root Cause

### 1. copyFile redeclared in this block

Go requires all identifiers within a package to be unique. The `cmd`
package had `copyFile(src, dst string) error` defined in **two** files:

- `cmd/update.go:204` (original, used by the self-update flow)
- `cmd/installobs.go:281` (duplicate, added with the OBS settings sync)

Both functions had identical signatures and implementations. When the Go
compiler processes all files in the `cmd` package, it sees two declarations
of the same identifier and rejects the build.

This is a known constraint of Go's flat package namespace -- all `.go`
files in the same directory share a single scope, so helper functions
must be defined exactly once regardless of which file uses them.

### 2. Unused import

`installvscode.go` imported `github.com/alimtvnetwork/gitmap-v27/gitmap-v27/constants` but never
referenced any symbol from it. Go treats unused imports as compile errors.

## Fix

1. **Removed the duplicate `copyFile`** from `installobs.go`. The function
   in `update.go` is shared across the entire `cmd` package and is
   accessible from all files without any import.

2. **Removed the unused `constants` import** from `installvscode.go`.

## Prevention

1. Before adding a helper function to a new file, search the `cmd` package
   for existing functions with the same name:
   `grep -rn "^func <name>" gitmap-v27/cmd/`

2. Run `go build ./...` locally before committing to catch redeclaration
   and unused-import errors early.

3. See `memory/tech/go-namespace-constraints` for the full set of Go
   namespace rules enforced in this project.

## Related

- `spec/01-app/86-settings-sync.md` -- settings sync specification
- `memory/tech/go-namespace-constraints` -- Go namespace rules
