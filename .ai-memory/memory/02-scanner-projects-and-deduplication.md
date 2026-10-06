# 02 — Scanner Engine, Project Discovery & Zero-Alloc Deduplication

- **Domain:** Repository Scanning, Polyglot Project Heuristics & Deduplication
- **Authoritative Specification:** [02-scanner-and-projects](../../02-spec/21-app/02-scanner-and-projects/01-architecture-spec.md)
- **Status:** Active & Ratified

---

## 1. Scanner Traversal Architecture

The GitMap scanner traverses deep filesystem directories at ultra-high speed to discover Git repositories and managed codebases:

- **Parallel Directory Walker:** Uses goroutine worker pools with bounded concurrency to walk file trees across local disks and mounted network shares.
- **Git Root Identification:** Directories containing a valid `.git/` directory (or git worktree pointer file) are indexed immediately as repository roots; traversal does not descend into `.git/` subdirectories.
- **Symlink Protection:** Follows symlinks only when explicitly configured, utilizing circular reference detection maps (`map[string]bool`) to prevent infinite recursion loops.
- **OS Path Normalization:** Paths are converted to forward slashes internally on all operating systems, with case-folding collation for Windows and macOS filesystems.

---

## 2. Polyglot Project Detection Heuristics

Once a Git repository is identified, the scanner detects internal project stacks and frameworks via signature heuristic probes:

| Stack / Language | Primary Signature File | Secondary Framework Indicators |
| :--- | :--- | :--- |
| **Go** | `go.mod` | `go.sum`, `main.go`, `cmd/` |
| **Python** | `pyproject.toml`, `requirements.txt` | `setup.py`, `Pipfile`, `venv/`, `.venv/` |
| **TypeScript / Node** | `package.json` | `tsconfig.json`, `pnpm-lock.yaml`, `vite.config.ts` |
| **Rust** | `Cargo.toml` | `Cargo.lock`, `src/main.rs`, `src/lib.rs` |
| **PHP** | `composer.json` | `artisan`, `composer.lock` |
| **C# / .NET** | `*.csproj`, `*.sln` | `Program.cs`, `appsettings.json` |

- **Multi-Project Workspaces:** Repositories containing monorepos or multiple sub-projects record parent-child relationships in the database, allowing granular sub-project targeting.

---

## 3. Fast File Reader & Binary Hygiene

- **Fast Streaming Buffer:** The scanner inspects file headers without reading full files into memory, checking magic bytes to classify binary vs text files.
- **Size Guards:** Files exceeding 1MB are automatically bypassed during content searches to prevent out-of-memory spikes.
- **`.gitmapignore` Specification:** Supports `.gitmapignore` files placed at repository or workspace roots, adhering to standard `.gitignore` syntax to skip build caches (`dist/`, `target/`, `node_modules/`, `vendor/`).

---

## 4. Zero-Allocation Deduplication & Database Cache

- **Fast Hashing:** Generates 64-bit non-cryptographic hashes for path strings and folder structures to detect duplicates with zero memory allocations.
- **Split-DB Cache Lifecycle:**
  - Discovered repositories are cached in `gitmap.db` within the `Repository` and `Project` tables.
  - Periodic scans execute incremental delta syncs using file modification timestamps (`mtime`), updating only altered repositories.
  - Stale repositories (directories deleted from disk) are marked inactive or pruned based on user retention flags (`--prune`).
