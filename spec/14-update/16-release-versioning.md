# Release Versioning

**Version:** 3.2.0
**Updated:** 2026-04-16

---

## Purpose

Define how version numbers are resolved, tags are created, and changelog entries are extracted for release automation.

---

## Version Resolution

The version is resolved with a 3-tier priority:

| Priority | Source | Example |
|----------|--------|---------|
| 1 | Explicit CLI argument | `release v1.2.0` |
| 2 | Bump flag | `release --bump minor` (1.1.0 → 1.2.0) |
| 3 | Current version from source | Read from `constants.go` or `version.go` |

### Normalization

All versions are normalized to ensure consistency:

```
1.2.0   → v1.2.0  (auto-prefix v)
v1.2.0  → v1.2.0  (already correct)
v01.2.0 → v1.2.0  (strip zero-padding)
  v1.2.0  → v1.2.0  (trim whitespace)
v1.2.0-beta.1 → v1.2.0-beta.1  (pre-release preserved)
```

### Go Implementation — NormalizeVersion

```go
import (
    "strconv"
    "strings"
)

// NormalizeVersion ensures consistent "vMAJOR.MINOR.PATCH" format.
func NormalizeVersion(v string) string {
    v = strings.TrimSpace(v)
    v = strings.TrimPrefix(v, "v")

    // Split and strip zero-padding from each component
    parts := strings.SplitN(v, "-", 2) // separate pre-release suffix
    components := strings.Split(parts[0], ".")
    for i, c := range components {
        if n, err := strconv.Atoi(c); err == nil {
            components[i] = strconv.Itoa(n) // "01" → "1"
        }
    }

    normalized := "v" + strings.Join(components, ".")
    if len(parts) > 1 {
        normalized += "-" + parts[1] // re-attach pre-release
    }
    return normalized
}

// CompareVersions returns true if a and b represent the same version
// after normalization.
func CompareVersions(a, b string) bool {
    return NormalizeVersion(a) == NormalizeVersion(b)
}
```

### Bash Implementation — normalize_version

```bash
normalize_version() {
    local v="$1"
    v="${v#v}"                  # strip v prefix
    v="$(echo "$v" | xargs)"   # trim whitespace
    echo "v$v"
}

# Usage in post-update verification

old_version="1.2.0"
new_version=$(<binary> version | xargs)

if [[ "$(normalize_version "$new_version")" == "$(normalize_version "$old_version")" ]]; then
    echo " !! Warning: version unchanged after update ($old_version)"
fi
```

---

## Version Detection from Source

The release system reads the current version from the Go source code. The version constant must follow this pattern:

### Expected Format

```go
// In constants.go or version.go
const Version = "1.2.0"
```

### Detection Script

```bash
detect_version() {
    local source_file="${1:-constants.go}"
    local version

    # Extract version from Go source
    version=$(grep -oP 'const\s+Version\s*=\s*"([^"]+)"' "$source_file" | \
              grep -oP '"[^"]+"' | tr -d '"')

    if [ -z "$version" ]; then
        echo "::error::Could not detect version from $source_file"
        exit 1
    fi

    echo "$version"
}
```

### Go Implementation — DetectVersion

```go
import (
    "fmt"
    "os"
    "regexp"
    "strings"
)

var versionPattern = regexp.MustCompile(`const\s+Version\s*=\s*"([^"]+)"`)

// DetectVersion reads the Version constant from a Go source file.
func DetectVersion(filePath string) (string, error) {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return "", fmt.Errorf("read %s: %w", filePath, err)
    }

    matches := versionPattern.FindSubmatch(data)
    if matches == nil {
        return "", fmt.Errorf("no Version constant found in %s", filePath)
    }

    return NormalizeVersion(string(matches[1])), nil
}
```

---

## Auto-Bump Logic

When `--bump <level>` is used instead of an explicit version, the system calculates the next version:

### Go Implementation — BumpVersion

```go
import (
    "fmt"
    "strconv"
    "strings"
)

// BumpLevel represents the semantic version component to increment.
type BumpLevel int

const (
    BumpPatch BumpLevel = iota
    BumpMinor
    BumpMajor
)

// ParseBumpLevel converts a string to a BumpLevel.
func ParseBumpLevel(s string) (BumpLevel, error) {
    switch strings.ToLower(s) {
    case "patch":
        return BumpPatch, nil
    case "minor":
        return BumpMinor, nil
    case "major":
        return BumpMajor, nil
    default:
        return 0, fmt.Errorf("invalid bump level: %q (must be major, minor, or patch)", s)
    }
}

