# Component Specification: Workdir Resolution, Terminal UI Dashboard & Help Integration

> **Specification Identifier:** `02-spec/21-app/69-nodes-cfr-ui-async-delegation/02-component-spec.md`  
> **Parent Initiative:** `69-nodes-cfr-ui-async-delegation.md`  
> **Architecture Baseline:** `01-architecture-spec.md`  
> **Target Subtasks:**  
> - Subtask 03: Workdir Resolution, Relative Subdir Mirroring & Custom Destination  
> - Subtask 04: Terminal UI Dashboard Modernization, Help Documentation & Suggestion Footers  
> **Target Packages:** `cli/cmdnodes`, `cli/helptext`  
> **Release Target:** Minor Version Bump (`v6.454.0`)  
> **Status:** APPROVED & SPECIFIED  

---

## 1. Executive Summary & Problem Analysis

In distributed repository management across a fleet of heterogeneous developer workstations and remote servers, operators interact with GitMap via the `gitmap nodes cfr`, `gitmap nodes clone`, and `gitmap nodes cfrp` commands. While basic remote command dispatch exists, real-world fleet workflows exposed four fundamental component deficiencies:

1. **Unintelligent Work Directory Placement & Blind Mirroring Failures:**  
   When an engineer operates locally inside a nested project subfolder (e.g. `./presentations-repos` or `~/work/infrastructure/tools`), running `gitmap nodes cfr` dispatches clones blindly into remote root work folders (`D:\work` on Windows or `~/work` on Linux) instead of preserving the relative hierarchy (`./presentations-repos` and `~/work/presentations-repos`). If the remote target directory does not yet exist, remote shell execution immediately crashes with path-not-found errors.
2. **Missing Custom Destination Flags & Positional Flexibility:**  
   Operators have no unified mechanism to specify explicit destination directories. Standard directory flags (`-d`, `--dest`, `--dir`, `--target-dir`) and positional destination overrides (`gitmap nodes cfr <target> [dest]`) are either unparsed or collide with target arguments.
3. **Rudimentary Terminal UI & Missing Execution Route Visibility:**  
   The terminal interface does not expose pre-flight readiness metrics, machine reachability (`online`, `offline`, `auth_failed`), or remote GitMap binary versions before launching long operations. Crucially, operators are never shown the destination route mapping (`alias -> host:path`), leaving them blind to where files are being cloned.
4. **Undocumented CLI Subsystem & Absence of Actionable Guidance:**  
   The help screens in `cli/cmdnodes/nodes_clone_help.go` lack concrete examples for custom destinations, except-self modes, host key acceptance (`gitmap ssh trust`), and remote machine telemetry inspection (`gitmap machine --ssh`, `gitmap nodes ping`). Furthermore, `cli/helptext/catalog.go` fails to index `nodes-clone`, `nodes-cfr`, and related aliases. Post-execution results lack contextual next-step recommendations.

This specification establishes the technical component contracts, data structures, algorithms, CLI parsing logic, and presentation designs to resolve these challenges.

---

## 2. Component Architecture & Module Distribution

```
cli/
├── cmdnodes/
│   ├── nodes_clone_types.go     # Data contracts: NodesCloneDestination, NodesCloneOptions, PreFlightNodeInfo
│   ├── nodes_clone.go           # CLI flag parsing, option extraction, dispatch orchestration
│   ├── nodes_clone_file.go      # CWD relative path resolution, workdir detection, manifest staging
│   ├── nodes_clone_remote.go    # Remote command string synthesis with mkdir -p / New-Item pre-creation
│   ├── nodes_clone_table.go     # Pre-flight readiness table, banner box, route blueprint, results table, footer
│   ├── nodes_clone_help.go      # Modernized CLI help with rich examples and parameter documentation
│   └── nodes_clone_test.go      # Unit test suite verifying path resolution, flags, and table rendering
└── helptext/
    └── catalog.go               # Topic summaries registering nodes-clone, nodes-cfr, nodes-cfrp, cfr, cfrp
```

