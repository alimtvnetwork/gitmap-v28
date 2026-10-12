# gitmap clean-repos

Purge all excluded, cache, and corrupt repositories from the gitmap database.

## Usage

```
gitmap clean-repos
gitmap exclude clean
```

Aliases: `cleanrepos`, `repo-clean`, `repoclean`.

## Behavior

Walks every tracked repository and deletes the database row for any repo that
is on the exclusion list, missing from disk, or corrupt. Only database records
are removed — on-disk folders are never touched.

Use it after moving or deleting repo folders outside gitmap, or to drop
stale entries left behind by `gitmap rm --db-only`.

## Examples

### Purge stale database records

```bash
gitmap clean-repos
```

### Purge via the exclude subcommand

```bash
gitmap exclude clean
```
