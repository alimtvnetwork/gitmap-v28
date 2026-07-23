# gitmap history-purge

Remove file(s) and folder(s) from **all** Git history. Wraps
`git filter-repo --invert-paths` in a mirror-clone sandbox so your
working repository is never rewritten in place.

**Alias:** `hp`

Spec: `spec/04-generic-cli/16-history-rewrite.md`.

## Synopsis

```
gitmap history-purge <path> [<path> ...] [flags]
gitmap hp            <path> [<path> ...] [flags]
```

Multiple paths may be passed as separate args, joined by `,`, or
joined by `, ` (comma-space). Files and folders are both accepted.

## Flags

| Flag             | Description                                                  |
|------------------|--------------------------------------------------------------|
| `-y`, `--yes`    | Skip the push confirmation prompt; force-push on success     |
| `--no-push`      | Stop after verification; print the manual push command       |
| `--dry-run`      | Run the rewrite + verification, then exit without pushing    |
| `--message <s>`  | Rewrite the message of every touched commit to `<s>`         |
| `--keep-sandbox` | Don't delete the temp mirror-clone on exit                   |
| `-q`, `--quiet`  | Suppress per-phase progress lines                            |

## Example

```
$ gitmap hp secret.env build/cache.bin --message "history cleanup"

▸ history-rewrite: identifying origin remote
▸ history-rewrite: mirror-cloning git@github.com:acme/repo.git into /tmp/gitmap-history-rewrite-1a2b
▸ history-rewrite: running filter-repo (purge) for 2 path(s)
▸ history-rewrite: verifying sandbox
✓ history-rewrite: verification passed

────────────────────────────────────────────────────────────
  ✓ Verification PASSED
  Mode      : history-purge
  Paths     : 2
  Sandbox   : /tmp/gitmap-history-rewrite-1a2b
  Remote    : git@github.com:acme/repo.git
  Action    : git push --mirror --force-with-lease
  WARNING   : This rewrites published history. Downstream
              clones will need to re-clone or hard-reset.
────────────────────────────────────────────────────────────
Type 'yes' to force-push to git@github.com:acme/repo.git (anything else aborts): yes
▸ history-rewrite: pushing to git@github.com:acme/repo.git with --force-with-lease
✓ history-rewrite: push complete
```

## Exit Codes

`0` ok / `2` not-in-repo / `3` filter-repo-not-installed /
`4` bad-args / `5` filter-repo-failed / `6` verify-failed /
`7` push-failed.

## See Also

- `history-pin` (`hpin`) — Pin file content across all history
- `fix-repo` (`fr`) — Rewrite versioned-token strings, no history rewrite
## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter history-purge
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).

## Examples

```bash
gitmap history-purge
```
