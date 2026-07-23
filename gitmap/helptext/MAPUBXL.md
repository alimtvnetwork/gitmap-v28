# gitmap MAPUBXL

Uppercase shorthand for `gitmap make-all-public-except-latest`.

```
gitmap MAPUBXL <owner-or-url> <patterns> \
    [-Y|--yes] [--verbose] [--parallel=N] [--cache-ttl=SECONDS]
```

Behaves identically to `make-all-public-except-latest`: every matched
repo whose name ends in `-v<digits>` is grouped by its base prefix,
and the highest version per group is preserved. Repos without a `-vN`
suffix flow through untouched.

## Examples

```
gitmap MAPUBXL alice "demo-v*"
gitmap MAPUBXL alice "myapp-v*,proto-v*" -Y --parallel=16
```

## See also

- `gitmap make-all-public-except-latest` — canonical long form.
- `gitmap MAPRIXL` — shorthand for `make-all-private-except-latest`.
- `gitmap MAPUB` — flip every match, no preservation.
