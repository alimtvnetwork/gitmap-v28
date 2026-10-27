# Subtask 336.02 — CLI Command Registration, Dispatch and Themed Help Documentation

> **Sequence ID:** 336  
> **Parent Plan:** [336-automatic-binary-swap-and-skill-synchronization.md](file:///d:/work/gitmap/.ai-memory/plans/pending/336-automatic-binary-swap-and-skill-synchronization.md)  
> **Target Packages:** `cli/`, `cli/cmdswap/`, `cli/helpdoc/`  
> **Status:** Pending  

---

## 1. Objectives

1. Register `gitmap swap-binary` command and aliases in `cli/root.go`:
   - Primary: `swap-binary`
   - Aliases: `swap-bin`, `bin-swap`, `replace-bin`, `replace-binary`, `binary swap`, `swap`
   - Group: System & Maintenance commands.
2. Implement CLI flag parsing in `cli/cmdswap/swap_flags.go`:
   - `--from` (source binary path)
   - `--to` / `--target` / `-t` (destination deployed binary path)
   - `--pid` (specific target PID to terminate)
   - `--kill` / `-k` (terminate active locking processes)
   - `--force` / `-f` (forceful process termination)
   - `--restart` / `--reopen` (execute verification/UI after swap)
   - `--script-only` (print generated handoff script without running)
   - `--json` (emit machine-readable JSON envelope)
3. Implement Lipgloss terminal UI status card in `cli/cmdswap/swap_ui.go`:
   - Render header, source/target paths, detected PIDs, and handoff script execution status.
4. Author themed help documentation entry in `cli/helpdoc/`:
   - Register `swap-binary` in command help dictionary with examples, flags, and cross-platform behavior.
   - Add examples for common workflows (e.g., swapping freshly built `./cli/gitmap.exe` into active PATH binary).

---

## 2. File Manifest

- `cli/cmdswap/swap_cmd.go` — Command entry point and flag binding.
- `cli/cmdswap/swap_flags.go` — Parsing and normalization of CLI flags.
- `cli/cmdswap/swap_ui.go` — Lipgloss terminal UI card formatting and JSON response writer.
- `cli/root.go` — Root command dispatch mapping for `swap-binary` and aliases.
- `cli/helpdoc/helpdoc_swap_binary.go` — Help documentation, examples, and flag descriptions.

---

## 3. Verification Criteria

- [ ] `gitmap swap-binary --help` displays full usage, flag options, and examples.
- [ ] Aliases (`gitmap swap`, `gitmap bin-swap`) correctly route to `cmdswap.RunSwapCommand`.
- [ ] `--json` flag formats structured output complying with `SwapResult` contract.
- [ ] UI summary card renders cleanly with Lipgloss theme formatting in interactive terminals.
