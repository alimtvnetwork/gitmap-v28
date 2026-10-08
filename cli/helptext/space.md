# gitmap space

> **Generated — do not hand-edit.** This topic is rendered from the
> `cmdspace.SpaceHelpDisplay` HelpDisplay struct. Regenerate with:
>
> `gitmap py 03-ai-scripts/51-helptext-generator.py --topics space`

```text
Usage: gitmap space <subcommand> [flags]

Subcommands
  common  Apply the curated common baselines (.gitignore, .gitattributes, .prettierignore, .prettierrc) and run 'git lfs install --local' — same logic as 'gitmap commons'.
    e.g. gitmap space common --dry-run
  backup-branch  Create a backup branch (backup/<slug>) from the current HEAD for the given task string, then push it to origin.
    e.g. gitmap space backup-branch "CLI help displayer overhaul"
Flags (with common)
  --dry-run, -n  Print planned additions without touching disk
    e.g. gitmap space common --dry-run
  --force, -f  Overwrite conflicting JSON values in .prettierrc
    e.g. gitmap space common --force
Flags (with backup-branch)
  --no-push  Skip pushing the new branch to origin
  --force  Recreate the branch at current HEAD if it already exists
Examples
  gitmap space common  Apply the curated baselines in the current repo
  gitmap space common --dry-run  Preview planned additions without touching disk
  gitmap space backup-branch "CLI help displayer overhaul"  Snapshot HEAD into backup/<slug> and push to origin
  gitmap space backup-branch "hotfix" --no-push  Snapshot HEAD without pushing

* Subcommand help: gitmap space --help <subcommand> (try: common, backup-branch)
```
