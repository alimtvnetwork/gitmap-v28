---
name: ent-goal-feasibility-and-email-command-integration
description: Conduct Ent Go ORM feasibility analysis and implement multi-format email command integration with cross-platform lemon service daemon.
---

# ENT Goal Feasibility and Email Command Integration Skill

> **Invocation:** `@[ent-goal-feasibility-and-email-command-integration](file;.agents/skills/ent-goal-feasibility-and-email-command-integration)`

## Mission & Purpose
This skill provides authoritative architectural guidance, code navigation, and execution rules for:
1. **Ent Go (entgo.io) Feasibility & Migration Analysis:** Benchmarking query latency, memory allocation overhead, binary size, compilation duration, and SQLite Split-DB compatibility against GitMap's existing type-safe data access layer.
2. **Email Command Integration:** Multi-format inbound/outbound email parsing (plain text, JSON envelopes, structured action commands), sender authentication, receipt acknowledgment (`ACK Received`), and execution pipeline dispatch inspired by Antigravity Manager / Antigravity Tools.
3. **Cross-Platform Lemon Service Daemon:** Background runner for Windows, Ubuntu (Linux systemd), and macOS (launchd) with configurable polling (5–8s interval, default 8s, min 5s guard), graceful shutdown, and bidirectional email command processing.
4. **Compile-Time Type Safety:** Complete elimination of dynamic SQL string concatenation in favor of parameterized, compile-time typed queries.

---

## Key Architectural Components & Code Map

| Component | Primary Location | Key Responsibilities |
|---|---|---|
| **App Spec & Feasibility Study** | `02-spec/21-app/ent-goal-feasibility-and-email-command-integration/` | Architecture specs, Ent Go benchmark report, email protocol spec. |
| **Execution Plans & Subtasks** | `.ai-memory/plans/ent-goal-feasibility-and-email-command-integration.md` | Master plan and decomposed subtasks in `subtasks/ent-goal-feasibility-and-email-command-integration/`. |
| **Email Subsystem** | `internal/email/` | Message parsers, MIME decoders, SMTP/IMAP adapters, ACK repliers, command envelopes. |
| **Lemon Daemon Service** | `internal/lemon/` or `cli/cmdlemon/` | Background polling loop (5-8s interval), service lifecycle hooks, command execution dispatch. |
| **Database & Type Safety** | `internal/storage/` or `internal/db/` | Split-DB SQLite query implementations, parameterized statements, type-safe repositories. |
| **CLI Commands** | `cli/` / `cmd/gitmap/` | Subcommands for `gitmap lemon`, `gitmap email`, flag bindings, and help system. |
| **End-to-End Tests** | `tests/e2e/` | Isolated integration tests for email parsing, ACK generation, daemon dispatch, and database safety. |

---

## Essential Commands

```bash
# Start the lemon daemon in foreground with 8s polling interval
gitmap lemon run --interval 8s

# Check status of the lemon service
gitmap lemon status

# Send a test command via email subsystem
gitmap email send --to user@example.com --subject "gitmap: status" --body '{"command": "status"}'

# Execute unit and end-to-end test suites
go test -v ./internal/email/...
go test -v ./internal/lemon/...
```
