# Specification: Component & CLI Contracts — Machine OS Version, HD Docs UI & Merge-AI

- **Document ID:** `02-spec/21-app/git-prime-machine-os-version-and-ui-help-docs-improvement/02-component-and-cli-spec.md`
- **Parent Plan:** `.ai-memory/plans/git-prime-machine-os-version-and-ui-help-docs-improvement.md`
- **Domain:** CLI Interface Contracts / OS Identity Schema / Help Docs UI / Merge-AI Modalities
- **Classification:** Core Application Component Specification (`02-spec/21-app`)
- **Status:** APPROVED & RATIFIED
- **Target Release:** `v6.534.0` (Minor Version Bump)

---

## 1. Executive Summary & Scope

This specification details the concrete interface contracts, JSON schemas, terminal rendering rules, web UI component architectures, and acceptance criteria for milestone `git-prime-machine-os-version-and-ui-help-docs-improvement`.

### Addressed Core Deficiencies:
1. **Unstructured Machine OS Version:** `gitmap machine --json` previously exposed a flat, unstructured composite string for `osVersion`. Downstream agents, orchestration scripts, and cluster schedulers could not reliably inspect whether a host is Windows Server vs. Windows Pro, or programmatically query major/minor/build integer numbers.
2. **Missing New Commands in Help Docs UI:** The web dashboard (`src/pages/HelpIndex.tsx` / `gitmap help-dashboard`) lacked documentation for newly integrated commands: `merge-ai`, `agent-ai`, `heatmap`, `machine`, `full-summary`, `new-commands`, `tree`, `which-format`, and `cpar`.
3. **Under-Documented Merge-AI CLI:** `gitmap merge-ai` (`ma`) lacked help flag interception (`-h`, `--help`), named parameter parsing (`--destination`, `--source`), and parameterized multi-modal documentation covering positional arguments, flag lists, comma-delimited strings, text manifests, and JSON configuration files.

---

## 2. Command Contract: `gitmap machine` & `git prime machine`

### 2.1 CLI Grammar & Flag Routing

The machine identity command is invokable via top-level `gitmap machine`, the legacy alias `gitmap alias`, or the unified prime route `git prime machine`.

```bash
gitmap machine [options]
git prime machine [options]
gitmap alias [options]
```

#### Supported Flags:

| Flag | Short | Default | Description |
|---|---|---|---|
| `--json` | `-j` | `false` | Emits structured JSON payload conforming to the `MachineIdentity` schema with nested `osVersionDetails`. |
| `--ssh` | `-s` | `false` | Queries and renders identity across all configured remote SSH cluster nodes in addition to the local host. |
| `--help` | `-h` | `false` | Displays terminal usage, flag descriptions, and structured schema documentation. |

---

### 2.2 Terminal UI Block Formatting (Human ANSI Rendering)

When invoked without `--json`, `gitmap machine` renders a structured, colorized terminal block using PTerm / ANSI escape codes. The operating system section breaks down the host into distinct, readable attributes:

```text
● Machine Identity:
  • Alias:           192.168.1.20
  • Hostname:        AI-Main
  • IP Address:      192.168.1.20
  • OS Platform:     Windows (windows)
  • OS Release:      Windows Server 2022 Standard (21H2)
  • OS Build:        20348.1607 (Major: 10, Minor: 0, Build: 20348, UBR: 1607)
  • OS Edition:      Server Standard (Windows Server: Yes | Pro: No | Home: No)
  • Architecture:    amd64 (20 CPUs)
  • User:            Administrator
  • GitMap Version:  6.534.0
  • Scope:           local
```

On Linux distributions (e.g. Ubuntu):
```text
● Machine Identity:
  • Alias:           192.168.1.50
  • Hostname:        dev-linux-01
  • IP Address:      192.168.1.50
  • OS Platform:     Ubuntu / Linux (linux)
  • OS Release:      Ubuntu 22.04.4 LTS (Jammy Jellyfish)
  • OS Build:        5.15.0-94-generic (Major: 22, Minor: 4, Build: 4)
  • OS Edition:      Server (Windows Server: No | Pro: No | Home: No)
  • Architecture:    amd64 (16 CPUs)
  • User:            ubuntu
  • GitMap Version:  6.534.0
  • Scope:           local
```

