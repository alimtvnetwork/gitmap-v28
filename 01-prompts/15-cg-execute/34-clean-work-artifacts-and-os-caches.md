[/goal](slashCommand:goal) Autonomously scan, plan (Dry-Run preview), and purge all unnecessary build artifacts, binaries, and multi-layer caches across the work directory (`d:/work` or `~/work`) and the operating system — including work repo `dist/`/`build/`/`target/`/`tmp/`/`.cache/` and untracked binaries, Go caches (`GOCACHE` & `GOMODCACHE`), npm/pnpm/Yarn/Bun/Node.js caches (while preserving `node_modules` packages), Browser/VS Code/Antigravity DevTools caches, OS temporary directories (`%TEMP%`, `Windows/Temp`, `/tmp`), Windows `SoftwareDistribution/Download`, Recycle Bin / Trash, and Git caches — using `python 03-ai-scripts/44-work-and-system-cache-cleaner.py` and `scripts-fixer` (`clean` / `clear` subcommands).

[/learn](slashCommand:learn) Enforce the Top-Instruction Priority Mandate: whatever directives, target folders, or flags are provided ABOVE this prompt are HIGHEST PRIORITY and strictly NON-NEGOTIABLE. Always run Plan Mode (`--dry-run`) first or use the built-in Plan preview of `03-ai-scripts/44-work-and-system-cache-cleaner.py` before executing with `-y`, and report the exact disk space reclaimed per layer.

> **Top-Instruction Priority Mandate (Above Precedence):**
> Whatever directives, constraints, or user instructions are given ABOVE this prompt are HIGHEST PRIORITY and MUST BE FOLLOWED as strictly NON-NEGOTIABLE.

```text
N = 100 (Total self-loop steps budget)
A = 2   (Number of spawned autonomous subagents running concurrently, default: 2)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
```

---

## 1. Eight-Layer Work & System Cache Cleanup Architecture

`03-ai-scripts/44-work-and-system-cache-cleaner.py` (and `scripts-fixer` `clean` / `clear`) systematically plans and cleans 8 layers across Windows, macOS, and Linux/Unix:

1. **`work-artifacts` — Work Directory Build Elements, Binaries & Repo Caches:**
   - Scans `d:/work` (or `--work-dir <path>`) for build/cache directories (`dist/`, `build/`, `target/`, `.next/`, `.nuxt/`, `.turbo/`, `.parcel-cache/`, `.vite/`, `.cache/`, `__pycache__/`, `.pytest_cache/`, `.mypy_cache/`, `.ruff_cache/`, `tmp/`, `.tmp/`, `.tanstack/tmp/`, `coverage/`, `node_modules/.cache`, `node_modules/.vite`) and untracked binary/build files (`*.exe`, `*.dll`, `*.so`, `*.dylib`, `*.test`, `*.out`, `*.pdb`, `coverage.out`, `tsconfig.tsbuildinfo`, `.eslintcache`).
   - **Strict Preservation:** Keeps `node_modules` installed packages, `.venv`/`venv` virtual environments, and git-tracked files intact.
2. **`go-cache` — Golang Build, Test, Fuzz & Module Caches:**
   - Executes `go clean -cache -testcache -fuzzcache -modcache` and purges `GOCACHE` and `GOMODCACHE` (`dev-tool/go/cache`, `dev-tool/go/pkg/mod`, `~/go/pkg/mod`, `~/.cache/go-build`).
3. **`npm-pnpm-node-cache` — npm, pnpm Store, Yarn, Bun & Node.js Caches:**
   - Executes `npm cache clean --force` and `pnpm store prune`, and purges `npm-cache`, `.pnpm-store`, `dev-tool/pnpm/store`, `node-gyp`, V8 compile caches, Yarn cache, and Bun install cache (`node_modules` folders are kept).
4. **`devtools-cache` — Browser DevTools, VS Code & Antigravity Caches:**
   - Purges Chrome/Edge/Brave DevTools & GPU caches (`Cache`, `Code Cache`, `GPUCache`, `ShaderCache`), VS Code caches (`Cache`, `CachedData`, `CachedExtensionVSIXs`), and Antigravity scratch/crash/temp media caches.
5. **`temp-dirs` — OS & User Temporary Directories:**
   - Purges `%TEMP%`, `%LOCALAPPDATA%/Temp`, `C:/Windows/Temp`, `/tmp`, and `/var/tmp`.
6. **`windows-update` — Windows `SoftwareDistribution/Download` & OS Package Caches:**
   - Purges `%WINDIR%/SoftwareDistribution/Download` and `DeliveryOptimization` on Windows (or `/var/cache/apt/archives` / `Homebrew` cache on Linux/macOS).
7. **`recycle-bin` — System Recycle Bin / Trash:**
   - Empties Windows Recycle Bin (`Clear-RecycleBin -Force` & `$Recycle.Bin`), macOS `~/.Trash`, and Linux `~/.local/share/Trash`.
8. **`git-cache` — Git Caches, `.gitmap` Staging & Stale Pack Temp Files:**
   - Purges `~/.gitcache`, `.gitmap/temp`, `.gitmap/purge`, `.gitmap/downloads`, `.gitmap/sandbox`, and orphaned `.git/objects/pack/tmp_*` files.

---

## 2. Execution Protocol (Plan First -> Confirm or `-y` -> Verify Space Saved)

### Step 1: Run Plan / Dry-Run Preview
Always inspect what will be removed and how much disk space will be reclaimed:
```powershell
python 03-ai-scripts/44-work-and-system-cache-cleaner.py --dry-run
```

### Step 2: Execute Cleanup & Reclaim Space
Run interactively (prompts `[y/N]` after displaying the plan) or automatically with `-y`:
```powershell
# Automatic execution after displaying the Plan table:
python 03-ai-scripts/44-work-and-system-cache-cleaner.py -y

# Or target specific layers only:
python 03-ai-scripts/44-work-and-system-cache-cleaner.py --only work-artifacts,go-cache,npm-pnpm-node-cache,devtools-cache,temp-dirs,windows-update,recycle-bin,git-cache -y
```

### Step 3: `scripts-fixer` Cross-Platform CLI Commands (Windows, macOS, Linux/Unix)
You can also invoke the cleaner directly from `scripts-fixer` on any OS:
```powershell
# Windows (PowerShell):
.\run.ps1 clean artifacts --dry-run
.\run.ps1 clean artifacts -y
.\run.ps1 clean all -y
.\run.ps1 clean --help

# macOS & Linux / Unix (Bash):
./run.sh clean artifacts --dry-run
./run.sh clean artifacts -y
./run.sh clean all -y
./run.sh clean --help
```

---

## 3. Mandatory Output Summary
After running the script, always output the **Final Cleanup Summary Table** showing:
- Category ID & Description
- Space Before Cleanup
- Space After Cleanup (Locked / Active OS files)
- Total Disk Space Saved (`MB` / `GB`) and Total Files Removed.
