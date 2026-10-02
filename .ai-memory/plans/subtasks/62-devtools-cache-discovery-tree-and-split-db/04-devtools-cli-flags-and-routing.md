# Subtask 04: Devtools CLI Flags Integration & Routing

> **Parent Plan:** [62-devtools-cache-discovery-tree-and-split-db](../../completed/62-devtools-cache-discovery-tree-and-split-db.md)
> **Tracking Spec:** [199-devtools-cache-dynamic-discovery-tree-view-and-split-db.md](../../../../02-spec/21-app/199-devtools-cache-dynamic-discovery-tree-view-and-split-db.md)
> **Primary File Targets:** `cli/cmdos/os_dev_clean.go`, `cli/cmdos/os_dev_clean_flags.go`, `cli/cmdos/devclean_cmd.go`, `cli/osclean/dev_types.go`
> **Status:** `COMPLETED`
> **Owner:** Worker 02

---

## 1. Objective & Context

The developer tools cache subsystem (`gitmap clear devtools`, `gitmap clean devtools`, `gitmap devtool clear`) requires full parity with user expectations for flexible automation, visual tree inspection, and cache control:
1. **`--force` (`-f`):** Forces a deep multi-tier re-probe across runtime binaries and filesystems, invalidating Split-DB cache records.
2. **`--tree` (`-t`):** Requests hierarchical directory tree visualization rendering parent category nodes, root paths, sub-caches, file counts, and sizes.
3. **`--dry-run` (`-n`, `-d`):** Previews targeted paths and calculate reclaimable space without modifying or deleting files.
4. **`--only <cats>`:** Restricts discovery and cleaning to specific comma-separated ecosystems (e.g. `--only go,npm,pnpm`).
5. **`-y`, `--yes`, `/y`:** Bypasses interactive terminal confirmation for headless CI/CD execution and agent workflows.
6. **`--json`:** Emits structured JSON summary payload for programmatic consumers and AI subagents.
7. **`-v`, `--verbose`:** Outputs detailed notes, skipped paths, and subprocess diagnostic errors.

To preserve GitMap's strict architectural limit of `< 100 lines per file` and `<= 15 lines per function`, flag parsing must be cleanly extracted to `cli/cmdos/os_dev_clean_flags.go`.

---

## 2. CLI Flags Specification & Supported Variations

| Flag | Short / Variations | Value Type | Description |
| :--- | :--- | :--- | :--- |
| `--force` | `-f`, `force` | Boolean | Invalidate Split-DB cache and force fresh dynamic discovery |
| `--tree` | `-t`, `tree` | Boolean | Render hierarchical directory tree view with file counts and sizes |
| `--dry-run`| `-n`, `-d`, `dry-run` | Boolean | Preview space and files without performing actual file deletion |
| `--yes` | `-y`, `/y`, `yes` | Boolean | Bypass interactive confirmation prompt |
| `--json` | `json` | Boolean | Output structured JSON envelope to stdout |
| `--verbose`| `-v`, `verbose` | Boolean | Print individual path diagnostics and tool execution notes |
| `--only` | `--only=<cats>`, `--only <cats>` | String list | Comma-delimited list of ecosystems to filter (e.g. `go,npm`) |
| `--help` | `-h`, `help`, `/?` | Boolean | Print usage banner and flag documentation |

---

## 3. Data Structures & Flag Parsing Architecture

### 3.1 Extended Options Model (`cli/osclean/dev_types.go`)

```go
package osclean

// DevCleanOptions configures developer tools cache sweeping.
type DevCleanOptions struct {
	IsDryRun       bool     `json:"isDryRun"`
	HasAutoYes     bool     `json:"hasAutoYes"`
	IsJSON         bool     `json:"isJson"`
	IsVerbose      bool     `json:"isVerbose"`
	IsForce        bool     `json:"isForce"`
	IsTree         bool     `json:"isTree"`
	OnlyCategories []string `json:"onlyCategories,omitempty"`
}
```

### 3.2 Modular Flag Parsing (`cli/cmdos/os_dev_clean_flags.go`)

Flag parsing logic handles hyphen variations, flag arguments, and ignorable command tokens (`clear`, `clean`, `devtools`, etc.):

