# Next 18 Task — `cfr`/`cfrp` honor existing folder origin transport (plan 03 step 3, partial)

Executed the `cfr`/`cfrp` half of step 3 in `.lovable/plans/pending/03-reclone-transport-and-vscode-open.md`:

- New `preferExistingFolderTransport` reads the destination folder's `remote.origin.url` (when `.git/` exists) and rewrites the positional URL to match the existing transport before `executeDirectClone` runs. Closes the silent HTTPS-downgrade gap identified in `.lovable/audits/2026-06-07-reclone-pickers.md`.
- New tests cover the classifier + both rewrite arms + the fresh-clone untouched case.
- Bumped to v6.26.0 across `gitmap/constants/constants.go`, `src/constants/index.ts`, `README.md`, `CHANGELOG.md`.

Pending in step 3: persist transport on the `Repo` row (depends on step 2 migration 007) and write a `RecloneHistory` event row.
