# gitmap MLPRI

Uppercase shorthand for `gitmap make-last-private`.

```
gitmap MLPRI <owner-or-url> <base> [-Y|--yes]
```

Resolves the highest `-vN` sibling under `<base>` and flips that
single repo private. See `gitmap help make-last-private` for the
full resolution order and exit codes.

## See also

- `gitmap make-last-private` — canonical long form.
- `gitmap MLPUB` — shorthand for `make-last-public`.
- `gitmap MAPRI` — bulk wildcard variant.

## Examples

```bash
gitmap MLPRI
```
