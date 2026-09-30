# Subtask 02: CFR and Clone Manifest Auto-Detection and Missing-Only Execution

> **Parent Plan:** `200-nodes-display-cfr-manifest-gitmap-json-and-vscode-optimize.md`  
> **Status:** `PENDING`  
> **Target Files:**
> - `cli/cmdclone/clonefixrepo.go`
> - `cli/cmdclone/clone.go`
> - `cli/cmdclone/clone_only_missing.go`

---

## Technical Specification

1. **`gitmap cfr` Manifest Auto-Detection**:
   - In `runCloneFixRepo(args)` / `runCloneFixRepoPipeline(args)`:
     - Case A: `len(args) == 0` (or only flags without URL positional):
       - If `gitmap.json` (or `.gitmap/gitmap.json`) exists in CWD:
         - Route to `RunCloneOnlyMissingFromJSON("gitmap.json", ...)` or `runCloneOnlyMissingPipeline`!
         - Inform user:
           `ℹ No target URL specified; discovered local manifest "gitmap.json".`
           `↪ Cloning missing repositories only...`
     - Case B: The first positional argument is a `.json` file (`strings.HasSuffix(f.url, ".json")` or `isRegularFile(f.url)`):
       - Do NOT pass it to `ResolveRepoSlug`!
       - Route directly to manifest missing-only clone pipeline!
       - Each missing repository that is cloned will have its fix-repo post-step run.

2. **`gitmap clone` Optimization for Existing Repos**:
   - In `cli/cmdclone/clone.go`:
     - When `cf.Source` is a manifest file (`gitmap.json`), and `--missing-only` is set OR default clone is run:
       - Check which records already exist on disk before resolving SSH/HTTPS auth:
         - If directory `filepath.Join(cf.TargetDir, rec.RelativePath)` exists and has `.git`, mark as existing.
         - Skip running network SSH/token probe (`ResolveRepoAuth`) for already-existing repos!
     - If all repositories in `cf.Source` already exist on disk, do not exit silently. Emit:
       `✓ All %d repository(ies) in %s already exist on disk. Nothing to clone.`
     - When run with zero arguments and `gitmap.json` exists in CWD, default to cloning from `gitmap.json`.

3. **`gitmap clone-only-missing`**:
   - In `cli/cmdclone/clone_only_missing.go`:
     - If no args provided, auto-discover `gitmap.json` and clone missing only.
