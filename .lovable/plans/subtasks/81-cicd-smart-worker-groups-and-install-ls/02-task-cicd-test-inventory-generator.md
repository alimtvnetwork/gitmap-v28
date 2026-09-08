# Subtask 02: CI/CD Test Inventory Generator & Execution Timings Cache

## 1. Description
Implement a test inventory generator at the very start of `03-ai-scripts/06-cicd-local-runner.py` that discovers all existing tests in the codebase (all Go unit tests `Test*` across `*_test.go`, linter gates, and E2E suites), catalogs them into `.lovable/test-inventory.json`, and records test execution timings.

## 2. Files to Modify
- `03-ai-scripts/06-cicd-local-runner.py`:
  - Define `TEST_INVENTORY_PATH = Path(".lovable/test-inventory.json")`
  - Implement `discover_test_inventory(repo_root: Path) -> dict[str, Any]`:
    - Walk `gitmap/` for all `*_test.go` files.
    - Extract all `func (Test[A-Za-z0-9_]+)\(` function declarations.
    - Catalog test name, test file, package path, last status, last execution duration, and last run timestamp.
  - Implement `save_test_inventory(path: Path, inventory: dict[str, Any]) -> None`.
  - Add CLI flag `--inventory-only` to display/generate test catalog.
  - Update `update_test_timings` to persist each test's duration in seconds.

## 3. Invariants
- Discovery must be fast (< 250ms for 2,277+ tests).
- JSON format must be clean and human-readable with indentation.
- Timings must record floating-point seconds.
