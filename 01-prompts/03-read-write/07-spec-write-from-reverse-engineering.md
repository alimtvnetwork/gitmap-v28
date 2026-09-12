# Reverse Engineering Spec Writer & System Architecture Discovery — Workflow (must follow)

Trigger Keywords & Aliases: `reverse-engineer`, `spec-write`, `reverse-engineering-spec`, `codebase-to-spec`, `generate-app-spec`, `reverse-spec`, `audit-and-spec`

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

```text
N = 80
```

N = total self-loop steps budget that the agents will perform (default: 80 steps for deep scanning, brain memory generation, and complete specification mapping).

/goal Autonomously scan, reverse-engineer, and synthesize an exhaustive, multi-file architectural specification of any target codebase into `spec/21-app/` (or `<spec-folder>/21-app/`), detecting programming languages, isolating architectural boundaries, analyzing security vulnerabilities and risk tiers, and structuring all findings with strictly lowercase naming and relative git paths without stopping until 100% green.

---

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Step 0 (Dynamic Environment Bootstrap): Detect if `.lovable/` and `.lovable/temp/` exist; if missing, dynamically create `.lovable/temp/` to hold scratch files and inventories.
2. [ ] /goal Step 0 (Spec Folder Initialization): Detect or create the target spec folder (defaulting to `spec/21-app/` or `<XX>-spec/21-app/`) adhering to two-digit monotonic numbering and strictly lowercase directory conventions.
3. [ ] /goal Phase 0 (Antigravity Skill Bootstrap): Check if `.agents/skills/spec-reverse-engineering/skill.md` exists; if missing, auto-scaffold the skill with YAML frontmatter for progressive memory disclosure.
4. [ ] /goal Step 1 (Tooling Script Creation): Create or verify the Python inventory tool (`03-ai-scripts/<NN>-codebase-file-lister.py` or `.lovable/<temp>/list_files.py`) to crawl all repository files while excluding `.git`, `node_modules`, `vendor`, and build artifacts.
5. [ ] /goal Step 2 (Codebase Crawl & Language Topology): Execute the inventory script to generate `.lovable/<temp>/files-inventory.json` with relative paths, file sizes, and language breakdown percentages.
6. [ ] /goal Phase 1 (Planning & File Management Ledger, Steps 1..N/2): Partition discovered files into balanced batches across 2–3 concurrent sub-agents, recording active locks in `.lovable/<temp>/file-assignments.json`.
7. [ ] /goal Phase 1 (Micro-Tasking Subtasks): Decompose file analysis into granular subtask files under `.lovable/plans/subtasks/<reverse-engineering>/` tracking which sub-agent processes which file paths.
8. [ ] /goal Phase 1 (Zero-Stop Transition): Immediately upon completing file inventory and assignment planning, self-loop and transition directly into Phase 2 execution mode without pausing or requesting user input.
9. [ ] /goal Phase 2 (Parallel File Reverse Engineering, Steps N/2+1..N): Dispatch 2–3 execution sub-agents in parallel on disjoint file sets to analyze code semantics, exported types, data flow, functions, and external dependencies.
10. [ ] /goal Phase 2 (Component Specification Generation): Write dedicated, modular specification files under `spec/21-app/` (e.g., `spec/21-app/XX-core-engine.md`, `spec/21-app/XX-data-models.md`, `spec/21-app/XX-api-contracts.md`) documenting all reverse-engineered logic.
11. [ ] /goal Phase 3 (Master Index Synthesis): Author `spec/21-app/01-index.md` summarizing the overall application architecture, behavior, technology stack, architectural health score, and component topology.
12. [ ] /goal Phase 3 (Security Audit & Risk Assessment): Identify hardcoded credentials, unauthenticated endpoints, input sanitization flaws, and dependency vulnerabilities, publishing an exhaustive risk evaluation in `spec/21-app/xx-security-and-risks.md`.
13. [ ] /goal Phase 3 (Final Structure Communication): Output a clean, viewable markdown/ASCII folder tree in the final chat response illustrating the complete generated specification layout.
14. [ ] /learn Ingest `.lovable/memory/01-index.md` for project memory index and past learnings.
15. [ ] /learn Ingest `.lovable/strictly-avoid.md` for banned anti-patterns and strict constraints.
16. [ ] /learn Ingest `spec/02-coding-guidelines/08-file-folder-naming/` for lowercase naming and continuous file sequencing.
17. [ ] /learn Ingest `spec/02-coding-guidelines/01-cross-language/01-index.md` for strict relative path citation requirements.
18. [ ] /learn Ingest `spec/21-app/01-index.md` for baseline application documentation standards.
19. [ ] /goal Verify zero absolute paths or `file:///` URIs exist in generated specification markdown.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Environment Bootstrap, Inventory Script, File Discovery, Language Topology, Concurrency Ledger)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: 2–3 Agent Parallel File Analysis, Module Specs, 01-index.md Synthesis, Security Audit, Tree Output)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## Dedicated Section: Dynamic Environment & Spec Folder Conventions

