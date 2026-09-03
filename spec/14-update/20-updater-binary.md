# Updater Binary

**Version:** 3.2.0
**Updated:** 2026-04-16

---

## Purpose

Define the architecture, interface, and behavior of the standalone updater binary (`<binary>-updater`). This is a separate Go module that handles downloading and installing updates when the main binary cannot replace itself (e.g., on Windows due to file locks).

---

## Why a Separate Binary?

The main binary cannot overwrite itself on Windows while it is running. The updater binary solves this by:

1. Being invoked by the main binary's `update` command
2. Downloading the new version
3. Replacing the main binary using the rename-first strategy
4. Verifying the update succeeded

The updater is also cross-compiled and distributed as a release asset alongside the main binary.

---

## Module Structure

The updater is a **separate Go module** with its own `go.mod`:

```
<binary>-updater/
├── go.mod                  # Independent module (not a workspace member)
├── go.sum
├── main.go                 # Entry point
├── cmd/
│   └── root.go             # CLI command definitions
├── updater/
│   ├── download.go         # GitHub API + asset download
│   ├── install.go          # Rename-first deploy logic
│   └── verify.go           # Post-install version verification
├── version/
│   └── version.go          # Embedded version constants
├── winres.json             # Windows icon/version resources
└── assets/
    └── icon-updater.png    # Distinct icon from main binary
```

---

## CLI Interface

```
<binary>-updater [flags]

Flags:
  --install-dir <path>    Target directory for installation (required)
  --version <ver>         Specific version to install (default: latest)
  --repo <owner/repo>     GitHub repository (embedded at build time)
  --binary-name <name>    Name of the binary to update (embedded at build time)
  --skip-checksum         Skip SHA-256 verification (not recommended)
  --verbose               Enable detailed output
  --help                  Show usage
```

### Example Invocations

```bash

# Install latest version to a specific directory

<binary>-updater --install-dir /usr/local/bin/<binary>

# Install a specific version

<binary>-updater --install-dir /usr/local/bin/<binary> --version v1.3.0

# Called by the main binary's update command (typical flow)

<binary>-updater --install-dir "$INSTALL_DIR" --version "$TARGET_VERSION"
```

---

## Main Flow (`main()`)

```
1. Parse CLI flags
2. Resolve target version
   a. If --version provided → use it
   b. Otherwise → query GitHub API for latest release tag
3. Construct download URLs for:
   a. Binary archive (<binary>-<os>-<arch>.<ext>)
   b. checksums.txt
   c. docs-site.zip (if applicable)
4. Download all assets to a temp directory
5. Verify SHA-256 checksum of the binary archive
6. Extract binary from archive
7. Deploy using rename-first strategy:
   a. Rename existing binary to <binary>.old
   b. Move new binary to install directory
8. Download and extract docs-site.zip (if applicable)
9. Verify: run "<binary> version" and check output
10. Clean up temp directory and .old files
11. Print success summary
```

---

## Version Resolution (GitHub API)

When `--version` is not provided, query the GitHub Releases API:

```go
func GetLatestVersion(repo string) (string, error) {
    url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return "", err
    }
    req.Header.Set("Accept", "application/vnd.github.v3+json")
    req.Header.Set("User-Agent", "<binary>-updater")

    resp, err := httpClient.Do(req)
    if err != nil {
        return "", fmt.Errorf("GitHub API request failed: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode == 403 {
        return "", fmt.Errorf("GitHub API rate limited — try again later or use --version flag")
    }
    if resp.StatusCode != 200 {
        return "", fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
    }

    var release struct {
        TagName string `json:"tag_name"`
    }
    if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
        return "", fmt.Errorf("parse GitHub response: %w", err)
    }

    return release.TagName, nil
}
```

---

## Network Requirements

See [14-network-requirements.md](./21-network-requirements.md) for full details. Summary:

| Requirement | Value |
|-------------|-------|
| HTTP client timeout | 30 seconds (connect), 5 minutes (download) |
| Retry policy | 3 attempts with exponential backoff (1s, 2s, 4s) |
| User-Agent | `<binary>-updater/<version>` |
| Proxy support | Respect `HTTP_PROXY`, `HTTPS_PROXY`, `NO_PROXY` |
| TLS | System certificate store; no custom CA pinning |
| Redirects | Follow up to 10 redirects (GitHub release URLs redirect) |

---

## Build Integration

The updater is built as part of the release pipeline, using the **main module's `dist/` folder** as output:

```bash

# Build from updater module, output to main module's dist/

cd <binary>-updater
CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" \
  go build -ldflags "$LDFLAGS" \
  -o "../<binary>/dist/<binary>-updater-${os}-${arch}${ext}" .
```

### LDFLAGS

The updater embeds fewer constants than the main binary:

| Constant | Purpose |
|----------|---------|
| `Version` | Updater version (matches main binary release version) |
| `Commit` | Short commit SHA |
| `RepoPath` | GitHub repository (`owner/repo`) — used for API queries |
| `BinaryName` | Name of the main binary to update |

---

## Error Handling

| Error | Response |
|-------|----------|
| GitHub API unreachable | Print error with suggestion to use `--version` flag; exit 1 |
| GitHub API rate limited (403) | Print rate limit message; suggest waiting or using `--version`; exit 1 |
| Download fails (network error) | Retry up to 3 times with backoff; exit 1 if all fail |
| Checksum mismatch | Print expected vs actual hash; delete downloaded file; exit 1 |
| Rename-first fails (all retries) | Print error; leave existing binary untouched; exit 1 |
| Version verification fails | Warn but don't roll back (binary may have different output format) |
| Install directory doesn't exist | Create it with `os.MkdirAll` |
| Docs-site download fails | Warn but continue — docs are optional |

---

## Constraints

- The updater is a **separate Go module** — not part of the main binary's module
- It has its own `go.mod` and can have different dependencies
- Output binaries go to the **main module's `dist/` folder** during CI builds
- The updater never updates itself — only the main binary
- The updater has its own `winres.json` and icon for Windows branding
- The updater uses the same version number as the main binary release

---

## Cross-References

- [Self-Update Overview](./02-self-update-overview.md) — Overall update flow and strategy selection
- [Rename-First Deploy](./04-rename-first-deploy.md) — File replacement strategy used by the updater
- [Handoff Mechanism](./06-handoff-mechanism.md) — How the main binary invokes the updater on Windows
- [Network Requirements](./21-network-requirements.md) — HTTP client configuration and retry policies
- [Release Assets](./14-release-assets.md) — Asset naming convention for updater archives
- [Complete Workflow Reference](../12-cicd-pipeline-workflows/02-go-binary-deploy/04-complete-workflow-reference.md) — How the updater is built in CI

---

*Updater binary — v3.2.0 — 2026-04-13*
