# Subtask 05: Documentation, Help, Completion & UI Registry

- **Subtask ID:** `SUBTASK-05`
- **Parent Plan:** `.ai-memory/plans/gitmap-nodes-running-only.md`
- **Component Spec:** `02-spec/21-app/gitmap-nodes-running-only/02-component-spec.md`
- **Status:** `PENDING`
- **Priority:** Normal
- **Dependencies:** Subtask 02 (`02-gitmap-nodes-running-only-nro-dispatcher-and-cli.md`), Subtask 03 (`03-running-only-delegation-and-fast-ssh-execution.md`)

---

## 1. Objective & Context

The introduction of `gitmap nodes-running-only` and its canonical alias `gitmap nro` requires full integration across documentation engines, terminal help renderers, shell completion generators, and the web UI command catalog.

This subtask ensures complete discoverability and developer ergonomics by:
1. Authoring the comprehensive Markdown help document in `cli/helpdoc/nodes-running-only.md`.
2. Registering command aliases in `helpAliases` in `cli/helpdoc/print.go`.
3. Registering command catalog topic summaries in `cli/helpdoc/catalog.go`.
4. Registering command definitions in the TypeScript web UI catalog (`src/data/commands.ts`).
5. Registering shell tab completion candidates in `cli/completion/completion.go`.
6. Declaring standardized constants in `cli/constants/constants_commands.go`.

---

## 2. Owned Files & Target Modifications

| File Path | Description of Changes |
|---|---|
| `cli/helpdoc/nodes-running-only.md` | Author detailed Markdown help document covering synopsis, flags, subcommands, output formats, and practical examples. |
| `cli/helpdoc/print.go` | Register `nro`, `nodes-running-only`, `nodesrunningonly`, `nodes-running`, `nodes-online`, `nodes-open`, `nodes-active` in `helpAliases`. |
| `cli/helpdoc/catalog.go` | Register summary entries for `nodes-running-only` and `nro` in `topicSummaries`. |
| `src/data/commands.ts` | Add `nodes-running-only` command definition object with category `cluster`, flags, examples, and seeAlso links. |
| `cli/completion/completion.go` | Add `"nodes-running-only"`, `"nro"`, `"nodesrunningonly"`, `"nodes-running"` to `manualExtras` slice. |
| `cli/constants/constants_commands.go` | Declare `CmdNodesRunningOnly` and `CmdNodesRunningOnlyAlias`. |

---

## 3. Detailed Technical Design & Code Specifications

### 3.1 Constants Declaration (`cli/constants/constants_commands.go`)

Add constants for the new command and alias:

```go
	CmdNodesRunningOnly      = "nodes-running-only"
	CmdNodesRunningOnlyAlias = "nro"
```

### 3.2 Shell Tab Completion Registration (`cli/completion/completion.go`)

Add the command tokens to `manualExtras`:

```go
var manualExtras = []string{
	"db", "start-fresh", "find-duplicates", "pe", "pd", "ping",
	"nodes", "node", "gitignore", "agm", "ta",
	"cpu", "processor", "ram", "disk",
	"nodes-running-only", "nro", "nodesrunningonly", "nodes-running",
}
```

### 3.3 Terminal Help Alias Mapping (`cli/helpdoc/print.go`)

In `var helpAliases`:

```go
	"nodes-running-only":          "nodes-running-only",
	"nro":                         "nodes-running-only",
	"nodesrunningonly":            "nodes-running-only",
	"nodes-running":               "nodes-running-only",
	"nodes-online":                "nodes-running-only",
	"nodes-open":                  "nodes-running-only",
	"nodes-active":                "nodes-running-only",
```

### 3.4 Topic Summaries in Catalog (`cli/helpdoc/catalog.go`)

In `var topicSummaries`:

```go
	"nodes-running-only": "Discover and execute commands strictly on online cluster nodes with 10s liveness caching and zero offline delay.",
	"nro":                "Discover and execute commands strictly on online cluster nodes with 10s liveness caching and zero offline delay.",
```

