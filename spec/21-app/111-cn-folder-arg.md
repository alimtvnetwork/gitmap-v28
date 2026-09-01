# spec/01-app/111-cn-folder-arg.md — `gitmap cn` accepts a folder path

## Status

Proposed → implemented in v3.117.0.

## Motivation

Today `gitmap cn` only operates on the current working directory's repo:

    cd /repos/macro-ahk-v11
    gitmap cn v+1                   # → clones macro-ahk-v12

To `cn` a different folder a user must `cd` first, which breaks the
flow when scripting batch upgrades or operating from a launcher /
file-manager terminal where the cwd is fixed (e.g. `~`). Two new
forms eliminate the `cd` step:

    gitmap cn v+1 "/repos/macro-ahk-v11"   # explicit version, explicit folder
    gitmap cn "/repos/macro-ahk-v11"       # folder only — defaults to v++

Both forms reuse the existing `clonenext` pipeline; only the
dispatcher in `gitmap/cmd/clonenext.go` changes.

## Forms (additive — existing forms unchanged)

| Form | Semantics |
|---|---|
| `gitmap cn vX` / `cn v+1` / `cn v++` | **Existing**: in-place on cwd. |
| `gitmap cn <alias> vX` | **Existing**: cross-dir via release-alias resolver. |
| `gitmap cn vX <folder>` | **NEW**: chdir to `<folder>`, run `cn vX`, chdir back. |
| `gitmap cn v+1 <folder>` | **NEW**: same, with version-bump shortcut. |
| `gitmap cn <folder>` | **NEW**: chdir to `<folder>`, run `cn v++`, chdir back. |

`<folder>` may be a relative path, an absolute path, `~`-prefixed, or
a bare folder name resolved against the cwd. Quoted paths with spaces
are honored verbatim (the shell does the unquoting).

## Disambiguation rules

The dispatcher walks the positional args left-to-right and classifies
each token using the existing `looksLikeVersion` regex extended to
recognize `v++` and `v+N` shortcuts:

    versionPattern = ^(v?\d+(\.\d+\.\d+)?([.\-+].+)?|v\+\+|v\+\d+)$

Classification table for two-positional invocations:

| pos[0] | pos[1] | Resolution |
|---|---|---|
| version | folder-shaped | NEW — `cn version folder` |
| folder-shaped | version | EXISTING — `cn alias version` (cross-dir release-alias) |
| version | version | error: ambiguous (refuse, exit 1) |
| folder-shaped | folder-shaped | error: ambiguous (refuse, exit 1) |

For one positional:

| pos[0] | Resolution |
|---|---|
| version | EXISTING — in-place `cn version` |
| folder-shaped, exists on disk | NEW — `cn <folder>` defaulting to v++ |
| folder-shaped, NOT on disk | EXISTING — release-alias path resolution (kept for back-compat) |

A token is "folder-shaped" if it is NOT version-shaped AND
(`os.Stat` succeeds OR the token contains `/`, `\`, or starts with `~`).
The path-separator heuristic is the disambiguator that prevents
release-alias names like `gitmap` from being mistakenly stat'd in cwd.

## Resolution

`<folder>` is resolved in this order:

1. Absolute path → use as-is.
2. `~`-prefixed → `os.UserHomeDir()` expansion.
3. Otherwise → `filepath.Join(cwd, folder)`.

After resolution, `os.Stat` MUST succeed and the result MUST be a
directory, OR the dispatcher exits 1 with the canonical message:

    Error: cn: folder not found or not a directory: <resolved-path>

The dispatcher does NOT check for a `.git` subdirectory at this stage
— that's the responsibility of the existing `gitutil.RemoteURL` call
inside `runCloneNext`, which already produces a clear error for
non-repo folders.

## Pipeline

Both NEW forms reuse `performCrossDirCloneNext` from
`clonenextcrossdir.go`:

1. Capture `originalDir = os.Getwd()`.
2. `os.Chdir(<resolved-folder>)`.
3. `defer os.Chdir(originalDir)`.
4. `runCloneNext([]string{version, ...flags})` — the existing in-place
   pipeline runs unmodified, including `gitutil.RemoteURL` lookup,
   parent-folder flatten, version-history recording, shell handoff.
5. Print `MsgCNXReturnedFmt` so the user knows control returned.

The version arg defaults to `"v++"` for the single-positional form;
otherwise the user-supplied version is forwarded verbatim.

## Flag forwarding

All flags from the original argv (e.g. `-f`, `--keep`, `--no-desktop`,
`--ssh-key`) are forwarded into the inner `runCloneNext` call via the
existing `extractFlagArgs` helper. The dispatcher itself consumes no
new flags.

## Out of scope

- Batch mode interaction: `--csv` and `--all` continue to take
  precedence over positional args, exactly as today.
- Help text auto-generation: the help file
  `gitmap/helptext/clone-next.md` is updated by hand in the same PR.

## Error contract

| Condition | Exit code | Message |
|---|---|---|
| Folder not found | 1 | `cn: folder not found or not a directory: <path>` |
| Both positionals look like versions | 1 | `cn: ambiguous arguments — both look like version strings` |
| Both positionals look like folders | 1 | `cn: ambiguous arguments — neither looks like a version (use vN, v+N, or v++)` |
| `gitutil.RemoteURL` failure inside target | 1 | (existing) `clone-next: no remote configured: <err>` |

## Test matrix

`gitmap/cmd/clonenextfolderdispatch_test.go` covers:

- single-positional folder → resolved + dispatched with `v++`
- two-positional version-then-folder → resolved + dispatched with version
- two-positional folder-then-version → falls through to existing alias path
- absolute, relative, `~`-prefixed, and bare-name folder resolution
- non-existent folder → exit 1 with canonical message
- folder + flags (`-f`, `--keep`) forwarded into inner call
- ambiguous double-version refusal
- ambiguous double-folder refusal

## References

- Existing cross-dir form: `spec/01-app/107-cn-find-next-bridge.md`
- `looksLikeVersion`: `gitmap/cmd/releaserebase.go:22-28`
- Pipeline reuse: `gitmap/cmd/clonenextcrossdir.go:39-57`
- Version resolver: `gitmap/clonenext/version.go:37-58`
