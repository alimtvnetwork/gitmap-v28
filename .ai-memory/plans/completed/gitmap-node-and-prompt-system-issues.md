# Execution Plan: GitMap Node & Prompt System Issues Remediation

**Task Slug:** `gitmap-node-and-prompt-system-issues`  
**Parent Task ID:** `task-20261019063719-gitmap-node-and-prompt-system-issues`  
**Task DB:** `.ai-memory/temp-agents/gitmap-node-and-prompt-system-issues/agent-task.db`  
**Status:** Completed  
**Spec Reference:** `02-spec/21-app/gitmap-node-and-prompt-system-issues/spec.md`

---

## 1. Subtask Breakdown

| Subtask ID | Title | Scope & Owned Files | Target Outcome |
| :--- | :--- | :--- | :--- |
| **subtask-01** | Node Fleet Telemetry & SSH Ping | `cli/cmdnodes/`, `cli/cmdssh/`, `cli/cmdui/ui_server_ssh.go`, `cli/cmdui/ui_assets_js_nodes.go`, `cli/cmdui/ui_types.go`, `cli/cmdui/ui_assets_markup_ops.go` | Display GitHub node aliasing, account mapping, instance count & list across CLI & UI; fix Web UI SSH ping (`testNode`) |
| **subtask-02** | Prompt Tree Suite & Telegram Priority | `cli/cmdagy/agy_ui.go`, `cli/cmdagy/agy_ui_html.go`, `cli/cmdagy/agy_telegram_subcommands.go`, `cli/cmdui/ui_assets_js_prompts_help.go`, `cli/cmdui/ui_assets_markup_ops.go` | Interactive visual prompt tree with Send/Backup/Restore/Reset/Refresh buttons, Telegram machine identification, running & queued prompts pinned to top |
| **subtask-03** | Supabase Backoff, Throttling & Header Details | `cli/cmdsupabase/`, `cli/cmdagy/agy_account_switch_engine.go`, `cli/cmdagy/agy_conv_scanner.go`, `cli/cmdprompttemplate/template_import_export.go`, `cli/cmdui/` | 5-minute Supabase error backoff, 5m UI auto-refresh, $\ge 30$s credit sync, Prompt Preview header details collapse, CLM cache, and template import/export re-indexing |
| **subtask-04** | Isolated Testing, FastGate & Minor Release | Tests, `03-ai-scripts/50-fastgate.py`, `03-ai-scripts/37-bump-version.py` | Verify test suite with isolated test instance (never close `default-copy-814`), execute minor bump release ceremony, monitor CI/CD green |

---

## 2. Dependency Graph & Worker Concurrency (Tree Swarm A = 3)
```
          [Subtask 01: Node & SSH Telemetry] (Worker 01)
         /
[Plan] -- [Subtask 02: Prompt Tree & Telegram] (Worker 02)
         \
          [Subtask 03: Supabase & Header Details] (Worker 03)
                           |
                           v
          [Subtask 04: Testing, FastGate & Minor Bump] (Lead Agent)
```

---

## 3. Worker Safety Invariants
1. **Total Ban on Git in Subagents:** Workers must NEVER execute `git` commands.
2. **Total Ban on Full Builds / Suites:** NEVER execute `go build` or `go test ./...`. Targeted file tests only.
3. **Total Ban on Raw Search Tools:** Total ban on `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, and `findstr`. Use `gitmap aum search` or `gitmap find` exclusively.
4. **Strict Relative Git Paths:** Only relative paths starting from repo root.
5. **Active Process Preservation:** Never kill or close `default-copy-814`.
