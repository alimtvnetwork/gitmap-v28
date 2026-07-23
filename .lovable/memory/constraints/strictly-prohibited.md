---
name: strictly-prohibited
description: Append-only numbered registry of forbidden actions; mirrors spec/03-general/10-strictly-prohibited.md.
type: constraint
---
# Strictly-Prohibited Registry (append-only, numbered)

Entries are FORBIDDEN forever. Never deleted, never renumbered. Refuse and cite the entry by number.

1. Never put time/date/clock/timestamp data in `readme.txt` (any format, any timezone).
2. Never suggest adding "last updated" / "git update time" / time fields to any README.
3. **commit-in:** Never store file content, file hashes, or diff payloads in any commit-in table. Only `RelativePath` strings. **Why:** `git cat-file` already gives byte-exact content on demand. See `spec/03-commit-in/01-overview-and-glossary.md` INV-06.
4. **commit-in:** Never rewrite history of `<source>`. Append-only to current branch tip. **Why:** `<source>` is a live working repo; rewrites would invalidate every collaborator's clone. History rewrites belong to `gitmap history-purge` / `history-pin`. See INV-09.