---

### 2.3 Structured JSON Schema Contract (`gitmap machine --json`)

When invoked with `--json`, the command emits a JSON array of `MachineIdentity` objects. The `osVersionDetails` object provides deep machine introspection with positive affirmative booleans adhering to `02-spec/02-coding-guidelines/03-golang/02-boolean-standards.md`.

```json
[
  {
    "sequence": 1,
    "nodeId": "local-01",
    "ipAddress": "192.168.1.20",
    "alias": "192.168.1.20",
    "machineName": "AI-Main",
    "osHostname": "AI-Main",
    "osPlatform": "Windows (windows)",
    "osVersion": "Windows Server 2022 Standard 21H2 (Build 20348.1607)",
    "osVersionDetails": {
      "osType": "Windows",
      "marketingName": "Windows Server 2022",
      "edition": "Standard",
      "displayVersion": "21H2",
      "fullVersionString": "Windows Server 2022 Standard 21H2 (Build 20348.1607)",
      "isWindowsServer": true,
      "isWindowsPro": false,
      "isWindowsHome": false,
      "isLinux": false,
      "isMacOS": false,
      "breakdown": {
        "major": 10,
        "minor": 0,
        "build": 20348,
        "revision": 1607
      },
      "kernel": "10.0.20348.1607",
      "installationType": "Server"
    },
    "architecture": "amd64",
    "cpuCores": 20,
    "currentUser": "Administrator",
    "gitmapVersion": "6.534.0",
    "scope": "local"
  }
]
```

#### Go Struct Definitions (`cli/cmdos/os_info_types.go`):

```go
// OSVersionBreakdown represents numerical version integers.
type OSVersionBreakdown struct {
	Major    int `json:"major"`
	Minor    int `json:"minor"`
	Build    int `json:"build"`
	Revision int `json:"revision,omitempty"`
}

// OSVersionInfo provides deep structured operating system introspection.
type OSVersionInfo struct {
	OSType            string             `json:"osType"`
	MarketingName     string             `json:"marketingName"`
	Edition           string             `json:"edition"`
	DisplayVersion    string             `json:"displayVersion"`
	FullVersionString string             `json:"fullVersionString"`
	IsWindowsServer   bool               `json:"isWindowsServer"`
	IsWindowsPro      bool               `json:"isWindowsPro"`
	IsWindowsHome     bool               `json:"isWindowsHome"`
	IsLinux           bool               `json:"isLinux"`
	IsMacOS           bool               `json:"isMacOS"`
	Breakdown         OSVersionBreakdown `json:"breakdown"`
	Kernel            string             `json:"kernel,omitempty"`
	InstallationType  string             `json:"installationType,omitempty"`
}
```

---

## 3. Command Contract: `gitmap merge-ai` (Alias: `ma`)

### 3.1 Invocation Grammar & Help Routing

`gitmap merge-ai` (shorthand `ma`) stages multiple Git repositories into a unified target directory, sequences colliding files into `01_` base and `02_` incoming variants, generates `merge-ai-manifest.json`, and crafts an actionable `instruction.md` consolidation checklist for AI agents.

```bash
gitmap merge-ai [destination] [sources...] [flags]
gitmap ma [destination] [sources...] [flags]
```

#### Help Interception Contract:
Passing any help flag as the first argument, or without arguments, displays the full parameterized help manual with exit code 0:
- `gitmap merge-ai -h`
- `gitmap merge-ai --help`
- `gitmap merge-ai help`
- `gitmap ma --help`
- `gitmap ma` (without arguments: renders guidance and syntax examples)

---

### 3.2 Named Parameter Flags Table

| Flag | Shorthand | Type | Required | Description |
|---|---|---|---|---|
| `--destination <dir>` | `-d <dir>` | string | Yes | Target directory or repository slug where files will be merged. |
| `--source <url\|path>` | `-s <url\|path>` | string[] | Yes | Source repository URL or local path. Can be repeated multiple times. |
| `--url <url>` | `-u <url>` | string[] | No | Shorthand alias for `--source`. Can be repeated multiple times. |
| `--config <file.json>` | `-c <file.json>` | string | No | Path to JSON configuration file specifying destination and sources. |
| `--help` | `-h` | boolean | No | Print comprehensive usage manual and parameterized examples. |

