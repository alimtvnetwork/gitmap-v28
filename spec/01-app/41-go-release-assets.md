# Go Release Assets — Cross-Compiled Binaries

## Overview

When `gitmap-v28 release` runs in a repository that contains a Go project
(detected by `go.mod`), the release process should automatically
cross-compile binaries for common OS/arch pairs and attach them as
GitHub release assets — without requiring any external download tools
(no `gh`, `curl`, or `wget`).

---

## Goals

1. **Zero external dependencies** — everything built using Go's native
   cross-compilation (`GOOS`/`GOARCH` env vars) and the GitHub Releases
   API (already used for tag creation).
2. **Automatic detection** — if the repo has `go.mod` and a buildable
   `main` package, assets are produced. No extra flags required.
3. **Opt-out via flag** — `--no-assets` skips binary compilation.
4. **Configurable targets** — a `release.targets` key in `config.json`
   can override the default OS/arch matrix.

---

## Default Target Matrix

| GOOS    | GOARCH | Filename suffix          |
|---------|--------|--------------------------|
| windows | amd64  | `_windows_amd64.exe`     |
| windows | arm64  | `_windows_arm64.exe`     |
| linux   | amd64  | `_linux_amd64`           |
| linux   | arm64  | `_linux_arm64`           |
| darwin  | amd64  | `_darwin_amd64`          |
| darwin  | arm64  | `_darwin_arm64`          |

Binary naming pattern: `<module-name>_<version>_<goos>_<goarch>[.exe]`

Example: `gitmap_v2.14.0_windows_amd64.exe`

---

## Build Step Integration

The release workflow (10-step lifecycle) gains a new sub-step between
step 6 (Tagging) and step 8 (Push):

### Step 7a — Cross-Compile (Go only)

1. Read `go.mod` to extract module name.
2. Locate the main package entry point:
   - Root `main.go`, or
   - `cmd/<name>/main.go` (prefer single cmd if only one exists).
3. For each target in the matrix, run:
   ```
   GOOS=<os> GOARCH=<arch> go build -ldflags "-s -w -X main.version=<ver>" -o <output> .
   ```
4. Collect all binaries into a staging directory (`release-assets/`).
5. Optionally create SHA256 checksums file (`checksums.txt`).

### Step 7b — Upload Assets

After the GitHub release is created (via the existing API call), upload
each binary using the GitHub Releases Upload API:

```
POST https://uploads.github.com/repos/{owner}/{repo}/releases/{id}/assets?name={filename}
Content-Type: application/octet-stream
Authorization: token <GITHUB_TOKEN>
```

This endpoint is already accessible with the same token used for release
creation — no additional CLI tools needed.

---

## Configuration

### config.json (optional)

```json
{
  "release": {
    "targets": [
      {"goos": "windows", "goarch": "amd64"},
      {"goos": "linux",   "goarch": "amd64"},
      {"goos": "darwin",  "goarch": "arm64"}
    ],
    "checksums": true,
    "compress": false
  }
}
```

### Flags

| Flag           | Description                              |
|----------------|------------------------------------------|
| `--no-assets`  | Skip binary cross-compilation            |
| `--targets`    | Comma-separated list: `windows/amd64,linux/arm64` |
| `--compress`   | Wrap each binary in a `.tar.gz` (Linux/macOS) or `.zip` (Windows) |
| `--checksums`  | Generate `checksums.txt` with SHA256 hashes (default: true) |

---

## Package Structure

| File                          | Responsibility                     |
|-------------------------------|------------------------------------|
| `release/assets.go`          | Cross-compile orchestration        |
| `release/assetstargets.go`   | Target matrix + config parsing     |
| `release/assetsupload.go`    | GitHub API upload                  |
| `release/assetschecksum.go`  | SHA256 checksum generation         |
| `constants/constants_assets.go` | Asset-related constants         |

---

## Dry-Run Support

When `--dry-run` is active, the asset step prints:

```
[dry-run] Would cross-compile 6 binaries:
  → gitmap_v2.14.0_windows_amd64.exe
  → gitmap_v2.14.0_windows_arm64.exe
  → gitmap_v2.14.0_linux_amd64
  → gitmap_v2.14.0_linux_arm64
  → gitmap_v2.14.0_darwin_amd64
  → gitmap_v2.14.0_darwin_arm64
[dry-run] Would upload 6 assets + checksums.txt
```

---

## Error Handling

- If `go build` fails for a specific target, log the error and continue
  with remaining targets. Report failures in the summary.
- If asset upload fails, retry once. On second failure, log and continue.
- Never abort the entire release because of an asset failure.

---

## Acceptance Criteria

1. `gitmap-v28 release v1.0.0` in a Go repo produces 6 binaries + checksums.
2. `gitmap-v28 release --dry-run` lists all planned binaries without building.
3. `gitmap-v28 release --no-assets` skips binary compilation entirely.
4. `gitmap-v28 release --targets windows/amd64,linux/amd64` builds only 2 binaries.
5. `gitmap-v28 release --compress` produces `.zip`/`.tar.gz` archives.
6. Non-Go repos are unaffected — no asset step runs.
7. Build failures for individual targets do not abort the release.
8. All uploaded assets appear on the GitHub release page.

## Cross-References (Generic Specifications)

| Topic | Generic Spec | Covers |
|-------|-------------|--------|
| Cross-compilation | [01-cross-compilation.md](../07-generic-release/01-cross-compilation.md) | Multi-platform Go build targets and naming |
| Release pipeline | [02-release-pipeline.md](../07-generic-release/02-release-pipeline.md) | Build-once constraint, compression, checksums, publish |
| Release assets | [05-release-assets.md](../07-generic-release/05-release-assets.md) | Asset naming conventions, archive format |
| Checksums | [04-checksums-verification.md](../07-generic-release/04-checksums-verification.md) | SHA-256 generation and verification |
