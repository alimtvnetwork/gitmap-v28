# Release Management & Versioning (`gitmap release`)

Automate release ceremonies, Git tagging, changelog generation, and asset distribution.

<div align="center">

<img src="../../assets/release.svg" alt="GitMap Release Terminal Demo" width="850">

</div>

## Commands & Flags

### 1. `gitmap release [ver]`
* **Alias:** `r`
* **Flags:**
  * `--bump major|minor|patch`: Auto-increments version number.
  * `--draft`: Creates an unpublished draft release on GitHub.
  * `--dry-run`: Previews release actions without executing git tag or push.
  * `--bin` (`-b`): Cross-compiles Go binaries and includes in release assets.
  * `--targets <list>`: Cross-compile targets (e.g. `windows/amd64,linux/arm64`).
  * `--compress`: Compresses binaries into `.zip` (Windows) or `.tar.gz`.
  * `--checksums`: Generates `SHA256` checksums file.
  * `--yes` (`-y`): Auto-confirms prompts.

#### Flag Examples:
```bash
# Minor bump with cross-compiled binaries and compression
gitmap release --bump minor --bin --compress --yes

# Create draft patch release
gitmap release --bump patch --draft

# Preview release steps safely
gitmap release --dry-run
```

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
