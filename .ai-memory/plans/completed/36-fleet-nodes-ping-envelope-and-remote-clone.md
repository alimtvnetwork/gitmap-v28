# Milestone 36: Fleet Nodes Ping, Typed JSON Envelopes, and Remote Clone

- **Slug:** `fleet-nodes-ping-envelope-and-remote-clone`
- **Milestone Index:** `36`
- **Status:** `COMPLETED`
- **Source Plans Merged:** Plans 183, 184, 185, 186, 188, 191, 192, 194, 196
- **Folded Subtask Folders:**
  - `193-typed-json-envelope` (6 files)
  - `194-import-all-json-what-configs` (6 files)
  - `200-nodes-display-cfr-manifest-gitmap-json-and-vscode-optimize` (6 files)
  - `201-nodes-clone-cfr-cfrp-fleet-manifest-staging-and-suggestions` (4 files)
  - `67-nodes-cfr-remote-fleet-clone-enhancement` (6 files)
- **Target Subsystems:** `cli/cmdnodes/`, `cli/jsonenvelope/`, `cli/render/`, `cli/cmd/`

---

## 1. Domain Context & Architectural Problem

As GitMap expanded from single-machine repository orchestration to multi-node SSH fleet operations, four core architectural challenges emerged:
1. **Network Probing Limitations:** Basic ICMP echo requests failed or hung when traversing corporate networks where Windows Defender Firewall dropped ICMP packets while keeping SSH port 22 open. A dual-stack ICMP and TCP port probing mechanism was required to guarantee accurate machine reachability reporting (`● ONLINE`, `● REACHABLE (TCP)`, `○ OFFLINE`).
2. **Configuration & Data Ingestion Fragility:** Importing diverse configurations (SSH nodes, macro definitions, UI theme settings, commit profiles) required multiple fragmented CLI commands. A unified JSON Envelope specification (`which-format`, `import-all-json`, `what-configs` / `wc`) was needed to provide strongly typed schema validation and seamless batch loading.
3. **Cross-Node Fleet Cloning:** Remote repository cloning across distributed nodes suffered from self-targeting loops and misaligned terminal tables. Commands like `nodes clone --except-self`, Clone-From-Remote (CFR), and Clone-From-Remote-Parallel (CFRP) required deterministic node manifest generation and clean, double-bordered visual tables.
4. **Visual & Formatting Parity:** Terminal outputs suffered from raw unescaped JSON, trailing whitespace, and inconsistent box rendering when displaying remote execution results.

---

## 2. Synthesized Architectural Outcomes

### 2.1 Dual-Stack Fleet Reachability Probing (`gitmap nodes ping` / `gitmap ping`)
- **Engine Implementation:** Implemented in `cli/cmd/nodes_ping_cmd.go` and `cli/cmdnodes/ping.go`.
- **Cross-Platform Execution:** Automatically detects host platform and dispatches `ping -n <count>` on Windows or `ping -c <count>` on Linux/macOS.
- **Concurrent TCP Probing:** When ICMP ping encounters packet loss or firewall drops, the engine immediately initiates a concurrent TCP handshake (`net.DialTimeout`) to the configured SSH port (default 22).
- **Status Enum Classification:**
  - `StatusOnline`: Both ICMP echo and TCP handshake succeed.
  - `StatusReachableTCP`: ICMP echo dropped (firewall guard active), but TCP port 22 connects within timeout.
  - `StatusOffline`: Neither ICMP nor TCP handshake respond within timeout.
  - `StatusDegraded`: Packet loss observed (>0% and <100%) or RTT latency exceeds 350ms.
- **Formatted Box Rendering:** Renders a clean ANSI box table showing Node ID, Alias, IP/Hostname, Ping RTT (avg/min/max), TCP Status, and Machine Health indicator.

### 2.2 Typed JSON Envelope V2 & Unified Config Importer
- **Envelope Contracts:** Designed and implemented strongly typed schema envelopes in `cli/jsonenvelope/` supporting polymorphic payloads:
  - `EnvelopeCommitConfig`: Auto-stage, branch rules, templates DB configuration.
  - `EnvelopeMacroConfig`: Recorded macro steps, delay timings, target shells.
  - `EnvelopeSSHNodes`: Fleet node definitions, SSH credentials, jump hosts.
  - `EnvelopeUISettings`: ANSI colors, box borders, padding margins.
