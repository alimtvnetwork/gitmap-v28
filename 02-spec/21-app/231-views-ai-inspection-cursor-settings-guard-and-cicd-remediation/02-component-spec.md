# Component Specification: Views AI Inspection, Cursor Settings Guard & CI/CD Remediation

> **Document Version:** 1.0.0  
> **Status:** Active  
> **Scope:** Core CLI & Web UI Components  
> **Traceability:** Task-231  

---

## 1. Component Details

### Component A: Cursor Settings Validator (`cli/cmdcursor`)
- **Package**: `cmdcursor`
- **Source File**: `cli/cmdcursor/cursor_settings.go`
- **Test File**: `cli/cmdcursor/cursor_settings_test.go`
- **Function**: `requireCursorSettingsFile(path string) error`
  - Positive existence check via `isPathPresent(path)`.
  - Guard early return `nil`.
  - Construct typed AppError with diagnostic remediation text:
    `Cursor settings.json not found at <path>. Run 'gitmap cursor settings apply' to create it, then retry sync.`

### Component B: Web UI Settings Controller (`src/pages/Settings.tsx`)
- **Source File**: `src/pages/Settings.tsx`
- **Hook / Function**: `fetchInstances()`
- **Remediation**:
  - Replace nested conditional block with guard clause:
    ```ts
    if (!response.ok) {
      return;
    }
    const json = await response.json();
    if (json.isSuccess && Array.isArray(json.instances)) {
      setInstances(json.instances);
      return;
    }
    ```

### Component C: Discovery Engine Spacing Hygiene (`cli/cmdagy`)
- **Source File**: `cli/cmdagy/agy_instance_discovery.go`
- **Functions**: `unmarshalWinProcItems`, `truncateSliceByLimit`
- **Enforcement**: Vertical blank line spacing after all block closures (`}`).

---

## 2. Release & Governance Criteria

1. **Working Tree**: Zero untracked or dirty files prior to release.
2. **Version Bump**: `v6.505.0` -> `v6.506.0` (Minor bump per Rule 0).
3. **Commit Convention**: Single atomic commit with hyphenated format:
   `gitmap cpf "cursor - settings guard validation views ai cleanup and minor release v6.506.0"`
4. **CI/CD Health**: Verify via `gitmap pe -t` that all workflows succeed.
