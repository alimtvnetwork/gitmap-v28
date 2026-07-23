# 01 — Cross-Compilation

## Purpose

Build a single CLI binary for every supported OS and architecture from
one CI runner. The output is a set of static, self-contained executables
with no runtime dependencies.

---

## Default Targets

A production CLI should support at minimum 6 targets:

| OS | Architecture | Binary Name | Archive Format |
|----|-------------|-------------|----------------|
| `windows` | `amd64` | `<binary>.exe` | `.zip` |
| `windows` | `arm64` | `<binary>.exe` | `.zip` |
| `linux` | `amd64` | `<binary>` | `.tar.gz` |
| `linux` | `arm64` | `<binary>` | `.tar.gz` |
| `darwin` | `amd64` | `<binary>` | `.tar.gz` |
| `darwin` | `arm64` | `<binary>` | `.tar.gz` |

Windows uses `.zip` because it is natively supported by PowerShell and
Explorer. Unix platforms use `.tar.gz` for permission preservation.

---

## Build Command

For Go CLI tools, the standard cross-compilation command is:

```bash
CGO_ENABLED=0 GOOS=<os> GOARCH=<arch> go build \
  -ldflags "-X '<module>/constants.Version=<version>' \
            -X '<module>/constants.RepoPath=<repo-path>'" \
  -o <output-path> .
```

### Key Flags

| Flag | Purpose |
|------|---------|
| `CGO_ENABLED=0` | Produce a fully static binary (no C dependencies) |
| `GOOS` | Target operating system |
| `GOARCH` | Target CPU architecture |
| `-ldflags -X` | Embed build-time constants (version, repo path, commit SHA) |
| `-o` | Output binary path |

### Embedded Constants

Embed at minimum these values at build time:

| Constant | Purpose |
|----------|---------|
| `Version` | Semantic version string (e.g., `1.2.0`) |
| `RepoPath` | Absolute path to the source repo (enables self-update) |
| `CommitSHA` | Git commit hash (optional, for diagnostics) |
| `BuildDate` | ISO 8601 build timestamp (optional) |

---

## Build Loop (CI)

The CI pipeline should iterate over all targets in a single job:

```yaml
strategy:
  matrix:
    include:
      - os: windows
        arch: amd64
      - os: windows
        arch: arm64
      - os: linux
        arch: amd64
      - os: linux
        arch: arm64
      - os: darwin
        arch: amd64
      - os: darwin
        arch: arm64
```

Or in a shell loop:

```bash
TARGETS="windows/amd64 windows/arm64 linux/amd64 linux/arm64 darwin/amd64 darwin/arm64"

for target in $TARGETS; do
    os="${target%/*}"
    arch="${target#*/}"
    ext=""
    [[ "$os" == "windows" ]] && ext=".exe"

    output="dist/<binary>-${os}-${arch}${ext}"

    CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" \
        go build -ldflags "$LDFLAGS" -o "$output" .

    echo "Built: $output ($(stat -c%s "$output" 2>/dev/null || stat -f%z "$output") bytes)"
done
```

---

## Output Directory Structure

All binaries are placed in a single `dist/` directory:

```
dist/
├── <binary>-windows-amd64.exe
├── <binary>-windows-arm64.exe
├── <binary>-linux-amd64
├── <binary>-linux-arm64
├── <binary>-darwin-amd64
└── <binary>-darwin-arm64
```

The `dist/` directory is the **single source of truth** for all
downstream packaging, checksumming, and publishing steps.

---

## Verification

After building each target, verify it is a valid executable:

```bash
file "$output"       # Should show correct architecture
ls -la "$output"     # Should have non-zero size
```

For the host-architecture binary, run a quick smoke test:

```bash
./<binary> version   # Should print the embedded version
```

---

## Constraints

- Never rebuild binaries after the build stage — all downstream steps
  operate on the already-built artifacts.
- Use the same Go toolchain version across all targets.
- Do not use `go install` for cross-compilation — use `go build`.
- Binary names must include OS and architecture for disambiguation.
- The `.exe` extension is mandatory for Windows targets only.

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
