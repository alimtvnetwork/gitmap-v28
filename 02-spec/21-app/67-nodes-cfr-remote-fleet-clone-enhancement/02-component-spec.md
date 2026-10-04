# Component Technical Specification: Work Directory Preservation, Destination Overrides, Unified Help & Actionable Suggestions

- **Spec Document:** `02-spec/21-app/67-nodes-cfr-remote-fleet-clone-enhancement/02-component-spec.md`
- **Spec ID:** 67-02
- **Status:** Approved / Ready for Implementation
- **Author:** Spec Author 02 (Component Spec & Tasks 04–06)
- **Parent Plan:** [.ai-memory/plans/pending/67-nodes-cfr-remote-fleet-clone-enhancement.md](../../../.ai-memory/plans/pending/67-nodes-cfr-remote-fleet-clone-enhancement.md)
- **Target Subtasks:**
  - [Subtask 67.04: Intelligent Work Directory & Relative Path Resolution](../../../.ai-memory/plans/subtasks/67-nodes-cfr-remote-fleet-clone-enhancement/04-work-dir-and-relative-path.md)
  - [Subtask 67.05: Nodes Help Menu, Examples & Command Suggestions Alignment](../../../.ai-memory/plans/subtasks/67-nodes-cfr-remote-fleet-clone-enhancement/05-nodes-help-and-suggestions.md)
  - [Subtask 67.06: Verification, Minor Bump & Release Ceremony](../../../.ai-memory/plans/subtasks/67-nodes-cfr-remote-fleet-clone-enhancement/06-verification-and-release.md)

---

## 1. Executive Summary & Problem Analysis

### 1.1 Context & Incident Description
GitMap's fleet orchestration command `gitmap nodes cfr` (and related commands `gitmap nodes clone`, `gitmap nodes cfrp`) coordinates repository cloning and remediation across a local machine and remote SSH fleet nodes.

In practical everyday operations, engineers encountered several major usability, directory placement, and discoverability deficiencies:

1. **Blind Remote Directory Placement:**
   Currently, remote nodes default unconditionally to root work folders (`D:\work` on Windows, `~/work` on Linux). When an engineer is operating inside a specialized workspace subfolder locally (for example, `./presentations-repos`), executing `gitmap nodes cfr` places repositories into `D:\work` or `~/work` on remote machines instead of mirroring the local subfolder hierarchy (`./presentations-repos` and `~/work/presentations-repos`).
2. **Missing Destination Overrides:**
   Users lacked flexible CLI mechanisms to specify an explicit target directory. There was no support for standard destination flags (`-d`, `--dir`, `--dest`, `--target-dir`) or an intuitive optional positional destination argument (`gitmap nodes cfr <target> [dest]`).
3. **Remote Execution Failure on Missing Subdirectories:**
   Remote shell command generation constructed `Set-Location "<dir>"` (PowerShell) or `cd <dir>` (Bash) without ensuring the destination directory actually existed on remote nodes beforehand. If a subfolder did not already exist, the remote command failed immediately with path not found errors.
4. **Undocumented CLI Subsystem & Missing Help Pages:**
   While `gitmap help` indexed `nodes` in `cli/helptext/catalog.go`, no corresponding `cli/helptext/nodes.md` markdown file existed. Running `gitmap help nodes` or accessing embedded documentation yielded missing help errors or partial fallback output. Furthermore, `printUnifiedNodesHelp` in `cli/cmd/nodes_cmd.go` did not explain `cfr` destination mechanics, accept keys, or remote machine inspection.
5. **Absence of Post-Execution Next Steps:**
   Upon completing fleet clone operations, the terminal displayed a basic summary count but provided zero actionable suggestions for next steps (such as querying machine health, connecting via SSH, or inspecting clone history).

### 1.2 Architectural Scope & Component Responsibilities
This component specification defines the architectural contracts, data structures, algorithms, and CLI interfaces for:
- **Intelligent Work Directory & Relative Path Preservation** in `cli/cmdnodes/nodes_clone_file.go`, `cli/cmdnodes/nodes_clone.go`, and `cli/cmdnodes/nodes_clone_types.go`.
- **Remote Pre-Creation of Target Directories** across Windows (`New-Item -ItemType Directory -Force`) and Linux (`mkdir -p`) nodes.
- **Destination Flag and Positional Argument Resolution** with collision-free target parsing.
- **Embedded Help Architecture** in `cli/helptext/nodes.md` and terminal help parity in `cli/cmd/nodes_cmd.go`.
- **Actionable Footer Suggestions** in `cli/cmdnodes/nodes_clone_table.go`.
- **Quality Gates, Guideline Compliance, and SemVer Minor Bump** (`v6.454.0`).

