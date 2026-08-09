# CI/CD Pipeline & Release Skew (Extend V2)

This document is a mandatory read for any AI agents orchestrating releases, interacting with GitHub Actions, or modifying the continuous integration pipeline for this repository. It serves to outline exact procedures to prevent "release skew", where Git tags become out of sync with the codebase version.

## 1. Root Cause Analysis: The `v6.87.1` Incident
In the past, an agent pushed a new Git tag (`v6.87.1`) to trigger a release, but failed to bump the internal source code references first. The canonical file `constants.go` remained hardcoded at `6.87.0`. 

**What Happened:**
1. **GitHub Actions Triggered:** Pushing the tag `v6.87.1` successfully triggered the `.github/workflows/release.yml` pipeline.
2. **Missing Notes:** The release workflow uses an `awk` script to extract release notes from `changelog.md` by searching for `[v6.87.1]`. Since the file was never updated, the script fell back to outputting: `"No changelog entry found"`.
3. **Smoke Test Failure:** The workflow builds the Go binaries using the hardcoded values in the repo. It then runs an `Installer Smoke Test` using PowerShell. The script executes the built binary expecting it to output `v6.87.1` (matching the tag). Instead, the binary output `v6.87.0`. The script forcefully threw `Version mismatch. expected: v6.87.1, actual: v6.87.0` and killed the CI run.
4. **The Result:** A broken, partially complete GitHub release devoid of release notes.

**The Fix / Core Rule:**
**NEVER push a tag to trigger a release before securely pushing the modified canonical version files to the `main` branch first.**

## 2. Release Ceremony Blueprint (Agent Instructions)
To correctly execute a release, an AI agent MUST follow this exact synchronous sequence.

### Step 1: PowerShell Version Sweeping
You must bump the version synchronously across all canonical files. Do not use manual multi-edits if PowerShell can do it cleaner.
```powershell
$newVer = "X.Y.Z"
$files = @(
  ".gitmap/release/latest.json", 
  ".lovable/plan.md", 
  "gitmap/constants/constants.go", 
  "readme.md", 
  "src/constants/index.ts", 
  "package.json"
)

foreach ($f in $files) { 
  $c = Get-Content $f -Raw
  $c = $c -replace 'OLD\.VER\.SION', $newVer
  Set-Content $f -Value $c -NoNewline 
}
```
*Note: Always handle `.gitmap/release/vOLD.json` by copying it to `vNEW.json` and replacing its internals.*

### Step 2: Changelog Updates
Update `changelog.md` by injecting the `## [vX.Y.Z] YYYY-MM-DD` block at the top. The `release.yml` uses an `awk` command to slice exactly between this header and the next `## [v` header to generate the GitHub Release Notes.

### Step 3: Git Synchronization
Group these changes into a commit and push to `main` **FIRST**.
```bash
git add .
git commit -m "chore(release): bump version to vX.Y.Z"
git push origin main
```
*Wait a few seconds to ensure GitHub registers the commit on `main`.*

### Step 4: Tag Execution
Only now can you safely push the tag.
```bash
git tag vX.Y.Z
git push origin vX.Y.Z
```

## 3. Workflow Integration Insights
The `.github/workflows/release.yml` file is heavily guarded. 

- **Changelog Extraction**:
  ```bash
  changelog=$(awk -v tag="[${GITHUB_REF_NAME}]" '
    $0 ~ "^## " tag {flag=1; next}
    /^## \[v/ && flag {flag=0; exit}
    flag {print}
  ' changelog.md)
  ```
  If `changelog.md` doesn't contain the exact tag name natively, this silent fails.

- **Installer Smoke Test**:
  The workflow compiles the binary, places it in a `$tmpDir`, and runs it:
  ```powershell
  $actualVersion = & $binary version
  if ($actualVersion -notmatch $expectedVersion) {
      Write-Error "Version mismatch"
  }
  ```
  If `constants.go` wasn't updated in Step 1, this step will instantly crash the entire CI/CD release chain.

## 4. Conclusion
AI Agents modifying this repository must treat versioning as a rigid, multi-file transaction. **Source control dictates the tag, the tag does not dictate source control.**
