# Specification: Git Prime Machine OS Version and UI Help Docs Improvement

- **Slug:** `git-prime-machine-os-version-and-ui-help-docs-improvement`
- **Domain:** CLI & Node Identity / Terminal UI / Help Documentation / Merge-AI
- **Status:** APPROVED & ACTIVE
- **Target Release:** `v6.534.0` (Minor Version Bump)

---

## 1. Overview

This specification addresses three critical developer experience and runtime inspection improvements in GitMap:
1. **Accurate OS Version Object & JSON Breakdown:**
   `gitmap machine` and `gitmap machine --json` provide deep, structured operating system detection across Windows, Linux, and macOS. The OS version is represented as a structured object (`OSVersionInfo`) containing operating system type, positive affirmative booleans (`isWindowsServer`, `isWindowsPro`, `isWindowsHome`), version breakdown (`major`, `minor`, `build`), marketing release, edition, and exact composite string.
2. **Help Docs (HD) Web UI Command Catalog Expansion:**
   The Help Dashboard web interface (`gitmap help-dashboard` / `hd`) is updated with complete documentation and navigation for all recently introduced CLI commands: `merge-ai`, `agent-ai`, `heatmap`, `machine`, `full-summary`, `new-commands`, `tree`, `which-format`, and `cpar`.
3. **Parameterized CLI Examples for Merge-AI:**
   `gitmap merge-ai` (`ma`) receives comprehensive CLI help flag routing (`-h`, `--help`, `help`), named flag parsing (`--destination`, `-d`, `--source`, `-s`, `--url`, `--config`), and clear parameterized examples across all 5 input modalities (positional, named flags, comma-delimited, text manifest `sources.txt`, and JSON config `config.json`), documenting the non-destructive `01_`/`02_` collision sequencing and AI consolidation handoff.

---

## 2. Specification Structure

| Document | Purpose |
|---|---|
| [01-architecture-spec.md](01-architecture-spec.md) | High-level system architecture, `OSVersionInfo` struct design, Windows registry/Unix probing, HD docs UI architecture, and invariants |
| [02-component-and-cli-spec.md](02-component-and-cli-spec.md) | Command contracts, CLI flag syntax, JSON schemas, UI component specs, parameterized examples, and Acceptance Criteria (AC-01 – AC-12) |

---

## 3. Related Specifications & References

- `02-spec/21-app/agent-ai-synchronization-and-task-management/`: AI agent task management and Split-DB telemetry.
- `02-spec/21-app/241-gitmap-root-ai-suggestions-and-llm-train-skills/`: AI suggestions and command discovery.
- `02-spec/02-coding-guidelines/03-golang/02-boolean-standards.md`: Affirmative boolean naming rules (`is*` and `has*`).
