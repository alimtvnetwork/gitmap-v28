# gitmap gomod

Rename the Go module path across the entire repository with branch safety.

## Alias

gm

## Usage

    gitmap gomod <new-module-path> [flags]

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| --ext \<exts\> | all files | Comma-separated extensions to filter |
| --dry-run | false | Preview changes without modifying |
| --no-merge | false | Stay on feature branch after commit |
| --no-tidy | false | Skip go mod tidy after replacement |
| --verbose | false | Print each file path as modified |

## Prerequisites

- Must be inside a Go project with go.mod

## Examples

### Example 1: Rename module path (full flow)

    gitmap gomod "github.com/neworg/myproject"

**Output:**

    Current module: github.com/oldorg/myproject
    New module:     github.com/neworg/myproject
    Creating branch gomod-rename...
    Replacing in 24 files...
      go.mod
      main.go
      cmd/root.go
      cmd/scan.go
      store/store.go
      ...
    Running go mod tidy... done
    Committing changes... done
    Merging into main... done
    ✓ Module renamed across 24 files

### Example 2: Dry-run with specific extensions

    gitmap gm "github.com/new/name" --ext "*.go,*.md" --dry-run

**Output:**

    [DRY RUN] Current module: github.com/old/name
    [DRY RUN] New module:     github.com/new/name
    [DRY RUN] Would modify:
      18 .go files
       3 .md files
    [DRY RUN] Would run go mod tidy
    [DRY RUN] Would commit and merge
    No changes made.

### Example 3: Rename without merge, verbose output

    gitmap gomod "github.com/new/name" --no-merge --verbose

**Output:**

    Current module: github.com/old/name
    Creating branch gomod-rename...
    [verbose] Replacing in cmd/root.go... 3 occurrences
    [verbose] Replacing in cmd/scan.go... 2 occurrences
    [verbose] Replacing in store/store.go... 1 occurrence
    [verbose] Replacing in main.go... 1 occurrence
    ...
    Running go mod tidy... done
    Committed on branch 'gomod-rename' (not merged)
    ✓ 24 files updated — merge manually when ready

### Example 4: Rename and skip go mod tidy

    gitmap gomod "github.com/new/name" --no-tidy

**Output:**

    Current module: github.com/old/name
    Replacing in 24 files... done
    Skipping go mod tidy (--no-tidy).
    Committing changes... done
    Merging into main... done
    ✓ Module renamed (tidy skipped)

## See Also

- [go-repos](go-repos.md) — List detected Go projects
- [scan](scan.md) — Scan directories to detect Go projects

## Scripting (JSON)

Discover this command from a script using the machine-readable help payload:

```bash
gitmap help --json --filter gomod
```

The JSON schema is published at `spec/08-json-schemas/help-json.schema.json` (v5.43.0+).
