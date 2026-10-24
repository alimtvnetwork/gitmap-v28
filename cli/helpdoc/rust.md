# Rust Toolchain Manager

Native Rust compiler and toolchain manager providing unified access to `rustc`, `rustup`, and `cargo` with automated LLVM/libclang path resolution.

## Synopsis

```bash
gitmap rust [subcommand|command] [args...]
gitmap rst [subcommand|command] [args...]
gitmap rustc [args...]
```

## Description

`gitmap rust` provides a unified cross-platform frontend for the Rust toolchain. On Windows, it automatically discovers installed LLVM/libclang directories (e.g. `D:\LLVM\LLVM\bin`, `C:\Program Files\LLVM\bin`) and injects `LIBCLANG_PATH` and `PATH` into the execution environment so that crates using `bindgen` compile seamlessly without manual environment variable setup.

## Subcommands

| Subcommand | Aliases | Description |
|------------|---------|-------------|
| `status` | `st`, `info` | Display toolchain dashboard (rustc, cargo, rustup, libclang versions and paths) |
| `c` | `rustc` | Forward arguments directly to `rustc` with enriched environment |
| `up` | `rustup` | Forward arguments directly to `rustup` |
| `build` | `b` | Forward build command to `cargo build` |
| `test` | `t` | Forward test command to `cargo test` |
| `check` | `k` | Forward check command to `cargo check` |
| `clippy` | `clip` | Forward clippy linter to `cargo clippy` |

## Aliases

- `gitmap rst`
- `gitmap rustc`

## Examples

### View Rust and LLVM Toolchain Status

```bash
gitmap rust status
```

### Run Rust Compiler Directly

```bash
gitmap rust c main.rs -o main.exe
```

### Build Workspace with Enriched LLVM Environment

```bash
gitmap rust build --release
```

### Run Rust Linter

```bash
gitmap rust clippy
```

## See Also

- `gitmap hd cargo` — Native Cargo package manager runner.
- `gitmap hd kill` — Process termination utility.
