# Subtask 19: Heatmap Execution Tree Visual Styling

> **Subtask ID:** Subtask-19  
> **Parent Plan:** `.ai-memory/plans/289-improvements-planning-and-execution.md`  
> **Component Spec:** `02-spec/21-app/289-improvements-planning-and-execution/02-component-spec.md`  
> **Target Path:** `cli/heatmap/heatmap_cmd_render.go`  
> **Status:** COMPLETED  

---

## 1. Objectives & Scope
1. Format tree branches (`├── `, `└── `, `│   `) in `cli/heatmap/heatmap_cmd_render.go` with muted subtle colors (`TokyoMuted` / `RoleFGSubtle`).
2. Accentuate action log IDs (`#%d`) with `TokyoAccent` / `RoleAccent`.
3. Highlight execution reasoning states (`PASS`, `FAIL`, `RETRY`, etc.) with semantic status colors.
4. Update both `renderTreeText` and `renderRcTreeText` execution paths.
5. Adhere to strict 15-line limit by modularizing node line formatting.

---

## 2. Detailed Technical Plan

### 2.1 Branch & ID Styling Helpers
- Add dedicated line formatting helpers to `cli/heatmap/heatmap_cmd_render.go`:
  ```go
  func formatTreeBranch(branch string) string {
  	return termout.RoleSGR(termout.RoleFGSubtle) + branch + termout.Reset()
  }

  func formatNodeLogID(id int64) string {
  	return termout.RoleSGRBold(termout.RoleAccent) + fmt.Sprintf("#%d", id) + termout.Reset()
  }

  func formatReasoningState(state string) string {
  	role := termout.RoleFG
  	switch state {
  	case "PASS", "SUCCESS", "OPTIMAL":
  		role = termout.RoleOK
  	case "FAIL", "ERROR", "CRITICAL":
  		role = termout.RoleError
  	case "RETRY", "WARN":
  		role = termout.RoleWarn
  	}

  	return termout.RoleSGR(role) + state + termout.Reset()
  }
  ```

### 2.2 Refactor `renderTreeText`
- Use formatted branch and node elements:
  ```go
  func renderTreeText(out io.Writer, nodes []store.HeatmapExecutionTreeNode, indent string) {
  	for i, node := range nodes {
  		branch, nextIndent := resolveBranchIndent(i, len(nodes), indent)
  		lineInfo := formatFileLineInfo(node.FilePath, node.StartLine, node.EndLine)
  		styledBranch := formatTreeBranch(branch)
  		styledID := formatNodeLogID(node.HeatmapActionLogId)
  		styledState := formatReasoningState(node.ReasoningState)

  		fmt.Fprintf(out, "%s%s%s (%s) %s%s: %s\n", indent, styledBranch, styledID,
  			styledState, node.ActionType, lineInfo, node.ReasoningNotes)

  		if len(node.Children) > 0 {
  			renderTreeText(out, node.Children, nextIndent)
  		}
  	}
  }
  ```

### 2.3 Refactor `renderRcTreeText`
- Mirror the formatting in `renderRcTreeText` for root cause execution trees.

---

## 3. Strict Coding Guidelines Compliance
- **Function Limit:** Recursive render functions decomposed to <= 15 lines.
- **Guard Clauses:** Early returns in helper functions.
- **Theme Safety:** Uses `termout.RoleSGR` instead of hardcoded escape codes.
- **Return Formatting:** Single blank line before return statements.

---

## 4. Verification Criteria
- [x] Tree glyphs (`├── `, `└── `) display in subtle muted tone.
- [x] Action Log IDs (`#102`, `#103`) pop with distinct accent color.
- [x] Status labels reflect success (green), warning (amber), and failure (red).
- [x] Indentation hierarchy remains visually aligned across nested children.
