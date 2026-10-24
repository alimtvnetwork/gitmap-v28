# Subtask 318-01: Native Toolchain Commands (Cargo, Rust, Kill) & LLVM Auto-Discovery

Status: COMPLETED

## Objectives
1. Implement `cli/cmdtoolchain/` package:
   - `toolchain_llvm.go`: Automatic detection of `LIBCLANG_PATH` on Windows (scanning `D:\LLVM\LLVM\bin`, `C:\Program Files\LLVM\bin`, etc.) and injection into child environment.
   - `toolchain_env.go`: Environment enrichment (`PATH` augmentation with LLVM and cargo binaries).
   - `toolchain_cargo.go`: `gitmap cargo` execution with LLVM auto-injection.
   - `toolchain_rust.go`: `gitmap rust` execution with `status`, `build`, `test`, `check`, `clippy` forwarding.
   - `toolchain_kill.go`: Cross-platform process termination by PID, name, or port with safety guards.
2. Register commands in `cli/constants/constants_commands.go` (`CmdCargo`, `CmdRust`, `CmdKill`) and in `cli/cmd/roottooling.go`.
3. Unit tests in `cli/cmdtoolchain/toolchain_test.go`.

## Verification Results
- All unit tests in `cli/cmdtoolchain` pass with 100% success.
- CLI constants AST uniqueness tests in `cli/constants` pass.
- Dispatched tooling commands in `cli/cmd` pass.
- Coding guidelines verified: <= 15 lines per function, positive booleans (`is*`, `has*`), zero nested ifs, leaf error return sovereignty.
