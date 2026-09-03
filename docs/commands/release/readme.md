# Release Management & Versioning (`gitmap release`)

Automate release ceremonies, Git tagging, changelog generation, and asset distribution.

## Commands

### `gitmap release [ver]`
* **Alias:** `r`
* Validates working tree, bumps version, generates changelog, builds cross-compiled binaries, creates Git tag, pushes release branch, and publishes to GitHub.
* Flags: `--bump major|minor|patch`, `--draft`, `--dry-run`, `--assets <path>`, `--bin` (`-b`), `--targets <list>`, `--compress`, `--checksums`, `--yes` (`-y`).

### `gitmap pull-release [ver]`
* **Alias:** `pr`
* Pulls latest commits (`--ff-only|--rebase|--merge`) and immediately triggers release pipeline.

### `gitmap release-self`
* **Alias:** `rs`
* Releases GitMap itself from any directory.

### `gitmap release-branch`
* **Alias:** `rb`
* Completes release ceremony from an existing release branch.

### `gitmap changelog [ver]`
* **Alias:** `cl`
* Views concise release notes (supports `--open`, `--source`).

### `gitmap changelog-gen`
* **Alias:** `cg`
* Auto-generates changelog markdown between two tags.

### `gitmap list-versions`
* **Alias:** `lv`
* Displays all release tags sorted by SemVer (highest first).

### `gitmap list-releases`
* **Alias:** `lr`
* Shows releases recorded in `.gitmap/release/` metadata files or database.

### `gitmap prune`
* **Alias:** `prn`
* Deletes stale local release branches that have already been tagged.
