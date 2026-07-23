# gitmap history-pin

Pin file(s) to their **current** content across **all** Git history.
Every past commit that touched the file will appear to have always
contained the present-day bytes. Wraps `git filter-repo --blob-callback`
in a mirror-clone sandbox so your working repository is never
rewritten in place.

**Alias:** `hpin`

Spec: `spec/04-generic-cli/16-history-rewrite.md`.

## Synopsis

```
gitmap history-pin <path> [<path> ...] [flags]
gitmap hpin        <path> [<path> ...] [flags]
```

Multiple paths may be passed as separate args, joined by `,`, or
joined by `, ` (comma-space). Each path's current bytes are read once
from the working tree, then substituted for every historical blob.

## Flags

Same as `history-purge`. See `gitmap history-purge --help`.

## Example

```
$ gitmap hpin docs/README.md --dry-run

▸ history-rewrite: identifying origin remote
▸ history-rewrite: mirror-cloning ... into /tmp/gitmap-history-rewrite-9f2a
▸ history-rewrite: running filter-repo (pin) for 1 path(s)
▸ history-rewrite: verifying sandbox
✓ history-rewrite: verification passed
✓ history-rewrite: dry-run complete; sandbox at /tmp/gitmap-history-rewrite-9f2a
```

## Verification (automatic)

For each requested path `P`, the command asserts:

```
for sha in $(git log --all --pretty=format:%H -- P); do
  git show "$sha:P" | sha256sum
done | sort -u | wc -l   # must be exactly 1
```

A failed verification exits **6** and never pushes.

## Exit Codes

Same as `history-purge`.

## See Also

- `history-purge` (`hp`) — Remove file(s) from all history
- `fix-repo` (`fr`) — Rewrite versioned-token strings, no history rewrite
## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter history-pin
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).

## Examples

```bash
gitmap history-pin
```
