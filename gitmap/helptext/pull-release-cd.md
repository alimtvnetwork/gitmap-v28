# gitmap pull-release-cd

> **Multi-repo, one-shot pull-release runner.**
> Chdirs into each target repository, runs the standard `pull-release`
> pipeline with `-y` (auto-commit), then moves on. New in **v5.31.0**.

## Aliases

| Form | Notes |
|------|-------|
| `prc` | **Primary short alias** |
| `pull-release-cd` | Canonical long form |

## Synopsis

```
gitmap pull-release-cd <name-or-url> <version>[, <name-or-url> <version> ...]
gitmap prc            <name-or-url> <version>[, ...]
gitm   prc            <name-or-url> <version>[, ...]
```

## Argument syntax

A **comma-separated list of `<name-or-url> <version>` pairs**. The comma may
appear as its own shell token or attached to the previous one — both forms
parse identically.

| Token | Resolution |
|-------|------------|
| Slug (`gitmap`, `marco`) | Looked up in the gitmap DB via `FindBySlug`. Must already be registered (run `gitmap scan` first). |
| URL (`https://…`, `git@…:…`) | Cloned via `gitmap clone <url>`, then the derived slug is resolved. |

## Behavior

For every parsed entry, sequentially:

1. Resolve the target absolute path (slug lookup, or clone-then-lookup for URLs).
2. Spawn `gitmap pull-release <version> -y` as a subprocess with `cwd` set to that path.
3. Stream stdout/stderr live.
4. Collect success/failure; never abort the batch on a single failure.
5. Print a summary table at the end.

`-y` is **implicit and non-negotiable** — the point of `prc` is unattended
multi-repo release, so the post-release auto-commit prompt is always answered
"yes" for that repo's modified files (notably `.gitmap/release/latest.json`).

## Exit codes

| Code | Meaning |
|------|---------|
| `0`  | Every entry released successfully. |
| `1`  | One or more entries failed (see summary). |
| `2`  | Argument parse error. |

## Examples

```
# Three already-registered repos
gitm prc gitmap v5.31.0, marco v2.5.0, other-rep v3.5.0

# Mix of registered slug and a fresh URL to clone-then-release
gitm prc gitmap v5.31.0, https://github.com/me/url-git v3.1.0
```

## See also

- **`gitmap pull-release`** / **`pr`** — single-repo cwd version.
- **`gitmap release-alias-pull`** / **`rap`** — alias-based pull-then-release.

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter pull-release-cd
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