// BumpVersion increments the specified component of a semantic version.
// Input version may or may not have a "v" prefix.
// Pre-release suffixes are stripped during bump.
func BumpVersion(current string, level BumpLevel) (string, error) {
    normalized := NormalizeVersion(current)
    clean := strings.TrimPrefix(normalized, "v")

    // Strip pre-release suffix for bump calculation
    base := strings.SplitN(clean, "-", 2)[0]
    parts := strings.Split(base, ".")

    if len(parts) != 3 {
        return "", fmt.Errorf("version %q does not have 3 components", current)
    }

    major, err := strconv.Atoi(parts[0])
    if err != nil {
        return "", fmt.Errorf("invalid major version: %s", parts[0])
    }
    minor, err := strconv.Atoi(parts[1])
    if err != nil {
        return "", fmt.Errorf("invalid minor version: %s", parts[1])
    }
    patch, err := strconv.Atoi(parts[2])
    if err != nil {
        return "", fmt.Errorf("invalid patch version: %s", parts[2])
    }

    switch level {
    case BumpMajor:
        major++
        minor = 0
        patch = 0
    case BumpMinor:
        minor++
        patch = 0
    case BumpPatch:
        patch++
    }

    return fmt.Sprintf("v%d.%d.%d", major, minor, patch), nil
}
```

### Bash Implementation — bump_version

```bash
bump_version() {
    local current="$1"
    local level="$2"

    # Strip v prefix and pre-release suffix
    local base="${current#v}"
    base="${base%%-*}"

    IFS='.' read -r major minor patch <<< "$base"

    case "$level" in
        major)
            major=$((major + 1))
            minor=0
            patch=0
            ;;
        minor)
            minor=$((minor + 1))
            patch=0
            ;;
        patch)
            patch=$((patch + 1))
            ;;
        *)
            echo "::error::Invalid bump level: $level (must be major, minor, or patch)"
            exit 1
            ;;
    esac

    echo "v${major}.${minor}.${patch}"
}

# Usage:

# bump_version "v1.2.3" "minor"  →  "v1.3.0"

# bump_version "1.0.0" "major"   →  "v2.0.0"

```

---

## Semantic Versioning

Versions follow [SemVer 2.0.0](https://semver.org):

```
v<major>.<minor>.<patch>[-<prerelease>]
```

| Bump | When | Example |
|------|------|---------|
| Major | Breaking changes | `v1.0.0` → `v2.0.0` |
| Minor | New features (backward-compatible) | `v1.1.0` → `v1.2.0` |
| Patch | Bug fixes only | `v1.2.0` → `v1.2.1` |

### Pre-Release Versions

Pre-release versions use a hyphen suffix following SemVer:

```
v1.2.0-alpha.1    # Early development
v1.2.0-beta.1     # Feature-complete, testing
v1.2.0-rc.1       # Release candidate
```

Pre-release detection in CI:

```bash
IS_PRERELEASE="false"
if [[ "$VERSION" == *-* ]]; then
    IS_PRERELEASE="true"
fi
```

---

## Tagging

After version resolution, the release system:

1. **Verifies** the tag does not already exist locally or remotely.
2. **Creates** a lightweight Git tag at `HEAD`.
3. **Pushes** the tag to the remote.

### Tag Verification Implementation

```bash
verify_tag_available() {
    local tag="$1"

    # Check local tags
    if git tag -l "$tag" | grep -q "$tag"; then
        echo "::error::Tag $tag already exists locally"
        exit 1
    fi

    # Check remote tags
    if git ls-remote --tags origin "$tag" | grep -q "$tag"; then
        echo "::error::Tag $tag already exists on remote"
        exit 1
    fi

    echo "✅ Tag $tag is available"
}

create_and_push_tag() {
    local tag="$1"

    verify_tag_available "$tag"

    git tag "$tag"

    if ! git push origin "$tag"; then
        echo "::error::Failed to push tag $tag — cleaning up"
        git tag -d "$tag"
        exit 1
    fi

    echo "✅ Tag $tag pushed successfully"
}

# Usage:

# create_and_push_tag "v1.3.0"

```

### Go Implementation — TagVerification

```go
import (
    "fmt"
    "os/exec"
    "strings"
)