```go
package cmdos

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/osclean"
)

// ParseDevCleanOptions parses CLI arguments into DevCleanOptions.
func ParseDevCleanOptions(args []string) osclean.DevCleanOptions {
	opts := osclean.DevCleanOptions{}
	for i := 0; i < len(args); i++ {
		parseDevCleanFlagToken(args[i], args, &i, &opts)
	}
	return opts
}

func parseDevCleanFlagToken(raw string, args []string, idx *int, opts *osclean.DevCleanOptions) {
	a := strings.ToLower(raw)
	if isDevCleanIgnoredToken(a) {
		return
	}
	if parseDevCleanBooleans(a, opts) {
		return
	}
	parseDevCleanComplex(a, args, idx, opts)
}

func parseDevCleanBooleans(a string, opts *osclean.DevCleanOptions) bool {
	switch {
	case a == "--force" || a == "-f" || a == "force":
		opts.IsForce = true
		return true
	case a == "--tree" || a == "-t" || a == "tree":
		opts.IsTree = true
		return true
	case a == "--dry-run" || a == "-n" || a == "-d" || a == "dry-run":
		opts.IsDryRun = true
		return true
	case a == "--yes" || a == "-y" || a == "/y" || a == "yes":
		opts.HasAutoYes = true
		return true
	case a == "--json":
		opts.IsJSON = true
		return true
	case a == "--verbose" || a == "-v":
		opts.IsVerbose = true
		return true
	}
	return false
}

func parseDevCleanComplex(a string, args []string, idx *int, opts *osclean.DevCleanOptions) {
	if strings.HasPrefix(a, "--only=") {
		opts.OnlyCategories = strings.Split(strings.TrimPrefix(a, "--only="), ",")
		return
	}
	if a == "--only" && *idx+1 < len(args) {
		*idx++
		opts.OnlyCategories = strings.Split(args[*idx], ",")
	}
}
```

---

## 4. Command Dispatch & Routing Flow (`os_dev_clean.go`)

### 4.1 Orchestration Sequence

```
           [RunOSDevClean(args)]
                     |
            Is help requested?
            /               \
         YES                 NO
         /                     \
[Print Usage & Exit]     [Parse Options via ParseDevCleanOptions(args)]
                               |
                   Is interactive prompt needed?
                   (!IsDryRun && !HasAutoYes && !IsJSON)
                               |
                    Yes -> confirmDevClean()
                               |
              [Execute Clean / Discovery Workflow]
              • Pass opts.IsForce to Split-DB layer
              • Measure paths BEFORE wipe (prevent 0.01 MB bug)
              • If !IsDryRun: wipe caches
                               |
                   [Render Output Stage]
                   /         |          \
           IsJSON?       IsTree?       Standard ANSI
             /               |               \
   [json.MarshalIndent]  [RenderTree]  [RenderEnhancedCards]
```

### 4.2 Updated Help Documentation (`printOSDevCleanUsage`)

Both `devclean_cmd.go` and `os_dev_clean.go` must document `--force` and `--tree`:

```text
Usage: gitmap clear devtools [flags]
       gitmap devtools clear [flags]
       gitmap clean-dev [flags]

Flags:
  -f, --force        Invalidate Split-DB cache and force fresh dynamic discovery
  -t, --tree         Render hierarchical directory tree view with file counts and sizes
  -n, --dry-run      Preview space reclaimed without deleting
  -y, --yes          Bypass confirmation prompt
      --json         Output structured JSON summary
  -v, --verbose      Display individual subpaths and command notes
      --only <cats>  Comma-separated categories to clean (e.g. go,npm,pnpm,cargo)
  -h, --help         Show this documentation
```

---

## 5. Architectural & Coding Guideline Guardrails

1. **File Decomposition:** Keep `os_dev_clean.go` strictly under 100 lines by delegating parsing to `os_dev_clean_flags.go` and tree rendering to `os_dev_clean_tree.go`.
2. **Short Functions (<= 15 lines):** `RunOSDevClean`, `parseDevCleanBooleans`, and prompt checks must remain under 15 lines.
3. **Affirmative Booleans:** Strictly use positive prefixes (`IsForce`, `IsTree`, `IsDryRun`, `HasAutoYes`). Never use negative fields (`NoPrompt`, `SkipCache`).
4. **Command Aliases Support:** All existing aliases (`devtool`, `devtools`, `dev-tools`, `devtools-cache`, `clean-dev`) must seamlessly accept and forward all flags.

---

## 6. Verification & Acceptance Criteria

- **AC-1:** Running `gitmap clear devtools -f` or `--force` correctly sets `opts.IsForce = true`.
- **AC-2:** Running `gitmap clear devtools -t` or `--tree` correctly sets `opts.IsTree = true` and invokes tree rendering.
- **AC-3:** Running `gitmap clear devtools -f -t -n` combines force refresh, tree view, and dry-run without confirmation prompt.
- **AC-4:** Flag `--only go,npm` properly parses comma-separated lists and filters target categories.
- **AC-5:** No `go build` or `go test` invocations per Rule R1; passes `python 03-ai-scripts/05-guideline-autofixer.py` with zero guideline violations.
