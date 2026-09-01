# Project Type Detection — Overview

## Purpose

During `scan` and `rescan`, gitmap detects project types inside each
discovered Git repository. Detection results are written to dedicated
JSON files and persisted in the SQLite database. Users can query
detected projects by type using dedicated commands.

---

## Supported Project Types

| Type   | Key         | Description                   |
|--------|-------------|-------------------------------|
| Go     | `go`        | Go modules / packages         |
| Node   | `node`      | Node.js projects              |
| React  | `react`     | React applications            |
| C++    | `cpp`       | C/C++ projects                |
| C#     | `csharp`    | .NET / C# projects            |

---

## Spec File Index

| File                          | Contents                          |
|-------------------------------|-----------------------------------|
| `01-overview.md`              | This file — scope and types       |
| `02-detection-rules.md`       | Per-type detection heuristics     |
| `03-data-model.md`            | Tables, relationships, SQL        |
| `04-go-metadata.md`           | Go-specific metadata and runnables|
| `05-csharp-metadata.md`       | C#-specific metadata and key files|
| `06-json-output.md`           | JSON export format and behavior   |
| `07-commands.md`              | CLI commands and flags            |
| `08-scan-integration.md`      | Pipeline integration and flow     |
| `09-package-structure.md`     | New files and modified files      |
| `10-acceptance-criteria.md`   | Acceptance criteria               |
| `11-constants-project.md`     | All constants fully defined       |

---

## Constraints

- All code style rules from `spec/03-general/06-code-style-rules.md`.
- Functions 8–15 lines. Files under 200 lines.
- All string literals in `constants` package.
- Positive conditions only.
- Blank line before `return`.
- PascalCase for DB table/column names.