---

## 3. Data Contracts: `cli/cmdnodes/nodes_clone_types.go`

### 3.1 `NodesCloneDestination`

`NodesCloneDestination` encapsulates resolved destination directory metadata across both the local host and remote target machines:

```go
// NodesCloneDestination describes resolved destination directory settings for clone operations.
type NodesCloneDestination struct {
	TargetDir          string `json:"targetDir,omitempty"`          // Normalized explicit target dir or empty if default
	RelativeSubdir     string `json:"relativeSubdir,omitempty"`      // Relative subpath if CWD is inside work directory
	HasCustomTargetDir bool   `json:"hasCustomTargetDir"`            // True if user specified -d/--dest/--dir/positional dest
	IsInsideWorkDir    bool   `json:"isInsideWorkDir"`               // True if local CWD is inside registered work directory
}
```

#### Invariants & Field Semantics:
- `IsInsideWorkDir`: `true` if and only if `filepath.Clean(cwd)` is equal to or a child subdirectory of `filepath.Clean(localWorkDir)`.
- `RelativeSubdir`: Normalized relative subfolder path using forward slashes (`/`) for cross-platform manipulation (e.g. `presentations-repos` or `internal/tools`). Empty `""` if CWD is the work directory root or outside work directory.
- `HasCustomTargetDir`: `true` if the user passed `-d`, `--dest`, `--dir`, `--target-dir`, `--dest=<path>`, `--dir=<path>`, `--target-dir=<path>`, or an explicit positional destination argument.
- `TargetDir`: Cleaned destination string if `HasCustomTargetDir` is true; otherwise empty `""` (indicating automatic resolution based on `RelativeSubdir`).

---

### 3.2 Extended `NodesCloneOptions`

`NodesCloneOptions` carries all parsed CLI arguments and operational configuration:

```go
// NodesCloneOptions stores configuration for fleet clone execution.
type NodesCloneOptions struct {
	Kind               NodesCloneKind        `json:"kind"`
	TargetFilter       string                `json:"targetFilter,omitempty"`
	ExcludeFilter      string                `json:"excludeFilter,omitempty"`
	ExceptOS           string                `json:"exceptOS,omitempty"`
	TargetOS           string                `json:"targetOS,omitempty"`
	TargetDir          string                `json:"targetDir,omitempty"`
	RelativeSubdir     string                `json:"relativeSubdir,omitempty"`
	HasCustomTargetDir bool                  `json:"hasCustomTargetDir"`
	IsInsideWorkDir    bool                  `json:"isInsideWorkDir"`
	Destination        NodesCloneDestination `json:"destination"`
	IsDryRun           bool                  `json:"isDryRun"`
	IsJSON             bool                  `json:"isJSON"`
	IsSkipLocal        bool                  `json:"isSkipLocal"`
	RawArgs            []string              `json:"rawArgs"`
	PassArgs           []string              `json:"passArgs"`
	DetectedFile       string                `json:"detectedFile,omitempty"`
	HasFile            bool                  `json:"hasFile"`
}
```

---

## 4. Intelligent Work Directory & Path Resolution Engine

### 4.1 Local Work Directory Detection

GitMap discovers the local active work directory through database configuration with platform fallbacks:

```go
// ResolveLocalWorkDir returns the active local work directory.
func ResolveLocalWorkDir() string {
	dbConn, err := store.OpenDefault()
	if err == nil {
		defer dbConn.Close()
		wd, errGet := dbConn.GetDefaultWorkDir()
		if errGet == nil && wd != nil && wd.AbsolutePath != "" {
			return filepath.Clean(wd.AbsolutePath)
		}
	}
	if isWindowsLocalHost() {
		return `D:\work`
	}
	home := os.Getenv("HOME")
	if home == "" {
		return "/tmp"
	}
	return filepath.Join(home, "work")
}
```

---

