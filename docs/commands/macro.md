# gitmap macro

Record, manage, and replay multi-step command sequences as reusable macros, with support for direct root-level execution.

```bash
gitmap macro <subcommand> [arguments]
```

Aliases: `m`

---

## Direct Root-Level Macro Execution

Any macro saved in GitMap can be run **directly as a root-level command** by typing `gitmap <macro-name>`:

```bash
# Save a macro
gitmap macro add test-all "go test ./..." "npm run test"

# Run it directly as a native command!
gitmap test-all

# Pass flags directly to the macro
gitmap test-all --verbose
gitmap test-all --dry-run
```

---

## Subcommands & Root Utility Commands

| Subcommand | Root-Level Alias | Description |
|---|---|---|
| `gitmap macro list` | `gitmap macro-list`, `gitmap macro-ls` | List all saved macros and step counts (`--json`, `--yaml`) |
| `gitmap macro add <name> <steps...>` | `gitmap macro-add` | Create a new macro directly from command line arguments |
| `gitmap macro run <name>` | `gitmap macro-run`, `gitmap execute` | Replay a macro with timing and step progress |
| `gitmap macro record <name>` | `gitmap macro-record`, `gitmap record` | Record an interactive shell session as a macro |
| `gitmap macro show <name>` | `gitmap macro-show` | Inspect the steps and parameters of a macro |
| `gitmap macro rm <name>` | `gitmap macro-rm` | Delete a saved macro |
| `gitmap macro run-until-succeed <name>` | `gitmap retry`, `gitmap loop` | Retry macro until success with AI diagnostics |
