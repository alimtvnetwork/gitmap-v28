# Specification 17 — Chapter 6: Replace Engine Diagnostics & Fixes

## 1. Diagnostic Analysis of Prior Failures
In previous versions, `gitmap replace <old> <new>` failed to find matches in certain repos (e.g. `gitmap replace errorwrapper-v3 errorwrapper-v4` found 0 matches across 634 files) due to:
1. **Path Normalization Discrepancies**:
   - On Windows, `git rev-parse --show-toplevel` returns forward slashes `/d/work/...` or `D:/work/...`, while `filepath.WalkDir` produces `D:\work\...`. Prefix exclusions (`.gitmap/release`) failed prefix tests or skipped entire trees.
2. **Overly Restrictive Extension / Binary Filtering**:
   - Binary sniffing heuristic classified files with UTF-16 BOM or zero bytes erroneously.
3. **String Token vs Boundary Variations**:
   - Code files frequently contained variations: `errorwrapper-v3`, `errorwrapper/v3`, `errorwrapper_v3`, `errorwrapper.v3`, or `errorwrapper v3`.

## 2. Comprehensive Replace Engine Enhancements
1. **Robust Slash & Path Normalization**:
   - All paths are passed through `filepath.ToSlash` and `filepath.Clean` before checking exclusion patterns.
2. **Smart Token Variations & Match Reporting**:
   - When literal replacement is given, the engine checks exact match, hyphenated variations, and slash variations, clearly logging per-file diff snippets.
3. **Preview & Diffs**:
   - Displays exact line numbers and colorized before/after snippets for matches.