### 4.2 CWD Relative Subdirectory Calculation

The resolver inspects `os.Getwd()` relative to `ResolveLocalWorkDir()`:

```go
// ResolveCwdRelativeSubDir inspects CWD and returns relative path if inside workdir.
func ResolveCwdRelativeSubDir(workDir string) (string, bool) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", false
	}
	cleanCwd := filepath.Clean(cwd)
	cleanWork := filepath.Clean(workDir)

	if !isPathInsideBase(cleanCwd, cleanWork) {
		return "", false
	}
	rel, errRel := filepath.Rel(cleanWork, cleanCwd)
	if errRel != nil || rel == "." || strings.HasPrefix(rel, "..") {
		return "", true
	}
	return filepath.ToSlash(rel), true
}

func isPathInsideBase(target, base string) bool {
	if strings.EqualFold(target, base) {
		return true
	}
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(rel, "..") && rel != "."
}
```

---

### 4.3 Remote Target Directory Synthesis & Slash Normalization

For each candidate remote machine, the destination path is computed based on its registered OS:

```go
// ResolveRemoteWorkDir returns the default base work directory for a remote OS.
func ResolveRemoteWorkDir(osType string) string {
	if isWindowsOS(osType) {
		return `D:\work`
	}
	return "~/work"
}

// ResolveRemoteTargetDir constructs the target directory on a remote node.
func ResolveRemoteTargetDir(osType string, dest NodesCloneDestination) string {
	if dest.HasCustomTargetDir && dest.TargetDir != "" {
		return formatRemoteCustomPath(osType, dest.TargetDir)
	}
	baseDir := ResolveRemoteWorkDir(osType)
	if dest.RelativeSubdir == "" {
		return baseDir
	}
	if isWindowsOS(osType) {
		cleanSub := strings.ReplaceAll(dest.RelativeSubdir, "/", "\\")
		return baseDir + "\\" + cleanSub
	}
	cleanSub := strings.Trim(dest.RelativeSubdir, "/")
	return strings.TrimRight(baseDir, "/") + "/" + cleanSub
}

func formatRemoteCustomPath(osType, customPath string) string {
	if isWindowsOS(osType) {
		return strings.ReplaceAll(customPath, "/", "\\")
	}
	return strings.ReplaceAll(customPath, "\\", "/")
}

func isWindowsOS(osType string) bool {
	low := strings.ToLower(osType)
	return low == "windows" || low == "win"
}
```

---

### 4.4 Decision Matrix: Resolution Scenarios

| Scenario | Local CWD | Local WorkDir | User Flag / Dest | Remote Windows Node | Remote Linux Node |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Inside WorkDir Subfolder** | `./presentations-repos` | `D:\work` | *(none)* | `./presentations-repos` | `~/work/presentations-repos` |
| **Inside Nested Subfolder** | `./internal\tool` | `D:\work` | *(none)* | `./internal\tool` | `~/work/internal/tool` |
| **Inside WorkDir Root** | `D:\work` | `D:\work` | *(none)* | `D:\work` | `~/work` |
| **Outside WorkDir** | `C:\Users\Admin` | `D:\work` | *(none)* | `D:\work` | `~/work` |
| **Outside WorkDir (Linux)** | `/home/dev/downloads` | `~/work` | *(none)* | `D:\work` | `~/work` |
| **Explicit Flag Override** | `./presentations-repos` | `D:\work` | `-d ./custom` | `./custom` | `~/work/custom` |
| **Explicit Flag Override** | `./presentations-repos` | `D:\work` | `--dest=~/work/special` | `./special` | `~/work/special` |
| **Positional Destination** | `C:\Temp` | `D:\work` | `./special` | `./special` | `~/work/special` |

---

### 4.5 Remote Target Directory Pre-Creation

To prevent remote shell failures when target folders do not yet exist, remote execution strings synthesize pre-creation commands:

