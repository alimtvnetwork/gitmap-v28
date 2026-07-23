# 27 — `fix-repo` Command (Go-native)

## Purpose

`gitmap-v27 fix-repo` (alias `gitmap-v27 fr`) rewrites prior versioned-repo-name
tokens to the current version across all tracked text files. It is the
Go-native re-implementation of the existing `fix-repo.ps1` /
`fix-repo.sh` shell scripts, with byte-for-byte identical default
behavior, exit codes, and config schema. The PowerShell + Bash scripts
remain as bootstrap helpers; the Go command is the canonical entry
point and is shipped inside the `gitmap-v27` binary.

## Synopsis

```
gitmap-v27 fix-repo [-2 | -3 | -5 | --all] [--dry-run] [--verbose] [--config <path>]
gitmap-v27 fr                                                       # short alias
```

PowerShell-style single-dash flags (`-2`, `-3`, `-5`, `-All`,
`-DryRun`, `-Verbose`, `-Config <p>`) are also accepted as aliases so
existing muscle memory continues to work.

## Behavior

1. Resolve repo identity from `git`:
   - `git rev-parse --show-toplevel` → repo root (else `E_NOT_A_REPO`).
   - `git config --get remote.origin.url` → remote URL (else `E_NO_REMOTE`).
   - Parse remote URL into `{host, owner, repo}` for HTTPS, SSH
     (`git@host:owner/repo`), and `ssh://` forms.
   - Split `repo` into `{base, version}` using regex `^(.+)-v(\d+)$`
     (else `E_NO_VERSION_SUFFIX`). Version must be ≥ 1
     (else `E_BAD_VERSION`).
2. Compute target span:
   - Default (no flag): last 2 prior versions.
   - `-2 | -3 | -5`: last N prior versions.
   - `--all`: every prior version `1..current-1`.
   - Targets clamp to `[max(1, current-span) .. current-1]`.
3. Load config (default `<repoRoot>/fix-repo.config.json`,
   override with `--config <path>`):
   - `ignoreDirs`: array of repo-relative directory prefixes to skip.
   - `ignorePatterns`: array of glob patterns. `**` matches across
     segments; `*` matches within a single segment; `?` matches one
     non-`/` char. Missing config file is non-fatal; explicit
     `--config <missing>` is fatal (`E_BAD_CONFIG`).
4. Enumerate tracked files via `git ls-files`. For each path:
   - Skip if matched by `ignoreDirs` / `ignorePatterns`.
   - Skip reparse points and files larger than 5 MiB.
   - Skip files whose extension is in the binary allow-list:
     `.png .jpg .jpeg .gif .webp .ico .pdf .zip .tar .gz .tgz .bz2
     .xz .7z .rar .woff .woff2 .ttf .otf .eot .mp3 .mp4 .mov .wav
     .ogg .webm .class .jar .so .dylib .dll .exe .pyc`.
   - Skip files whose first 8 KiB contain a NUL byte.
5. For each surviving file, replace every literal occurrence of
   `{base}-v{N}` (for each `N` in targets) with `{base}-v{current}`,
   guarded by a negative-lookahead so `{base}-v1` does **not** match
   inside `{base}-v10`. Counts every replacement.
6. In `--dry-run`, files are not written; counts are still reported.
7. Print header + summary identical to the PowerShell script:

   ```
   fix-repo  base=<base>  current=v<N>  mode=<mode>
   targets:  v1, v2, ...
   host:     <host>  owner=<owner>

   scanned: <S> files
   changed: <C> files (<R> replacements)
   mode:    write|dry-run
   ```

## Exit codes

| Code | Constant              | Meaning                                  |
|------|-----------------------|------------------------------------------|
| 0    | `ExitOk`              | Success.                                 |
| 2    | `ExitNotARepo`        | Not inside a git work tree.              |
| 3    | `ExitNoRemote`        | Missing or unparseable `origin` URL.     |
| 4    | `ExitNoVersionSuffix` | Repo name has no `-vN` suffix.           |
| 5    | `ExitBadVersion`      | Version ≤ 0.                             |
| 6    | `ExitBadFlag`         | Unknown / conflicting CLI flag.          |
| 7    | `ExitWriteFailed`     | At least one file write failed.          |
| 8    | `ExitBadConfig`       | Explicit `--config` missing or invalid.  |

