# Subtask 03: CI/CD Verification, Minor Bump & Release Ceremony

> **Subtask ID:** Subtask-03  
> **Parent Plan:** `.ai-memory/plans/231-views-ai-inspection-cursor-settings-guard-and-cicd-remediation.md`  
> **Owned Files:**  
> - `version.json`  
> - `package.json`  
> - `01-prompts/README.md`  
> - `02-spec/README.md`  
> - `.ai-memory/README.md`  
> - `README.md`  

---

- [x] 1. Verify all unit tests across the codebase pass cleanly:
   - `go test -v ./cli/cmdcursor/...` (10/10 PASS)
   - `go test -v ./cli/cmdagy/...` (PASS)
- [x] 2. Run linters to guarantee 100% compliance:
   - `python linter-scripts/check-nested-ifs.py` (0 violations)
   - `python linter-scripts/check-enum-and-boolean.py` (0 violations)
   - `python linter-scripts/check-relative-paths.py` (0 violations)
- [x] 3. Execute Minor Version Bump (`v6.505.0` -> `v6.506.0`):
   - `python 03-ai-scripts/37-bump-version.py --tier minor --scope "cursor settings guard validation views ai cleanup and minor release v6.506.0"`
- [ ] 4. Atomic commit via GitMap:
   - `gitmap cpf "cursor - settings guard validation views ai cleanup and minor release v6.506.0"`
- [ ] 5. Release Orchestration:
   - `python 03-ai-scripts/29-release-orchestrator.py -v 6.506.0 -s "cursor settings guard validation views ai cleanup and minor release v6.506.0" --skip-tests`
- [ ] 6. CI/CD Monitoring:
   - Inspect GitHub Actions workflow runs via `gitmap pe -t` to ensure all workflows pass.

---

## 2. Verification Evidence

- `check-nested-ifs.py`: 0 violations.
- `check-enum-and-boolean.py`: 0 violations across 3,228 files.
- `check-relative-paths.py`: 0 violations across 7,782 files.
- Version bumped from 6.505.0 to 6.506.0 in `version.json`, `package.json`, `readme.md`, `what-to-read.md`, and `cli/constants/constants.go`.