// VerifyTagAvailable checks that a tag does not exist locally or remotely.
func VerifyTagAvailable(tag string) error {
    // Check local
    out, _ := exec.Command("git", "tag", "-l", tag).Output()
    if strings.TrimSpace(string(out)) == tag {
        return fmt.Errorf("tag %s already exists locally", tag)
    }

    // Check remote
    out, _ = exec.Command("git", "ls-remote", "--tags", "origin", tag).Output()
    if strings.Contains(string(out), tag) {
        return fmt.Errorf("tag %s already exists on remote", tag)
    }

    return nil
}

// CreateAndPushTag creates a lightweight tag and pushes it.
// On push failure, the local tag is deleted to prevent stale state.
func CreateAndPushTag(tag string) error {
    if err := VerifyTagAvailable(tag); err != nil {
        return err
    }

    if err := exec.Command("git", "tag", tag).Run(); err != nil {
        return fmt.Errorf("create tag %s: %w", tag, err)
    }

    if err := exec.Command("git", "push", "origin", tag).Run(); err != nil {
        // Cleanup: remove local tag on push failure
        _ = exec.Command("git", "tag", "-d", tag).Run()
        return fmt.Errorf("push tag %s: %w", tag, err)
    }

    return nil
}
```

---

## Changelog Extraction

The release pipeline extracts the relevant section from `changelog.md` for the release body:

```bash
extract_changelog() {
    local version="$1"
    local changelog="${2:-changelog.md}"

    if [ ! -f "$changelog" ]; then
        echo "Release $version"
        return
    fi

    local entry
    entry=$(awk -v ver="$version" '
        /^## / {
            if (found) exit
            if (index($0, ver)) found=1
        }
        found { print }
    ' "$changelog" 2>/dev/null)

    if [ -z "$entry" ]; then
        echo "Release $version"
    else
        echo "$entry"
    fi
}

# Usage:

# extract_changelog "v1.3.0" > /tmp/changelog-entry.md

```

### Changelog Format

```markdown

## v1.2.0 — Feature Title (2026-04-08)

### Improvements

- Added feature X for better performance.
- Updated Y to support Z.

### Bug Fixes

- Fixed crash when input is empty.
```

### Synchronization Requirement

Three sources must always be in sync:

| Source | Location | Purpose |
|--------|----------|---------|
| `Version` constant | Source code (`constants.go`) | Compiled into binary |
| `changelog.md` | Repository root | Human-readable history |
| Release metadata | `.release/latest.json` or tags | CI/CD and tooling |

When bumping a version:

1. Update the `Version` constant in source code.
2. Add the new section to `changelog.md`.
3. Update any metadata files (e.g., `latest.json`).

All three changes must happen in the **same commit** that is tagged.

---

## Version Source Update

After resolving the new version (explicit or auto-bumped), update the source constant:

### Go Implementation — UpdateVersionInSource

```go
import (
    "fmt"
    "os"
    "regexp"
    "strings"
)

// UpdateVersionInSource replaces the Version constant in a Go source file.
func UpdateVersionInSource(filePath, newVersion string) error {
    // Strip v prefix for source constant (source uses "1.2.0", not "v1.2.0")
    clean := strings.TrimPrefix(NormalizeVersion(newVersion), "v")

    data, err := os.ReadFile(filePath)
    if err != nil {
        return fmt.Errorf("read %s: %w", filePath, err)
    }

    pattern := regexp.MustCompile(`(const\s+Version\s*=\s*)"[^"]+"`)
    if !pattern.Match(data) {
        return fmt.Errorf("no Version constant found in %s", filePath)
    }

    updated := pattern.ReplaceAll(data, []byte(fmt.Sprintf(`${1}"%s"`, clean)))
    if err := os.WriteFile(filePath, updated, 0644); err != nil {
        return fmt.Errorf("write %s: %w", filePath, err)
    }

    return nil
}
```

### Bash Implementation

```bash
update_version_in_source() {
    local file="$1"
    local new_version="$2"

    # Strip v prefix for source constant
    local clean="${new_version#v}"

    sed -i "s/const Version = \"[^\"]*\"/const Version = \"$clean\"/" "$file"

    # Verify the update
    if ! grep -q "const Version = \"$clean\"" "$file"; then
        echo "::error::Failed to update version in $file"
        exit 1
    fi

    echo "✅ Updated version to $clean in $file"
}
```

---

## Release Metadata Files

Maintain a `.release/latest.json` file for programmatic version queries. This file is the **primary mechanism** for the updater binary and build scripts to detect the current published version.

### Schema

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["Version", "Tag", "Date", "Commit"],
  "properties": {
    "Version": {
      "type": "string",
      "description": "Clean semantic version without v prefix",
      "example": "1.2.0"
    },
    "Tag": {
      "type": "string",
      "description": "Git tag with v prefix",
      "example": "v1.2.0"
    },
    "Date": {
      "type": "string",
      "format": "date",
      "description": "UTC release date in YYYY-MM-DD format"
    },
    "Commit": {
      "type": "string",
      "description": "Full or short commit SHA",
      "example": "abc123def456"
    },
    "Assets": {
      "type": "array",
      "description": "List of release asset filenames (optional)",
      "items": { "type": "string" }
    }
  },
  "additionalProperties": false
}
```

