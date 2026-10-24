# GitMap v6.609.0

## What's Changed in v6.609.0

- **Plan 325 enum constants refactor and CI green release**
- Whole-repo enum/constants refactor: every enum in its own package `enums/` subpackage following the cmdcargo exemplar (type, `String()`, `IsValid()`, `Parse` returning `*apperror.AppError`), `String()` instead of `+` concatenation on enum values, app errors enriched with variables/paths, map/slice lookups instead of `==` flag chains, lazy regex (`lazyregex.New`) replacing package-level `regexp.MustCompile`, and shared constants centralized in `cli/constants`.
- Real commit-suite `--dry-run`: the flag now previews files, message, branch, and push target and changes nothing (previously it committed and pushed). Fixed binary verified zero-mutation; see the gitmap skill doc update in coding-guidelines-v24.
- Enum-per-package type convention now a hard rule in the CG execute prompt (coding-guidelines-v24).
- CI fixes on the release path: `lazyregex` data race fixed (all compile-state access under the instance lock; `lazyregex.go` split 471 → 175 lines with new `lazyregex_match.go` at 292 lines), cmdspace typed-nil `*AppError`-as-`error` fix, unused-code removals, nested-if flattens.
- CI green on the release base: CI, Cross-Platform Build, CI Beacon, History Rewrite Smoke, race-detector, all success.

### Quick Install One-Liners

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.609.0/install.ps1 | iex
```

**Linux / macOS (Bash):**
```bash
curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/v6.609.0/install.sh | sh
```
