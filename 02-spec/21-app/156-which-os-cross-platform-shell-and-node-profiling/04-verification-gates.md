# Specification 156: Verification Gates & Invariants

## Invariants

1. **Enum Reusability Invariant:**
   `which-os` MUST reuse canonical constants defined in `cli/constants/os_targets.go`. No ad-hoc string literals like `"win32"`, `"ubuntu-linux"`, or lowercase mismatch in output mappings.

2. **JSON Consistency Invariant:**
   `gitmap which-os --json` MUST output valid JSON matching `OSInfoReport` schema with all fields populated (`osType`, `osGroup`, `osVersion`, `buildVersion`, `architecture`, `platform`, `hostname`, `numCpu`, `hasGit`, `hasBash`, `hasPowerShell`).

3. **Fallback & Graceful Suggestions Invariant:**
   Missing binaries (`git`, `bash`, `pwsh`) MUST NEVER cause panic, unhandled error, or exit without explicit, actionable installation suggestions.

4. **Split-DB Persistence Invariant:**
   When node discovery occurs or Git path is verified, data MUST be persisted in SQLite Split-DB (`gitmap.db`), enabling zero-roundtrip retrieval for future commands.

5. **OS Filtering Invariant:**
   Installers by default skip mismatched OS targets (e.g. Unix scripts skip on Windows, Windows scripts skip on POSIX) unless `--force-all` is explicitly set.

6. **Function Size and Guidelines Invariant:**
   All newly created or refactored Go functions MUST adhere to <= 8-15 lines, single return / Result types, Unix LF endings, affirmative booleans (`hasGit`, `isWin`), and domain `*appfault.AppError`.
