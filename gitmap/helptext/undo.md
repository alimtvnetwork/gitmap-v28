# gitmap undo

Restore the latest `gitmap fix-repo` backup snapshot for the current
repo + current version. Every `gitmap fix-repo` write creates a
timestamped snapshot under
`.gitmap/backup/<repo>/v<current>/fix-repo/<UTC-ts>/` (v5.40.0+);
`undo` rolls those bytes back onto your working tree.

## Synopsis

```
gitmap undo                       # restore the latest snapshot
gitmap undo --list                # list snapshots (newest first)
gitmap undo --snapshot <ts>       # restore a specific UTC stamp
gitmap undo --dry-run             # preview without writing
gitmap ud                         # short alias
```

## Backup layout

```
<repoRoot>/.gitmap/backup/
  <repo>/                         # e.g. gitmap-v28
    v<current>/                   # current version at backup time
      fix-repo/
        <UTC-timestamp>/          # one dir per fix-repo invocation
          manifest.json           # list of restored rel paths + meta
          files/<rel/path>        # verbatim pre-rewrite bytes
```

Snapshots are scoped to the current repo + current version so an
`undo` inside `gitmap-v28` never touches a `gitmap-v28` snapshot.

## Example

```
$ gitmap fix-repo --all
fix-repo  base=gitmap  current=v20  mode=--all
...
fix-repo: backed up 7 file(s) → .gitmap/backup/gitmap-v28/v20/fix-repo/20260519T134210Z

$ gitmap undo --list
undo: snapshots under .gitmap/backup/gitmap-v28/v20/fix-repo (2 total)
  * 20260519T134210Z  (7 files)
    20260519T112055Z  (3 files)

$ gitmap undo
undo: restoring snapshot …/20260519T134210Z — 7 file(s) [mode: write]
  [restored] README.md
  [restored] docs/install.md
  ...
undo: restored 7 file(s), 0 failure(s)
```

## Exit codes

`0` ok / `6` bad-flag / `7` write-failed / `8` bad-config (manifest
missing/malformed).

See `spec/04-generic-cli/27-fix-repo-command.md` §"Backup & undo".

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter undo
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).

## Examples

```bash
gitmap undo
```