---

## 2. Component Architecture & Data Contracts

### 2.1 File Distribution & Module Boundaries
```
cli/
├── cmdnodes/
│   ├── nodes_clone_types.go     # Options, kinds, results & configuration types
│   ├── nodes_clone_file.go      # CWD relative path resolution & remote path construction
│   ├── nodes_clone.go           # CLI flag parsing (-d, --dir, --dest, --target-dir)
│   ├── nodes_clone_remote.go    # Remote command string building with mkdir-p / New-Item
│   ├── nodes_clone_table.go     # Terminal UI table & actionable footer suggestions
│   └── nodes_clone_help.go      # Interactive help screens for nodes clone/cfr/cfrp
├── cmd/
│   └── nodes_cmd.go             # Unified nodes CLI entrypoint & printUnifiedNodesHelp
└── helptext/
    └── nodes.md                 # Embedded markdown documentation for gitmap help nodes
```

### 2.2 Data Contracts: `cli/cmdnodes/nodes_clone_types.go`

`NodesCloneOptions` represents the normalized configuration for fleet clone execution. The struct is extended with fields tracking work directory status, relative subfolder paths, and explicit destination overrides:

```go
package cmdnodes

import (
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdclone"
)

// NodesCloneKind defines the clone operation mode.
type NodesCloneKind string

const (
	// CloneKindClone executes standard multi-repo or single-repo clone.
	CloneKindClone NodesCloneKind = "clone"
	// CloneKindCFR executes clone-fix-repo (remediation & auto-setup).
	CloneKindCFR NodesCloneKind = "cfr"
	// CloneKindCFRP executes clone-fix-repo with public visibility promotion.
	CloneKindCFRP NodesCloneKind = "cfrp"
)

// NodesCloneOptions stores configuration for fleet clone execution.
type NodesCloneOptions struct {
	Kind             NodesCloneKind `json:"kind"`
	TargetFilter     string         `json:"targetFilter,omitempty"`
	ExcludeFilter    string         `json:"excludeFilter,omitempty"`
	ExceptOS         string         `json:"exceptOS,omitempty"`
	TargetOS         string         `json:"targetOS,omitempty"`
	TargetDir        string         `json:"targetDir,omitempty"`        // Explicit or resolved target dir
	RelativeSubDir   string         `json:"relativeSubDir,omitempty"`  // Relative subpath inside workdir
	IsInsideWorkDir  bool           `json:"isInsideWorkDir"`           // True if CWD is inside work dir
	HasCustomDest    bool           `json:"hasCustomDest"`             // True if user specified -d/--dir/dest
	IsDryRun         bool           `json:"isDryRun"`
	IsJSON           bool           `json:"isJSON"`
	IsSkipLocal      bool           `json:"isSkipLocal"`
	RawArgs          []string       `json:"rawArgs"`
	PassArgs         []string       `json:"passArgs"`
	DetectedFile     string         `json:"detectedFile,omitempty"`
	HasFile          bool           `json:"hasFile"`
}

// RemoteCloneNodeResult captures the execution output from an individual node.
type RemoteCloneNodeResult struct {
	Alias      string        `json:"alias"`
	Host       string        `json:"host"`
	Role       string        `json:"role"`
	Status     string        `json:"status"`
	Duration   time.Duration `json:"duration"`
	DurationMs int64         `json:"durationMs"`
	Stdout     string        `json:"stdout,omitempty"`
	Stderr     string        `json:"stderr,omitempty"`
	Error      string        `json:"error,omitempty"`
	Details    string        `json:"details,omitempty"`
	TargetDir  string        `json:"targetDir,omitempty"`
}
```