---

### 3.3 The 5 Parameterized Input Modalities

To accommodate diverse developer workflows, automated pipelines, and AI orchestrators, `gitmap merge-ai` supports 5 distinct input modalities:

#### Modality 1: Positional Syntax (Standard CLI)
Target destination followed by space-separated source URLs or directory paths.
```bash
gitmap merge-ai ./unified-app https://github.com/org/repo-auth.git https://github.com/org/repo-billing.git
```
*Alias shorthand:*
```bash
gitmap ma ./unified-app ../repo-auth ../repo-billing ../repo-gateway
```

#### Modality 2: Named Parameter Flags (Explicit Routing)
Explicit `-d` / `--destination` and repeated `-s` / `--source` / `--url` flags.
```bash
gitmap merge-ai --destination ./unified-app --source https://github.com/org/repo-auth.git --source https://github.com/org/repo-billing.git
```
*Short flags:*
```bash
gitmap ma -d ./unified-app -s ../repo-auth -s ../repo-billing -s ../repo-gateway
```

#### Modality 3: Comma-Separated List (Compact Pipe)
Pass multiple sources as a single comma-delimited string argument.
```bash
gitmap merge-ai --destination ./unified-app --source "https://github.com/org/repo-auth.git,https://github.com/org/repo-billing.git,https://github.com/org/repo-gateway.git"
```
*Positional comma variation:*
```bash
gitmap ma ./unified-app "repo-auth,repo-billing,repo-gateway"
```

#### Modality 4: Text Manifest File (`sources.txt`)
Pass a `.txt` file containing newline-delimited repository URLs or paths (ignores comments starting with `#` and empty lines).
```bash
gitmap merge-ai --destination ./unified-app sources.txt
```
*Example `sources.txt`:*
```text
# Core services to merge into unified monorepo
https://github.com/org/service-auth.git
https://github.com/org/service-billing.git
https://github.com/org/service-notifications.git
```

#### Modality 5: JSON Configuration File (`config.json`)
Pass a structured `.json` configuration file via positional argument or `--config`.
```bash
gitmap merge-ai --config merge-config.json
# Or direct positional JSON:
gitmap ma merge-config.json
```
*Example `merge-config.json`:*
```json
{
  "destination": "./unified-app",
  "sources": [
    "https://github.com/org/service-auth.git",
    "https://github.com/org/service-billing.git",
    "https://github.com/org/service-notifications.git"
  ]
}
```

---

### 3.4 File Collision Architecture & Staging Engine

When merging repositories, files with identical relative paths must never overwrite each other silently.

```
Source 1 (Base, Chronologically Earlier)
└── src/config.go  ───────────┐
                              ▼
Source 2 (Incoming, Later)  [ Collision Detector ]
└── src/config.go  ───────────┘
                              │
                              ├──► 01_config.go (Source 1 Base Variant)
                              ├──► 02_config.go (Source 2 Incoming Variant)
                              ├──► merge-ai-manifest.json (Collision Metadata)
                              └──► instruction.md (AI Consolidation Checklist)
```

1. **Chronological Ordering:** Source repositories are sorted by their latest commit timestamp (`LastCommitDate`). Earliest repository is designated Sequence `01` (Base).
2. **Direct Placement:** Files unique to a repository are copied directly into the destination with their canonical relative path.
3. **Collision Sequencing:** When an incoming repository contains a file that already exists in the destination:
   - The existing base file is renamed from `<filename>` to `01_<filename>`.
   - The incoming file is staged as `02_<filename>` (or `03_<filename>` for subsequent collisions).
   - An entry is recorded in `merge-ai-manifest.json` under `fileCollisions`.
4. **Handoff Artifacts:**
   - `merge-ai-manifest.json`: Machine-readable collision registry.
   - `instruction.md`: Human- and AI-readable step-by-step checklist.

---

### 3.5 Generated Manifest Schema (`merge-ai-manifest.json`)

