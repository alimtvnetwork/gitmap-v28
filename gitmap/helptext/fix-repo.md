# gitmap fix-repo

Rewrite prior `{base}-vN` versioned-repo-name tokens in every tracked
text file to the current version. Go-native re-implementation of
`fix-repo.ps1` / `fix-repo.sh` with byte-identical exit codes and
config schema.

## Synopsis

```
gitmap fix-repo [-2 | -3 | -5 | --all] [--dry-run] [--verbose] [--strict] [--config <path>]
gitmap fr                                                                # short alias
```

PowerShell-style flags (`-DryRun`, `-Verbose`, `-Strict`, `-Config <p>`,
`-All`) are also accepted.

## Behavior

1. Read repo identity from `git`. Repo name must end with `-vN`.
2. Default mode rewrites the last 2 prior versions. `-3` / `-5`
   widen the window; `--all` rewrites every prior version.
3. Enumerate tracked files via `git ls-files`. Skip ignored paths,
   reparse points, > 5 MiB files, binary extensions, and files with
   a NUL byte in the first 8 KiB.
4. Replace `{base}-vN` with `{base}-v<current>` (negative-lookahead
   guard so `-v1` never matches inside `-v10`).
5. Print a summary; in `--dry-run` no file is written and each
   would-be-modified file gets a `[dry-run]` preview line listing the
   total replacements plus a per-rule breakdown (e.g.
   `v1×3, v2×1, bare×2`). The breakdown surfaces every numbered
   `{base}-vN` target plus the v1→v2 `bare` sweep distinctly so you
   can vet the plan without `--verbose`.
6. **Strict mode (`--strict`)**: after the rewrite + `gofmt -w` step,
   derive the unique set of touched Go packages from the modified
   `.go` files and run `go test ./pkgA ./pkgB …`. Catches semantic
   desyncs the byte-level rewriter cannot see — e.g. a hard-coded
   sibling literal that drifted from its `{base}-vN` neighbor across
   a width-crossing bump (the v9→v10/v12 failure mode closed by
   v4.12.0). Skips safely when `go` is not on PATH or when no `.go`
   file was modified, so `--strict` is safe to leave on across mixed
   environments. Test failure exits with code **9**.

## Example

```
$ gitmap fix-repo --dry-run --verbose
fix-repo  base=myrepo  current=v3  mode=-2
targets:  v1, v2
host:     github.com  owner=acme

  [dry-run] would rewrite README.md (4 replacements): v1×3, v2×1
  [dry-run] would rewrite docs/install.md (1 replacements): v2×1

scanned: 87 files
changed: 2 files (5 replacements)
mode:    dry-run
```

## Before / after by current version

The bare-base sweep only fires on the v1→v2 bump. Use `--restrict no-version`
(`-r nv`) to suppress it even there.

### Inside `gitmap-v27` (current=v2, target includes v1) — bare sweep ACTIVE

```
BEFORE                                      AFTER (gitmap fix-repo)
gitmap          → script body                gitmap-v27
gitmap-v27       → install URL                gitmap-v27
gitmap-v27       → already current            gitmap-v27  (no-op)
gitmap.js       → filename, word-boundary    gitmap.js  (skipped)
```

With `--restrict no-version` (`-r nv`):

```
BEFORE                                      AFTER (gitmap fr -r nv)
gitmap          → bare token                 gitmap     (PRESERVED)
gitmap-v27       → versioned                  gitmap-v27
```

### Inside `gitmap-v27` or higher (current≥v3) — bare sweep SKIPPED

```
BEFORE                                      AFTER (gitmap fix-repo --all)
gitmap                  → binary / brand     gitmap         (PRESERVED)
https://…/owner/gitmap  → upstream URL       …/owner/gitmap (PRESERVED)
gitmap-v27               → prior version      gitmap-v<cur>
gitmap-v27               → prior version      gitmap-v<cur>
gitmap-v27              → unrelated future   gitmap-v27     (negative-lookahead)
```

Rule of thumb: once you've shipped past v2, bare `{base}` is almost always
the binary, package, or brand name — never the pre-versioned remote — so
the rewriter leaves it alone.


## Exit codes

`0` ok / `2` not-a-repo / `3` no-remote / `4` no-version-suffix /
`5` bad-version / `6` bad-flag / `7` write-failed / `8` bad-config /
`9` tests-failed (`--strict` only).

See `spec/04-generic-cli/27-fix-repo-command.md` for the full spec.

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter fix-repo
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).

## Examples

```bash
gitmap fix-repo
```
