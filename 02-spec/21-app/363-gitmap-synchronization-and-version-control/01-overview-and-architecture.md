# Spec 363: GitMap Synchronization & Version Control — Overview and Architecture

- **Spec Number:** 363
- **Document:** 01
- **Title:** Overview and Architecture: GitMap as Canonical Root Authority & 3-Tier Asset Provider
- **Slug:** `363-gitmap-synchronization-and-version-control`
- **File Path:** `02-spec/21-app/363-gitmap-synchronization-and-version-control/01-overview-and-architecture.md`
- **Status:** Approved
- **Created Date:** 31-Oct-26
- **Architecture Standards:** GitMap V7 Engine, 3-Tier Asset Provider, Leaf Error Sovereignty, Positive Booleans, <= 15 Line Functions, Zero Absolute Paths

---

## 1. System Intent & Architectural Inversion

### 1.1 Historical Context & Problem Statement
In earlier iterations of the multi-repository ecosystem, the external repository `coding-guidelines` was treated as the nominal source of truth for guidelines, prompts, and skills. However, operational reality evolved:
1. **Tooling Nexus Disconnect**: `gitmap` is the active development engine where CLI tools, high-speed scanners, AST parsers, and AI task orchestrators are engineered, tested, and released.
2. **Brittle Source Discovery**: The existing `ResolveSourceRoot` implementation in `cli/cmdsync/defaults.go` was rigidly coupled to directory names containing `"coding-guidelines"`. When running `gitmap sync` inside the `gitmap` repository or in standalone installations without a local `coding-guidelines` folder, resolution failed or required awkward workarounds.
3. **Distribution Gaps**: Standalone installations of the `gitmap` binary across fleet nodes (e.g., Ubuntu servers, macOS workstations, Windows developer machines) had no native way to seed or update guidelines without cloning the entire source repository.

### 1.2 System Intent: GitMap as Canonical Single Root Authority
This specification inverts the authority hierarchy:
- **GitMap is the Single Canonical Root Authority (`SOURCE_ROOT`)** for all:
  - Core specification standards: `02-spec/01-spec-authoring-guide/` through `02-spec/20-app-spec-audit/`.
  - Canonical AI prompts: `01-prompts/`.
  - Canonical agent skills: `.agents/skills/` and `.cursor/skills/`.
  - Shared AI automation scripts: `03-ai-scripts/` and `.agents/scripts/`.
- **Downstream Fleet Targets**:
  - The `coding-guidelines` repository is transitioned into a downstream target that mirrors GitMap's canonical specifications.
  - All 43+ connected fleet repositories (e.g., `wp-*`, `seo-*`, `alim-*`, `mobile-*`, `desktop-*`) receive synchronized updates originating from GitMap.
- **Strict Boundary Isolation**:
  - Application-specific specs (`02-spec/21-app/*` and higher) are strictly local to each repository and are never overwritten or pushed downstream.
  - Runtime memory (`.ai-memory/`), local configuration (`.env`), and version control metadata (`.git/`) remain strictly isolated per repository.

---

## 2. Architectural Topology

### 2.1 Multi-Tier Distribution Topology Diagram

```mermaid
flowchart TD
    subgraph Authority ["GitMap Canonical Root Authority (SOURCE_ROOT)"]
        direction TB
        GM_Specs["02-spec/01-20 (Canonical Specs)"]
        GM_Prompts["01-prompts/ (Master Prompts)"]
        GM_Skills[".agents/skills/ & .cursor/skills/"]
        GM_Scripts["03-ai-scripts/ & .agents/scripts/"]
        GM_Manifest["version.json (Master SSOT)"]
    end

    subgraph Providers ["3-Tier Asset Provider Engine"]
        direction TB
        Tier1["Tier 1: Live Filesystem Provider\n(Local GitMap Repo)"]
        Tier2["Tier 2: Embedded Go Asset Provider\n(//go:embed Compiled Binary)"]
        Tier3["Tier 3: Downloaded Cache Provider\n(~/.gitmap/cache/assets/vX.Y.Z)"]
    end

    subgraph Engine ["GitMap Sync & 3-Way Hash Engine"]
        direction TB
        SyncCore["cli/cmdsync/sync.go"]
        EvalCore["cli/cmdsync/evaluator.go (3-Way Hash)"]
        Protect["Non-Destructive Human Edit Protection"]
    end

    subgraph DownstreamTargets ["Fleet Repositories (Downstream Targets)"]
        direction TB
        CG_Repo["coding-guidelines (Downstream Spec Mirror)"]
        Fleet1["wp-exam, wp-git-log, wp-onboarding"]
        Fleet2["alim-seo-writing, seo-packages"]
        Fleet3["antigravity-manager, gitmap fleet nodes"]
    end

    Authority --> Providers
    Tier1 --> SyncCore
    Tier2 --> SyncCore
    Tier3 --> SyncCore
    SyncCore --> EvalCore
    EvalCore --> Protect
    Protect --> CG_Repo
    Protect --> Fleet1
    Protect --> Fleet2
    Protect --> Fleet3
```

