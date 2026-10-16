# Spec 271 — `gitmap space PR-merge`

## Overview

Merge the current branch into a target branch from the CLI, with the same
safety posture as `space backup-branch` (clean tree required, no history
rewrites). Born from the cat-my-v12 PR #10 merge: the user ordered a regular
merge (all commits, no squash) on their behalf, then asked for this ability
as a first-class gitmap command.

## Command

`gitmap space PR-merge [target] [--squash] [--single-commit] [--flat] [--no-push] [--force]`

Lives in the existing `space` namespace alongside `space common` and
`space backup-branch`.

## Target resolution

1. If `[target]` is given, merge the current branch into it.
2. If omitted, resolve the default branch: `git symbolic-ref refs/remotes/origin/HEAD`
   (→ `origin/main`), strip the `origin/` prefix. Fallback order when the
   symbolic ref is missing: `main`, then `master`.
3. Refuse to merge a branch into itself.

## Behavior

1. Verify the working tree is clean (`git status --porcelain` empty, ignoring
   `.gitmap/` diagnostic state like `backup-branch` does). Dirty → error.
2. Fetch `origin` first so the target is current.
3. Default (no flags): `git merge --no-ff <source>` into the target —
   regular merge commit, all commits preserved. This is the user's standing
   preference ("merge it with all the commits. Do not do a squash merge.").
4. `--squash`: `git merge --squash`, then commit with a composed message
   (see "Squash message" below).
5. `--single-commit` / `--flat`: flat single-commit merge variant —
   `--squash` semantics with a one-line message (`<source> -> <target>`)
   and no PR-body embedding.
6. After a successful local merge, push the target with
   `git push origin <target>` unless `--no-push` is passed.
7. On merge conflicts: abort the merge (`git merge --abort`), print the
   conflicted files, and exit non-zero. Never auto-resolve, never force-push.

## Squash message

When `--squash` is used and a GitHub PR exists for the source branch
(`gh pr view <source> --json title,body,number`):

```
<PR title> (#<number>)

<PR body, trimmed>

---
Changed files (heatmap, top 12 by added+deleted lines):
  src/components/vision/white-box/PatternCanvas.tsx  +42 -18
  ...
```

- If no PR exists (or `gh` is unavailable), fall back to
  `Merge <source> into <target> (squash)` plus the heatmap.
- The heatmap comes from `git diff --numstat <merge-base>..<source>`,
  sorted by churn, top 12, rendered as `path  +A -D` — the "short, nice"
  file list the user asked for.

## Flags

- `--squash` — squash-merge with PR-details message.
- `--single-commit` — flat single commit, one-line message.
- `--flat` — alias of `--single-commit`.
- `--no-push` — merge locally only, do not push the target.
- `--force` — reserved; merge never overwrites history so this is a no-op
  guard (errors if passed, to avoid implying destructive behavior).

## Help output

`gitmap space PR-merge --help` prints usage via the same help-text
convention as `space backup-branch`.

## Non-goals

- No rebase, no fast-forward-only mode, no PR creation (use `gh pr create`),
  no conflict auto-resolution.
