# `gitmap pull-release-cd` (alias `prc`)

> Multi-repo, one-shot `pull-release` runner.

Introduced in **v5.31.0**.

Lets you release multiple repositories from a single machine with one command —
GitMap chdirs into each target repo, runs the standard `pull-release` pipeline
with `-y` (auto-commit), then moves to the next entry.

## Synopsis

```
gitmap pull-release-cd <name-or-url> <version>[, <name-or-url> <version> ...]
gitmap prc            <name-or-url> <version>[, ...]
gitm   prc            <name-or-url> <version>[, ...]
```

## Argument syntax

A **comma-separated list of `<name-or-url> <version>` pairs**. The comma may
appear as its own shell token (`gitmap`, `marco`) or attached to the previous
token (`gitmap,marco`) — the parser normalises both forms.

```
gitm prc gitmap v5.30, marco v2.5, other-rep v3.5, url-git v3.1
```

| Token | Resolution |
|-------|------------|
| Plain slug (`gitmap`, `marco`) | Looked up in the gitmap SQLite DB via `FindBySlug`. Must already be registered (run `gitmap scan` first). |
| HTTPS/SSH URL (`https://…`, `git@…:…`) | Cloned via `gitmap clone <url>` into the active workdir, then registered in the DB automatically by the cloner. The derived slug is used for the subsequent `pull-release`. |

## Per-entry behavior

For every parsed entry, in order:

1. Resolve the target absolute path (slug lookup, or clone-then-lookup for URLs).
2. Spawn `gitmap pull-release <version> -y` as a subprocess with `cwd` set to that path. Subprocess isolation means one failing repo never corrupts cwd or aborts the remaining entries (status is collected and surfaced at the end).
3. Stream stdout/stderr live, prefixed with `[<slug>] ` so multi-repo output stays readable.
4. Record success / failure and continue.

`-y` is **implicit and non-negotiable** in this command — the whole point is
unattended multi-repo release. The auto-commit prompt described in
`.lovable/memory/features/post-release-commit.md` is answered "yes"
automatically for the `.gitmap/release/latest.json` (and any other modified
files) of every repo in the batch.

## Exit codes

| Code | Meaning |
|------|---------|
| `0`  | Every entry released successfully. |
| `1`  | One or more entries failed. A summary table is printed to stderr listing each `[slug] FAIL <reason>`. |
| `2`  | Argument parse error (missing version, empty list, malformed pair). |

## Examples

```bash
# Three already-registered repos in one shot
gitm prc gitmap v5.31.0, marco v2.5.0, other-rep v3.5.0

# Mix of registered slugs and a fresh URL clone-then-release
gitm prc gitmap v5.31.0, https://github.com/me/url-git v3.1.0
```

## Non-goals (deferred)

- Parallel execution. v1 is strictly sequential — concurrent releases against
  a shared gitmap DB would race on the `releaseAcrossRepos` writer.
- Per-entry pull mode override (`--rebase` / `--merge`). v1 inherits the
  `pull-release` default (`--ff-only`).
- Dry-run preview (`--dry-run`). Add in a follow-up if needed.

## See also

- `gitmap pull-release` / `pr` — single-repo cwd version.
- `gitmap release-alias-pull` / `rap` — alias-based pull-then-release from any dir.