### Example

```json
{
    "Version": "1.2.0",
    "Tag": "v1.2.0",
    "Date": "2026-04-08",
    "Commit": "abc123def456",
    "Assets": [
        "gitmap-linux-amd64.tar.gz",
        "gitmap-windows-amd64.zip",
        "gitmap-updater-linux-amd64.tar.gz",
        "docs-site.zip",
        "checksums.txt",
        "install.ps1",
        "install.sh"
    ]
}
```

### Hosting Location

The file lives in the repository at `.release/latest.json` and is committed alongside version bumps. It is **not** a release asset — it is read from the repository's default branch.

### How It's Used

| Consumer | Usage |
|----------|-------|
| Build scripts | Detect current version without parsing Go source |
| Updater binary | **Primary**: uses GitHub API `releases/latest`. **Fallback**: reads `.release/latest.json` from the raw GitHub URL |
| CI pipeline | Verify version synchronization between source, changelog, and metadata |

### Primary vs. Fallback for Updates

The updater binary's version resolution order:

1. **CLI flag** `--version v1.3.0` → use directly
2. **GitHub API** `GET /repos/{owner}/{repo}/releases/latest` → parse `tag_name`
3. **Fallback**: fetch `https://raw.githubusercontent.com/{owner}/{repo}/main/.release/latest.json` → parse `Tag`

The GitHub API is preferred because it always reflects the latest published release. The `latest.json` file is a fallback for environments where the API is blocked or rate-limited.

### Update Script

```bash
update_release_metadata() {
    local version="$1"
    local tag="v${version#v}"
    local commit="${GITHUB_SHA:-$(git rev-parse HEAD)}"
    local date="$(date -u '+%Y-%m-%d')"

    mkdir -p .release
    cat > .release/latest.json << EOF
{
    "Version": "${version#v}",
    "Tag": "$tag",
    "Date": "$date",
    "Commit": "$commit"
}
EOF
    echo "✅ Updated .release/latest.json"
}
```

---

## Release Branch Strategy

For release preparation:

```
main → release/1.2.0 → tag v1.2.0 → merge back to main
```

1. Create `release/1.2.0` branch from `main`.
2. Bump version, update changelog on the branch.
3. Push the branch — CI runs tests.
4. Tag the branch head as `v1.2.0`.
5. CI publishes the release.
6. Merge back to `main`.

### Rollback

If the release push fails:

1. Switch back to the original branch.
2. Force-delete the local release branch.
3. Delete the local tag.

```bash
git checkout main
git branch -D release/1.2.0
git tag -d v1.2.0
```

---

## Constraints

- Version must be bumped in source code **before** tagging.
- Tags must be lightweight (not annotated) unless signing is required.
- Changelog, version constant, and metadata must be updated atomically.
- Pre-release versions use a hyphen suffix: `v1.2.0-beta.1`.
- Version resolution must never fall back to "latest tag" without explicit verification that the tag matches the source `Version`.
- Auto-bump must strip pre-release suffixes before calculating the next version.
- Version detection must fail loudly if the source constant is missing or malformed.

---

## Cross-References

- [Release Assets](./14-release-assets.md) — Asset naming uses the resolved version
- [Checksums & Verification](./15-checksums-verification.md) — Checksum file versioning
- [CI/CD Release Body & Changelog](../12-cicd-pipeline-workflows/14-release-body-and-changelog.md) — How the release body uses changelog data
- [CI/CD GitHub Release Standard](../12-cicd-pipeline-workflows/04-github-release-standard.md) — Pre-release detection
- [Self-Update Overview](./02-self-update-overview.md) — How the CLI consumes versioned releases
- [Deploy Path Resolution](./03-deploy-path-resolution.md) — Where the versioned binary is deployed

---

*Release versioning — v3.2.0 — 2026-04-10*
