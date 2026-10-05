# Subtask Plan 01: Oh-My-Zsh & Default Directory Scanner Exclusions with Force Override

- **Subtask Slug:** `01-omz-scanner-exclusions-and-force-flag`
- **Parent Task:** `219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging`
- **Target Files:**
  - `cli/constants/constants_scan.go`
  - `cli/scanner/scanner.go`
  - `cli/model/record.go`
  - `cli/fsutil/child_repos.go`
  - `cli/fsutil/recursive_top_level.go`
  - `cli/cmdscan/flags.go`
  - `cli/cmdscan/scan.go`
  - `cli/cmd/rootadd.go`

---

## 1. Overview & Objectives

On Linux/Ubuntu workstations and servers, scanning user home directories (`~` or `/home/<user>`) inadvertently descends into Oh-My-Zsh configurations (`~/.oh-my-zsh`), vendor directories, and virtual environments, indexing dozens of internal zsh plugins and themes into `gitmap.db`.

This subtask implements:
1. Centralized baseline exclusion list in `cli/constants/constants_scan.go` covering Oh-My-Zsh variants (`.oh-my-zsh`, `ohmyzsh`, `omz`, `.omz`, `omizssh`) and common package/tool directories (`node_modules`, `vendor`, `.cache`, `.venv`, `.git`).
2. Automatic merging of baseline exclusions into `cli/scanner/scanner.go:buildExcludeSet` and `cli/model/record.go:DefaultConfig`.
3. Alignment of filesystem utilities in `cli/fsutil/child_repos.go` and `cli/fsutil/recursive_top_level.go` to respect baseline exclusions.
4. Force override flags: `--force-include` on `gitmap scan` and `--force` on `gitmap add` allowing intentional inclusion of excluded directories.

---

## 2. Step-by-Step Implementation Instructions

### Step 1: Define Constants in `cli/constants/constants_scan.go`
1. Create or update `cli/constants/constants_scan.go`:
   ```go
   package constants

   // DefaultScanExcludeDirs defines canonical directories automatically ignored by scanner and discovery.
   var DefaultScanExcludeDirs = []string{
       ".oh-my-zsh",
       "oh-my-zsh",
       "ohmyzsh",
       ".ohmyzsh",
       "omz",
       ".omz",
       "omizssh",
       "node_modules",
       "vendor",
       ".cache",
       ".venv",
       "venv",
       ".git",
       ".terraform",
       ".next",
       ".turbo",
       "dist",
       "build",
       "bin",
       "obj",
       "target",
   }

   const (
       FlagScanForceInclude      = "force-include"
       FlagScanForceIncludeAlias = "fi"
       FlagDescScanForceInclude  = "Comma-separated directories to include despite default exclusions (or 'all' to include everything)"
       FlagAddForce              = "force"
       FlagAddForceAlias         = "f"
       FlagDescAddForce          = "Force tracking of directory even if matching default exclusions"
   )
   ```

### Step 2: Update Scanner Engine in `cli/scanner/scanner.go`
1. Add `ForceIncludeDirs []string` to `ScanOptions` struct:
   ```go
   type ScanOptions struct {
       ExcludeDirs      []string
       ForceIncludeDirs []string
       Workers          int
       Progress         func(ScanProgress)
       MaxDepth         int
       OnDirError       func(path string, err error)
   }
   ```
2. Refactor `buildExcludeSet` to take both `excludeDirs` and `forceIncludeDirs`:
   ```go
   func buildExcludeSet(customExcludes []string, forceIncludes []string) map[string]bool {
       // Check if user requested to bypass all exclusions
       for _, fi := range forceIncludes {
           if fi == "all" || fi == "*" {
               return make(map[string]bool)
           }
       }

       forceSet := make(map[string]bool, len(forceIncludes))
       for _, fi := range forceIncludes {
           forceSet[fi] = true
           forceSet[strings.ToLower(fi)] = true
       }

       set := make(map[string]bool)
       // 1. Add baseline defaults unless explicitly force-included
       for _, d := range constants.DefaultScanExcludeDirs {
           if !forceSet[d] && !forceSet[strings.ToLower(d)] {
               set[d] = true
           }
       }
       // 2. Add caller custom excludes
       for _, d := range customExcludes {
           if !forceSet[d] && !forceSet[strings.ToLower(d)] {
               set[d] = true
           }
       }

       return set
   }
   ```
3. Update `ScanDirWithOptions` to pass `buildExcludeSet(opts.ExcludeDirs, opts.ForceIncludeDirs)` into `walkParallel`.

### Step 3: Align `cli/model/record.go`
1. Update `DefaultConfig()` in `cli/model/record.go` so `ExcludeDirs` defaults to a copy of `constants.DefaultScanExcludeDirs`.

### Step 4: Align Filesystem Utilities
1. **`cli/fsutil/child_repos.go`**:
   - In `DiscoverChildGitRepos(parentDir string)`:
     * Before checking if an entry contains `.git`, verify that `entry.Name()` is not in `constants.DefaultScanExcludeDirs`.
     * Skip excluded subtrees to prevent scanning inside `.oh-my-zsh/plugins/` or `node_modules/`.
2. **`cli/fsutil/recursive_top_level.go`**:
   - In `DiscoverTopLevelGitRepos(rootDir string)`:
     * Replace `name == ".git" || name == "node_modules" || name == ".cache" || name == "vendor"` with a helper check matching `constants.DefaultScanExcludeDirs`.

### Step 5: Wire CLI Flags in `cli/cmdscan/` and `cli/cmd/`
1. In `cli/cmdscan/flags.go`:
   - Add `forceIncludeFlag *string` to `scanFlagPointers`.
   - Register in `registerScanStringFlags`:
     ```go
     flagPtrs.forceIncludeFlag = fs.String(constants.FlagScanForceInclude, "", constants.FlagDescScanForceInclude)
     ```
2. In `cli/cmdscan/scan.go`:
   - Parse comma-separated tokens from `*flagPtrs.forceIncludeFlag` into `opts.ForceIncludeDirs`.
3. In `cli/cmd/rootadd.go`:
   - Inspect flags for `--force` / `-f`.
   - If user attempts to add an excluded directory (e.g. `gitmap add ~/.oh-my-zsh`), display an informational warning unless `--force` is present.

---

## 3. CLI Usage & Examples

```bash
# 1. Standard scan ignores Oh-My-Zsh and tool dependencies:
gitmap scan ~

# 2. Scanning with explicit force inclusion:
gitmap scan ~ --force-include .oh-my-zsh

# 3. Forcing manual registration of an excluded directory:
gitmap add ~/.oh-my-zsh --force
```

---

## 4. Verification & Quality Gates

1. Run relative path hygiene:
   `python linter-scripts/check-relative-paths.py`
2. Run coding guideline autofixer check:
   `python 03-ai-scripts/05-guideline-autofixer.py cli/constants cli/scanner cli/fsutil cli/cmdscan cli/cmd --check-only`
3. Verify unit test expectations:
   - Scanning a directory tree with `.oh-my-zsh` ignores it by default.
   - Supplying `--force-include .oh-my-zsh` discovers repositories inside `.oh-my-zsh`.
