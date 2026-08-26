# Runner Contract

**Version:** 1.0.0
**Updated:** 2026-04-19

---

## Purpose

Define the sub-command surface of the **repo-root** runner scripts (`run.sh` and `run.ps1`). These scripts are the user's single entry point for the most common operations.

---

## Sub-command table

| Invocation | Effect | Inner script |
|------------|--------|--------------|
| `./run.ps1` (no positional args) | `git pull` → run the Go coding-guidelines validator on `src/` | `linter-scripts/run.ps1` |
| `./run.ps1 lint [-Path …] [-MaxLines …] [-Json] [-d]` | Explicit lint form. Same effect as no-args. | `linter-scripts/run.ps1` |
| `./run.ps1 slides` | `git pull` → build & preview the slide deck → open browser | inline (see §"Slides sub-command") |
| `./run.ps1 help` | Print the sub-command table and exit 0 | inline |

`run.sh` MUST implement the same surface in Bash with the same positional convention.

---

## Default behavior — preserved

When invoked with **no positional arguments**, the runner MUST behave exactly as before this spec was introduced: forward `-Path`, `-MaxLines`, `-Json`, and `-d` to `linter-scripts/run.ps1` (or `linter-scripts/run.sh`).

> **Why:** Existing users (CI jobs, local muscle memory) rely on `./run.ps1` triggering the Go validator. Adding a sub-command MUST NOT break them.

---

## Slides sub-command

`./run.ps1 slides` and `./run.sh slides` MUST:

1. Print a banner: `▸ slides — building offline deck and opening in browser`.
2. `git pull` (best-effort; warn but continue on failure — the deck may still be buildable from local state).
3. Verify `slides-app/` exists. If not, abort with a clear message pointing to the slides spec.
4. Verify `bun` is available (`bun --version`). If not, fall back to `pnpm`. If neither, abort with installation instructions.
5. Run, with `slides-app/` as the working directory:
   - `bun install --frozen-lockfile || bun install`
   - `bun run build`
   - `bun run preview &` (background)
6. Wait up to 10 seconds for the preview server to be reachable on `http://localhost:4173`.
7. Open `http://localhost:4173/` in the default browser:
   - PowerShell: `Start-Process "http://localhost:4173/"`
   - Bash: `xdg-open` (Linux), `open` (macOS), `start` (Git-Bash on Windows)
8. Print: `▸ slides — preview running. Press Ctrl-C to stop.`
9. `wait` on the preview process so Ctrl-C reaches it.

---

## Help sub-command

`./run.ps1 help` MUST print a table with the same shape as the §"Sub-command table" above and exit 0.

`./run.ps1 -h`, `./run.ps1 --help`, and `./run.ps1 -?` are aliases.

---

## Argument-parsing rule

The first positional argument is the sub-command. If it starts with `-` (e.g. `-Path src/`), treat the invocation as the legacy lint form and forward all arguments unchanged.

```
./run.ps1                       → lint (default)
./run.ps1 -Path cmd             → lint -Path cmd
./run.ps1 lint -Path cmd        → lint -Path cmd
./run.ps1 slides                → slides sub-command
./run.ps1 help                  → print help
./run.ps1 unknown-subcommand    → error: unknown sub-command, exit 2
```

---

## Exit codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Lint validator reported violations OR slides build failed |
| 2 | Unknown sub-command or bad CLI flags |
| 130 | User pressed Ctrl-C |

---

## Anti-requirements

- MUST NOT silently swap the default lint behavior.
- MUST NOT require the user to pre-install `bun` for the lint sub-command (only the slides sub-command needs it).
- MUST NOT auto-update or mutate `slides-app/package.json`.
- MUST NOT leave a background `bun preview` process running after the user exits.

---

## Cross-references

- [`./00-overview.md`](./00-overview.md) — Distribution overview
- [`./01-install-contract.md`](./01-install-contract.md) — Install contract
- [`spec-slides/`](../../spec-slides/) — Slides app spec
- [`spec/13-generic-cli/`](../13-generic-cli/) — Generic CLI conventions

---

*Runner contract — v1.0.0 — 2026-04-19*
