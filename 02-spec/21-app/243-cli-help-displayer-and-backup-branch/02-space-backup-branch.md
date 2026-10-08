# Spec 243.2 — `gitmap space backup-branch`

## Command
`gitmap space backup-branch "<task string>" [--no-push] [--force]`

Lives in the existing `space` namespace alongside `space common`.

## Slug generation
From the task string: lowercase, spaces/underscores → hyphens, strip everything
except `[a-z0-9-]`, collapse repeats, trim leading/trailing hyphens, max 60 chars.
Empty result → error asking for a better string.

Example: `"CLI help displayer overhaul"` → `cli-help-displayer-overhaul`.

## Branch
`backup/<slug>`, created from current HEAD.

## Behavior
1. Verify the working tree is clean (`git status --porcelain` empty). If dirty → error:
   "Working tree has uncommitted changes. Commit or stash first — backups snapshot HEAD only."
   (Dirty tree is always refused; there is no override for this.)
2. If `backup/<slug>` already exists and no `--force` → error naming the existing
   branch (suggest a different string, or `--force` to recreate it at current HEAD).
3. Create the branch; print `backup/<slug> @ <short-sha>`.
4. Push with `git push -u origin backup/<slug>` unless `--no-push` is passed.

## Help output
Standard `--help` via the existing help builder until spec 03's displayer lands,
then migrate this command's help to the new displayer like any other command.

## Future (not this spec)
`space backup-list`, `space backup-restore`, `space backup-prune`.
