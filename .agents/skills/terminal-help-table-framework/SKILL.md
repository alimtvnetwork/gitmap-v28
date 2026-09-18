---
name: terminal-help-table-framework
description: >-
  Autonomously design, implement, and audit dry terminal coloring, help rendering, auto-aligned command listings, and middle-ellipsized table display frameworks across GitMap.
---

# Terminal Help & Table Display Framework Skill

Autonomously implement unified terminal help, auto-aligned command layouts, and structured table formatters across GitMap CLI commands adhering to repository coding guidelines.

## Core Capabilities

1. **Dry Terminal Help & Theme System**:
   - Centralize terminal ANSI styling, banners, headers, and subcommand arrow pointers (`->`).
   - Eliminate duplicated command help formatting across CLI packages.

2. **Auto-Aligned Two-Column Command Struct**:
   - Left-column command tokens and right-column descriptions dynamically sized and aligned.
   - Automatic indentation and visual hierarchy.

3. **Middle-Truncated Table Formatter**:
   - Column and row struct abstraction with max width constraints.
   - Intelligent middle ellipsis (`foo...bar`) truncation preserving prefix and suffix context.

4. **Pipeline RCA & Coding Guideline Compliance**:
   - Zero nested ifs (depth <= 1).
   - Functions <= 8-15 lines.
   - Error management via `*AppError` wrappers.
