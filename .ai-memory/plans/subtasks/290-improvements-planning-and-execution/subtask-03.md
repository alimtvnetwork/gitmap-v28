# Subtask Plan 03: Eradicate `cliexit.Reportf` Returning `nil` in `cli/cmdconfig/config_export.go`

- **Subtask ID:** `subtask-03`
- **Subtask Name:** `03-config-export-cliexit-eradication`
- **Parent Spec:** `02-spec/21-app/290-improvements-planning-and-execution/01-architecture-spec.md`
- **Assigned Agent Role:** Spec Author 01 / Implementation Worker 01
- **Target Files:**
  - `cli/cmdconfig/config_export.go`
- **Status:** Completed

---

## 1. Targeted Problem & Architectural Rationale

In `cli/cmdconfig/config_export.go`, functions handling configuration export violate **Rule 1 (Error Return Sovereignty)**:

1. In `dispatchExportArgs(args []string) error`:
   ```go
   if !isSupportedExportTool(canonicalTool) {
       cliexit.Reportf("export-config", "export", inputTool, fmt.Errorf("unsupported tool '%s'", inputTool))

       return nil
   }
   ```
   When an unsupported tool is provided, `cliexit.Reportf` is invoked, but the function returns `nil`! This tricks the upstream caller into believing the command succeeded, corrupting command audit logging.
2. In `exportSingleToolConfig(canonicalTool, inputTool, targetArg string) error`:
   ```go
   err := writeBundleJSON(bundle, finalPath)
   if err != nil {
       cliexit.Reportf("export-config", "export", canonicalTool, err)

       return err
   }
   ```
   `cliexit.Reportf` is invoked inside a leaf helper function rather than allowing the root dispatcher to handle failure reporting.
3. Leaf functions must NEVER invoke exit handlers or report exit statuses and return `nil`. All errors must be wrapped and returned directly so the top-level dispatcher in `cli/cmd/` captures them in `finishCommandAudit`.

---

## 2. Component Boundaries & Code Signatures

### 2.1 Refactor `dispatchExportArgs`
- **File:** `cli/cmdconfig/config_export.go`
- **Signature:** `func dispatchExportArgs(args []string) error`
- **Requirements:**
  - When `!isSupportedExportTool(canonicalTool)`, immediately return `apperror.NewSimple(fmt.Sprintf("unsupported tool '%s'", inputTool), "E1002")` or `fmt.Errorf("unsupported tool '%s'", inputTool)`.
  - Remove all calls to `cliexit.Reportf`.
  - Keep function length <= 12 lines.

```go
func dispatchExportArgs(args []string) error {
	canonicalTool, inputTool, targetArg := parseExportArgs(args)
	if !isSupportedExportTool(canonicalTool) {
		return apperror.NewSimple(fmt.Sprintf("unsupported tool '%s'", inputTool), "E1002")
	}

	if canonicalTool == "all" {
		return exportAllToolConfigs(targetArg)
	}

	return exportSingleToolConfig(canonicalTool, inputTool, targetArg)
}
```

### 2.2 Refactor `exportSingleToolConfig`
- **File:** `cli/cmdconfig/config_export.go`
- **Signature:** `func exportSingleToolConfig(canonicalTool, inputTool, targetArg string) error`
- **Requirements:**
  - Remove `cliexit.Reportf`.
  - If `writeBundleJSON` fails, wrap and return the error directly: `return apperror.WrapSimple(err, "write bundle: "+canonicalTool)`.
  - Function length <= 15 lines.

```go
func exportSingleToolConfig(canonicalTool, inputTool, targetArg string) error {
	defaultName := resolveDefaultConfigFileName(inputTool, canonicalTool)
	finalPath := resolveConfigFilePath(targetArg, defaultName)
	srcDir := resolveToolConfigDir(canonicalTool)
	files, extensions := resolveExportFilesAndExtensions(canonicalTool, srcDir)
	bundle := buildConfigBundle(canonicalTool, srcDir, files, extensions)
	err := writeBundleJSON(bundle, finalPath)
	if err != nil {
		return apperror.WrapSimple(err, "write bundle: "+canonicalTool)
	}

	reportExportSuccess(canonicalTool, finalPath, len(files))

	return nil
}
```

### 2.3 Clean Up Imports
- **File:** `cli/cmdconfig/config_export.go`
- **Requirements:**
  - Remove unused import `"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"`.
  - Ensure `"github.com/alimtvnetwork/gitmap-v28/cli/apperror"` is imported.

---

## 3. Implementation Verification Checklist

- [x] All `cliexit.Reportf` calls removed from `cli/cmdconfig/config_export.go`.
- [x] Unsupported tool condition returns explicit error instead of `nil`.
- [x] `writeBundleJSON` failure returns wrapped `apperror` directly.
- [x] Unused `cliexit` package import removed.
- [x] Function lengths <= 15 lines each.
- [x] Zero nested `if` statements.
- [x] Blank line before `return` and after `}` preserved.