### 2.3 Field Specifications & Invariants
- `IsInsideWorkDir`: Must be set to `true` if and only if the current working directory is equal to or a subfolder of the local configured/default work directory (`fsutil.IsInsideWorkDir(cwd, localWorkDir)`).
- `RelativeSubDir`: Represents the clean, relative subfolder path from the local work directory to the current working directory (e.g. `presentations-repos` or `sub/nested`). If CWD is the work directory root or outside the work directory, `RelativeSubDir` is empty `""`.
- `TargetDir`:
  - When `HasCustomDest` is true: contains the raw or normalized custom directory requested by the user.
  - When `HasCustomDest` is false and `IsInsideWorkDir` is true: contains the relative subfolder or local directory path for local execution.
  - When `HasCustomDest` is false and `IsInsideWorkDir` is false: defaults to empty `""` (meaning default work root).
- `HasCustomDest`: Must be `true` if the user provided `-d`, `--dir`, `--dest`, `--target-dir`, or an explicit positional destination argument.

---

## 3. Intelligent Work Directory & Relative Path Resolution

### 3.1 Local Work Directory Discovery
To detect whether the user is inside a registered work directory, GitMap consults the database configuration first, falling back to standard platform defaults:

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
	return filepath.Join(os.Getenv("HOME"), "work")
}
```

### 3.2 CWD Relative Subdirectory Computation
When a command is invoked without an explicit destination argument, GitMap inspects the current process working directory:

```go
// ResolveCwdRelativeSubDir determines if CWD is inside the local work directory
// and computes the relative subpath if applicable.
func ResolveCwdRelativeSubDir(workDir string) (string, bool) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", false
	}
	cleanCwd := filepath.Clean(cwd)
	cleanWork := filepath.Clean(workDir)

	if !fsutil.IsInsideWorkDir(cleanCwd, cleanWork) {
		return "", false
	}

	rel, errRel := filepath.Rel(cleanWork, cleanCwd)
	if errRel != nil || rel == "." || strings.HasPrefix(rel, "..") {
		return "", true // Inside work directory root, no subfolder
	}

	// Normalize to forward slashes for cross-platform portability
	normalizedRel := filepath.ToSlash(rel)
	return normalizedRel, true
}
```

### 3.3 Remote Target Directory Construction
For each target remote node, GitMap resolves the appropriate root work directory (`D:\work` for Windows nodes, `~/work` for Linux/Unix nodes) and appends the relative subfolder if preserved:

```go
// ResolveRemoteWorkDir returns the default fleet work directory based on node OS.
func ResolveRemoteWorkDir(osType string) string {
	isWin := strings.EqualFold(osType, "windows") || strings.EqualFold(osType, "win")
	if isWin {
		return `D:\work`
	}
	return "~/work"
}

// ResolveRemoteTargetDir constructs the target execution directory on a remote node.
func ResolveRemoteTargetDir(osType, relativeSubDir, customDest string) string {
	if customDest != "" {
		return formatRemoteCustomDest(osType, customDest)
	}
	baseDir := ResolveRemoteWorkDir(osType)
	if relativeSubDir == "" {
		return baseDir
	}
	isWin := strings.EqualFold(osType, "windows") || strings.EqualFold(osType, "win")
	if isWin {
		cleanSub := strings.ReplaceAll(relativeSubDir, "/", "\\")
		return baseDir + "\\" + cleanSub
	}
	cleanSub := strings.Trim(relativeSubDir, "/")
	return strings.TrimRight(baseDir, "/") + "/" + cleanSub
}
```

### 3.4 Resolution Decision Matrix
| Execution Scenario | Local CWD | Local WorkDir | User Flag / Dest | Remote Windows Node | Remote Linux Node |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Inside WorkDir Subfolder** | `./presentations-repos` | `D:\work` | *(none)* | `./presentations-repos` | `~/work/presentations-repos` |
| **Inside Nested Subfolder** | `./clients\alpha` | `D:\work` | *(none)* | `./clients\alpha` | `~/work/clients/alpha` |
| **Inside WorkDir Root** | `D:\work` | `D:\work` | *(none)* | `D:\work` | `~/work` |
| **Outside WorkDir** | `C:\Users\Admin` | `D:\work` | *(none)* | `D:\work` | `~/work` |
| **Outside WorkDir** | `/home/user/downloads` | `~/work` | *(none)* | `D:\work` | `~/work` |
| **Explicit Flag Override** | `./presentations-repos` | `D:\work` | `-d ./custom` | `./custom` | `~/work/custom` |
| **Explicit Positional Override** | `C:\Temp` | `D:\work` | `./special` | `./special` | `~/work/special` |

### 3.5 Pre-Creation of Remote Target Directories
Prior to navigating or executing `gitmap <kind>` on remote nodes, the target directory must exist. Shell execution strings are constructed to safely ensure the directory exists:

#### Windows PowerShell Command Pattern:
```powershell
if (!(Test-Path -Path '<targetDir>')) { New-Item -ItemType Directory -Path '<targetDir>' -Force | Out-Null }; Set-Location '<targetDir>'; gitmap <kind> <args> --json
```

#### Linux / Unix Bash Command Pattern:
```bash
mkdir -p <targetDir> && cd <targetDir> && gitmap <kind> <args> --json
```

#### Helper Implementation Contract in `cli/cmdnodes/nodes_clone_remote.go`:
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

## 4. CLI Flags & Destination Argument Parsing

### 4.1 Supported CLI Syntax
The command line supports multiple variations for specifying destination paths:

1. **Explicit Flags:**
   - `-d <path>`
   - `--dir <path>`
   - `--dest <path>`
   - `--target-dir <path>`
   - `--dir=<path>`, `--dest=<path>`, `--target-dir=<path>`
2. **Positional Arguments:**
   - `gitmap nodes cfr <target> <destination>`
   - `gitmap nodes clone <repo-url> <destination>`
   - `gitmap nodes cfrp <manifest.json> <destination>`
   - `gitmap nodes cfr except-self <target> <destination>`

### 4.2 Distinguishing Targets from Destinations
Arguments must be classified deterministically:
- An argument is recognized as a **Target** if:
  - It matches a known file on disk (e.g. `gitmap.json`, `repos.json`).
  - It is a URL starting with `http://`, `https://`, or `git@`.
  - It contains a comma separating multiple items (`repo1,repo2`).
  - It is the first non-flag positional argument in the target position.
