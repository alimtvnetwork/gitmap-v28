# gitmap doctor fix-repo

Targeted health check for the `fix-repo → gofmt` pipeline. Introduced
in v6.80.1 after Windows users hit
`fork/exec ...gofmt.exe: The filename or extension is too long`
on large repos with hundreds of touched `.go` files.

## Usage

```
gitmap doctor fix-repo [--json] [--budget N]
gitmap doctor fr       [--json] [--budget N]
```

## Probes

| Probe             | Pass means                                            |
|-------------------|-------------------------------------------------------|
| gofmt-present     | `gofmt` resolves on PATH                              |
| gofmt-runs        | `gofmt -l` executes on a scratch file                 |
| argv-budget       | Configured budget ≤ measured Windows argv cap         |
| chunker-selftest  | `chunkPathsForGofmt` invariants hold on synthetic input |

Exit 0 when every probe passes, 1 otherwise. `--budget N` overrides
the argv budget used during the argv-budget measurement and the
chunker self-test.

## Examples

```
gitmap doctor fix-repo
gitmap doctor fix-repo --budget 8000
gitmap doctor fix-repo --json | jq '.results[] | select(.ok==false)'
```