> [!IMPORTANT]
> **DYNAMIC ADAPTATION TO ANY CODEBASE:**
> This workflow is designed to execute against any target codebase, including legacy projects or repositories lacking predefined directory structures. The AI must dynamically adapt:

1. **`.lovable/` Directory Handling:**
   - If the target workspace lacks a `.lovable/` folder, create `.lovable/` and `.lovable/temp/`.
   - All ephemeral file inventories, crawler outputs, and subagent file locks MUST reside in `.lovable/temp/`.
2. **Spec Directory Conventions:**
   - Locate the primary specification folder (e.g., `spec/`, `<spec>/`, or create `spec/` if none exists).
   - Create or populate the application-specific spec folder: `<spec-root>/21-app/` (e.g. `spec/21-app/`).
   - If sibling folders for issues exist or are needed, reference `<spec-root>/22-app-issues/` or create `spec/21-app/xx-security-and-risks.md`.
3. **Strict Lowercase File Naming:**
   - All files created MUST use strictly lowercase characters, hyphens, and numeric prefixes (e.g., `01-index.md`, `xx-security-and-risks.md`, `XX-core-architecture.md`). Uppercase letters are strictly banned.
4. **Strict Relative Git Paths:**
   - All markdown links and citations MUST be relative paths from repository root (e.g., `spec/21-app/01-index.md`). NEVER write absolute paths (`C:\...`, `/home/...`) or `file:///` URIs.

---

## Phase 0: Antigravity Skill Bootstrap (Progressive Memory Disclosure)

Before crawling the codebase, the AI must check if this workflow is registered as a native Antigravity Skill.

1. Inspect `.agents/skills/spec-reverse-engineering/skill.md`.
2. If missing, create the skill directory and file with YAML frontmatter:
   ```yaml
   ---
   name: spec-reverse-engineering
   description: Autonomously scan, reverse-engineer, and synthesize comprehensive architectural specs and security audits from any codebase.
   ---
   ```
3. Copy the core execution instructions into the skill file to enable progressive disclosure for future execution sessions.

---

## Step 0 & 1: Inventory Tool Creation & Polyglot Discovery

To prevent context-window exhaustion on large repositories, the AI must NOT attempt to list files via raw recursive shell output. Instead, an automated Python crawler must be generated and executed:

1. **Tooling Script Location:**
   - In standard meta-repositories: `03-ai-scripts/<NN>-codebase-file-lister.py`.
   - In bare workspaces: `.lovable/<temp>/list_files.py`.
2. **Crawler Capabilities:**
   - Scans directory tree using non-blocking traversal.
   - Strictly ignores `.git/`, `node_modules/`, `vendor/`, `dist/`, `build/`, `bin/`, `.venv/`, and binary file extensions (`.png`, `.jpg`, `.exe`, `.tar`, `.zip`, `.dll`, `.so`).
   - Counts lines of code (LOC) and detects primary programming languages based on file extensions (`.go`, `.ts`, `.tsx`, `.py`, `.rs`, `.cs`, `.php`, `.java`, `.c`, `.cpp`, `.sql`, etc.).
   - Outputs JSON payload to `.lovable/<temp>/files-inventory.json`:
     ```json
     {
       "total_files": 42,
       "languages": {
         "Go": 28,
         "Python": 10,
         "TypeScript": 4
       },
       "files": [
         {"path": "pkg/auth/token.go", "language": "Go", "lines": 120, "size_bytes": 3450}
       ]
     }
     ```

---

## Phase 1: Planning, File Management System & Concurrency Ledger (Steps 1 .. N/2)

The AI orchestrator must establish a rigorous File Management System before assigning tasks to sub-agents:

1. **Concurrency Ledger (`.lovable/<temp>/file-assignments.json`):**
   - Partition the discovered file list into 2 or 3 balanced disjoint batches.
   - Structure the assignment ledger:
     ```json
     {
       "agent_1": {
         "status": "in_progress",
         "assigned_files": ["pkg/api/router.go", "pkg/api/handlers.go"]
       },
       "agent_2": {
         "status": "in_progress",
         "assigned_files": ["pkg/models/user.go", "pkg/models/session.go"]
       },
       "agent_3": {
         "status": "in_progress",
         "assigned_files": ["pkg/storage/db.go", "pkg/storage/queries.go"]
       }
     }
     ```
