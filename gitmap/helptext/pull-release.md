# gitmap pull-release

> **Pull first, then release** — sugar for `gitmap release` that runs
> `git pull` in the current repo (with your chosen merge strategy)
> *before* delegating to the standard release pipeline.

Renamed in **v5.6.0**. The legacy names `release-pull`, `relp`, and
`rlp` still work as aliases — they will not be removed.

---

## Aliases

| Form | Notes |
|------|-------|
| `pr` | **Primary short alias** (new in v5.6.0) |
| `pull-release` | Canonical long form |
| `release-pull` | Legacy long form — still routes here |
| `relp`, `rlp` | Legacy short aliases — still work |

## Synopsis

```
gitmap pull-release [--ff-only | --rebase | --merge] [--dry-run] [--verbose] \
                    [version] [release flags...]
gitmap pr           [...]
```

## Pull modes (mutually exclusive)

| Flag | Behavior |
|------|----------|
| `--ff-only` *(default)* | **Safe.** Fast-forward only. Hard-fails on any divergent history so we never tag on top of a divergent tree. |
| `--rebase` | Rebases your local commits on top of upstream. On conflict, auto-runs `git rebase --abort` and exits non-zero — your working tree is never left mid-rebase. |
| `--merge` | Classic merge. Passes `--no-rebase` so it overrides any user-level `pull.rebase=true` config. Creates a merge commit on divergence. |

## Other flags

| Flag | Behavior |
|------|----------|
| `--dry-run` | Print the `git pull` command that would run, then **skip the pull** and forward to release. |
| `--verbose` | Echo the git invocation to stderr before running it. |

All remaining args (version, `--bump`, `--bin`, `--draft`, `--dry-run`,
`-y`, …) are **forwarded verbatim** to `gitmap release`.

> **Note:** a top-level `--dry-run` here applies to the **pull step
> only**. `gitmap release`'s own `--dry-run` is forwarded separately
> and previews the release plan.

## Behavior

1. Verify the current directory is inside a git repository.
2. Resolve the pull mode (default `--ff-only`).
3. Run `git pull <mode-flag>` in cwd. On failure:
   - `--ff-only` / `--merge`: exit `1`.
   - `--rebase`: attempt `git rebase --abort` first, then exit `1`.
4. Delegate the remaining args to `runRelease`.

## Examples

```
# Default: safe fast-forward, then release v1.4.0
gitmap pr v1.4.0

# Allow divergence: rebase local commits onto upstream, then release
gitmap pr --rebase v1.4.0

# Classic merge (overrides pull.rebase=true), build binaries, draft
gitmap pr --merge v2.0.0 --bin --draft

# See exactly which `git pull` would run, skip it, preview the release
gitmap pr --rebase --dry-run --dry-run

# Legacy names still work
gitmap release-pull v1.4.0
gitmap relp --merge v2.0.0
```

## See also

- **`gitmap release`** — the underlying release pipeline.
- **`gitmap release-alias-pull`** — pull-then-release for a registered alias from any directory.
- **`gitmap fix-repo`** — rewrite `{base}-vN` tokens after bumping versions.

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter pull-release
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