#### Windows PowerShell Command Pattern:
```powershell
if (!(Test-Path -Path '<targetDir>')) { New-Item -ItemType Directory -Path '<targetDir>' -Force | Out-Null }; Set-Location '<targetDir>'; gitmap <kind> <args> --json
```

#### Linux / Unix Bash Command Pattern:
```bash
mkdir -p '<targetDir>' && cd '<targetDir>' && gitmap <kind> <args> --json
```

#### Implementation Contracts in `cli/cmdnodes/nodes_clone_remote.go`:
```go
func buildWindowsWorkDirExecString(kindStr, args, targetDir string) string {
	execDir := `D:\work`
	if targetDir != "" {
		execDir = targetDir
	}
	execArgs := appendJSONFlag(args)
	return fmt.Sprintf(
		"if (!(Test-Path -Path \"%s\")) { New-Item -ItemType Directory -Path \"%s\" -Force | Out-Null }; Set-Location \"%s\"; gitmap %s %s",
		execDir, execDir, execDir, kindStr, execArgs,
	)
}

func buildUnixWorkDirExecString(kindStr, args, targetDir string) string {
	execDir := "~/work"
	if targetDir != "" {
		execDir = targetDir
	}
	execArgs := appendJSONFlag(args)
	return fmt.Sprintf("mkdir -p %s && cd %s && gitmap %s %s", execDir, execDir, kindStr, execArgs)
}
```

---

## 5. CLI Flags & Destination Parsing Specification

### 5.1 Supported CLI Syntax & Flag Variations
The argument parser must handle all standard flag permutations and positional arguments:
1. **Flag Variations:**
   - `-d <path>`, `--dest <path>`, `--dir <path>`, `--target-dir <path>`
   - `--dest=<path>`, `--dir=<path>`, `--target-dir=<path>`
2. **Positional Arguments:**
   - `gitmap nodes cfr <target> <destination>`
   - `gitmap nodes clone <repo-url> <destination>`
   - `gitmap nodes cfrp <manifest.json> <destination>`
   - `gitmap nodes cfr except-self <target> <destination>`

### 5.2 Disambiguation Algorithm: Target vs Destination

To prevent targets from being erroneously consumed as destinations:
- **Target Identification:**
  - An argument ending in `.json` that exists on disk is a **Manifest Target**.
  - An argument beginning with `http://`, `https://`, or `git@` is a **URL Target**.
  - An argument containing commas (e.g. `repo1,repo2`) is a **Multi-Repo Target**.
  - The first non-flag, non-subcommand token is identified as the **Primary Target**.
- **Destination Identification:**
  - Any argument following `-d`, `--dest`, `--dir`, or `--target-dir` is an **Explicit Destination**.
  - Any second non-flag positional argument (when the first argument is a Target and does not start with `-`) is an **Explicit Positional Destination**.

```go
func parseDestFlagValue(raw []string, idx int) (string, int, bool) {
	arg := raw[idx]
	for _, prefix := range []string{"--dir=", "--dest=", "--target-dir="} {
		if strings.HasPrefix(arg, prefix) {
			return strings.TrimPrefix(arg, prefix), idx, true
		}
	}
	if isDestFlag(arg) && idx+1 < len(raw) {
		return raw[idx+1], idx + 1, true
	}
	return "", idx, false
}

func isDestFlag(arg string) bool {
	return arg == "-d" || arg == "--dir" || arg == "--dest" || arg == "--target-dir"
}
```

---

## 6. Modernized Terminal UI Dashboard (`cli/cmdnodes/nodes_clone_table.go`)

### 6.1 Visual Design & Observability Flow