### 3.5 Web UI Catalog Definition (`src/data/commands.ts`)

Insert after `gitmap nodes tasks` in `Commands` array:

```typescript
  {
    name: "gitmap nodes-running-only",
    alias: "gitmap nro, nodes-running, nodes-online, nodes-open",
    description: "Discover and execute commands strictly on online fleet nodes, skipping offline machines with 10s liveness caching and zero timeout latency",
    usage: "gitmap nodes-running-only [--json] [--refresh] [exec|tasks|resources|ping|pull|cp|scan]",
    category: "cluster",
    flags: [
      { flag: "--json, -j", description: "Output active nodes as a clean JSON array with is_online: true guaranteed" },
      { flag: "--refresh, -r", description: "Bypass 10s liveness cache and trigger immediate network reachability probe" },
      { flag: "--target, -t <alias>", description: "Target a specific node alias or IP address if online" },
      { flag: "--timeout <duration>", description: "Custom timeout limit for probe or remote execution" },
    ],
    examples: [
      { command: "gitmap nro", description: "Display ANSI table of currently active and reachable cluster nodes" },
      { command: "gitmap nro --json", description: "Emit machine-readable JSON array of online nodes for automated scripts" },
      { command: "gitmap nro exec \"uptime\"", description: "Execute remote command across active nodes only without offline delays" },
      { command: "gitmap nro tasks", description: "Query active AI agent task queues exclusively on online machines" },
      { command: "gitmap nro resources", description: "Probe CPU and RAM utilization exclusively on reachable nodes" },
    ],
    notes: [
      "Results are cached in OS temporary directory for 10 seconds to prevent repetitive TCP probe delays during sequential CLI calls.",
      "Delegation subcommands bypass offline node connection attempts completely, preventing warning notices and connection hang.",
    ],
    seeAlso: [
      { name: "nodes", description: "Unified fleet node discovery and management" },
      { name: "nodes resources", description: "Probe CPU and RAM telemetry across fleet" },
      { name: "gitmap nodes tasks", description: "Inspect active and pending AI agent tasks across fleet" },
    ],
  },
```

### 3.6 Help Document Authoring (`cli/helpdoc/nodes-running-only.md`)

Create `cli/helpdoc/nodes-running-only.md` containing:
1. Executive Synopsis and command syntax.
2. Complete table of delegation subcommands (`exec`, `tasks`, `resources`, `ping`, `pull`, `cp`, `scan`).
3. Flags list (`--json`, `-j`, `--refresh`, `-r`, `--target`, `-t`, `--timeout`).
4. Output examples: ANSI table format and JSON format (`is_online: true` only).
5. Automation and scripting workflows.

---

## 4. Step-by-Step Implementation Plan

1. **Step 1: Update `cli/constants/constants_commands.go`**:
   - Define `CmdNodesRunningOnly` and `CmdNodesRunningOnlyAlias`.
2. **Step 2: Update `cli/completion/completion.go`**:
   - Append `nodes-running-only`, `nro`, `nodesrunningonly`, and `nodes-running` to `manualExtras`.
3. **Step 3: Update `cli/helpdoc/print.go` and `catalog.go`**:
   - Add alias entries to `helpAliases` in `print.go`.
   - Add summary entries to `topicSummaries` in `catalog.go`.
4. **Step 4: Author `cli/helpdoc/nodes-running-only.md`**:
   - Write comprehensive Markdown document following existing helpdoc conventions.
5. **Step 5: Update `src/data/commands.ts`**:
   - Insert `gitmap nodes-running-only` command definition in the exported `Commands` list.

---

## 5. Invariants & Verification Checklist

- [ ] **Relative Paths**: Zero absolute paths or `file:///` URLs.
- [ ] **Catalog Parity**: Help aliases match the dispatch aliases in `cli/cmdnodes/nodes_running_only.go` and `cli/cmd/rootcore.go`.
- [ ] **No Build/Test Execution**: Document and author without intermediate test/build runs.
- [ ] **Positive Booleans**: Use positive prefix identifiers throughout code additions.