### 2.2 Synchronization Boundary Matrix

| Asset Category | Path Pattern | Propagation Direction | Sync Policy |
| :--- | :--- | :--- | :--- |
| **Canonical Specs** | `02-spec/01-*` to `02-spec/20-*` | GitMap $\to$ Downstream | 3-Way Hash with Human Edit Protection |
| **Application Specs** | `02-spec/21-*` through `99-*` | Project-Local Only | **NEVER SYNCED** (Excluded) |
| **AI Master Prompts**| `01-prompts/*.md` | GitMap $\to$ Downstream | 3-Way Hash with Human Edit Protection |
| **Agent Skills** | `.agents/skills/*` | GitMap $\to$ Downstream | 3-Way Hash with Sidecar Collision Guard |
| **IDE Skills** | `.cursor/skills/*` | GitMap $\to$ Downstream | 3-Way Hash with Sidecar Collision Guard |
| **Automation Scripts**| `03-ai-scripts/*.py` | GitMap $\to$ Downstream | Add-only / Upstream Update (Non-destructive) |
| **Agent Scripts** | `.agents/scripts/*` | GitMap $\to$ Downstream | Add-only / Upstream Update |
| **Project Memory** | `.ai-memory/*` | Project-Local Only | **NEVER SYNCED** (Excluded) |
| **Local Secrets** | `.env*`, `*.key`, `*.pem` | Project-Local Only | **NEVER SYNCED** (Strictly Excluded) |

---

## 3. Decoupling Source Root Resolution

### 3.1 Deficiency of Legacy Resolution
The legacy `ResolveSourceRoot()` function relied on the following hardcoded check:
```go
// Legacy fragile logic
if err == nil && isCanonicalSourceRepo(cwd) && strings.Contains(strings.ToLower(cwd), "coding-guidelines") {
    return filepath.Clean(cwd), true
}
```
This design was flawed because:
- It prevented `gitmap` itself from acting as the source root.
- It broke if the repository folder was renamed or cloned under another name (e.g., `gitmap-v28`, `cg-repo`).
- It provided zero fallback when executing a compiled binary outside of a source tree.

### 3.2 Autonomous Discovery Algorithm
The decoupled `ResolveSourceRoot()` operates across 5 discrete discovery stages:

```mermaid
flowchart TD
    Start["Call ResolveSourceRoot()"] --> Step1{"Flag or ENV Set?\n(--source-root or GITMAP_SOURCE_ROOT)"}
    Step1 -- Yes --> Validate1["Validate Path has version.json & Canonical Dirs"]
    Validate1 -- Valid --> ReturnCustom["Return Custom Source Root"]
    Validate1 -- Invalid --> Step2
    Step1 -- No --> Step2{"Current Repo is GitMap?\n(Check repo identity & canonical markers)"}
    Step2 -- Yes --> ReturnCwd["Return CWD / RepoRoot as SOURCE_ROOT"]
    Step2 -- No --> Step3{"Find Git Root\n(cmdagent.FindRepoRoot)"}
    Step3 -- GitMap Found --> ReturnGitRoot["Return Git Root as SOURCE_ROOT"]
    Step3 -- Not GitMap --> Step4{"Sibling / Workspace Discovery\n(../gitmap, ../coding-guidelines)"}
    Step4 -- Sibling Found --> ReturnSibling["Return Sibling Path as SOURCE_ROOT"]
    Step4 -- Not Found --> Step5{"Check 3-Tier Fallback\n(Embedded FS or Downloaded Cache)"}
    Step5 -- Available --> ReturnVirtual["Return Virtual Provider (Embedded/Cache)"]
    Step5 -- Missing --> ErrorExit["Return apperror.NewNotFoundError"]
```

#### Canonical Authority Signature Validation
A directory qualifies as an authoritative `SOURCE_ROOT` if and only if:
1. `version.json` exists and its `RepoSlug` equals `"gitmap-v28"` or `Title` equals `"GitMap"`.
2. Required canonical directories exist:
   - `01-prompts/`
   - `02-spec/02-coding-guidelines/`
   - `.agents/skills/`

---

## 4. 3-Tier Asset Provider Architecture

To support offline execution, lightweight binary installs, and developer source trees, GitMap introduces a pluggable, layered asset provider system:

```mermaid
flowchart LR
    Request["Asset Request\n(relPath: '01-prompts/01-unified.md')"] --> Tier1{"Tier 1:\nLive Filesystem\nAvailable?"}
    Tier1 -- Yes --> Serve1["Serve from Live Filesystem\n(fastest, active dev)"]
    Tier1 -- No --> Tier2{"Tier 2:\nEmbedded Assets\nAvailable in Binary?"}
    Tier2 -- Yes --> Serve2["Serve from Go embed.FS\n(zero external dependencies)"]
    Tier2 -- No --> Tier3{"Tier 3:\nCached Remote\nAvailable on Disk?"}
    Tier3 -- Yes --> Serve3["Serve from ~/.gitmap/cache/assets/\n(version-pinned)"]
    Tier3 -- No --> Error["Return Asset Not Found Error"]
```

