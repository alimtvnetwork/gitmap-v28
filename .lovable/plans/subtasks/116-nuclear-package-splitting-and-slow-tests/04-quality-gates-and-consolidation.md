# Subtask 04: Quality Gates, Test Inventory Sync & Subtask Consolidation

## Objective
Update the test inventory manifest, run quality linters, ensure zero circular dependencies, consolidate all subtasks into completed plan, and push to main.

## Steps
1. Re-scan and record all modified files via `python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`.
2. Format Go code with `python 03-ai-scripts/26-go-code-formatter.py`.
3. Check return newline styling with `python 03-ai-scripts/04-newline-fixer.py --fix gitmap/cmd`.
4. Run `go vet ./...` to guarantee 100% clean compilation and zero import cycles.
5. Consolidate all subtasks into `.lovable/plans/completed/116-nuclear-package-splitting-and-slow-tests.md`.
6. Update `.lovable/plans/01-index.md`.
7. Commit and push cleanly to `origin main`.