```json
{
  "attributes": {
    "generatedAt": "2026-10-11T19:00:00Z",
    "gitMapVersion": "6.534.0",
    "tool": "gitmap merge-ai",
    "collisionStrategy": "prefixed_sequence_01_02"
  },
  "data": {
    "destination": {
      "targetType": "local_directory",
      "slug": "unified-app",
      "resolvedPath": "./unified-app"
    },
    "repoSequence": [
      {
        "sequenceOrder": 1,
        "repoUrl": "https://github.com/org/repo-auth.git",
        "branch": "main",
        "commitRange": "a1b2c3d..e4f5g6h",
        "lastCommitDate": "2026-09-15T14:30:00Z",
        "description": "feat: auth provider refactor",
        "folderTree": ["src/main.go", "src/auth.go", "README.md"]
      },
      {
        "sequenceOrder": 2,
        "repoUrl": "https://github.com/org/repo-billing.git",
        "branch": "main",
        "commitRange": "b2c3d4e..f5g6h7i",
        "lastCommitDate": "2026-10-01T09:15:00Z",
        "description": "feat: stripe checkout upgrade",
        "folderTree": ["src/main.go", "src/billing.go", "README.md"]
      }
    ],
    "fileCollisions": [
      {
        "canonicalPath": "src/main.go",
        "status": "pending_ai_consolidation",
        "variants": [
          {
            "sequence": "01",
            "file": "src/01_main.go",
            "sourceRepo": "repo-auth",
            "commit": "e4f5g6h"
          },
          {
            "sequence": "02",
            "file": "src/02_main.go",
            "sourceRepo": "repo-billing",
            "commit": "f5g6h7i"
          }
        ]
      },
      {
        "canonicalPath": "README.md",
        "status": "pending_ai_consolidation",
        "variants": [
          {
            "sequence": "01",
            "file": "01_README.md",
            "sourceRepo": "repo-auth",
            "commit": "e4f5g6h"
          },
          {
            "sequence": "02",
            "file": "02_README.md",
            "sourceRepo": "repo-billing",
            "commit": "f5g6h7i"
          }
        ]
      }
    ],
    "uniqueFilesDirectlyPlaced": 2
  }
}
```

---

### 3.6 Generated Instruction Checklist (`instruction.md`)

```markdown
# AI Merge Consolidation Instruction & Verification Guide

This repository has been prepared by `gitmap merge-ai`. Multiple source repositories have been assembled into this staging tree.

## Non-Negotiable Instructions for the Consolidating AI Agent:
1. **Preserve All Business Logic:** Do NOT delete features, endpoints, or data models. Consolidate overlapping files (`01_*`, `02_*`) into the single canonical file name (e.g. synthesize `01_main.go` and `02_main.go` into `main.go`).
2. **Flatten Sequence Prefixes:** Once consolidated, remove temporary `01_` and `02_` files so the repository strictly reflects canonical filenames.
3. **Reference the Manifest:** Open and read `merge-ai-manifest.json` to inspect the source repositories, commit histories, and collision map.
4. **Zero Code Deletion:** When resolving duplicate utilities, favor additive composition, interface unions, and modern patterns.

## Detected File Collisions:
- **src/main.go**:
  - Variant: `src/01_main.go` (Source: `repo-auth`, Commit: `e4f5g6h`)
  - Variant: `src/02_main.go` (Source: `repo-billing`, Commit: `f5g6h7i`)
- **README.md**:
  - Variant: `01_README.md` (Source: `repo-auth`, Commit: `e4f5g6h`)
  - Variant: `02_README.md` (Source: `repo-billing`, Commit: `f5g6h7i`)

## Step-by-Step AI Action Checklist:
- [ ] Read `merge-ai-manifest.json` and inventory all entries in `fileCollisions`.
- [ ] For each collision:
  - [ ] Analyze differences between `01_<file>` and `02_<file>`.
  - [ ] Synthesize unified implementation into `<file>`.
  - [ ] Delete `01_<file>` and `02_<file>` after verified synthesis.
- [ ] Ensure non-colliding files maintain correct imports and package names.
- [ ] Verify that no `01_` or `02_` prefixed files remain in the repository tree.
- [ ] Update `merge-ai-manifest.json` setting `fileCollisions[].status = "consolidated"`.
- [ ] Stage and commit the consolidated repository:
      `gitmap cpf "feat: consolidate multi-repo merge into unified codebase"`
```