- **Format Inspection (`gitmap which-format` / `gitmap wf`):** Inspects arbitrary JSON files and identifies envelope schema, target version compatibility, and required fields.
- **Batch Config Importer (`gitmap import-all-json` / `gitmap wc`):**
  - Accepts glob patterns (`*.json`, `specific.json`) and safely loads valid configuration records.
  - Applies atomic transaction semantics via Split-DB connection pooling, rolling back on schema mismatches.
- **Terminal Clear Primitives:** Enforces clean terminal viewport clears (`termclear`) prior to interactive configuration displays.

### 2.3 Remote Fleet Cloning (CFR, CFRP & Except-Self)
- **Except-Self Guard:** Added `--except-self` flag to `gitmap nodes clone`, comparing local machine fingerprints, hostname, and primary network interfaces against target node descriptors to prevent local-to-local recursive cloning.
- **CFR (Clone From Remote):** Dispatches streaming SSH commands to clone target repositories on remote nodes using optimized shallow depth flags (`--depth 1`) when specified.
- **CFRP (Clone From Remote Parallel):** Worker pool distributing repository clone requests across multiple target fleet nodes simultaneously with per-node concurrency throttles.
- **Fleet Manifest Staging:** Generates and validates `.gitmap/manifest/fleet-clone.json` prior to execution, allowing interactive confirmation and dry-run execution (`--dry-run`).

---

## 3. Go Type Contracts & Architecture

```go
// NodePingResult captures machine reachability metrics
type NodePingResult struct {
    NodeID        string        `json:"nodeId"`
    Alias         string        `json:"alias"`
    Host          string        `json:"host"`
    Port          int           `json:"port"`
    IsICMPOnline  bool          `json:"isIcmpOnline"`
    IsTCPOnline   bool          `json:"isTcpOnline"`
    AvgRTT        time.Duration `json:"avgRtt"`
    PacketLossPct float64       `json:"packetLossPct"`
    StatusText    string        `json:"statusText"`
}

// JSONEnvelope represents the root polymorphic payload container
type JSONEnvelope struct {
    SchemaVersion string          `json:"schemaVersion"`
    Kind          string          `json:"kind"` // "nodes", "macros", "ui", "commit"
    Timestamp     string          `json:"timestamp"`
    Payload       json.RawMessage `json:"payload"`
}

// FleetCloneOptions encapsulates remote clone parameters
type FleetCloneOptions struct {
    TargetNodes  []string `json:"targetNodes"`
    RepoURL      string   `json:"repoUrl"`
    Destination  string   `json:"destination"`
    IsExceptSelf bool     `json:"isExceptSelf"`
    IsParallel   bool     `json:"isParallel"`
    Concurrency  int      `json:"concurrency"`
}
```

All functions return `*appfault.AppError` for failure cases and utilize `Result[T]` wrappers for immutable data propagation.

---

## 4. Subtask Verification Ledger

| Folded Subtask Directory | Source Subtask Files | Verified Criteria & Delivered Artifacts |
| :--- | :--- | :--- |
| `193-typed-json-envelope` | 6 subtask files | Typed JSON envelope architecture, schema validation, `which-format` command, PascalCase fields. |
| `194-import-all-json-what-configs` | 6 subtask files | Batch configuration import (`import-all-json`), `what-configs` (`wc`) CLI, Split-DB safe transactions. |
| `200-nodes-display-cfr-manifest-gitmap-json-and-vscode-optimize` | 6 subtask files | Nodes display CFR manifest, `.gitmap/output/gitmap.json` synchronization, VS Code workspace launch optimization. |
| `201-nodes-clone-cfr-cfrp-fleet-manifest-staging-and-suggestions` | 4 subtask files | Remote fleet clone engine (`gitmap nodes clone`), parallel CFR runner, remote manifest staging, suggestion hints. |
| `67-nodes-cfr-remote-fleet-clone-enhancement` | 6 subtask files | Remote fleet clone error handling, reachability pre-flight checks, non-zero exit code propagation. |

---

## 5. Quality & Coding Guideline Compliance

- **Positive Booleans Only:** Enforced `isIcmpOnline`, `isTcpOnline`, `isExceptSelf`, `isParallel`, `hasFirewallDrop`.
- **Error Propagation:** Swallowed errors completely eliminated; zero silent ignores across network dials and JSON decoders.
- **Relative Path Hygiene:** Zero absolute filesystem paths across all manifest outputs and documentation references.
- **Compiler Cleanliness:** Verified with static analysis; all packages formatted to Unix LF line endings.