- An argument is recognized as an **Explicit Destination** if:
  - It follows an explicit flag (`-d`, `--dir`, `--dest`, `--target-dir`).
  - Or it is a subsequent non-flag positional argument that does not start with `-`, is not a URL, and does not contain commas.

### 4.3 Parsing Logic in `parseNodesCloneOptions`
```go
func isDestFlag(arg string) bool {
	return arg == "-d" || arg == "--dir" || arg == "--dest" || arg == "--target-dir"
}

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
```

---

## 5. Help Documentation Structure & Actionable Suggestions

### 5.1 Embedded Help Specification: `cli/helptext/nodes.md`
To ensure `gitmap help nodes` functions identically to existing help topics (`ssh.md`, `cluster.md`, `scan.md`), `cli/helptext/nodes.md` must be created with the canonical layout:

```markdown
# Unified Fleet & Infrastructure Nodes

Aggregates, discovers, monitors, and executes commands across all registered machines in the GitMap fleet.

## Usage

    gitmap nodes [subcommand|filter] [flags]
    gitmap nodes clone [flags] <targets> [dest]
    gitmap nodes cfr [flags] [targets] [dest]
    gitmap nodes cfrp [flags] [targets] [dest]
    gitmap nodes ping [target]

## Subcommands

| Subcommand | Alias | Description |
| :--- | :--- | :--- |
| *(none)* | | Display all registered nodes across SSH, Cluster DB, and Server-Client networks |
| ping | probe | Ping and probe network latency and liveness across fleet nodes |
| clone | | Clone repository or manifest asynchronously locally and across fleet nodes |
| cfr | clone-fix-repo | Clone, fix repository, and run auto-setup locally and across fleet nodes |
| cfrp | clone-fix-repo-pub | Clone, fix repository, and promote visibility to public across fleet nodes |
| history | histories | Display past fleet execution audit history |
| help | -h, --help | Display this comprehensive help manual |

## Fleet Clone & Remediation Flags

| Flag | Short | Description |
| :--- | :--- | :--- |
| --target | -t | Dispatch only to a specific node alias or IP |
| --exclude | | Exclude specific nodes by alias or IP |
| --except-self | --no-self | Execute purely on remote fleet nodes, skipping the local master host |
| --dest | -d, --dir, --target-dir | Explicit destination directory for cloned repositories |
| --dry-run | | Preview operations and targets without modifying filesystems |
| --json | -j | Emit machine-readable JSON telemetry |

## Work Directory & Path Resolution

GitMap intelligently detects your current working directory:
1. **Inside Work Directory Subfolder:** If executed from `./presentations-repos`, repositories are cloned into `<remote-work-dir>/presentations-repos` on all remote nodes.
2. **Outside Work Directory:** If executed outside a work directory, repositories default to remote work roots (`D:\work` on Windows, `~/work` on Linux).
3. **Explicit Override:** Specifying `[dest]` or `--dest <path>` directs clones to that explicit directory across all nodes.
4. **Auto-Directory Creation:** Target folders are created automatically (`mkdir -p` / `New-Item -ItemType Directory`) prior to cloning.

## Examples

    # Clone a repository across all fleet nodes
    gitmap nodes cfr ChrisTitusTech/winutil

    # Clone into a custom target directory
    gitmap nodes cfr https://github.com/user/project.git ./custom

    # Clone strictly across remote nodes, preserving relative subfolder
    gitmap nodes cfr except-self user/project

    # Transfer and clone a local manifest across fleet nodes
    gitmap nodes cfr gitmap.json

    # Probe fleet reachability and GitMap binary versions
    gitmap nodes ping
    gitmap machine
```

