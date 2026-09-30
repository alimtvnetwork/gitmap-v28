# Subtask 01: Nodes Display Parity and Supported Commands Table

> **Parent Plan:** `200-nodes-display-cfr-manifest-gitmap-json-and-vscode-optimize.md`  
> **Status:** `PENDING`  
> **Target Files:**
> - `cli/cmd/nodes_cmd.go`
> - `cli/cmd/nodes_cmd_test.go`

---

## Technical Specification

1. **Table 1: Primary Nodes Table**:
   - Matches `gitmap ssh nodes` styling exactly.
   - Header:
     `  ALIAS            ROLE           HOST (IP:PORT)         USER           STATUS                ENROLLED`
   - Divider:
     `  --------------------------------------------------------------------------------------------------------------` (110 chars)
   - Rows:
     - Alias: Bold White (width 16)
     - Role: Colored (width 14)
     - Host:Port: Plain White (width 22)
     - User: Dim (width 14)
     - Status: Colored icon + text (width 21)
     - Enrolled: Dim timestamp (width 19)
   - Subtotal:
     `  Total: %d registered node(s)`

2. **Table 2: Supported Commands Matrix**:
   - Header:
     `  NODE (ALIAS)     ROLE           SUBSYSTEMS           SUPPORTED COMMANDS`
   - Divider:
     `  --------------------------------------------------------------------------------------------------------------`
   - Rows:
     - Subsystems: "SSH, Cluster, SC"
     - Supported Commands:
       - If has SSH: `gitmap ssh <alias>, gitmap exec <alias> "<cmd>"`
       - If has Cluster: `gitmap cluster exec <alias> "<cmd>"`
       - If has SC: `gitmap sc exec "<cmd>"`
       - If ready: `gitmap ssh deploy-keys, gitmap nodes ping <alias>`

3. **Footer Suggestions**:
   - Clean, professional suggestions block:
     `[tip] Fleet Operations & Supported Command Suggestions:`
     `  • Shell Access:            gitmap ssh <alias>                (or: gitmap ssh login <alias>)`
     `  • Command Execution:       gitmap exec <alias> "<cmd>"       (or: gitmap ssh exec <alias> "<cmd>")`
     `  • Cluster Command:         gitmap cluster exec <alias> "<cmd>"`
     `  • Broadcast to All Nodes:  gitmap sc exec "<cmd>"            (or: gitmap ssh exec all "<cmd>")`
     `  • Sync Keys & Ping Nodes:  gitmap ssh deploy-keys            | gitmap nodes ping`

4. **Tests**:
   - Verify table output format in `cli/cmd/nodes_cmd_test.go`.
