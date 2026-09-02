---
name: gitmap-cli-command
description: >-
  Autonomously design, author, register, and test Gitmap CLI subcommands adhering to AST parity,
  help text simulation constraints, cliexit error reporting, and code generation.
---

# Gitmap CLI Command Authoring & AST Parity Skill

Autonomously implement, refactor, and verify CLI commands in the Gitmap Go toolchain adhering to `spec/01-app/`, `spec/17-consolidated-guidelines/23-generic-cli.md`, and `.lovable/strictly-avoid.md`.

## Core Checkpoints & Mandatory Invariants

1. **AST Parity & Registry Synchronization:**
   - Define top-level command constant in `gitmap/constants/constants_cli.go` inside a block marked with `// gitmap:cmd top-level`.
   - If a command is an internal alias or exempt, mark with `// gitmap:cmd skip`.
   - Update `topLevelCmds()` in `gitmap/constants/cmd_constants_test.go` whenever adding or changing command constants.
   - Run `go test ./gitmap/constants/... -run TestTopLevelCmdRegistryMatchesAST -count=1` to guarantee AST parity.

2. **Command Help Text Standards:**
   - Dedicated markdown file in `gitmap/helptext/<command>.md`.
   - Maximum **120 lines** total.
   - Realistic execution simulation block: **3 to 8 lines**.
   - ALWAYS use fenced code blocks (` ``` `); NEVER use 4-space indentations (prevents golden test failures).
   - Verify golden tests: `go test ./gitmap/helptext/... -run Golden -count=1`.

3. **Standardized Error Reporting:**
   - NEVER use bare `fmt.Fprintln(os.Stderr, err)` in `gitmap/cmd/`.
   - Always route error exits through `cliexit.Reportf` or `cliexit.Fail`.
   - Wrap internal errors in `apperror.AppError` with operational context.

4. **Code Generation:**
   - NEVER modify CLI constants, shell commands, or help text without running `go generate ./...` in the `gitmap/` directory.
   - Verify zero drift using `git diff --exit-code .`.

5. **Size & Function Limits:**
   - Functions: <= 15 lines max.
   - Files: <= 200 lines max.
   - Break complex command logic into dedicated helper files (e.g. `<command>ops.go`, `<command>dispatch.go`).
