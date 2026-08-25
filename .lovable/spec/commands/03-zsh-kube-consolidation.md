# 03-zsh-kube-consolidation: Detailed Specification (Steps 1-100)

## Overview
This specification details the 100-step blueprint to consolidate the newly provided `scripts/kubernetes` (specifically `02-ubuntu-install` and `01-base-shell-scripts`) into native Go-based `gitmap` commands.

## Sub-System 1: File Line Injector (Idempotent & Sequence-Aware)
- **Goal:** Safely append or insert lines into `.zshrc` and `authorized_keys`.
- **Logic:** Read file line-by-line. If the target line or block exists, skip. If it must follow a specific sequence (e.g., "after `plugins=(...)`"), find the sequence marker and insert. Write to a temporary file first, then atomically swap to prevent corruption.

## Sub-System 2: ZSH Management (`gitmap zsh`)
- `gitmap zsh install`: Cross-platform (Ubuntu, CentOS) installation of ZSH and Oh-My-Zsh.
- `gitmap zsh theme-change <theme>`: Modifies the `ZSH_THEME` line in `.zshrc` intelligently.
- `gitmap zsh completion`: Generates and dynamically sources ZSH completions.

## Sub-System 3: OS Utilities
- `gitmap os kill <process>`: Cross-platform process termination.
- `gitmap os install aria2c`: Automated installation of aria2c.

## Sub-System 4: SSH Auth Integration (Tasks 40-50 Roll-in)
- `gitmap ssh-join add-auth`: Uses the new intelligent File Line Injector to parse the remote `authorized_keys` line-by-line and safely inject the local public key via sudo without duplicate entries.
