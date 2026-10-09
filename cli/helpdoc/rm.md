# gitmap rm

Remove one or more repositories — **including the on-disk folder** —
and untrack them in the gitmap database.

## Usage

```
gitmap rm     [-y|--yes] <target>[,<target>...] [<target>...]
gitmap remove [-y|--yes] <target>[,<target>...] [<target>...]
gitmap del    [-y|--yes] <target>[,<target>...] [<target>...]
```

Aliases: `remove`, `del`.

## Target forms

A target may be any of:

| Form | Example | Behavior |
|------|---------|----------|
| Repo slug | `my-repo` | Exact-match on `Slug` |
| Path | `./projects/foo`, `.\macro-ahk`, `/abs/path` | `filepath.Abs` → match on `AbsolutePath` |
| Glob | `macro*`, `gitmap-v?`, `dev-[ab]*` | `filepath.Match` over slug **and** path basename |
| Comma list | `macro*,gitmap*` | Split on commas, each part expanded independently |

Overlapping targets that resolve to the same repo are de-duplicated.

## Confirmation

By default each matched repo is confirmed with a `[y/N]` prompt that
shows the slug and absolute folder. Pass `-y` / `--yes` (anywhere in
the args) to skip every prompt and delete unconditionally.

## What gets deleted

1. The on-disk folder at `AbsolutePath` (`os.RemoveAll`).
2. The `Repo` row in the gitmap database.

If the folder is missing the DB row is still cleared.

## Examples

```
$ gitmap rm my-repo
Delete folder and untrack my-repo
  /home/me/code/my-repo ? [y/N] y
removed: my-repo (/home/me/code/my-repo)

$ gitmap rm macro*
Delete folder and untrack macro-ahk
  /home/me/code/macro-ahk ? [y/N] y
removed: macro-ahk (/home/me/code/macro-ahk)
Delete folder and untrack macro-utils
  /home/me/code/macro-utils ? [y/N] n
skip: macro-utils

$ gitmap rm macro*,gitmap* -y
removed: macro-ahk (/home/me/code/macro-ahk)
removed: macro-utils (/home/me/code/macro-utils)
removed: gitmap-v28 (/home/me/code/gitmap-v28)

$ gitmap del .\macro-ahk
Delete folder and untrack macro-ahk
  C:\src\macro-ahk ? [y/N] y
removed: macro-ahk (C:\src\macro-ahk)
```

## Exit codes

| Code | Meaning |
|------|---------|
| 0    | Every matched repo was removed (or skipped at the prompt) |
| 1    | A target matched nothing, or a deletion failed |

## See also

- `gitmap list` (`gitmap ls`) — show every tracked repo
- `gitmap rescan` — re-discover repos under a scan root
