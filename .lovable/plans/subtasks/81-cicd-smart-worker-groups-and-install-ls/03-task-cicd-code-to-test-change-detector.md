# Subtask 03: Code-to-Test Mapping & Impact-Based Incremental Change Detection

## 1. Description
Map each test to its corresponding source implementation file and function. Compute SHA256 hashes of both target code/functions and test functions. Implement incremental skip logic so that a test executes **if and only if** its target function/code or test function has changed since the last passing run.

## 2. Files to Modify
- `03-ai-scripts/06-cicd-local-runner.py`:
  - Implement `extract_go_function_hashes(filepath: Path) -> dict[str, str]` to parse function bodies and compute SHA256 hashes.
  - Implement `map_test_to_code(test_pkg: str, test_func: str, test_file: Path, source_funcs: dict[Path, dict[str, str]]) -> tuple[str, str, str]`:
    - Matches test function name (e.g. `TestResolveToolStatusFromDB` -> `resolveToolStatus`).
    - Maps to target source file (e.g. `installlist.go`) and function body hash.
    - Fallback to source file hash when 1:1 function name correlation is absent.
  - In `should_run_test(test_info: dict, source_funcs: dict) -> bool`:
    - Compare current code hash with cached `code_hash`.
    - Compare current test hash with cached `test_hash`.
    - If either hash differs or previous status was failed -> MUST run test.
    - If hashes match and previous run was pass -> skip test (`[cached]`).

## 3. Invariants
- Fast change checking (< 200ms across whole repository).
- Zero false negatives: if a function changes, its test MUST run.
