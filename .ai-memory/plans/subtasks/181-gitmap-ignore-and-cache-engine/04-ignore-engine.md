---
status: PENDING
---

# Subtask: Ignore Engine

## Objective
Implement and optimize the ignore engine for gitmap to handle ignore groups and evaluate ignore rules efficiently.

## Requirements
- Optimize the ignore engine to avoid heavy regex evaluation on every file path.
- The ignore engine should correctly identify binary files and mark them as ignored.
- Connect ignore rules logic with the cache engine so that results are reused.
- Support ignore group semantics as required by gitmap workflows.

## Next Steps
1. Audit the current ignore implementation in `cli/cmdignore/` and `cli/cmd/ignore/`.
2. Refactor the ignore evaluation logic to be compatible with caching (return a boolean instead of filtering out immediately).
3. Add tests to ensure binary files are ignored correctly without repeated evaluation.
