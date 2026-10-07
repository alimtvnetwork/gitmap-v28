# Subtask 02: Settings.tsx Nested If Flattening & Discovery Spacing Hygiene

> **Subtask ID:** Subtask-02  
> **Parent Plan:** `.ai-memory/plans/231-views-ai-inspection-cursor-settings-guard-and-cicd-remediation.md`  
> **Owned Files:**  
> - `src/pages/Settings.tsx`  
> - `cli/cmdagy/agy_instance_discovery.go`  

---

- [x] 1. In `src/pages/Settings.tsx`:
   - Invert `response.ok` into an early guard clause:
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
- [x] 2. In `cli/cmdagy/agy_instance_discovery.go`:
   - Add vertical blank lines after closing braces in `unmarshalWinProcItems` and `truncateSliceByLimit`.
   - Format with `gofmt -w`.
- [x] 3. Run linters:
   - `python linter-scripts/check-nested-ifs.py` (verify 0 violations).
   - `python linter-scripts/check-enum-and-boolean.py` (verify 0 violations).

---

## 2. Verification Evidence

- `check-nested-ifs.py --changed-only`: 0 nested if violations across scanned files.
- `check-enum-and-boolean.py`: 0 violations across 3,228 files.
- `check-relative-paths.py`: 0 absolute path violations across 7,782 files.
- `go test -v ./cmdagy/...`: PASS in 6.038s.
- Status: **DONE**