2. **Zero Collision Guarantee:**
   - No file may ever be assigned to more than one sub-agent concurrently.
   - When a subagent completes analysis of a file, it marks the file as `"completed"` in the ledger.
3. **Subtask Decomposition:**
   - Write microscopic subtasks in `.lovable/plans/subtasks/<reverse-engineering>/` defining file boundaries and expected specification sections.
4. **Mandatory Auto-Loop:**
   - Transition immediately into Phase 2 execution mode without prompting the user.

---

## Phase 2: Parallel Reverse Engineering & Spec Authoring (Steps N/2+1 .. N)

Spawn 2 to 3 sub-agents concurrently (max 2 threads each) to process the assigned batches:

1. **File-Level Code Analysis:**
   - For every source file, extract:
     - **Module Identity:** Primary responsibility and role in the system.
     - **Key Types & Contracts:** Core structs, interfaces, enums, classes, and signatures.
     - **Data Flow & Algorithms:** Input sources, mutation steps, validation routines, and output transformations.
     - **Dependencies:** Internal package imports and external third-party libraries.
     - **Error Handling:** Fault paths, error types returned, and recovery mechanisms.
2. **Modular Spec Authoring (`spec/21-app/`):**
   - Group related files into coherent domain specifications:
     - `spec/21-app/xx-domain-models.md`
     - `spec/21-app/xx-service-layer.md`
     - `spec/21-app/xx-api-endpoints-and-protocols.md`
     - `spec/21-app/XX-database-and-persistence.md`
3. **Micro-Tasking Execution:**
   - Subagents must process 3 to 5 files per turn, updating their status and self-looping until their entire partition is documented.

---

## Phase 3: System Overview, Security Audit & Tree Output

Upon completion of all modular specifications, the master orchestrator synthesizes the top-level documentation:

### 1. Master Specification Index (`spec/21-app/01-index.md`)
The index file must provide an executive synthesis:
- **Application Overview:** Purpose of the software, high-level architecture, user flows, and core features.
- **Language & Stack Matrix:** Exhaustive breakdown of languages, frameworks, runtimes, and libraries discovered.
- **Architectural Health Score:** Objective rating (1–10) assessing code readability, separation of concerns, modularity, test coverage presence, and adherence to clean architecture principles.
- **System Topography Diagram:** ASCII or Mermaid diagram illustrating component interactions and data flow.
- **Specification Directory Index:** Complete table of contents linking to all generated module specs.

### 2. Security Audit & Risk Evaluation (`spec/21-app/xx-security-and-risks.md`)
Every reverse-engineered codebase must undergo a thorough security inspection:
- **Hardcoded Secrets & Sensitive Data:** Scan for embedded API keys, JWT secrets, passwords, private keys, or exposed test tokens.
- **Authentication & Authorization Posture:** Analyze session handling, token validation, permission checks, and privilege escalation risks.
- **Input Validation & Injection Vectors:** Check for raw SQL queries, command execution, path traversal, and unescaped user inputs.
- **Third-Party Dependency Vulnerabilities:** Flag deprecated or high-risk libraries identified in package manifests.
- **Project Risk Level:** Categorize overall risk as **LOW**, **MEDIUM**, **HIGH**, or **CRITICAL**, with immediate remediation recommendations.

### 3. Visual Tree Communication in Chat Response
In the final turn after all files are generated, the orchestrator MUST output a clear, viewable ASCII folder tree in chat:

```text
📁 spec/21-app/
├── 📄 01-index.md                     # Application overview, tech stack, health score & architecture
├── 📄 02-security-and-risks.md        # Security audit, flaw analysis, secrets scan & risk rating
├── 📄 03-domain-models.md             # Data structures, enums, interfaces, and entities
├── 📄 04-business-logic.md            # Service layer operations, workflows, and state machines
├── 📄 05-api-and-transports.md        # HTTP/gRPC/CLI endpoints, routing, and serialization
└── 📄 06-persistence-and-storage.md   # Database schemas, migrations, caching, and filesystem I/O
```

---

## Strict Discipline & Anti-Hallucination Mandates

- **NEVER Guess Implementation Details:** If a function or module is obfuscated or external, inspect its callers or mark it explicitly as external.
- **NEVER Attempt Everything in One Turn:** Enforce strict sequential self-looping across the $N=80$ budget.
- **Strict Relative Git Paths:** TOTAL BAN on absolute filesystem paths or `file:///` URIs.
- **Zero CI/CD Regressions:** Generated documentation must comply with Markdown gap linters, sequence integrity checkers, and doc path validators.
