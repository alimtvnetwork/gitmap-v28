# gitmap MLPUB

Uppercase shorthand for `gitmap make-last-public`.

```
gitmap MLPUB <owner-or-url> <base> [-Y|--yes]
```

Resolves the highest `-vN` sibling under `<base>` and flips that
single repo public. See `gitmap help make-last-public` for the full
resolution order and exit codes.

## See also

- `gitmap make-last-public` — canonical long form.
- `gitmap MLPRI` — shorthand for `make-last-private`.
- `gitmap MAPUB` — bulk wildcard variant.

## Examples

```bash
gitmap MLPUB
```