```
┌────────────────────────────────────────────────────────────────────────────────────────┐
│  Phase 1: Pre-Flight Readiness Summary Table                                           │
│  ▸ Fleet Readiness: 3 registered node(s) | 2 online | 1 offline                        │
│                                                                                        │
│    NODE (ALIAS)     HOST               OS         VERSION        DESTINATION   STATUS  │
│    ----------------------------------------------------------------------------------  │
│    w1               node-main       windows    v6.454.0       ./sub   ● online│
│    w2               gateway-node1       linux      v6.453.0       ~/work/sub    ● online│
│    w3               node-w3       linux      -              -             ○ off   │
└────────────────────────────────────────┬───────────────────────────────────────────────┘
                                         ▼
┌────────────────────────────────────────────────────────────────────────────────────────┐
│  Phase 2: Dispatch Banner & Route Mapping Card                                         │
│    ┌────────────────────────────────────────────────────────────────────────────────┐  │
│    │ GITMAP FLEET NODES CFR (CLONE-FIX-REPO) DISPATCH                               │  │
│    └────────────────────────────────────────────────────────────────────────────────┘  │
│    • Mode:        cfr (clone-fix-repo)                                                 │
│    • Target:      https://github.com/org/repo.git                                      │
│    • Workdir:     D:\work (preserved relative: sub)                                    │
│    • Scope:       Local host + 2 active remote worker(s)                               │
│    • Dispatch:    w1 -> node-main:./sub                                       │
│                   w2 -> gateway-node1:~/work/sub                                        │
└────────────────────────────────────────┬───────────────────────────────────────────────┘
                                         ▼
┌────────────────────────────────────────────────────────────────────────────────────────┐
│  Phase 3: Execution Telemetry Table & Dynamic Suggestions Footer                       │
│    NODE (ALIAS)     HOST        ROLE    STATUS     VERSION   DURATION   DETAILS        │
│    ----------------------------------------------------------------------------------  │
│    local (current)  127.0.0.1   master  ● success  v6.454.0  1250ms     cloned ok      │
│    w1               gateway-node worker  ● success  v6.454.0  1840ms     cloned ok      │
│    w2               node-alias worker  ● success  v6.453.0  2120ms     cloned ok      │
│    ----------------------------------------------------------------------------------  │
│  ✔ Fleet CFR Summary: 3/3 node(s) completed successfully (0 failed)                    │
│                                                                                        │
│  [tip] Fleet Operations & Suggested Commands:                                          │
│    • Ping Fleet Nodes:          gitmap nodes ping                                      │
│    • Query Machine Telemetry:   gitmap machine --ssh                                   │
│    • Inspect Node Connection:   gitmap ssh test <alias>                                │
│    • Rerun Except Local Host:   gitmap nodes cfr <target> --except-self                │
└────────────────────────────────────────────────────────────────────────────────────────┘
```

---

### 6.2 Pre-Flight Readiness Table Specification

- **Header Line:** `  ▸ Fleet Readiness: %d registered node(s) | %d online | %d offline\n\n`
- **Columns & Fixed Minimum Widths:**
  - `NODE (ALIAS)`: 16 chars
  - `HOST`: 18 chars
  - `OS`: 10 chars
  - `VERSION`: 14 chars (discovered version or `"-"`)
  - `DESTINATION`: 32 chars (resolved remote destination path or `"-"`)
  - `STATUS`: 14 chars (formatted ANSI badge)
- **Status Badges:**
  - `● online` (`constants.ColorGreen` + `● online` + `constants.ColorReset`)
  - `○ offline` (`constants.ColorYellow` + `○ offline` + `constants.ColorReset`)
  - `▲ auth_failed` (`constants.ColorMagenta` + `▲ auth_failed` + `constants.ColorReset`)
  - `✗ unreachable` (`constants.ColorRed` + `✗ unreachable` + `constants.ColorReset`)

---

### 6.3 Dispatch Banner & Route Mapping Card