---

## 4. Help Docs (HD) Web UI Component Contracts

### 4.1 Command Catalog Additions (`src/data/commands.ts`)

The command catalog in `src/data/commands.ts` must be extended with 9 new/updated `CommandDef` records:

| Command | Category | Alias | Core Purpose |
|---|---|---|---|
| `merge-ai` | `movemerge` | `ma` | Multi-repo amalgamation with non-destructive `01_`/`02_` collision sequencing and AI handoff. |
| `machine` | `tools` | `mach` | Local and SSH fleet machine identity, OS version breakdown, and affirmative booleans. |
| `agent-ai` | `antigravity` | `agm` | AI agent task manager, look-ahead file reservation, and collision detection. |
| `heatmap` | `history` | — | Git workspace commit velocity, file churn, and modification hotspots. |
| `full-summary` | `monitoring` | `fs` | Workspace release summary tree with pipeline error diagnostics (`fspe`). |
| `new-commands` | `tools` | — | Highlights and tours recently added GitMap CLI commands and features. |
| `tree` | `navigation` | — | High-speed ANSI directory tree visualizer with Git status badges. |
| `which-format` | `scanning` | — | Inspects detected output format capabilities (terminal, json, csv, markdown). |
| `cpar` | `cloning` | — | Ultra-fast multi-threaded parallel cloning engine with progress bars. |

#### Complete `CommandDef` for `merge-ai`:

```typescript
{
  category: "movemerge",
  name: "merge-ai",
  alias: "ma",
  description: "Merge multiple Git repositories into a unified destination with automatic collision sequencing (01_/02_) and AI consolidation checklist",
  usage: "gitmap merge-ai <destination> <sources...> [flags] or gitmap ma -d <dest> -s <src1> -s <src2>",
  flags: [
    { flag: "-d, --destination <dir>", description: "Target directory or repository slug for unified output" },
    { flag: "-s, --source <url|path>", description: "Source repository URL or local directory (repeatable)" },
    { flag: "-u, --url <url>", description: "Alias for --source (repeatable)" },
    { flag: "-c, --config <file.json>", description: "Path to JSON config defining destination and sources" },
    { flag: "-h, --help", description: "Show comprehensive help with parameterized examples" },
  ],
  examples: [
    { command: "gitmap merge-ai ./target-repo https://github.com/org/repo-a.git https://github.com/org/repo-b.git", description: "Positional syntax with remote Git URLs" },
    { command: "gitmap ma -d ./target-repo -s ../repo-a -s ../repo-b -s ../repo-c", description: "Named flags syntax with multiple local directories" },
    { command: "gitmap ma --destination ./target-repo --source \"repo-a,repo-b,repo-c\"", description: "Comma-separated source list syntax" },
    { command: "gitmap merge-ai --destination ./target-repo sources.txt", description: "Batch merge using text manifest file" },
    { command: "gitmap merge-ai config.json", description: "Configuration-driven merge using JSON file" },
  ],
  notes: [
    "Files with identical relative paths are non-destructively preserved as 01_<filename> and 02_<filename>.",
    "A merge-ai-manifest.json file is generated in the root of the destination describing all sequenced collisions.",
    "An instruction.md checklist is authored to guide downstream AI agents in synthesizing unified implementations.",
  ],
  seeAlso: [
    { name: "machine", description: "Inspect machine identity and OS version" },
    { name: "full-summary", description: "Inspect active repository statuses" },
    { name: "cpar", description: "Parallel repository cloning" },
  ],
}
```

---

### 4.2 Dedicated Interactive Docs Page (`src/pages/MergeAIPage.tsx`)

A dedicated interactive page is registered at route `/docs/merge-ai` and `/merge-ai`:

#### Component Architecture:
1. **Interactive Parameter Form:**
   - Destination Target input (with suggestions from `projects.json` / Split-DB).
   - Dynamic Source Repository adder (supports URL, relative path, or comma-separated pasting).
   - Modality Selector Tabs: `Positional CLI`, `Named Flags`, `Comma-Separated`, `Text Manifest`, `JSON Config`.
