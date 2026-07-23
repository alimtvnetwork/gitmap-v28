# Version Display and Help System

## Overview

Every CLI tool must provide consistent version output and a comprehensive help system. This document specifies the exact patterns for version display, help formatting, and how these integrate with the CI/CD pipeline.

---

## Version Display

### Version Constant

The version is compiled into the binary via `ldflags`:

```go
// constants/constants.go
const Version = "1.3.0"
```

Build command:

```bash
LDFLAGS="-s -w -X '<module>/constants.Version=$VERSION'"
go build -ldflags "$LDFLAGS" -o <binary> .
```

### `version` Command

The version command prints the current version and exits. It must produce **clean, machine-parseable output** on stdout (no decorations, no color):

```
$ <tool> version
1.3.0
```

Alias: `v` or `-v`

### Version at Startup

For commands that perform significant work (e.g., `scan`, `release`), print the version at the beginning:

```
$ <tool> scan

  <tool> v1.3.0

  Scanning repositories...
```

### Version at End of Operations

After long-running operations complete, print a summary that includes the version:

```
  Done. Scanned 42 repositories in 3.2s.
  <tool> v1.3.0
```

### Version Synchronization Checklist

Before every release, verify these are in sync:

1. `constants.Version` in source code
2. `CHANGELOG.md` has a matching entry
3. `.gitmap/release/<version>.json` metadata file
4. Documentation site changelog (if applicable)
5. Git tag matches the version

---

## Help System

### Top-Level Help

When the tool is run with no arguments or with `help`, print a usage summary:

```
$ <tool>

  <tool> v1.3.0 — <one-line description>

  Usage: <tool> <command> [flags]

  Core Commands:
    scan              Scan and index Git repositories
    clone             Clone repositories from your GitHub account
    ls                List indexed repositories

  Release Commands:
    release           Create a versioned release
    deploy            Deploy the latest build

  Utility Commands:
    config            View or edit configuration
    doctor            Run system health checks
    version (v)       Print version

  Run '<tool> <command> --help' for details on any command.
```

### Command-Level Help

Each command responds to `--help` or `-h` with structured documentation:

```
$ <tool> scan --help

  <tool> scan — Scan and index Git repositories

  Alias: s

  Usage:
    <tool> scan [flags]

  Flags:
    --path <dir>       Root directory to scan (default: current directory)
    --depth <n>        Maximum directory depth (default: 5)
    --verbose          Show detailed output

  Examples:
    $ <tool> scan
    $ <tool> scan --path ~/projects --depth 3
    $ <tool> scan --verbose

  See Also:
    ls — List indexed repositories
    clone — Clone repositories
```

### Help Text Source

Help text is embedded from Markdown files using `go:embed`:

```go
//go:embed helptext/help-scan.md
var helpScan string
```

Each help file follows this structure:

```markdown
# <tool> scan

<One-line description>

**Alias:** `s`

## Usage

\`\`\`
<tool> scan [flags]
\`\`\`

## Flags

| Flag | Description |
|------|-------------|
| `--path <dir>` | Root directory to scan (default: current directory) |
| `--depth <n>` | Maximum directory depth (default: 5) |

## Examples

    $ <tool> scan
    $ <tool> scan --path ~/projects --depth 3

## See Also

- ls — List indexed repositories
- clone — Clone repositories
```

### Help File Constraints

- Maximum 120 lines per file
- 2–3 realistic examples per command
- Examples show 3–8 lines of terminal simulation
- Standard headers: Alias, Usage, Flags, Examples, See Also
- Stored in `helptext/` directory, embedded via `go:embed`

---

## Integration with CI/CD

### Version in CI Builds

CI builds use `dev-<sha>` versioning:

```bash
VERSION="dev-${GITHUB_SHA::10}"
```

Release builds use the semantic version:

```bash
VERSION="${GITHUB_REF_NAME#release/}"
```

Both are injected via the same `-X` ldflag.

### Version in Release Body

The release pipeline includes the version in the GitHub Release body and verifies it matches the tag:

```yaml
- name: Verify version
  run: |
    BINARY_VERSION=$(./dist/<binary>-*-linux-amd64 version 2>/dev/null || echo "unknown")
    TAG_VERSION="${{ steps.version.outputs.version }}"
    if [ "v$BINARY_VERSION" != "$TAG_VERSION" ]; then
      echo "::error::Version mismatch: binary=$BINARY_VERSION tag=$TAG_VERSION"
      exit 1
    fi
```

### Help in Automated Tests

Test that every command's `--help` flag produces output without errors:

```go
func TestAllCommandsHaveHelp(t *testing.T) {
    commands := []string{"scan", "clone", "ls", "release", "config", "doctor"}
    for _, cmd := range commands {
        t.Run(cmd, func(t *testing.T) {
            // Verify help text is non-empty and contains expected headers
        })
    }
}
```

---

## Terminal Output Samples

### Version

```
$ <tool> version
1.3.0

$ <tool> v
1.3.0
```

### Help (no args)

```
$ <tool>

  <tool> v1.3.0 — Git repository manager and release automation toolkit

  Usage: <tool> <command> [flags]

  Core Commands:
    scan (s)          Scan and index Git repositories
    clone (cl)        Clone repositories from your account
    ls (l)            List indexed repositories
    cd                Jump to a repository directory

  Release Commands:
    release (r)       Create a versioned release
    deploy (dp)       Deploy the latest build
    list-releases     List all releases

  Tooling Commands:
    config (cfg)      View or edit configuration
    doctor (dr)       Run system health checks
    setup             Configure shell completions and ignore rules
    update (up)       Self-update to the latest version

  Info Commands:
    version (v)       Print version
    help              Show this help message
    docs              Open documentation website

  Run '<tool> <command> --help' for details on any command.
```

### Unknown Command

```
$ <tool> foobar

  Error: unknown command "foobar"

  Run '<tool> help' for a list of available commands.
```

---

## Constraints

- `version` output must be a single line with no prefix (e.g., `1.3.0`, not `version: 1.3.0`)
- Help text must work in 80-column terminals (no line wrapping issues)
- All command aliases are shown in parentheses: `scan (s)`
- Group commands logically: Core, Release, Tooling, Info
- Help text uses 2-space indentation for visual hierarchy
- `--help` and `-h` are intercepted before any command logic runs
- Unknown commands print to stderr and exit with code 1