### 4.1 Tier 1: Live Filesystem Provider (`LiveFsProvider`)
- **Use Case**: Active development within or adjacent to the `gitmap` repository.
- **Characteristics**:
  - Direct file I/O operations via OS filesystem.
  - Instant reflection of changes without requiring binary recompilation.
  - Streaming SHA-256 calculation directly from file descriptors.

### 4.2 Tier 2: Embedded Go Asset Provider (`EmbeddedAssetProvider`)
- **Use Case**: Standalone compiled `gitmap` binaries distributed across fleet servers, CI/CD runners, and developer machines where the GitMap source repository is not cloned.
- **Characteristics**:
  - Compiles canonical prompts, skills, and specs directly into the Go binary using `//go:embed`.
  - Zero filesystem dependencies; cannot be broken by missing sibling directories.
  - Provides immutable baseline assets for seeding new repositories or restoring corrupted assets.

### 4.3 Tier 3: Downloaded Cache Provider (`CachedRemoteProvider`)
- **Use Case**: Slim binary builds (where assets are omitted to reduce binary size) or dynamically fetching updated asset packs from GitHub Releases without replacing the binary.
- **Characteristics**:
  - Downloads release asset archives (`gitmap-assets-vX.Y.Z.tar.gz`) from GitHub.
  - Unpacks into user cache directory (`~/.gitmap/cache/assets/vX.Y.Z/`).
  - Cryptographically verifies the archive against published SHA-256 checksums before unpacking.

### 4.4 Pluggable Go Interface Contract
All asset providers implement a unified interface in `cli/cmdsync/`:

```go
// AssetProvider abstracts asset retrieval across filesystem, embedded binary, and remote cache.
type AssetProvider interface {
	ReadAsset(relPath string) ([]byte, error)
	ListAssets(prefix string) ([]string, error)
	HasAsset(relPath string) bool
	ProviderType() enums.AssetProviderType
}
```

---

## 5. Fleet Synchronization Workflow

When a developer or agent executes `gitmap sync`, the execution lifecycle proceeds through six deterministic phases:

```mermaid
sequenceDiagram
    autonumber
    participant CLI as Developer / Agent (CLI)
    participant Dispatch as cmdsync.RunSync
    participant Resolver as Source & Provider Resolver
    participant Target as Target Project (e.g. wp-exam)
    participant Evaluator as 3-Way Hash State Engine
    participant Output as Terminal & Manifest Writer

    CLI->>Dispatch: gitmap sync [flags]
    Dispatch->>Resolver: ResolveSourceRoot() & ResolveAssetProvider()
    Resolver-->>Dispatch: Active AssetProvider (Tier 1/2/3)
    Dispatch->>Target: Load target version.json & sync manifest
    loop For each canonical asset
        Dispatch->>Evaluator: Compare Baseline (B), Upstream (U), Local (L)
        Evaluator-->>Dispatch: State (CleanUpdate, LocalEdited, Collision, etc.)
        alt Clean Upstream Update
            Dispatch->>Target: Fast-forward write asset (U -> Disk)
        else Human Edit Detected (L != B && U == B)
            Dispatch->>Target: Preserve local edit (No-op)
        else Collision (L != B && U != B)
            Dispatch->>Target: Write sidecar asset (.upstream.md)
        end
    end
    Dispatch->>Target: Update version.json sync manifest with new baseline hashes
    Dispatch->>Output: Render Tokyo Night summary docket & collision report
    Output-->>CLI: Complete with ExitSuccessType
```

---

## 6. Acceptance Criteria & Invariants

1. **Root Authority Invariant**: `gitmap` is the authoritative origin for all coding guidelines (`02-spec/01-20`), master prompts (`01-prompts`), and canonical skills (`.agents/skills`, `.cursor/skills`).
2. **Decoupled Resolution**: `ResolveSourceRoot` must succeed when executed from within `gitmap`, from an external repo with sibling `gitmap`, or via `--source-root` override. Zero reliance on hardcoded `"coding-guidelines"` string.
3. **Multi-Tier Fallback**: When no live filesystem source root exists, `gitmap sync` gracefully falls back to Tier 2 (Embedded) or Tier 3 (Cache) without crashing.
4. **Non-Destructive Guarantee**: Local human edits are protected. Collisions produce sidecars (`.upstream.md`) and never destroy local code.
5. **Guideline Compliance**:
   - Zero absolute paths across all documentation and code.
   - All Go functions $\le$ 15 lines.
   - All booleans strictly follow `is*` or `has*` naming.
   - All errors returned via `apperror.Wrap` or typed app errors.
