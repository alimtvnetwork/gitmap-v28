---
name: macro-schedule-service-suite
description: Autonomously implement and verify macro output padding, async execution, recursion, advanced scheduling with interactive modes, OS service management, OS group/vmware/cron integration, and terminal table visual alignment across GitMap.
---

# Macro, Schedule, Service, OS & Table Alignment Suite

## Core Capabilities

1. **Macro Output Padding & Dynamic Newlines**:
   - Automatically pad every line of command execution output with default 2-space indentation.
   - Smart newline padding: Add a single newline before and after command execution unless output already starts or ends with a newline, preventing double blank lines.

2. **Macro Async Subcommands & Recursion**:
   - Support `async ps "..." -t <sec>`, `async bash "..." -t <sec>`, `async shell "..." -t <sec>`, `async "<cmd>" -t <sec>`.
   - Non-blocking detached execution with background monitoring.
   - Macro recursion support with execution delays.

3. **Schedule Engine Enhancements**:
   - Advanced intervals: `daily (d)`, `weekly (w)`, `startup (s)`, `startup-once (so)`, `startup-weekly (sw)`, `startup-monthly (sm)`, `startup-yearly (sy)`, `yearly (y)`, `every-hour (eh/h)`, `every (e) [n] [hours(h)/daily(d)/weekly(w)/yearly(y)]`.
   - Subcommands: `ls`, `ls <name>`, `status <name>`, `on/off/enable/disable <name>`, `edit`, `rm`, `export`, `import`, `export-all`, `import-all`.
   - Interactive schedule definition reusing macro recording engine.

4. **OS Service Management (`gitmap service`)**:
   - Subcommands: `ls`, `status`, `on`, `off`, `create`, `rm`, `export`, `import`.
   - Native systemd integration on Linux, Service Control Manager on Windows, launchd on macOS.

5. **OS Subcommands Integration**:
   - `gitmap os group`: `add`, `rm`, `ls`.
   - `gitmap os vmware`: shared folder discovery and mounting.
   - `gitmap os cron`: crontab inspection and manipulation.

6. **Terminal Table Visual Alignment**:
   - True visual column width calculation via `runewidth.StringWidth(stripANSI(s))`.
   - Precise space padding for Unicode runes (`✔`, `→`, `●`, `✖`) and ANSI color codes.
