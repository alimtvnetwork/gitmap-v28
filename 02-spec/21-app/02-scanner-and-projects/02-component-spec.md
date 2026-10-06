# 02-scanner-and-projects: Scanner Components, Detection Heuristics & Formatter Specification

- **Spec ID:** `02-scanner-and-projects/02-component-spec.md`
- **Status:** `APPROVED`
- **Version:** `1.0.0`
- **Subsystem:** Scanner Components, Polyglot Project Classifier, Bookmarks, Aliases
- **Dependencies:** `cli/scanner`, `cli/detector`, `cli/formatter`, `cli/repodb`
- **Version Baseline:** `v6.498.0`
- **Target Version:** `v6.499.0`

---

## 1. Component Breakdown

The Scanner and Projects cluster comprises four primary operational components:

```
02-scanner-and-projects/
├── 01-architecture-spec.md
└── 02-component-spec.md
```

| Component | Target Source Path | Primary Role |
| :--- | :--- | :--- |
| **Scanner Engine** | `cli/scanner/scanner.go`, `cli/scanner/walk.go` | Concurrent directory tree walking, ignore rules evaluation, `.git` discovery. |
| **Project Detector** | `cli/detector/detector.go`, `cli/detector/rules.go` | Polyglot signature inspection, dependency manifest parsing, language scoring. |
| **Bookmarks & Aliases** | `cli/cmd/bookmarks.go`, `cli/cmd/repo_aliases.go` | Repository bookmark pins, short aliases, fast terminal switching (`gcd`). |
| **Output Formatter** | `cli/formatter/json.go`, `cli/formatter/table.go` | Multi-format rendering (JSON, CSV, Markdown, ANSI Box Table). |

---

## 2. Polyglot Project Detection Engine

### 2.1 Detection Signatures & Scoring
The project detector probes root and 1st-level subdirectory manifests to determine project classification:

```go
package detector

type ProjectType string

const (
    ProjectGo     ProjectType = "go"
    ProjectNode   ProjectType = "node"
    ProjectPython ProjectType = "python"
    ProjectRust   ProjectType = "rust"
    ProjectPHP    ProjectType = "php"
    ProjectDotNet ProjectType = "dotnet"
    ProjectDocker ProjectType = "docker"
    ProjectMixed  ProjectType = "mixed"
)

type DetectionResult struct {
    PrimaryType   ProjectType   `json:"primaryType"`
    DetectedTypes []ProjectType `json:"detectedTypes"`
    ManifestPath  string        `json:"manifestPath"`
    Version       string        `json:"version,omitempty"`
    HasTests      bool          `json:"hasTests"`
    HasDocker     bool          `json:"hasDocker"`
}
```

### 2.2 Detection Heuristics Table

| Project Type | File Signatures Probed | Version Extraction Source |
| :--- | :--- | :--- |
| **Go** | `go.mod`, `go.sum`, `main.go` | Go directive in `go.mod` (e.g. `go 1.24`) |
| **Node.js** | `package.json`, `pnpm-lock.yaml`, `yarn.lock` | `version` field in `package.json` |
| **Python** | `pyproject.toml`, `requirements.txt`, `setup.py` | `pyproject.toml` tool.poetry/project version |
| **Rust** | `Cargo.toml`, `Cargo.lock` | `[package]` version in `Cargo.toml` |
| **PHP** | `composer.json`, `composer.lock` | `composer.json` version string |
| **.NET / C#** | `*.csproj`, `*.sln`, `global.json` | TargetFramework tag |

---

## 3. Bookmarks & Repo Aliases Component

### 3.1 Bookmarks Database Schema
Bookmarks and aliases are persisted in SQLite `installation.db` within the `Bookmark` table:

```sql
CREATE TABLE IF NOT EXISTS Bookmark (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    alias TEXT NOT NULL UNIQUE,
    repoPath TEXT NOT NULL,
    description TEXT,
    isPinned INTEGER NOT NULL DEFAULT 0,
    createdAt DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### 3.2 Operational CLI Commands
- `gitmap bookmark add <alias> [path]`: Binds a short alias to a repository directory.
- `gitmap bookmark list`: Displays all configured bookmarks with branch status.
- `gitmap bookmark remove <alias>`: Deletes an existing bookmark.
- `gitmap cd <alias>` (`gcd <alias>`): Switches working directory instantaneously.

---

## 4. Repository Output Formatters

The formatter layer transforms scanned repository models into multiple presentation formats:
- **`--format=table` (Default):** ANSI colorized terminal table with middle-ellipsized filesystem paths.
- **`--format=json`:** Structured JSON array conformant with `jsonenv.EnvelopeV2`.
- **`--format=csv`:** Comma-separated export including Path, RemoteURL, Branch, CleanStatus, Type.
- **`--format=markdown`:** GitHub Flavored Markdown table for documentation and PR reports.

---

## 5. Verification & Acceptance Criteria

```yaml
verificationGates:
  isPolyglotDetectionVerified: true
  isBookmarkSchemaActive: true
  isMultiFormatOutputSupported: true
  isRelativePathEnforced: true
```

- [x] Accurate language classification across polyglot workspaces.
- [x] Fast bookmark lookup with zero external dependencies.
- [x] Standardized CSV, JSON, Markdown, and ANSI table formatters.