2. **Live Command Generator:**
   - Dynamically renders the exact shell command based on user inputs.
   - One-click copy button with visual feedback.
3. **Collision Simulator Visualizer:**
   - Interactive tree rendering demonstrating how `01_main.go` and `02_main.go` are generated.
   - Interactive preview of `merge-ai-manifest.json` and `instruction.md`.
4. **AI Consolidation Guide:**
   - Markdown step-by-step tutorial on triggering AI agents to consolidate sequenced files.

---

## 5. Acceptance Criteria (AC-01 through AC-12)

| ID | Category | Requirement Specification | Pass Verification Criterion |
|---|---|---|---|
| **AC-01** | OS Identity | Structured OS Version Object in `MachineIdentity` | `gitmap machine --json` output contains `osVersionDetails` object with `osType`, `marketingName`, `edition`, and `fullVersionString`. |
| **AC-02** | OS Identity | Positive Affirmative Booleans | `osVersionDetails` contains `isWindowsServer`, `isWindowsPro`, and `isWindowsHome`. On Windows Server 2022, `isWindowsServer=true`, `isWindowsPro=false`, `isWindowsHome=false`. |
| **AC-03** | OS Identity | SemVer Version Breakdown | `osVersionDetails.breakdown` provides `major` (int), `minor` (int), `build` (int), and optional `revision` (int). |
| **AC-04** | OS Identity | Cross-Platform Parity | On Linux, parses `/etc/os-release` (`isLinux=true`). On macOS, parses `sw_vers` (`isMacOS=true`). |
| **AC-05** | Terminal UI | Human ANSI Block Layout | `gitmap machine` displays structured OS breakdown lines: `OS Platform`, `OS Release`, `OS Build`, `OS Edition`. |
| **AC-06** | Merge-AI CLI | Help Interception | `gitmap merge-ai --help`, `-h`, `help`, and `gitmap ma --help` intercept and print help manual with exit code 0. |
| **AC-07** | Merge-AI CLI | Named Parameter Flags | Supports `-d`/`--destination`, `-s`/`--source`, `-u`/`--url`, `-c`/`--config`. |
| **AC-08** | Merge-AI CLI | 5 Input Modalities | Accepts positional syntax, named flags, comma-separated lists, text files (`sources.txt`), and JSON configs (`config.json`). |
| **AC-09** | Merge-AI CLI | Collision Staging Engine | Non-destructive sequencing creates `01_` base, `02_` incoming variants, `merge-ai-manifest.json`, and `instruction.md`. |
| **AC-10** | HD Docs UI | Command Catalog Expansion | `src/data/commands.ts` includes `merge-ai`, `agent-ai`, `heatmap`, `machine`, `full-summary`, `new-commands`, `tree`, `which-format`, and `cpar`. |
| **AC-11** | HD Docs UI | Interactive `MergeAIPage.tsx` | Dedicated interactive page at `/docs/merge-ai` with parameter builder, live CLI preview, and collision visualizer. |
| **AC-12** | Release Gate | FastGate & Version Bump | Passes FastGate `<1.5s`, zero absolute paths, version bumped to `v6.534.0` across all metadata files. |

---

## 6. Architectural Invariants & Anti-Patterns

1. **Strict Relative Git Paths Only:** All file paths referenced in CLI outputs, JSON schemas, web UI docs, and tests must use forward-slash relative Git paths (e.g. `02-spec/21-app/...`). Absolute drive paths (e.g. `C:\...`, `D:\...`) are strictly prohibited.
2. **Positive Boolean Naming:** All booleans must use positive prefixes (`is*`, `has*`, `can*`). Never emit inverted flags (e.g. `isNotServer`, `disabled`).
3. **Non-Destructive Staging:** `merge-ai` must NEVER overwrite existing files in the destination. Conflicting files MUST be sequenced into `01_` and `02_` prefixes until an explicit consolidation agent runs.
4. **Zero Magic Numbers:** Version numbers must be dynamically queried from `constants.Version` and `version.json`.