- **Header Box:** Enclosed in double-byte box borders (`┌───┐`, `│   │`, `└───┘`).
- **Metadata Fields:**
  - `• Mode:` Formatted kind (`clone (multi-node clone)`, `cfr (clone-fix-repo)`, `cfrp (clone-fix-repo-pub)`).
  - `• Target:` Target repository, URL, or manifest name.
  - `• Workdir:` Formatted workdir with relative note (e.g. `D:\work (preserved relative: presentations-repos)` or `./custom (custom destination)`).
  - `• Scope:` Worker breakdown (`Local host + %d active remote worker(s)` or `%d active remote worker(s) (remote-only)`).
  - `• Dispatch:` Direct route mapping for each online node:
    ```text
    • Dispatch:    w1 -> node-main:./presentations-repos
                   w2 -> gateway-node1:~/work/presentations-repos
    ```

---

### 6.4 Results Telemetry Table Specification

- **Columns & Layout:**
  - `NODE (ALIAS)`: 16 chars (`local (current)` for local machine, alias for remote nodes).
  - `HOST`: 18 chars (`127.0.0.1` for local, IP address for remote).
  - `ROLE`: 10 chars (`master` for local, `worker` for remote).
  - `STATUS`: 14 visual chars (`● success`, `○ skipped`, `✗ failed`, `▲ auth_failed`).
  - `VERSION`: 12 chars (local version from `constants.Version`, remote version from `nodeVersionCache`).
  - `DURATION`: 12 chars (`%dms` formatted duration).
  - `DETAILS`: Truncated, sanitized status string ($\le 55$ chars).

---

### 6.5 Dynamic Suggestions Footer Specification

Rendered at the termination of every execution to provide guided next steps:

```go
func renderFleetFooterSuggestions(out io.Writer, opts NodesCloneOptions, hasFailures bool) {
	target := "repo"
	if len(opts.PassArgs) > 0 {
		target = opts.PassArgs[0]
	}
	fmt.Fprintln(out, "  [tip] Fleet Operations & Suggested Commands:")
	fmt.Fprintln(out, "    • Ping Fleet Nodes:          gitmap nodes ping")
	fmt.Fprintln(out, "    • Query Machine Telemetry:   gitmap machine --ssh")
	if hasFailures {
		fmt.Fprintln(out, "    • Inspect Node Connection:   gitmap ssh test <alias>")
		fmt.Fprintln(out, "    • Accept Remote Host Keys:   gitmap ssh trust <alias>")
	} else {
		fmt.Fprintln(out, "    • Inspect Node Connection:   gitmap ssh test <alias>")
	}
	fmt.Fprintf(out, "    • Rerun Except Local Host:   gitmap nodes %s %s --except-self\n\n", opts.Kind, target)
}
```

---

### 6.6 ANSI Escape Handling & Width Computation Rules

Because ANSI color escape sequences occupy bytes but zero terminal columns, direct string slicing and length calculations (`len()`) cause severe table column misalignments.
- `stripANSI(s string) string`: Strips all ANSI CSI escape sequences matching `\x1b[...[mKHJ]`.
- `visualLen(s string) int`: Computes true visual rune width after stripping ANSI escapes.
- `padVisual(s string, width int) string`: Appends `width - visualLen(s)` space characters to ensure exact column alignment.

---

## 7. Help System & Documentation Synchronization

### 7.1 Enhanced Interactive Help Screen (`cli/cmdnodes/nodes_clone_help.go`)

`PrintNodesCloneHelp` must render structured, comprehensive guidance across five distinct sections:
1. **Header Box:** Displays command title and operation mode.
2. **Usage Section:**
   ```text
   Usage:
     gitmap nodes <kind> [flags] <repo|url|file> [dest]
     gitmap nodes <kind> except-self <repo|url|file> [dest]
     gitmap nodes <kind> [flags] <repo1,repo2,...> [dest]
   ```
3. **Description Section:** Explicitly documents:
   - Subfolder Mirroring behavior when invoked inside work directory subfolders.
   - Target Directory Pre-Creation (`mkdir -p` and `New-Item -ItemType Directory`).
   - Custom Destination Overrides via flags or positional arguments.
   - Except-Self Mode skipping local execution.
   - Manifest Staging copying `.json` files to remote workers before cloning.
