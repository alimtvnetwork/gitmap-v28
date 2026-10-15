# 05 — `merge write` / `mr` Command Design

## 1. Purpose

Create a pull/merge request from the CLI. Ranked P0 because the user's
Lovable migration needs PR creation without leaving the terminal (Lovable
apps live on GitHub; migration flows open PRs against migrated repos).

## 2. Surface

- `gitmap merge write` — primary verb.
- `gitmap mr` — first-class alias (kept forever per alias rule).

Help example (canonical format, spec 03):

```text
gitmap merge write — help example

WHAT IT DOES
  Opens a pull request from a source branch into a target branch on the
  repo's host (GitHub), optionally as a draft.

INVOCATION
  $ gitmap merge write --title "Migrate Lovable app to Vite shell" \
      --body "Ports the Lovable export into the TanStack Start shell." \
      --source lovable-export --target main --draft

EXPECTED OUTPUT
  pr #42 opened: lovable-export -> main (draft)
  https://github.com/alimtvnetwork/my-app/pull/42
```

## 3. Flags

| Flag | Default | Meaning |
|------|---------|---------|
| `--title` | required | PR title |
| `--body` | `""` | PR body (markdown allowed) |
| `--source` | current branch | Head branch |
| `--target` | repo default branch | Base branch |
| `--draft` | false | Open as draft PR |
| `--yes` | false | Skip the confirmation summary |

## 4. Behavior

1. Resolve repo + host from the current directory (git remote; GitHub only
   for v1 — other hosts return a clear "unsupported host" error, tracked in
   the error module per the no-silent-errors rule).
2. Auth: use the existing credential path (`gh` CLI if present, else the
   GitHub App token flow already used by gitmap). NEVER invent a new token
   store.
3. Print a confirmation summary (source→target, title, draft state); proceed
   on `--yes`, otherwise prompt.
4. On success print the PR number + URL. On failure surface the host's error
   verbatim (no swallowing).

## 5. File box (WS3)

New package `cli/cmdmerge/` (one file, ≤300 lines): dispatch entry,
flags, host call, help text, and the example-registry entry (spec 03).
No changes to pull/commit internals. Dispatch entry added in the merge
family's existing dispatch table (surgical edit only).

## 6. Explicitly out of scope

- Merge conflict resolution, auto-merge, PR review/approve commands.
- Non-GitHub hosts (documented as a future proposal for WS4).
