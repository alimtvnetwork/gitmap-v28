# Learned Memory: plan-coding-guideline-audit-v4

## Overview

Comprehensive audit and enforcement of v1.4.5 coding guidelines across Golang backend packages, TypeScript/React UI pages, and PowerShell CI/build scripts.

## Key Patterns & Discovered Conventions

### 1. Parity Marker Tests (`cmd_constants_parity_test.go`)

- **Pattern**: Any `const` block containing identifiers starting with `Cmd*` must be annotated with `// gitmap:cmd top-level` or each non-top-level spec must be tagged with `// gitmap:cmd skip`.
- **Learned**: Node/subcommand execution verbs (e.g. `LifecycleCmdShutdown`, `LifecycleCmdReboot`) should avoid the naked `Cmd*` prefix unless intended as top-level CLI commands, or be named with domain-specific prefixes like `LifecycleCmd*`.

### 2. Helptext Coverage Guard (`coverage_test.go`)

- **Pattern**: Every primary `Cmd*` identifier in `constants_cli.go` requires a corresponding Markdown file in `helptext/<id>.md`.
- **Learned**: Subcommands that belong to a parent command group (e.g. `groups` under `gitmap ls` / `gitmap group`) should use `SubCmd*` prefix (e.g. `SubCmdGroups = "groups"`), avoiding stray `CmdGroups` declarations that trigger coverage test failures.

### 3. PowerShell Error Logging & Warning Formatting

- **Pattern**: Never leave `catch {}` empty. Log actionable operational warnings without throwing terminating exceptions for optional cleanup routines:
  ```powershell
  catch {
      Write-Warning "[$operation] $_"
  }
  ```
- **Learned**: Prevents silent swallowing of diagnostic errors in CI pipelines while maintaining smooth script execution.

### 4. Positive Boolean Extractions & Guard Clauses

- **Pattern**: Negations (`!isSuccess`, `!ok`, `!open`, `!api`) must be extracted into explicitly named positive boolean variables:
  ```go
  isMissingEntry := !ok
  if isMissingEntry == true {
      return
  }
  ```
  ```typescript
  const isClosed = !open;
  if (isClosed) return;
  ```

### 5. Function Size Decomposition ($\le 15$ lines cap, $\le 8$ lines preferred)

- **Pattern**: Large functions should be decomposed into single-responsibility helpers using Shell + Wire or Table Dispatch patterns.
- **Learned**: Reduces cyclomatic complexity, enforces isolated testability, and complies with hard line limits.
