# Subtask 01: Nodes Table Sequence Inversion & Commands Matrix Role Pruning

> **Parent Plan:** [.ai-memory/plans/pending/65-gitmap-update-all-zip-and-fixes.md](../../pending/65-gitmap-update-all-zip-and-fixes.md)  
> **Spec Reference:** [02-spec/21-app/65-gitmap-update-all-zip-and-fixes/01-architecture-spec.md](../../../../02-spec/21-app/65-gitmap-update-all-zip-and-fixes/01-architecture-spec.md)  
> **Status:** `PENDING`  
> **Target Files:**
> - `cli/cmd/nodes_cmd.go`
> - `cli/cmd/nodes_cmd_test.go`

---

## 1. Technical Objective

Reorder the rendering hierarchy of the `gitmap nodes` command so that:
1. The **Supported Commands & Capabilities Matrix** table is rendered at the **TOP**.
2. Redundant **`ROLE`** column is removed from the Commands Matrix table, expanding the command suggestions width.
3. The **Primary Node Status Table** is rendered at the **BOTTOM**, directly followed by the summary line `Total: X registered node(s)`.
4. Unit tests in `cli/cmd/nodes_cmd_test.go` verify table ordering, column absence, and summary counts.

---

## 2. Implementation Scope & File Edits

### Target 1: `cli/cmd/nodes_cmd.go`
- **Invert Table Invocations in `renderUnifiedNodesTable`**:
  ```go
  func renderUnifiedNodesTable(out io.Writer, nodes []UnifiedFleetNode) error {
      if len(nodes) == 0 {
          return renderEmptyUnifiedNodesNotice(out)
      }
      // Render Commands Matrix Table first (TOP)
      if err := renderCommandsMatrixTable(out, nodes); err != nil {
          return err
      }
      // Render Primary Nodes Status Table second (BOTTOM)
      if err := renderPrimaryNodesTable(out, nodes); err != nil {
          return err
      }
      renderUnifiedNodesFooter(out)
      return nil
  }
  ```
- **Prune `ROLE` from `renderCommandsMatrixHeader`**:
  - Remove `padCell("ROLE", 14)`.
  - Rebalance column widths:
    - `colAlias`: 16 chars
    - `colSys`: 22 chars
    - `colCmds`: 70 chars
    - Total: 110 characters matching the standard divider line.
  - Header line format: `fmt.Sprintf("  %s %s %s\n", colAlias, colSys, colCmds)`
- **Prune `ROLE` from `renderCommandsMatrixRow`**:
  - Remove `role := formatCellRole(n.Role, 14)` formatting.
  - Row format: `fmt.Fprintf(out, "  %s %s %s\n", alias, subsys, cmds)`
- **Preserve Primary Nodes Table**:
  - Retain `renderPrimaryNodesTable` with full columns (`ALIAS`, `ROLE`, `HOST (IP:PORT)`, `USER`, `STATUS`, `ENROLLED`).
  - Keep `Total: %d registered node(s)` subtotal line at the bottom of the status table.

### Target 2: `cli/cmd/nodes_cmd_test.go`
- **Update `TestUnifiedNodes_RenderTable`**:
  - Assert that `SUPPORTED COMMANDS & CLUSTERS` appears before `HOST (IP:PORT)` in the rendered string output:
    ```go
    matrixIdx := strings.Index(out, "SUPPORTED COMMANDS & CLUSTERS")
    primaryIdx := strings.Index(out, "HOST (IP:PORT)")
    if matrixIdx == -1 || primaryIdx == -1 {
        t.Fatalf("expected both matrix and primary headers to be present")
    }
    if matrixIdx > primaryIdx {
        t.Errorf("expected commands matrix (index %d) to be rendered above primary nodes table (index %d)", matrixIdx, primaryIdx)
    }
    ```
  - Assert that `Total: 2 registered node(s)` is present after the primary table.
  - Assert that the matrix table header does not contain `ROLE` between `NODE (ALIAS)` and `SUBSYSTEMS`.

---

## 3. Verification Protocol

- Run unit test suite:
  ```bash
  go test -v ./cli/cmd -run TestUnifiedNodes
  ```
- Verify static analysis:
  ```bash
  go vet ./cli/cmd
  ```
- Manual inspection check:
  - Run `gitmap nodes` and verify that the capabilities matrix appears first and the node status list appears second with the total count at the bottom.