4. **Parameters & Flags Section:** Fully documents `-t/--target`, `--exclude`, `-d/--dest/--dir/--target-dir`, `--except-self/--no-self/--skip-local`, `--dry-run`, `-j/--json`, `-h/--help`.
5. **Rich Concrete Examples:**
   - Clone repo by name across fleet: `gitmap nodes cfr ChrisTitusTech/winutil`
   - Clone URL into custom target directory: `gitmap nodes clone https://github.com/u/repo ./custom`
   - Clone strictly across remote nodes: `gitmap nodes cfr except-self ChrisTitusTech/winutil`
   - Query machine network and SSH telemetry: `gitmap machine --ssh`
   - Probe fleet reachability and latency: `gitmap nodes ping`
   - Accept remote host keys and enroll: `gitmap ssh trust <alias>`
   - Staged JSON manifest clone: `gitmap nodes cfr gitmap.json`

---

### 7.2 Help Catalog Registration (`cli/helptext/catalog.go`)

To enable `gitmap help <topic>` lookup for fleet clone commands, `cli/helptext/catalog.go` is updated with corresponding entries in `topicSummaries`:

```go
var topicSummaries = map[string]string{
    // ... existing entries ...
    "nodes-clone":        "Asynchronously clone repositories or manifests across local host and remote SSH fleet nodes with subfolder preservation.",
    "nodes-cfr":          "Asynchronously clone, remediate, and auto-setup repositories across local host and remote SSH fleet nodes.",
    "nodes-cfrp":         "Asynchronously clone, remediate, and promote repository visibility to public across local host and remote fleet nodes.",
    "cfr":                "Asynchronously clone, remediate, and auto-setup repositories across local host and remote SSH fleet nodes.",
    "cfrp":               "Asynchronously clone, remediate, and promote repository visibility to public across local host and remote fleet nodes.",
    "clone-except-self":  "Execute repository cloning strictly across remote SSH fleet nodes, skipping the local master workstation.",
}
```

---

## 8. Quality Gates, Linters & Release Governance

### 8.1 Architectural & Coding Rules
1. **15-Line Function Limit:** No newly authored or modified function may exceed 15 lines of executable statements. Complex rendering or parsing routines must be decomposed into atomic helper subroutines.
2. **Positive Boolean Naming:** All booleans must use positive prefixes (`isOnline`, `hasCustomTargetDir`, `isInsideWorkDir`, `isLocalSuccess`).
3. **Structured Error Wrapping:** Errors must be handled cleanly or wrapped with contextual descriptions (`apperror.WrapSimple` or `fmt.Errorf`).
4. **Absolute Path Ban:** No hardcoded absolute paths outside user-configurable defaults. All internal references must use clean relative paths or dynamic discovery.
5. **Zero Raw Git Command Ban:** Neither subagents nor execution scripts may invoke `git add`, `git commit`, `git push`, etc. All version control operations must proceed through `gitmap cpf`.

### 8.2 Linter Remediation Suite
Before marking subtasks complete, the following linters must pass with zero violations:
```powershell
python 03-ai-scripts/05-guideline-autofixer.py --files cli/cmdnodes/nodes_clone_types.go cli/cmdnodes/nodes_clone.go cli/cmdnodes/nodes_clone_file.go cli/cmdnodes/nodes_clone_remote.go cli/cmdnodes/nodes_clone_table.go cli/cmdnodes/nodes_clone_help.go cli/helptext/catalog.go
python 03-ai-scripts/11-check-relative-paths.py
python 03-ai-scripts/21-check-forbidden-strings.py
```

### 8.3 Release Ceremony
- **Target Bump:** Minor SemVer increment (`v6.454.0`).
- **Automation Script:** `python 03-ai-scripts/37-bump-version.py --tier minor`.
- **Commit Headline:** `feat(nodes): overhaul cfr terminal ui, preflight probing and async delegation`.