These match `fix-repo.ps1` and `fix-repo.sh` 1:1 so CI scripts that
inspect the exit code keep working when they switch invocation from
the script to the binary.

## Idempotency & safety

- Replacing `{base}-v{N}` with the same `{base}-v{current}` is
  idempotent: a clean repo with no prior versions yields
  `changed: 0`.
- The negative-lookahead (`(?!\d)`) prevents partial-token rewrites
  inside larger version numbers (`v10`, `v123`).
- Tracked-file enumeration via `git ls-files` automatically respects
  `.gitignore`, so untracked build output is never touched.
- The 5 MiB / NUL-byte / extension guards prevent corrupting
  binary assets that happen to be tracked.

## Naming + ownership

- Top-level command name: `fix-repo` (kebab-case, like `desktop-sync`).
- Short alias: `fr`.
- Constants live in `gitmap-v27/constants/constants_fixrepo.go` (package-
  domain ownership, per the constants-ownership rule).
- Implementation is split across `gitmap-v27/cmd/fixrepo*.go` files,
  each ≤ 200 lines, functions ≤ 15 lines, positive conditionals only,
  no swallowed errors (logged to `os.Stderr` via the standard format).

## Bare-base scope rule (v5.38.0+)

The rewrite engine has an extra pass that substitutes standalone
`{base}` tokens (no `-vN` suffix) for the case where the original
repository shipped without a `-v1` suffix and downstream references
read the bare name.

**This pass is restricted to the v1→v2 transition only.** Concretely,
the bare-base sweep runs if and only if both of these are true:

1. `1` is in the target version set (i.e. v1 is being rewritten), AND
2. The current repo version is exactly `2` (`current == 2`).

For any current version ≥ 3 the bare-base pass is skipped even when v1
is in the target span. Rationale: once the project has shipped past
v2, a bare `{base}` token in source / docs / scripts is overwhelmingly
NOT the pre-versioned origin URL — it is the binary name, package
identifier, brand string, or an unrelated repo reference, and rewriting
it to `{base}-v{current}` silently corrupts the repo.

### Before / after table — `gitmap-v27 fix-repo` semantics by current version

| Working repo | Target span | Token in source       | Rewritten to        | Why |
|--------------|-------------|-----------------------|---------------------|-----|
| `gitmap-v27`  | v1          | `gitmap-v27`              | `gitmap-v27`         | Bare sweep ACTIVE (n=1, current=2) |
| `gitmap-v27`  | v1          | `gitmap-v27`           | `gitmap-v27`         | Versioned token rewrite |
| `gitmap-v27`  | v1          | `gitmap.js`           | `gitmap.js`         | Word-boundary guard |
| `gitmap-v27`  | v1 (`-r nv`)| `gitmap-v27`              | `gitmap-v27`            | Bare sweep suppressed by `--restrict no-version` |
| `gitmap-v27`  | v1, v2      | `gitmap-v27`              | `gitmap-v27`            | Bare sweep SKIPPED (current ≥ 3) |
| `gitmap-v27`  | v1, v2      | `gitmap-v27`/`-v2`     | `gitmap-v27`         | Versioned token rewrite |
| `gitmap-v27`  | v1..v3      | `https://…/gitmap-v27`    | `https://…/gitmap-v27`  | Bare upstream URL preserved |
| `gitmap-v27`  | v1..v3      | `gitmap-v27`          | `gitmap-v27`        | Negative-lookahead vs `-v1` |
| `gitmap-v27`  | v1..v3      | `gitmap-v27..v3`       | `gitmap-v27`         | Versioned token rewrite |

Pre-v5.38.0 regression: running inside `gitmap-v27` rewrote every bare
`gitmap-v27` mention (binary name, upstream URL, `gitmap-cli` package
descriptor, brand strings) to `gitmap-v27`. The v5.38.0 gate restored
correctness; the table above is the contract going forward.

