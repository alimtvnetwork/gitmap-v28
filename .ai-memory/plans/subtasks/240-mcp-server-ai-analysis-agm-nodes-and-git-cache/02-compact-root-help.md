# Subtask 02: Minimal Root Help Invariant & Compact CLI UX

- **Parent Task:** `240-mcp-server-ai-analysis-agm-nodes-and-git-cache`
- **Subtask ID:** `02-compact-root-help`
- **Spec Reference:** `02-spec/21-app/240-mcp-server-ai-analysis-agm-nodes-and-git-cache/01-architecture-spec.md` (Section 2.4 & Section 7)
- **Status:** `PENDING`
- **Assigned Subagent:** Implementation Subagent (Phase 2)

---

## 1. Objective & Scope

Implement the **Minimal Root Help Invariant** across the GitMap command-line interface.

### The Problem
Currently, executing bare `gitmap` with no positional arguments (`len(os.Args) < 2`) triggers `printUsage()`, which outputs hundreds of lines of categorized help menus, subcommands, flags, and trailers. For human developers and automated AI tools inspecting the binary environment, this massive dump creates terminal noise, disrupts CLI screen flow, and inflates LLM token costs.

### The Solution
Enforce the Minimal Root Help Invariant:
1. **Bare `gitmap` Invocation:** When executed with no arguments, GitMap strictly outputs a compact, elegant 3-part card:
   - **Version Block:** Compact brand and SemVer version identifier.
   - **Location Triplet:** The three core runtime filesystem paths: Executable Binary, Config File, and Database/Data Directory (emitted via `PrintBinaryLocations()`).
   - **Concise Footer:** Directs the user to `gitmap help` or `gitmap -h` for the full catalog, providing two quick-start commands (`gitmap scan`, `gitmap clone <url>`).
2. **Full Catalog Reservation:** The rich, categorized command catalog (Get Started, Repos, Release, Projects, Advanced) remains fully accessible when the user explicitly requests help via `gitmap help`, `gitmap -h`, or `gitmap --help`.

---

## 2. Concrete Files to Create / Modify

| File | Nature | Purpose |
| :--- | :--- | :--- |
| `cli/cmd/root.go` | Modify | Update the bare invocation checks (`len(os.Args) < 2`) to call `printCompactRootSummary()` instead of `printUsage()`. |
| `cli/cmd/rootusagecompact.go` | Modify / Enhance | Implement `printCompactRootSummary()` rendering the version banner, location triplet, and quick guidance footer. |
| `cli/cmd/root_no_args_test.go` | Modify / Extend | Add regression tests verifying bare `gitmap` output line count (<15 lines) and presence of the location triplet. |

---

## 3. Detailed Implementation Requirements

### 3.1 `printCompactRootSummary()` Design
```go
// printCompactRootSummary outputs the concise version, location triplet,
// and quick help hints for bare `gitmap` invocations.
func printCompactRootSummary() {
	fmt.Printf("GitMap v%s - Autonomous Developer Companion & AI MCP Server\n\n", constants.Version)
	PrintBinaryLocations()
	fmt.Println()
	fmt.Println("Quick Start:")
	fmt.Println("  gitmap scan                  Discover and index all repositories")
	fmt.Println("  gitmap clone <url>           Clone and register a git repository")
	fmt.Println("  gitmap ai-analysis           Inspect autonomous AI tasks & staged files")
	fmt.Println()
	fmt.Println("For complete command catalog and options:")
	fmt.Println("  Run 'gitmap help' or 'gitmap -h'")
}
```

### 3.2 Dispatch Logic in `cli/cmd/root.go`
```go
if len(os.Args) < 2 {
    printCompactRootSummary()
    return
}
```
Ensure that early exit occurs cleanly without panicking, initializing unnecessary background routines, or executing heavy migrations.

---

## 4. Acceptance Criteria

- [ ] Invoking `gitmap` with no arguments produces less than 15 lines of output.
- [ ] Output strictly contains:
  1. The version banner.
  2. The location triplet (Binary, Config, Data).
  3. The quick start tips and help command hint (`gitmap help` / `-h`).
- [ ] Invoking `gitmap help`, `gitmap -h`, or `gitmap --help` continues to display the full categorized command catalog without regression.
- [ ] Automated unit test `TestBareRootOutputCompact` passes.
- [ ] Existing test `TestRunNoArgsDoesNotPanic` continues to pass.

---

## 5. Verification Commands

```powershell
# 1. Run unit tests for root command no-args behavior
go test -v ./cli/cmd -run "TestRunNoArgs|TestBareRootOutput"

# 2. Build and verify bare output line count manually
go build -o gitmap.exe .
$output = .\gitmap.exe
$output
($output -split "`n").Count # Must be <= 15 lines
```
