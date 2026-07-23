# gitmap commit-right

> **Status (v3.76.0):** **LIVE.** Phase 1 of spec/01-app/106 has shipped.
> `commit-left` and `commit-both` are still scaffolds.

Replay LEFT's commits onto RIGHT as a fresh, cleaned commit sequence.
The "-right" suffix names the **destination**, exactly like
`merge-right` writes files to RIGHT.

## Alias

cmr

> Spec §13 reserved `cr`, but `cr` is already taken by `cpp-repos`.
> Use `cmr` instead. The long-form `commit-right` always works.

## Usage

    gitmap commit-right LEFT RIGHT [flags]

LEFT and RIGHT use the same endpoint syntax as the merge-* family:
either a local folder path or an `https://` / `git@` URL with optional
`:branch` suffix.

## Flags (planned, full set)

| Flag | Default | Description |
|------|---------|-------------|
| --mirror | false | Delete target files not present in source commit |
| --include-merges | false | Include merge commits in the replay set |
| --limit N | 0 | Replay at most N source commits (oldest first) |
| --since <sha\|date> | (auto) | Override the divergence base |
| --strip <regex> | (config) | Add a strip pattern (repeatable) |
| --no-strip | false | Disable all strip patterns |
| --drop <regex> | (config) | Add a drop pattern (repeatable) |
| --no-drop | false | Replay every commit |
| --conventional | (config) | Force conventional-commit normalization |
| --no-conventional | false | Disable conventional-commit normalization |
| --provenance | true | Append `gitmap-replay:` footer |
| --no-provenance | false | Skip provenance footer |
| --prefer-source | false | Source side wins file conflicts |
| --prefer-target | false | Target side wins file conflicts |
| --force-replay | false | Replay even commits already carrying a footer |
| --dry-run | false | Print plan + cleaned messages; no writes |
| --yes / -y | false | Skip the confirmation prompt |
| --no-push | false | Stop after local commit (skip push) |
| --no-commit | false | Copy files but skip both commit and push |

## Prerequisites

- Both endpoints resolvable (local folder OR clonable URL).
- Target side must be on a writable branch.
- Spec §18 Phase 1 implementation merged.

## Examples (planned UX, from spec §3)

    gitmap commit-right ./repo-A ./repo-B

Output:

    [commit-right] replaying 7 commits from LEFT onto RIGHT:
      [1/7] a3f2c1d  feat: add OAuth flow
      [2/7] b7e4a9f  fix: handle expired tokens
      ...
    [commit-right] proceed? [y/N]

## See Also

- [commit-left](commit-left.md) — opposite direction
- [commit-both](commit-both.md) — bidirectional
- [merge-right](merge-right.md) — file-state mirror (no commit replay)
- spec/01-app/106-commit-left-right-both.md — full design

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter commit-right
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
