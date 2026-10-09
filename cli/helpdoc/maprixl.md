# gitmap MAPRIXL

Uppercase shorthand for `gitmap make-all-private-except-latest`.

```
gitmap MAPRIXL <owner-or-url> <patterns> \
    [-Y|--yes] [--verbose] [--parallel=N] [--cache-ttl=SECONDS]
```

Behaves identically to `make-all-private-except-latest`: every matched
repo whose name ends in `-v<digits>` is grouped by its base prefix,
and the highest version per group is preserved. Repos without a `-vN`
suffix flow through untouched.

## Examples

```
gitmap MAPRIXL alice "demo-v*"
gitmap MAPRIXL alice "myapp-v*,proto-v*" -Y --parallel=16
```

## See also

- `gitmap make-all-private-except-latest` — canonical long form.
- `gitmap MAPUBXL` — shorthand for `make-all-public-except-latest`.
- `gitmap MAPRI` — flip every match, no preservation.
