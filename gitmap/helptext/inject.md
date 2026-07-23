# gitmap inject (alias: inj)

Inject an existing on-disk folder into your tooling: register it
with GitHub Desktop, open it in VS Code, and (when a remote `origin`
is configured) record it in the gitmap SQLite database so it appears
in `cd`, `list`, etc.

## Usage

    gitmap inject              # inject the current working directory
    gitmap inject <folder>     # inject the given folder
    gitmap inj   <folder>      # short alias
    gitmap inject --force      # re-register Desktop + VS Code (alias: -f)

## Idempotency (--force / -f)

`inject` (and its sibling `gitmap open`) stamp two timestamps on the
`Repo` row each time they register the folder with an external tool:

| Column | Set when |
|--------|----------|
| `LastInjectedDesktopAt` | GitHub Desktop registration ran |
| `LastInjectedVSCodeAt`  | VS Code window opened |

On subsequent invocations, if a stamp is already non-empty the
matching action is skipped with a one-line notice naming the tool +
its previous timestamp. Pass `--force` (`-f`) to bypass both gates
and re-stamp to "now" after the side effects run. The check is
per-tool: it's valid to skip Desktop while still re-opening VS Code
(or vice versa).


## Supported `<folder>` formats

`<folder>` is optional. When omitted, the current working directory
is used. When provided, all of the following are accepted:

| Format    | Example                          | Resolved against         |
|-----------|----------------------------------|--------------------------|
| absolute  | `C:\dev\macro-ahk-v11`           | filesystem root          |
| absolute  | `/home/me/dev/some-repo`         | filesystem root          |
| relative  | `../sibling-repo`                | current working dir      |
| relative  | `./projects/api`                 | current working dir      |
| home `~`  | `~/dev/macro-ahk-v11`            | `$HOME` (or USERPROFILE) |
| bare name | `macro-ahk-v11` (next to cwd)    | current working dir      |

The path must resolve to an existing directory. A missing or
non-directory target aborts with a clear error before any side
effects run.

## What happens, step by step

1. **Database upsert (conditional).** gitmap runs
   `git remote get-url origin` inside the target folder.
   - If a remote URL is returned, the repo is upserted into SQLite.
     SSH URLs (`git@…` or `ssh://…`) are stored in the `SSHUrl`
     column; everything else is stored in `HTTPSUrl`. After this,
     the repo is reachable via `gitmap cd <name>`, `gitmap list`,
     etc.
   - If `origin` is missing (local-only repo, brand-new sandbox,
     or non-repo folder), gitmap prints
     `no remote configured, skipping database` and continues.
     **No error, no exit** — the database step is best-effort.
2. **GitHub Desktop registration.** The folder is registered with
   Desktop. Non-repo folders are silently ignored by Desktop.
3. **VS Code open.** Opens the folder in a new VS Code window.
   When `code` is not on PATH, gitmap prints a warning and moves
   on — the command still succeeds.
4. **Shell handoff.** The parent shell `cd`s into the injected
   folder (same UX as `clone`, `cn`, and `cd`). Skipped silently
   when the shell wrapper isn't installed.

## Examples

### cwd default (no argument)

    cd ~/dev/some-repo
    gitmap inject
      → registers cwd, opens in VS Code, cds you back in.

    cd /c/work/api-service
    gitmap inj
      → short alias, same behavior on the cwd.

### Positional with `~` (home-relative)

    gitmap inject ~/dev/macro-ahk-v11
      → expands ~ to $HOME, registers, opens.

    gitmap inj ~/projects/api
      → short alias + ~ path, resolved against $HOME (or USERPROFILE).

    gitmap inject ~
      → injects your home directory itself (DB step skipped, no origin).

### Positional with relative paths

    cd ~/dev
    gitmap inject ./macro-ahk-v11
      → explicit ./ relative to cwd.

    gitmap inj ../sibling-repo
      → parent-relative path resolved against cwd.

    gitmap inject some-repo
      → bare name resolved as ./some-repo next to cwd.

    gitmap inj ../../monorepo/packages/ui
      → deep relative path, fully resolved before any side effects.

### Positional with absolute paths

    gitmap inject /home/me/dev/some-repo
      → absolute POSIX path.

    gitmap inject C:\dev\macro-ahk-v11
      → absolute Windows path.

    gitmap inject C:\sandbox\plain-folder
      → no .git/, no origin → DB step is skipped, but Desktop +
        VS Code still proceed and you're cd'd into the folder.

## Notes

- No `.git/` is required. VS Code happily opens any folder, and
  Desktop silently skips non-repos. The DB upsert is the only
  step that requires `origin`, and it fails open.
- For a fresh clone instead of injecting an existing folder, use
  `gitmap clone <url>`.

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter inject
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