### 5.2 Terminal Help Synchronization: `cli/cmd/nodes_cmd.go`
`printUnifiedNodesHelp()` must be synchronized with `cli/helptext/nodes.md`, including:
- Clear documentation for `[dest]` in `gitmap nodes clone`, `cfr`, and `cfrp`.
- Subfolder mirroring explanations.
- Examples showing `gitmap machine` and accept keys.

### 5.3 Actionable Footer Suggestions: `cli/cmdnodes/nodes_clone_table.go`
Upon completing `gitmap nodes clone/cfr/cfrp`, the terminal table summary must output context-sensitive next-step recommendations:

```go
func renderFleetSuggestionsFooter(out io.Writer, opts NodesCloneOptions, results []RemoteCloneNodeResult) {
	fmt.Fprintln(out, "  💡 Suggestions & Next Steps:")
	hasFailures := false
	for _, r := range results {
		if r.Status != "success" {
			hasFailures = true
			break
		}
	}
	if hasFailures {
		fmt.Fprintln(out, "    • Probe fleet reachability:       gitmap nodes ping")
		fmt.Fprintln(out, "    • Inspect node connection logs:   gitmap ssh check")
		fmt.Fprintln(out, "    • Accept host keys / re-enroll:   gitmap ssh trust <alias>")
	} else {
		fmt.Fprintln(out, "    • Inspect remote node status:     gitmap nodes")
		fmt.Fprintln(out, "    • Query local machine identity:   gitmap machine")
	}
	if len(results) > 0 {
		fmt.Fprintf(out, "    • Connect to remote worker node:  gitmap ssh %s\n", results[0].Alias)
	}
	fmt.Fprintln(out, "    • Review fleet execution history: gitmap nodes history")
	fmt.Fprintln(out)
}
```

---

## 6. Verification, Autofixing & Release Governance

### 6.1 Coding Guidelines & Rule Enforcement
Every file authored or modified under this specification must strictly comply with:
1. **Positive Boolean Naming:** Use `is*` and `has*` prefixes (`isInsideWorkDir`, `hasCustomDest`, `isLocalSuccess`).
2. **Function Line Limits:** Every function must be $\le 15$ lines of executable code.
3. **Structured Error Wrapping:** Wrap all errors using `apperror.WrapSimple(err, "<caller>")` or `apperror.WrapWithDetails()`.
4. **No Git Commands:** Subagents and execution agents are strictly forbidden from running raw `git` commands. All commits must be made via `gitmap cpf`.

### 6.2 Guideline Linter Validation
Before release, run the project's canonical guideline autofixer:
```powershell
python 03-ai-scripts/05-guideline-autofixer.py --files cli/cmdnodes/nodes_clone_types.go cli/cmdnodes/nodes_clone_file.go cli/cmdnodes/nodes_clone.go cli/cmdnodes/nodes_clone_remote.go cli/cmdnodes/nodes_clone_table.go cli/cmd/nodes_cmd.go
```

### 6.3 Version Bump & Atomic Release
- **Bump Target:** Minor SemVer release `v6.454.0` in `version.json`.
- **Root README Update:** Synchronize `readme.md` version badges and release highlights.
- **Atomic Commit Message:** `gitmap cpf "nodes - modernize cfr ui remote version probe and async dispatch"`.
