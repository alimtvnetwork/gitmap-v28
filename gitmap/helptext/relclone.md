# gitmap relclone (alias of `reclone`)

`relclone` is a **legacy alias** of `gitmap reclone` — the canonical
command that re-runs `git clone` against `gitmap scan` artifacts.

Behavior, flags, exit codes, and arguments are identical. The
`relclone` and `rc` spellings will keep working forever; new docs,
examples, and tab-completion use `reclone` / `rec`.

```
gitmap reclone  <file> [flags]  # canonical
gitmap relclone <file> [flags]  # this alias
gitmap rc       <file> [flags]  # short alias
```

For the full reference, see:

```
gitmap help reclone
```

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter relclone
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).

## Examples

```bash
gitmap relclone
```
