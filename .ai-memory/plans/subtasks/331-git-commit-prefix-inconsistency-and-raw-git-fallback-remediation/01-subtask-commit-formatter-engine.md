# Subtask 331.1 — Standardize Commit Formatter Engine for Releases

> **Parent Plan:** `.ai-memory/plans/pending/331-git-commit-prefix-inconsistency-and-raw-git-fallback-remediation.md`  
> **Target Package:** `cli/cmdcommit/`  
> **Status:** Pending  

---

## 1. Objectives
1. Inspect `cli/cmdcommit/commit_formatter.go`:
   - Identify `resolveHyphenModuleAndItem(prefix, mod, item)` and `isRedundantModule(prefix, upperMod)`.
   - Prevent stripping of `RELEASE` module when `prefix` is `"Release: "`.
   - If input is `gitmap cpr "v6.615.0 <scope>"`, automatically format as `Release: RELEASE - v6.615.0 <scope>`.
   - If input is `gitmap cpr "VERSION - v6.615.0 <scope>"`, normalize to `Release: RELEASE - v6.615.0 <scope>`.
2. Update unit tests in `cli/cmdcommit/commit_formatter_test.go` to test these scenarios.
3. Run `go test -C cli ./cmdcommit` to verify.