Implementation: `applyAllTargets` in `gitmap-v27/cmd/fixrepo_rewrite.go`,
guarded by `if n == 1 && current == 2`. Regression locks:
`TestApplyAllTargets_BareBase_SkippedAtV3Plus` and
`TestApplyAllTargets_BareBase_SkippedAtV4WithV1InTargets` in
`fixrepo_rewrite_barebase_test.go`.

## Restrict modes (v5.39.0+)

`--restrict <mode>` (alias `-r <mode>`) narrows the rewrite scope.
Currently the only defined value is `no-version` (short: `nv`).

When `--restrict no-version` is set, the v1→v2 bare-base sweep is
suppressed. The engine rewrites ONLY versioned `{base}-vN` tokens —
bare `{base}` occurrences are left untouched even during a v1→v2
bump. Useful when the project already used the `{base}-v1` form on
its first remote and the bare token is a legitimate unrelated
identifier (binary name, package, brand) that must be preserved.

Examples:

```
gitmap-v27 fix-repo -2 --restrict no-version
gitmap-v27 fr -2 -r nv
```

Implementation: `applyAllTargetsR` in
`gitmap-v27/cmd/fixrepo_rewrite.go` gates the bare-base call on
`!restrictNoVersion`. Flag parsing in `fixrepo_flags.go`
(`consumeFixRepoRestrictArg`) accepts `--restrict`, `-restrict`,
`-r`, plus the `=value` forms; unknown values exit with
`E_BAD_FLAG` (code 6).

## Cross-references

- PowerShell script: `fix-repo.ps1` + `scripts/fix-repo/*.ps1`.
- POSIX script: `fix-repo.sh` + `scripts/fix-repo/*.sh`.
- Constants ownership rule: `spec/12-consolidated-guidelines/02-go-code-style.md`.
- Strictly-prohibited registry: `spec/03-general/10-strictly-prohibited.md`
  (no time/date in `readme.txt`, no manual edits to `.gitmap/release/`).

## Backup & undo (v5.40.0+)

Every non-dry-run `gitmap-v27 fix-repo` invocation that modifies ≥ 1 file
writes a backup snapshot under the repo root:

```
<repoRoot>/.gitmap/backup/<repo-name>/v<current>/fix-repo/<UTC-ts>/
  manifest.json          schemaVersion, repo, currentVersion, timestamp,
                         gitmapVersion, files[]
  files/<rel/path>       verbatim pre-rewrite bytes
```

Rules:

1. **One snapshot per invocation.** Timestamp format
   `20060102T150405Z` so lexical sort == chronological sort.
2. **Scoped to `<repo>/v<current>`.** An `undo` inside `gitmap-v27`
   never sees a `gitmap-v27` snapshot. This prevents cross-version
   restores that would re-introduce stale `{base}-vN` tokens.
3. **Lazy + idempotent.** The snapshot directory is created on the
   FIRST `BackupFile` call so dry-run / no-op sweeps leave no trace.
   Within one snapshot the first observation of a file wins
   (subsequent rewrites in the same invocation never overwrite the
   original backup).
4. **Manifest-driven restore.** `gitmap-v27 undo` reads `manifest.json`
   and copies each listed `files/<rel>` back to `<repoRoot>/<rel>`.

### `gitmap-v27 undo` (alias `ud`)

```
gitmap-v27 undo                       # restore latest snapshot
gitmap-v27 undo --list                # list snapshots, newest first (* = latest)
gitmap-v27 undo --snapshot <UTC-ts>   # restore a specific stamp
gitmap-v27 undo --dry-run             # preview without writing
```

Exit codes: `0` ok / `6` bad-flag / `7` write-failed / `8` bad-config
(manifest missing or malformed).

Implementation: `gitmap-v27/cmd/fixrepo_backup.go` (session +
`BackupFile` / `Finalize`), `gitmap-v27/cmd/undo.go` (list / pick /
restore). Wired into the sweep by splitting `rewriteOneFile` into a
pure-compute step plus `persistRewrittenFile`, which calls
`backup.BackupFile(rel)` strictly BEFORE `os.WriteFile` so the
snapshot always captures the pre-rewrite bytes.
