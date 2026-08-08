# gitmap MAPUB

Uppercase shorthand for `gitmap make-all-public`.

```
gitmap MAPUB <owner-or-url-or-folder> <patterns> [-Y|--yes] [--verbose]
```

`MAPUB` uses the same resolver, wildcard matcher, confirmation flow,
provider checks, audit tables, and exit codes as `make-all-public`.

## Examples

```
gitmap MAPUB alice "demo-*"
gitmap MAPUB https://github.com/alice "demo-*,!demo-secret" -Y
gitmap MAPUB . "archive-*" --verbose
```

## See also

- `gitmap make-all-public` — canonical long form.
- `gitmap MAPRI` — shorthand for `make-all-private`.
- `gitmap visibility-undo` (`vu`) — reverse a completed run.