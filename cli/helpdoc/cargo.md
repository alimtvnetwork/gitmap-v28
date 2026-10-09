# Rust Cargo Runner

Rust Cargo package manager and compiler toolchain runner with automatic install suggestions.

## Alias

crg

## Usage

```bash
gitmap cargo [subcommand] [args...]
gitmap crg [subcommand] [args...]
gitmap cargo <args...> [--install]
```

## Description

`gitmap cargo` provides seamless interaction with the Rust Cargo toolchain across Windows, macOS, and Linux. When Cargo is installed, commands are forwarded directly to the `cargo` binary. If Cargo or Rust is missing from the environment, GitMap displays formatted installation suggestions and multi-platform install options, or automatically installs it when `--install` is supplied.

## Subcommands

| Subcommand | Aliases | Description |
|------------|---------|-------------|
| `status` | `st`, `info` | Inspect Rust and Cargo toolchain installation status, version, and binary path |
| `install` | `in` | Install Rust and Cargo toolchain via `gitmap install cargo` |

## Examples

### Check Cargo and Rust toolchain status

```bash
gitmap cargo status
```

### Run Cargo with install suggestions if missing

```bash
gitmap cargo build
```

### Auto-install Cargo if missing and run command

```bash
gitmap cargo --install check
```

### Manually trigger Cargo installation

```bash
gitmap cargo install
```
