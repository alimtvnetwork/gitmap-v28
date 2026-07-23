# gitmap MAPRI

Uppercase shorthand for `gitmap make-all-private`.

```
gitmap MAPRI <owner-or-url-or-folder> <patterns> [-Y|--yes] [--verbose]
```

`MAPRI` uses the same resolver, wildcard matcher, confirmation flow,
provider checks, audit tables, and exit codes as `make-all-private`.

## Examples

```
gitmap MAPRI alice "demo-*"
gitmap MAPRI https://github.com/alice "demo-*,!demo-public" -Y
gitmap MAPRI . "archive-*" --verbose
```

## See also

- `gitmap make-all-private` — canonical long form.
- `gitmap MAPUB` — shorthand for `make-all-public`.
- `gitmap visibility-undo` (`vu`) — reverse a completed run.