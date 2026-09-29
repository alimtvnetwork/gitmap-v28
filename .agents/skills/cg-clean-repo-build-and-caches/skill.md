---
name: cg-clean-repo-build-and-caches
description: Autonomously clean repository build issues, purge temporary directories, OS temp folders, Golang build/test caches, pnpm/npm stores, Vite/TS incremental caches, and polyglot caches while enforcing lowercase hygiene and .gitignore sync.
---

[/goal](slashCommand:goal) Autonomously diagnose and clean the repository to resolve build issues, purge stale build artifacts, and systematically clean all temporary directories and multi-language caches — including repository temp folders, OS temp directories, Golang build/test caches, pnpm/npm/Vite caches, Python/Rust/PHP caches, and untracked artifact pollution — with strict no-build and no-test execution (NEVER run heavy build commands like `go build` or `npm run build`, and NEVER run test suites like `go test ./...` or `pytest` during routine turns; compilation and testing are verified later in CI/CD). Spawn autonomous subagents (A = 2, H = 2) for parallel inspection and cache purging, leverage GitMap high-speed tooling as primary, synchronize `.gitignore` rules, and finalize any repository hygiene fixes with an atomic push.

[/learn](slashCommand:learn) Enforce the Bottom-Instruction Priority Mandate: whatever directives, custom cache targets, build issue symptoms, or user instructions are appended BELOW this prompt (following the `--` divider border at the bottom) are HIGHEST PRIORITY and strictly NON-NEGOTIABLE. Execute the multi-layer cache and temp cleanup protocol across repository, OS, Go, pnpm/npm, and polyglot toolchains, and persist hygiene findings into `.ai-memory/plans/`.

## 1. When to Use

Activate this skill when:
- The user asks to clean the repository for build issues, clear temporary files, or purge OS, Golang, pnpm/npm, or polyglot caches.
- Executing `01-prompts/v4/15-cg-execute/30-clean-repo-build-and-caches.md`.

## 2. Multi-Layer Cache & Temp Cleanup Protocol

1. **Layer 1 — Repository Build Artifacts:**
   - Remove `.tmp/`, `tmp/`, `__pycache__/`, `.pytest_cache/`, `.mypy_cache/`, `.ruff_cache/`, `.vite/`, `.turbo/`, `.parcel-cache/`, `*.pyc`, `*.pyo`, `*.tmp`, `*.test`, `coverage.out`, `tsconfig.tsbuildinfo`, `.eslintcache`.
   - Run `python 03-ai-scripts/19-artifact-remover.py --execute` (if present).
2. **Layer 2 — OS Temporary Directories:**
   - Clean stale `$env:TEMP\go-build*`, `$env:TEMP\vite*`, `$env:TEMP\npm-*`, `/tmp/go-build*` directories.
   - Run `gitmap update-cleanup`.
3. **Layer 3 — Golang Caches:**
   - Execute `go clean -cache -testcache -fuzzcache` and `golangci-lint cache clean`.
4. **Layer 4 — pnpm, npm & Frontend Caches:**
   - Execute `pnpm store prune`, `npm cache clean --force`, and remove `node_modules/.cache` and `node_modules/.vite`.
5. **Layer 5 — Git Hygiene & Lowercase Enforcement:**
   - Run `gitmap commons` to sync `.gitignore` and `.gitattributes`.
   - Run `gitmap lcf` and `gitmap lowercase-readme` to eliminate case-sensitivity build blockers.
   - Run `gitmap storage` (`gitmap stor`) to verify disk space recovery.
   - Commit any repository hygiene fixes atomically via `gitmap cpb`.
